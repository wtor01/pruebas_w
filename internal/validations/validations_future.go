package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

/*
recorrer por periodo

*/

type ValidatorFuture struct {
	ValidatorBase
}

var validDateFutureKeys = map[string]struct{}{
	StartDate:   {},
	EndDate:     {},
	MeasureDate: {},
	FX:          {},
}

func NewValidatorFuture(v ValidationData, status measures.Status) (ValidatorFuture, error) {
	if len(v.Config) != 0 {
		return ValidatorFuture{}, newErrorConfigProperties()
	}

	validator := ValidatorFuture{
		ValidatorBase: NewValidatorBase(v, status),
	}

	if err := validator.IsValidKeys(validDateFutureKeys); err != nil {
		return ValidatorFuture{}, err
	}

	return validator, nil
}

func (v ValidatorFuture) getValueKey(m MeasureValidatable, key string) time.Time {
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

func (v ValidatorFuture) verifyFuture(t time.Time) bool {
	return t.Unix() < time.Now().Unix()
}

func (v ValidatorFuture) Validate(m MeasureValidatable) *ValidatorBase {
	if m.Type != v.MeasureType {
		return nil
	}
	for _, k := range v.Keys {
		if !m.isInWhiteList(k) {
			continue
		}
		value := v.getValueKey(m, k)
		if !v.verifyFuture(value) {
			return &v.ValidatorBase
		}
	}
	return nil
}
