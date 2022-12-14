// Package smarkia provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package smarkia

import (
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for RecoverSmarkiaMeasuresProcessName.
const (
	RecoverSmarkiaMeasuresProcessNameClose RecoverSmarkiaMeasuresProcessName = "close"

	RecoverSmarkiaMeasuresProcessNameCurve RecoverSmarkiaMeasuresProcessName = "curve"
)

// day_types object
type RecoverSmarkiaMeasures struct {
	Cups          string                            `json:"cups"`
	Date          openapi_types.Date                `json:"date"`
	DistributorId string                            `json:"distributor_id"`
	ProcessName   RecoverSmarkiaMeasuresProcessName `json:"process_name"`
}

// RecoverSmarkiaMeasuresProcessName defines model for RecoverSmarkiaMeasures.ProcessName.
type RecoverSmarkiaMeasuresProcessName string

// RecoverSmarkiaMeasuresJSONRequestBody defines body for RecoverSmarkiaMeasures for application/json ContentType.
type RecoverSmarkiaMeasuresJSONRequestBody RecoverSmarkiaMeasures
