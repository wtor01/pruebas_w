package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"context"
)

//GetCalendarByIdDto struct send code for Get
type GetCalendarByIdDto struct {
	Id string
}
type GetOneCalendarService struct {
	CalendarRepository calendar.RepositoryCalendar
}

//NewGetOneCalendarService create get one struct
func NewGetOneCalendarService(calendarRepository calendar.RepositoryCalendar) GetOneCalendarService {
	return GetOneCalendarService{CalendarRepository: calendarRepository}
}

//Handler handle get one geographic service
func (s GetOneCalendarService) Handler(ctx context.Context, dto GetCalendarByIdDto) (calendar.Calendar, error) {
	res, err := s.CalendarRepository.GetCalendar(ctx, dto.Id)

	return res, err
}
