// Package dashboard provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package dashboard

import (
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for ServicePointDashboardMagnitudes.
const (
	ServicePointDashboardMagnitudesAE ServicePointDashboardMagnitudes = "AE"

	ServicePointDashboardMagnitudesAI ServicePointDashboardMagnitudes = "AI"

	ServicePointDashboardMagnitudesR1 ServicePointDashboardMagnitudes = "R1"

	ServicePointDashboardMagnitudesR2 ServicePointDashboardMagnitudes = "R2"

	ServicePointDashboardMagnitudesR3 ServicePointDashboardMagnitudes = "R3"

	ServicePointDashboardMagnitudesR4 ServicePointDashboardMagnitudes = "R4"
)

// Defines values for ServicePointDashboardPeriods.
const (
	ServicePointDashboardPeriodsP1 ServicePointDashboardPeriods = "P1"

	ServicePointDashboardPeriodsP2 ServicePointDashboardPeriods = "P2"

	ServicePointDashboardPeriodsP3 ServicePointDashboardPeriods = "P3"

	ServicePointDashboardPeriodsP4 ServicePointDashboardPeriods = "P4"

	ServicePointDashboardPeriodsP5 ServicePointDashboardPeriods = "P5"

	ServicePointDashboardPeriodsP6 ServicePointDashboardPeriods = "P6"
)

// Defines values for ServicePointDashboardResponseType.
const (
	ServicePointDashboardResponseTypeOTHER ServicePointDashboardResponseType = "OTHER"

	ServicePointDashboardResponseTypeTLG ServicePointDashboardResponseType = "TLG"

	ServicePointDashboardResponseTypeTLM ServicePointDashboardResponseType = "TLM"
)

// Curve defines model for Curve.
type Curve struct {
	P1     *CurveValues `json:"P1,omitempty"`
	P2     *CurveValues `json:"P2,omitempty"`
	P3     *CurveValues `json:"P3,omitempty"`
	P4     *CurveValues `json:"P4,omitempty"`
	P5     *CurveValues `json:"P5,omitempty"`
	P6     *CurveValues `json:"P6,omitempty"`
	Status string       `json:"status"`
}

// CurveProcessServicePoint defines model for CurveProcessServicePoint.
type CurveProcessServicePoint struct {
	Date   string      `json:"date"`
	Status string      `json:"status"`
	Values CurveValues `json:"values"`
}

// CurveValues defines model for CurveValues.
type CurveValues struct {
	AE float64 `json:"AE"`
	AI float64 `json:"AI"`
	R1 float64 `json:"R1"`
	R2 float64 `json:"R2"`
	R3 float64 `json:"R3"`
	R4 float64 `json:"R4"`
}

// DailyClosure defines model for DailyClosure.
type DailyClosure struct {
	P0     *CurveValues `json:"P0,omitempty"`
	P1     *CurveValues `json:"P1,omitempty"`
	P2     *CurveValues `json:"P2,omitempty"`
	P3     *CurveValues `json:"P3,omitempty"`
	P4     *CurveValues `json:"P4,omitempty"`
	P5     *CurveValues `json:"P5,omitempty"`
	P6     *CurveValues `json:"P6,omitempty"`
	Status string       `json:"status"`
}

// DashboardCups defines model for DashboardCups.
type DashboardCups struct {
	Cups    string              `json:"cups"`
	Curve   DashboardCupsValues `json:"curve"`
	Daily   DashboardCupsValues `json:"daily"`
	Monthly DashboardCupsValues `json:"monthly"`
}

// DashboardCupsValues defines model for DashboardCupsValues.
type DashboardCupsValues struct {
	Invalid   int `json:"invalid"`
	None      int `json:"none"`
	ShouldBe  int `json:"should_be"`
	Supervise int `json:"supervise"`
	Total     int `json:"total"`
	Valid     int `json:"valid"`
}

// DashboardProcessMeasureData defines model for DashboardProcessMeasureData.
type DashboardProcessMeasureData struct {
	Invalid          int `json:"invalid"`
	MeasuresShouldBe int `json:"measures_should_be"`
	Supervise        int `json:"supervise"`
	Valid            int `json:"valid"`
}

// DashboardProcessMeasureOthers defines model for DashboardProcessMeasureOthers.
type DashboardProcessMeasureOthers struct {
	Closing DashboardProcessMeasureData `json:"closing"`
}

// DashboardProcessMeasureTLG defines model for DashboardProcessMeasureTLG.
type DashboardProcessMeasureTLG struct {
	Closing DashboardProcessMeasureData `json:"closing"`
	Curva   DashboardProcessMeasureData `json:"curva"`
	Resumen DashboardProcessMeasureData `json:"resumen"`
}

// DashboardProcessMeasureTLM defines model for DashboardProcessMeasureTLM.
type DashboardProcessMeasureTLM struct {
	Closing DashboardProcessMeasureData `json:"closing"`
	Curva   DashboardProcessMeasureData `json:"curva"`
}

// MonthlyClosure defines model for MonthlyClosure.
type MonthlyClosure struct {
	P0       *MonthlyValues `json:"P0,omitempty"`
	P1       *MonthlyValues `json:"P1,omitempty"`
	P2       *MonthlyValues `json:"P2,omitempty"`
	P3       *MonthlyValues `json:"P3,omitempty"`
	P4       *MonthlyValues `json:"P4,omitempty"`
	P5       *MonthlyValues `json:"P5,omitempty"`
	P6       *MonthlyValues `json:"P6,omitempty"`
	EndDate  string         `json:"end_date"`
	Id       string         `json:"id"`
	InitDate string         `json:"init_date"`
	Status   string         `json:"status"`
}

// MonthlyValues defines model for MonthlyValues.
type MonthlyValues struct {
	AE  float64 `json:"AE"`
	AEi float64 `json:"AEi"`
	AI  float64 `json:"AI"`
	AIi float64 `json:"AIi"`
	R1  float64 `json:"R1"`
	R1i float64 `json:"R1i"`
	R2  float64 `json:"R2"`
	R2i float64 `json:"R2i"`
	R3  float64 `json:"R3"`
	R3i float64 `json:"R3i"`
	R4  float64 `json:"R4"`
	R4i float64 `json:"R4i"`
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

// ServicePointDashboard defines model for ServicePointDashboard.
type ServicePointDashboard struct {
	Curve           Curve                             `json:"curve"`
	DailyClosure    DailyClosure                      `json:"daily_closure"`
	Date            string                            `json:"date"`
	MagnitudeEnergy string                            `binding:"oneof='AI' 'AE'" json:"magnitude_energy"`
	Magnitudes      []ServicePointDashboardMagnitudes `json:"magnitudes"`
	MonthlyClosure  MonthlyClosure                    `json:"monthly_closure"`
	Periods         []ServicePointDashboardPeriods    `json:"periods"`
}

// ServicePointDashboardMagnitudes defines model for ServicePointDashboard.Magnitudes.
type ServicePointDashboardMagnitudes string

// ServicePointDashboardPeriods defines model for ServicePointDashboard.Periods.
type ServicePointDashboardPeriods string

// Service point process measures object by type
type ServicePointDashboardResponse struct {
	// Service point process measures object
	Days *ServicePointDashboardResponseDays `json:"days,omitempty"`
	Type *ServicePointDashboardResponseType `json:"type,omitempty"`
}

// ServicePointDashboardResponseType defines model for ServicePointDashboardResponse.Type.
type ServicePointDashboardResponseType string

// Service point process measures object
type ServicePointDashboardResponseDays []ServicePointDashboard

// DashboardCupsProcessMeasure defines model for DashboardCupsProcessMeasure.
type DashboardCupsProcessMeasure struct {
	// Embedded struct due to allOf(#/components/schemas/Pagination)
	Pagination `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Results []DashboardCups `json:"results"`
}

// DashboardProcessMeasure defines model for DashboardProcessMeasure.
type DashboardProcessMeasure struct {
	Daily []struct {
		Date        openapi_types.Date            `json:"date"`
		Others      DashboardProcessMeasureOthers `json:"others"`
		Telegestion DashboardProcessMeasureTLG    `json:"telegestion"`
		Telemedida  DashboardProcessMeasureTLM    `json:"telemedida"`
	} `json:"daily"`
	Totals struct {
		Others      DashboardProcessMeasureOthers `json:"others"`
		Telegestion DashboardProcessMeasureTLG    `json:"telegestion"`
		Telemedida  DashboardProcessMeasureTLM    `json:"telemedida"`
	} `json:"totals"`
}

// Service point process measures object by type
type ServicePointProcessDashboard ServicePointDashboardResponse

// GetProcessMeasureDashboardListParams defines parameters for GetProcessMeasureDashboardList.
type GetProcessMeasureDashboardListParams struct {
	// Limit of Meters to view
	Limit int `json:"limit"`

	// Meters to skip
	Offset *int `json:"offset,omitempty"`

	// Start Date
	StartDate openapi_types.Date `json:"start_date"`

	// End Date
	EndDate openapi_types.Date `json:"end_date"`

	// Distributor Id
	DistributorId string `json:"distributor_id"`

	// Type of meter
	Type GetProcessMeasureDashboardListParamsType `json:"type"`
}

// GetProcessMeasureDashboardListParamsType defines parameters for GetProcessMeasureDashboardList.
type GetProcessMeasureDashboardListParamsType string

// GetCurveProcessServicePointParams defines parameters for GetCurveProcessServicePoint.
type GetCurveProcessServicePointParams struct {
	// The service point distributor_id whose the data will be taken
	Distributor string `json:"distributor"`

	// The service point cups whose the data will be taken
	Cups string `json:"cups"`

	// Start date of the process measures
	StartDate openapi_types.Date `json:"start_date"`

	// End date of the process measures
	EndDate   openapi_types.Date                         `json:"end_date"`
	CurveType GetCurveProcessServicePointParamsCurveType `json:"curve_type"`
}

// GetCurveProcessServicePointParamsCurveType defines parameters for GetCurveProcessServicePoint.
type GetCurveProcessServicePointParamsCurveType string

// GetProcessMeasureDashboardParams defines parameters for GetProcessMeasureDashboard.
type GetProcessMeasureDashboardParams struct {
	StartDate openapi_types.Date `json:"start_date"`
	EndDate   openapi_types.Date `json:"end_date"`

	// ditributor id
	DistributorId string `json:"distributor_id"`
}

// GetDashboardProcessServicePointParams defines parameters for GetDashboardProcessServicePoint.
type GetDashboardProcessServicePointParams struct {
	// The cups whose the data will be taken
	Cups string `json:"cups"`

	// The distributor whose the data will be taken
	Distributor string `json:"distributor"`

	// Start date of the process measures
	StartDate openapi_types.Date `json:"start_date"`

	// End date of the process measures
	EndDate openapi_types.Date `json:"end_date"`
}