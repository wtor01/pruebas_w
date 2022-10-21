package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"math"
	"time"
)

type FiscalBillingMeasuresDashboard struct {
	Id                 string               `json:"id"`
	Cups               string               `json:"cups"`
	Type               measures.MeterType   `json:"type"`
	StartDate          string               `json:"start_date"`
	EndDate            string               `json:"end_date"`
	LastMvDate         string               `json:"last_mv_date"`
	PrincipalMagnitude measures.Magnitude   `json:"principal_magnitude"`
	Status             Status               `json:"status"`
	Periods            []measures.PeriodKey `json:"periods"`
	Magnitudes         []measures.Magnitude `json:"magnitudes"`
	Summary            Summary              `json:"summary"`
	ExecutionSummary   ExecutionSummary     `json:"execution_summary"`
	CalendarCurve      []CalendarCurve      `json:"calendar_curve"`
	Balance            Balance              `json:"balance"`
	GraphHistory       map[string]*Graph    `json:"graph_history"`
	Curve              []Curve              `json:"curve"`
}

func (obj *FiscalBillingMeasuresDashboard) SetBalancePeriod(period measures.PeriodKey, balancePeriod BalancePeriod) {
	obj.Balance.setBalancePeriod(balancePeriod, period)
}

func (obj *FiscalBillingMeasuresDashboard) SetCalendarCurve(curves []CalendarCurve) {
	obj.CalendarCurve = curves
}

func (obj *FiscalBillingMeasuresDashboard) SetCurveGraph(curves []Curve) {
	obj.Curve = curves
}

func (obj *FiscalBillingMeasuresDashboard) AddSummary(period measures.PeriodKey, method GeneralEstimateMethod) {
	summaryPeriod := obj.Summary.getSummaryItemPeriod(method, period)
	value := summaryPeriod + 1
	switch method {
	case GeneralReal:
		obj.Summary.Real.setPeriodValue(period, value)
	case GeneralAdjusted:
		obj.Summary.Adjusted.setPeriodValue(period, value)
	case GeneralOutlined:
		obj.Summary.Outlined.setPeriodValue(period, value)
	case GeneralCalculated:
		obj.Summary.Calculated.setPeriodValue(period, value)
	case GeneralEstimated:
		obj.Summary.Estimated.setPeriodValue(period, value)
	}
}

func (obj *FiscalBillingMeasuresDashboard) AddConsum(period measures.PeriodKey, value float64) {
	obj.Summary.Consum.addConsum(period, value)
}

func (obj *FiscalBillingMeasuresDashboard) CalculatePercentages(periods []measures.PeriodKey) {

	methods := []GeneralEstimateMethod{GeneralReal, GeneralAdjusted, GeneralOutlined, GeneralCalculated, GeneralEstimated}
	periodValue := make(map[measures.PeriodKey]float64)
	for _, period := range periods {
		periodValue[period] = obj.Summary.getSummaryPeriodTotal(period)
	}

	for _, method := range methods {
		for _, period := range periods {
			totalPeriod := periodValue[period]
			obj.Summary.calcSummaryPeriodPercentage(method, period, totalPeriod)
		}

		obj.Summary.calcSummaryTotalPercentage(method, periods)
	}

}

func (obj *FiscalBillingMeasuresDashboard) CalcTotalConsum(periods []measures.PeriodKey) {
	obj.Summary.Consum.sumConsumPeriods(periods)
}

type Summary struct {
	Real       SummaryItem `json:"real"`
	Adjusted   SummaryItem `json:"adjusted"`
	Outlined   SummaryItem `json:"outlined"`
	Calculated SummaryItem `json:"calculated"`
	Estimated  SummaryItem `json:"estimated"`
	Consum     SummaryItem `json:"consum"`
}

func (obj Summary) getSummaryItemPeriod(method GeneralEstimateMethod, period measures.PeriodKey) float64 {
	switch method {
	case GeneralReal:
		return obj.Real.getPeriod(period)
	case GeneralAdjusted:
		return obj.Adjusted.getPeriod(period)
	case GeneralOutlined:
		return obj.Outlined.getPeriod(period)
	case GeneralCalculated:
		return obj.Calculated.getPeriod(period)
	case GeneralEstimated:
		return obj.Estimated.getPeriod(period)
	default:
		return 0
	}
}

func (obj *Summary) calcSummaryPeriodPercentage(method GeneralEstimateMethod, period measures.PeriodKey, total float64) {
	switch method {
	case GeneralReal:
		obj.Real.calcPercentagePeriod(period, total)
	case GeneralAdjusted:
		obj.Adjusted.calcPercentagePeriod(period, total)
	case GeneralOutlined:
		obj.Outlined.calcPercentagePeriod(period, total)
	case GeneralCalculated:
		obj.Calculated.calcPercentagePeriod(period, total)
	case GeneralEstimated:
		obj.Estimated.calcPercentagePeriod(period, total)
	}
}

func (obj *Summary) calcSummaryTotalPercentage(method GeneralEstimateMethod, periods []measures.PeriodKey) {
	switch method {
	case GeneralReal:
		obj.Real.calcPercentageTotal(periods)
	case GeneralAdjusted:
		obj.Adjusted.calcPercentageTotal(periods)
	case GeneralOutlined:
		obj.Outlined.calcPercentageTotal(periods)
	case GeneralCalculated:
		obj.Calculated.calcPercentageTotal(periods)
	case GeneralEstimated:
		obj.Estimated.calcPercentageTotal(periods)
	}
}

func (obj Summary) getSummaryPeriodTotal(period measures.PeriodKey) float64 {
	methods := []GeneralEstimateMethod{GeneralReal, GeneralAdjusted, GeneralOutlined, GeneralCalculated, GeneralEstimated}
	var total float64
	for _, m := range methods {
		summaryPeriodValue := obj.getSummaryItemPeriod(m, period)
		total += summaryPeriodValue
	}

	return total
}

type SummaryItem struct {
	Total float64 `json:"total"`
	P1    float64 `json:"p1"`
	P2    float64 `json:"p2"`
	P3    float64 `json:"p3"`
	P4    float64 `json:"p4"`
	P5    float64 `json:"p5"`
	P6    float64 `json:"p6"`
}

func (obj *SummaryItem) addTotalValue(value float64) {
	obj.Total += value
}

func (obj SummaryItem) getPeriod(period measures.PeriodKey) float64 {
	switch period {
	case measures.P1:
		return obj.P1
	case measures.P2:
		return obj.P2
	case measures.P3:
		return obj.P3
	case measures.P4:
		return obj.P4
	case measures.P5:
		return obj.P5
	case measures.P6:
		return obj.P6
	default:
		return 0
	}
}

func (obj *SummaryItem) setPeriodValue(period measures.PeriodKey, value float64) {
	switch period {
	case measures.P1:
		obj.P1 = value
	case measures.P2:
		obj.P2 = value
	case measures.P3:
		obj.P3 = value
	case measures.P4:
		obj.P4 = value
	case measures.P5:
		obj.P5 = value
	case measures.P6:
		obj.P6 = value
	}
}

func (obj *SummaryItem) calcPercentagePeriod(period measures.PeriodKey, total float64) {
	if total == 0 {
		total = 1
	}

	periodValue := obj.getPeriod(period)
	percentage := math.Round((periodValue/total*100)*100) / 100
	obj.addTotalValue(percentage)
	obj.setPeriodValue(period, percentage)
}

func (obj *SummaryItem) calcPercentageTotal(periods []measures.PeriodKey) {
	totalPeriods := float64(len(periods))

	if totalPeriods == 0 {
		totalPeriods = 1
	}
	percentage := math.Round((obj.Total/totalPeriods)*100) / 100
	obj.Total = percentage

}

func (obj *SummaryItem) addConsum(period measures.PeriodKey, toSum float64) {
	periodValue := obj.getPeriod(period)
	obj.setPeriodValue(period, periodValue+toSum)
}

func (obj *SummaryItem) sumConsumPeriods(periods []measures.PeriodKey) {
	var total float64
	for _, period := range periods {
		periodValue := obj.getPeriod(period)
		total += periodValue
	}

	obj.Total = total
}

type CalendarCurve struct {
	Date   string                `json:"date"`
	Status GeneralEstimateMethod `json:"status"`
}

type Curve struct {
	Date   string   `json:"date"`
	Values []Values `json:"values"`
}

type Values struct {
	Hour   string                `json:"date"`
	Status GeneralEstimateMethod `json:"status"`
	AI     float64               `json:"ai"`
	AE     float64               `json:"ae"`
	AiAuto *float64              `json:"ai_auto"`
	AeAuto *float64              `json:"ae_auto"`
	R1     float64               `json:"r1"`
	R2     float64               `json:"r2"`
	R3     float64               `json:"r3"`
	R4     float64               `json:"r4"`
}

type Balance struct {
	Origin BalanceOriginType     `json:"origin"`
	Method GeneralEstimateMethod `json:"method"`
	P0     *BalancePeriod        `json:"p0"`
	P1     *BalancePeriod        `json:"p1"`
	P2     *BalancePeriod        `json:"p2"`
	P3     *BalancePeriod        `json:"p3"`
	P4     *BalancePeriod        `json:"p4"`
	P5     *BalancePeriod        `json:"p5"`
	P6     *BalancePeriod        `json:"p6"`
}

func (obj *Balance) setBalancePeriod(balancePeriod BalancePeriod, periodKey measures.PeriodKey) {
	switch periodKey {
	case measures.P0:
		obj.P0 = &balancePeriod
	case measures.P1:
		obj.P1 = &balancePeriod
	case measures.P2:
		obj.P2 = &balancePeriod
	case measures.P3:
		obj.P3 = &balancePeriod
	case measures.P4:
		obj.P4 = &balancePeriod
	case measures.P5:
		obj.P5 = &balancePeriod
	case measures.P6:
		obj.P6 = &balancePeriod
	}
}

func (obj *Balance) SetMetadata(origin BalanceOriginType, method GeneralEstimateMethod) {
	obj.Origin = origin
	obj.Method = method
}

type BalancePeriod struct {
	AE            float64               `json:"ae"`
	BalanceTypeAE GeneralEstimateMethod `json:"balance_type_ae"`
	AI            float64               `json:"ai"`
	BalanceTypeAI GeneralEstimateMethod `json:"balance_type_ai"`
	R1            float64               `json:"r1"`
	BalanceTypeR1 GeneralEstimateMethod `json:"balance_type_r1"`
	R2            float64               `json:"r2"`
	BalanceTypeR2 GeneralEstimateMethod `json:"balance_type_r2"`
	R3            float64               `json:"r3"`
	BalanceTypeR3 GeneralEstimateMethod `json:"balance_type_r3"`
	R4            float64               `json:"r4"`
	BalanceTypeR4 GeneralEstimateMethod `json:"balance_type_r4"`
}

func NewFiscalBillingMeasuresDashboard(bM, LbM BillingMeasure) FiscalBillingMeasuresDashboard {
	return FiscalBillingMeasuresDashboard{
		Id:                 bM.Id,
		Cups:               bM.CUPS,
		Type:               bM.MeterType,
		StartDate:          bM.InitDate.Format("02-01-2006"),
		EndDate:            bM.EndDate.Format("02-01-2006"),
		LastMvDate:         LbM.EndDate.Format("02-01-2006 15:04"),
		PrincipalMagnitude: bM.GetPrincipalMagnitude(),
		Status:             bM.Status,
		Periods:            bM.Periods,
		Magnitudes:         bM.Magnitudes,
		Summary:            Summary{},
		CalendarCurve:      nil,
		ExecutionSummary:   bM.ExecutionSummary,
		Balance: Balance{
			Origin: bM.GetBalanceOrigin(measures.P0, bM.GetPrincipalMagnitude()),
			Method: bM.GetBalanceGeneralType(measures.P0, bM.GetPrincipalMagnitude()),
		},
		GraphHistory: bM.GraphHistory,
		Curve:        nil,
	}
}

func NewBillingMeasureDashboardResumeClosure(bM BillingMeasureResumeClosureResponse) BillingMeasureDashboardResumeClosure {
	toReturn := BillingMeasureDashboardResumeClosure{
		ActualReadingClose:   newMetaDataReadingCloseResume(bM, "actual"),
		PreviousReadingClose: newMetaDataReadingCloseResume(bM, "previous"),
	}
	return toReturn
}

type BillingMeasureResumeClosureResponse struct {
	ActualReadingClose   measures.DailyReadingClosure `bson:"actual_reading_closure"`
	PreviousReadingClose measures.DailyReadingClosure `bson:"previous_reading_closure"`
	Periods              []measures.PeriodKey         `bson:"periods"`
	Magnitudes           []measures.Magnitude         `bson:"magnitudes"`
}

type BillingMeasureDashboardResumeClosure struct {
	ActualReadingClose   ReadingClosureResume `json:"actual_reading_close"`
	PreviousReadingClose ReadingClosureResume `json:"previous_reading_close"`
}
type ReadingClosureResume struct {
	MeterSerialNumber string             `json:"meter_serial_number"`
	ClosureType       string             `json:"closure_type"`
	Origin            string             `json:"origin"`
	InitDate          string             `json:"init_date"`
	EndDate           string             `json:"end_date"`
	P0                PeriodsByMagnitude `json:"p0"`
	P1                PeriodsByMagnitude `json:"p1"`
	P2                PeriodsByMagnitude `json:"p2"`
	P3                PeriodsByMagnitude `json:"p3"`
	P4                PeriodsByMagnitude `json:"p4"`
	P5                PeriodsByMagnitude `json:"p5"`
	P6                PeriodsByMagnitude `json:"p6"`
}
type PeriodsByMagnitude struct {
	AE PeriodReadingAndConsum `json:"ae"`
	AI PeriodReadingAndConsum `json:"ai"`
	R1 PeriodReadingAndConsum `json:"r1"`
	R2 PeriodReadingAndConsum `json:"r2"`
	R3 PeriodReadingAndConsum `json:"r3"`
	R4 PeriodReadingAndConsum `json:"r4"`
}
type PeriodReadingAndConsum struct {
	Reading float64 `json:"reading"`
	Consum  float64 `json:"consum"`
}

func newMetaDataReadingCloseResume(billingMeasure BillingMeasureResumeClosureResponse, flag string) ReadingClosureResume {
	var billingClose measures.DailyReadingClosure
	var initDate, endDate string

	if flag == "actual" {
		billingClose = billingMeasure.ActualReadingClose
	}
	if flag == "previous" {
		billingClose = billingMeasure.PreviousReadingClose
	}

	initDate = billingClose.InitDate.Format("2006-01-02")
	endDate = billingClose.EndDate.Format("2006-01-02")

	toReturn := ReadingClosureResume{
		MeterSerialNumber: billingClose.MeterSerialNumber,
		ClosureType:       string(billingClose.ClosureType),
		Origin:            string(billingClose.Origin),
		InitDate:          initDate,
		EndDate:           endDate,
		P0:                setPeriodsByMagnitude(billingClose.GetCalendarPeriod(measures.P0)),
		P1:                setPeriodsByMagnitude(billingClose.GetCalendarPeriod(measures.P1)),
		P2:                setPeriodsByMagnitude(billingClose.GetCalendarPeriod(measures.P2)),
		P3:                setPeriodsByMagnitude(billingClose.GetCalendarPeriod(measures.P3)),
		P4:                setPeriodsByMagnitude(billingClose.GetCalendarPeriod(measures.P4)),
		P5:                setPeriodsByMagnitude(billingClose.GetCalendarPeriod(measures.P5)),
		P6:                setPeriodsByMagnitude(billingClose.GetCalendarPeriod(measures.P6)),
	}
	return toReturn
}

func setPeriodsByMagnitude(billingCalendarPeriod *measures.DailyReadingClosureCalendarPeriod) PeriodsByMagnitude {
	if billingCalendarPeriod == nil {
		return PeriodsByMagnitude{}
	}
	return PeriodsByMagnitude{
		AE: PeriodReadingAndConsum{
			Reading: billingCalendarPeriod.Values.AE,
			Consum:  0,
		},
		AI: PeriodReadingAndConsum{
			Reading: billingCalendarPeriod.Values.AI,
			Consum:  0,
		},
		R1: PeriodReadingAndConsum{
			Reading: billingCalendarPeriod.Values.R1,
			Consum:  0,
		},
		R2: PeriodReadingAndConsum{
			Reading: billingCalendarPeriod.Values.R2,
			Consum:  0,
		},
		R3: PeriodReadingAndConsum{
			Reading: billingCalendarPeriod.Values.R3,
			Consum:  0,
		},
		R4: PeriodReadingAndConsum{
			Reading: billingCalendarPeriod.Values.R4,
			Consum:  0,
		},
	}
}

type FiscalMeasureSummary struct {
	MeterType     measures.MeterType   `json:"meter_type"`
	BalanceType   BalanceTypeSummary   `json:"balance_type"`
	BalanceOrigin BalanceOriginSummary `json:"balance_origin"`
	CurveType     CurveTypeSummary     `json:"curve_type"`
	CurveStatus   CurveStatusSummary   `json:"curve_status"`
	Total         int                  `json:"total"`
}
type TypeSummary struct {
	Real       int `json:"real"`
	Calculated int `json:"calculated"`
	Estimated  int `json:"estimated"`
}

type BalanceTypeSummary struct {
	TypeSummary
}

type BalanceOriginSummary struct {
	Monthly   int `json:"monthly"`
	Daily     int `json:"daily"`
	Other     int `json:"other"`
	NoClosure int `json:"no_closure"`
}

type CurveTypeSummary struct {
	TypeSummary
	Adjusted int `json:"adjusted"`
	Outlined int `json:"outlined"`
}

type CurveStatusSummary struct {
	Completed    int `json:"complete"`
	NotCompleted int `json:"not_completed"`
	Absent       int `json:"absent"`
}

type GroupFiscalMeasureSummaryQuery struct {
	DistributorId string
	MeterType     measures.MeterType
	StartDate     time.Time
	EndDate       time.Time
}

type QueryBillingMeasuresTax struct {
	Offset        *int      `json:"offset,omitempty"`
	Limit         int       `json:"limit"`
	DistributorId string    `json:"distributor_id"`
	MeasureType   string    `json:"measure_type"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
}

type BillingMeasuresTax struct {
	Cups             string           `json:"cups" bson:"cups"`
	DistributorId    string           `json:"distributor_id" bson:"distributor_id"`
	StartDate        time.Time        `json:"init_date" bson:"init_date"`
	EndDate          time.Time        `json:"end_date" bson:"end_date"`
	ExecutionSummary ExecutionSummary `json:"execution_summary" bson:"execution_summary"`
}

type BillingMeasuresTaxResult struct {
	Data  []BillingMeasuresTax `bson:"data"`
	Count int                  `bson:"total"`
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=BillingMeasuresDashboardRepository
type BillingMeasuresDashboardRepository interface {
	SearchFiscalBillingMeasures(ctx context.Context, cups string, distributorId string, startDate time.Time, endDate time.Time) ([]BillingMeasure, error)
	SearchLastBillingMeasures(ctx context.Context, cups string, distributorId string) (BillingMeasure, error)
	SearchBillingMeasureClosureResume(ctx context.Context, billingMeasureID string) (BillingMeasureResumeClosureResponse, error)
	GroupFiscalMeasureSummary(ctx context.Context, query GroupFiscalMeasureSummaryQuery) (FiscalMeasureSummary, error)
	GetBillingMeasuresTax(ctx context.Context, query QueryBillingMeasuresTax) (BillingMeasuresTaxResult, error)
}
