package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

type OriginType string

const (
	ScheduledOrigin OriginType = "SCHEDULED"
	ExternalOrigin  OriginType = "EXTERNAL"
)

type Status string

const (
	Pending     Status = "PENDING"
	Calculating Status = "CALCULATING"
	Calculated  Status = "CALCULATED"
	PendingCau  Status = "PENDING_CAU"
	ReadyToBill Status = "READY_TO_BILL"
	Billed      Status = "BILLED"
	Overrided   Status = "OVERRIDED"
	Cancelled   Status = "CANCELLED"
	Supervision Status = "SUPERVISION"
)

type DescriptionStatus string

const (
	NoMatchLenLoadCurve    DescriptionStatus = "NO_MATCH_LEN_BILLING_LOAD_CURVE"
	NotValidReadingClosure DescriptionStatus = "NOT_VALID_READING_CLOSURE"
	NoLastBillingMeasure   DescriptionStatus = "NO_LAST_BILLING_MEASURE"
)

type BalanceType string

type GeneralEstimateMethod string

const (
	GeneralReal       GeneralEstimateMethod = "REAL"
	GeneralOutlined   GeneralEstimateMethod = "OUTLINED"
	GeneralAdjusted   GeneralEstimateMethod = "ADJUSTED"
	GeneralCalculated GeneralEstimateMethod = "CALCULATED"
	GeneralEstimated  GeneralEstimateMethod = "ESTIMATED"
)

var GeneralMethodValue = map[GeneralEstimateMethod]int{
	GeneralEstimated:  0,
	GeneralCalculated: 1,
	GeneralAdjusted:   2,
	GeneralOutlined:   3,
	GeneralReal:       4,
}

func (g GeneralEstimateMethod) Value(methodValues map[GeneralEstimateMethod]int) int {
	value, _ := methodValues[g]

	return value
}

func (g GeneralEstimateMethod) Compare(compareMethods map[GeneralEstimateMethod]int, b GeneralEstimateMethod) bool {
	return g.Value(compareMethods) < b.Value(compareMethods)
}

var RealBalanceTypes = []BalanceType{
	RealBalance,
	RealByRemoteRead,
	RealByAbsLocalRead,
	RealByAutoRead,
	RealValidMeasure,
	RealMeasureAdjustment,
}

const (
	EstimateBalance                               BalanceType = "ESTIMATE"
	RealBalance                                   BalanceType = "REAL"
	RealByRemoteRead                              BalanceType = "REAL_BY_REMOTE_READ"
	RealByAbsLocalRead                            BalanceType = "REAL_BY_ABSOLUTE_LOCAL_READ"
	RealByAutoRead                                BalanceType = "REAL_BY_AUTO_READ"
	CalculatedBalance                             BalanceType = "CALCULATED"
	CalculatedByCloseSum                          BalanceType = "CALCULATED_BY_CLOSE_SUM"
	CalculatedByCloseBalance                      BalanceType = "CALCULATED_BY_CLOSE_AND_BALANCE"
	Adjustment                                    BalanceType = "ADJUSTED"
	ProfileMeasure                                BalanceType = "PROFILE_MEASURE"
	ProfileMeasureAutoReading                     BalanceType = "PROFILE_MEASURE_AUTO_READING"
	EstimateHistoryProfile                        BalanceType = "ESTIMATE_HISTORY_PROFILE"
	EstimateFactorUsed                            BalanceType = "ESTIMATE_FACTOR_USED"
	RealValidMeasure                              BalanceType = "REAL_VALID_MEASURE"
	EstimatedContractPower                        BalanceType = "ESTIMATED_CONTRACT_POWER"
	EstimateOnlyHistoric                          BalanceType = "ESTIMATE_ONLY_HISTORIC"
	FirmMainConfig                                BalanceType = "FIRM_MAIN_CONFIG"
	FirmRedundantConfig                           BalanceType = "FIRM_REDUNDANT_CONFIG"
	FirmReceiptConfig                             BalanceType = "FIRM_RECEIPT_CONFIG"
	EstimatedByFlatProfile                        BalanceType = "ESTIMATED_BY_FLAT_PROFILE"
	RealMeasureAdjustment                         BalanceType = "REAL_MEASURE_ADJUSTMENT"
	PowerUseFactor                                BalanceType = "POWER_USE_FACTOR"
	EstimateByHistoricLastYear                    BalanceType = "ESTIMATE_BY_HISTORIC_LAST_YEAR"
	EstimateByHistoricConsumLastYear              BalanceType = "ESTIMATE_BY_HISTORIC_CONSUM_LAST_YEAR"
	EstimateHistoricMainMeasurePoint              BalanceType = "ESTIMATE_HISTORIC_MAIN_MEASURE_POINT"
	EstimateWhosePenaltiesForClientsTypeOneAndTwo BalanceType = "ESTIMATE_WHOSE_PENALTIES_FOR_CLIENTS_TYPE_ONE_AND_TWO"
	ObtainedByCurve                               BalanceType = "OBTAINED_BY_CURVE"
	Outlined                                      BalanceType = "OUTLINED"
	AutoOutlined                                  BalanceType = "AUTO_OUTLINED"
	EstimatedByHistoricOutlinedConsum             BalanceType = "ESTIMATED_BY_HISTORIC_OUTLINE_CONSUM"
	EstimatedByOutlinedFactor                     BalanceType = "ESTIMATED_BY_OUTLINED_FACTOR"
	EstimateHistoricMainMeasurePointBalance       BalanceType = "ESTIMATE_HISTORIC_MAIN_MEASURE_POINT_BALANCE"
)

type BalanceMeasureType string

const (
	FirmBalanceMeasure        BalanceMeasureType = "FIRM"
	ProvisionalBalanceMeasure BalanceMeasureType = "PROVISIONAL"
)

type BalanceOriginType string

const (
	RemoteBalanceOrigin BalanceOriginType = "REMOTE"
	AutoBalanceOrigin   BalanceOriginType = "AUTO"
	EstimateOrigin      BalanceOriginType = "ESTIMATE"
	TlgOrigin           BalanceOriginType = "TLG"
	TlmOrigin           BalanceOriginType = "TLM"
	LocalOrigin         BalanceOriginType = "Local"
)

type CoefficientType string

const (
	CoefficientA CoefficientType = "A"
	CoefficientB CoefficientType = "B"
	CoefficientC CoefficientType = "C"
	CoefficientD CoefficientType = "D"
)

type AtrVsCurve struct {
	P0 *measures.Values `json:"p0" bson:"p0,omitempty"`
	P1 *measures.Values `json:"p1" bson:"p1,omitempty"`
	P2 *measures.Values `json:"p2" bson:"p2,omitempty"`
	P3 *measures.Values `json:"p3" bson:"p3,omitempty"`
	P4 *measures.Values `json:"p4" bson:"p4,omitempty"`
	P5 *measures.Values `json:"p5" bson:"p5,omitempty"`
	P6 *measures.Values `json:"p6" bson:"p6,omitempty"`
}

type ExecutionSummary struct {
	BalanceType   GeneralEstimateMethod `json:"balance_type" bson:"balance_type"`
	CurveType     GeneralEstimateMethod `json:"curve_type" bson:"curve_type"`
	CurveStatus   measures.QualityCode  `json:"curve_status" bson:"curve_status"`
	BalanceOrigin measures.ClosureType  `json:"balance_origin" bson:"balance_origin"`
}

func (e *ExecutionSummary) setCurveStatus(status measures.QualityCode) {
	e.CurveStatus = status
}

func (e *ExecutionSummary) setCurveType(method GeneralEstimateMethod) {
	if e.CurveType == "" {
		e.CurveType = method
		return
	}

	if !e.CurveType.Compare(GeneralMethodValue, method) {
		e.CurveType = method
	}
}

func (e *ExecutionSummary) setBalanceOrigin(closureType measures.ClosureType) {
	e.BalanceOrigin = closureType
}

func (e *ExecutionSummary) setBalanceType(method GeneralEstimateMethod) {
	if e.BalanceType == "" {
		e.BalanceType = method
		return
	}

	if !e.BalanceType.Compare(GeneralMethodValue, method) {
		e.BalanceType = method
	}
}

type BillingMeasure struct {
	Id                     string                       `json:"id" bson:"_id"`
	DistributorCode        string                       `json:"distributor_code" bson:"distributor_code"`
	DistributorID          string                       `json:"distributor_id" bson:"distributor_id"`
	Coefficient            CoefficientType              `json:"coefficient" bson:"coefficient"`
	CUPS                   string                       `json:"cups" bson:"cups"`
	PointType              measures.PointType           `json:"point_type" bson:"point_type"`
	Inaccessible           bool                         `json:"inaccessible" bson:"inaccessible"`
	RegisterType           measures.RegisterType        `json:"register_type" bson:"register_type"`
	ReadingType            measures.Type                `json:"reading_type" bson:"reading_type"`
	Technology             string                       `json:"technology" bson:"technology"`
	GenerationDate         time.Time                    `json:"generation_date" bson:"generation_date"`
	EndDate                time.Time                    `json:"end_date" bson:"end_date"`
	InitDate               time.Time                    `json:"init_date" bson:"init_date"`
	Version                string                       `json:"version" bson:"version"`
	Origin                 OriginType                   `json:"origin" bson:"origin"`
	Status                 Status                       `json:"status" bson:"status"`
	DescriptionStatus      *DescriptionStatus           `json:"description_status" bson:"description_status"`
	PreviousReadingClosure measures.DailyReadingClosure `json:"previous_reading_closure" bson:"previous_reading_closure"`
	ActualReadingClosure   measures.DailyReadingClosure `json:"actual_reading_closure" bson:"actual_reading_closure"`
	BillingBalance         BillingBalance               `json:"billing_balance" bson:"billing_balance"`
	BillingLoadCurve       []BillingLoadCurve           `json:"billing_load_curve" bson:"billing_load_curve"`
	AtrVsCurve             AtrVsCurve                   `json:"atr_vs_curve" bson:"atr_vs_curve"`
	Periods                []measures.PeriodKey         `json:"periods" bson:"periods"`
	Magnitudes             []measures.Magnitude         `json:"magnitudes" bson:"magnitudes"`
	TariffID               string                       `json:"tariff_id" bson:"tariff_id"`
	CalendarCode           string                       `json:"calendar_code" bson:"calendar_code"`
	P1Demand               float64                      `json:"p1_demand" bson:"p1_demand"`
	P2Demand               float64                      `json:"p2_demand" bson:"p2_demand"`
	P3Demand               float64                      `json:"p3_demand" bson:"p3_demand"`
	P4Demand               float64                      `json:"p4_demand" bson:"p4_demand"`
	P5Demand               float64                      `json:"p5_demand" bson:"p5_demand"`
	P6Demand               float64                      `json:"p6_demand" bson:"p6_demand"`
	GraphHistory           map[string]*Graph            `json:"graph_history" bson:"graph_history"`
	ExecutionSummary       ExecutionSummary             `json:"execution_summary" bson:"execution_summary"`
	MeterType              measures.MeterType           `json:"meter_type" bson:"meter_type"`
}

func (b BillingMeasure) GetPeriodsString() []string {
	return utils.MapSlice(b.Periods, func(item measures.PeriodKey) string {
		return string(item)
	})
}

func (b BillingMeasure) GetMagnitudesString() []string {
	return utils.MapSlice(b.Magnitudes, func(item measures.Magnitude) string {
		return string(item)
	})
}

func (b BillingMeasure) isPointType(pointTypeOptions []measures.PointType) bool {
	for _, pointTypeOp := range pointTypeOptions {
		if pointTypeOp == b.PointType {
			return true
		}
	}
	return false
}

func (b BillingMeasure) IsLoadCurveOrigin(index int, originTypeOptions []measures.OriginType) bool {
	if index > len(b.BillingLoadCurve) {
		return false
	}
	curve := b.BillingLoadCurve[index]
	originType := curve.getOriginType()
	for _, option := range originTypeOptions {
		if originType == option {
			return true
		}
	}
	return false
}

func (b *BillingMeasure) SetLoadCurveEstimatedCode(index int, magnitude measures.Magnitude, estimateCode int) {
	if index > len(b.BillingLoadCurve) {
		return
	}
	curve := &b.BillingLoadCurve[index]
	curve.setEstimatedCode(magnitude, estimateCode)
}

func (b BillingMeasure) GetPrincipalMagnitude() measures.Magnitude {

	readingClosure := b.ActualReadingClosure

	if readingClosure.EndDate.IsZero() {
		readingClosure = b.PreviousReadingClosure
	}

	if readingClosure.ServiceType == string(measures.DcServiceType) {
		return measures.AI
	}
	if readingClosure.ServiceType == string(measures.GdServiceType) {
		return measures.AE
	}

	return ""
}

func (b BillingMeasure) GetLoadCurveEstimatedCode(index int, magnitude measures.Magnitude) (int, error) {
	if index > len(b.BillingLoadCurve) {
		return 0, errors.New("index loadCurve not found")
	}
	curve := b.BillingLoadCurve[index]
	estimatedCode, err := curve.getEstimateCode(magnitude)
	if err != nil {
		return 0, errors.New("unexpected magnitude")
	}
	return estimatedCode, nil
}

func (b BillingMeasure) GetLoadCurveGeneralMagnitudeValue(index int, magnitude measures.Magnitude) (float64, error) {
	if index > len(b.BillingLoadCurve) {
		return 0.0, errors.New("index loadCurve not fount")
	}
	curve := b.BillingLoadCurve[index]
	estimatedGeneralMethod, err := curve.getMagnitude(magnitude)
	if err != nil {
		return 0.0, err
	}
	return estimatedGeneralMethod, nil
}

func (b BillingMeasure) GetActualReadingClose() measures.DailyReadingClosure {
	return b.ActualReadingClosure
}
func (b BillingMeasure) GetPreviousReadingClose() measures.DailyReadingClosure {
	return b.PreviousReadingClosure
}

func (b BillingMeasure) GetLoadCurvePeriod(index int) (measures.PeriodKey, error) {
	if index > len(b.BillingLoadCurve) {
		return measures.P1, errors.New("unexpected index value")
	}
	curve := b.BillingLoadCurve[index]
	return curve.getPeriod(), nil
}

func (b BillingMeasure) GetLoadCurveGeneralEstimatedMethod(index int, magnitude measures.Magnitude) (GeneralEstimateMethod, error) {
	if index > len(b.BillingLoadCurve) {
		return GeneralOutlined, errors.New("index loadCurve not fount")
	}
	curve := b.BillingLoadCurve[index]
	estimatedGeneralMethod, err := curve.getGeneralEstimatedMethod(magnitude)
	if err != nil {
		return GeneralOutlined, err
	}
	return estimatedGeneralMethod, nil
}

func (b *BillingMeasure) SetLoadCurveGeneralEstimatedMethod(index int, magnitude measures.Magnitude, method GeneralEstimateMethod) {
	if index > len(b.BillingLoadCurve) {
		return
	}
	curve := &b.BillingLoadCurve[index]
	curve.setGeneralEstimatedMethod(magnitude, method)
}

func (b *BillingMeasure) SetLoadCurveEstimatedMethod(index int, magnitude measures.Magnitude, estimatedMethod BalanceType) {
	if index > len(b.BillingLoadCurve) {
		return
	}
	curve := &b.BillingLoadCurve[index]
	curve.setEstimatedMethod(magnitude, estimatedMethod)
}

func (b *BillingMeasure) SetLoadCurveMeasureType(index int, magnitude measures.Magnitude, measureType BalanceMeasureType) {
	if index > len(b.BillingLoadCurve) {
		return
	}
	curve := &b.BillingLoadCurve[index]
	curve.setMeasureType(magnitude, measureType)
}

func (b *BillingMeasure) SumLoadCurvePeriodMagnitude(index int, magnitude measures.Magnitude, value float64) {
	if index > len(b.BillingLoadCurve) {
		return
	}
	curve := &b.BillingLoadCurve[index]
	curve.SumMagnitude(magnitude, value)
}

func (b *BillingMeasure) GetLoadCurvePeriodMagnitude(index int, magnitude measures.Magnitude) float64 {
	if index > len(b.BillingLoadCurve) {
		return .0
	}

	return b.BillingLoadCurve[index].GetMagnitude(magnitude)
}

func (b *BillingMeasure) SetLoadCurvePeriodMagnitude(index int, magnitude measures.Magnitude, value float64) {
	if index > len(b.BillingLoadCurve) {
		return
	}
	curve := &b.BillingLoadCurve[index]
	curve.SetMagnitude(magnitude, value)
}

func (b *BillingMeasure) SetLoadCurveSelfConsumptionMagnitude(index int, magnitude SelfConsumptionMagnitude, value *float64) {
	if index > len(b.BillingLoadCurve) {
		return
	}
	curve := &b.BillingLoadCurve[index]
	curve.SetSelfConsumptionMagnitude(magnitude, value)
}

func (b *BillingMeasure) SumLoadCurveSelfConsumptionMagnitude(index int, magnitude SelfConsumptionMagnitude, value *float64) {
	if index > len(b.BillingLoadCurve) {
		return
	}
	curve := &b.BillingLoadCurve[index]
	curve.SumSelfConsumptionMagnitude(magnitude, value)
}

func (b *BillingMeasure) GetLoadCurveOrigin(index int) measures.OriginType {
	if index > len(b.BillingLoadCurve) {
		return ""
	}
	return b.BillingLoadCurve[index].Origin
}

func (b *BillingMeasure) GetLoadCurveMeasurePointType(index int) measures.MeasurePointType {
	if index > len(b.BillingLoadCurve) {
		return ""
	}
	return b.BillingLoadCurve[index].MeasurePointType
}

func (b *BillingMeasure) GetLoadCurveEquipment(index int) measures.EquipmentType {
	if index > len(b.BillingLoadCurve) {
		return ""
	}
	return b.BillingLoadCurve[index].Equipment
}

func (b *BillingMeasure) GetBalancePeriod(period measures.PeriodKey) *BillingBalancePeriod {
	switch period {
	case measures.P0:
		return b.BillingBalance.P0
	case measures.P1:
		return b.BillingBalance.P1
	case measures.P2:
		return b.BillingBalance.P2
	case measures.P3:
		return b.BillingBalance.P3
	case measures.P4:
		return b.BillingBalance.P4
	case measures.P5:
		return b.BillingBalance.P5
	case measures.P6:
		return b.BillingBalance.P6
	default:
		return nil
	}
}

func (b *BillingMeasure) GetBalanceOrigin(period measures.PeriodKey, magnitude measures.Magnitude) BalanceOriginType {
	balance := b.GetBalancePeriod(period)

	if balance == nil {
		return ""
	}

	return balance.getBalanceOrigin(magnitude)
}

func (b *BillingMeasure) GetBalanceGeneralType(period measures.PeriodKey, magnitude measures.Magnitude) GeneralEstimateMethod {
	balance := b.GetBalancePeriod(period)

	if balance == nil {
		return ""
	}

	return balance.getBalanceGeneralType(magnitude)
}

func (b *BillingMeasure) SetBalanceGeneralType(period measures.PeriodKey, magnitude measures.Magnitude, generalType GeneralEstimateMethod) {
	balance := b.GetBalancePeriod(period)

	if balance == nil {
		return
	}

	balance.setBalanceGeneralType(magnitude, generalType)
}

func (b *BillingMeasure) SetWorstBalanceGeneralType(magnitude measures.Magnitude, method GeneralEstimateMethod) {
	balanceMethod := b.GetBalanceGeneralType(measures.P0, magnitude)

	if balanceMethod == "" || !balanceMethod.Compare(GeneralMethodValue, method) {
		b.SetBalanceGeneralType(measures.P0, magnitude, method)
	}
}

func (b *BillingMeasure) GetBalanceStatus(period measures.PeriodKey, magnitude measures.Magnitude) measures.Status {
	balance := b.GetBalancePeriod(period)

	if balance == nil {
		return ""
	}

	return balance.GetStatus(magnitude)
}

func (b *BillingMeasure) SetBalanceStatus(period measures.PeriodKey, magnitude measures.Magnitude, status measures.Status) {
	balance := b.GetBalancePeriod(period)

	if balance == nil {
		return
	}

	balance.SetStatus(magnitude, status)
}

func (b *BillingMeasure) SetBalanceEstimateCode(period measures.PeriodKey, magnitude measures.Magnitude, estimateCode int) {
	balance := b.GetBalancePeriod(period)

	if balance == nil {
		return
	}
	balance.setEstimateCode(magnitude, estimateCode)
}

func (b *BillingMeasure) GetBalanceEstimateCode(period measures.PeriodKey, magnitude measures.Magnitude) int {
	balance := b.GetBalancePeriod(period)

	if balance == nil {
		return 0
	}
	return balance.getEstimateCode(magnitude)
}

func (b *BillingMeasure) SetBalanceType(period measures.PeriodKey, magnitude measures.Magnitude, balanceType BalanceType) {
	balance := b.GetBalancePeriod(period)

	if balance == nil {
		return
	}
	balance.setBalanceType(magnitude, balanceType)
}

func (b *BillingMeasure) SetBalanceMeasureType(period measures.PeriodKey, magnitude measures.Magnitude, measureType BalanceMeasureType) {
	balance := b.GetBalancePeriod(period)

	if balance == nil {
		return
	}
	balance.setBalanceMeasureType(magnitude, measureType)

}

func (b *BillingMeasure) SetBalanceOrigin(period measures.PeriodKey, magnitude measures.Magnitude, val BalanceOriginType) {
	balance := b.GetBalancePeriod(period)

	if balance == nil {
		return
	}
	balance.setBalanceOrigin(magnitude, val)
}

func (b *BillingMeasure) SumBalancePeriodMagnitude(period measures.PeriodKey, magnitude measures.Magnitude, value float64) {

	balance := b.GetBalancePeriod(period)

	if balance == nil {
		return
	}

	balance.sumPeriodMagnitude(magnitude, value)
}

func (b *BillingMeasure) SetBalancePeriodMagnitude(period measures.PeriodKey, magnitude measures.Magnitude, value float64) {

	balance := b.GetBalancePeriod(period)

	if balance == nil {
		return
	}

	switch magnitude {
	case measures.AI:
		balance.AI = value
	case measures.AE:
		balance.AE = value
	case measures.R1:
		balance.R1 = value
	case measures.R2:
		balance.R2 = value
	case measures.R3:
		balance.R3 = value
	case measures.R4:
		balance.R4 = value
	}
}

func (b *BillingMeasure) GetBalancePeriodMagnitude(period measures.PeriodKey, magnitude measures.Magnitude) float64 {

	balance := b.GetBalancePeriod(period)

	if balance == nil {
		return .0
	}

	switch magnitude {
	case measures.AI:
		return balance.AI
	case measures.AE:
		return balance.AE
	case measures.R1:
		return balance.R1
	case measures.R2:
		return balance.R2
	case measures.R3:
		return balance.R3
	case measures.R4:
		return balance.R4
	default:
		return 0
	}
}

func (b *BillingMeasure) GetPeriodPowerDemanded(period measures.PeriodKey) float64 {
	switch period {
	case measures.P1:
		return b.P1Demand
	case measures.P2:
		return b.P2Demand
	case measures.P3:
		return b.P3Demand
	case measures.P4:
		return b.P4Demand
	case measures.P5:
		return b.P5Demand
	case measures.P6:
		return b.P6Demand
	default:
		return 0
	}
}

func (b BillingMeasure) IsValidBillingBalancePeriod(period measures.PeriodKey, magnitude measures.Magnitude) bool {

	var billingBalancePeriod *BillingBalancePeriod

	switch period {
	case measures.P0:
		billingBalancePeriod = b.BillingBalance.P0
	case measures.P1:
		billingBalancePeriod = b.BillingBalance.P1
	case measures.P2:
		billingBalancePeriod = b.BillingBalance.P2
	case measures.P3:
		billingBalancePeriod = b.BillingBalance.P3
	case measures.P4:
		billingBalancePeriod = b.BillingBalance.P4
	case measures.P5:
		billingBalancePeriod = b.BillingBalance.P5
	case measures.P6:
		billingBalancePeriod = b.BillingBalance.P6
	}

	if billingBalancePeriod == nil {
		return false
	}

	return billingBalancePeriod.GetStatus(magnitude) == measures.Valid
}

func (b BillingMeasure) IsAtrVsCurveValid(period measures.PeriodKey, magnitude measures.Magnitude) bool {
	var values *measures.Values
	switch period {
	case measures.P0:
		values = b.AtrVsCurve.P0
	case measures.P1:
		values = b.AtrVsCurve.P1
	case measures.P2:
		values = b.AtrVsCurve.P2
	case measures.P3:
		values = b.AtrVsCurve.P3
	case measures.P4:
		values = b.AtrVsCurve.P4
	case measures.P5:
		values = b.AtrVsCurve.P5
	case measures.P6:
		values = b.AtrVsCurve.P6
	}

	if values == nil {
		return false
	}

	return math.Abs(values.GetMagnitude(magnitude)) <= 1000
}

func (b *BillingMeasure) CalcAtrBalance() {
	b.BillingBalance.Origin = b.ActualReadingClosure.Origin
	b.BillingBalance.EndDate = b.ActualReadingClosure.EndDate

	previousMagnitudes := utils.NewSet(b.PreviousReadingClosure.Magnitudes)

	for _, period := range append(b.Periods, measures.P0) {
		for _, magnitude := range b.Magnitudes {

			var balancePeriod *BillingBalancePeriod
			var actualCalendarPeriod *measures.DailyReadingClosureCalendarPeriod
			var previousCalendarPeriod *measures.DailyReadingClosureCalendarPeriod

			switch period {
			case measures.P0:
				balancePeriod = b.BillingBalance.P0
				actualCalendarPeriod = b.ActualReadingClosure.CalendarPeriods.P0
				previousCalendarPeriod = b.PreviousReadingClosure.CalendarPeriods.P0
			case measures.P1:
				balancePeriod = b.BillingBalance.P1
				actualCalendarPeriod = b.ActualReadingClosure.CalendarPeriods.P1
				previousCalendarPeriod = b.PreviousReadingClosure.CalendarPeriods.P1
			case measures.P2:
				balancePeriod = b.BillingBalance.P2
				actualCalendarPeriod = b.ActualReadingClosure.CalendarPeriods.P2
				previousCalendarPeriod = b.PreviousReadingClosure.CalendarPeriods.P2
			case measures.P3:
				balancePeriod = b.BillingBalance.P3
				actualCalendarPeriod = b.ActualReadingClosure.CalendarPeriods.P3
				previousCalendarPeriod = b.PreviousReadingClosure.CalendarPeriods.P3
			case measures.P4:
				balancePeriod = b.BillingBalance.P4
				actualCalendarPeriod = b.ActualReadingClosure.CalendarPeriods.P4
				previousCalendarPeriod = b.PreviousReadingClosure.CalendarPeriods.P4
			case measures.P5:
				balancePeriod = b.BillingBalance.P5
				actualCalendarPeriod = b.ActualReadingClosure.CalendarPeriods.P5
				previousCalendarPeriod = b.PreviousReadingClosure.CalendarPeriods.P5
			case measures.P6:
				balancePeriod = b.BillingBalance.P6
				actualCalendarPeriod = b.ActualReadingClosure.CalendarPeriods.P6
				previousCalendarPeriod = b.PreviousReadingClosure.CalendarPeriods.P6
			}

			if b.ReadingType == measures.Incremental {
				balancePeriod.SetAtrBalance(actualCalendarPeriod, magnitude)
				continue
			}
			if !previousMagnitudes.Has(magnitude) {
				previousCalendarPeriod = nil
			}
			balancePeriod.CalcAtrBalance(actualCalendarPeriod, previousCalendarPeriod, magnitude)
		}
	}
}

func (balancePeriod *BillingBalancePeriod) CalcAtrBalance(actual, previous *measures.DailyReadingClosureCalendarPeriod, magnitude measures.Magnitude) {
	balancePeriod.SetStatus(magnitude, measures.Invalid)

	if previous == nil || actual == nil {
		return
	}

	balancePeriod.setPeriodMagnitude(magnitude, actual.GetMagnitude(magnitude)-previous.GetMagnitude(magnitude))
	if actual.ValidationStatus == measures.Valid && previous.ValidationStatus == measures.Valid {

		balancePeriod.SetStatus(magnitude, measures.Valid)
	}
}

func (balancePeriod *BillingBalancePeriod) SetAtrBalance(actual *measures.DailyReadingClosureCalendarPeriod, magnitude measures.Magnitude) {
	balancePeriod.SetStatus(magnitude, measures.Invalid)

	if actual == nil {
		return
	}

	balancePeriod.setPeriodMagnitude(magnitude, actual.GetMagnitude(magnitude))
	if actual.ValidationStatus == measures.Valid {
		balancePeriod.SetStatus(magnitude, measures.Valid)
	}
}

func (b BillingMeasure) IsCurvePeriodCompletedByPeriod(period measures.PeriodKey) bool {
	for _, curve := range b.BillingLoadCurve {
		if curve.Period == period && curve.Origin == measures.Filled {
			return false
		}
	}

	return true
}

func (b BillingMeasure) AreSomeCurveMeasureForPeriod(period measures.PeriodKey) bool {
	for _, curve := range b.BillingLoadCurve {
		if curve.Period == period && curve.Origin != measures.Filled {
			return true
		}
	}

	return false
}

func (b BillingMeasure) GetPeriods() []measures.PeriodKey {
	return append(b.Periods, measures.P0)
}

func (b *BillingMeasure) CalcAtrVsCurve() {
	for _, p := range b.GetPeriods() {
		var balance *BillingBalancePeriod
		var atrVsCurve *measures.Values
		switch p {
		case measures.P0:
			balance = b.BillingBalance.P0
			atrVsCurve = b.AtrVsCurve.P0
		case measures.P1:
			balance = b.BillingBalance.P1
			atrVsCurve = b.AtrVsCurve.P1
		case measures.P2:
			balance = b.BillingBalance.P2
			atrVsCurve = b.AtrVsCurve.P2
		case measures.P3:
			balance = b.BillingBalance.P3
			atrVsCurve = b.AtrVsCurve.P3
		case measures.P4:
			balance = b.BillingBalance.P4
			atrVsCurve = b.AtrVsCurve.P4
		case measures.P5:
			balance = b.BillingBalance.P5
			atrVsCurve = b.AtrVsCurve.P5
		case measures.P6:
			balance = b.BillingBalance.P6
			atrVsCurve = b.AtrVsCurve.P6
		}
		atrVsCurve.AI = balance.AI
		atrVsCurve.AE = balance.AE
		atrVsCurve.R1 = balance.R1
		atrVsCurve.R2 = balance.R2
		atrVsCurve.R3 = balance.R3
		atrVsCurve.R4 = balance.R4
	}

	validPeriods := utils.NewSet(b.GetPeriods())

	for _, curve := range b.BillingLoadCurve {

		if curve.Origin == measures.Filled || !validPeriods.Has(curve.Period) {
			continue
		}

		switch curve.Period {
		case measures.P1:
			{
				b.AtrVsCurve.P1.AI -= curve.AI
				b.AtrVsCurve.P1.AE -= curve.AE
				b.AtrVsCurve.P1.R1 -= curve.R1
				b.AtrVsCurve.P1.R2 -= curve.R2
				b.AtrVsCurve.P1.R3 -= curve.R3
				b.AtrVsCurve.P1.R4 -= curve.R4
			}
		case measures.P2:
			{
				b.AtrVsCurve.P2.AI -= curve.AI
				b.AtrVsCurve.P2.AE -= curve.AE
				b.AtrVsCurve.P2.R1 -= curve.R1
				b.AtrVsCurve.P2.R2 -= curve.R2
				b.AtrVsCurve.P2.R3 -= curve.R3
				b.AtrVsCurve.P2.R4 -= curve.R4
			}
		case measures.P3:
			{
				b.AtrVsCurve.P3.AI -= curve.AI
				b.AtrVsCurve.P3.AE -= curve.AE
				b.AtrVsCurve.P3.R1 -= curve.R1
				b.AtrVsCurve.P3.R2 -= curve.R2
				b.AtrVsCurve.P3.R3 -= curve.R3
				b.AtrVsCurve.P3.R4 -= curve.R4
			}
		case measures.P4:
			{
				b.AtrVsCurve.P4.AI -= curve.AI
				b.AtrVsCurve.P4.AE -= curve.AE
				b.AtrVsCurve.P4.R1 -= curve.R1
				b.AtrVsCurve.P4.R2 -= curve.R2
				b.AtrVsCurve.P4.R3 -= curve.R3
				b.AtrVsCurve.P4.R4 -= curve.R4
			}
		case measures.P5:
			{
				b.AtrVsCurve.P5.AI -= curve.AI
				b.AtrVsCurve.P5.AE -= curve.AE
				b.AtrVsCurve.P5.R1 -= curve.R1
				b.AtrVsCurve.P5.R2 -= curve.R2
				b.AtrVsCurve.P5.R3 -= curve.R3
				b.AtrVsCurve.P5.R4 -= curve.R4
			}
		case measures.P6:
			{
				b.AtrVsCurve.P6.AI -= curve.AI
				b.AtrVsCurve.P6.AE -= curve.AE
				b.AtrVsCurve.P6.R1 -= curve.R1
				b.AtrVsCurve.P6.R2 -= curve.R2
				b.AtrVsCurve.P6.R3 -= curve.R3
				b.AtrVsCurve.P6.R4 -= curve.R4
			}
		default:
			break
		}

		b.AtrVsCurve.P0.AI -= curve.AI
		b.AtrVsCurve.P0.AE -= curve.AE
		b.AtrVsCurve.P0.R1 -= curve.R1
		b.AtrVsCurve.P0.R2 -= curve.R2
		b.AtrVsCurve.P0.R3 -= curve.R3
		b.AtrVsCurve.P0.R4 -= curve.R4
	}
}

func (b *BillingMeasure) SetBillingLoadCurve(curve []BillingLoadCurve) {
	b.BillingLoadCurve = curve
}

func (b *BillingMeasure) SetActualReadingClosure(a measures.DailyReadingClosure) {
	b.ActualReadingClosure = a
}

func (b *BillingMeasure) SetPreviousReadingClosure(p measures.DailyReadingClosure) {
	b.PreviousReadingClosure = p
}

func (b *BillingMeasure) SetContractInfo(c measures.MeterConfig) {
	if b.ActualReadingClosure.DistributorID == "" {
		b.TariffID = c.ContractualSituations.TariffID
		b.CalendarCode = c.CalendarID
		b.P1Demand = c.ContractualSituations.P1Demand
		b.P2Demand = c.ContractualSituations.P2Demand
		b.P3Demand = c.ContractualSituations.P3Demand
		b.P4Demand = c.ContractualSituations.P4Demand
		b.P5Demand = c.ContractualSituations.P5Demand
		b.P6Demand = c.ContractualSituations.P6Demand
		return
	}

	b.TariffID = b.ActualReadingClosure.TariffId
	b.CalendarCode = b.ActualReadingClosure.CalendarCode
	b.P1Demand = b.ActualReadingClosure.P1Demand
	b.P2Demand = b.ActualReadingClosure.P2Demand
	b.P3Demand = b.ActualReadingClosure.P3Demand
	b.P4Demand = b.ActualReadingClosure.P4Demand
	b.P5Demand = b.ActualReadingClosure.P5Demand
	b.P6Demand = b.ActualReadingClosure.P6Demand
}

func (b *BillingMeasure) SetCoefficient(tariff clients.Tariffs) {
	b.Coefficient = CoefficientType(b.ActualReadingClosure.Coefficient)

	if b.ActualReadingClosure.DistributorID == "" {
		b.Coefficient = CoefficientType(tariff.Coef)
	}
}

func (b *BillingMeasure) SetExecutionBalanceOrigin(a measures.DailyReadingClosure) {
	var origin measures.ClosureType
	switch {
	case a.EndDate.IsZero():
		origin = measures.NoClosure
	case a.Origin == measures.STG && a.ClosureType == measures.Daily:
		origin = measures.Daily
	case a.Origin == measures.STG && a.ClosureType == measures.Monthly,
		a.Origin == measures.STM:
		origin = measures.Monthly
	default:
		origin = measures.Other
	}

	b.ExecutionSummary.setBalanceOrigin(origin)
}

func (b BillingMeasure) IsBalanceGreaterCch(period measures.PeriodKey, magnitude measures.Magnitude) bool {
	var values *measures.Values
	switch period {
	case measures.P0:
		values = b.AtrVsCurve.P0
	case measures.P1:
		values = b.AtrVsCurve.P1
	case measures.P2:
		values = b.AtrVsCurve.P2
	case measures.P3:
		values = b.AtrVsCurve.P3
	case measures.P4:
		values = b.AtrVsCurve.P4
	case measures.P5:
		values = b.AtrVsCurve.P5
	case measures.P6:
		values = b.AtrVsCurve.P6
	}

	if values == nil {
		return false
	}

	return values.GetMagnitude(magnitude) >= 0
}

func (b BillingMeasure) IsCchComplete() bool {
	for _, curve := range b.BillingLoadCurve {
		if curve.Origin == measures.Filled {
			return false
		}
	}
	return true
}

func (b BillingMeasure) IsChhValid() bool {
	for _, curve := range b.BillingLoadCurve {
		if curve.Origin != measures.Filled {
			return true
		}
	}
	return false
}

func (b BillingMeasure) IsHourly() bool {
	return b.RegisterType == measures.Hourly || b.RegisterType == measures.Both
}

func (b *BillingMeasure) GenerateID() {
	id := fmt.Sprintf("%s%s%s%s%s", b.DistributorCode, b.CUPS, b.InitDate.Format("2006-01-02_15:04"), b.EndDate.Format("2006-01-02_15:04"), b.Version)

	b.Id = fmt.Sprintf("%x", sha256.Sum256([]byte(id)))
}

func NewBillingMeasure(
	cups string,
	initDate time.Time,
	endDate time.Time,
	distributorCode string,
	distributorId string,
	periods []measures.PeriodKey,
	magnitudes []measures.Magnitude,
	meterType measures.MeterType,
) BillingMeasure {
	b := BillingMeasure{
		DistributorID:    distributorId,
		DistributorCode:  distributorCode,
		CUPS:             cups,
		InitDate:         initDate,
		EndDate:          endDate,
		Status:           Calculating,
		Periods:          periods,
		GraphHistory:     map[string]*Graph{},
		Version:          "0",
		Magnitudes:       magnitudes,
		GenerationDate:   time.Now().UTC(),
		ExecutionSummary: ExecutionSummary{},
		MeterType:        meterType,
	}

	b.GenerateID()

	for _, p := range append(b.Periods, measures.P0) {
		switch p {
		case measures.P0:
			b.BillingBalance.P0 = &BillingBalancePeriod{}
			b.AtrVsCurve.P0 = &measures.Values{}
		case measures.P1:
			b.BillingBalance.P1 = &BillingBalancePeriod{}
			b.AtrVsCurve.P1 = &measures.Values{}
		case measures.P2:
			b.BillingBalance.P2 = &BillingBalancePeriod{}
			b.AtrVsCurve.P2 = &measures.Values{}
		case measures.P3:
			b.BillingBalance.P3 = &BillingBalancePeriod{}
			b.AtrVsCurve.P3 = &measures.Values{}
		case measures.P4:
			b.BillingBalance.P4 = &BillingBalancePeriod{}
			b.AtrVsCurve.P4 = &measures.Values{}
		case measures.P5:
			b.BillingBalance.P5 = &BillingBalancePeriod{}
			b.AtrVsCurve.P5 = &measures.Values{}
		case measures.P6:
			b.BillingBalance.P6 = &BillingBalancePeriod{}
			b.AtrVsCurve.P6 = &measures.Values{}
		}
	}

	return b
}

func (b *BillingMeasure) CalculateVersionByPreviousBillingMeasure(previous BillingMeasure) {
	if previous.Id == "" {
		return
	}
	newVersion := previous.Version
	if previous.Status == ReadyToBill || previous.Status == Billed {
		newVersion, _ = previous.IncrementVersion()
	}

	b.setVersion(newVersion)
}

func (b *BillingMeasure) setVersion(newVersion string) {
	b.Version = newVersion
	b.GenerateID()
}
func (b BillingMeasure) IncrementVersion() (string, error) {
	version, err := strconv.Atoi(b.Version)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", version+1), nil
}

func (b BillingMeasure) IsEmptyCentralHoursCch() bool {
	if b.BillingLoadCurve[0].Origin == measures.Filled || b.BillingLoadCurve[len(b.BillingLoadCurve)-1].Origin == measures.Filled {
		return false
	} else {
		return true
	}
}

func (b BillingMeasure) IsThereChhWindows() bool {
	nGapFilled := 0
	for _, curve := range b.BillingLoadCurve {
		if nGapFilled >= 4 {
			return true
		}
		if curve.Origin != measures.Filled {
			nGapFilled = 0
		} else {
			nGapFilled++
		}
	}
	return false
}

func (b BillingMeasure) AreAllCloseAtrEmpty(magnitude measures.Magnitude) bool {

	for _, period := range b.Periods {
		if b.GetBalanceStatus(period, magnitude) == measures.Valid {
			return false
		}
	}

	return true
}

func (b BillingMeasure) IsBalanceValid(magnitude measures.Magnitude) bool {
	return b.GetBalanceStatus(measures.P0, magnitude) == measures.Valid
}

func (b BillingMeasure) IsHouseClose() bool {
	return b.Inaccessible
}

func (b *BillingMeasure) ShouldExecuteGraph() bool {
	hours := b.EndDate.Sub(b.InitDate).Hours()

	if b.RegisterType != measures.NoneType && float64(len(b.BillingLoadCurve)) != hours {
		b.ExecutionSummary.setCurveStatus(measures.Absent)
		b.Supervision()
		b.SetDescriptionStatus(NoMatchLenLoadCurve)
		return false
	}

	return true
}

func (b *BillingMeasure) AfterExecuteGraph() {
	principalMagnitude := b.GetPrincipalMagnitude()
	curveStatus := measures.Complete
	for i, curve := range b.BillingLoadCurve {
		if curve.Origin == measures.Filled {
			curveStatus = measures.Incomplete
		}
		curveMethod, _ := b.GetLoadCurveGeneralEstimatedMethod(i, principalMagnitude)
		b.ExecutionSummary.setCurveType(curveMethod)

	}

	b.ExecutionSummary.setCurveStatus(curveStatus)
	balanceMethod := b.GetBalanceGeneralType(measures.P0, principalMagnitude)
	b.ExecutionSummary.setBalanceType(balanceMethod)
}

func (b *BillingMeasure) BeforeSave(config measures.MeterConfig) {
	for graphKey, g := range b.GraphHistory {
		for k, node := range g.Dict {
			if !node.Done {
				delete(g.Dict, k)
				continue
			}

			if len(node.Algorithms) > 0 {
				b.GraphHistory[graphKey].Algorithms = utils.MapSlice(node.Algorithms, func(item Algorithm) string {
					return item.ID()
				})
			}
		}
	}

	b.fillActualReadingClosure(config)
}

func (b *BillingMeasure) fillActualReadingClosure(config measures.MeterConfig) {
	if !b.ActualReadingClosure.EndDate.IsZero() {
		return
	}

	calendar := measures.NewDailyReadingClosureCalendar(b.GetPeriods())
	for _, period := range b.GetPeriods() {
		for _, magnitude := range b.Magnitudes {
			calendar.SetMagnitudeToPeriod(
				period,
				magnitude,
				b.PreviousReadingClosure.GetCalendarPeriodMagnitude(period, magnitude)+b.GetBalancePeriodMagnitude(period, magnitude),
			)
		}

		calendar.SetFilledToPeriod(period)
	}

	n := measures.DailyReadingClosure{
		ClosureType:       measures.Other,
		DistributorCode:   b.DistributorCode,
		DistributorID:     b.DistributorID,
		CUPS:              b.CUPS,
		InitDate:          b.InitDate,
		EndDate:           b.EndDate,
		MeterSerialNumber: "",
		GenerationDate:    b.GenerationDate,
		ReadingDate:       time.Time{},
		ServiceType:       string(config.ServiceType()),
		PointType:         string(b.PointType),
		Origin:            measures.Filled,
		MeasurePointType:  config.MeasurePointType(),
		ContractNumber:    "",
		CalendarPeriods:   *calendar,
		QualityCode:       "",
		ValidationStatus:  "",
		TariffId:          b.TariffID,
		CalendarCode:      b.CalendarCode,
		Coefficient:       string(b.Coefficient),
		P1Demand:          b.P1Demand,
		P2Demand:          b.P2Demand,
		P3Demand:          b.P3Demand,
		P4Demand:          b.P4Demand,
		P5Demand:          b.P5Demand,
		P6Demand:          b.P6Demand,
		Magnitudes:        b.Magnitudes,
	}

	b.ActualReadingClosure = n
}

func (b *BillingMeasure) Calculated() {
	b.Status = Calculated

}

func (b *BillingMeasure) Supervision() {
	b.Status = Supervision
}

func (b *BillingMeasure) SetDescriptionStatus(status DescriptionStatus) {
	b.DescriptionStatus = &status
}
