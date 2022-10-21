package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"context"
)

//DeletePeriodByCodeDto dto to get code
type DeletePeriodByCodeDto struct {
	Code string
}
type DeletePeriodService struct {
	CalendarRepository calendar.RepositoryCalendar
	publisher          event.PublisherCreator
	topic              string
}

//NewDeletePeriodService create delete service
func NewDeletePeriodService(calendarRepository calendar.RepositoryCalendar, publisher event.PublisherCreator, topic string) DeletePeriodService {
	return DeletePeriodService{
		CalendarRepository: calendarRepository,
		publisher:          publisher,
		topic:              topic,
	}
}

//Handler handle delete
func (s DeletePeriodService) Handler(ctx context.Context, dto DeletePeriodByCodeDto) error {
	err := s.CalendarRepository.DeletePeriod(ctx, dto.Code)

	if err != nil {
		return err
	}
	err = event.PublishEvent(ctx, s.topic, s.publisher, master_tables.NewCalendarRedisGenerateEvent())

	return err
}
