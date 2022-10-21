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

// Data defines model for Data.
type Data struct {
	Invalid          int `json:"invalid"`
	MeasuresShouldBe int `json:"measures_should_be"`
	Supervise        int `json:"supervise"`
	Valid            int `json:"valid"`
}

// Others defines model for Others.
type Others struct {
	Closing Data `json:"closing"`
}

// Telegestion defines model for Telegestion.
type Telegestion struct {
	Closing Data `json:"closing"`
	Curva   Data `json:"curva"`
	Resumen Data `json:"resumen"`
}

// Telemedida defines model for Telemedida.
type Telemedida struct {
	Closing Data `json:"closing"`
	Curva   Data `json:"curva"`
}

// DashboardMeasure defines model for DashboardMeasure.
type DashboardMeasure struct {
	Daily []struct {
		Date        openapi_types.Date `json:"date"`
		Others      Others             `json:"others"`
		Telegestion Telegestion        `json:"telegestion"`
		Telemedida  Telemedida         `json:"telemedida"`
	} `json:"daily"`
	Totals struct {
		Others      Others      `json:"others"`
		Telegestion Telegestion `json:"telegestion"`
		Telemedida  Telemedida  `json:"telemedida"`
	} `json:"totals"`
}

// GetMeasureDashboardParams defines parameters for GetMeasureDashboard.
type GetMeasureDashboardParams struct {
	StartDate openapi_types.Date `json:"start_date"`
	EndDate   openapi_types.Date `json:"end_date"`

	// ditributor id
	DistributorId string `json:"distributor_id"`
}