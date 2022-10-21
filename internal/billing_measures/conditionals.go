package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
)

type Condition interface {
	Eval(ctx context.Context) bool
}

type Simple struct {
	R bool
}

func (s Simple) Eval(ctx context.Context) bool {
	return s.R
}

type IsTlg struct {
	ID string
	b  *BillingMeasure `bson:"-"`
}

func NewIsTlg(b *BillingMeasure) *IsTlg {
	return &IsTlg{b: b, ID: "IS_TLG"}
}

func (i IsTlg) Eval(ctx context.Context) bool {
	return i.b.MeterType == measures.TLG
}

type Conditional interface {
	Condition(ctx context.Context, t *BillingMeasure) bool
}

type IsPointType struct {
	ID         string
	PointTypes []measures.PointType `bson:"-"`
	b          *BillingMeasure      `bson:"-"`
}

func NewIsPointType(pointTypes []measures.PointType, b *BillingMeasure) *IsPointType {
	return &IsPointType{PointTypes: pointTypes, b: b, ID: "IS_POINT_TYPE"}
}

func (i IsPointType) Eval(ctx context.Context) bool {
	for _, p := range i.PointTypes {
		if p == i.b.PointType {
			return true
		}
	}
	return false
}

type IsRegisterType struct {
	ID           string
	RegisterType measures.RegisterType `bson:"-"`
	b            *BillingMeasure       `bson:"-"`
}

func NewIsRegisterType(registerType measures.RegisterType, b *BillingMeasure) *IsRegisterType {
	return &IsRegisterType{RegisterType: registerType, b: b, ID: "IS_REGISTER_TYPE"}
}

func (i IsRegisterType) Eval(ctx context.Context) bool {
	return i.RegisterType == i.b.RegisterType
}

type IsValidBillingAtrByPeriod struct {
	ID        string
	Period    measures.PeriodKey `bson:"-"`
	Magnitude measures.Magnitude `bson:"-"`
	b         *BillingMeasure    `bson:"-"`
}

func NewIsValidBillingAtrByPeriod(period measures.PeriodKey, magnitude measures.Magnitude, b *BillingMeasure) *IsValidBillingAtrByPeriod {
	return &IsValidBillingAtrByPeriod{
		Period: period,
		b:      b, ID: "SALDO_ATR_VALIDO",
		Magnitude: magnitude,
	}
}

func (i IsValidBillingAtrByPeriod) Eval(ctx context.Context) bool {
	return i.b.IsValidBillingBalancePeriod(i.Period, i.Magnitude)
}

type Negate struct {
	ID string
	C  Condition `bson:"condition"`
}

func NewNegate(c Condition) *Negate {
	return &Negate{C: c, ID: "NEGATE"}
}

func (n Negate) Eval(ctx context.Context) bool {
	return !n.C.Eval(ctx)
}

type SeveralEval struct {
	ID string
	C  []Condition `bson:"condition"`
}

func NewSeveralEval(c ...Condition) *SeveralEval {
	return &SeveralEval{C: c, ID: "SEVERAL"}
}

func (n SeveralEval) Eval(ctx context.Context) bool {

	for _, c := range n.C {
		if c.Eval(ctx) == false {
			return false
		}
	}
	return true
}

type IsCurvePeriodCompletedByPeriod struct {
	ID     string
	Period measures.PeriodKey `bson:"-"`
	b      *BillingMeasure    `bson:"-"`
}

func NewIsCurvePeriodCompletedByPeriod(period measures.PeriodKey, b *BillingMeasure) *IsCurvePeriodCompletedByPeriod {
	return &IsCurvePeriodCompletedByPeriod{Period: period, b: b, ID: "CCH_COMPLETA"}
}

func (i IsCurvePeriodCompletedByPeriod) Eval(ctx context.Context) bool {
	return i.b.IsCurvePeriodCompletedByPeriod(i.Period)
}

type IsAtrVsCurveValid struct {
	ID        string
	Period    measures.PeriodKey `bson:"-"`
	Magnitude measures.Magnitude `bson:"-"`
	b         *BillingMeasure    `bson:"-"`
}

func NewIsAtrVsCurveValidByPeriod(period measures.PeriodKey, b *BillingMeasure, magnitude measures.Magnitude) *IsAtrVsCurveValid {
	return &IsAtrVsCurveValid{Period: period, b: b, ID: "ISaldo ATR - CCH < 1KW", Magnitude: magnitude}
}

func (conditional IsAtrVsCurveValid) Eval(ctx context.Context) bool {
	return conditional.b.IsAtrVsCurveValid(conditional.Period, conditional.Magnitude)
}

type AreSomeCurveMeasureForPeriod struct {
	ID     string
	Period measures.PeriodKey `bson:"-"`
	b      *BillingMeasure    `bson:"-"`
}

func NewAreSomeCurveMeasureForPeriod(period measures.PeriodKey, b *BillingMeasure) *AreSomeCurveMeasureForPeriod {
	return &AreSomeCurveMeasureForPeriod{Period: period, b: b, ID: "CCH_SOME_CURVE_MEASURE_FOR_PERIOD"}
}

func (i AreSomeCurveMeasureForPeriod) Eval(ctx context.Context) bool {
	return i.b.AreSomeCurveMeasureForPeriod(i.Period)
}

type HasBillingHistory struct {
	ID                       string
	b                        *BillingMeasure          `bson:"-"`
	billingMeasureRepository BillingMeasureRepository `bson:"-"`
	Period                   measures.PeriodKey       `bson:"-"`
	Magnitude                measures.Magnitude       `bson:"-"`
	GraphContext             *GraphContext            `bson:"-"`
}

func NewHasBillingHistory(
	b *BillingMeasure,
	repository BillingMeasureRepository,
	period measures.PeriodKey,
	magnitude measures.Magnitude,
	contextTlg *GraphContext,
) *HasBillingHistory {
	return &HasBillingHistory{
		ID:                       "HAVE_BILLING_HISTORY",
		b:                        b,
		billingMeasureRepository: repository,
		Period:                   period,
		Magnitude:                magnitude,
		GraphContext:             contextTlg,
	}
}

func (conditional HasBillingHistory) Eval(ctx context.Context) bool {
	if !conditional.GraphContext.IsLastHistoryRequested {
		billingHistory, err := conditional.billingMeasureRepository.LastHistory(ctx, QueryLastHistory{
			CUPS:       conditional.b.CUPS,
			InitDate:   conditional.b.InitDate,
			EndDate:    conditional.b.EndDate,
			Periods:    []measures.PeriodKey{conditional.Period},
			Magnitudes: []measures.Magnitude{conditional.Magnitude},
		})
		if err == nil {
			conditional.GraphContext.LastHistory = billingHistory
		}
		conditional.GraphContext.IsLastHistoryRequested = true

	}

	return conditional.GraphContext.LastHistory.Id != ""
}

type IsBalanceGreaterCch struct {
	ID        string
	Period    measures.PeriodKey `bson:"-"`
	b         *BillingMeasure    `bson:"-"`
	Magnitude measures.Magnitude `bson:"-"`
}

func NewIsBalanceGreaterCch(b *BillingMeasure, period measures.PeriodKey, magnitude measures.Magnitude) *IsBalanceGreaterCch {
	return &IsBalanceGreaterCch{
		ID:        "IS_BALANCE_GREATER_CCH",
		Period:    period,
		b:         b,
		Magnitude: magnitude,
	}
}

func (i IsBalanceGreaterCch) Eval(ctx context.Context) bool {
	return i.b.IsBalanceGreaterCch(i.Period, i.Magnitude)
}

type IsChhComplete struct {
	ID string
	b  *BillingMeasure `bson:"-"`
}
