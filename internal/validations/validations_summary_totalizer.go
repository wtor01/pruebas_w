package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
)

type ValidatorSummaryTotalizer struct {
	ValidatorBase
	Tolerance float64
}

var validSummaryTotalizer = map[string]struct{}{
	AI: {},
}

func NewValidatorSummaryTotalizer(v ValidationData, status measures.Status) (ValidatorSummaryTotalizer, error) {

	if len(v.Config) != 1 {
		return ValidatorSummaryTotalizer{}, newErrorConfigProperties()
	}

	validator := ValidatorSummaryTotalizer{
		ValidatorBase: NewValidatorBase(v, status),
	}

	if err := validator.IsValidKeys(validSummaryTotalizer); err != nil {
		return ValidatorSummaryTotalizer{}, err
	}
	tol, err := utils.ParseMapKey(v.Config, Tolerance, utils.ParseFloat64)
	if err != nil {
		return ValidatorSummaryTotalizer{}, err
	}
	validator.Tolerance = tol
	return validator, nil
}

func (v ValidatorSummaryTotalizer) getValueKey(m MeasureValidatable, key string) float64 {
	switch key {
	case AI:
		return m.AI
	default:
		return 0
	}
}

func (v ValidatorSummaryTotalizer) Validate(m MeasureValidatable) *ValidatorBase {
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
	checkValue := v.getPeriodSummatory(m)
	if !v.comparePeriods(m.P0.AI, checkValue.AI) {
		return &v.ValidatorBase
	}
	if !v.comparePeriods(m.P0.AE, checkValue.AE) {
		return &v.ValidatorBase
	}
	if !v.comparePeriods(m.P0.R1, checkValue.R1) {
		return &v.ValidatorBase
	}
	if !v.comparePeriods(m.P0.R2, checkValue.R2) {
		return &v.ValidatorBase
	}
	if !v.comparePeriods(m.P0.R3, checkValue.R3) {
		return &v.ValidatorBase
	}
	if !v.comparePeriods(m.P0.R4, checkValue.R4) {
		return &v.ValidatorBase
	}

	return nil
}

func (v ValidatorSummaryTotalizer) getPeriodSummatory(m MeasureValidatable) measures.Values {

	return measures.Values{

		AI: m.P1.AI + m.P2.AI + m.P3.AI + m.P4.AI + m.P5.AI + m.P6.AI,
		AE: m.P1.AE + m.P2.AE + m.P3.AE + m.P4.AE + m.P5.AE + m.P6.AE,
		R1: m.P1.R1 + m.P2.R1 + m.P3.R1 + m.P4.R1 + m.P5.R1 + m.P6.R1,
		R2: m.P1.R2 + m.P2.R2 + m.P3.R2 + m.P4.R2 + m.P5.R2 + m.P6.R2,
		R3: m.P1.R3 + m.P2.R3 + m.P3.R3 + m.P4.R3 + m.P5.R3 + m.P6.R3,
		R4: m.P1.R4 + m.P2.R4 + m.P3.R4 + m.P4.R4 + m.P5.R4 + m.P6.R4,
	}

}
func (v ValidatorSummaryTotalizer) comparePeriods(p0 float64, sum float64) bool {
	min := p0 - v.Tolerance
	max := p0 + v.Tolerance
	if sum >= min && sum <= max {
		return true
	}
	return false
}
