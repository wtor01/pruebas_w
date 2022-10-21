package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
	"math"
)

type OnSchedulerDistributorEvent struct {
	publish             event.PublisherCreator
	inventoryRepository measures.InventoryRepository
	topic               string
	limit               int
	tracer              trace.Tracer
}

func NewOnSchedulerDistributorEvent(publisher event.PublisherCreator, inventoryClient measures.InventoryRepository, topic string, limit int) *OnSchedulerDistributorEvent {
	return &OnSchedulerDistributorEvent{
		publish:             publisher,
		inventoryRepository: inventoryClient,
		topic:               topic,
		limit:               limit,
		tracer:              telemetry.GetTracer(),
	}
}

func (svc OnSchedulerDistributorEvent) PaginateEvents(
	ctx context.Context,
	dto measures.SchedulerEventPayload,
	generateEvent func(dto measures.SchedulerEventPayload) measures.SchedulerEvent,
) error {
	ctx, span := svc.tracer.Start(ctx, "OnSchedulerDistributorEvent[T] - PaginateEvents")
	defer span.End()

	count, err := svc.inventoryRepository.CountMeterConfigByDate(ctx, measures.ListMeterConfigByDateQuery{
		DistributorID: dto.DistributorId,
		ServiceType:   dto.ServiceType,
		PointType:     dto.PointType,
		MeterType:     dto.MeterType,
		ReadingType:   dto.ReadingType,
		Date:          dto.Date,
	})
	offset := 0

	if err != nil {
		return err
	}

	msg := make([]measures.SchedulerEvent, 0, int(math.Ceil(float64(count)/float64(svc.limit))))

	for count > 0 {
		count -= svc.limit
		msg = append(
			msg,
			generateEvent(measures.SchedulerEventPayload{
				ID:            dto.ID,
				DistributorId: dto.DistributorId,
				Name:          dto.Name,
				Description:   dto.Description,
				ServiceType:   dto.ServiceType,
				PointType:     dto.PointType,
				MeterType:     dto.MeterType,
				ReadingType:   dto.ReadingType,
				Format:        dto.Format,
				Date:          dto.Date,
				Limit:         svc.limit,
				Offset:        offset,
			}),
		)
		offset += svc.limit
	}

	err = event.PublishAllEvents(ctx, svc.topic, svc.publish, msg)

	return err
}

func (svc OnSchedulerDistributorEvent) ListMeterConfig(ctx context.Context, dto measures.SchedulerEventPayload) ([]measures.MeterConfig, error) {
	ctx, span := svc.tracer.Start(ctx, "OnSchedulerDistributorEvent[T] - ListMeterConfig")
	defer span.End()

	servicePointConfigs, err := svc.inventoryRepository.ListMeterConfigByDate(ctx, measures.ListMeterConfigByDateQuery{
		DistributorID: dto.DistributorId,
		ServiceType:   dto.ServiceType,
		PointType:     dto.PointType,
		MeterType:     dto.MeterType,
		ReadingType:   dto.ReadingType,
		Date:          dto.Date,
		Limit:         dto.Limit,
		Offset:        dto.Offset,
	})

	return servicePointConfigs, err
}

func (svc OnSchedulerDistributorEvent) PublishAllEvents(ctx context.Context, msg []measures.ProcessMeasureEvent) error {
	ctx, span := svc.tracer.Start(ctx, "OnSchedulerDistributorEvent[T] - PublishAllEvents")
	defer span.End()

	return event.PublishAllEvents(ctx, svc.topic, svc.publish, msg)

}

func (svc OnSchedulerDistributorEvent) Handle(
	ctx context.Context,
	dto measures.SchedulerEventPayload,
	generatePaginateEvent func(dto measures.SchedulerEventPayload) measures.SchedulerEvent,
	generatePublishEvent func(dto measures.SchedulerEventPayload, meterConfig measures.MeterConfig) []measures.ProcessMeasureEvent,

) error {
	ctx, span := svc.tracer.Start(ctx, "Handle")
	defer span.End()

	// todavia no se ha paginado el evento
	if dto.Limit == 0 {
		return svc.PaginateEvents(ctx, dto, generatePaginateEvent)
	}

	servicePointConfigs, err := svc.ListMeterConfig(ctx, dto)

	if err != nil {
		return err
	}

	msg := make([]measures.ProcessMeasureEvent, 0)
	for _, s := range servicePointConfigs {
		eventToPublish := generatePublishEvent(dto, s)
		msg = append(msg, eventToPublish...)
	}

	err = svc.PublishAllEvents(ctx, msg)

	return err
}
