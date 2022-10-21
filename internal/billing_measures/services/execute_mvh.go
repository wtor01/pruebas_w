package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"

	"golang.org/x/net/context"
	"time"
)

type ServiceExecuteMvh struct {
	inventoryRepository measures.InventoryRepository
	publisher           event.PublisherCreator
	topic               string
}

type DtoServiceExecuteMvh struct {
	Cups          string
	DistributorId string
	Date          time.Time
	Location      *time.Location
}

func NewServiceExecuteMvh(
	inventoryRepository measures.InventoryRepository,
	publisher event.PublisherCreator,
	topic string,
) *ServiceExecuteMvh {
	return &ServiceExecuteMvh{
		inventoryRepository: inventoryRepository,
		publisher:           publisher,
		topic:               topic,
	}
}

func (s ServiceExecuteMvh) Handler(ctx context.Context, dto DtoServiceExecuteMvh) error {
	dateTime := time.Date(dto.Date.Year(), dto.Date.Month(), dto.Date.Day(), 0, 0, 0, 0, dto.Location).UTC()
	meterConfig, err := s.inventoryRepository.GetMeterConfigByCupsAPI(ctx, measures.GetMeterConfigByCupsQuery{
		CUPS:        dto.Cups,
		Time:        dateTime,
		Distributor: dto.DistributorId,
	})
	if err != nil {
		return err
	}
	ev := billing_measures.NewProcessMvhEvent(dateTime, meterConfig)
	pub, err := s.publisher(ctx)
	if err != nil {
		return err
	}
	msg, err := ev.Marshal()
	if err != nil {
		return err
	}
	err = pub.Publish(ctx, s.topic, msg, ev.GetAttributes())

	return err

}
