package seasons

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"context"
	uuid "github.com/satori/go.uuid"
)

type Search struct {
	Q             string
	Limit         int
	Offset        *int
	Sort          map[string]string
	CurrentUser   auth.User
	DistributorID string
}

//go:generate mockery --case=snake --outpkg=mocks --output=../../platform/mocks --name=RepositorySeasons
type RepositorySeasons interface {
	GetAllSeasons(ctx context.Context, search Search) ([]Seasons, int, error)
	GetSeason(ctx context.Context, seasonId string) (Seasons, error)
	InsertSeason(ctx context.Context, season Seasons) error
	ModifySeason(ctx context.Context, seasonId uuid.UUID, season Seasons) error
	DeleteSeason(ctx context.Context, seasonId uuid.UUID) error

	GetAllDayTypes(ctx context.Context, seasonId string, search Search) ([]DayTypes, int, error)
	GetDayType(ctx context.Context, dayTypeId string) (DayTypes, error)
	GetDayTypeByMonth(ctx context.Context, month int, isFestive bool) (DayTypes, error)
	InsertDayTypes(ctx context.Context, seasonId string, dayT DayTypes) error
	ModifyDayTypes(ctx context.Context, dayTypeId uuid.UUID, dayT DayTypes) error
	DeleteDayType(ctx context.Context, dayTypeId uuid.UUID) error
}
