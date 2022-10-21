package register

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/services"
	inventory_services "bitbucket.org/sercide/data-ingestion/internal/inventory/services"
	mongo_repos "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	validation_services "bitbucket.org/sercide/data-ingestion/internal/validations/services"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/log"
	"bitbucket.org/sercide/data-ingestion/pkg/storage"
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

type StorageEvent struct {
	Kind                    string    `json:"kind"`
	Id                      string    `json:"id"`
	SelfLink                string    `json:"selfLink"`
	Name                    string    `json:"name"`
	Bucket                  string    `json:"bucket"`
	Generation              string    `json:"generation"`
	Metageneration          string    `json:"metageneration"`
	ContentType             string    `json:"contentType"`
	TimeCreated             time.Time `json:"timeCreated"`
	Updated                 time.Time `json:"updated"`
	StorageClass            string    `json:"storageClass"`
	TimeStorageClassUpdated time.Time `json:"timeStorageClassUpdated"`
	Size                    string    `json:"size"`
	Md5Hash                 string    `json:"md5Hash"`
	MediaLink               string    `json:"mediaLink"`
	Crc32C                  string    `json:"crc32c"`
	Etag                    string    `json:"etag"`
}

const OBJECT_FINALIZE_EVENT_NAME = "OBJECT_FINALIZE"

func Register(ctx context.Context, db *gorm.DB, cnf config.Config, mongoClient *mongo.Client, redis *redis.Client) error {
	serviceName := "gross_measures"
	if os.Getenv("SERVICE") != "" {
		serviceName = os.Getenv("SERVICE")
	}

	inventoryRepository := postgres.NewInventoryPostgres(db, redis_repos.NewDataCacheRedis(redis))
	adminRepository := postgres.NewValidationMeasurePostgresRepository(db, redis_repos.NewDataCacheRedis(redis), cnf.ShutdownTimeout)
	measureRepository := mongo_repos.NewGrossMeasureRepositoryMongo(mongoClient, cnf.MeasureDB, cnf.MeasureCollection, cnf.LocalLocation)
	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	client, err := gcp_pubsub.NewClient(ctx, cnf.ProjectID)
	inventoryRepositoryMongo := mongo_repos.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)
	processRepository := mongo_repos.NewProcessMeasureRepository(
		mongoClient,
		cnf.MeasureDB,
		cnf.CollectionLoadCurve,
		cnf.CollectionDailyLoadCurve,
		cnf.CollectionMonthlyClosure,
		cnf.CollectionDailyClosure,
		cnf.LocalLocation,
	)

	if err != nil {
		return err
	}

	sub := client.Subscription(os.Getenv("SUB_READ"))

	if sub == nil {
		return errors.New(fmt.Sprintf("invalid subscription %s", os.Getenv("SUB_READ")))
	}

	svcs := services.NewServices(
		measureRepository,
		publisher,
		storage.NewStorageGCP,
		cnf,
		inventory.NewInventory(inventory_services.New(inventoryRepository, inventoryRepositoryMongo), redis_repos.NewDataCacheRedis(redis)),
		internal_clients.NewValidation(validation_services.NewServices(adminRepository, processRepository)),
	)

	logger := log.New(serviceName)
	defer logger.Sync()

	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 10

	tp := telemetry.NewTracerProvider(serviceName)

	defer tp.ForceFlush(ctx)
	tracer := telemetry.SetTracer("gross_measures/" + serviceName + "/pubsub")

	err = sub.Receive(ctx, func(ctx context.Context, msg *gcp_pubsub.Message) {

		eventType := msg.Attributes[event.EventTypeKey]

		if eventType == "" {
			eventType = msg.Attributes["eventType"]
		}

		ctx, span := tracer.Start(ctx, eventType)

		defer span.End()

		var errHandler error

		defer func() {
			l := logger.Ctx(ctx)
			if rec := recover(); rec != nil {
				msg.Nack()
				logger.Ctx(ctx).Errorf("panic", rec)
				span.RecordError(errHandler)
				span.SetStatus(codes.Error, fmt.Sprintf("%v", rec))
			}
			if errHandler != nil {
				msg.Nack()
				l.Errorf(errHandler.Error())
				span.RecordError(errHandler)
				span.SetStatus(codes.Error, errHandler.Error())
			} else {
				msg.Ack()
				l.Infof("success %s", eventType)
				span.SetStatus(codes.Ok, fmt.Sprintf("success %s", eventType))
			}
			span.End()
		}()

		span.SetAttributes(attribute.String("event_data", string(msg.Data)))

		switch eventType {
		case OBJECT_FINALIZE_EVENT_NAME:
			{
				var ev StorageEvent

				errHandler = json.Unmarshal(msg.Data, &ev)

				if errHandler != nil {
					return
				}

				ctx, _ = tracer.Start(ctx, "ParseFilesService")

				errHandler = svcs.ParseFilesService.Handle(ctx, gross_measures.HandleFileDTO{
					FilePath: ev.Bucket + "/" + ev.Name,
				})
				break
			}
		case gross_measures.InsertMeasureCurveEventType:
			{
				var ev gross_measures.InsertMeasureCurveEvent

				errHandler = json.Unmarshal(msg.Data, &ev)

				if errHandler != nil {
					return
				}
				ctx, _ = tracer.Start(ctx, "InsertMeasureCurve")

				errHandler = svcs.InsertMeasureCurve.Handle(ctx, ev.Payload)

				if errHandler != nil {
					return
				}
				break
			}
		case gross_measures.InsertMeasureCloseEventType:
			{
				var ev gross_measures.InsertMeasureCloseEvent

				errHandler = json.Unmarshal(msg.Data, &ev)

				if errHandler != nil {
					return
				}

				ctx, _ = tracer.Start(ctx, "InsertMeasureClose")

				errHandler = svcs.InsertMeasureClose.Handle(ctx, ev.Payload)

				if errHandler != nil {
					return
				}
				break
			}
		default:

			break
		}
	})

	return err

}
