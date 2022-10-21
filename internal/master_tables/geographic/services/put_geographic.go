package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/geographic"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
)

type ModifyGeographicsService struct {
	GeographicRepository geographic.RepositoryGeographic
	publisher            event.PublisherCreator
	topic                string
}

//NewModifyGeographicsService create new Modify service
func NewModifyGeographicsService(geographicRepository geographic.RepositoryGeographic, publisher event.PublisherCreator, topic string) ModifyGeographicsService {
	return ModifyGeographicsService{
		GeographicRepository: geographicRepository,
		publisher:            publisher,
		topic:                topic,
	}
}

//Handler handle modify service
func (s ModifyGeographicsService) Handler(ctx context.Context, geographicId string, zones geographic.GeographicZones) error {
	err := s.GeographicRepository.ModifyGeographicZone(ctx, geographicId, zones)

	if err != nil {
		return err
	}

	err = event.PublishAllEvents(ctx, s.topic, s.publisher, []master_tables.CalendarPeriodsEvent{master_tables.NewCalendarRedisGenerateEvent(), master_tables.NewFestiveDaysGenerateEvent()})
	return err
}
