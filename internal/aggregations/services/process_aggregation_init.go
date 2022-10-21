package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
)

type ProcessAggregationInitService struct {
	publisher       event.PublisherCreator
	inventoryClient clients.Inventory
	topic           string
}

func NewProcessAggregationInitService(publisher event.PublisherCreator, inventoryClient clients.Inventory, topic string) *ProcessAggregationInitService {
	return &ProcessAggregationInitService{publisher: publisher, inventoryClient: inventoryClient, topic: topic}
}

func (svc ProcessAggregationInitService) Handler(ctx context.Context, dto aggregations.ConfigScheduler) error {
	if dto.StartDate.After(dto.Date) || (!dto.EndDate.IsZero() && dto.EndDate.Before(dto.Date)) {
		return nil
	}

	distributors, err := svc.inventoryClient.GetAllDistributors(ctx)

	if err != nil {
		return err
	}

	schedulersToPublish := make([]aggregations.SchedulerEvent, 0, cap(distributors))

	for _, d := range distributors {
		dto.DistributorId = d.ID
		dto.DistributorCDOS = d.CDOS
		schedulersToPublish = append(schedulersToPublish, aggregations.NewSchedulerDistributorEvent(dto))
	}

	err = event.PublishAllEvents(ctx, svc.topic, svc.publisher, schedulersToPublish)

	return err
}
