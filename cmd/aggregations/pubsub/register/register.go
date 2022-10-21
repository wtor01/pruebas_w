package register

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	services "bitbucket.org/sercide/data-ingestion/internal/aggregations/services"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	inventory_services "bitbucket.org/sercide/data-ingestion/internal/inventory/services"
	mongo_r "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
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
	serviceName := "aggregations"
	if os.Getenv("SERVICE") != "" {
		serviceName = os.Getenv("SERVICE")
	}

	client, err := gcp_pubsub.NewClient(ctx, cnf.ProjectID)

	if err != nil {
		return err
	}
	sub := client.Subscription(cnf.SubscriptionAggregations)
	if sub == nil {
		return errors.New(fmt.Sprintf("invalid subscription %s", cnf.SubscriptionAggregations))
	}

	inventoryRepository := postgres.NewInventoryPostgres(db, redis_repos.NewDataCacheRedis(rd))
	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	inventoryRepositoryMongo := mongo_r.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)
	aggregationRepositoryMongo := mongo_r.NewAggregateRepositoryMongo(mongoClient, cnf.MeasureDB)

	srv := services.NewPubsub(
		publisher,
		cnf.TopicAggregations,
		inventory.NewInventory(inventory_services.New(inventoryRepository, inventoryRepositoryMongo), redis_repos.NewDataCacheRedis(rd)),
		aggregationRepositoryMongo,
	)
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 10

	logger := log.New(serviceName)
	defer logger.Sync()

	tp := telemetry.NewTracerProvider(serviceName)

	defer tp.ForceFlush(ctx)
	tracer := telemetry.SetTracer(serviceName + "/pubsub")

	err = sub.Receive(ctx, func(ctx context.Context, message *gcp_pubsub.Message) {

		eventType := message.Attributes[event.EventTypeKey]

		ctx, span := tracer.Start(ctx, eventType)

		defer span.End()

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
		case aggregations.SchedulerEventType:
			{
				ctx, _ = tracer.Start(ctx, "ProcessAggregationInitService")

				var ev aggregations.SchedulerEvent
				errHandler = json.Unmarshal(message.Data, &ev)

				if errHandler != nil {
					break
				}

				today := time.Now().In(cnf.LocalLocation)
				dateTime := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, cnf.LocalLocation)

				if ev.Payload.Date.IsZero() {
					ev.Payload.Date = dateTime
				}

				errHandler = srv.ProcessAggregationInitService.Handler(ctx, ev.Payload)

				break
			}
		case aggregations.SchedulerDistributorEventType:
			{
				ctx, _ = tracer.Start(ctx, "ProcessAggregationByDistributorService")

				var ev aggregations.SchedulerEvent
				errHandler = json.Unmarshal(message.Data, &ev)

				if errHandler != nil {
					break
				}

				errHandler = srv.ProcessAggregationByDistributorService.Handler(ctx, ev.Payload)

				break

			}
		default:
			break
		}

	})

	return err

}
