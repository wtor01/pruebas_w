package calendar

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"context"
	uuid "github.com/satori/go.uuid"
	"time"
)

//Search search used in Get all
type Search struct {
	Q             string
	Limit         int
	Offset        *int
	Sort          map[string]string
	CurrentUser   auth.User
	DistributorID string
}

//go:generate mockery --case=snake --outpkg=mocks --output=../../platform/mocks --name=RepositoryCalendar
type RepositoryCalendar interface {
	GetAllCalendars(ctx context.Context, search Search) ([]Calendar, int, error)
	InsertCalendars(ctx context.Context, calendar Calendar) error
	ModifyCalendar(ctx context.Context, calendarId string, zones Calendar) error
	DeleteCalendar(ctx context.Context, calendarId string) error
	GetCalendar(ctx context.Context, calendarId string) (Calendar, error)

	ModifyPeriodicCalendar(ctx context.Context, code string, period PeriodCalendar) error
	DeletePeriod(ctx context.Context, code string) error
	InsertPeriod(ctx context.Context, code string, period PeriodCalendar) error
	GetAllPeriodCalendars(ctx context.Context, code string, search Search) ([]PeriodCalendar, int, error)
	GetPeriodsByDate(ctx context.Context, calendarCode string, date time.Time) ([]PeriodCalendar, error)
	GetPeriodsActive(ctx context.Context) ([]PeriodCalendar, error)
	GetPeriodById(ctx context.Context, periodId uuid.UUID) (PeriodCalendar, error)
}
