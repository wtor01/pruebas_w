package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
)

type ValidatorCloseMeter struct {
	ValidatorBase
}

var validCloseMeter = map[string]struct{}{
	AI: {},
}

func NewValidatorCloseMeter(v ValidationData, status measures.Status) (ValidatorCloseMeter, error) {

	validator := ValidatorCloseMeter{
		ValidatorBase: NewValidatorBase(v, status),
	}

	if err := validator.IsValidKeys(validSummaryCalendar); err != nil {
		return ValidatorCloseMeter{}, err
	}

	return validator, nil
}

func (v ValidatorCloseMeter) getValueKey(m MeasureValidatable, key string) float64 {
	switch key {
	case AI:
		return m.AI
	default:
		return 0
	}
}

func (v ValidatorCloseMeter) Validate(m MeasureValidatable) *ValidatorBase {
	if !v.compareMeter(m.P0, m.P0LM) {
		return &v.ValidatorBase
	}
	if !v.compareMeter(m.P1, m.P1LM) {
		return &v.ValidatorBase
	}
	if !v.compareMeter(m.P2, m.P2LM) {
		return &v.ValidatorBase
	}
	if !v.compareMeter(m.P3, m.P3LM) {
		return &v.ValidatorBase
	}
	if !v.compareMeter(m.P4, m.P4LM) {
		return &v.ValidatorBase
	}
	if !v.compareMeter(m.P5, m.P5LM) {
		return &v.ValidatorBase
	}
	if !v.compareMeter(m.P6, m.P6LM) {
		return &v.ValidatorBase
	}

	return nil
}

func (v ValidatorCloseMeter) compareMeter(p Periods, pLM Periods) bool {
	if p.AI < pLM.AI {
		return false
	}
	if p.AE < pLM.AE {
		return false
	}
	if p.R1 < pLM.R1 {
		return false
	}
	if p.R2 < pLM.R2 {
		return false
	}
	if p.R3 < pLM.R3 {
		return false
	}
	if p.R4 < pLM.R4 {
		return false
	}

	return true
}
