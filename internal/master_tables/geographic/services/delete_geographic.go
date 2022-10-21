package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/geographic"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
)

//DeleteGeographicByCodeDto dto to get code
type DeleteGeographicByIdDto struct {
	Id string
}
type DeleteGeographicsService struct {
	GeographicRepository geographic.RepositoryGeographic
	publisher            event.PublisherCreator
	topic                string
}

//NewDeleteGeographicsService create delete service
func NewDeleteGeographicsService(geographicRepository geographic.RepositoryGeographic, publisher event.PublisherCreator, topic string) DeleteGeographicsService {
	return DeleteGeographicsService{
		GeographicRepository: geographicRepository,
		publisher:            publisher,
		topic:                topic,
	}
}

//Handler handle delete
func (s DeleteGeographicsService) Handler(ctx context.Context, dto DeleteGeographicByIdDto) error {
	err := s.GeographicRepository.DeleteGeographicZone(ctx, dto.Id)

	if err != nil {
		return err
	}

	err = event.PublishAllEvents(ctx, s.topic, s.publisher, []master_tables.CalendarPeriodsEvent{master_tables.NewCalendarRedisGenerateEvent(), master_tables.NewFestiveDaysGenerateEvent()})
	return err
}
