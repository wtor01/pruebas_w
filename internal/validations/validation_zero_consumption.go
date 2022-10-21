package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
)

type ValidatorZeroConsumption struct {
	ValidatorBase
}

var validZeroConsumption = map[string]struct{}{
	AI: {},
}

func NewValidatorZeroConsumption(v ValidationData, status measures.Status) (ValidatorZeroConsumption, error) {
	if len(v.Config) != 0 {
		return ValidatorZeroConsumption{}, newErrorConfigProperties()
	}

	validator := ValidatorZeroConsumption{
		ValidatorBase: NewValidatorBase(v, status),
	}

	if err := validator.IsValidKeys(validZeroConsumption); err != nil {
		return ValidatorZeroConsumption{}, err
	}

	return validator, nil
}

func (v ValidatorZeroConsumption) getValueKey(m MeasureValidatable, key string) float64 {
	switch key {
	case AI:
		return m.AI
	default:
		return 0
	}
}

func (v ValidatorZeroConsumption) Validate(m MeasureValidatable) *ValidatorBase {
	if m.Type != v.MeasureType {
		return nil
	}
	if m.ServiceType != measures.DcServiceType && m.ServiceType != measures.GdServiceType {
		return nil
	}
	if m.TechnologyType != "TLG" {
		return nil
	}
	if m.PointType != "3" && m.PointType != "4" {
		return nil
	}

	if !v.comparePeriods(m.P0, m) {
		return &v.ValidatorBase
	}
	return nil
}

func (v ValidatorZeroConsumption) comparePeriods(p0 Periods, p MeasureValidatable) bool {
	if p0.AI > 0 {
		if p.P1.AI == 0 {
			return false
		}
		if p.P2.AI == 0 {
			return false
		}
		if p.P3.AI == 0 {
			return false
		}
		if p.P4.AI == 0 {
			return false
		}
		if p.P5.AI == 0 {
			return false
		}
		if p.P6.AI == 0 {
			return false
		}
	}
	if p0.AE > 0 {
		if p.P1.AE == 0 {
			return false
		}
		if p.P2.AE == 0 {
			return false
		}
		if p.P3.AE == 0 {
			return false
		}
		if p.P4.AE == 0 {
			return false
		}
		if p.P5.AE == 0 {
			return false
		}
		if p.P6.AE == 0 {
			return false
		}
	}
	return true
}
