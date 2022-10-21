package validations

/*

recorrer por periodo
*/

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

type ValidatorDaily struct {
	ValidatorBase
}

var validDailyKeys = map[string]struct{}{
	StartDate:   {},
	EndDate:     {},
	MeasureDate: {},
	FX:          {},
}

func NewValidatorDaily(v ValidationData, status measures.Status) (ValidatorDaily, error) {

	if len(v.Config) != 0 {
		return ValidatorDaily{}, newErrorConfigProperties()
	}

	validator := ValidatorDaily{
		ValidatorBase: NewValidatorBase(v, status),
	}

	if err := validator.IsValidKeys(validDailyKeys); err != nil {
		return ValidatorDaily{}, err
	}

	return validator, nil
}

func (v ValidatorDaily) getValueKey(m MeasureValidatable, key string) time.Time {
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

func (v ValidatorDaily) verifyDaily(t time.Time) bool {
	return t.UTC().Hour() == 00 && t.UTC().Minute() == 00 && t.UTC().Second() == 00
}

func (v ValidatorDaily) Validate(m MeasureValidatable) *ValidatorBase {
	if m.Type != v.MeasureType {
		return nil
	}
	for _, k := range v.Keys {
		if !m.isInWhiteList(k) {
			continue
		}
		value := v.getValueKey(m, k)
		if !v.verifyDaily(value) {
			return &v.ValidatorBase
		}
	}
	return nil
}
