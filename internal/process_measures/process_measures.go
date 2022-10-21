package process_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidDistributorId = errors.New("invalid distributor id")
	ErrInvalidMeasureData   = errors.New("invalid measure data")
	ErrInvalidMeasureDate   = errors.New("invalid measure date")
)

type QueryProcessedLoadCurveByCups struct {
	CUPS      string
	StartDate time.Time
	EndDate   time.Time
	CurveType measures.MeasureCurveReadingType
	Status    measures.Status
}

type QueryClosedCupsMeasureOnDate struct {
	CUPS string
	Date time.Time
}

type QueryMonthlyClosedMeasures struct {
	CUPS      string
	StartDate time.Time
	EndDate   time.Time
}

type QueryHistoryDailyClosureByCups struct {
	CUPS      string
	StartDate time.Time
	EndDate   time.Time
}

type QueryHistoryLoadCurve struct {
	CUPS          string
	StartDate     time.Time
	EndDate       time.Time
	Count         int
	Periods       []measures.PeriodKey
	Magnitude     measures.Magnitude
	WithCriterias bool
	IsFuture      bool
	Type          measures.RegisterType
}

type ProcessMeasureBase interface {
	SetStatusMeasure(validation validations.ValidatorBase)
	ToValidatable() validations.MeasureValidatable
}

type ToMeasuresBase interface {
	MeasuresBase() []ProcessMeasureBase
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=ProcessedMeasureRepository
type ProcessedMeasureRepository interface {
	SaveDailyClosure(ctx context.Context, processed ProcessedDailyClosure) error
	SaveMonthlyClosure(ctx context.Context, processed ProcessedMonthlyClosure) error
	SaveAllProcessedLoadCurve(ctx context.Context, processed []ProcessedLoadCurve) error
	SaveProcessedDailyLoadCurve(ctx context.Context, processed ProcessedDailyLoadCurve) error
	ProcessedLoadCurveByCups(ctx context.Context, q QueryProcessedLoadCurveByCups) ([]ProcessedLoadCurve, error)
	GetMonthlyClosureByCup(ctx context.Context, q QueryClosedCupsMeasureOnDate) (ProcessedMonthlyClosure, error)
	GetProcessedDailyClosureByCup(ctx context.Context, q QueryClosedCupsMeasureOnDate) (ProcessedDailyClosure, error)
	ListHistoryLoadCurve(ctx context.Context, q QueryHistoryLoadCurve) ([]ProcessedLoadCurve, error)
	GetLoadCurveByID(ctx context.Context, id string) (ProcessedLoadCurve, error)
	GetDailyClosureByID(ctx context.Context, id string) (ProcessedDailyClosure, error)
	GetMonthlyClosureByID(ctx context.Context, id string) (ProcessedMonthlyClosure, error)
	ProcessedDailyClosureByCups(ctx context.Context, q QueryHistoryDailyClosureByCups) ([]ProcessedDailyClosure, error)
	GetMonthlyClosureMeasuresByCup(ctx context.Context, q QueryMonthlyClosedMeasures) ([]ProcessedMonthlyClosure, error)
}

type Repository interface {
	ProcessedMeasureRepository
	ProcessMeasureDashboardRepository
}

func GenerateEvents(
	readingType measures.ReadingType,
	date time.Time,
	meterConfig measures.MeterConfig,
) []measures.ProcessMeasureEvent {
	var ev []measures.ProcessMeasureEvent
	switch readingType {
	case measures.Curve:
		{
			if utils.InSlice(meterConfig.CurveType, []measures.RegisterType{measures.Hourly, measures.Both}) {
				ev = append(ev, NewProcessCurveEvent(date, meterConfig, measures.HourlyMeasureCurveReadingType))
			}

			if utils.InSlice(meterConfig.CurveType, []measures.RegisterType{measures.QuarterHour, measures.Both}) {
				ev = append(ev, NewProcessCurveEvent(date, meterConfig, measures.QuarterMeasureCurveReadingType))
			}

		}
	case measures.BillingClosure:
		{
			ev = append(ev, NewProcessBillingClosureEvent(date, meterConfig))
		}
	case measures.DailyClosure:
		{
			ev = append(ev, NewProcessDailyClosureEvent(date, meterConfig))
		}
	}
	return ev
}
