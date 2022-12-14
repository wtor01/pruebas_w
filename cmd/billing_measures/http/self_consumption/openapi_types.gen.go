// Package self_consumption provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package self_consumption

import (
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// BillingCalendarConsumption defines model for BillingCalendarConsumption.
type BillingCalendarConsumption struct {
	Date   string                             `json:"date"`
	Energy float64                            `json:"energy"`
	Values []BillingCalendarConsumptionValues `json:"values"`
}

// BillingCalendarConsumptionValues defines model for BillingCalendarConsumptionValues.
type BillingCalendarConsumptionValues struct {
	Energy float64 `json:"energy"`
	Hour   string  `json:"hour"`
}

// BillingCauInfo defines model for BillingCauInfo.
type BillingCauInfo struct {
	// Cau config type
	ConfigType string `json:"config_type"`

	// Cau Id
	Id string `json:"id"`

	// Cau name
	Name string `json:"name"`

	// Cau Number Points
	Points int `json:"points"`

	// Cau Unit type
	UnitType string `json:"unit_type"`
}

// BillingNetGeneration defines model for BillingNetGeneration.
type BillingNetGeneration struct {
	Date      string  `json:"date"`
	Excedente float64 `json:"excedente"`
	Net       float64 `json:"net"`
}

// Billing Cups points list
type BillingSelfConsumptionCups struct {
	Consumption float64 `json:"consumption"`
	Cups        string  `json:"cups"`
	EndDate     string  `json:"end_date"`
	Generation  float64 `json:"generation"`
	PsType      string  `json:"ps_type"`
	StartDate   string  `json:"start_date"`
}

// BillingSelfConsumptionUnitInfo defines model for BillingSelfConsumptionUnitInfo.
type BillingSelfConsumptionUnitInfo struct {
	CalendarConsumption []BillingCalendarConsumption `json:"calendar_consumption"`
	CauInfo             BillingCauInfo               `json:"cau_info"`
	CupsList            []BillingSelfConsumptionCups `json:"cups_list"`
	EndDate             time.Time                    `json:"end_date"`
	NetGeneration       []BillingNetGeneration       `json:"net_generation"`
	StartDate           time.Time                    `json:"start_date"`
	Status              string                       `json:"status"`
	Totals              BillingTotals                `json:"totals"`
	UnitConsumption     []BillingUnitConsumption     `json:"unit_consumption"`
}

// BillingSelfConsumptionUnits defines model for BillingSelfConsumptionUnits.
type BillingSelfConsumptionUnits []BillingSelfConsumptionUnitInfo

// BillingTotals defines model for BillingTotals.
type BillingTotals struct {
	// Auxiliar services consumption
	AuxConsumption float64 `json:"aux_consumption"`

	// Gross Generation
	GrossGeneration float64 `json:"gross_generation"`

	// Net Generation
	NetGeneration float64 `json:"net_generation"`

	// Network consumption
	NetworkConsumption float64 `json:"network_consumption"`

	// Self consumption consumed
	SelfConsumption float64 `json:"self_consumption"`
}

// BillingUnitConsumption defines model for BillingUnitConsumption.
type BillingUnitConsumption struct {
	Aux     float64 `json:"aux"`
	Date    string  `json:"date"`
	Network float64 `json:"network"`
	Self    float64 `json:"self"`
}

// Pagination defines model for Pagination.
type Pagination struct {
	Links struct {
		// url for request next list
		Next *string `json:"next,omitempty"`

		// url for request previous list
		Prev *string `json:"prev,omitempty"`

		// url for request current list
		Self string `json:"self"`
	} `json:"_links"`
	Count  int  `json:"count"`
	Limit  int  `json:"limit"`
	Offset *int `json:"offset,omitempty"`
	Size   int  `json:"size"`
}

// SelfConsumptionConfig defines model for SelfConsumptionConfig.
type SelfConsumptionConfig struct {
	AntivertType          *string    `json:"antivert_type,omitempty"`
	CnmcTypeDesc          *string    `json:"cnmc_type_desc,omitempty"`
	CnmcTypeId            *int       `json:"cnmc_type_id,omitempty"`
	CnmcTypeName          *string    `json:"cnmc_type_name,omitempty"`
	Compensation          *bool      `json:"compensation,omitempty"`
	ConfType              *string    `json:"conf_type,omitempty"`
	ConfTypeDescription   *string    `json:"conf_type_description,omitempty"`
	ConnType              *string    `json:"conn_type,omitempty"`
	ConsumerType          *string    `json:"consumer_type,omitempty"`
	EndDate               *time.Time `json:"end_date,omitempty"`
	Excedents             *bool      `json:"excedents,omitempty"`
	GenerationPot         *float32   `json:"generation_pot,omitempty"`
	GroupSubgroup         *int       `json:"group_subgroup,omitempty"`
	Id                    *string    `json:"id,omitempty"`
	InitDate              *time.Time `json:"init_date,omitempty"`
	ParticipantNumber     *int       `json:"participant_number,omitempty"`
	SolarZoneId           *int       `json:"solar_zone_id,omitempty"`
	SolarZoneName         *string    `json:"solar_zone_name,omitempty"`
	SolarZoneNum          *int       `json:"solar_zone_num,omitempty"`
	StatusId              *int       `json:"status_id,omitempty"`
	StatusName            *string    `json:"status_name,omitempty"`
	TechnologyDescription *string    `json:"technology_description,omitempty"`
	TechnologyId          *string    `json:"technology_id,omitempty"`
}

// SelfConsumptionPoint defines model for SelfConsumptionPoint.
type SelfConsumptionPoint struct {
	CUPS             *string    `json:"CUPS,omitempty"`
	EndDate          *time.Time `json:"end_date,omitempty"`
	Exent1Flag       *int       `json:"exent1_flag,omitempty"`
	Exent2Flag       *int       `json:"exent2_flag,omitempty"`
	Id               *string    `json:"id,omitempty"`
	InitDate         *time.Time `json:"init_date,omitempty"`
	InstalationFlag  *int       `json:"instalation_flag,omitempty"`
	PartitionCoeff   *float32   `json:"partition_coeff,omitempty"`
	ServicePointType *string    `json:"service_point_type,omitempty"`
	WithoutmeterFlag *int       `json:"withoutmeter_flag,omitempty"`
}

// Self-consumption unit config
type SelfConsumptionUnitConfig struct {
	CAU           string                  `json:"CAU"`
	Ccaa          string                  `json:"ccaa"`
	CcaaId        int                     `json:"ccaa_id"`
	Configs       []SelfConsumptionConfig `json:"configs"`
	DistributorId string                  `json:"distributor_id"`
	EndDate       time.Time               `json:"end_date"`
	Id            string                  `json:"id"`
	InitDate      time.Time               `json:"init_date"`
	Name          string                  `json:"name"`
	Points        []SelfConsumptionPoint  `json:"points"`
	StatusId      int                     `json:"status_id"`
	StatusName    string                  `json:"status_name"`
}

// BillingSelfConsumptionUnitResponse defines model for BillingSelfConsumptionUnitResponse.
type BillingSelfConsumptionUnitResponse BillingSelfConsumptionUnits

// ListSelfConsumptionUnitConfigs defines model for ListSelfConsumptionUnitConfigs.
type ListSelfConsumptionUnitConfigs struct {
	// Embedded struct due to allOf(#/components/schemas/Pagination)
	Pagination `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Results []SelfConsumptionUnitConfig `json:"results"`
}

// SelfConsumptionUnitConfigResponse defines model for SelfConsumptionUnitConfigResponse.
type SelfConsumptionUnitConfigResponse []SelfConsumptionUnitConfig

// SearchActivesSelfConsumptionUnitConfigByDistributorIdParams defines parameters for SearchActivesSelfConsumptionUnitConfigByDistributorId.
type SearchActivesSelfConsumptionUnitConfigByDistributorIdParams struct {
	// The number of items to skip before starting to collect the result set
	Offset *int `json:"offset,omitempty"`

	// The numbers of items to return
	Limit int `json:"limit"`

	// Date of configuration
	Date openapi_types.Date `json:"date"`
}

// SearchSelfConsumptionUnitConfigParams defines parameters for SearchSelfConsumptionUnitConfig.
type SearchSelfConsumptionUnitConfigParams struct {
	// The distributor cups whose the data will be taken
	Cups string `json:"cups"`

	// Date of configuration
	Date openapi_types.Date `json:"date"`
}

// GetSelfConsumptionByCauParams defines parameters for GetSelfConsumptionByCau.
type GetSelfConsumptionByCauParams struct {
	StartDate openapi_types.Date `json:"start_date"`
	EndDate   openapi_types.Date `json:"end_date"`
}
