package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

/*
recorrer por periodo

*/

type ValidatorQuarterly struct {
	ValidatorBase
}

var validQuarterlyKeys = map[string]struct{}{
	StartDate:   {},
	EndDate:     {},
	MeasureDate: {},
}

func NewValidatorQuarterly(v ValidationData, status measures.Status) (ValidatorQuarterly, error) {

	if len(v.Config) != 0 {
		return ValidatorQuarterly{}, newErrorConfigProperties()
	}

	validator := ValidatorQuarterly{
		ValidatorBase: NewValidatorBase(v, status),
	}

	if err := validator.IsValidKeys(validQuarterlyKeys); err != nil {
		return ValidatorQuarterly{}, err
	}

	return validator, nil
}

func (v ValidatorQuarterly) getValueKey(m MeasureValidatable, key string) time.Time {
	switch key {
	case StartDate:
		return m.StartDate
	case EndDate:
		return m.EndDate
	case MeasureDate:
		return m.ReadingDate
	default:
		return time.Time{}
	}
}

func (v ValidatorQuarterly) verifyQuarterly(t time.Time) bool {
	return t.UTC().Hour() == 00 && (t.UTC().Minute() == 00 || t.UTC().Minute() == 15 ||
		t.UTC().Minute() == 30 || t.UTC().Minute() == 45) && t.UTC().Second() == 00
}
func (v ValidatorQuarterly) Validate(m MeasureValidatable) *ValidatorBase {
	if m.Type != v.MeasureType {
		return nil
	}
	for _, k := range v.Keys {
		if !m.isInWhiteList(k) {
			continue
		}
		value := v.getValueKey(m, k)
		if !v.verifyQuarterly(value) {
			return &v.ValidatorBase
		}
	}
	return nil
}
