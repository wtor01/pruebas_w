package register

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/services"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/log"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	gcp_pubsub "cloud.google.com/go/pubsub"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
	"os"
)

func Register(ctx context.Context, cnf config.Config, db *gorm.DB, rd *redis.Client) error {
	serviceName := "calendar_periods"
	if os.Getenv("SERVICE") != "" {
		serviceName = os.Getenv("SERVICE")
	}

	client, err := gcp_pubsub.NewClient(ctx, cnf.ProjectID)
	subscription := cnf.SubscriptionCalendarPeriods

	if err != nil {
		return err
	}
	sub := client.Subscription(subscription)
	if sub == nil {
		return errors.New(fmt.Sprintf("invalid subscription %s", subscription))
	}

	calendarRepository := postgres.NewCalendarPostgres(db, redis_repos.NewDataCacheRedis(rd))
	calendarPeriodRepository := redis_repos.NewProcessMeasureFestiveDays(rd)
	festiveDaysRepository := postgres.NewFestiveDaysPostgres(db, redis_repos.NewDataCacheRedis(rd))
	srv := services.NewCalendarPeriodsPubSubServices(calendarRepository, calendarPeriodRepository, festiveDaysRepository)

	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 10

	logger := log.New(serviceName)
	defer logger.Sync()

	tp := telemetry.NewTracerProvider(serviceName)

	defer tp.ForceFlush(ctx)
	tracer := telemetry.SetTracer("master_tables/" + serviceName + "/pubsub")

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
		case master_tables.CalendarPeriodGenerate:
			{
				ctx, _ = tracer.Start(ctx, "GenerateCalendarPeriods")

				errHandler = srv.CalendarPeriodGenerateService.Handler(ctx)

				break
			}
		case master_tables.FestiveDaysGenerate:
			ctx, _ = tracer.Start(ctx, "GenerateFestiveDays")

			errHandler = srv.FestiveDaysGenerateService.Handler(ctx)
		default:
			break
		}

	})

	return err
}
