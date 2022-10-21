package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
	"time"
)

type PublishDistributorDTO struct {
	ProcessName string
	Date        time.Time
}

type PublishDistributorService struct {
	publisher       event.PublisherCreator
	inventoryClient clients.Inventory
	topic           string
}

func NewPublishDistributorService(publisher event.PublisherCreator, inventoryClient clients.Inventory, topic string) *PublishDistributorService {
	return &PublishDistributorService{publisher: publisher, inventoryClient: inventoryClient, topic: topic}
}

func (svc PublishDistributorService) Handle(ctx context.Context, dto PublishDistributorDTO) error {

	distributors, _, err := svc.inventoryClient.ListDistributors(ctx, clients.ListDistributorsDto{
		Limit: 1000,
		Values: map[string]string{
			"is_smarkia_active": "1",
		},
	})

	if err != nil {
		return err
	}

	messages := make([]smarkia.MessageDistributorProcess, 0, cap(distributors))
	for _, d := range distributors {
		messages = append(messages, smarkia.NewMessageDistributorProcess(d.ID, d.SmarkiaID, dto.ProcessName, d.CDOS, dto.Date))
	}

	err = event.PublishAllEvents(ctx, svc.topic, svc.publisher, messages)

	return err
}
