package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
	"time"
)

type DtoServiceExecute struct {
	Cups          string
	DistributorId string
	StartDate     time.Time
	EndDate       time.Time
	ReadingType   measures.ReadingType
}

type ExecuteServices struct {
	inventoryRepository measures.InventoryRepository
	publisher           event.PublisherCreator
	topic               string
}

func NewExecuteServices(inventoryRepository measures.InventoryRepository, publisher event.PublisherCreator, topic string) *ExecuteServices {
	return &ExecuteServices{
		inventoryRepository: inventoryRepository,
		publisher:           publisher,
		topic:               topic,
	}
}
func (s ExecuteServices) Handler(ctx context.Context, dto DtoServiceExecute) error {

	meterConfigs, err := s.inventoryRepository.ListMeterConfigByCups(ctx, measures.ListMeterConfigByCups{
		CUPS:          dto.Cups,
		StartDate:     dto.StartDate,
		EndDate:       dto.EndDate,
		DistributorId: dto.DistributorId,
	})
	if err != nil {
		return err
	}
	evs := s.generateEvents(dto, meterConfigs)

	err = event.PublishAllEvents(ctx, s.topic, s.publisher, evs)

	return err

}
func (s ExecuteServices) generateEvents(dto DtoServiceExecute, meterConfigs []measures.MeterConfig) []measures.ProcessMeasureEvent {
	evs := make([]measures.ProcessMeasureEvent, 0)

	for sd := dto.StartDate; sd.Before(dto.EndDate) || sd.Equal(dto.EndDate); sd = sd.AddDate(0, 0, 1) {
		var meterConfig measures.MeterConfig
		for _, mc := range meterConfigs {
			if mc.StartDate.Before(sd) && mc.ContractualSituations.InitDate.Before(sd) && mc.EndDate.After(sd) && mc.ContractualSituations.EndDate.After(sd) {
				meterConfig = mc
			}
		}
		ev := process_measures.GenerateEvents(dto.ReadingType, sd, meterConfig)
		evs = append(evs, ev...)
	}
	return evs

}
