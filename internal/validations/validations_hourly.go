package validations

/*
recorrer por periodo

*/

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

type ValidatorHourly struct {
	ValidatorBase
}

var validDateHourlyKeys = map[string]struct{}{
	StartDate:   {},
	EndDate:     {},
	MeasureDate: {},
	FX:          {},
}

func NewValidatorHourly(v ValidationData, status measures.Status) (ValidatorHourly, error) {
	if len(v.Config) != 0 {
		return ValidatorHourly{}, newErrorConfigProperties()
	}

	validator := ValidatorHourly{
		ValidatorBase: NewValidatorBase(v, status),
	}

	if err := validator.IsValidKeys(validDateHourlyKeys); err != nil {
		return ValidatorHourly{}, err
	}

	return validator, nil
}

func (v ValidatorHourly) verifyHourly(t time.Time) bool {
	return t.Minute() == 0 && t.Second() == 0
}

func (v ValidatorHourly) getValueKey(m MeasureValidatable, key string) time.Time {
	switch key {
	case StartDate:
		return m.StartDate
	case EndDate:
		return m.EndDate
	case MeasureDate:
		return m.ReadingDate
	case FX:
		return m.FX
	default:
		return time.Time{}
	}
}

func (v ValidatorHourly) Validate(m MeasureValidatable) *ValidatorBase {
	if m.Type != v.MeasureType {
		return nil
	}
	for _, k := range v.Keys {
		if !m.isInWhiteList(k) {
			continue
		}
		value := v.getValueKey(m, k)
		if !v.verifyHourly(value) {
			return &v.ValidatorBase
		}
	}
	return nil
}
