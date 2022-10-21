package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"context"
	"strings"
)

type InsertCalendarService struct {
	calendarRepository calendar.RepositoryCalendar
}

//NewInsertCalendarService create a service struct
func NewInsertCalendarService(calendarRepository calendar.RepositoryCalendar) InsertCalendarService {
	return InsertCalendarService{calendarRepository: calendarRepository}
}

//Handler handle insert service
func (s InsertCalendarService) Handler(ctx context.Context, calendars calendar.Calendar) error {
	//convert to uppercase
	calendars.Code = strings.ToUpper(calendars.Code)
	err := s.calendarRepository.InsertCalendars(ctx, calendars)

	return err
}
