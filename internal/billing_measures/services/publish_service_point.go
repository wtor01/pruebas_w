package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	measures_services "bitbucket.org/sercide/data-ingestion/internal/measures/services"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"go.opentelemetry.io/otel/trace"
)

type PublishServicePointService struct {
	*measures_services.OnSchedulerDistributorEvent
	tracer trace.Tracer
}

func NewServicePointService(
	publisher event.PublisherCreator,
	inventoryClient measures.InventoryRepository,
	topic string,
) *PublishServicePointService {
	return &PublishServicePointService{
		OnSchedulerDistributorEvent: measures_services.NewOnSchedulerDistributorEvent(publisher, inventoryClient, topic, 100),
		tracer:                      telemetry.GetTracer(),
	}
}

func (svc PublishServicePointService) Handle(ctx context.Context, dto measures.SchedulerEventPayload) error {
	ctx, span := svc.tracer.Start(ctx, "PublishServicePointService - Handle")
	defer span.End()

	err := svc.OnSchedulerDistributorEvent.Handle(
		ctx,
		dto,
		billing_measures.NewProcessByDistributorEvent,
		func(dto measures.SchedulerEventPayload, meterConfig measures.MeterConfig) []measures.ProcessMeasureEvent {

			if !utils.InSlice(dto.ServiceType, []measures.ServiceType{measures.DcServiceType, measures.GdServiceType}) {
				return []measures.ProcessMeasureEvent{}
			}

			return []measures.ProcessMeasureEvent{billing_measures.NewProcessMvhEvent(dto.Date, meterConfig)}
		},
	)

	return err

}
