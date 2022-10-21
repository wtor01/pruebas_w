// Package validations provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package validations

import (
	"encoding/json"
	"fmt"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for MeasureForValidateStatus.
const (
	MeasureForValidateStatusINV MeasureForValidateStatus = "INV"

	MeasureForValidateStatusVAL MeasureForValidateStatus = "VAL"
)

// Defines values for ParamsType.
const (
	ParamsTypeSimple ParamsType = "simple"
)

// Defines values for ValidationKeys.
const (
	ValidationKeysAE ValidationKeys = "AE"

	ValidationKeysAI ValidationKeys = "AI"

	ValidationKeysE ValidationKeys = "E"

	ValidationKeysEndDate ValidationKeys = "end_date"

	ValidationKeysFX ValidationKeys = "FX"

	ValidationKeysMX ValidationKeys = "MX"

	ValidationKeysMeasureDate ValidationKeys = "measure_date"

	ValidationKeysQualifier ValidationKeys = "qualifier"

	ValidationKeysR1 ValidationKeys = "R1"

	ValidationKeysR2 ValidationKeys = "R2"

	ValidationKeysR3 ValidationKeys = "R3"

	ValidationKeysR4 ValidationKeys = "R4"

	ValidationKeysStartDate ValidationKeys = "start_date"
)

// Defines values for ValidationType.
const (
	ValidationTypeCloseHose ValidationType = "close_hose"

	ValidationTypeCloseMeter ValidationType = "close_meter"

	ValidationTypeCurveLimit ValidationType = "curve_limit"

	ValidationTypeDailyDate ValidationType = "daily_date"

	ValidationTypeExcesiveConsumption ValidationType = "excesive_consumption"

	ValidationTypeFutureDate ValidationType = "future_date"

	ValidationTypeHourDate ValidationType = "hour_date"

	ValidationTypeQualifier ValidationType = "qualifier"

	ValidationTypeQuarterlyDate ValidationType = "quarterly_date"

	ValidationTypeSummaryCalendar ValidationType = "summary_calendar"

	ValidationTypeSummaryTotalizer ValidationType = "summary_totalizer"

	ValidationTypeThreshold ValidationType = "threshold"

	ValidationTypeZeroConsumption ValidationType = "zero_consumption"
)

// Defines values for ValidationMeasureBaseAction.
const (
	ValidationMeasureBaseActionALERT ValidationMeasureBaseAction = "ALERT"

	ValidationMeasureBaseActionINV ValidationMeasureBaseAction = "INV"

	ValidationMeasureBaseActionNONE ValidationMeasureBaseAction = "NONE"

	ValidationMeasureBaseActionSUPERV ValidationMeasureBaseAction = "SUPERV"
)

// Defines values for ValidationMeasureBaseMeasureType.
const (
	ValidationMeasureBaseMeasureTypeABS ValidationMeasureBaseMeasureType = "ABS"

	ValidationMeasureBaseMeasureTypeABSCLO ValidationMeasureBaseMeasureType = "ABS_CLO"

	ValidationMeasureBaseMeasureTypeINC ValidationMeasureBaseMeasureType = "INC"

	ValidationMeasureBaseMeasureTypeINCCLO ValidationMeasureBaseMeasureType = "INC_CLO"
)

// Defines values for ValidationMeasureBaseType.
const (
	ValidationMeasureBaseTypeCOHE ValidationMeasureBaseType = "COHE"

	ValidationMeasureBaseTypeINM ValidationMeasureBaseType = "INM"

	ValidationMeasureBaseTypePROC ValidationMeasureBaseType = "PROC"
)

// Defines values for ValidationMeasureConfigBaseAction.
const (
	ValidationMeasureConfigBaseActionALERT ValidationMeasureConfigBaseAction = "ALERT"

	ValidationMeasureConfigBaseActionINV ValidationMeasureConfigBaseAction = "INV"

	ValidationMeasureConfigBaseActionNONE ValidationMeasureConfigBaseAction = "NONE"

	ValidationMeasureConfigBaseActionSUPERV ValidationMeasureConfigBaseAction = "SUPERV"
)

// Defines values for ValidationMeasureConfigCreateAction.
const (
	ValidationMeasureConfigCreateActionALERT ValidationMeasureConfigCreateAction = "ALERT"

	ValidationMeasureConfigCreateActionINV ValidationMeasureConfigCreateAction = "INV"

	ValidationMeasureConfigCreateActionNONE ValidationMeasureConfigCreateAction = "NONE"

	ValidationMeasureConfigCreateActionSUPERV ValidationMeasureConfigCreateAction = "SUPERV"
)

// ExtraConfig defines model for ExtraConfig.
type ExtraConfig struct {
	Id                   string            `json:"id"`
	Name                 string            `json:"name"`
	AdditionalProperties map[string]string `json:"-"`
}

// Measure for validation object
type MeasureForValidate struct {
	ID               string                   `json:"ID"`
	InvalidationCode string                   `json:"invalidation_code"`
	Status           MeasureForValidateStatus `json:"status"`
}

// MeasureForValidateStatus defines model for MeasureForValidate.Status.
type MeasureForValidateStatus string

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

// Params defines model for Params.
type Params struct {
	Type        ParamsType   `binding:"oneof='simple'" json:"type"`
	Validations []Validation `binding:"required,dive" json:"validations"`
}

// ParamsType defines model for Params.Type.
type ParamsType string

// Validation defines model for Validation.
type Validation struct {
	Config   *Validation_Config `json:"config,omitempty"`
	Id       string             `json:"id"`
	Keys     []ValidationKeys   `binding:"dive,oneof='start_date' 'end_date' 'measure_date' 'qualifier' 'AI' 'AE' 'R1' 'R2' 'R3' 'R4' 'MX' 'FX' 'E'" json:"keys"`
	Name     string             `json:"name"`
	Required bool               `json:"required"`
	Type     ValidationType     `binding:"oneof='threshold' 'qualifier' 'daily_date' 'hour_date' 'quarterly_date' 'future_date' 'curve_limit' 'summary_totalizer' 'summary_calendar' 'zero_consumption' 'close_hose' 'excesive_consumption' 'close_meter'" json:"type"`
}

// Validation_Config defines model for Validation.Config.
type Validation_Config struct {
	AdditionalProperties map[string]string `json:"-"`
}

// ValidationKeys defines model for Validation.Keys.
type ValidationKeys string

// ValidationType defines model for Validation.Type.
type ValidationType string

// ValidationMeasure defines model for ValidationMeasure.
type ValidationMeasure struct {
	// Embedded fields due to inline allOf schema
	Id string `json:"id"`
	// Embedded struct due to allOf(#/components/schemas/ValidationMeasureBase)
	ValidationMeasureBase `yaml:",inline"`
}

// Validation measure
type ValidationMeasureBase struct {
	Action      ValidationMeasureBaseAction      `binding:"oneof='INV' 'SUPERV' 'ALERT' 'NONE'" json:"action"`
	Code        string                           `json:"code"`
	Description *string                          `json:"description,omitempty"`
	Enabled     bool                             `json:"enabled"`
	MeasureType ValidationMeasureBaseMeasureType `binding:"oneof='INC' 'ABS' 'INC_CLO' 'ABS_CLO'" json:"measure_type"`
	Message     string                           `json:"message"`
	Name        string                           `json:"name"`
	Params      Params                           `json:"params"`
	Type        ValidationMeasureBaseType        `binding:"oneof='INM' 'PROC' 'COHE'" json:"type"`
}

// ValidationMeasureBaseAction defines model for ValidationMeasureBase.Action.
type ValidationMeasureBaseAction string

// ValidationMeasureBaseMeasureType defines model for ValidationMeasureBase.MeasureType.
type ValidationMeasureBaseMeasureType string

// ValidationMeasureBaseType defines model for ValidationMeasureBase.Type.
type ValidationMeasureBaseType string

// ValidationMeasureConfig defines model for ValidationMeasureConfig.
type ValidationMeasureConfig struct {
	// Embedded fields due to inline allOf schema
	Id string `json:"id"`
	// Embedded struct due to allOf(#/components/schemas/ValidationMeasureConfigBase)
	ValidationMeasureConfigBase `yaml:",inline"`
}

// Validation measure config
type ValidationMeasureConfigBase struct {
	Action            ValidationMeasureConfigBaseAction `binding:"oneof='INV' 'SUPERV' 'ALERT' 'NONE'" json:"action"`
	Enabled           bool                              `json:"enabled"`
	ExtraConfig       *[]ExtraConfig                    `json:"extra_config,omitempty"`
	ValidationMeasure ValidationMeasure                 `json:"validation_measure"`
}

// ValidationMeasureConfigBaseAction defines model for ValidationMeasureConfigBase.Action.
type ValidationMeasureConfigBaseAction string

// Validation measure config create
type ValidationMeasureConfigCreate struct {
	Action              ValidationMeasureConfigCreateAction `binding:"oneof='INV' 'SUPERV' 'ALERT' 'NONE'" json:"action"`
	Enabled             bool                                `json:"enabled"`
	ExtraConfig         *[]ExtraConfig                      `json:"extra_config,omitempty"`
	ValidationMeasureId string                              `json:"validation_measure_id"`
}

// ValidationMeasureConfigCreateAction defines model for ValidationMeasureConfigCreate.Action.
type ValidationMeasureConfigCreateAction string

// CreateValidationMeasureConfig defines model for CreateValidationMeasureConfig.
type CreateValidationMeasureConfig ValidationMeasureConfig

// CreateValidationMeasures defines model for CreateValidationMeasures.
type CreateValidationMeasures ValidationMeasure

// GetValidationMeasure defines model for GetValidationMeasure.
type GetValidationMeasure ValidationMeasure

// GetValidationMeasureConfig defines model for GetValidationMeasureConfig.
type GetValidationMeasureConfig ValidationMeasureConfig

// ListValidationMeasureConfig defines model for ListValidationMeasureConfig.
type ListValidationMeasureConfig struct {
	// Embedded struct due to allOf(#/components/schemas/Pagination)
	Pagination `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Results []ValidationMeasureConfig `json:"results"`
}

// ListValidationMeasures defines model for ListValidationMeasures.
type ListValidationMeasures struct {
	// Embedded struct due to allOf(#/components/schemas/Pagination)
	Pagination `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Results []ValidationMeasure `json:"results"`
}

// Validation measure
type CreateValidationMeasure ValidationMeasureBase

// ListValidationsMeasureConfigParams defines parameters for ListValidationsMeasureConfig.
type ListValidationsMeasureConfigParams struct {
	// The type of the validations measure config
	Type *ListValidationsMeasureConfigParamsType `json:"type,omitempty"`
}

// ListValidationsMeasureConfigParamsType defines parameters for ListValidationsMeasureConfig.
type ListValidationsMeasureConfigParamsType string

// ListValidationsMeasureParams defines parameters for ListValidationsMeasure.
type ListValidationsMeasureParams struct {
	// The number of items to skip before starting to collect the result set
	Offset *int `json:"offset,omitempty"`

	// The numbers of items to return
	Limit int `json:"limit"`
}

// PutMeasurementValidationParamsMeasureType defines parameters for PutMeasurementValidation.
type PutMeasurementValidationParamsMeasureType string

// CreateValidationsMeasureConfigJSONRequestBody defines body for CreateValidationsMeasureConfig for application/json ContentType.
type CreateValidationsMeasureConfigJSONRequestBody CreateValidationMeasureConfig

// UpdateValidationsMeasureConfigJSONRequestBody defines body for UpdateValidationsMeasureConfig for application/json ContentType.
type UpdateValidationsMeasureConfigJSONRequestBody CreateValidationMeasureConfig

// CreateValidationsMeasureJSONRequestBody defines body for CreateValidationsMeasure for application/json ContentType.
type CreateValidationsMeasureJSONRequestBody CreateValidationMeasure

// UpdateValidationsMeasureJSONRequestBody defines body for UpdateValidationsMeasure for application/json ContentType.
type UpdateValidationsMeasureJSONRequestBody CreateValidationMeasure

// PutMeasurementValidationJSONRequestBody defines body for PutMeasurementValidation for application/json ContentType.
type PutMeasurementValidationJSONRequestBody MeasureForValidate

// Getter for additional properties for ExtraConfig. Returns the specified
// element and whether it was found
func (a ExtraConfig) Get(fieldName string) (value string, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for ExtraConfig
func (a *ExtraConfig) Set(fieldName string, value string) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]string)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for ExtraConfig to handle AdditionalProperties
func (a *ExtraConfig) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["id"]; found {
		err = json.Unmarshal(raw, &a.Id)
		if err != nil {
			return fmt.Errorf("error reading 'id': %w", err)
		}
		delete(object, "id")
	}

	if raw, found := object["name"]; found {
		err = json.Unmarshal(raw, &a.Name)
		if err != nil {
			return fmt.Errorf("error reading 'name': %w", err)
		}
		delete(object, "name")
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for ExtraConfig to handle AdditionalProperties
func (a ExtraConfig) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	object["id"], err = json.Marshal(a.Id)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'id': %w", err)
	}

	object["name"], err = json.Marshal(a.Name)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'name': %w", err)
	}

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for Validation_Config. Returns the specified
// element and whether it was found
func (a Validation_Config) Get(fieldName string) (value string, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for Validation_Config
func (a *Validation_Config) Set(fieldName string, value string) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]string)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for Validation_Config to handle AdditionalProperties
func (a *Validation_Config) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for Validation_Config to handle AdditionalProperties
func (a Validation_Config) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}
