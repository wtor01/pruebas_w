package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
)

type ValidatorExcessiveConsumption struct {
	ValidatorBase
}

var validExcessiveConsumption = map[string]struct{}{
	AI: {},
	AE: {},
}

func NewValidatorExcessiveConsumption(v ValidationData, status measures.Status) (ValidatorExcessiveConsumption, error) {
	if len(v.Config) != 0 {
		return ValidatorExcessiveConsumption{}, newErrorConfigProperties()
	}

	validator := ValidatorExcessiveConsumption{
		ValidatorBase: NewValidatorBase(v, status),
	}

	if err := validator.IsValidKeys(validExcessiveConsumption); err != nil {
		return ValidatorExcessiveConsumption{}, err
	}

	return validator, nil

}

func (v ValidatorExcessiveConsumption) getValueKey(m MeasureValidatable, key string) float64 {
	switch key {
	case AI:
		return m.AI
	case AE:
		return m.AE
	default:
		return 0
	}
}

func (v ValidatorExcessiveConsumption) Validate(m MeasureValidatable) *ValidatorBase {

	if true { //FALTA INDICAR SI HAY O NO MAXÍMETRO
		switch m.Period {
		case measures.P1:
			if !v.comparePower(m.RegisterType, float64(m.P1Demand), m.P1) {
				return &v.ValidatorBase
			}
		case measures.P2:
			if !v.comparePower(m.RegisterType, float64(m.P2Demand), m.P2) {
				return &v.ValidatorBase
			}
		case measures.P3:
			if !v.comparePower(m.RegisterType, float64(m.P3Demand), m.P3) {
				return &v.ValidatorBase
			}
		case measures.P4:
			if !v.comparePower(m.RegisterType, float64(m.P4Demand), m.P4) {
				return &v.ValidatorBase
			}
		case measures.P5:
			if !v.comparePower(m.RegisterType, float64(m.P5Demand), m.P5) {
				return &v.ValidatorBase
			}
		case measures.P6:
			if !v.comparePower(m.RegisterType, float64(m.P6Demand), m.P6) {
				return &v.ValidatorBase
			}

		}
	}
	return nil
}

func (v ValidatorExcessiveConsumption) comparePower(registerType measures.RegisterType, demand float64, p Periods) bool {
	if registerType == measures.Hourly {
		if demand >= p.AI && demand >= p.AE {
			return true
		}
	}
	if registerType == measures.QuarterHour {
		if (p.AI*4) <= demand && (p.AE*4) <= demand {
			return true
		}
	}
	return false
}
