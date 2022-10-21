package register

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	inventory_services "bitbucket.org/sercide/data-ingestion/internal/inventory/services"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	mongo_repos "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services"
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
)

func Register(ctx context.Context, cnf config.Config, mongoClient *mongo.Client, db *gorm.DB, rd *redis.Client) error {
	serviceName := "re_process_measures"
	if os.Getenv("SERVICE") != "" {
		serviceName = os.Getenv("SERVICE")
	}

	client, err := gcp_pubsub.NewClient(ctx, cnf.ProjectID)

	if err != nil {
		return err
	}
	sub := client.Subscription(cnf.ReProcessedMeasuresSub)
	if sub == nil {
		return errors.New(fmt.Sprintf("invalid subscription %s", cnf.ReProcessedMeasuresSub))
	}

	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)

	inventoryRepository := postgres.NewInventoryPostgres(db, redis_repos.NewDataCacheRedis(rd))
	schedulerRepository := postgres.NewProcessMeasureSchedulerPostgres(db)

	grossRepository := mongo_repos.NewGrossMeasureRepositoryMongo(mongoClient, cnf.MeasureDB, cnf.MeasureCollection, cnf.LocalLocation)
	inventoryMeasuresRepository := mongo_repos.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)

	reprocessingDateRepository := redis_repos.NewDateRedis(redis_repos.NewDataCacheRedis(rd))

	srv := services.NewReprocessingServicesPubSub(
		publisher,
		cnf.ProcessedMeasureTopic,
		schedulerRepository,
		inventory.NewInventory(inventory_services.New(inventoryRepository, inventoryMeasuresRepository), redis_repos.NewDataCacheRedis(rd)),
		inventoryMeasuresRepository,
		grossRepository,
		reprocessingDateRepository,
	)
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 10

	logger := log.New(serviceName)
	defer logger.Sync()

	tp := telemetry.NewTracerProvider(serviceName)

	defer tp.ForceFlush(ctx)
	tracer := telemetry.SetTracer("re_process_measures/" + serviceName + "/pubsub")

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
		case process_measures.SchedulerReprocessingEventType:
			{
				ctx, _ = tracer.Start(ctx, "ReprocessingDistributorPointService")

				var ev measures.SchedulerEvent

				errHandler = json.Unmarshal(message.Data, &ev)
				if errHandler != nil {
					break
				}
				errHandler = srv.ReprocessingPublishDistributorService.Handle(ctx, ev.Payload)

				break
			}
		case process_measures.TypeReprocessingMeterProcess:
			{
				ctx, _ = tracer.Start(ctx, "ReprocessingMeterService")
				var ev process_measures.ReSchedulerMeterEvent
				errHandler = json.Unmarshal(message.Data, &ev)
				if errHandler != nil {
					break
				}
				errHandler = srv.PublishReprocessingMeterService.Handle(ctx, ev.Payload)

				break
			}

		case process_measures.TypeReprocessingDistributorProcess:
			{
				ctx, _ = tracer.Start(ctx, "PublishReprocessingServicePointService")

				var ev process_measures.ReSchedulerEvent

				errHandler = json.Unmarshal(message.Data, &ev)
				if errHandler != nil {
					break
				}
				errHandler = srv.PublishReprocessingServicePointService.Handle(ctx, ev.Payload)

				break
			}
		default:
			break
		}

	})

	return err

}
