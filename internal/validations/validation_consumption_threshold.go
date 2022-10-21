package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
)

type ValidatorConsumptionThreshold struct {
	ValidatorBase
	valRepo ValidationMongoRepository
}

var ValidValidatorConsumptionThreshold = map[string]struct{}{
	AI: {},
	AE: {},
}

func NewValidatorConsumptionThreshold(v ValidationData, status measures.Status, valRepo ValidationMongoRepository) (ValidatorConsumptionThreshold, error) {

	if len(v.Config) != 1 {
		return ValidatorConsumptionThreshold{}, newErrorConfigProperties()
	}
	validator := ValidatorConsumptionThreshold{
		ValidatorBase: NewValidatorBase(v, status),
		valRepo:       valRepo,
	}
	if err := validator.IsValidKeys(ValidHistoricalExcessive); err != nil {
		return ValidatorConsumptionThreshold{}, err
	}

	return validator, nil
}

func (v ValidatorConsumptionThreshold) getValueKey(m MeasureValidatable, key string) float64 {
	switch key {
	case AI:
		return m.AI
	default:
		return 0
	}
}

func (v ValidatorConsumptionThreshold) Validate(m MeasureValidatable) *ValidatorBase {
	if m.Type != v.MeasureType {
		return nil
	}

	pmc, err := v.getCurves(m)
	if err != nil {
		return nil
	}
	p1, p2, p3, p4, p5, p6 := v.comparePeriods(pmc)
	if !v.compare(m.P1Demand, p1, m.P1i.AI, m.P1i.AE) {
		return &v.ValidatorBase

	}
	if !v.compare(m.P2Demand, p2, m.P2i.AI, m.P2i.AE) {
		return &v.ValidatorBase

	}
	if !v.compare(m.P3Demand, p3, m.P3i.AI, m.P3i.AE) {
		return &v.ValidatorBase

	}
	if !v.compare(m.P4Demand, p4, m.P4i.AI, m.P4i.AE) {
		return &v.ValidatorBase

	}
	if !v.compare(m.P5Demand, p5, m.P5i.AI, m.P5i.AE) {
		return &v.ValidatorBase

	}
	if !v.compare(m.P6Demand, p6, m.P6i.AI, m.P6i.AE) {
		return &v.ValidatorBase

	}
	return nil
}
func (v ValidatorConsumptionThreshold) compare(demand int, hours float64, AI float64, AE float64) bool {
	max := float64(demand) * hours
	if AI > max {
		return false
	}
	if AE > max {
		return false
	}
	return true
}
func (v ValidatorConsumptionThreshold) comparePeriods(pmc []ProcessedLoadCurve) (float64, float64, float64, float64, float64, float64) {
	var p1 float64
	var p2 float64
	var p3 float64
	var p4 float64
	var p5 float64
	var p6 float64

	for _, val := range pmc {
		switch val.Period {
		case measures.P1:
			if val.Type == measures.QuarterHour {
				p1 += 0.25
			}
			if val.Type == measures.Hourly {
				p1 += 1
			}
		case measures.P2:
			if val.Type == measures.QuarterHour {
				p2 += 0.25
			}
			if val.Type == measures.Hourly {
				p2 += 1
			}
		case measures.P3:
			if val.Type == measures.QuarterHour {
				p3 += 0.25
			}
			if val.Type == measures.Hourly {
				p3 += 1
			}
		case measures.P4:
			if val.Type == measures.QuarterHour {
				p4 += 0.25
			}
			if val.Type == measures.Hourly {
				p4 += 1
			}
		case measures.P5:
			if val.Type == measures.QuarterHour {
				p5 += 0.25
			}
			if val.Type == measures.Hourly {
				p5 += 1
			}
		case measures.P6:
			if val.Type == measures.QuarterHour {
				p6 += 0.25
			}
			if val.Type == measures.Hourly {
				p6 += 1
			}
		}

	}
	return p1, p2, p3, p4, p5, p6
}
func (v ValidatorConsumptionThreshold) getCurves(m MeasureValidatable) ([]ProcessedLoadCurve, error) {

	ctx := context.Background()

	pmc, err := v.valRepo.GetLoadCurveByQuery(ctx, QueryCurveCupsMeasureOnDate{CUPS: m.CUPS, StartDate: m.StartDate, EndDate: m.EndDate})

	return pmc, err

}
