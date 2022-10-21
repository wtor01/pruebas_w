package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"errors"
)

// type
const (
	Immediate string = "INM"
	Process   string = "PROC"
	Coherence string = "COHE"
)

var ValidTypesValidationMeasureParamsValidationsType = map[string]struct{}{
	Threshold:     {},
	Qualifier:     {},
	DailyDate:     {},
	HourDate:      {},
	QuarterlyDate: {},
	FutureDate:    {},
}

const (
	Simple string = "simple"
)

var ValidTypesValidationMeasureParamsType = map[string]struct{}{
	Simple: {},
}

var ValidTypesKeysValidate = map[string]struct{}{
	StartDate:   {},
	EndDate:     {},
	MeasureDate: {},
	Qualifier:   {},
	AI:          {},
	AE:          {},
	R1:          {},
	R2:          {},
	R3:          {},
	R4:          {},
	MX:          {},
	FX:          {},
	E:           {},
}

func isValidType(validTypes map[string]struct{}, str string) bool {
	_, ok := validTypes[str]

	return ok
}

type ValidationMeasureParamsValidations struct {
	Id       string            `json:"id"`
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Keys     []string          `json:"keys"`
	Required bool              `json:"required"`
	Config   map[string]string `json:"config"`
}

func NewValidationMeasureParamsValidations(id string, name string,
	Type string, keys []string, required bool, config map[string]string) (ValidationMeasureParamsValidations, error) {
	if !isValidType(ValidTypesValidationMeasureParamsValidationsType, Type) {
		return ValidationMeasureParamsValidations{}, errors.New("invalid type")
	}
	if len(keys) == 0 {
		return ValidationMeasureParamsValidations{}, errors.New("invalid key")
	}
	for _, k := range keys {
		if !isValidType(ValidTypesKeysValidate, k) {
			return ValidationMeasureParamsValidations{}, errors.New("invalid key")
		}

	}

	_, err := NewValidator(
		ValidationData{
			Type:   Type,
			Keys:   keys,
			Config: config,
		},
		measures.Valid,
		nil,
	)

	if err != nil {
		return ValidationMeasureParamsValidations{}, err
	}

	return ValidationMeasureParamsValidations{Id: id, Name: name, Type: Type, Keys: keys, Required: required, Config: config}, nil
}

type ValidationMeasureParams struct {
	Type        string                               `json:"type"`
	Validations []ValidationMeasureParamsValidations `json:"validations"`
}

func NewValidationMeasureParams(Type string, validations []Validation) (ValidationMeasureParams, error) {

	if !isValidType(ValidTypesValidationMeasureParamsType, Type) {
		return ValidationMeasureParams{}, errors.New("invalid type")
	}

	vvv := make([]ValidationMeasureParamsValidations, 0)

	for _, v := range validations {
		vv, err := NewValidationMeasureParamsValidations(v.Id, v.Name, v.Type, v.Keys, v.Required, v.Config)

		if err != nil {
			return ValidationMeasureParams{}, err
		}

		vvv = append(vvv, vv)

	}

	return ValidationMeasureParams{Type: Type, Validations: vvv}, nil
}

type ValidationMeasure struct {
	Id          string
	Name        string
	Action      string
	Enabled     bool
	MeasureType string
	Type        string
	Code        string
	Message     string
	Description string
	Params      ValidationMeasureParams
	CreatedByID string
	UpdatedByID string
}

type Validation struct {
	Id       string
	Name     string
	Type     string
	Keys     []string
	Required bool
	Config   map[string]string
}

type Params struct {
	Type        string
	Validations []Validation
}

func NewValidationMeasure(
	Id string,
	userId string,
	Name string,
	Action string,
	Enabled bool,
	MeasureType string,
	Type string,
	Code string,
	Message string,
	Description string,
	params Params,
) (ValidationMeasure, error) {

	p, err := NewValidationMeasureParams(params.Type, params.Validations)

	if err != nil {
		return ValidationMeasure{}, err
	}

	return ValidationMeasure{
		Id:          Id,
		CreatedByID: userId,
		Name:        Name,
		Action:      Action,
		Enabled:     Enabled,
		MeasureType: MeasureType,
		Type:        Type,
		Code:        Code,
		Message:     Message,
		Description: Description,
		Params:      p,
	}, nil
}

func (v ValidationMeasure) UpdateValidationMeasure(
	userId string,
	Name string,
	Action string,
	Enabled bool,
	MeasureType string,
	Type string,
	Code string,
	Message string,
	Description string,
	params Params,
) (ValidationMeasure, error) {

	p, err := NewValidationMeasureParams(params.Type, params.Validations)

	if err != nil {
		return ValidationMeasure{}, err
	}

	return ValidationMeasure{
		Id:          v.Id,
		CreatedByID: v.CreatedByID,
		UpdatedByID: userId,
		Name:        Name,
		Action:      Action,
		Enabled:     Enabled,
		MeasureType: MeasureType,
		Type:        Type,
		Code:        Code,
		Message:     Message,
		Description: Description,
		Params:      p,
	}, nil
}
