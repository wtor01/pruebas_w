package measures

import (
	"context"
	"time"
)

type Meter struct {
	ID                string    `json:"id" bson:"_id"`
	SmarkiaID         string    `json:"smarkia_id" bson:"smarkia_id"`
	SerialNumber      string    `json:"serial_number" bson:"serial_number"`
	Brand             string    `json:"brand" bson:"brand"`
	Model             string    `json:"model" bson:"model"`
	ActiveConstant    float64   `json:"active_constant" bson:"active_constant"`
	ReactiveConstant  float64   `json:"reactive_constant" bson:"reactive_constant"`
	MaximeterConstant float64   `json:"maximeter_constant" bson:"maximeter_constant"`
	Type              MeterType `json:"type" bson:"type"`
	Technology        string    `json:"technology" bson:"technology"`
}

type ServicePoint struct {
	ID                  string      `json:"id" bson:"_id"`
	Name                string      `json:"name" bson:"name"`
	Cups                string      `json:"cups" bson:"cups"`
	StartDate           time.Time   `json:"start_date" bson:"start_date"`
	ServiceType         ServiceType `json:"service_type" bson:"service_type"`
	PointType           PointType   `json:"point_type" bson:"point_type"`
	PointTensionLevel   string      `json:"point_tension_level" bson:"point_tension_level"`
	MeasureTensionLevel string      `json:"measure_tension_level" bson:"measure_tension_level"`
	TensionSection      int         `json:"tension_section" bson:"tension_section"`
	TransfNomDemand     int         `json:"transf_nom_demand" bson:"transf_nom_demand"`
	Disabled            bool        `json:"disabled" bson:"disabled"`
	SolarZoneID         string      `json:"solar_zone_id" bson:"solar_zone_id"`
	SolarZoneGroupID    string      `json:"solar_zone_group_id" bson:"solar_zone_group_id"`
	TechnologyID        string      `json:"technology_id" bson:"technology_id"`
}

type MeasurePoint struct {
	ID         string           `json:"id" bson:"_id"`
	Type       MeasurePointType `json:"type" bson:"type"`
	LossesPerc float64          `json:"losses_perc" bson:"losses_perc"`
	LossesCoef float64          `json:"losses_coef" bson:"losses_coef"`
}

type ContractualSituations struct {
	ID           string    `json:"id" bson:"_id"`
	TariffID     string    `json:"tariff_id" bson:"tariff_id"`
	P1Demand     float64   `json:"p1_demand" bson:"p1_demand"`
	P2Demand     float64   `json:"p2_demand" bson:"p2_demand"`
	P3Demand     float64   `json:"p3_demand" bson:"p3_demand"`
	P4Demand     float64   `json:"p4_demand" bson:"p4_demand"`
	P5Demand     float64   `json:"p5_demand" bson:"p5_demand"`
	P6Demand     float64   `json:"p6_demand" bson:"p6_demand"`
	RetailerCode int       `json:"retailer_code" bson:"retailer_code"`
	InitDate     time.Time `json:"init_date" bson:"init_date"`
	EndDate      time.Time `json:"end_date" bson:"end_date"`
	RetailerCdos string    `json:"retailer_cdos" bson:"retailer_cdos"`
	RetailerName string    `json:"retailer_name" bson:"retailer_name"`
}

type MeterConfig struct {
	ID                    string                `json:"id" bson:"id"`
	StartDate             time.Time             `json:"start_date" bson:"start_date"`
	EndDate               time.Time             `json:"end_date" bson:"end_date"`
	CurveType             RegisterType          `json:"curve_type" bson:"curve_type"`
	ReadingType           string                `json:"reading_type" bson:"reading_type"`
	Inaccessible          bool                  `json:"inaccessible" bson:"inaccessible"`
	PriorityContract      PriorityContract      `json:"priority_contract" bson:"priority_contract"`
	CalendarID            string                `json:"calendar_id" bson:"calendar_id"`
	AI                    int                   `json:"ai" bson:"ai"`
	AE                    int                   `json:"ae" bson:"ae"`
	R1                    int                   `json:"r1" bson:"r1"`
	R2                    int                   `json:"r2" bson:"r2"`
	R3                    int                   `json:"r3" bson:"r3"`
	R4                    int                   `json:"r4" bson:"r4"`
	M                     int                   `json:"m" bson:"m"`
	E                     int                   `json:"e" bson:"e"`
	RentingPrice          float64               `json:"renting_price" bson:"renting_price"`
	DistributorID         string                `json:"distributor_id" bson:"distributor_id"`
	DistributorCode       string                `json:"distributor_cdos" bson:"distributor_cdos"`
	TlgCode               MeterConfigTlgCode    `json:"tlg_code" bson:"tlg_code"`
	Type                  MeterType             `json:"type" bson:"type"`
	Meter                 Meter                 `json:"meter" bson:"meter"`
	ServicePoint          ServicePoint          `json:"service_point" bson:"service_point"`
	MeasurePoint          MeasurePoint          `json:"measure_point" bson:"measure_point"`
	ContractualSituations ContractualSituations `json:"contractual_situations" bson:"contractual_situations"`
}

type ListMeterConfigByDateQuery struct {
	Date          time.Time
	DistributorID string
	ServiceType   ServiceType
	PointType     string
	MeterType     []string
	ReadingType   ReadingType
	Limit         int
	Offset        int
}

type GetMeterConfigByMeterQuery struct {
	MeterSerialNumber string
	Date              time.Time
}

type GetMeterConfigByCupsQuery struct {
	CUPS        string
	Time        time.Time
	Distributor string
}
type ListMeterConfigByCups struct {
	DistributorId string
	StartDate     time.Time
	EndDate       time.Time
	ReadingType   ReadingType
	CUPS          string
}

type GroupMetersByTypeQuery struct {
	DistributorId string
	StartDate     time.Time
	EndDate       time.Time
}

type GetMetersAndCountByDistributorIdQuery struct {
	DistributorId string
	Type          MeterType
	StartDate     time.Time
	EndDate       time.Time
	Limit         int64
	Offset        int64
}
type GetMetersAndCountMeter struct {
	CurveType RegisterType `bson:"curve_type"`
	EndDate   time.Time    `bson:"end_date"`
	StartDate time.Time    `bson:"start_date"`
}
type GetMetersAndCountData struct {
	Cups   string                   `bson:"_id"`
	Meters []GetMetersAndCountMeter `bson:"meters"`
}
type GetMetersAndCountResult struct {
	Data  []GetMetersAndCountData `bson:"data"`
	Count int                     `bson:"total"`
}

func (m MeterConfig) GetMagnitudesActive() []Magnitude {
	magnitudes := make([]Magnitude, 0, 6)
	if m.AI > 0 {
		magnitudes = append(magnitudes, AI)
	}
	if m.AE > 0 {
		magnitudes = append(magnitudes, AE)
	}
	if m.R1 > 0 {
		magnitudes = append(magnitudes, R1)
	}
	if m.R2 > 0 {
		magnitudes = append(magnitudes, R2)
	}
	if m.R3 > 0 {
		magnitudes = append(magnitudes, R3)
	}
	if m.R4 > 0 {
		magnitudes = append(magnitudes, R4)
	}

	return magnitudes
}

func (m MeterConfig) IsCurveTypeHourly() bool {
	return Hourly == m.CurveType
}

func (m MeterConfig) Cups() string {
	return m.ServicePoint.Cups
}

func (m MeterConfig) SerialNumber() string {
	return m.Meter.SerialNumber
}

func (m MeterConfig) ServiceType() ServiceType {
	return m.ServicePoint.ServiceType
}

func (m MeterConfig) PointType() PointType {
	return m.ServicePoint.PointType
}

func (m MeterConfig) MeasurePointType() MeasurePointType {
	return m.MeasurePoint.Type
}

func (m MeterConfig) EnergyMagnitude() Magnitude {
	if m.ServiceType() == DcServiceType {
		return AI
	}
	if m.ServiceType() == GdServiceType {
		return AE
	}

	return ""
}

func (m MeterConfig) TariffID() string {
	return m.ContractualSituations.TariffID
}
func (m MeterConfig) MeterType() MeterType {
	return m.Type
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=InventoryRepository
type InventoryRepository interface {
	ListMeterConfigByDate(ctx context.Context, query ListMeterConfigByDateQuery) ([]MeterConfig, error)
	CountMeterConfigByDate(ctx context.Context, query ListMeterConfigByDateQuery) (int, error)
	GetMeterConfigByMeter(ctx context.Context, query GetMeterConfigByMeterQuery) (MeterConfig, error)
	GetMeterConfigByCups(ctx context.Context, query GetMeterConfigByCupsQuery) (MeterConfig, error)
	GetMeterConfigByCupsAPI(ctx context.Context, query GetMeterConfigByCupsQuery) (MeterConfig, error)
	ListMeterConfigByCups(ctx context.Context, query ListMeterConfigByCups) ([]MeterConfig, error)
	GroupMetersByType(ctx context.Context, query GroupMetersByTypeQuery) (map[MeterType]map[RegisterType]MeasureCount, error)
	GetMetersAndCountByDistributorId(ctx context.Context, query GetMetersAndCountByDistributorIdQuery) (GetMetersAndCountResult, error)
}
