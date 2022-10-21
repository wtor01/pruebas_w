package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
)

type ModifyPeriodCalendarService struct {
	CalendarRepository calendar.RepositoryCalendar
	publisher          event.PublisherCreator
	topic              string
}

//NewModifyPeriodCalendarService create new Modify service
func NewModifyPeriodCalendarService(calendarRepository calendar.RepositoryCalendar, publisher event.PublisherCreator, topic string) ModifyPeriodCalendarService {
	return ModifyPeriodCalendarService{
		CalendarRepository: calendarRepository,
		publisher:          publisher,
		topic:              topic,
	}
}

//Handler handle modify service
func (s ModifyPeriodCalendarService) Handler(ctx context.Context, code string, cal calendar.PeriodCalendar) error {
	err := s.CalendarRepository.ModifyPeriodicCalendar(ctx, code, cal)

	if err != nil {
		return err
	}
	err = event.PublishEvent(ctx, s.topic, s.publisher, master_tables.NewCalendarRedisGenerateEvent())

	return err
}
