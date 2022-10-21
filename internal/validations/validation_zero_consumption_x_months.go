package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
)

type ValidatorZeroConsumptionXMonths struct {
	ValidatorBase
	Months  float64
	valRepo ValidationMongoRepository
}

var ValidValidatorZeroConsumptionXMonths = map[string]struct{}{
	AI: {},
	AE: {},
}

func NewValidatorZeroConsumptionXMonths(v ValidationData, status measures.Status, valRepo ValidationMongoRepository) (ValidatorZeroConsumptionXMonths, error) {

	if len(v.Config) != 1 {
		return ValidatorZeroConsumptionXMonths{}, newErrorConfigProperties()
	}
	validator := ValidatorZeroConsumptionXMonths{
		ValidatorBase: NewValidatorBase(v, status),
		valRepo:       valRepo,
	}
	if err := validator.IsValidKeys(ValidHistoricalExcessive); err != nil {
		return ValidatorZeroConsumptionXMonths{}, err
	}
	months, err := utils.ParseMapKey(v.Config, Months, utils.ParseFloat64)
	if err != nil || months < 2 || months > 12 {
		return ValidatorZeroConsumptionXMonths{}, err
	}
	validator.Months = months

	return validator, nil
}

func (v ValidatorZeroConsumptionXMonths) getValueKey(m MeasureValidatable, key string) float64 {
	switch key {
	case AI:
		return m.AI
	default:
		return 0
	}
}

func (v ValidatorZeroConsumptionXMonths) Validate(m MeasureValidatable) *ValidatorBase {
	if m.Type != v.MeasureType {
		return nil
	}
	if m.TechnologyType == "INACCESIBLE" || m.TechnologyType == "TLM" || m.TechnologyType == "TLG" {
		return nil
	}
	values := v.getLastPeriods(m)
	if len(values) < int(v.Months) {
		return nil
	}
	if v.comparePeriods(values) {
		return &v.ValidatorBase
	}
	return nil
}
func (v ValidatorZeroConsumptionXMonths) comparePeriods(values []measures.Values) bool {
	ai := true
	ae := true
	r1 := true
	r2 := true
	r3 := true
	r4 := true
	for _, val := range values {
		if val.AI != 0 {
			ai = false
		}
		if val.AE != 0 {
			ae = false
		}
		if val.R1 != 0 {
			r1 = false
		}
		if val.R2 != 0 {
			r2 = false
		}
		if val.R3 != 0 {
			r3 = false
		}
		if val.R4 != 0 {
			r4 = false
		}
	}
	if ai || ae || r1 || r2 || r3 || r4 {
		return true
	}
	return false
}
func (v ValidatorZeroConsumptionXMonths) getLastPeriods(m MeasureValidatable) []measures.Values {

	mv := make([]measures.Values, 0)
	ctx := context.Background()
	for i := -1; i >= -int(v.Months); i-- {
		pmc, err := v.valRepo.GetMonthlyClosureByCup(ctx, QueryClosedCupsMeasureOnDate{CUPS: m.CUPS, Date: m.EndDate.AddDate(0, i, 0)})
		if err != nil {
			return mv
		}
		mv = append(mv, measures.Values{AI: pmc.CalendarPeriods.P0.AIi, AE: pmc.CalendarPeriods.P0.AEi, R1: pmc.CalendarPeriods.P0.R1i, R2: pmc.CalendarPeriods.P0.R2i, R3: pmc.CalendarPeriods.P0.R3i, R4: pmc.CalendarPeriods.P0.R4i})

	}
	return mv

}
