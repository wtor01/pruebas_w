package inventory

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"context"
	"strings"
	"time"
)

type Search struct {
	Q             string
	Limit         int
	Offset        *int
	Sort          map[string]string
	CurrentUser   auth.User
	DistributorID string
	Values        map[string]string
}

type PaginationConfig struct {
	Limit  int
	Offset int
}

func NewSort(arr []string) map[string]string {
	sort := make(map[string]string)

	for _, param := range arr {
		sortParam := strings.Split(param, " ")

		if len(sortParam) != 2 {
			continue
		}

		if sortParam[1] != "asc" && sortParam[1] != "desc" {
			continue
		}

		sort[sortParam[0]] = sortParam[1]
	}
	return sort
}

type GetMeasureEquipmentByMeterQuery struct {
	MeterID           string
	MeterSerialNumber string
	Time              time.Time
}
type GetServicePointByCupsQuery struct {
	CUPS string
	Time time.Time
}

type GetServicePointContractQuery struct {
	CUP       string
	StartDate time.Time
	EndDate   time.Time
}

type GetMeterConfigByMeterQuery struct {
	MeterSerialNumber string
	Date              time.Time
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=RepositoryInventory
type RepositoryInventory interface {
	ListDistributors(ctx context.Context, search Search) ([]Distributor, int, error)
	GetDistributorBySmarkiaID(ctx context.Context, id string) (Distributor, error)
	GetDistributorByID(ctx context.Context, id string) (Distributor, error)
	GetDistributorByCdos(ctx context.Context, cdos string) (Distributor, error)
	ListMeters(ctx context.Context, search Search) ([]MeasureEquipment, int, error)
	GroupMetersByType(ctx context.Context, distributorID string) (map[string]int, error)
	GetMeterConfig(ctx context.Context, search GetMeterConfigByMeterQuery) (MeterConfig, error)
	GetMeasureEquipmentByID(ctx context.Context, distributorID string, id string) (MeasureEquipment, error)
	GetMeasureEquipmentBySmarkiaID(ctx context.Context, id string) (MeasureEquipment, error)
	GetServicePoints(ctx context.Context, dto ServicePointSchedulerDto) ([]ServicePointScheduler, error)
	GetServicePointByMeter(ctx context.Context, query GetMeasureEquipmentByMeterQuery) (ServicePointScheduler, error)
}
