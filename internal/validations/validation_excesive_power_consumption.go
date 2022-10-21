package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
)

type ValidatorExcesivePowerConsumption struct {
	ValidatorBase
	MaxLimit float64
	MinLimit float64
	valRepo  ValidationMongoRepository
}

var ValidValidatorExcesivePowerConsumption = map[string]struct{}{
	AI: {},
	AE: {},
}

func NewValidatorValidatorExcesivePowerConsumption(v ValidationData, status measures.Status, valRepo ValidationMongoRepository) (ValidatorExcesivePowerConsumption, error) {

	if len(v.Config) != 1 {
		return ValidatorExcesivePowerConsumption{}, newErrorConfigProperties()
	}
	validator := ValidatorExcesivePowerConsumption{
		ValidatorBase: NewValidatorBase(v, status),
		valRepo:       valRepo,
	}
	if err := validator.IsValidKeys(ValidHistoricalExcessive); err != nil {
		return ValidatorExcesivePowerConsumption{}, err
	}
	maxLimit, err := utils.ParseMapKey(v.Config, Maxlimit, utils.ParseFloat64)
	if err != nil {
		return ValidatorExcesivePowerConsumption{}, err
	}
	validator.MaxLimit = maxLimit
	minLimit, err := utils.ParseMapKey(v.Config, Minlimit, utils.ParseFloat64)
	if err != nil {
		return ValidatorExcesivePowerConsumption{}, err
	}
	validator.MinLimit = minLimit

	return validator, nil
}

func (v ValidatorExcesivePowerConsumption) getValueKey(m MeasureValidatable, key string) float64 {
	switch key {
	case AI:
		return m.AI
	default:
		return 0
	}
}

func (v ValidatorExcesivePowerConsumption) Validate(m MeasureValidatable) *ValidatorBase {
	if m.Type != v.MeasureType {
		return nil
	}
	if m.TechnologyType == "INACCESIBLE" || m.TechnologyType == "TLM" || m.TechnologyType == "TLG" {
		return nil
	}
	values := v.getLastMonth(m)
	if v.comparePeriods(values) {
		return &v.ValidatorBase
	}
	return nil
}
func (v ValidatorExcesivePowerConsumption) comparePeriods(val measures.Values) bool {
	ai := true
	ae := true
	r1 := true
	r2 := true
	r3 := true
	r4 := true

	if val.AI > v.MaxLimit || val.AI < v.MinLimit {
		ai = false
	}
	if val.AE > v.MaxLimit || val.AE < v.MinLimit {
		ae = false
	}
	if val.R1 > v.MaxLimit || val.R1 < v.MinLimit {
		r1 = false
	}
	if val.R2 > v.MaxLimit || val.R2 < v.MinLimit {
		r2 = false
	}
	if val.R3 > v.MaxLimit || val.R3 < v.MinLimit {
		r3 = false
	}
	if val.R4 > v.MaxLimit || val.R4 < v.MinLimit {
		r4 = false
	}

	if ai || ae || r1 || r2 || r3 || r4 {
		return true
	}
	return false
}
func (v ValidatorExcesivePowerConsumption) getLastMonth(m MeasureValidatable) measures.Values {

	ctx := context.Background()

	pmc, err := v.valRepo.GetMonthlyClosureByCup(ctx, QueryClosedCupsMeasureOnDate{CUPS: m.CUPS, Date: m.EndDate.AddDate(0, -1, 0)})
	if err != nil {
		return measures.Values{}
	}
	mv := measures.Values{AI: pmc.CalendarPeriods.P0.AIi, AE: pmc.CalendarPeriods.P0.AEi, R1: pmc.CalendarPeriods.P0.R1i, R2: pmc.CalendarPeriods.P0.R2i, R3: pmc.CalendarPeriods.P0.R3i, R4: pmc.CalendarPeriods.P0.R4i}

	return mv

}
