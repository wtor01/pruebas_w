package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
)

/*
solo se aplican a curva

*/

type ValidatorThreshold struct {
	ValidatorBase
	Max float64
	Min float64
}

var validThresholdKeys = map[string]struct{}{
	AI: {},
	AE: {},
	R1: {},
	R2: {},
	R3: {},
	R4: {},
	MX: {},
	E:  {},
}

func NewValidatorThreshold(v ValidationData, status measures.Status) (ValidatorThreshold, error) {

	if len(v.Config) != 2 {
		return ValidatorThreshold{}, newErrorConfigProperties()
	}

	validator := ValidatorThreshold{
		ValidatorBase: NewValidatorBase(v, status),
	}

	if err := validator.IsValidKeys(validThresholdKeys); err != nil {
		return ValidatorThreshold{}, err
	}

	max, err := utils.ParseMapKey(v.Config, MaxKey, utils.ParseFloat64)

	if err != nil {
		return ValidatorThreshold{}, err
	}
	min, err := utils.ParseMapKey(v.Config, MinKey, utils.ParseFloat64)

	if err != nil {
		return ValidatorThreshold{}, err
	}
	validator.Min = min
	validator.Max = max

	if validator.Max < validator.Min {
		return ValidatorThreshold{}, newErrorRanged(validator.Min, validator.Max)
	}

	return validator, nil
}

func (v ValidatorThreshold) getValueKey(m MeasureValidatable, key string) float64 {
	switch key {
	case AI:
		return m.AI
	case AE:
		return m.AE
	case R1:
		return m.R1
	case R2:
		return m.R2
	case R3:
		return m.R3
	case R4:
		return m.R4
	case MX:
		return m.MX
	case E:
		return m.E
	default:
		return 0
	}
}

func (v ValidatorThreshold) Validate(m MeasureValidatable) *ValidatorBase {
	if m.Type != v.MeasureType {
		return nil
	}
	for _, k := range v.Keys {
		if !m.isInWhiteList(k) {
			continue
		}
		value := v.getValueKey(m, k)
		if value >= v.Max || value <= v.Min {
			return &v.ValidatorBase
		}
	}
	return nil
}
