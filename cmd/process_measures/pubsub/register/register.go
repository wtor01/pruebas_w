package register

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	inventory_services "bitbucket.org/sercide/data-ingestion/internal/inventory/services"
	tariff_services "bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff/services"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	mongo_repos "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services"
	validation_services "bitbucket.org/sercide/data-ingestion/internal/validations/services"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/log"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	gcp_pubsub "cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
	"os"
	"time"
)

func Register(ctx context.Context, cnf config.Config, mongoClient *mongo.Client, db *gorm.DB, rd *redis.Client) error {
	serviceName := "process_measures"
	if os.Getenv("SERVICE") != "" {
		serviceName = os.Getenv("SERVICE")
	}

	client, err := gcp_pubsub.NewClient(ctx, cnf.ProjectID)

	if err != nil {
		return err
	}
	sub := client.Subscription(cnf.ProcessedMeasuresSub)
	if sub == nil {
		return errors.New(fmt.Sprintf("invalid subscription %s", cnf.ProcessedMeasuresSub))
	}

	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)

	inventoryRepository := postgres.NewInventoryPostgres(db, redis_repos.NewDataCacheRedis(rd))
	schedulerRepository := postgres.NewProcessMeasureSchedulerPostgres(db)
	adminRepository := postgres.NewValidationMeasurePostgresRepository(db, redis_repos.NewDataCacheRedis(rd), cnf.ShutdownTimeout)
	calendarRepository := postgres.NewCalendarPostgres(db, redis_repos.NewDataCacheRedis(rd))
	tariffRepository := postgres.NewTariffPostgres(db, redis_repos.NewDataCacheRedis(rd))
	seasonRepository := postgres.NewSeasonPostgres(db, redis_repos.NewDataCacheRedis(rd))

	processRepository := mongo_repos.NewProcessMeasureRepository(
		mongoClient,
		cnf.MeasureDB,
		cnf.CollectionLoadCurve,
		cnf.CollectionDailyLoadCurve,
		cnf.CollectionMonthlyClosure,
		cnf.CollectionDailyClosure,
		cnf.LocalLocation,
	)
	validationRepository := mongo_repos.NewValidationMeasureRepository(
		mongoClient,
		cnf.MeasureDB,
		cnf.CollectionLoadCurve,
		cnf.CollectionDailyLoadCurve,
		cnf.CollectionMonthlyClosure,
		cnf.CollectionDailyClosure,
		cnf.LocalLocation,
	)

	grossRepository := mongo_repos.NewGrossMeasureRepositoryMongo(mongoClient, cnf.MeasureDB, cnf.MeasureCollection, cnf.LocalLocation)
	inventoryMeasuresRepository := mongo_repos.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)

	repositoryFestiveDays := redis_repos.NewProcessMeasureFestiveDays(rd)

	masterTablesClient := master_tables.NewMasterTables(tariff_services.NewTariffServices(tariffRepository))

	srv := services.NewServicesPubsub(
		publisher,
		cnf.ProcessedMeasureTopic,
		cnf.TopicBillingMeasures,
		schedulerRepository,
		inventory.NewInventory(inventory_services.New(inventoryRepository, inventoryMeasuresRepository), redis_repos.NewDataCacheRedis(rd)),
		inventoryMeasuresRepository,
		grossRepository,
		processRepository,
		repositoryFestiveDays,
		seasonRepository,
		cnf.LocalLocation,
		internal_clients.NewValidation(validation_services.NewServices(adminRepository, processRepository)),
		calendarRepository,
		masterTablesClient,
		validationRepository,
	)
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 10

	logger := log.New(serviceName)
	defer logger.Sync()

	tp := telemetry.NewTracerProvider(serviceName)

	defer tp.ForceFlush(ctx)
	tracer := telemetry.SetTracer("gross_measures/" + serviceName + "/pubsub")

	err = sub.Receive(ctx, func(ctx context.Context, message *gcp_pubsub.Message) {

		eventType := message.Attributes[event.EventTypeKey]

		ctx, span := tracer.Start(ctx, eventType)

		var errHandler error

		defer func() {
			l := logger.Ctx(ctx)
			if rec := recover(); rec != nil {
				message.Nack()
				logger.Ctx(ctx).Errorf("panic", rec)
				span.RecordError(errHandler)
				span.SetStatus(codes.Error, fmt.Sprintf("%v", rec))
			}
			if errHandler != nil {
				message.Nack()
				l.Errorf(errHandler.Error())
				span.RecordError(errHandler)
				span.SetStatus(codes.Error, errHandler.Error())
			} else {
				message.Ack()
				l.Infof("success %s", eventType)
				span.SetStatus(codes.Ok, fmt.Sprintf("success %s", eventType))
			}
			span.End()
		}()

		span.SetAttributes(attribute.String("event_data", string(message.Data)))

		switch eventType {
		case process_measures.SchedulerEventType:
			{
				ctx, _ = tracer.Start(ctx, "PublishDistributorService")

				var ev measures.SchedulerEvent
				errHandler = json.Unmarshal(message.Data, &ev)

				if errHandler != nil {
					break
				}

				today := time.Now().In(cnf.LocalLocation)
				dateTime := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, cnf.LocalLocation).AddDate(0, 0, -1)

				if ev.Payload.Date.IsZero() {
					ev.Payload.Date = dateTime
				}

				errHandler = srv.PublishDistributorService.Handle(ctx, ev.Payload)

				break
			}
		case process_measures.TypeDistributorProcess:
			{
				ctx, _ = tracer.Start(ctx, "PublishServicePointService")

				var ev measures.SchedulerEvent

				errHandler = json.Unmarshal(message.Data, &ev)

				if errHandler != nil {
					break
				}

				errHandler = srv.PublishServicePointService.Handle(ctx, ev.Payload)

				break
			}
			//CURVA
		case process_measures.ProcessCurveEventType:
			{
				ctx, _ = tracer.Start(ctx, "ProcessCurve")

				var ev measures.ProcessMeasureEvent
				errHandler = json.Unmarshal(message.Data, &ev)

				if errHandler != nil {
					break
				}

				errHandler = srv.ProcessCurve.Handle(ctx, ev.Payload)

				break
			}
			//CIERRE DIARIO
		case process_measures.ProcessDailyClosureEventType:
			{
				ctx, _ = tracer.Start(ctx, "ProcessDailyClosure")

				var ev measures.ProcessMeasureEvent
				errHandler = json.Unmarshal(message.Data, &ev)
				if errHandler != nil {
					break
				}
				errHandler = srv.ProcessDailyClosure.Handle(ctx, ev.Payload)
				break
			}
			//CIERRE MENSUAL
		case process_measures.ProcessBillingClosureEventType:
			{
				ctx, _ = tracer.Start(ctx, "ProcessMonthlyClosure")

				var ev measures.ProcessMeasureEvent
				errHandler = json.Unmarshal(message.Data, &ev)

				if errHandler != nil {
					break
				}

				errHandler = srv.ProcessMonthlyClosure.Handle(ctx, ev.Payload)

				break
			}
		default:
			break
		}

	})

	return err

}
