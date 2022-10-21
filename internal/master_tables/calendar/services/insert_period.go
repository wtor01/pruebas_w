package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
)

type InsertPeriodService struct {
	CalendarRepository calendar.RepositoryCalendar
	publisher          event.PublisherCreator
	topic              string
}

//NewInsertPeriodService create a service struct
func NewInsertPeriodService(calendarRepository calendar.RepositoryCalendar, publisher event.PublisherCreator, topic string) InsertPeriodService {
	return InsertPeriodService{
		CalendarRepository: calendarRepository,
		publisher:          publisher,
		topic:              topic,
	}
}

//Handler handle insert service
func (s InsertPeriodService) Handler(ctx context.Context, calendarId string, pc calendar.PeriodCalendar) error {

	err := s.CalendarRepository.InsertPeriod(ctx, calendarId, pc)

	if err != nil {
		return err
	}
	err = event.PublishEvent(ctx, s.topic, s.publisher, master_tables.NewCalendarRedisGenerateEvent())

	return err
}
