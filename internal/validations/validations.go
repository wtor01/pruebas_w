package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"errors"
	"fmt"
	"time"
)

const (
	Threshold                      string = "threshold"
	QualifierValidation            string = "qualifier"
	DailyDate                      string = "daily_date"
	FutureDate                     string = "future_date"
	CurveLimit                     string = "curve_limit"
	SummaryTotalizer               string = "summary_totalizer"
	SummaryCalendar                string = "summary_calendar"
	ZeroConsmuption                string = "zero_consumption"
	ClosedHouse                    string = "closed_house"
	ExcessiveConsumption           string = "excessive_consumption"
	CloseMeter                     string = "close_meter"
	HourDate                       string = "hour_date"
	QuarterlyDate                  string = "quarterly_date"
	ExcessivePowerConsumption      string = "excessive_power_consumption"
	ZeroConsumptionXMonths         string = "zero_consumption_x_months"
	HistoricalExcessiveConsumption string = "historical_excessive_consumption"
	ConsumptionThreshold           string = "consumption_threshold"
)
const (
	MaxKey = "max"
	MinKey = "min"
)

const (
	Type1 = "type1"
	Type2 = "type2"
	Type3 = "type3"
	Type4 = "type4"
	Type5 = "type5"
)

const (
	Tolerance = "tolerance"
	Limit     = "limit"
	Months    = "months"
	Maxlimit  = "max_limit"
	Minlimit  = "min_limit"
)

const (
	StartDate   string = "start_date"
	EndDate     string = "end_date"
	MeasureDate string = "measure_date"
	AI          string = "AI"
	AE          string = "AE"
	R1          string = "R1"
	R2          string = "R2"
	R3          string = "R3"
	R4          string = "R4"
	MX          string = "MX"
	E           string = "E"
	FX          string = "FX"
	Qualifier   string = "qualifier"
)

func newErrorInvalidKeyForValidation(key string, validationName string, validKeys []string) error {
	return errors.New(fmt.Sprintf("invalid key %s for %s, valid keys %v", key, validationName, validKeys))
}

func newErrorRanged(minValue, maxValue float64) error {
	return errors.New(fmt.Sprintf("%v min value cant be greater than %v max value", minValue, maxValue))
}
func newErrorConfigProperties() error {
	return errors.New(fmt.Sprintf("Error config properties number"))
}

type ValidationData struct {
	Type        string            `json:"type"`
	Keys        []string          `json:"keys"`
	Config      map[string]string `json:"config"`
	Code        string            `json:"code"`
	MeasureType measures.Type     `json:"measure_type"`
}

type ValidatorConfig struct {
	Action measures.Status
	Params struct {
		Type        string           `json:"type"`
		Validations []ValidationData `json:"validations"`
	} `json:"params"`
}

type ValidatorBase struct {
	Type        string            `json:"type"`
	Keys        []string          `json:"keys"`
	Config      map[string]string `json:"config"`
	Action      measures.Status
	Code        string        `json:"code"`
	MeasureType measures.Type `json:"measure_type"`
}

func NewValidatorBase(v ValidationData, status measures.Status) ValidatorBase {
	return ValidatorBase{
		Type:        v.Type,
		Keys:        v.Keys,
		Config:      v.Config,
		Action:      status,
		Code:        v.Code,
		MeasureType: v.MeasureType,
	}
}

type MeasureUpdatableStatus interface {
	ChangeStatus(status measures.Status)
	AppendInvalidation(invalidation string)
}

func (b ValidatorBase) SetStatusMeasure(m MeasureUpdatableStatus) {
	m.ChangeStatus(b.Action)
	m.AppendInvalidation(b.Code)
}

func (b ValidatorBase) IsValidKeys(validKeys map[string]struct{}) error {
	validKeysFormat := make([]string, 0, len(validKeys))

	for k, _ := range validKeys {
		validKeysFormat = append(validKeysFormat, k)
	}

	for _, k := range b.Keys {
		if _, ok := validKeys[k]; !ok {
			return newErrorInvalidKeyForValidation(k, b.Type, validKeysFormat)
		}
	}

	return nil
}

type MeasureValidatable struct {
	CUPS           string
	Type           measures.Type
	StartDate      time.Time
	EndDate        time.Time
	ReadingDate    time.Time
	RegisterType   measures.RegisterType
	AI             float64
	AE             float64
	R1             float64
	R2             float64
	R3             float64
	R4             float64
	MX             float64
	E              float64
	FX             time.Time
	Qualifier      string
	ServiceType    measures.ServiceType
	P0             Periods
	P1             Periods
	P2             Periods
	P3             Periods
	P4             Periods
	P5             Periods
	P6             Periods
	P0i            Periods
	P1i            Periods
	P2i            Periods
	P3i            Periods
	P4i            Periods
	P5i            Periods
	P6i            Periods
	P0LM           Periods
	P1LM           Periods
	P2LM           Periods
	P3LM           Periods
	P4LM           Periods
	P5LM           Periods
	P6LM           Periods
	CalendarCode   string
	PeriodsNumbers []string
	PointType      string
	TechnologyType string
	WhiteListKeys  map[string]struct{}
	Period         measures.PeriodKey
	P1Demand       int
	P2Demand       int
	P3Demand       int
	P4Demand       int
	P5Demand       int
	P6Demand       int
}

func (m MeasureValidatable) isInWhiteList(key string) bool {
	_, ok := m.WhiteListKeys[key]

	return ok
}

type ValidatorI interface {
	Validate(m MeasureValidatable) *ValidatorBase
}

func NewValidator(v ValidationData, status measures.Status, validationRepository ValidationMongoRepository) (ValidatorI, error) {
	switch v.Type {
	case Threshold:
		{
			return NewValidatorThreshold(v, status)
		}

	case Qualifier:
		{
			return NewValidatorQualifier(v, status)
		}

	case HourDate:
		{
			return NewValidatorHourly(v, status)
		}

	case DailyDate:
		{
			return NewValidatorDaily(v, status)
		}

	case FutureDate:
		{
			return NewValidatorFuture(v, status)
		}
	case QuarterlyDate:
		{
			return NewValidatorQuarterly(v, status)
		}
	case CurveLimit:
		return NewValidatorCurveLimit(v, status)
	case SummaryTotalizer:
		return NewValidatorSummaryTotalizer(v, status)
	case SummaryCalendar:
		return NewValidatorSummaryCalendar(v, status)
	case ZeroConsmuption:
		return NewValidatorZeroConsumption(v, status)
	case CloseMeter:
		return NewValidatorCloseMeter(v, status)
	case ClosedHouse:
		return NewValidatorClosedHouse(v, status)
	case ExcessiveConsumption:
		return NewValidatorExcessiveConsumption(v, status)
	case ExcessivePowerConsumption:
		return NewValidatorValidatorExcesivePowerConsumption(v, status, validationRepository)
	case ZeroConsumptionXMonths:
		return NewValidatorZeroConsumptionXMonths(v, status, validationRepository)
	case ConsumptionThreshold:
		return NewValidatorConsumptionThreshold(v, status, validationRepository)
	case HistoricalExcessiveConsumption:
		return NewValidatorHistoricalExcessive(v, status, validationRepository)

	default:
		return nil, errors.New("")
	}
}

func NewValidatorFromClient(config clients.ValidationConfig, valRepo ValidationMongoRepository) ([]ValidatorI, error) {

	if !config.Enabled {
		return []ValidatorI{}, nil
	}

	mapConfigs := make(map[string]struct {
		Id     string            `json:"id"`
		Name   string            `json:"name"`
		Config map[string]string `json:"config"`
	})

	for _, e := range config.ExtraConfig {
		mapConfigs[e.Id] = struct {
			Id     string            `json:"id"`
			Name   string            `json:"name"`
			Config map[string]string `json:"config"`
		}{
			Id:     e.Id,
			Name:   e.Name,
			Config: e.Config,
		}
	}

	validators := make([]ValidatorI, 0)

	for _, p := range config.ValidationMeasure.Params.Validations {
		conf, _ := mapConfigs[p.Id]

		utils.Merge(p.Config, conf.Config)

		v := ValidationData{
			Type:        p.Type,
			Keys:        p.Keys,
			Config:      p.Config,
			Code:        config.ValidationMeasure.Code,
			MeasureType: measures.Type(config.ValidationMeasure.MeasureType),
		}

		validator, err := NewValidator(v, measures.Status(config.Action), valRepo)
		if err != nil {
			continue
		}

		validators = append(validators, validator)
	}

	return validators, nil
}
