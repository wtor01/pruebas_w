package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
)

type ValidatorSummaryCalendar struct {
	ValidatorBase
}

var validSummaryCalendar = map[string]struct{}{
	AI: {},
}

func NewValidatorSummaryCalendar(v ValidationData, status measures.Status) (ValidatorSummaryCalendar, error) {

	validator := ValidatorSummaryCalendar{
		ValidatorBase: NewValidatorBase(v, status),
	}

	if err := validator.IsValidKeys(validSummaryCalendar); err != nil {
		return ValidatorSummaryCalendar{}, err
	}

	return validator, nil
}

func (v ValidatorSummaryCalendar) getValueKey(m MeasureValidatable, key string) float64 {
	switch key {
	case AI:
		return m.AI
	default:
		return 0
	}
}

func (v ValidatorSummaryCalendar) Validate(m MeasureValidatable) *ValidatorBase {
	if m.Type != v.MeasureType {
		return nil
	}
	if !v.isElementExist(m.PeriodsNumbers, "P1") {
		if !v.emptyData(m.P1i) {
			return &v.ValidatorBase
		}
	}
	if !v.isElementExist(m.PeriodsNumbers, "P2") {
		if !v.emptyData(m.P2i) {
			return &v.ValidatorBase
		}
	}
	if !v.isElementExist(m.PeriodsNumbers, "P3") {
		if !v.emptyData(m.P3i) {
			return &v.ValidatorBase
		}
	}
	if !v.isElementExist(m.PeriodsNumbers, "P4") {
		if !v.emptyData(m.P4i) {
			return &v.ValidatorBase
		}
	}
	if !v.isElementExist(m.PeriodsNumbers, "P5") {
		if !v.emptyData(m.P5i) {
			return &v.ValidatorBase
		}
	}
	if !v.isElementExist(m.PeriodsNumbers, "P6") {
		if !v.emptyData(m.P6i) {
			return &v.ValidatorBase
		}
	}

	return nil
}
func (v ValidatorSummaryCalendar) isElementExist(s []string, str string) bool {
	for _, p := range s {
		if p == str {
			return true
		}
	}
	return false
}
func (v ValidatorSummaryCalendar) emptyData(p Periods) bool {
	if p.AI != 0 {
		return false
	}
	if p.AE != 0 {
		return false
	}
	if p.R1 != 0 {
		return false
	}
	if p.R2 != 0 {
		return false
	}
	if p.R3 != 0 {
		return false
	}
	if p.R4 != 0 {
		return false
	}
	return true
}
