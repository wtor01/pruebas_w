package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
)

/*
solo se aplican a curva

*/

type ValidatorCurveLimit struct {
	ValidatorBase
	Type1 float64
	Type2 float64
	Type3 float64
	Type4 float64
	Type5 float64
}

var validCurveLimit = map[string]struct{}{
	AI: {},
}

func NewValidatorCurveLimit(v ValidationData, status measures.Status) (ValidatorCurveLimit, error) {

	if len(v.Config) != 5 {
		return ValidatorCurveLimit{}, newErrorConfigProperties()
	}

	validator := ValidatorCurveLimit{
		ValidatorBase: NewValidatorBase(v, status),
	}

	if err := validator.IsValidKeys(validCurveLimit); err != nil {
		return ValidatorCurveLimit{}, err
	}

	type1, err := utils.ParseMapKey(v.Config, Type1, utils.ParseFloat64)
	if err != nil {
		return ValidatorCurveLimit{}, err
	}
	type2, err := utils.ParseMapKey(v.Config, Type2, utils.ParseFloat64)
	if err != nil {
		return ValidatorCurveLimit{}, err
	}
	type3, err := utils.ParseMapKey(v.Config, Type3, utils.ParseFloat64)
	if err != nil {
		return ValidatorCurveLimit{}, err
	}
	type4, err := utils.ParseMapKey(v.Config, Type4, utils.ParseFloat64)
	if err != nil {
		return ValidatorCurveLimit{}, err
	}
	type5, err := utils.ParseMapKey(v.Config, Type5, utils.ParseFloat64)

	if err != nil {
		return ValidatorCurveLimit{}, err
	}

	validator.Type1 = type1
	validator.Type2 = type2
	validator.Type3 = type3
	validator.Type4 = type4
	validator.Type5 = type5

	return validator, nil
}

func (v ValidatorCurveLimit) getValueKey(m MeasureValidatable, key string) float64 {
	switch key {
	case AI:
		return m.AI
	default:
		return 0
	}
}

func (v ValidatorCurveLimit) Validate(m MeasureValidatable) *ValidatorBase {
	if m.Type != v.MeasureType {
		return nil
	}
	if m.ServiceType != measures.DcServiceType {
		return nil
	}
	for _, k := range v.Keys {
		if !m.isInWhiteList(k) {
			continue
		}
		value := v.getValueKey(m, k)
		checkValue := v.getCompareType(m)
		if value > checkValue {
			return &v.ValidatorBase
		}
	}
	return nil
}
func (v ValidatorCurveLimit) getCompareType(m MeasureValidatable) float64 {
	switch m.PointType {
	case "1":
		return v.Type1
	case "2":
		return v.Type2
	case "3":
		return v.Type3
	case "4":
		return v.Type4
	case "5":
		return v.Type5
	default:
		return 0
	}

}
