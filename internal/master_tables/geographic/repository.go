package geographic

import (
	"context"
)

//go:generate mockery --case=snake --outpkg=mocks --output=../../platform/mocks --name=RepositoryGeographic
type RepositoryGeographic interface {
	GetAllGeographicZones(ctx context.Context, search Search) ([]GeographicZones, int, error)
	GetGeographicZone(ctx context.Context, geographicId string) (GeographicZones, error)
	InsertGeographicZone(ctx context.Context, zones GeographicZones) error
	ModifyGeographicZone(ctx context.Context, geographicId string, zones GeographicZones) error
	DeleteGeographicZone(ctx context.Context, geographicId string) error
}
