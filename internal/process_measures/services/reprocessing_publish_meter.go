package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	measures_services "bitbucket.org/sercide/data-ingestion/internal/measures/services"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
)

type PublishReprocessingMeterService struct {
	*measures_services.OnSchedulerDistributorEvent
	inventoryClient measures.InventoryRepository
	publisher       event.PublisherCreator
	topic           string
}

func NewReprocessingMeterService(
	publisher event.PublisherCreator,
	inventoryClient measures.InventoryRepository,
	topic string,
) *PublishReprocessingMeterService {
	return &PublishReprocessingMeterService{
		inventoryClient: inventoryClient,
		publisher:       publisher,
		topic:           topic,
	}
}

// Handle busca el meterConfig por meterSerialNumber y genera un evento dependiendo del ReadingType
func (svc PublishReprocessingMeterService) Handle(ctx context.Context, dto process_measures.ReSchedulerMeterPayload) error {
	meterConfig, err := svc.inventoryClient.GetMeterConfigByMeter(ctx, measures.GetMeterConfigByMeterQuery{
		MeterSerialNumber: dto.MeterSerialNumber,
		Date:              dto.Date,
	})
	if err != nil {
		return nil
	}

	events := process_measures.GenerateEvents(dto.ReadingType, dto.Date, meterConfig)

	return event.PublishAllEvents(ctx, svc.topic, svc.publisher, events)
}
