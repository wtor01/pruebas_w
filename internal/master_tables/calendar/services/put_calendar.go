package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
)

type ModifyCalendarService struct {
	CalendarRepository calendar.RepositoryCalendar
	publisher          event.PublisherCreator
	topic              string
}

//NewModifyCalendarService create new Modify service
func NewModifyCalendarService(calendarRepository calendar.RepositoryCalendar,
	publisher event.PublisherCreator, topic string) ModifyCalendarService {
	return ModifyCalendarService{
		CalendarRepository: calendarRepository,
		publisher:          publisher, topic: topic,
	}
}

//Handler handle modify service
func (s ModifyCalendarService) Handler(ctx context.Context, calendarId string, cal calendar.Calendar) error {
	err := s.CalendarRepository.ModifyCalendar(ctx, calendarId, cal)

	if err != nil {
		return err
	}
	err = event.PublishEvent(ctx, s.topic, s.publisher, master_tables.NewCalendarRedisGenerateEvent())

	return err
}
