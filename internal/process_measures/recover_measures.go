package process_measures

import (
	"context"
	"time"
)

type RecoverMeasures struct {
	ID                string
	MeterId           string
	Date              time.Time
	DistributorID     string
	DistributorCode   string
	CUPS              string
	MeterSerialNumber string
	MeasurePointType  string
	ContractNumber    string
	MeterConfigId     string
	ReadingType       string

	ServiceType string
	PointType   string
}

type PaginationRecoverMeasures struct {
	Limit  int
	Offset *int
}

type SearchRecoverMeasures struct {
	MeterId           string
	Date              time.Time
	DistributorID     string
	DistributorCode   string
	CUPS              string
	MeterSerialNumber string
	MeasurePointType  string
	ContractNumber    string
	MeterConfigId     string
	ReadingType       string

	ServiceType string
	PointType   string
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=RecoverMeasuresRepository
type RecoverMeasuresRepository interface {
	AllRecoverMeasures(ctx context.Context, s PaginationRecoverMeasures) ([]RecoverMeasures, int, error)
	SaveRecoverMeasures(ctx context.Context, rm RecoverMeasures) error
	SearchRecoverMeasures(ctx context.Context, s RecoverMeasures) (RecoverMeasures, int, error)
	DeleteRecoverMeasures(ctx context.Context, id string) error
	GetRecoverMeasures(ctx context.Context, id string) (RecoverMeasures, error)
}
