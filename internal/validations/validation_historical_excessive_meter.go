package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
)

type ValidatorHistoricalExcessive struct {
	ValidatorBase
	Limit   float64
	valRepo ValidationMongoRepository
}

var ValidHistoricalExcessive = map[string]struct{}{
	AI: {},
	AE: {},
}

func NewValidatorHistoricalExcessive(v ValidationData, status measures.Status, valRepo ValidationMongoRepository) (ValidatorHistoricalExcessive, error) {

	if len(v.Config) != 1 {
		return ValidatorHistoricalExcessive{}, newErrorConfigProperties()
	}

	validator := ValidatorHistoricalExcessive{
		ValidatorBase: NewValidatorBase(v, status),
		valRepo:       valRepo,
	}

	if err := validator.IsValidKeys(ValidHistoricalExcessive); err != nil {
		return ValidatorHistoricalExcessive{}, err
	}
	limit, err := utils.ParseMapKey(v.Config, Limit, utils.ParseFloat64)
	if err != nil {
		return ValidatorHistoricalExcessive{}, err
	}
	validator.Limit = limit
	return validator, nil
}

func (v ValidatorHistoricalExcessive) getValueKey(m MeasureValidatable, key string) float64 {
	switch key {
	case AI:
		return m.AI
	default:
		return 0
	}
}

func (v ValidatorHistoricalExcessive) Validate(m MeasureValidatable) *ValidatorBase {
	if m.Type != v.MeasureType {
		return nil
	}

	checkValue := v.getLastPeriods(m)
	if len(checkValue) < 1 {
		return nil
	}
	if m.P0i.AI > v.Limit || m.P0i.AE > v.Limit {
		return &v.ValidatorBase
	}
	for _, value := range checkValue {
		if value.AE > v.Limit || value.AI > v.Limit {
			return &v.ValidatorBase
		}
	}
	return nil
}

func (v ValidatorHistoricalExcessive) getLastPeriods(m MeasureValidatable) []measures.Values {
	mv := make([]measures.Values, 0)
	ctx := context.Background()
	for i := -1; i >= -12; i-- {
		pmc, err := v.valRepo.GetMonthlyClosureByCup(ctx, QueryClosedCupsMeasureOnDate{CUPS: m.CUPS, Date: m.EndDate.AddDate(0, i, 0)})
		if err != nil {
			return mv
		}
		mv = append(mv, measures.Values{AI: pmc.CalendarPeriods.P0.AIi, AE: pmc.CalendarPeriods.P0.AEi})

	}
	return mv

}
