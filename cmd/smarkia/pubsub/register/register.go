package register

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia/services"
	inventory_services "bitbucket.org/sercide/data-ingestion/internal/inventory/services"
	mongo_repos "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/internal/platform/smarkia_api"
	"bitbucket.org/sercide/data-ingestion/pkg/log"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	gcp_pubsub "cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
	"os"
	"time"
)

func Register(ctx context.Context, db *gorm.DB, cnf config.Config, redis *redis.Client, mongoClient *mongo.Client) error {
	serviceName := "smarkia"
	if os.Getenv("SERVICE") != "" {
		serviceName = os.Getenv("SERVICE")
	}

	client, err := gcp_pubsub.NewClient(ctx, cnf.ProjectID)
	inventoryRepository := postgres.NewInventoryPostgres(db, redis_repos.NewDataCacheRedis(redis))
	inventoryMeasuresRepository := mongo_repos.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)
	if err != nil {
		return err
	}

	sub := client.Subscription(cnf.SmarkiaSub)

	if sub == nil {
		return errors.New(fmt.Sprintf("invalid subscription %s", cnf.SubscriptionMeasures))
	}

	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)

	smarkiaApi := smarkia_api.NewApi(cnf.SmarkiaToken, cnf.SmarkiaHost, cnf.LocalLocation)

	srv := services.NewServices(
		publisher,
		inventory.NewInventory(inventory_services.New(inventoryRepository, inventoryMeasuresRepository), redis_repos.NewDataCacheRedis(redis)),
		smarkiaApi,
		cnf,
		inventoryMeasuresRepository,
	)

	logger := log.New(serviceName)
	defer logger.Sync()

	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 100
	tp := telemetry.NewTracerProvider(serviceName)

	defer tp.ForceFlush(ctx)

	tracer := telemetry.SetTracer("gross_measures/" + serviceName + "/pubsub")

	err = sub.Receive(ctx, func(ctx context.Context, msg *gcp_pubsub.Message) {
		eventType := msg.Attributes["type"]

		ctx, span := tracer.Start(ctx, eventType)

		defer span.End()

		defer func() {
			l := logger.Ctx(ctx)
			if rec := recover(); rec != nil {
				msg.Nack()
				logger.Ctx(ctx).Errorf("panic", rec)
				span.RecordError(err)
				span.SetStatus(codes.Error, fmt.Sprintf("%v", rec))
			}
			if err != nil {
				msg.Nack()
				l.Errorf(err.Error())
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				msg.Ack()
				l.Infof("success %s", eventType)
				span.SetStatus(codes.Ok, string(msg.Data))
			}
		}()

		switch eventType {
		case smarkia.TypeRequestSmarkia:
			{
				var event smarkia.RequestSmarkiaEvent
				err = json.Unmarshal(msg.Data, &event)

				if err != nil {
					return
				}

				date := event.Payload.Date

				loc, _ := time.LoadLocation(cnf.TimeZone)

				if date.IsZero() {
					date = time.Now().In(loc)
				}

				date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)

				date.Format(time.RFC3339)
				ctx, _ = tracer.Start(ctx, "PublishDistributorService")

				err = srv.PublishDistributorService.Handle(ctx, services.PublishDistributorDTO{
					ProcessName: event.Payload.ProcessName,
					Date:        date,
				})
				if err != nil {
					return
				}
			}
		case smarkia.TypeDistributorProcess:
			{
				var event smarkia.MessageDistributorProcess

				err = json.Unmarshal(msg.Data, &event)
				if err != nil {
					return
				}
				ctx, _ = tracer.Start(ctx, "PublishCtService")

				err = srv.PublishCtService.Handle(ctx, services.PublishCtDTO{
					ProcessName:     event.Payload.ProcessName,
					DistributorId:   event.Payload.DistributorId,
					DistributorCDOS: event.Payload.DistributorCDOS,
					SmarkiaId:       event.Payload.SmarkiaId,
					Date:            event.Payload.Date,
				})
				if err != nil {
					return
				}
			}
		case smarkia.TypeCtProcess:
			{
				var event smarkia.CtProcessEvent
				err = json.Unmarshal(msg.Data, &event)
				if err != nil {
					return
				}
				ctx, _ = tracer.Start(ctx, "PublishEquipmentService")

				err = srv.PublishEquipmentService.Handle(ctx, services.PublishEquipmentDto{
					ProcessName:     event.Payload.ProcessName,
					DistributorId:   event.Payload.DistributorId,
					SmarkiaId:       event.Payload.SmarkiaId,
					CtId:            event.Payload.CtId,
					DistributorCDOS: event.Payload.DistributorCDOS,
					Date:            event.Payload.Date,
				})

				if err != nil {
					return
				}
			}
		case smarkia.TypeEquipmentProcess:
			{
				var event smarkia.EquipmentProcessEvent
				err = json.Unmarshal(msg.Data, &event)
				if err != nil {
					return
				}
				ctx, _ = tracer.Start(ctx, "PublishEquipmentService")

				err = srv.PublishMeasuresFromSmarkiaService.Handle(ctx, event)
				if err != nil {
					return
				}
			}

		}
	})

	return err

}
