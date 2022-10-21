package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"context"
)

type GetPeriodCalendarService struct {
	calendarRepository calendar.RepositoryCalendar
}

// NewGetPeriodCalendarService Generate Get distributor service
func NewGetPeriodCalendarService(calendarRepository calendar.RepositoryCalendar) GetPeriodCalendarService {
	return GetPeriodCalendarService{calendarRepository: calendarRepository}
}

// Handler Get data, count and errors from db
func (s GetPeriodCalendarService) Handler(ctx context.Context, calendarId string, dto GetAllDto) ([]calendar.PeriodCalendar, int, error) {
	ds, count, err := s.calendarRepository.GetAllPeriodCalendars(ctx, calendarId, calendar.Search{
		Q:           dto.Q,
		Limit:       dto.Limit,
		Offset:      dto.Offset,
		Sort:        dto.Sort,
		CurrentUser: dto.CurrentUser,
	})

	return ds, count, err
}
