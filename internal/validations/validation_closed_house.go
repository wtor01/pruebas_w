package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
)

type ValidatorClosedHouse struct {
	ValidatorBase
}

var validClosedHouse = map[string]struct{}{
	AI: {},
	AE: {},
	R1: {},
	R2: {},
	R3: {},
	R4: {},
}

func NewValidatorClosedHouse(v ValidationData, status measures.Status) (ValidatorClosedHouse, error) {
	if len(v.Config) != 0 {
		return ValidatorClosedHouse{}, newErrorConfigProperties()
	}

	validator := ValidatorClosedHouse{
		ValidatorBase: NewValidatorBase(v, status),
	}

	if err := validator.IsValidKeys(validClosedHouse); err != nil {
		return ValidatorClosedHouse{}, err
	}

	return validator, nil
}

func (v ValidatorClosedHouse) getValueKey(m MeasureValidatable, key string) float64 {
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
	default:
		return 0
	}
}

func (v ValidatorClosedHouse) Validate(m MeasureValidatable) *ValidatorBase {
	if m.TechnologyType != "INACCESIBLE" {
		return nil
	}
	if !v.noPowerInClosedHouse(m.P0, m) {
		return &v.ValidatorBase
	}
	return nil
}

func (v ValidatorClosedHouse) noPowerInClosedHouse(p0 Periods, m MeasureValidatable) bool {
	if p0.AI > 0 {
		if m.P1.AI > 0 {
			return false
		}
		if m.P2.AI > 0 {
			return false
		}
		if m.P3.AI > 0 {
			return false
		}
		if m.P4.AI > 0 {
			return false
		}
		if m.P5.AI > 0 {
			return false
		}
		if m.P6.AI > 0 {
			return false
		}
	}
	if p0.AE > 0 {
		if m.P1.AE > 0 {
			return false
		}
		if m.P2.AE > 0 {
			return false
		}
		if m.P3.AE > 0 {
			return false
		}
		if m.P4.AE > 0 {
			return false
		}
		if m.P5.AE > 0 {
			return false
		}
		if m.P6.AE > 0 {
			return false
		}
	}
	if p0.R1 > 0 {
		if m.P1.R1 > 0 {
			return false
		}
		if m.P2.R1 > 0 {
			return false
		}
		if m.P3.R1 > 0 {
			return false
		}
		if m.P4.R1 > 0 {
			return false
		}
		if m.P5.R1 > 0 {
			return false
		}
		if m.P6.R1 > 0 {
			return false
		}
	}
	if p0.R2 > 0 {
		if m.P1.R2 > 0 {
			return false
		}
		if m.P2.R2 > 0 {
			return false
		}
		if m.P3.R2 > 0 {
			return false
		}
		if m.P4.R2 > 0 {
			return false
		}
		if m.P5.R2 > 0 {
			return false
		}
		if m.P6.R2 > 0 {
			return false
		}
	}
	if p0.R3 > 0 {
		if m.P1.R3 > 0 {
			return false
		}
		if m.P2.R3 > 0 {
			return false
		}
		if m.P3.R3 > 0 {
			return false
		}
		if m.P4.R3 > 0 {
			return false
		}
		if m.P5.R3 > 0 {
			return false
		}
		if m.P6.R3 > 0 {
			return false
		}
	}
	if p0.R4 > 0 {
		if m.P1.R4 > 0 {
			return false
		}
		if m.P2.R4 > 0 {
			return false
		}
		if m.P3.R4 > 0 {
			return false
		}
		if m.P4.R4 > 0 {
			return false
		}
		if m.P5.R4 > 0 {
			return false
		}
		if m.P6.R4 > 0 {
			return false
		}
	}
	return true

}
