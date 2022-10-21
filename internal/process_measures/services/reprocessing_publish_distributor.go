package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
	"errors"
	"time"
)

type ReprocessingPublishDistributorService struct {
	publisher                  event.PublisherCreator
	inventoryClient            clients.Inventory
	schedulerClient            *SearchScheduler
	topic                      string
	generatorDate              func() time.Time
	ReprocessingDateRepository measures.ReprocessingDateRepository
}

const getReprocessingDistributorRedis = "ReprocessingTime"

func NewReprocessingPublishDistributorService(
	publisher event.PublisherCreator,
	schedulerClient *SearchScheduler,
	inventoryClient clients.Inventory,
	topic string,
	redisRepo measures.ReprocessingDateRepository,
	generatorDate func() time.Time,
) *ReprocessingPublishDistributorService {
	return &ReprocessingPublishDistributorService{
		publisher:                  publisher,
		inventoryClient:            inventoryClient,
		schedulerClient:            schedulerClient,
		topic:                      topic,
		generatorDate:              generatorDate,
		ReprocessingDateRepository: redisRepo,
	}
}

// Handle busca en redis la fecha si no hay la guarda y da ko si hay fecha busca los distribuidores genera eventos y actualiza fecha
func (svc ReprocessingPublishDistributorService) Handle(ctx context.Context, dto measures.SchedulerEventPayload) error {
	schedulersToPublish := make([]process_measures.ReSchedulerEvent, 0)
	err, startDate := svc.ReprocessingDateRepository.GetDate(ctx, getReprocessingDistributorRedis)
	if err != nil {
		err = svc.ReprocessingDateRepository.SetDate(ctx, getReprocessingDistributorRedis, svc.generatorDate())
		if err != nil {
			return err
		}
		return errors.New("no redis cache start date")
	}

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
					process_measures.NewReprocessingProcessByDistributorEvent(process_measures.ReSchedulerEventPayload{
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
						StartDate:     startDate,
						EndDate:       svc.generatorDate(),
					}),
				)
			}
		}
		err = svc.ReprocessingDateRepository.SetDate(ctx, getReprocessingDistributorRedis, svc.generatorDate())
		if err != nil {
			return err
		}

	} else {
		schedulersToPublish = append(
			schedulersToPublish,
			process_measures.NewReprocessingProcessByDistributorEvent(process_measures.ReSchedulerEventPayload{
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
				StartDate:     startDate,
				EndDate:       svc.generatorDate(),
			}),
		)
	}

	err = event.PublishAllEvents(ctx, svc.topic, svc.publisher, schedulersToPublish)

	return err
}
