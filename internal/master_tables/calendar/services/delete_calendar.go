package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
)

//DeleteCalendarByCodeDto dto to get code
type DeleteCalendarByIdDto struct {
	Id string
}
type DeleteCalendarService struct {
	CalendarRepository calendar.RepositoryCalendar
	publisher          event.PublisherCreator
	topic              string
}

//NewDeleteCalendarService create delete service
func NewDeleteCalendarService(calendarRepository calendar.RepositoryCalendar, publisher event.PublisherCreator, topic string) DeleteCalendarService {
	return DeleteCalendarService{
		CalendarRepository: calendarRepository,
		publisher:          publisher,
		topic:              topic,
	}
}

//Handler handle delete
func (s DeleteCalendarService) Handler(ctx context.Context, dto DeleteCalendarByIdDto) error {
	err := s.CalendarRepository.DeleteCalendar(ctx, dto.Id)

	if err != nil {
		return err
	}
	err = event.PublishEvent(ctx, s.topic, s.publisher, master_tables.NewCalendarRedisGenerateEvent())

	return err
}
