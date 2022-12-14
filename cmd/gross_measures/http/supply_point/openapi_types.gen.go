// Package supply_point provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package supply_point

import (
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for GrossMeasureServicePointMagnitudes.
const (
	GrossMeasureServicePointMagnitudesAE GrossMeasureServicePointMagnitudes = "AE"

	GrossMeasureServicePointMagnitudesAI GrossMeasureServicePointMagnitudes = "AI"

	GrossMeasureServicePointMagnitudesR1 GrossMeasureServicePointMagnitudes = "R1"

	GrossMeasureServicePointMagnitudesR2 GrossMeasureServicePointMagnitudes = "R2"

	GrossMeasureServicePointMagnitudesR3 GrossMeasureServicePointMagnitudes = "R3"

	GrossMeasureServicePointMagnitudesR4 GrossMeasureServicePointMagnitudes = "R4"
)

// Defines values for GrossMeasureServicePointPeriods.
const (
	GrossMeasureServicePointPeriodsP1 GrossMeasureServicePointPeriods = "P1"

	GrossMeasureServicePointPeriodsP2 GrossMeasureServicePointPeriods = "P2"

	GrossMeasureServicePointPeriodsP3 GrossMeasureServicePointPeriods = "P3"

	GrossMeasureServicePointPeriodsP4 GrossMeasureServicePointPeriods = "P4"

	GrossMeasureServicePointPeriodsP5 GrossMeasureServicePointPeriods = "P5"

	GrossMeasureServicePointPeriodsP6 GrossMeasureServicePointPeriods = "P6"
)

// Defines values for GrossMeasureServicePointType.
const (
	GrossMeasureServicePointTypeOTHER GrossMeasureServicePointType = "OTHER"

	GrossMeasureServicePointTypeTLG GrossMeasureServicePointType = "TLG"

	GrossMeasureServicePointTypeTLM GrossMeasureServicePointType = "TLM"
)

// Defines values for GrossMeasureValidationStatus.
const (
	GrossMeasureValidationStatusALERT GrossMeasureValidationStatus = "ALERT"

	GrossMeasureValidationStatusINV GrossMeasureValidationStatus = "INV"

	GrossMeasureValidationStatusNONE GrossMeasureValidationStatus = "NONE"

	GrossMeasureValidationStatusSUPERV GrossMeasureValidationStatus = "SUPERV"

	GrossMeasureValidationStatusVAL GrossMeasureValidationStatus = "VAL"
)

// CurveGrossMeasureMeter defines model for CurveGrossMeasureMeter.
type CurveGrossMeasureMeter struct {
	// Embedded struct due to allOf(#/components/schemas/ServicePointCalendarStatus)
	ServicePointCalendarStatus `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	File   string             `json:"file"`
	Values GrossMeasureValues `json:"values"`
}

// Curve process measures distributor object
type CurveGrossMeasureMeterList []CurveGrossMeasureMeter

// GrossMeasureMonthlyValues defines model for GrossMeasureMonthlyValues.
type GrossMeasureMonthlyValues struct {
	// Embedded struct due to allOf(#/components/schemas/GrossMeasureValues)
	GrossMeasureValues `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	AEi float64 `json:"AEi"`
	AIi float64 `json:"AIi"`
	R1i float64 `json:"R1i"`
	R2i float64 `json:"R2i"`
	R3i float64 `json:"R3i"`
	R4i float64 `json:"R4i"`
}

// GrossMeasureServicePoint defines model for GrossMeasureServicePoint.
type GrossMeasureServicePoint struct {
	CalendarCurve          []ServicePointCalendarStatus         `json:"calendar_curve"`
	CalendarDailyClosure   []ServicePointCalendarStatus         `json:"calendar_daily_closure"`
	CalendarMonthlyClosure []ServicePointCalendarStatus         `json:"calendar_monthly_closure"`
	Cups                   string                               `json:"cups"`
	ListDailyClosures      []ServicePointDailyValues            `json:"list_daily_closures"`
	ListMonthlyClosures    []ServicePointMonthlyValues          `json:"list_monthly_closures"`
	MagnitudeEnergy        string                               `binding:"oneof='AI' 'AE'" json:"magnitude_energy"`
	Magnitudes             []GrossMeasureServicePointMagnitudes `json:"magnitudes"`
	Periods                []GrossMeasureServicePointPeriods    `json:"periods"`
	SerialNumber           string                               `json:"serial_number"`
	Type                   GrossMeasureServicePointType         `json:"type"`
}

// GrossMeasureServicePointMagnitudes defines model for GrossMeasureServicePoint.Magnitudes.
type GrossMeasureServicePointMagnitudes string

// GrossMeasureServicePointPeriods defines model for GrossMeasureServicePoint.Periods.
type GrossMeasureServicePointPeriods string

// GrossMeasureServicePointType defines model for GrossMeasureServicePoint.Type.
type GrossMeasureServicePointType string

// GrossMeasureValidationStatus defines model for GrossMeasureValidationStatus.
type GrossMeasureValidationStatus string

// GrossMeasureValues defines model for GrossMeasureValues.
type GrossMeasureValues struct {
	AE float64 `json:"AE"`
	AI float64 `json:"AI"`
	R1 float64 `json:"R1"`
	R2 float64 `json:"R2"`
	R3 float64 `json:"R3"`
	R4 float64 `json:"R4"`
}

// ServicePointCalendarStatus defines model for ServicePointCalendarStatus.
type ServicePointCalendarStatus struct {
	Date   string                       `json:"date"`
	Status GrossMeasureValidationStatus `json:"status"`
}

// ServicePointDailyValues defines model for ServicePointDailyValues.
type ServicePointDailyValues struct {
	P0      GrossMeasureValues           `json:"P0"`
	P1      GrossMeasureValues           `json:"P1"`
	P2      GrossMeasureValues           `json:"P2"`
	P3      GrossMeasureValues           `json:"P3"`
	P4      GrossMeasureValues           `json:"P4"`
	P5      GrossMeasureValues           `json:"P5"`
	P6      GrossMeasureValues           `json:"P6"`
	EndDate string                       `json:"end_date"`
	File    string                       `json:"file"`
	Status  GrossMeasureValidationStatus `json:"status"`
}

// ServicePointMonthlyValues defines model for ServicePointMonthlyValues.
type ServicePointMonthlyValues struct {
	P0       GrossMeasureMonthlyValues    `json:"P0"`
	P1       GrossMeasureMonthlyValues    `json:"P1"`
	P2       GrossMeasureMonthlyValues    `json:"P2"`
	P3       GrossMeasureMonthlyValues    `json:"P3"`
	P4       GrossMeasureMonthlyValues    `json:"P4"`
	P5       GrossMeasureMonthlyValues    `json:"P5"`
	P6       GrossMeasureMonthlyValues    `json:"P6"`
	EndDate  string                       `json:"end_date"`
	File     string                       `json:"file"`
	InitDate string                       `json:"init_date"`
	Status   GrossMeasureValidationStatus `json:"status"`
}

// Curve process measures distributor object
type CurveGrossMeasureMeterResponse CurveGrossMeasureMeterList

// ServicePointGrossMeasureResponse defines model for ServicePointGrossMeasureResponse.
type ServicePointGrossMeasureResponse GrossMeasureServicePoint

// GetCurveGrossMeasureMeterParams defines parameters for GetCurveGrossMeasureMeter.
type GetCurveGrossMeasureMeterParams struct {
	// The service point distributor_id whose the data will be taken
	Distributor string `json:"distributor"`

	// Date of the gross curve measures
	Date      openapi_types.Date                       `json:"date"`
	CurveType GetCurveGrossMeasureMeterParamsCurveType `json:"curve_type"`
}

// GetCurveGrossMeasureMeterParamsCurveType defines parameters for GetCurveGrossMeasureMeter.
type GetCurveGrossMeasureMeterParamsCurveType string

// GetGrossMeasureServicePointParams defines parameters for GetGrossMeasureServicePoint.
type GetGrossMeasureServicePointParams struct {
	// The distributor whose the data will be taken
	Distributor string `json:"distributor"`

	// Start date of the process measures
	StartDate openapi_types.Date `json:"start_date"`

	// End date of the process measures
	EndDate openapi_types.Date `json:"end_date"`
}
