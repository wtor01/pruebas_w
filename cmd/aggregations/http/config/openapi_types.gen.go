// Package config provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package config

import (
	"time"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// AggregationConfig defines model for AggregationConfig.
type AggregationConfig struct {
	// Embedded fields due to inline allOf schema
	Id string `json:"id"`
	// Embedded struct due to allOf(#/components/schemas/AggregationConfigBase)
	AggregationConfigBase `yaml:",inline"`
}

// Aggregation configuration
type AggregationConfigBase struct {
	Description *string              `json:"description,omitempty"`
	EndDate     *time.Time           `json:"end_date"`
	Features    []AggregationFeature `json:"features"`
	Name        string               `json:"name"`
	Scheduler   string               `json:"scheduler"`
	StartDate   time.Time            `json:"start_date"`
}

// AggregationFeature defines model for AggregationFeature.
type AggregationFeature struct {
	// Embedded fields due to inline allOf schema
	Id string `json:"id"`
	// Embedded struct due to allOf(#/components/schemas/AggregationFeaturesBase)
	AggregationFeaturesBase `yaml:",inline"`
}

// Aggregation features
type AggregationFeaturesBase struct {
	Field string `json:"field"`
	Name  string `json:"name"`
}

// AggregationsConfig defines model for AggregationsConfig.
type AggregationsConfig struct {
	// Embedded struct due to allOf(#/components/schemas/Pagination)
	Pagination `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Results []AggregationConfig `json:"results"`
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

// CreateAggregationConfig defines model for CreateAggregationConfig.
type CreateAggregationConfig AggregationConfig

// GetAggregationConfig defines model for GetAggregationConfig.
type GetAggregationConfig AggregationConfig

// GetAggregationsConfig defines model for GetAggregationsConfig.
type GetAggregationsConfig AggregationsConfig

// UpdateAggregationConfig defines model for UpdateAggregationConfig.
type UpdateAggregationConfig AggregationConfig

// GetAllAggregationsConfigParams defines parameters for GetAllAggregationsConfig.
type GetAllAggregationsConfigParams struct {
	// Filter by query name
	Q *string `json:"q,omitempty"`

	// The number of items to skip before starting to collect the result set
	Offset *int `json:"offset,omitempty"`

	// The numbers of items to return
	Limit int `json:"limit"`
}

// CreateAggregationConfigJSONRequestBody defines body for CreateAggregationConfig for application/json ContentType.
type CreateAggregationConfigJSONRequestBody AggregationConfig

// UpdateAggregationConfigJSONRequestBody defines body for UpdateAggregationConfig for application/json ContentType.
type UpdateAggregationConfigJSONRequestBody AggregationConfig
