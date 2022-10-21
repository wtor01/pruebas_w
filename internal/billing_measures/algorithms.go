package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

var (
	ErrorPeriodKeyDoesNotMatch  = errors.New("your period key name is NOT VALID")
	ErrorPowerDemanded          = errors.New("your period power demand is NOT VALID")
	ErrorNotFoundConsumProfiles = errors.New("CONSUM PROFILES RECORDS NOT FOUND")
)

type Algorithm interface {
	ID() string
	Execute(_ context.Context) error
}

type Log struct {
	name string
}

func (l Log) Execute(_ context.Context) error {
	fmt.Print(l.name)

	return nil
}

func (l Log) ID() string {
	return fmt.Sprintf("LOG - %s", l.name)
}

type CchCompleted struct {
	Name      string
	b         *BillingMeasure `bson:"-"`
	Period    measures.PeriodKey
	Magnitude measures.Magnitude
}

func NewCchCompleted(b *BillingMeasure, period measures.PeriodKey, magnitude measures.Magnitude) *CchCompleted {
	return &CchCompleted{b: b, Period: period, Name: "CCH_COMPLETED", Magnitude: magnitude}
}

func (algorithm CchCompleted) ID() string {
	return algorithm.Name
}

func (algorithm CchCompleted) Execute(_ context.Context) error {
	for i, curve := range algorithm.b.BillingLoadCurve {
		if curve.Period != algorithm.Period {
			continue
		}
		algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 1)
		algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, RealBalance)
		algorithm.b.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralReal)
		algorithm.b.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)
	}
	return nil
}

type BalanceCalculatedByCch struct {
	Name      string
	b         *BillingMeasure `bson:"-"`
	Period    measures.PeriodKey
	Magnitude measures.Magnitude
}

func NewBalanceCalculatedByCch(b *BillingMeasure, period measures.PeriodKey, magnitude measures.Magnitude) *BalanceCalculatedByCch {
	return &BalanceCalculatedByCch{b: b, Period: period, Name: "BALANCE_CALCULATED_BY_CCH", Magnitude: magnitude}
}

func (algorithm BalanceCalculatedByCch) ID() string {
	return algorithm.Name
}

func (algorithm BalanceCalculatedByCch) Execute(_ context.Context) error {
	billingBalancePeriod := algorithm.b.GetBalancePeriod(algorithm.Period)

	if billingBalancePeriod == nil {
		return ErrorPeriodKeyDoesNotMatch
	}

	for _, curve := range algorithm.b.BillingLoadCurve {
		if curve.Period != algorithm.Period {
			continue
		}
		algorithm.b.SumBalancePeriodMagnitude(algorithm.Period, algorithm.Magnitude, curve.GetMagnitude(algorithm.Magnitude))
	}
	algorithm.b.SetBalanceEstimateCode(algorithm.Period, algorithm.Magnitude, 2)
	algorithm.b.SetBalanceMeasureType(algorithm.Period, algorithm.Magnitude, FirmBalanceMeasure)
	algorithm.b.SetBalanceType(algorithm.Period, algorithm.Magnitude, CalculatedBalance)
	algorithm.b.SetBalanceGeneralType(algorithm.Period, algorithm.Magnitude, GeneralCalculated)
	algorithm.b.SetBalanceOrigin(algorithm.Period, algorithm.Magnitude, TlgOrigin)
	algorithm.b.SetWorstBalanceGeneralType(algorithm.Magnitude, GeneralCalculated)
	return nil
}

type EstimatedBalanceByPowerDemand struct {
	Name      string
	b         *BillingMeasure `bson:"-"`
	Period    measures.PeriodKey
	Magnitude measures.Magnitude
}

func NewEstimatedBalanceByPowerDemand(b *BillingMeasure, period measures.PeriodKey, magnitude measures.Magnitude) *EstimatedBalanceByPowerDemand {
	return &EstimatedBalanceByPowerDemand{
		Name:      "BALANCE_ESTIMATE_BY_POWER_DEMAND",
		b:         b,
		Period:    period,
		Magnitude: magnitude,
	}
}

func (algorithm EstimatedBalanceByPowerDemand) ID() string {
	return algorithm.Name
}

func (algorithm EstimatedBalanceByPowerDemand) Execute(_ context.Context) error {

	balancePeriod := algorithm.b.GetBalancePeriod(algorithm.Period)
	if balancePeriod == nil {
		return ErrorPeriodKeyDoesNotMatch
	}

	powerDemanded := algorithm.b.GetPeriodPowerDemanded(algorithm.Period)

	nTime := 0.0

	for _, curve := range algorithm.b.BillingLoadCurve {
		if curve.Period == algorithm.Period {
			nTime++
		}
	}

	utilizationFactor := .33

	if powerDemanded <= 10 {
		nTime = nTime / 24
		switch algorithm.Period {
		case measures.P1:
			utilizationFactor = .5
		case measures.P2:
			utilizationFactor = .5
		case measures.P3:
			utilizationFactor = 2.7
		}
	}

	algorithm.b.SetBalanceEstimateCode(algorithm.Period, algorithm.Magnitude, 6)
	algorithm.b.SetBalanceType(algorithm.Period, algorithm.Magnitude, EstimatedContractPower)
	algorithm.b.SetBalanceMeasureType(algorithm.Period, algorithm.Magnitude, ProvisionalBalanceMeasure)
	algorithm.b.SetBalanceOrigin(algorithm.Period, algorithm.Magnitude, EstimateOrigin)
	algorithm.b.SetBalanceGeneralType(algorithm.Period, algorithm.Magnitude, GeneralEstimated)
	algorithm.b.SetWorstBalanceGeneralType(algorithm.Magnitude, GeneralEstimated)
	value := 0.0

	if algorithm.Magnitude == measures.AI {
		value = (powerDemanded * utilizationFactor * nTime) * 1000
	}
	algorithm.b.SetBalancePeriodMagnitude(algorithm.Period, algorithm.Magnitude, value)

	return nil
}

type CchPartialAdjustment struct {
	Name      string
	b         *BillingMeasure `bson:"-"`
	Period    measures.PeriodKey
	Magnitude measures.Magnitude
}

func NewCchPartialAdjustment(b *BillingMeasure, period measures.PeriodKey, magnitude measures.Magnitude) *CchPartialAdjustment {
	return &CchPartialAdjustment{
		Name:      "CCH_PARTIAL_ADJUST",
		b:         b,
		Period:    period,
		Magnitude: magnitude,
	}
}

func (algorithm CchPartialAdjustment) ID() string {
	return algorithm.Name
}

func (algorithm CchPartialAdjustment) Execute(_ context.Context) error {
	var cchTotals float64

	for _, curve := range algorithm.b.BillingLoadCurve {
		if algorithm.Period != curve.Period {
			continue
		}
		cchTotals += curve.GetMagnitude(algorithm.Magnitude)
	}

	for i, curve := range algorithm.b.BillingLoadCurve {
		if algorithm.Period != curve.Period {
			continue
		}
		algorithm.b.SetLoadCurvePeriodMagnitude(
			i,
			algorithm.Magnitude,
			getFormulaAdjustment(
				curve.GetMagnitude(algorithm.Magnitude),
				algorithm.b.GetBalancePeriodMagnitude(algorithm.Period, algorithm.Magnitude),
				cchTotals,
			),
		)
		algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, Adjustment)
		algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 3)
		algorithm.b.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralAdjusted)

		if algorithm.b.GetBalanceEstimateCode(algorithm.Period, algorithm.Magnitude) == 4 {
			algorithm.b.SetLoadCurveMeasureType(i, algorithm.Magnitude, ProvisionalBalanceMeasure)
		} else {
			algorithm.b.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)
		}
	}

	return nil
}

type CchCompleteAdjustment struct {
	Name      string
	b         *BillingMeasure `bson:"-"`
	Period    measures.PeriodKey
	Magnitude measures.Magnitude
}

func NewCchCompleteAdjustment(b *BillingMeasure, period measures.PeriodKey, magnitude measures.Magnitude) *CchCompleteAdjustment {
	return &CchCompleteAdjustment{
		Name:      "CCH_COMPLETE_ADJUST",
		b:         b,
		Period:    period,
		Magnitude: magnitude,
	}
}

func (algorithm CchCompleteAdjustment) ID() string {
	return algorithm.Name
}

func getFormulaAdjustment(curve, balance, cchTotals float64) float64 {
	if cchTotals == 0 {
		return 0
	}
	return math.Round(curve * (balance / cchTotals))
}

func (algorithm CchCompleteAdjustment) Execute(_ context.Context) error {
	var cchTotals float64
	balance := algorithm.b.GetBalancePeriod(algorithm.Period)

	if balance == nil {
		return ErrorPeriodKeyDoesNotMatch
	}

	for _, curve := range algorithm.b.BillingLoadCurve {
		if algorithm.Period != curve.Period {
			continue
		}
		cchTotals += curve.GetMagnitude(algorithm.Magnitude)

	}

	for i, curve := range algorithm.b.BillingLoadCurve {
		if algorithm.Period != curve.Period {
			continue
		}
		algorithm.b.SetLoadCurvePeriodMagnitude(
			i,
			algorithm.Magnitude,
			getFormulaAdjustment(
				curve.GetMagnitude(algorithm.Magnitude),
				algorithm.b.GetBalancePeriodMagnitude(algorithm.Period, algorithm.Magnitude),
				cchTotals,
			),
		)
		algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, Adjustment)
		algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 3)
		algorithm.b.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralAdjusted)

		if algorithm.b.GetBalanceEstimateCode(algorithm.Period, algorithm.Magnitude) == 4 {
			algorithm.b.SetLoadCurveMeasureType(i, algorithm.Magnitude, ProvisionalBalanceMeasure)
		} else {
			algorithm.b.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)
		}
	}

	return nil
}

type BalanceCompleted struct {
	Name      string
	b         *BillingMeasure `bson:"-"`
	Period    measures.PeriodKey
	Magnitude measures.Magnitude
}

func NewBalanceCompleted(b *BillingMeasure, period measures.PeriodKey, magnitude measures.Magnitude) *BalanceCompleted {
	return &BalanceCompleted{b: b, Period: period, Name: "BALANCE_COMPLETED", Magnitude: magnitude}
}

func (algorithm BalanceCompleted) ID() string {
	return algorithm.Name
}

func (algorithm BalanceCompleted) Execute(_ context.Context) error {

	if algorithm.b.ActualReadingClosure.Origin == measures.STG || algorithm.b.ActualReadingClosure.Origin == measures.STM {
		if algorithm.b.IsAtrVsCurveValid(algorithm.Period, algorithm.Magnitude) {
			algorithm.b.SetBalanceMeasureType(algorithm.Period, algorithm.Magnitude, FirmBalanceMeasure)
		} else {
			algorithm.b.SetBalanceMeasureType(algorithm.Period, algorithm.Magnitude, ProvisionalBalanceMeasure)
		}
		algorithm.b.SetBalanceType(algorithm.Period, algorithm.Magnitude, RealBalance)
		algorithm.b.SetBalanceOrigin(algorithm.Period, algorithm.Magnitude, RemoteBalanceOrigin)
		algorithm.b.SetBalanceEstimateCode(algorithm.Period, algorithm.Magnitude, 1)
		algorithm.b.SetBalanceGeneralType(algorithm.Period, algorithm.Magnitude, GeneralReal)
		algorithm.b.SetWorstBalanceGeneralType(algorithm.Magnitude, GeneralReal)
	}
	if algorithm.b.ActualReadingClosure.Origin == measures.TPL || algorithm.b.ActualReadingClosure.Origin == measures.Manual || algorithm.b.ActualReadingClosure.Origin == measures.Visual {
		algorithm.b.SetBalanceType(algorithm.Period, algorithm.Magnitude, RealBalance)
		algorithm.b.SetBalanceGeneralType(algorithm.Period, algorithm.Magnitude, GeneralReal)
		algorithm.b.SetWorstBalanceGeneralType(algorithm.Magnitude, GeneralReal)
		algorithm.b.SetBalanceMeasureType(algorithm.Period, algorithm.Magnitude, FirmBalanceMeasure)
		algorithm.b.SetBalanceOrigin(algorithm.Period, algorithm.Magnitude, RemoteBalanceOrigin)
		algorithm.b.SetBalanceEstimateCode(algorithm.Period, algorithm.Magnitude, 3)
	}
	if algorithm.b.ActualReadingClosure.Origin == measures.Auto {
		algorithm.b.SetBalanceType(algorithm.Period, algorithm.Magnitude, EstimateBalance)
		algorithm.b.SetBalanceGeneralType(algorithm.Period, algorithm.Magnitude, GeneralEstimated)
		algorithm.b.SetBalanceMeasureType(algorithm.Period, algorithm.Magnitude, ProvisionalBalanceMeasure)
		algorithm.b.SetBalanceOrigin(algorithm.Period, algorithm.Magnitude, AutoBalanceOrigin)
		algorithm.b.SetBalanceEstimateCode(algorithm.Period, algorithm.Magnitude, 4)
		algorithm.b.SetWorstBalanceGeneralType(algorithm.Magnitude, GeneralEstimated)
	}

	return nil
}

type EstimateBase struct {
	B         *BillingMeasure `bson:"-"`
	Period    measures.PeriodKey
	Magnitude measures.Magnitude
}

func (algorithm EstimateBase) getSumProfilesEstimate(keysProfiles map[time.Time]ConsumProfile) (sumCCH float64, sumProfiles float64) {

	for _, curve := range algorithm.B.BillingLoadCurve {
		if curve.Period != algorithm.Period {
			continue
		}
		if curve.Origin != measures.Filled {
			sumCCH += curve.GetMagnitude(algorithm.Magnitude)
			continue
		}
		if cf, ok := keysProfiles[curve.EndDate]; ok {
			sumProfiles += cf.GetValueByCoefficientType(algorithm.B.Coefficient)
		}

	}
	return sumCCH, sumProfiles

}

func (algorithm EstimateBase) getValueEstimate(valueCCH float64, profileValue float64, profileSum float64) float64 {
	if profileSum == 0.0 {
		return 0
	}

	return math.Round((algorithm.B.GetBalancePeriodMagnitude(algorithm.Period, algorithm.Magnitude) - valueCCH) * (profileValue / profileSum))
}

type CchPartialEstimation struct {
	EstimateBase
	Name     string
	Profiles ConsumProfileRepository `bson:"-"`
}

func NewCchPartialEstimation(b *BillingMeasure, period measures.PeriodKey, profiles ConsumProfileRepository, magnitude measures.Magnitude) *CchPartialEstimation {
	return &CchPartialEstimation{
		Name:     "CCH_PARTIAL_ESTIMATE",
		Profiles: profiles,
		EstimateBase: EstimateBase{
			B:         b,
			Period:    period,
			Magnitude: magnitude,
		},
	}
}

func getKeysValuesProfiles(cps []ConsumProfile) map[time.Time]ConsumProfile {
	keys := make(map[time.Time]ConsumProfile)

	for _, cp := range cps {
		keys[cp.Date] = cp
	}
	return keys
}

func (algorithm CchPartialEstimation) setMetaDataOrigin(index int) {
	estimatedCode := algorithm.B.GetBalanceEstimateCode(algorithm.Period, algorithm.Magnitude)
	algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, ProvisionalBalanceMeasure)
	switch estimatedCode {
	case 1, 3:
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, ProfileMeasure)
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 2)
	case 4:
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, ProfileMeasureAutoReading)
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 4)
	case 5:
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, EstimateHistoryProfile)
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 5)
	case 6:
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, EstimateFactorUsed)
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 6)
	}

	if estimatedCode <= 4 {
		algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralOutlined)
		return
	}

	algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralEstimated)
}

func (algorithm CchPartialEstimation) ID() string {
	return algorithm.Name
}

func (algorithm CchPartialEstimation) Execute(ctx context.Context) error {

	cProfiles, err := algorithm.Profiles.Search(ctx, QueryConsumProfile{StartDate: algorithm.B.InitDate, EndDate: algorithm.B.EndDate})
	if err != nil {
		return err
	}
	if len(cProfiles) == 0 {
		return ErrorNotFoundConsumProfiles

	}
	keysProfiles := getKeysValuesProfiles(cProfiles)
	if len(keysProfiles) == 0 {
		return ErrorNotFoundConsumProfiles
	}
	sumCCH, sumProfiles := algorithm.getSumProfilesEstimate(keysProfiles)

	if algorithm.B.GetBalancePeriod(algorithm.Period) == nil {
		return ErrorPeriodKeyDoesNotMatch
	}

	for i, curve := range algorithm.B.BillingLoadCurve {
		if curve.Period != algorithm.Period {
			continue
		}

		profileValue := 0.00

		if cf, ok := keysProfiles[curve.EndDate]; ok {
			profileValue = cf.GetValueByCoefficientType(algorithm.B.Coefficient)
		}

		if curve.Origin == measures.Filled {
			algorithm.B.SetLoadCurvePeriodMagnitude(i, algorithm.Magnitude, algorithm.getValueEstimate(sumCCH, profileValue, sumProfiles))
			algorithm.setMetaDataOrigin(i)
		} else {
			algorithm.B.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 1)
			algorithm.B.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, RealBalance)
			algorithm.B.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralReal)
			algorithm.B.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)
		}
	}

	return nil
}

type CchTotalEstimation struct {
	EstimateBase
	Name     string
	Profiles ConsumProfileRepository `bson:"-"`
}

func NewCchTotalEstimation(b *BillingMeasure, period measures.PeriodKey, profiles ConsumProfileRepository, magnitude measures.Magnitude) *CchTotalEstimation {
	return &CchTotalEstimation{
		Name:     "CCH_TOTAL_ESTIMATE",
		Profiles: profiles,
		EstimateBase: EstimateBase{
			B:         b,
			Period:    period,
			Magnitude: magnitude,
		},
	}
}

func (algorithm CchTotalEstimation) setMetaDataOrigin(index int) {
	estimatedCode := algorithm.B.GetBalanceEstimateCode(algorithm.Period, algorithm.Magnitude)
	algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, ProvisionalBalanceMeasure)

	switch estimatedCode {
	case 1, 3:
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, ProfileMeasure)
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 2)
	case 4:
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, ProfileMeasureAutoReading)
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 4)
	case 5:
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, EstimateHistoryProfile)
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 5)
	case 6:
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, EstimateFactorUsed)
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 6)
	}

	if estimatedCode <= 4 {
		algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralOutlined)
		return
	}

	algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralEstimated)
}

func (algorithm CchTotalEstimation) ID() string {
	return algorithm.Name
}

func (algorithm CchTotalEstimation) Execute(ctx context.Context) error {

	cProfiles, err := algorithm.Profiles.Search(ctx, QueryConsumProfile{StartDate: algorithm.B.InitDate, EndDate: algorithm.B.EndDate})
	if err != nil {
		return err
	}
	if len(cProfiles) == 0 {
		return ErrorNotFoundConsumProfiles

	}
	keysProfiles := getKeysValuesProfiles(cProfiles)
	if len(keysProfiles) == 0 {
		return ErrorNotFoundConsumProfiles
	}
	sumCCH, sumProfiles := algorithm.getSumProfilesEstimate(keysProfiles)

	if algorithm.B.GetBalancePeriod(algorithm.Period) == nil {
		return ErrorPeriodKeyDoesNotMatch
	}

	for i, curve := range algorithm.B.BillingLoadCurve {
		if curve.Period != algorithm.Period {
			continue
		}
		profileValue := 0.00
		if cf, ok := keysProfiles[curve.EndDate]; ok {
			profileValue = cf.GetValueByCoefficientType(algorithm.B.Coefficient)
		}
		algorithm.B.SetLoadCurvePeriodMagnitude(i, algorithm.Magnitude, algorithm.getValueEstimate(sumCCH, profileValue, sumProfiles))
		algorithm.setMetaDataOrigin(i)

	}
	return nil
}

type EstimatedHistoryTlg struct {
	Name       string
	b          *BillingMeasure `bson:"-"`
	Period     measures.PeriodKey
	Magnitude  measures.Magnitude
	ContextTlg *GraphContext `bson:"-"`
}

func NewEstimatedHistoryTlg(b *BillingMeasure, period measures.PeriodKey, magnitude measures.Magnitude, contextTlg *GraphContext) *EstimatedHistoryTlg {
	return &EstimatedHistoryTlg{
		Name:       "BALANCE_ESTIMATED_HISTORY_TLG",
		b:          b,
		Period:     period,
		Magnitude:  magnitude,
		ContextTlg: contextTlg,
	}
}

func (algorithm EstimatedHistoryTlg) ID() string {
	return algorithm.Name
}

func (algorithm EstimatedHistoryTlg) Execute(_ context.Context) error {
	var lastPeriodHours float64
	var actualPeriodHours float64
	billingHistory := algorithm.ContextTlg.LastHistory
	actualBillingPeriod := algorithm.b.GetBalancePeriod(algorithm.Period)
	lastBillingPeriod := billingHistory.GetBalancePeriod(algorithm.Period)

	if actualBillingPeriod == nil || lastBillingPeriod == nil {
		return errors.New("period cannot be nil")
	}

	for _, curve := range algorithm.b.BillingLoadCurve {
		if curve.Period != algorithm.Period {
			continue
		}
		if algorithm.b.RegisterType == measures.Hourly {
			actualPeriodHours += 1
		} else {
			actualPeriodHours += .25
		}
	}

	for _, curve := range billingHistory.BillingLoadCurve {
		if curve.Period != algorithm.Period {
			continue
		}
		if billingHistory.RegisterType == measures.Hourly {
			lastPeriodHours += 1
		} else {
			lastPeriodHours += .25
		}
	}

	if lastPeriodHours == 0 {
		lastPeriodHours = 1
	}

	algorithm.b.SetBalancePeriodMagnitude(
		algorithm.Period,
		algorithm.Magnitude,
		math.Round((billingHistory.GetBalancePeriodMagnitude(algorithm.Period, algorithm.Magnitude)/lastPeriodHours)*actualPeriodHours),
	)
	algorithm.b.SetBalanceEstimateCode(algorithm.Period, algorithm.Magnitude, 5)
	algorithm.b.SetBalanceType(algorithm.Period, algorithm.Magnitude, EstimateHistoryProfile)
	algorithm.b.SetBalanceMeasureType(algorithm.Period, algorithm.Magnitude, ProvisionalBalanceMeasure)
	algorithm.b.SetBalanceOrigin(algorithm.Period, algorithm.Magnitude, EstimateOrigin)
	algorithm.b.SetBalanceGeneralType(algorithm.Period, algorithm.Magnitude, GeneralEstimated)
	algorithm.b.SetWorstBalanceGeneralType(algorithm.Magnitude, GeneralEstimated)
	return nil
}
