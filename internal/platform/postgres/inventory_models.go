package postgres

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Distributor struct {
	ModelEntity
	Name              string `gorm:"column:name;index"`
	R1                string `gorm:"column:r1"`
	SmarkiaId         string `gorm:"column:smarkia_id"`
	IsSmarkiaActive   bool   `gorm:"column:is_smarkia_active;default:false;"`
	CreatedBy         string
	MeasureEquipments []Meter `gorm:"foreignKey:DistributorID"`
	Roles             []*Role `gorm:"many2many:user_distributor_roles;"`
	Users             []*User `gorm:"many2many:user_distributor_roles;"`
	CDOS              string  `gorm:"index;unique;column:cdos"`
	UpdatedByID       *string `gorm:"column:updated_by;type:uuid"  gorm:"foreignKey:UpdatedByID"`
	UpdatedBy         User    `gorm:"column:updated_by"`
}

type Meter struct {
	ModelEntity
	SmarkiaId         string      `gorm:"column:smarkia_id"  gorm:"uniqueIndex"`
	SerialNumber      string      `gorm:"column:serial_number" gorm:"uniqueIndex,not null"`
	Technology        string      `gorm:"column:technology;type:meters_technology" gorm:"not null"`
	Type              string      `gorm:"column:type;type:meters_type" gorm:"not null"`
	Brand             string      `gorm:"column:brand"`
	Model             string      `gorm:"column:model"`
	ActiveConstant    float64     `gorm:"column:active_constant" gorm:"default:1.0"`
	ReactiveConstant  float64     `gorm:"column:reactive_constant" gorm:"default:1.0"`
	MaximeterConstant float64     `gorm:"column:maximeter_constant" gorm:"default:1.0"`
	DistributorID     uuid.UUID   `gorm:"type:uuid;"`
	Distributor       Distributor `gorm:"distributor" gorm:"foreignKey:DistributorID"`
	ModelRegisterUserActions
}

type MeterConfig struct {
	ModelEntity
	StartDate time.Time `gorm:"column:start_date"`
	EndDate   time.Time `gorm:"column:end_date"`
	AI        int       `gorm:"column:ai"`
	AE        int       `gorm:"column:ae"`
	R1        int       `gorm:"column:r1"`
	R2        int       `gorm:"column:r2"`
	R3        int       `gorm:"column:r3"`
	R4        int       `gorm:"column:r4"`
	M         int       `gorm:"column:m"`
	E         int       `gorm:"column:e"`

	CurveType   string `gorm:"column:curve_type;type:meter_configs_curve_type"`
	ReadingType string `gorm:"column:reading_type;type:meter_configs_reading_type"`
	CalendarID  string `gorm:"column:calendar_id"`

	PriorityContract string `gorm:"column:priority_contract;type:meter_configs_contract_number" gorm:"not null"`

	MeterID string `gorm:"type:uuid"`
	Meter   Meter  `gorm:"meter" gorm:"foreignKey:MeterSerialNumber"`

	MeasurePointID string       `gorm:"type:uuid"`
	MeasurePoint   MeasurePoint `gorm:"measure_point_id" gorm:"foreignKey:MeasurePointID"`

	ServicePointID string       `gorm:"type:uuid"`
	ServicePoint   ServicePoint `gorm:"service_point_id" gorm:"foreignKey:ServicePointID"`

	TlgCode          string  `gorm:"column:tlg_code;type:meter_configs_tlg_type"`
	MeasurePointType string  `gorm:"column:measure_point_type;type:measure_points_type"`
	Type             string  `gorm:"column:tlg_type;type:meters_type"`
	RentingPrice     float64 `gorm:"column:renting_price"`

	SerialNumber  string    `gorm:"column:serial_number" gorm:"uniqueIndex,not null"`
	DistributorID uuid.UUID `gorm:"type:uuid;"`
	CUPS          string    `gorm:"cups" gorm:"index"`

	ModelRegisterUserActions
}

type MeasurePoint struct {
	ModelEntity
	Type       string  `gorm:"column:type;not null;type:measure_points_type"`
	LossesPerc float64 `gorm:"column:losses_perc;"`
	LossesCoef float64 `gorm:"column:losses_coef;"`

	ServicePoint   ServicePoint
	ServicePointID string `gorm:"type:uuid;"`
	ModelRegisterUserActions
}

type ServicePoint struct {
	ModelEntity
	Name                string         `gorm:"name" gorm:"index"`
	CUPS                string         `gorm:"cups" gorm:"index"`
	StartDate           time.Time      `gorm:"start_date"`
	ServiceType         string         `gorm:"service_type;not null;type:service_point_service_type"`
	PointType           string         `gorm:"point_type;not null;type:service_point_point_type;"`
	PointTensionLevel   string         `gorm:"point_tension_level;type:service_point_point_tension_level"`
	MeasureTensionLevel string         `gorm:"measure_tension_level;type:service_point_measure_tension_level"`
	TensionSection      int            `gorm:"tension_section;type:integer;"`
	TransfNomDemand     int            `gorm:"transf_nom_demand"`
	Disabled            bool           `gorm:"disabled"`
	SolarZoneID         string         `gorm:"solar_zone_id;type:uuid;"`
	SolarZoneGroupID    string         `gorm:"solar_zone_group_id;type:uuid;"`
	TechnologyID        string         `gorm:"technology_id;type:uuid;"`
	DistributorID       uuid.UUID      `gorm:"type:string;type:uuid;"`
	Distributor         Distributor    `gorm:"distributor" gorm:"foreignKey:DistributorID"`
	MeasurePoint        []MeasurePoint `gorm:"foreignKey:ServicePointID"`
	ModelRegisterUserActions
}

type ContractualSituation struct {
	ModelEntity
	TariffID       string       `gorm:"column:tariff_id"`
	Tariff         Tariff       `gorm:"foreignKey:TariffID;references:ID"`
	P1Demand       float64      `gorm:"column:p1_demand"`
	P2Demand       float64      `gorm:"column:p2_demand"`
	P3Demand       float64      `gorm:"column:p3_demand"`
	P4Demand       float64      `gorm:"column:p4_demand"`
	P5Demand       float64      `gorm:"column:p5_demand"`
	P6Demand       float64      `gorm:"column:p6_demand"`
	RetailerCode   int          `gorm:"column:retailer_code"`
	InitDate       time.Time    `gorm:"column:init_date"`
	EndDate        time.Time    `gorm:"column:end_date"`
	ServicePointID string       `gorm:"column:service_point_id"`
	ServicePoint   ServicePoint `gorm:"foreignKey:ServicePointID;references:ID"`
	ModelRegisterUserActions
}
