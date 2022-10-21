package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"context"
)

type GetAllDto struct {
	Q           string
	Limit       int
	Offset      *int
	Sort        map[string]string
	CurrentUser auth.User
}

type GetCalendarService struct {
	calendarRepository calendar.RepositoryCalendar
}

// NewGetCalendarService Generate Get distributor service
func NewGetCalendarService(calendarRepository calendar.RepositoryCalendar) GetCalendarService {
	return GetCalendarService{calendarRepository: calendarRepository}
}

// Handler Get data, count and errors from db
func (s GetCalendarService) Handler(ctx context.Context, dto GetAllDto) ([]calendar.Calendar, int, error) {
	ds, count, err := s.calendarRepository.GetAllCalendars(ctx, calendar.Search{
		Q:           dto.Q,
		Limit:       dto.Limit,
		Offset:      dto.Offset,
		Sort:        dto.Sort,
		CurrentUser: dto.CurrentUser,
	})

	return ds, count, err
}
