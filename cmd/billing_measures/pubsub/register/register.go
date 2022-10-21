package register

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/services"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	inventory_services "bitbucket.org/sercide/data-ingestion/internal/inventory/services"
	tariff_services "bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff/services"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	mongo_r "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/apperrors"
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

const EVENT_KEY_LOGGER = "event"
const EVENT_KEY_ID = "event_id"

func Register(ctx context.Context, cnf config.Config, mongoClient *mongo.Client, db *gorm.DB, rd *redis.Client) error {
	serviceName := "billing_measures"
	if os.Getenv("SERVICE") != "" {
		serviceName = os.Getenv("SERVICE")
	}

	client, err := gcp_pubsub.NewClient(ctx, cnf.ProjectID)

	if err != nil {
		return err
	}
	sub := client.Subscription(cnf.SubscriptionBillingMeasures)
	if sub == nil {
		return errors.New(fmt.Sprintf("invalid subscription %s", cnf.ProcessedMeasuresSub))
	}

	inventoryRepository := postgres.NewInventoryPostgres(db, redis_repos.NewDataCacheRedis(rd))
	schedulerRepository := postgres.NewProcessBillingSchedulerPostgres(db)
	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	billingMeasuresRepository := mongo_r.NewBillingMeasureRepositoryMongo(mongoClient, cnf.MeasureDB, "billing_measure")
	repositoryFestiveDays := redis_repos.NewProcessMeasureFestiveDays(rd)
	inventoryRepositoryMongo := mongo_r.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)
	selfConsumptionRepositoryRepositoryMongo := mongo_r.NewSelfConsumptionRepository(mongoClient, cnf.MeasureDB, cnf.LocalLocation)
	billingSelfConsumptionRepositoryRepositoryMongo := mongo_r.NewBillingSelfConsumptionRepository(mongoClient, cnf.MeasureDB, cnf.LocalLocation)
	tariffRepository := postgres.NewTariffPostgres(db, redis_repos.NewDataCacheRedis(rd))

	processRepository := mongo_r.NewProcessMeasureRepository(
		mongoClient,
		cnf.MeasureDB,
		cnf.CollectionLoadCurve,
		cnf.CollectionDailyLoadCurve,
		cnf.CollectionMonthlyClosure,
		cnf.CollectionDailyClosure,
		cnf.LocalLocation,
	)

	masterTablesClient := master_tables.NewMasterTables(tariff_services.NewTariffServices(tariffRepository))

	srv := services.NewServicesPubsub(
		publisher,
		cnf.TopicBillingMeasures,
		schedulerRepository,
		billingMeasuresRepository,
		selfConsumptionRepositoryRepositoryMongo,
		billingSelfConsumptionRepositoryRepositoryMongo,
		processRepository,
		inventory.NewInventory(inventory_services.New(inventoryRepository, inventoryRepositoryMongo), redis_repos.NewDataCacheRedis(rd)),
		inventoryRepositoryMongo,
		mongo_r.NewConsumProfileRepository(mongoClient, cnf.MeasureDB),
		repositoryFestiveDays,
		cnf.LocalLocation,
		masterTablesClient,
		mongo_r.ConsumCoefficientRepositoryMongo{},
	)
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 10

	logger := log.New(serviceName)
	defer logger.Sync()

	tp := telemetry.NewTracerProvider(serviceName)

	defer tp.ForceFlush(ctx)
	tracer := telemetry.SetTracer("billing_measures/" + serviceName + "/pubsub")

	err = sub.Receive(ctx, func(ctx context.Context, message *gcp_pubsub.Message) {

		eventType := message.Attributes[event.EventTypeKey]

		ctx, span := tracer.Start(ctx, eventType)

		defer span.End()

		var errHandler error

		defer func() {
			defer func() {
				span.End()
			}()
			l := logger.Ctx(ctx)
			if rec := recover(); rec != nil {
				message.Nack()
				logger.Ctx(ctx).Errorf("panic", rec)
				span.RecordError(errHandler)
				span.SetStatus(codes.Error, fmt.Sprintf("%v", rec))
				return
			}
			if errHandler != nil {
				l.Errorf(errHandler.Error())
				span.RecordError(errHandler)
				span.SetStatus(codes.Error, errHandler.Error())
				var appError apperrors.AppError
				if errors.As(errHandler, &appError) && !appError.ShouldRetry() {
					message.Ack()
					return
				}
				message.Nack()
				return
			}
			message.Ack()
			l.Infof("success %s", eventType)
			span.SetStatus(codes.Ok, fmt.Sprintf("success %s", eventType))
		}()

		span.SetAttributes(attribute.String("event_data", string(message.Data)))

		switch eventType {
		case billing_measures.SchedulerEventType:
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
		case billing_measures.TypeDistributorProcess:
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

		case billing_measures.ProcessMVH:
			{
				ctx, _ = tracer.Start(ctx, "ProcessMVH")

				var ev measures.ProcessMeasureEvent

				errHandler = json.Unmarshal(message.Data, &ev)

				if errHandler != nil {
					break
				}

				errHandler = srv.ProcessDcMeasureByCup.Handler(ctx, ev.Payload)
			}
		case billing_measures.SelfConsumptionEventType:
			{
				ctx, _ = tracer.Start(ctx, "ProcessMVH")

				var ev billing_measures.OnSaveBillingMeasureEvent

				errHandler = json.Unmarshal(message.Data, &ev)

				if errHandler != nil {
					break
				}

				errHandler = srv.ProcessSelfConsumption.Handler(ctx, ev.Payload)
			}
		default:
			break
		}

	})

	return err

}
