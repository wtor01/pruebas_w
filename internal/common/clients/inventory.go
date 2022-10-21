package clients

import (
	"context"
	"time"
)

type Distributor struct {
	ID        string
	Name      string
	SmarkiaID string
	CDOS      string
}

const (
	EquipmentTelegestionType string = "TLG"
	EquipmentTelemedidaType  string = "TLM"
	EquipmentOtherType       string = "OTHER"
)

type Equipment struct {
	ID                string
	Type              string
	SerialNumber      string
	Technology        string
	Brand             string
	Model             string
	ActiveConstant    float64
	ReactiveConstant  float64
	MaximeterConstant float64
}

type ServicePoint struct {
	DistributorId string
	ServiceType   string
	PointType     string
	MeterType     []string
	ReadingType   string
	time          time.Time
}

type GetServicePointConfigByMeterIdDto struct {
	MeterID           string
	MeterSerialNumber string
	Time              time.Time
}

type ServicePointScheduler struct {
	MeterConfigId          string
	MeterId                string
	DistributorId          string
	DistributorCode        string
	Cups                   string
	MeterSerialNumber      string
	ServiceType            string
	ReadingType            string
	PointType              string
	MeasurePointType       string
	PriorityContract       string
	TlgCode                string
	TlgType                string
	Calendar               string
	Type                   string
	StartDate              time.Time
	EndDate                time.Time
	CurveType              string
	Technology             string
	ContractualSituationId string
	TariffID               string
	Coefficient            string
	CalendarCode           string
	P1Demand               float64
	P2Demand               float64
	P3Demand               float64
	P4Demand               float64
	P5Demand               float64
	P6Demand               float64
	AI                     int
	AE                     int
	R1                     int
	R2                     int
	R3                     int
	R4                     int
	ContractStartDate      time.Time
	ContractEndDate        time.Time
	ServicePointID         string
	GeographicID           string
}

type MeterConfig struct {
	ID               string
	StartDate        time.Time
	EndDate          time.Time
	AI               int
	AE               int
	R1               int
	R2               int
	R3               int
	R4               int
	M                int
	E                int
	CurveType        string
	ReadingType      string
	CalendarID       string
	PriorityContract string
	MeterID          string
	MeasurePointID   string
	ServicePointID   string
	TlgCode          string
	Type             string
	RentingPrice     float64
	SerialNumber     string
	DistributorID    string
	CUPS             string
	MeasurePointType string
}

type ServicePointSchedulerDto struct {
	DistributorId string
	ServiceType   string
	PointType     string
	MeterType     []string
	ReadingType   string
	Time          time.Time
	Limit         int
	Offset        int
}

type ListDistributorsDto struct {
	Q      string
	Limit  int
	Offset *int
	Sort   map[string]string
	Values map[string]string
}

type GetMeterConfigByMeterDto struct {
	MeterSerialNumber string
	Date              time.Time
}

//go:generate mockery --case=snake --outpkg=mocks --output=./clients_mocks --name=Inventory
type Inventory interface {
	GetAllDistributors(ctx context.Context) ([]Distributor, error)
	ListDistributors(ctx context.Context, dto ListDistributorsDto) ([]Distributor, int, error)
	GetDistributorBySmarkiaId(ctx context.Context, smarikiaId string) (Distributor, error)
	GetDistributorById(ctx context.Context, id string) (Distributor, error)
	GetDistributorByCdos(ctx context.Context, cdos string) (Distributor, error)
	GetMeasureEquipmentBySmarkiaId(ctx context.Context, smarikiaId string) (Equipment, error)
	ListMeasureEquipmentByDistributorId(ctx context.Context, distributorID string) ([]Equipment, error)
	GroupMetersByType(ctx context.Context, distributorID string) (map[string]int, error)
	GetServicePoints(ctx context.Context, inventoryClient ServicePointSchedulerDto) ([]ServicePointScheduler, error)
	GetServicePointConfigByMeter(ctx context.Context, query GetServicePointConfigByMeterIdDto) (ServicePointScheduler, error)
	GetMeterConfigByMeterService(ctx context.Context, query GetMeterConfigByMeterDto) (MeterConfig, error)
}
