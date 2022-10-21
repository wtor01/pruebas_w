package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	measures_services "bitbucket.org/sercide/data-ingestion/internal/measures/services"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"math"

	"time"
)

type PublishReprocessingServicePointService struct {
	*measures_services.OnSchedulerDistributorEvent
	grossRepository gross_measures.GrossMeasureRepository
	publisher       event.PublisherCreator
	topic           string
	limit           int
}

func NewReprocessingServicePointService(
	publisher event.PublisherCreator,
	topic string,
	grossRepository gross_measures.GrossMeasureRepository,
	limit int,
) *PublishReprocessingServicePointService {
	return &PublishReprocessingServicePointService{
		publisher:       publisher,
		topic:           topic,
		grossRepository: grossRepository,
		limit:           limit,
	}
}

// Handle Busca entre un rango de fechas por distribuidor y genera un evento ReSchedulerMeterEvent
func (svc PublishReprocessingServicePointService) Handle(ctx context.Context, dto process_measures.ReSchedulerEventPayload) error {
	if dto.Limit == 0 {
		return svc.PaginateEvents(ctx, dto, process_measures.NewReprocessingProcessByDistributorEvent)
	}
	rescheduler, err := svc.grossRepository.ListGrossMeasuresFromGenerationDate(ctx, gross_measures.QueryListForProcessCurveGenerationDate{
		ReadingType:   dto.ReadingType,
		DistributorId: dto.DistributorId,
		StartDate:     dto.StartDate,
		EndDate:       dto.EndDate,
		Limit:         dto.Limit,
		Offset:        dto.Offset,
	})

	if err != nil {
		return err
	}

	events := utils.MapSlice(rescheduler, func(msn gross_measures.MeasureCurveMeterSerialNumber) process_measures.ReSchedulerMeterEvent {
		return process_measures.NewReprocessingMeterEvent(process_measures.ReSchedulerMeterPayload{
			MeterSerialNumber: msn.MeterSerialNumber,
			Date:              time.Date(msn.Year, time.Month(msn.Month), msn.Day, 0, 0, 0, 0, time.UTC),
			ReadingType:       dto.ReadingType,
		})
	})

	return event.PublishAllEvents(ctx, svc.topic, svc.publisher, events)
}
func (svc PublishReprocessingServicePointService) PaginateEvents(
	ctx context.Context,
	dto process_measures.ReSchedulerEventPayload,
	generateEvent func(dto process_measures.ReSchedulerEventPayload) process_measures.ReSchedulerEvent,
) error {

	count, err := svc.grossRepository.CountGrossMeasuresFromGenerationDate(ctx, gross_measures.QueryListForProcessCurveGenerationDate{
		ReadingType:   dto.ReadingType,
		DistributorId: dto.DistributorId,
		StartDate:     dto.StartDate,
		EndDate:       dto.EndDate,
	})
	offset := 0

	if err != nil {
		return err
	}

	msg := make([]process_measures.ReSchedulerEvent, 0, int(math.Ceil(float64(count)/float64(svc.limit))))

	for count > 0 {
		count -= svc.limit
		msg = append(
			msg,
			generateEvent(process_measures.ReSchedulerEventPayload{
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

	err = event.PublishAllEvents(ctx, svc.topic, svc.publisher, msg)

	return err
}
