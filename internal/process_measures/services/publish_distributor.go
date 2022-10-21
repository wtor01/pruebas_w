package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
)

type PublishDistributorService struct {
	publisher       event.PublisherCreator
	inventoryClient clients.Inventory
	schedulerClient *SearchScheduler
	topic           string
}

func NewPublishDistributorService(
	publisher event.PublisherCreator,
	schedulerClient *SearchScheduler,
	inventoryClient clients.Inventory,
	topic string,
) *PublishDistributorService {
	return &PublishDistributorService{
		publisher:       publisher,
		inventoryClient: inventoryClient,
		schedulerClient: schedulerClient,
		topic:           topic,
	}
}

func (svc PublishDistributorService) Handle(ctx context.Context, dto measures.SchedulerEventPayload) error {
	schedulersToPublish := make([]measures.SchedulerEvent, 0)

	if dto.DistributorId == "" {
		distributors, err := svc.inventoryClient.GetAllDistributors(ctx)
		if err != nil {
			return err
		}
		for _, d := range distributors {
			measuresScheduler, err := svc.schedulerClient.Handler(ctx, SearchSchedulerDTO{
				DistributorId: d.ID,
				ServiceType:   string(dto.ServiceType),
				PointType:     dto.PointType,
				MeterType:     dto.MeterType,
				ReadingType:   string(dto.ReadingType),
			})
			if err != nil {
				return err
			}

			if len(measuresScheduler) == 0 {
				schedulersToPublish = append(
					schedulersToPublish,
					process_measures.NewProcessByDistributorEvent(measures.SchedulerEventPayload{
						ID:            dto.ID,
						DistributorId: d.ID,
						Name:          dto.Name,
						Description:   dto.Description,
						ServiceType:   dto.ServiceType,
						PointType:     dto.PointType,
						MeterType:     dto.MeterType,
						ReadingType:   dto.ReadingType,
						Format:        dto.Format,
						Date:          dto.Date,
					}),
				)
			}

		}
	} else {
		schedulersToPublish = append(
			schedulersToPublish,
			process_measures.NewProcessByDistributorEvent(measures.SchedulerEventPayload{
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
			}),
		)
	}

	err := event.PublishAllEvents(ctx, svc.topic, svc.publisher, schedulersToPublish)

	return err
}
