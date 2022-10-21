package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"
)

type FlatCastNoTLG struct {
	Name      string
	b         *BillingMeasure `bson:"-"`
	Magnitude measures.Magnitude
}

func NewFlatCastNoTLG(b *BillingMeasure, magnitude measures.Magnitude) *FlatCastNoTLG {
	return &FlatCastNoTLG{
		Name:      "CCH_FLAT_CAST_NO_TLG",
		b:         b,
		Magnitude: magnitude,
	}
}

func (algorithm FlatCastNoTLG) ID() string {
	return algorithm.Name
}

type PeriodGapsAndSumFlatCast struct {
	Value float64
	nGap  int
}

func (algorithm FlatCastNoTLG) Execute(_ context.Context) error {
	var periodsLocated []measures.PeriodKey

	sumUpPeriods := make(map[measures.PeriodKey]*PeriodGapsAndSumFlatCast)
	for _, period := range algorithm.b.Periods {
		periodsLocated = append(periodsLocated, period)
		sumUpPeriods[period] = &PeriodGapsAndSumFlatCast{}
	}

	for _, curve := range algorithm.b.BillingLoadCurve {
		if _, ok := sumUpPeriods[curve.Period]; !ok {
			continue
		}
		sumUpPeriods[curve.Period].Value += curve.GetMagnitude(algorithm.Magnitude)
		if curve.Origin == measures.Filled {
			sumUpPeriods[curve.Period].nGap += 1
		}
	}
	for i, curve := range algorithm.b.BillingLoadCurve {
		if _, ok := sumUpPeriods[curve.Period]; !ok {
			continue
		}
		if curve.Origin == measures.Filled {
			algorithm.b.SetLoadCurvePeriodMagnitude(i, algorithm.Magnitude, algorithm.getValueFromFormula(curve.Period, *sumUpPeriods[curve.Period]))

		}
		algorithm.setCurveFeatures(i)
	}

	return nil
}

func (algorithm FlatCastNoTLG) getValueFromFormula(period measures.PeriodKey, periodGapsAndSumFlatCast PeriodGapsAndSumFlatCast) float64 {
	if periodGapsAndSumFlatCast.nGap == 0 {
		return 0
	}
	return (algorithm.b.GetBalancePeriodMagnitude(period, algorithm.Magnitude) - periodGapsAndSumFlatCast.Value) / float64(periodGapsAndSumFlatCast.nGap)
}

func (algorithm FlatCastNoTLG) setCurveFeatures(indexLoadCurve int) {
	if algorithm.b.isPointType([]measures.PointType{"1", "2"}) {
		if algorithm.b.GetLoadCurveOrigin(indexLoadCurve) == measures.Filled {
			algorithm.b.SetLoadCurveEstimatedCode(indexLoadCurve, algorithm.Magnitude, 8)
			algorithm.b.SetLoadCurveEstimatedMethod(indexLoadCurve, algorithm.Magnitude, EstimatedByFlatProfile)
			algorithm.b.SetLoadCurveMeasureType(indexLoadCurve, algorithm.Magnitude, ProvisionalBalanceMeasure)
			algorithm.b.SetLoadCurveGeneralEstimatedMethod(indexLoadCurve, algorithm.Magnitude, GeneralEstimated)
			return
		}

		switch algorithm.b.GetLoadCurveEquipment(indexLoadCurve) {
		case measures.Main:
			algorithm.b.SetLoadCurveEstimatedCode(indexLoadCurve, algorithm.Magnitude, 1)
			algorithm.b.SetLoadCurveEstimatedMethod(indexLoadCurve, algorithm.Magnitude, FirmMainConfig)
		case measures.Redundant:
			algorithm.b.SetLoadCurveEstimatedCode(indexLoadCurve, algorithm.Magnitude, 2)
			algorithm.b.SetLoadCurveEstimatedMethod(indexLoadCurve, algorithm.Magnitude, FirmRedundantConfig)
		case measures.Receipt:
			algorithm.b.SetLoadCurveEstimatedCode(indexLoadCurve, algorithm.Magnitude, 3)
			algorithm.b.SetLoadCurveEstimatedMethod(indexLoadCurve, algorithm.Magnitude, FirmReceiptConfig)
		}
		algorithm.b.SetLoadCurveGeneralEstimatedMethod(indexLoadCurve, algorithm.Magnitude, GeneralReal)
		algorithm.b.SetLoadCurveMeasureType(indexLoadCurve, algorithm.Magnitude, FirmBalanceMeasure)
		return

	}

	if algorithm.b.isPointType([]measures.PointType{"3", "4"}) {
		if algorithm.b.GetLoadCurveOrigin(indexLoadCurve) == measures.Filled {
			algorithm.b.SetLoadCurveEstimatedCode(indexLoadCurve, algorithm.Magnitude, 1)
			algorithm.b.SetLoadCurveEstimatedMethod(indexLoadCurve, algorithm.Magnitude, RealValidMeasure)
			algorithm.b.SetLoadCurveGeneralEstimatedMethod(indexLoadCurve, algorithm.Magnitude, GeneralReal)
			algorithm.b.SetLoadCurveMeasureType(indexLoadCurve, algorithm.Magnitude, FirmBalanceMeasure)
			return
		}

		algorithm.b.SetLoadCurveEstimatedCode(indexLoadCurve, algorithm.Magnitude, 3)
		algorithm.b.SetLoadCurveEstimatedMethod(indexLoadCurve, algorithm.Magnitude, RealMeasureAdjustment)
		algorithm.b.SetLoadCurveGeneralEstimatedMethod(indexLoadCurve, algorithm.Magnitude, GeneralReal)
		algorithm.b.SetLoadCurveMeasureType(indexLoadCurve, algorithm.Magnitude, ProvisionalBalanceMeasure)
		return
	}
}

type CchAverage struct {
	Name         string
	B            *BillingMeasure `bson:"-"`
	ContextNoTlg *GraphContext   `bson:"-"`
	Magnitude    measures.Magnitude
}

func NewCchAverage(b *BillingMeasure, magnitude measures.Magnitude, contextNoTlg *GraphContext) *CchAverage {
	return &CchAverage{
		Name:         "CCH_AVERAGE",
		B:            b,
		ContextNoTlg: contextNoTlg,
		Magnitude:    magnitude,
	}
}

func (algorithm CchAverage) ID() string {
	return algorithm.Name
}

func (algorithm CchAverage) transformCurve(loadCurve process_measures.ProcessedLoadCurve) BillingLoadCurve {
	return BillingLoadCurve{
		EndDate: loadCurve.EndDate,
		Origin:  loadCurve.Origin,
		AI:      loadCurve.AI,
		AE:      loadCurve.AE,
		R1:      loadCurve.R1,
		R2:      loadCurve.R2,
		R3:      loadCurve.R3,
		R4:      loadCurve.R4,
		Period:  loadCurve.Period,
	}
}

func (algorithm CchAverage) getAverageCurve(previous BillingLoadCurve, next BillingLoadCurve) float64 {
	return (previous.GetMagnitude(algorithm.Magnitude) + next.GetMagnitude(algorithm.Magnitude)) / 2
}

func (algorithm CchAverage) setNewValues(previousCurve BillingLoadCurve, nextCurve BillingLoadCurve, idxToSet []int) {
	for _, idx := range idxToSet {
		algorithm.B.SetLoadCurvePeriodMagnitude(idx, algorithm.Magnitude, algorithm.getAverageCurve(previousCurve, nextCurve))
	}
}
func (algorithm CchAverage) setMetadata(idxs ...int) {

	for _, idx := range idxs {
		switch algorithm.B.PointType {
		case "1", "2":
			{
				if algorithm.B.GetLoadCurveOrigin(idx) == measures.Filled {
					algorithm.B.SetLoadCurveEstimatedCode(idx, algorithm.Magnitude, 9)
					algorithm.B.SetLoadCurveEstimatedMethod(idx, algorithm.Magnitude, EstimateOnlyHistoric)
					algorithm.B.SetLoadCurveGeneralEstimatedMethod(idx, algorithm.Magnitude, GeneralEstimated)
					algorithm.B.SetLoadCurveMeasureType(idx, algorithm.Magnitude, ProvisionalBalanceMeasure)
				} else {
					algorithm.B.SetLoadCurveMeasureType(idx, algorithm.Magnitude, FirmBalanceMeasure)
					algorithm.B.SetLoadCurveGeneralEstimatedMethod(idx, algorithm.Magnitude, GeneralReal)
					switch algorithm.B.BillingLoadCurve[idx].MeasurePointType {
					case measures.MeasurePointTypeP:
						{
							algorithm.B.SetLoadCurveEstimatedMethod(idx, algorithm.Magnitude, FirmMainConfig)
							algorithm.B.SetLoadCurveEstimatedCode(idx, algorithm.Magnitude, 1)
						}
					case measures.MeasurePointTypeR:
						{
							algorithm.B.SetLoadCurveEstimatedMethod(idx, algorithm.Magnitude, FirmRedundantConfig)
							algorithm.B.SetLoadCurveEstimatedCode(idx, algorithm.Magnitude, 2)
						}
					case measures.MeasurePointTypeC:
						{
							algorithm.B.SetLoadCurveEstimatedMethod(idx, algorithm.Magnitude, FirmReceiptConfig)
							algorithm.B.SetLoadCurveEstimatedCode(idx, algorithm.Magnitude, 3)
						}

					}
				}
			}
		case "3", "4", "5":
			{
				if algorithm.B.GetLoadCurveOrigin(idx) == measures.Filled {
					algorithm.B.SetLoadCurveEstimatedCode(idx, algorithm.Magnitude, 5)
					algorithm.B.SetLoadCurveEstimatedMethod(idx, algorithm.Magnitude, EstimateOnlyHistoric)
					algorithm.B.SetLoadCurveGeneralEstimatedMethod(idx, algorithm.Magnitude, GeneralEstimated)
					algorithm.B.SetLoadCurveMeasureType(idx, algorithm.Magnitude, ProvisionalBalanceMeasure)
				} else {
					algorithm.B.SetLoadCurveEstimatedCode(idx, algorithm.Magnitude, 1)
					algorithm.B.SetLoadCurveEstimatedMethod(idx, algorithm.Magnitude, RealValidMeasure)
					algorithm.B.SetLoadCurveGeneralEstimatedMethod(idx, algorithm.Magnitude, GeneralReal)
					algorithm.B.SetLoadCurveMeasureType(idx, algorithm.Magnitude, FirmBalanceMeasure)

				}
			}
		}
	}
}

func (algorithm CchAverage) Execute(_ context.Context) error {
	periods := algorithm.B.Periods

	idxSliceMeasuresToEstimate := map[measures.PeriodKey][]int{}

	mapPreviousRealValueNew := map[measures.PeriodKey]*BillingLoadCurve{}
	mapNextRealValueNew := map[measures.PeriodKey]*BillingLoadCurve{}
	mapFutureCurve := map[measures.PeriodKey]*BillingLoadCurve{}

	for _, p := range periods {
		idxSliceMeasuresToEstimate[p] = make([]int, 0)
		mapPreviousRealValueNew[p] = nil
		mapNextRealValueNew[p] = nil
		mapFutureCurve[p] = nil
	}

	for i := len(algorithm.ContextNoTlg.SimpleHistoric.PreviousLoadCurve) - 1; i >= 0; i-- {
		if val, _ := mapPreviousRealValueNew[algorithm.ContextNoTlg.SimpleHistoric.PreviousLoadCurve[i].Period]; val == nil {
			m := algorithm.transformCurve(algorithm.ContextNoTlg.SimpleHistoric.PreviousLoadCurve[i])
			mapPreviousRealValueNew[algorithm.ContextNoTlg.SimpleHistoric.PreviousLoadCurve[i].Period] = &m
		}
	}

	for _, m := range algorithm.ContextNoTlg.SimpleHistoric.NextLoadCurve {
		if val, _ := mapFutureCurve[m.Period]; val == nil {
			mm := algorithm.transformCurve(m)
			mapFutureCurve[m.Period] = &mm
		}
	}

	for i := range algorithm.B.BillingLoadCurve {
		algorithm.setMetadata(i)
		curve := algorithm.B.BillingLoadCurve[i]
		if curve.Origin == measures.Filled {
			idxSliceMeasuresToEstimate[curve.Period] = append(idxSliceMeasuresToEstimate[curve.Period], i)
			continue
		}

		if len(idxSliceMeasuresToEstimate[curve.Period]) == 0 {
			mapPreviousRealValueNew[curve.Period] = &curve
		} else {
			mapNextRealValueNew[curve.Period] = &curve
		}

		if len(idxSliceMeasuresToEstimate[curve.Period]) != 0 && mapPreviousRealValueNew[curve.Period] != nil && mapNextRealValueNew[curve.Period] != nil {
			algorithm.setNewValues(*mapPreviousRealValueNew[curve.Period], *mapNextRealValueNew[curve.Period], idxSliceMeasuresToEstimate[curve.Period])
			idxSliceMeasuresToEstimate[curve.Period] = make([]int, 0)
			mapPreviousRealValueNew[curve.Period] = mapNextRealValueNew[curve.Period]
			mapNextRealValueNew[curve.Period] = nil
		}

	}

	for period, ids := range idxSliceMeasuresToEstimate {
		if len(ids) != 0 && mapPreviousRealValueNew[period] != nil && mapFutureCurve[period] != nil {
			algorithm.setNewValues(*mapPreviousRealValueNew[period], *mapFutureCurve[period], ids)
			idxSliceMeasuresToEstimate[period] = make([]int, 0)
			mapPreviousRealValueNew[period] = nil
		}
	}
	return nil
}

type BalanceCompleteNoTlg struct {
	Name      string
	b         *BillingMeasure
	Magnitude measures.Magnitude
}

func NewBalanceCompleteNoTlg(b *BillingMeasure, magnitude measures.Magnitude) *BalanceCompleteNoTlg {
	return &BalanceCompleteNoTlg{
		Name:      "BALANCE_COMPLETE",
		b:         b,
		Magnitude: magnitude,
	}
}

func (algorithm BalanceCompleteNoTlg) ID() string {
	return algorithm.Name
}

func (algorithm BalanceCompleteNoTlg) Execute(_ context.Context) error {

	for _, period := range algorithm.b.GetPeriods() {
		algorithm.b.SetBalanceGeneralType(period, algorithm.Magnitude, GeneralReal)
		switch algorithm.b.BillingBalance.Origin {
		case measures.STM:
			algorithm.b.SetBalanceEstimateCode(period, algorithm.Magnitude, 1)
			algorithm.b.SetBalanceOrigin(period, algorithm.Magnitude, TlmOrigin)
			algorithm.b.SetBalanceType(period, algorithm.Magnitude, RealByRemoteRead)
			algorithm.b.SetBalanceMeasureType(period, algorithm.Magnitude, FirmBalanceMeasure)

		case measures.TPL, measures.File, measures.Manual, measures.Visual:
			algorithm.b.SetBalanceEstimateCode(period, algorithm.Magnitude, 3)
			algorithm.b.SetBalanceOrigin(period, algorithm.Magnitude, LocalOrigin)
			algorithm.b.SetBalanceType(period, algorithm.Magnitude, RealByAbsLocalRead)
			algorithm.b.SetBalanceMeasureType(period, algorithm.Magnitude, FirmBalanceMeasure)

		case measures.Auto:
			algorithm.b.SetBalanceEstimateCode(period, algorithm.Magnitude, 4)
			algorithm.b.SetBalanceOrigin(period, algorithm.Magnitude, LocalOrigin)
			algorithm.b.SetBalanceType(period, algorithm.Magnitude, RealByAutoRead)
			algorithm.b.SetBalanceMeasureType(period, algorithm.Magnitude, ProvisionalBalanceMeasure)
		default:
			return fmt.Errorf("invalid origin %s", algorithm.b.BillingBalance.Origin)
		}

	}
	return nil
}

type CloseSum struct {
	Name      string
	b         *BillingMeasure `bson:"-"`
	Magnitude measures.Magnitude
}

func NewCloseSum(b *BillingMeasure, magnitude measures.Magnitude) *CloseSum {
	return &CloseSum{
		Name:      "BALANCE_CALCULATED_BY_CLOSE_SUM",
		b:         b,
		Magnitude: magnitude,
	}
}

func (algorithm CloseSum) ID() string {
	return algorithm.Name
}

func (algorithm CloseSum) Execute(_ context.Context) error {
	worsePeriod := 0

	for _, period := range algorithm.b.Periods {
		algorithm.b.SumBalancePeriodMagnitude(measures.P0, algorithm.Magnitude, algorithm.b.GetBalancePeriodMagnitude(period, algorithm.Magnitude))

		algorithm.b.SetBalanceGeneralType(period, algorithm.Magnitude, GeneralReal)
		switch algorithm.b.BillingBalance.Origin {
		case measures.STM:
			algorithm.b.SetBalanceEstimateCode(period, algorithm.Magnitude, 1)
			algorithm.b.SetBalanceType(period, algorithm.Magnitude, RealByRemoteRead)
			algorithm.b.SetBalanceOrigin(period, algorithm.Magnitude, TlmOrigin)
			algorithm.b.SetBalanceMeasureType(period, algorithm.Magnitude, FirmBalanceMeasure)
		case measures.TPL, measures.File, measures.Manual, measures.Visual:
			algorithm.b.SetBalanceEstimateCode(period, algorithm.Magnitude, 3)
			algorithm.b.SetBalanceType(period, algorithm.Magnitude, RealByAbsLocalRead)
			algorithm.b.SetBalanceOrigin(period, algorithm.Magnitude, LocalOrigin)
			algorithm.b.SetBalanceMeasureType(period, algorithm.Magnitude, FirmBalanceMeasure)
		case measures.Auto:
			algorithm.b.SetBalanceEstimateCode(period, algorithm.Magnitude, 4)
			algorithm.b.SetBalanceType(period, algorithm.Magnitude, RealByAutoRead)
			algorithm.b.SetBalanceOrigin(period, algorithm.Magnitude, LocalOrigin)
			algorithm.b.SetBalanceMeasureType(period, algorithm.Magnitude, ProvisionalBalanceMeasure)
		}

		if worsePeriod < algorithm.b.GetBalanceEstimateCode(period, algorithm.Magnitude) {
			worsePeriod = algorithm.b.GetBalanceEstimateCode(period, algorithm.Magnitude)
		}

	}

	switch worsePeriod {
	case 1:
		algorithm.b.SetBalanceEstimateCode(measures.P0, algorithm.Magnitude, 2)
		algorithm.b.SetBalanceOrigin(measures.P0, algorithm.Magnitude, TlmOrigin)

		algorithm.b.SetBalanceMeasureType(measures.P0, algorithm.Magnitude, FirmBalanceMeasure)
	case 3:
		algorithm.b.SetBalanceEstimateCode(measures.P0, algorithm.Magnitude, 3)
		algorithm.b.SetBalanceOrigin(measures.P0, algorithm.Magnitude, LocalOrigin)

		algorithm.b.SetBalanceMeasureType(measures.P0, algorithm.Magnitude, FirmBalanceMeasure)
	case 4:
		algorithm.b.SetBalanceEstimateCode(measures.P0, algorithm.Magnitude, 4)
		algorithm.b.SetBalanceOrigin(measures.P0, algorithm.Magnitude, LocalOrigin)

		algorithm.b.SetBalanceMeasureType(measures.P0, algorithm.Magnitude, ProvisionalBalanceMeasure)
	}
	algorithm.b.SetBalanceType(measures.P0, algorithm.Magnitude, CalculatedByCloseSum)
	algorithm.b.SetBalanceGeneralType(measures.P0, algorithm.Magnitude, GeneralCalculated)
	return nil
}

type CCHCompleteNoTlg struct {
	Name      string
	b         *BillingMeasure `bson:"-"`
	Magnitude measures.Magnitude
}

func NewCCHCompleteNoTlg(b *BillingMeasure, magnitude measures.Magnitude) *CCHCompleteNoTlg {
	return &CCHCompleteNoTlg{
		Name:      "CCH_COMPLETE",
		b:         b,
		Magnitude: magnitude,
	}
}

func (algorithm CCHCompleteNoTlg) ID() string {
	return algorithm.Name
}

func (algorithm CCHCompleteNoTlg) Execute(_ context.Context) error {

	validOrigins := map[measures.OriginType]struct{}{
		measures.STM:    {},
		measures.TPL:    {},
		measures.File:   {},
		measures.Manual: {},
		measures.Visual: {},
	}

	for i, curve := range algorithm.b.BillingLoadCurve {

		if _, ok := validOrigins[curve.Origin]; !ok {
			continue
		}

		switch algorithm.b.PointType {
		case "1", "2":

			switch curve.Equipment {
			case measures.Main:
				algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 1)
				algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, FirmMainConfig)
			case measures.Redundant:
				algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 2)
				algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, FirmRedundantConfig)
			case measures.Receipt:
				algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 3)
				algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, FirmReceiptConfig)
			}

		case "3", "4", "5":
			algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 1)
			algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, RealBalance)
		default:
			return errors.New("invalid point type")
		}

		algorithm.b.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralReal)
		algorithm.b.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)
	}

	return nil
}

type PowerUseFactorCCH struct {
	Name      string
	B         *BillingMeasure `bson:"-"`
	Magnitude measures.Magnitude
}

func NewPowerUseFactorCCH(b *BillingMeasure, magnitude measures.Magnitude) *PowerUseFactorCCH {
	return &PowerUseFactorCCH{
		Name:      "CCH_POWER_USE_FACTOR",
		B:         b,
		Magnitude: magnitude,
	}
}

func (algorithm PowerUseFactorCCH) ID() string {
	return algorithm.Name
}

func (algorithm PowerUseFactorCCH) setMetadata(index int) {
	switch algorithm.B.PointType {
	case "1", "2":
		{
			if algorithm.B.GetLoadCurveOrigin(index) == measures.Filled {
				algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 11)
				algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, PowerUseFactor)
				algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralCalculated)
				algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, ProvisionalBalanceMeasure)
			} else {
				algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, FirmBalanceMeasure)
				algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralReal)
				switch algorithm.B.GetLoadCurveMeasurePointType(index) {
				case measures.MeasurePointTypeP:
					{
						algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, FirmMainConfig)
						algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 1)

					}
				case measures.MeasurePointTypeR:
					{
						algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, FirmRedundantConfig)
						algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 2)
					}
				case measures.MeasurePointTypeC:
					{
						algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, FirmReceiptConfig)
						algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 3)
					}

				}
			}
		}
	case "3", "4", "5":
		{
			if algorithm.B.GetLoadCurveOrigin(index) == measures.Filled {
				algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 6)
				algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, PowerUseFactor)
				algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralCalculated)
				algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, ProvisionalBalanceMeasure)
			} else {
				algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 1)
				algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, RealValidMeasure)
				algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralReal)
				algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, FirmBalanceMeasure)

			}
		}
	}
}

func (algorithm PowerUseFactorCCH) Execute(_ context.Context) error {
	for i, curve := range algorithm.B.BillingLoadCurve {
		algorithm.setMetadata(i)
		if curve.Origin == measures.Filled {
			powerFactor := math.Round(float64(algorithm.B.GetPeriodPowerDemanded(curve.Period)) * 0.33)
			algorithm.B.SetLoadCurvePeriodMagnitude(i, algorithm.Magnitude, powerFactor)
		}
	}
	return nil
}

type PowerUseFactorBalance struct {
	Name      string
	Magnitude measures.Magnitude
	B         *BillingMeasure `bson:"-"`
}

func NewPowerUseFactorBalance(b *BillingMeasure, magnitude measures.Magnitude) *PowerUseFactorBalance {
	return &PowerUseFactorBalance{
		Name:      "BALANCE_POWER_USE_FACTOR",
		Magnitude: magnitude,
		B:         b,
	}
}

func (algorithm PowerUseFactorBalance) ID() string {
	return algorithm.Name
}
func (algorithm PowerUseFactorBalance) setMetadataValidClosing(period measures.PeriodKey, magnitude measures.Magnitude) {
	SetMetadataValidClosing(algorithm.B, period, magnitude)
}

func SetMetadataValidClosing(bm *BillingMeasure, period measures.PeriodKey, magnitude measures.Magnitude) {
	bm.SetBalanceGeneralType(period, magnitude, GeneralReal)
	switch bm.BillingBalance.Origin {
	case measures.STM:
		{
			bm.SetBalanceEstimateCode(period, magnitude, 1)
			bm.SetBalanceType(period, magnitude, RealByRemoteRead)
			bm.SetBalanceOrigin(period, magnitude, TlmOrigin)
			bm.SetBalanceMeasureType(period, magnitude, FirmBalanceMeasure)
		}
	case measures.TPL, measures.File, measures.Manual, measures.Visual:
		{

			bm.SetBalanceEstimateCode(period, magnitude, 3)
			bm.SetBalanceType(period, magnitude, RealByAbsLocalRead)
			bm.SetBalanceOrigin(period, magnitude, LocalOrigin)
			bm.SetBalanceMeasureType(period, magnitude, FirmBalanceMeasure)
		}
	case measures.Auto:
		{
			bm.SetBalanceEstimateCode(period, magnitude, 4)
			bm.SetBalanceType(period, magnitude, RealByAutoRead)
			bm.SetBalanceOrigin(period, magnitude, LocalOrigin)
			bm.SetBalanceMeasureType(period, magnitude, ProvisionalBalanceMeasure)
		}
	}
}

func (algorithm PowerUseFactorBalance) Execute(_ context.Context) error {
	periods := algorithm.B.Periods
	countMeasuresPeriod := map[measures.PeriodKey]int{}
	for _, curve := range algorithm.B.BillingLoadCurve {
		countMeasuresPeriod[curve.Period]++
	}

	sumClosing := 0.0
	for _, period := range periods {
		if algorithm.B.GetBalanceStatus(period, algorithm.Magnitude) == measures.Invalid {
			// TODO: 0.33 Tiene que ser un valor parametrizable, Variable de entorno?
			balanceValue := 0.0
			if algorithm.Magnitude == measures.AI {
				balanceValue = math.Round(0.33 * float64(algorithm.B.GetPeriodPowerDemanded(period)) * float64(countMeasuresPeriod[period]))
			}

			algorithm.B.SetBalancePeriodMagnitude(period, algorithm.Magnitude, balanceValue)
			algorithm.B.SetBalanceType(period, algorithm.Magnitude, PowerUseFactor)
			algorithm.B.SetBalanceOrigin(period, algorithm.Magnitude, EstimateOrigin)
			algorithm.B.SetBalanceMeasureType(period, algorithm.Magnitude, ProvisionalBalanceMeasure)
			algorithm.B.SetBalanceEstimateCode(period, algorithm.Magnitude, 6)
			algorithm.B.SetBalanceGeneralType(period, algorithm.Magnitude, GeneralEstimated)
		} else {
			algorithm.setMetadataValidClosing(period, algorithm.Magnitude)
		}
		sumClosing = sumClosing + algorithm.B.GetBalancePeriodMagnitude(period, algorithm.Magnitude)
	}
	valueBalance := 0.0
	if algorithm.Magnitude == measures.AI {
		valueBalance = sumClosing
	}
	algorithm.B.SetBalancePeriodMagnitude(measures.P0, algorithm.Magnitude, valueBalance)
	algorithm.B.SetBalanceType(measures.P0, algorithm.Magnitude, CalculatedByCloseSum)
	algorithm.B.SetBalanceOrigin(measures.P0, algorithm.Magnitude, EstimateOrigin)
	algorithm.B.SetBalanceMeasureType(measures.P0, algorithm.Magnitude, ProvisionalBalanceMeasure)
	algorithm.B.SetBalanceEstimateCode(measures.P0, algorithm.Magnitude, 6)
	algorithm.B.SetBalanceGeneralType(measures.P0, algorithm.Magnitude, GeneralCalculated)
	return nil
}

type FillOneCloseAtr struct {
	Name      string
	b         *BillingMeasure `bson:"-"`
	Magnitude measures.Magnitude
}

func NewFillOneCloseAtr(b *BillingMeasure, magnitude measures.Magnitude) *FillOneCloseAtr {
	return &FillOneCloseAtr{
		Name:      "BALANCE_FILL_ONE_CLOSE_ATR",
		b:         b,
		Magnitude: magnitude,
	}
}

func (algorithm FillOneCloseAtr) ID() string {
	return algorithm.Name
}

func (algorithm FillOneCloseAtr) Execute(_ context.Context) error {
	var balance float64
	var fillPeriod measures.PeriodKey

	for _, period := range algorithm.b.GetPeriods() {
		if algorithm.b.GetBalanceStatus(period, algorithm.Magnitude) == measures.Invalid {
			fillPeriod = period
			continue
		}

		algorithm.b.SetBalanceGeneralType(period, algorithm.Magnitude, GeneralReal)
		switch algorithm.b.BillingBalance.Origin {
		case measures.STM:
			algorithm.b.SetBalanceEstimateCode(period, algorithm.Magnitude, 1)
			algorithm.b.SetBalanceType(period, algorithm.Magnitude, RealByRemoteRead)
			algorithm.b.SetBalanceOrigin(period, algorithm.Magnitude, TlmOrigin)
			algorithm.b.SetBalanceMeasureType(period, algorithm.Magnitude, FirmBalanceMeasure)
		case measures.TPL, measures.File, measures.Manual, measures.Visual:
			algorithm.b.SetBalanceEstimateCode(period, algorithm.Magnitude, 3)
			algorithm.b.SetBalanceType(period, algorithm.Magnitude, RealByAbsLocalRead)
			algorithm.b.SetBalanceOrigin(period, algorithm.Magnitude, LocalOrigin)
			algorithm.b.SetBalanceMeasureType(period, algorithm.Magnitude, FirmBalanceMeasure)
		case measures.Auto:
			algorithm.b.SetBalanceEstimateCode(period, algorithm.Magnitude, 4)
			algorithm.b.SetBalanceType(period, algorithm.Magnitude, RealByAutoRead)
			algorithm.b.SetBalanceOrigin(period, algorithm.Magnitude, LocalOrigin)
			algorithm.b.SetBalanceMeasureType(period, algorithm.Magnitude, ProvisionalBalanceMeasure)

		}

		if period == measures.P0 {
			continue
		}

		balance += algorithm.b.GetBalancePeriodMagnitude(period, algorithm.Magnitude)
	}

	algorithm.b.SetBalancePeriodMagnitude(
		fillPeriod,
		algorithm.Magnitude,
		algorithm.b.GetBalancePeriodMagnitude(measures.P0, algorithm.Magnitude)-balance,
	)
	algorithm.b.SetBalanceEstimateCode(fillPeriod, algorithm.Magnitude, 2)
	algorithm.b.SetBalanceType(fillPeriod, algorithm.Magnitude, CalculatedByCloseBalance)
	algorithm.b.SetBalanceGeneralType(fillPeriod, algorithm.Magnitude, GeneralCalculated)
	algorithm.b.SetBalanceOrigin(fillPeriod, algorithm.Magnitude, TlmOrigin)
	algorithm.b.SetBalanceMeasureType(fillPeriod, algorithm.Magnitude, FirmBalanceMeasure)

	return nil
}

type ClosingWithBalance struct {
	Name      string
	Magnitude measures.Magnitude
	B         *BillingMeasure `bson:"-"`
}

func NewClosingWithBalance(b *BillingMeasure, Magnitude measures.Magnitude) *ClosingWithBalance {
	return &ClosingWithBalance{
		Name:      "BALANCE_CLOSING",
		B:         b,
		Magnitude: Magnitude,
	}
}

func (algorithm ClosingWithBalance) ID() string {
	return algorithm.Name
}

func (algorithm ClosingWithBalance) setMetadataValidClosing(period measures.PeriodKey, magnitude measures.Magnitude) {
	SetMetadataValidClosing(algorithm.B, period, magnitude)
}

func (algorithm ClosingWithBalance) Execute(_ context.Context) error {
	periods := algorithm.B.Periods
	sumClosing := 0.0
	sumPowerDemanded := 0.0
	for _, period := range periods {
		sumClosing = sumClosing + algorithm.B.GetBalancePeriodMagnitude(period, algorithm.Magnitude)
		sumPowerDemanded = sumPowerDemanded + algorithm.B.GetPeriodPowerDemanded(period)
	}

	for _, period := range periods {
		if algorithm.B.GetBalanceStatus(period, algorithm.Magnitude) == measures.Invalid {
			valueMagnitude := 0.
			if algorithm.Magnitude == measures.AI && sumPowerDemanded != 0.0 {
				valueMagnitude = math.Round((algorithm.B.GetBalancePeriodMagnitude(measures.P0, algorithm.Magnitude) - sumClosing) * algorithm.B.GetPeriodPowerDemanded(period) / sumPowerDemanded)
			}
			algorithm.B.SetBalancePeriodMagnitude(period, algorithm.Magnitude, valueMagnitude)
			algorithm.B.SetBalanceType(period, algorithm.Magnitude, EstimatedContractPower)
			algorithm.B.SetBalanceGeneralType(period, algorithm.Magnitude, GeneralEstimated)
			algorithm.B.SetBalanceOrigin(period, algorithm.Magnitude, EstimateOrigin)
			algorithm.B.SetBalanceMeasureType(period, algorithm.Magnitude, ProvisionalBalanceMeasure)
			algorithm.B.SetBalanceEstimateCode(period, algorithm.Magnitude, 6)
		} else {
			algorithm.setMetadataValidClosing(period, algorithm.Magnitude)
		}
	}
	algorithm.setMetadataValidClosing(measures.P0, algorithm.Magnitude)
	return nil
}

type FlatCastBalanceNoTLG struct {
	Name      string
	b         *BillingMeasure `bson:"-"`
	Magnitude measures.Magnitude
}

func NewFlatCastBalanceNoTLG(b *BillingMeasure, magnitude measures.Magnitude) *FlatCastBalanceNoTLG {
	return &FlatCastBalanceNoTLG{
		Name:      "CCH_FLAT_CAST_BALANCE_NO_TLG",
		b:         b,
		Magnitude: magnitude,
	}
}

func (algorithm FlatCastBalanceNoTLG) ID() string {
	return algorithm.Name
}

func (algorithm FlatCastBalanceNoTLG) Execute(_ context.Context) error {
	nGap := 0.0
	value := 0.0

	for i, curve := range algorithm.b.BillingLoadCurve {
		if curve.Origin != measures.Filled {
			value += algorithm.b.GetLoadCurvePeriodMagnitude(i, algorithm.Magnitude)
			continue
		}
		nGap++
	}

	if nGap == 0.0 {
		nGap = 1.0
	}

	value = (algorithm.b.GetBalancePeriodMagnitude(measures.P0, algorithm.Magnitude) - value) / nGap

	for i, curve := range algorithm.b.BillingLoadCurve {
		if curve.Origin == measures.Filled {
			algorithm.b.SetLoadCurvePeriodMagnitude(i, algorithm.Magnitude, value)
		}

		if algorithm.b.PointType == "1" || algorithm.b.PointType == "2" {
			if curve.Origin == measures.Filled {
				algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 8)
				algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, EstimatedByFlatProfile)
				algorithm.b.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralEstimated)
				algorithm.b.SetLoadCurveMeasureType(i, algorithm.Magnitude, ProvisionalBalanceMeasure)

				continue
			}
			switch curve.Equipment {
			case measures.Main:
				algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 1)
				algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, FirmMainConfig)
				algorithm.b.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)
			case measures.Redundant:
				algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 2)
				algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, FirmRedundantConfig)
				algorithm.b.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)
			default:
				algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 3)
				algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, FirmReceiptConfig)
				algorithm.b.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)
			}

			algorithm.b.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralReal)
			continue
		}

		if curve.Origin != measures.Filled {
			algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 1)
			algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, RealValidMeasure)
			algorithm.b.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)
			algorithm.b.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralReal)
			continue
		}

		algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 3)
		algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, RealMeasureAdjustment)
		algorithm.b.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralAdjusted)
		algorithm.b.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)

	}

	return nil
}

type SumHoursNoClosure struct {
	Name      string
	b         *BillingMeasure `bson:"-"`
	Magnitude measures.Magnitude
}

func NewSumHoursNoClosure(b *BillingMeasure, magnitude measures.Magnitude) *SumHoursNoClosure {
	return &SumHoursNoClosure{
		Name:      "BALANCE_SUM_HOURS_NO_CLOSURE",
		b:         b,
		Magnitude: magnitude,
	}
}
func (algorithm SumHoursNoClosure) ID() string {
	return algorithm.Name
}

func (algorithm SumHoursNoClosure) Execute(_ context.Context) error {
	worsePeriod := 0
	worstCurvePeriod := map[measures.PeriodKey]struct {
		estimatedCode int
		origin        measures.OriginType
	}{}

	for i, curve := range algorithm.b.BillingLoadCurve {

		curveMagnitudeEstimateCode, err := algorithm.b.GetLoadCurveEstimatedCode(i, algorithm.Magnitude)
		if err != nil {
			return err
		}

		balancePeriodValidation := algorithm.b.GetBalanceStatus(curve.Period, algorithm.Magnitude)

		if balancePeriodValidation == measures.Invalid {
			algorithm.b.SumBalancePeriodMagnitude(curve.Period, algorithm.Magnitude, curve.GetMagnitude(algorithm.Magnitude))
		}

		if _, ok := worstCurvePeriod[curve.Period]; ok {
			wp := worstCurvePeriod[curve.Period]
			if curveMagnitudeEstimateCode >= wp.estimatedCode {
				wp.estimatedCode = curveMagnitudeEstimateCode
				wp.origin = curve.Origin
			}
			continue
		}

		worstCurvePeriod[curve.Period] = struct {
			estimatedCode int
			origin        measures.OriginType
		}{
			estimatedCode: curveMagnitudeEstimateCode,
			origin:        curve.Origin,
		}
	}
	for _, period := range algorithm.b.Periods {

		algorithm.b.SumBalancePeriodMagnitude(measures.P0, algorithm.Magnitude, algorithm.b.GetBalancePeriodMagnitude(period, algorithm.Magnitude))
		balancePeriodValidation := algorithm.b.GetBalanceStatus(period, algorithm.Magnitude)

		if balancePeriodValidation == measures.Invalid {
			algorithm.SetBalanceValuesOnEmptyCase(period, worstCurvePeriod[period].estimatedCode, worstCurvePeriod[period].origin)
			if worsePeriod < algorithm.b.GetBalanceEstimateCode(period, algorithm.Magnitude) {
				worsePeriod = algorithm.b.GetBalanceEstimateCode(period, algorithm.Magnitude)
			}
			continue
		}

		if worsePeriod < algorithm.b.GetBalanceEstimateCode(period, algorithm.Magnitude) {
			worsePeriod = algorithm.b.GetBalanceEstimateCode(period, algorithm.Magnitude)
		}
		switch algorithm.b.BillingBalance.Origin {
		case measures.STM:
			algorithm.SetBalancePeriod(period, 1, TlmOrigin, FirmBalanceMeasure, RealByRemoteRead, GeneralReal)
		case measures.TPL, measures.File, measures.Manual, measures.Visual:
			algorithm.SetBalancePeriod(period, 3, LocalOrigin, FirmBalanceMeasure, RealByAbsLocalRead, GeneralReal)
		case measures.Auto:
			algorithm.SetBalancePeriod(period, 4, LocalOrigin, ProvisionalBalanceMeasure, RealByAutoRead, GeneralReal)
		default:
			return errors.New("unexpected case")
		}
	}

	balancePeriodValidation := algorithm.b.GetBalanceStatus(measures.P0, algorithm.Magnitude)

	if balancePeriodValidation != measures.Invalid {
		switch algorithm.b.BillingBalance.Origin {
		case measures.STM:
			algorithm.SetBalancePeriod(measures.P0, 1, TlmOrigin, FirmBalanceMeasure, RealByRemoteRead, GeneralReal)
		case measures.TPL, measures.File, measures.Manual, measures.Visual:
			algorithm.SetBalancePeriod(measures.P0, 3, LocalOrigin, FirmBalanceMeasure, RealByAbsLocalRead, GeneralReal)
		case measures.Auto:
			algorithm.SetBalancePeriod(measures.P0, 4, LocalOrigin, ProvisionalBalanceMeasure, RealByAutoRead, GeneralReal)
		}
		return nil
	}
	switch worsePeriod {
	case 1:
		algorithm.SetBalancePeriod(measures.P0, 2, TlmOrigin, FirmBalanceMeasure, CalculatedByCloseSum, GeneralCalculated)
	case 3:
		algorithm.SetBalancePeriod(measures.P0, 3, LocalOrigin, FirmBalanceMeasure, CalculatedByCloseSum, GeneralCalculated)
	case 4:
		algorithm.SetBalancePeriod(measures.P0, 4, LocalOrigin, ProvisionalBalanceMeasure, CalculatedByCloseSum, GeneralCalculated)
	case 5:
		algorithm.SetBalancePeriod(measures.P0, 5, EstimateOrigin, ProvisionalBalanceMeasure, CalculatedByCloseSum, GeneralCalculated)
	case 6:
		algorithm.SetBalancePeriod(measures.P0, 6, EstimateOrigin, ProvisionalBalanceMeasure, CalculatedByCloseSum, GeneralCalculated)
	}

	return nil
}

func (algorithm SumHoursNoClosure) SetBalancePeriod(
	period measures.PeriodKey,
	estimateCode int,
	balanceOrigin BalanceOriginType,
	balanceMeasureType BalanceMeasureType,
	balanceType BalanceType, generalType GeneralEstimateMethod) {
	algorithm.b.SetBalanceEstimateCode(period, algorithm.Magnitude, estimateCode)
	algorithm.b.SetBalanceOrigin(period, algorithm.Magnitude, balanceOrigin)
	algorithm.b.SetBalanceMeasureType(period, algorithm.Magnitude, balanceMeasureType)
	algorithm.b.SetBalanceType(period, algorithm.Magnitude, balanceType)
	algorithm.b.SetBalanceGeneralType(period, algorithm.Magnitude, generalType)
}

func (algorithm SumHoursNoClosure) SetBalanceValuesOnEmptyCase(period measures.PeriodKey, worstEstimatedCode int, worstOrigin measures.OriginType) error {
	if algorithm.b.isPointType([]measures.PointType{"1", "2"}) {
		switch {
		case utils.InSlice(worstOrigin, []measures.OriginType{measures.STM, measures.Filled}) &&
			utils.InSlice(worstEstimatedCode, []int{1, 2, 3, 8}):
			algorithm.SetBalancePeriod(period, 1, TlmOrigin, FirmBalanceMeasure, ObtainedByCurve, GeneralCalculated)

		case utils.InSlice(worstOrigin, []measures.OriginType{measures.TPL, measures.File, measures.Manual, measures.Visual}) &&
			utils.InSlice(worstEstimatedCode, []int{1, 2, 3}):
			algorithm.SetBalancePeriod(period, 3, LocalOrigin, FirmBalanceMeasure, ObtainedByCurve, GeneralCalculated)

		case utils.InSlice(worstOrigin, []measures.OriginType{measures.Filled}) &&
			utils.InSlice(worstEstimatedCode, []int{11}):
			algorithm.SetBalancePeriod(period, 6, EstimateOrigin, ProvisionalBalanceMeasure, ObtainedByCurve, GeneralCalculated)

		case utils.InSlice(worstOrigin, []measures.OriginType{measures.Filled}) &&
			utils.InSlice(worstEstimatedCode, []int{9, 7}):
			algorithm.SetBalancePeriod(period, 6, EstimateOrigin, ProvisionalBalanceMeasure, ObtainedByCurve, GeneralCalculated)

		default:
			return errors.New("unexpected Case")
		}

		return nil
	}

	if algorithm.b.isPointType([]measures.PointType{"3", "4"}) {
		switch {
		case utils.InSlice(worstOrigin, []measures.OriginType{measures.STM}) &&
			utils.InSlice(worstEstimatedCode, []int{1, 3}):
			algorithm.SetBalancePeriod(period, 1, TlmOrigin, FirmBalanceMeasure, ObtainedByCurve, GeneralCalculated)

		case utils.InSlice(worstOrigin, []measures.OriginType{measures.TPL, measures.File, measures.Manual, measures.Visual}) &&
			utils.InSlice(worstEstimatedCode, []int{1}):
			algorithm.SetBalancePeriod(period, 3, LocalOrigin, FirmBalanceMeasure, ObtainedByCurve, GeneralCalculated)

		case utils.InSlice(worstOrigin, []measures.OriginType{measures.Filled}) &&
			utils.InSlice(worstEstimatedCode, []int{5}):
			algorithm.SetBalancePeriod(period, 5, EstimateOrigin, ProvisionalBalanceMeasure, ObtainedByCurve, GeneralCalculated)

		case utils.InSlice(worstOrigin, []measures.OriginType{measures.Filled}) &&
			utils.InSlice(worstEstimatedCode, []int{6}):
			algorithm.SetBalancePeriod(period, 6, EstimateOrigin, ProvisionalBalanceMeasure, ObtainedByCurve, GeneralCalculated)

		default:
			return errors.New("unexpected Case")
		}
	}
	return errors.New("unexpected Case")
}

type BalanceZeroConsumption struct {
	Name      string
	Magnitude measures.Magnitude
	b         *BillingMeasure `bson:"-"`
}

func NewBalanceZeroConsumption(b *BillingMeasure, magnitude measures.Magnitude) *BalanceZeroConsumption {
	return &BalanceZeroConsumption{
		Name:      "BALANCE_ZERO_CONSUMPTION",
		Magnitude: magnitude,
		b:         b,
	}
}

func (algorithm BalanceZeroConsumption) ID() string {
	return algorithm.Name
}

func (algorithm BalanceZeroConsumption) Execute(_ context.Context) error {

	for _, period := range algorithm.b.GetPeriods() {
		algorithm.b.SetBalancePeriodMagnitude(period, algorithm.Magnitude, 0)
		algorithm.b.SetBalanceEstimateCode(period, algorithm.Magnitude, 5)
		algorithm.b.SetBalanceOrigin(period, algorithm.Magnitude, EstimateOrigin)
		algorithm.b.SetBalanceType(period, algorithm.Magnitude, EstimateOnlyHistoric)
		algorithm.b.SetBalanceMeasureType(period, algorithm.Magnitude, ProvisionalBalanceMeasure)
		algorithm.b.SetBalanceGeneralType(period, algorithm.Magnitude, GeneralEstimated)
	}
	return nil
}

type ConsumZeroClosedHouseNoTLG struct {
	Name      string
	b         *BillingMeasure `bson:"-"`
	Magnitude measures.Magnitude
}

func NewConsumZeroClosedHouseNoTLG(b *BillingMeasure, magnitude measures.Magnitude) *ConsumZeroClosedHouseNoTLG {
	return &ConsumZeroClosedHouseNoTLG{
		Name:      "CCH_CONSUM_ZERO_CLOSED_HOUSE_NO_TLG",
		b:         b,
		Magnitude: magnitude,
	}
}

func (algorithm ConsumZeroClosedHouseNoTLG) ID() string {
	return algorithm.Name
}

func (algorithm ConsumZeroClosedHouseNoTLG) Execute(_ context.Context) error {
	for i := range algorithm.b.BillingLoadCurve {

		estimateCode := 5
		estimateMethod := EstimateByHistoricConsumLastYear

		if algorithm.b.PointType == "1" || algorithm.b.PointType == "2" {
			estimateCode = 9
			estimateMethod = EstimateHistoricMainMeasurePoint
		}

		algorithm.b.SetLoadCurvePeriodMagnitude(i, algorithm.Magnitude, 0.0)
		algorithm.b.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, estimateCode)
		algorithm.b.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, estimateMethod)
		algorithm.b.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralEstimated)
		algorithm.b.SetLoadCurveMeasureType(i, algorithm.Magnitude, ProvisionalBalanceMeasure)
	}

	return nil
}

type CloseHistoryWithoutBalance struct {
	Name      string
	Magnitude measures.Magnitude
	Context   *GraphContext   `bson:"-"`
	b         *BillingMeasure `bson:"-"`
}

func NewCloseHistoryWithoutBalance(b *BillingMeasure, magnitude measures.Magnitude, context *GraphContext) *CloseHistoryWithoutBalance {
	return &CloseHistoryWithoutBalance{
		Name:      "BALANCE_CLOSE_HISTORIC_WITHOUT_BALANCE",
		b:         b,
		Context:   context,
		Magnitude: magnitude,
	}
}

func (algorithm CloseHistoryWithoutBalance) ID() string {
	return algorithm.Name
}

func (algorithm CloseHistoryWithoutBalance) calcTypicalDeviation(values []float64, periodAverage float64) float64 {
	typicalDeviation := 0.
	for _, value := range values {
		typicalDeviation += math.Pow(value+periodAverage, 2)
	}

	typicalDeviation /= float64(len(values))
	return typicalDeviation
}

func (algorithm CloseHistoryWithoutBalance) calcAverage(values []float64) float64 {
	average := 0.
	for _, value := range values {
		average += value
	}

	average /= float64(len(values))

	return average
}

func (algorithm CloseHistoryWithoutBalance) sumCurvePeriodHours(loadCurve []BillingLoadCurve, periodsToEstimate map[measures.PeriodKey][]float64) map[measures.PeriodKey]float64 {
	hoursPerPeriod := make(map[measures.PeriodKey]float64)

	for period := range periodsToEstimate {
		hoursPerPeriod[period] = 0
	}

	for _, curve := range loadCurve {
		if _, ok := periodsToEstimate[curve.Period]; !ok {
			continue
		}
		hoursPerPeriod[curve.Period] += 1
	}

	return hoursPerPeriod
}

func (algorithm CloseHistoryWithoutBalance) setMetadata(period measures.PeriodKey) {

	var (
		estimatedCode      = 5
		balanceType        = EstimateByHistoricLastYear
		balanceOrigin      = EstimateOrigin
		balanceMeasureType = ProvisionalBalanceMeasure
		generalType        = GeneralEstimated
	)

	if algorithm.b.GetBalanceStatus(period, algorithm.Magnitude) == measures.Valid {
		switch algorithm.b.BillingBalance.Origin {
		case measures.STM:
			estimatedCode = 1
			balanceType = RealByRemoteRead
			balanceOrigin = TlmOrigin
			balanceMeasureType = FirmBalanceMeasure
			generalType = GeneralReal
		case measures.TPL, measures.File, measures.Manual, measures.Visual:
			estimatedCode = 3
			balanceType = RealByAbsLocalRead
			balanceOrigin = LocalOrigin
			balanceMeasureType = FirmBalanceMeasure
			generalType = GeneralReal
		case measures.Auto:
			estimatedCode = 4
			balanceType = RealByAutoRead
			balanceOrigin = LocalOrigin
			balanceMeasureType = ProvisionalBalanceMeasure
			generalType = GeneralReal
		}
	}

	if period == measures.P0 {
		estimatedCode = 5
		balanceType = CalculatedByCloseSum
		balanceOrigin = EstimateOrigin
		balanceMeasureType = ProvisionalBalanceMeasure
		generalType = GeneralCalculated
	}

	algorithm.b.SetBalanceEstimateCode(period, algorithm.Magnitude, estimatedCode)
	algorithm.b.SetBalanceType(period, algorithm.Magnitude, balanceType)
	algorithm.b.SetBalanceOrigin(period, algorithm.Magnitude, balanceOrigin)
	algorithm.b.SetBalanceMeasureType(period, algorithm.Magnitude, balanceMeasureType)
	algorithm.b.SetBalanceGeneralType(period, algorithm.Magnitude, generalType)

}

func (algorithm CloseHistoryWithoutBalance) calcHourlyAverage(periodsToEstimate map[measures.PeriodKey][]float64) {
	for _, historic := range algorithm.Context.ClosedHistory {
		hoursPerPeriod := algorithm.sumCurvePeriodHours(historic.BillingLoadCurve, periodsToEstimate)

		for period, hours := range hoursPerPeriod {
			magnitudeValue := historic.GetBalancePeriodMagnitude(period, algorithm.Magnitude)

			if hours == 0 {
				hours = 1
			}

			periodsToEstimate[period] = append(periodsToEstimate[period], magnitudeValue/hours)
		}
	}
}

func (algorithm CloseHistoryWithoutBalance) calcAverageAndDeviation(periodsToEstimate map[measures.PeriodKey][]float64) (map[measures.PeriodKey]float64, map[measures.PeriodKey]float64, map[measures.PeriodKey]float64) {
	average := make(map[measures.PeriodKey]float64)
	maxPeriodValue := make(map[measures.PeriodKey]float64)
	minPeriodValue := make(map[measures.PeriodKey]float64)
	for period, values := range periodsToEstimate {
		averageValue := algorithm.calcAverage(values)
		typicalDeviationValue := algorithm.calcTypicalDeviation(values, averageValue)
		average[period] = averageValue
		maxPeriodValue[period] = averageValue + 2*typicalDeviationValue
		minPeriodValue[period] = averageValue - 2*typicalDeviationValue
	}

	return average, maxPeriodValue, minPeriodValue
}

func (algorithm CloseHistoryWithoutBalance) getNormalizedValues(periodsToEstimate map[measures.PeriodKey][]float64, maxValues, minValues map[measures.PeriodKey]float64) map[measures.PeriodKey]float64 {
	normalizedAverage := make(map[measures.PeriodKey]float64)

	for period, values := range periodsToEstimate {
		maxValue := maxValues[period]
		minValue := minValues[period]
		validValues := make([]float64, 0, 4)
		for _, value := range values {
			if value > maxValue || value < minValue {
				continue
			}
			validValues = append(validValues, value)
		}
		normalizedAverage[period] = algorithm.calcAverage(validValues)
	}

	return normalizedAverage
}

func (algorithm CloseHistoryWithoutBalance) getPeriodsToEstimate() map[measures.PeriodKey][]float64 {

	periodsToEstimate := make(map[measures.PeriodKey][]float64)
	for _, period := range algorithm.b.Periods {
		if algorithm.b.GetBalanceStatus(period, algorithm.Magnitude) != measures.Invalid {
			algorithm.setMetadata(period)
			continue
		}

		periodsToEstimate[period] = make([]float64, 0, 4)
	}

	return periodsToEstimate
}

func (algorithm CloseHistoryWithoutBalance) sumCloseAtr() {
	for _, period := range algorithm.b.Periods {
		balanceValue := algorithm.b.GetBalancePeriodMagnitude(period, algorithm.Magnitude)
		algorithm.b.SumBalancePeriodMagnitude(measures.P0, algorithm.Magnitude, balanceValue)
	}
}

func (algorithm CloseHistoryWithoutBalance) calcEstimatePeriods(period measures.PeriodKey, average float64, curveHours float64) {
	magnitudeValue := math.Round(average * curveHours)
	algorithm.b.SetBalancePeriodMagnitude(period, algorithm.Magnitude, magnitudeValue)
	algorithm.setMetadata(period)
}

func (algorithm CloseHistoryWithoutBalance) Execute(_ context.Context) error {
	periodsToEstimate := algorithm.getPeriodsToEstimate()

	algorithm.calcHourlyAverage(periodsToEstimate)

	average, maxPeriodValue, minPeriodValue := algorithm.calcAverageAndDeviation(periodsToEstimate)

	average = algorithm.getNormalizedValues(periodsToEstimate, maxPeriodValue, minPeriodValue)

	hoursPerPeriod := algorithm.sumCurvePeriodHours(algorithm.b.BillingLoadCurve, periodsToEstimate)

	for period := range periodsToEstimate {
		algorithm.calcEstimatePeriods(period, average[period], hoursPerPeriod[period])
	}

	algorithm.sumCloseAtr()
	algorithm.setMetadata(measures.P0)

	return nil
}

type ReeBalanceOutline struct {
	Name      string
	B         *BillingMeasure `bson:"-"`
	Magnitude measures.Magnitude
	Profiles  ConsumProfileRepository `bson:"-"`
}

func NewReeBalanceOutline(b *BillingMeasure, profiles ConsumProfileRepository, magnitude measures.Magnitude) *ReeBalanceOutline {
	return &ReeBalanceOutline{B: b, Name: "CCH_REE_BALANCE_OUTLINE", Profiles: profiles, Magnitude: magnitude}
}

func (algorithm ReeBalanceOutline) ID() string {
	return algorithm.Name
}

func (algorithm ReeBalanceOutline) Execute(ctx context.Context) error {
	worstCase := 0

	for _, period := range algorithm.B.Periods {
		if algorithm.B.GetBalanceEstimateCode(period, algorithm.Magnitude) > worstCase {
			worstCase = algorithm.B.GetBalanceEstimateCode(period, algorithm.Magnitude)
		}
	}

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
	periodSumCProfile := algorithm.getSumProfiles(keysProfiles)

	for i, curve := range algorithm.B.BillingLoadCurve {
		curveValue := 0.0
		balanceValue := algorithm.B.GetBalancePeriodMagnitude(curve.Period, algorithm.Magnitude)
		if cf, ok := keysProfiles[curve.EndDate]; ok {

			if periodSumCProfile[curve.Period] != 0 {
				coefValue := cf.GetValueByCoefficientType(algorithm.B.Coefficient)
				curveValue = balanceValue * (coefValue / periodSumCProfile[curve.Period])
			}
		}

		algorithm.B.SetLoadCurvePeriodMagnitude(i, algorithm.Magnitude, math.Round(curveValue))
		algorithm.setMetaData(i, worstCase)
	}

	return nil
}

func (algorithm ReeBalanceOutline) getSumProfiles(keysProfiles map[time.Time]ConsumProfile) map[measures.PeriodKey]float64 {

	sumPeriodProfiles := make(map[measures.PeriodKey]float64)

	for _, curve := range algorithm.B.BillingLoadCurve {
		if cf, ok := keysProfiles[curve.EndDate]; ok {
			sumPeriodProfiles[curve.Period] += cf.GetValueByCoefficientType(algorithm.B.Coefficient)
		}
	}
	return sumPeriodProfiles
}

func (algorithm ReeBalanceOutline) setMetaData(index int, worstCaseBalance int) {
	switch worstCaseBalance {
	case 1, 2, 3:
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 2)
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, Outlined)
		algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, ProvisionalBalanceMeasure)
	case 4:
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 4)
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, AutoOutlined)
		algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, ProvisionalBalanceMeasure)
	case 5:
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 5)
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, EstimatedByHistoricOutlinedConsum)
		algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, ProvisionalBalanceMeasure)
	case 6:
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 6)
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, EstimatedByOutlinedFactor)
		algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, ProvisionalBalanceMeasure)
	}

	if worstCaseBalance <= 4 {
		algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralOutlined)
		return
	}

	algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralEstimated)
}

type ClosingHistoryWithBalance struct {
	Name      string
	Magnitude measures.Magnitude
	B         *BillingMeasure `bson:"-"`
	Context   *GraphContext   `bson:"-"`
}

func NewClosingHistoryWithBalance(b *BillingMeasure, magnitude measures.Magnitude, context *GraphContext) *ClosingHistoryWithBalance {
	return &ClosingHistoryWithBalance{
		Name:      "BALANCE_CLOSE_HISTORY_WITH_BALANCE",
		Magnitude: magnitude,
		B:         b,
		Context:   context,
	}
}

func (algorithm ClosingHistoryWithBalance) ID() string {
	return algorithm.Name
}
func (algorithm ClosingHistoryWithBalance) setMetadataValidClosing(period measures.PeriodKey, magnitude measures.Magnitude) {
	SetMetadataValidClosing(algorithm.B, period, magnitude)
}

func (algorithm ClosingHistoryWithBalance) estimateClosePeriods() {
	closeWithoutBalance := NewCloseHistoryWithoutBalance(algorithm.B, algorithm.Magnitude, algorithm.Context)

	periodsToEstimate := closeWithoutBalance.getPeriodsToEstimate()
	closeWithoutBalance.calcHourlyAverage(periodsToEstimate)
	average, maxPeriodValue, minPeriodValue := closeWithoutBalance.calcAverageAndDeviation(periodsToEstimate)

	average = closeWithoutBalance.getNormalizedValues(periodsToEstimate, maxPeriodValue, minPeriodValue)

	hoursPerPeriod := closeWithoutBalance.sumCurvePeriodHours(algorithm.B.BillingLoadCurve, periodsToEstimate)

	for period := range periodsToEstimate {
		closeWithoutBalance.calcEstimatePeriods(period, average[period], hoursPerPeriod[period])
	}
}

func (algorithm ClosingHistoryWithBalance) Execute(_ context.Context) error {
	algorithm.estimateClosePeriods()
	periods := algorithm.B.Periods
	sumValidClosing := 0.0
	sumEstimatedClosing := 0.0
	for _, period := range periods {
		balanceValidation := algorithm.B.GetBalanceStatus(period, algorithm.Magnitude)
		magnitudeValue := algorithm.B.GetBalancePeriodMagnitude(period, algorithm.Magnitude)
		if balanceValidation == measures.Valid {
			sumValidClosing += magnitudeValue
			continue
		}

		sumEstimatedClosing += magnitudeValue
	}

	if sumEstimatedClosing == 0 {
		sumEstimatedClosing = 1
	}

	for _, period := range periods {
		balanceValidation := algorithm.B.GetBalanceStatus(period, algorithm.Magnitude)

		if balanceValidation == measures.Valid {
			algorithm.setMetadataValidClosing(period, algorithm.Magnitude)
			continue
		}

		balanceValue := algorithm.B.GetBalancePeriodMagnitude(measures.P0, algorithm.Magnitude)
		magnitudeValue := algorithm.B.GetBalancePeriodMagnitude(period, algorithm.Magnitude)
		estimatedValue := math.Round((balanceValue - sumValidClosing) * magnitudeValue / sumEstimatedClosing)
		algorithm.B.SetBalancePeriodMagnitude(period, algorithm.Magnitude, estimatedValue)

		algorithm.B.SetBalanceEstimateCode(period, algorithm.Magnitude, 5)
		algorithm.B.SetBalanceOrigin(period, algorithm.Magnitude, EstimateOrigin)
		algorithm.B.SetBalanceType(period, algorithm.Magnitude, EstimateByHistoricLastYear)
		algorithm.B.SetBalanceGeneralType(period, algorithm.Magnitude, GeneralEstimated)
		algorithm.B.SetBalanceMeasureType(period, algorithm.Magnitude, ProvisionalBalanceMeasure)
	}

	algorithm.setMetadataValidClosing(measures.P0, algorithm.Magnitude)
	return nil

}

type EstimatedHistoryNoTlg struct {
	Name         string
	b            *BillingMeasure `bson:"-"`
	Magnitude    measures.Magnitude
	GraphContext *GraphContext `bson:"-"`
}

func NewEstimatedHistoryNoTlg(
	b *BillingMeasure,
	magnitude measures.Magnitude,
	contextTlg *GraphContext,
) *EstimatedHistoryNoTlg {
	return &EstimatedHistoryNoTlg{
		Name:         "BALANCE_ESTIMATED_HISTORY_TLG",
		b:            b,
		Magnitude:    magnitude,
		GraphContext: contextTlg,
	}
}

func (algorithm EstimatedHistoryNoTlg) ID() string {
	return algorithm.Name
}

func (algorithm EstimatedHistoryNoTlg) Execute(ctx context.Context) error {

	for _, period := range algorithm.b.Periods {
		_ = NewEstimatedHistoryTlg(algorithm.b, period, algorithm.Magnitude, algorithm.GraphContext).Execute(ctx)
	}

	return nil
}

func calcAverage(values ...float64) float64 {
	var average float64
	totalValues := len(values)

	if totalValues == 0 {
		return 0
	}

	for _, value := range values {
		average += value
	}

	average /= float64(len(values))

	return average
}

func calcTypicalDeviation(average float64, values ...float64) float64 {
	typicalDeviation := 0.
	totalValues := len(values)

	if totalValues == 0 {
		return 0
	}

	for _, value := range values {
		typicalDeviation += math.Pow(value+average, 2)
	}

	typicalDeviation = math.Sqrt(typicalDeviation / float64(len(values)))
	return typicalDeviation
}

func getValidRangeValues(average float64, typicalDeviation float64) (float64, float64) {
	var maxValidValue, minValidValue float64

	maxValidValue = average + (2 * typicalDeviation)
	minValidValue = average - (2 * typicalDeviation)

	return maxValidValue, minValidValue
}

func getIsInValidRange(max, min float64, values ...float64) []float64 {
	validValues := make([]float64, 0, len(values))
	for _, value := range values {
		if !isInValidRange(max, min, value) {
			continue
		}

		validValues = append(validValues, value)
	}

	return validValues
}

func isInValidRange(maxValidValue, minValidValue, value float64) bool {
	return value <= maxValidValue && value >= minValidValue
}

type CCHWindows struct {
	Name         string
	B            *BillingMeasure `bson:"-"`
	Magnitude    measures.Magnitude
	GraphContext *GraphContext `bson:"-"`
}

func NewCCHWindows(b *BillingMeasure, magnitude measures.Magnitude, context *GraphContext) *CCHWindows {
	return &CCHWindows{
		Name:         "CCH_WINDOWS",
		B:            b,
		Magnitude:    magnitude,
		GraphContext: context,
	}
}

func (algorithm CCHWindows) ID() string {
	return algorithm.Name
}

func (algorithm CCHWindows) getPeriodHoursToEstimate() map[measures.PeriodKey]int {
	periodHoursToEstimate := make(map[measures.PeriodKey]int)
	for _, curve := range algorithm.B.BillingLoadCurve {
		if curve.Origin != measures.Filled {
			continue
		}

		periodHoursToEstimate[curve.Period] += 1
	}

	return periodHoursToEstimate
}

func (algorithm CCHWindows) getMagnitude(processedLoadCurve process_measures.ProcessedLoadCurve) float64 {
	switch algorithm.Magnitude {
	case measures.AI:
		return processedLoadCurve.AI
	case measures.AE:
		return processedLoadCurve.AE
	case measures.R1:
		return processedLoadCurve.R1
	case measures.R2:
		return processedLoadCurve.R2
	case measures.R3:
		return processedLoadCurve.R3
	case measures.R4:
		return processedLoadCurve.R4
	default:
		return 0
	}
}

func (algorithm CCHWindows) setValidMetadata(index int) {
	curveEquipment := algorithm.B.GetLoadCurveEquipment(index)

	algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, FirmBalanceMeasure)
	algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralReal)
	if algorithm.B.isPointType([]measures.PointType{"3", "4", "5"}) {
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 1)
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, RealBalance)
		return
	}

	switch curveEquipment {
	case measures.Main:
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 1)
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, FirmMainConfig)
	case measures.Redundant:
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 2)
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, FirmRedundantConfig)
	case measures.Receipt:
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 3)
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, FirmReceiptConfig)
	}
	return

}

func (algorithm CCHWindows) getUtilValues(values ...float64) []float64 {
	newValues := make([]float64, 0, len(values))

	if len(values) < 3 {
		return newValues
	}

	sort.Slice(values, func(i, j int) bool {

		return values[i] < values[j]
	})

	newValues = values[1 : len(values)-1]

	return newValues
}

func (algorithm CCHWindows) setFilledMetadata(index int) {

	algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, ProvisionalBalanceMeasure)
	algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralEstimated)
	if algorithm.B.isPointType([]measures.PointType{"1", "2"}) {
		algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 9)
		algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, EstimateHistoricMainMeasurePoint)
		return
	}

	algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 5)
	algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, EstimateByHistoricLastYear)

}

type CCHWindowsCurve struct {
	BillingLoadCurve
	criteria int
}

func (algorithm CCHWindows) generateCriterias(period measures.PeriodKey, seasonId string, dayTypeId string) (string, string, string) {
	criteria := fmt.Sprintf("%s%s%s", period, seasonId, dayTypeId)
	criteria2 := fmt.Sprintf("%s%s", period, seasonId)
	criteria3 := fmt.Sprintf("%s", period)

	return criteria, criteria2, criteria3
}

func (algorithm CCHWindows) groupByCriteria(historyLoadCurves []BillingLoadCurve) map[string][]CCHWindowsCurve {
	curveMap := make(map[string][]CCHWindowsCurve)
	for _, curve := range historyLoadCurves {
		criteria, criteria2, criteria3 := algorithm.generateCriterias(curve.Period, curve.SeasonId, curve.DayTypeId)

		curveWithCriteria := CCHWindowsCurve{
			BillingLoadCurve: curve,
			criteria:         2,
		}

		curveMap[criteria] = append(curveMap[criteria], curveWithCriteria)

		curveWithCriteria.criteria = 1
		curveMap[criteria2] = append(curveMap[criteria2], curveWithCriteria)

		curveWithCriteria.criteria = 0
		curveMap[criteria3] = append(curveMap[criteria3], curveWithCriteria)
	}

	return curveMap
}

func (algorithm CCHWindows) fillCurvesToEstimate(curveIndex []int, curvesEstimateValues map[int][]float64, curveMap map[string][]CCHWindowsCurve) {
	for _, index := range curveIndex {
		curve := algorithm.B.BillingLoadCurve[index]
		criteria, criteria2, criteria3 := algorithm.generateCriterias(curve.Period, curve.SeasonId, curve.DayTypeId)

		curvesToEstimate := curveMap[criteria]
		if curves := curveMap[criteria2]; len(curvesToEstimate) < 6 {
			curvesToEstimate = append(curvesToEstimate, curves...)
		}

		if curves := curveMap[criteria3]; len(curvesToEstimate) < 6 {
			curvesToEstimate = append(curvesToEstimate, curves...)
		}

		sort.Slice(curvesToEstimate, func(i, j int) bool {
			bestCriteria := curvesToEstimate[i].criteria > curvesToEstimate[j].criteria
			return bestCriteria
		})

		sort.Slice(curvesToEstimate, func(i, j int) bool {
			dateI := math.Abs(curvesToEstimate[i].EndDate.Sub(curve.EndDate).Hours())
			dateJ := math.Abs(curvesToEstimate[j].EndDate.Sub(curve.EndDate).Hours())
			return dateI < dateJ
		})

		curvesEstimateValues[index] = utils.MapSlice(curvesToEstimate[:6], func(item CCHWindowsCurve) float64 {
			return item.GetMagnitude(algorithm.Magnitude)
		})
	}

}

func (algorithm CCHWindows) calcEstimateCurve(curvesEstimateValues map[int][]float64) {
	for i, values := range curvesEstimateValues {

		//Nos Quedamos con la medida util
		utilValues := algorithm.getUtilValues(values...)

		//CALCULAMOS LA MEDIA DESVIACION TIPICA Y RANGOS MAX Y MIN (SE CALCULAN A PARTIR DE LOS VALORES UTILES)
		average := calcAverage(utilValues...)
		typicalDeviation := calcTypicalDeviation(average, utilValues...)
		maxValue, minValue := getValidRangeValues(average, typicalDeviation)

		//Comprobamos los valores en rangos validos
		validValues := getIsInValidRange(maxValue, minValue, values...)

		//Calculamos la media de los valores validos para obtener la estimacion final
		finalAverage := math.Round(calcAverage(validValues...))

		algorithm.B.SetLoadCurvePeriodMagnitude(i, algorithm.Magnitude, finalAverage)
		algorithm.setFilledMetadata(i)
	}
}

func (algorithm CCHWindows) Execute(ctx context.Context) error {
	curveIndex := make([]int, 0, 744)
	periodCurves := make(map[int][]float64)
	historyBillingCurve := make([]BillingLoadCurve, 0, len(algorithm.GraphContext.IterativeHistory)+740)
	historyBillingCurve = utils.MapSlice(algorithm.GraphContext.IterativeHistory, func(item process_measures.ProcessedLoadCurve) BillingLoadCurve {
		return BillingLoadCurve{
			EndDate:          item.EndDate,
			Origin:           item.Origin,
			Equipment:        "", // FALTA POR METER A LAS PROCESADAS
			Period:           item.Period,
			SeasonId:         item.SeasonId,
			DayTypeId:        item.DayTypeId,
			MeasurePointType: measures.MeasurePointType(item.MeasurePointType),
			AI:               item.AI,
			AE:               item.AE,
			R1:               item.R1,
			R2:               item.R2,
			R3:               item.R3,
			R4:               item.R4,
		}
	})

	/*
		PARA TODOS LOS CRITERIOS TIENE QUE SER EL MISMO PERIODO
		Criterio 1 Mismo DayType, Misma Season, Mismo Mes
		Criterio 2 Misma Season, Mismo DayType
		Criterio 3 Cualquier cosa
	*/

	for i, curve := range algorithm.B.BillingLoadCurve {
		if curve.Origin != measures.Filled {
			algorithm.setValidMetadata(i)
			historyBillingCurve = append(historyBillingCurve, curve)
			continue
		}
		periodCurves[i] = make([]float64, 0, 6)
		curveIndex = append(curveIndex, i)

	}

	//Group Curve By Criteria
	curveMap := algorithm.groupByCriteria(historyBillingCurve)

	//Fill Curves to estimate with best values
	algorithm.fillCurvesToEstimate(curveIndex, periodCurves, curveMap)

	//Calc Estimate value
	algorithm.calcEstimateCurve(periodCurves)

	return nil
}

type CCHWindowsBalanceModulated struct {
	Name      string
	Magnitude measures.Magnitude
	B         *BillingMeasure
	context   *GraphContext
}

func NewCCHWindowsBalanceModulated(b *BillingMeasure, magnitude measures.Magnitude, context *GraphContext) *CCHWindowsBalanceModulated {
	return &CCHWindowsBalanceModulated{
		Name:      "CCH_WINDOWS_BALANCE_MODULATED",
		Magnitude: magnitude,
		B:         b,
		context:   context,
	}
}

func (algorithm CCHWindowsBalanceModulated) ID() string {
	return algorithm.Name
}

func (algorithm CCHWindowsBalanceModulated) getCurveSum() (map[measures.PeriodKey]float64, map[measures.PeriodKey]float64) {
	curveRealSum := make(map[measures.PeriodKey]float64)
	curveEstimatedSum := make(map[measures.PeriodKey]float64)

	for i, curve := range algorithm.B.BillingLoadCurve {
		if curve.Origin != measures.Filled {
			curveRealSum[curve.Period] += algorithm.B.GetLoadCurvePeriodMagnitude(i, algorithm.Magnitude)
			continue
		}

		curveEstimatedSum[curve.Period] += algorithm.B.GetLoadCurvePeriodMagnitude(i, algorithm.Magnitude)
	}

	for _, period := range algorithm.B.Periods {
		curveSum, ok := curveRealSum[period]

		if ok && curveSum == 0 {
			curveRealSum[period] = 1
		}

		curveSum, ok = curveEstimatedSum[period]

		if ok && curveSum == 0 {
			curveEstimatedSum[period] = 1
		}

	}
	return curveRealSum, curveEstimatedSum
}

func (algorithm CCHWindowsBalanceModulated) setFilledMetadata(index int) {
	if !algorithm.B.isPointType([]measures.PointType{"1", "2"}) {
		return
	}

	algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralEstimated)
	algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 7)
	algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, EstimateHistoricMainMeasurePointBalance)
	algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, ProvisionalBalanceMeasure)
}

func (algorithm CCHWindowsBalanceModulated) Execute(ctx context.Context) error {
	cchWindows := NewCCHWindows(algorithm.B, algorithm.Magnitude, algorithm.context)
	err := cchWindows.Execute(ctx)

	if err != nil {
		return err
	}

	curveRealSum, curveEstimateSum := algorithm.getCurveSum()
	balance := algorithm.B.GetBalancePeriodMagnitude(measures.P0, algorithm.Magnitude)

	for i, curve := range algorithm.B.BillingLoadCurve {
		if curve.Origin != measures.Filled {
			continue
		}

		curveMagnitudeValue := algorithm.B.GetLoadCurvePeriodMagnitude(i, algorithm.Magnitude)
		estimateValue := math.Round(((balance - curveRealSum[curve.Period]) * curveMagnitudeValue) / curveEstimateSum[curve.Period])
		algorithm.B.SetLoadCurvePeriodMagnitude(i, algorithm.Magnitude, estimateValue)
		algorithm.setFilledMetadata(i)
	}

	return nil
}

type CCHWindowsCloseModulated struct {
	Name      string
	Magnitude measures.Magnitude
	B         *BillingMeasure
	context   *GraphContext
}

func NewCCHWindowsCloseModulated(b *BillingMeasure, magnitude measures.Magnitude, context *GraphContext) *CCHWindowsCloseModulated {
	return &CCHWindowsCloseModulated{
		Name:      "CCH_WINDOWS_CLOSE_MODULATED",
		Magnitude: magnitude,
		B:         b,
		context:   context,
	}
}

func (algorithm CCHWindowsCloseModulated) ID() string {
	return algorithm.Name
}

func (algorithm CCHWindowsCloseModulated) getCurveSum() (map[measures.PeriodKey]float64, map[measures.PeriodKey]float64) {
	curveRealSum := make(map[measures.PeriodKey]float64)
	curveEstimatedSum := make(map[measures.PeriodKey]float64)

	for i, curve := range algorithm.B.BillingLoadCurve {
		if curve.Origin != measures.Filled {
			curveRealSum[curve.Period] += algorithm.B.GetLoadCurvePeriodMagnitude(i, algorithm.Magnitude)
			continue
		}

		curveEstimatedSum[curve.Period] += algorithm.B.GetLoadCurvePeriodMagnitude(i, algorithm.Magnitude)
	}

	for _, period := range algorithm.B.Periods {

		curveSum, ok := curveRealSum[period]

		if ok && curveSum == 0 {
			curveRealSum[period] = 1
		}

		curveSum, ok = curveEstimatedSum[period]

		if ok && curveSum == 0 {
			curveEstimatedSum[period] = 1
		}
	}

	return curveRealSum, curveEstimatedSum
}

func (algorithm CCHWindowsCloseModulated) setFilledMetadata(index int) {
	if !algorithm.B.isPointType([]measures.PointType{"1", "2"}) {
		return
	}

	algorithm.B.SetLoadCurveGeneralEstimatedMethod(index, algorithm.Magnitude, GeneralEstimated)
	algorithm.B.SetLoadCurveEstimatedCode(index, algorithm.Magnitude, 7)
	algorithm.B.SetLoadCurveEstimatedMethod(index, algorithm.Magnitude, EstimateHistoricMainMeasurePointBalance)
	algorithm.B.SetLoadCurveMeasureType(index, algorithm.Magnitude, ProvisionalBalanceMeasure)
}

func (algorithm CCHWindowsCloseModulated) Execute(ctx context.Context) error {
	cchWindows := NewCCHWindows(algorithm.B, algorithm.Magnitude, algorithm.context)
	err := cchWindows.Execute(ctx)

	if err != nil {
		return err
	}

	curveRealSum, curveEstimateSum := algorithm.getCurveSum()

	for i, curve := range algorithm.B.BillingLoadCurve {
		closeAtr := algorithm.B.GetBalancePeriodMagnitude(curve.Period, algorithm.Magnitude)
		if curve.Origin != measures.Filled {
			continue
		}

		curveMagnitudeValue := algorithm.B.GetLoadCurvePeriodMagnitude(i, algorithm.Magnitude)
		estimateValue := math.Round(((closeAtr - curveRealSum[curve.Period]) * curveMagnitudeValue) / curveEstimateSum[curve.Period])
		algorithm.B.SetLoadCurvePeriodMagnitude(i, algorithm.Magnitude, estimateValue)
		algorithm.setFilledMetadata(i)
	}

	return nil
}

type CCHPenalty struct {
	Name      string
	Magnitude measures.Magnitude
	B         *BillingMeasure
	context   context.Context
	Period    measures.PeriodKey
}

func NewCCHPenalty(b *BillingMeasure, magnitude measures.Magnitude) *CCHPenalty {
	return &CCHPenalty{
		Name:      "CCH_PENALTY",
		Magnitude: magnitude,
		B:         b,
	}
}

func (algorithm CCHPenalty) ID() string {
	return algorithm.Name
}

func (algorithm CCHPenalty) Execute(_ context.Context) error {
	for i, _ := range algorithm.B.BillingLoadCurve {
		if algorithm.B.isPointType([]measures.PointType{"3", "4", "5"}) {
			if algorithm.B.IsLoadCurveOrigin(i, []measures.OriginType{measures.STM, measures.TPL, measures.File, measures.Manual, measures.Visual}) {
				algorithm.B.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 1)
				algorithm.B.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, RealValidMeasure)
				algorithm.B.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralReal)
				algorithm.B.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)
				continue
			}
			if algorithm.B.IsLoadCurveOrigin(i, []measures.OriginType{measures.Filled}) {
				algorithm.B.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 5)
				algorithm.B.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, EstimateByHistoricLastYear)
				algorithm.B.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralEstimated)
				algorithm.B.SetLoadCurveMeasureType(i, algorithm.Magnitude, ProvisionalBalanceMeasure)
				continue
			}
		}

		if algorithm.B.isPointType([]measures.PointType{"1", "2"}) {
			if algorithm.B.IsLoadCurveOrigin(i, []measures.OriginType{measures.STM, measures.TPL, measures.File, measures.Manual, measures.Visual}) {
				if algorithm.B.GetLoadCurveEquipment(i) == measures.Main {
					algorithm.B.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 1)
					algorithm.B.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, FirmMainConfig)
				}
				if algorithm.B.GetLoadCurveEquipment(i) == measures.Redundant {
					algorithm.B.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 2)
					algorithm.B.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, FirmRedundantConfig)
				}
				if algorithm.B.GetLoadCurveEquipment(i) == measures.Receipt {
					algorithm.B.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 3)
					algorithm.B.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, FirmReceiptConfig)
				}
				algorithm.B.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralReal)
				algorithm.B.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)
				continue
			}

			if algorithm.B.IsLoadCurveOrigin(i, []measures.OriginType{measures.Filled}) {
				algorithm.B.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 22)
				algorithm.B.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, EstimateWhosePenaltiesForClientsTypeOneAndTwo)
				algorithm.B.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralEstimated)
				algorithm.B.SetLoadCurveMeasureType(i, algorithm.Magnitude, ProvisionalBalanceMeasure)
				continue
			}
		}
	}
	return nil
}

type CCCHCompleteGD struct {
	Name                       string
	B                          *BillingMeasure                             `bson:"-"`
	processedMeasureRepository process_measures.ProcessedMeasureRepository `bson:"-"`
	Magnitude                  measures.Magnitude
}

func NewCCCHCompleteGD(b *BillingMeasure, magnitude measures.Magnitude, processedMeasureRepository process_measures.ProcessedMeasureRepository) *CCCHCompleteGD {
	return &CCCHCompleteGD{
		Name:                       "CCH_CCCH_COMPLETE_GD",
		B:                          b,
		processedMeasureRepository: processedMeasureRepository,
		Magnitude:                  magnitude,
	}
}

func (algorithm CCCHCompleteGD) SumQuarterMagnitudes(billingLoadCurve *BillingLoadCurve, curve *process_measures.ProcessedLoadCurve) {
	newValue := billingLoadCurve.GetMagnitude(algorithm.Magnitude) + curve.GetMagnitude(algorithm.Magnitude)
	billingLoadCurve.SetMagnitude(algorithm.Magnitude, newValue)
}

func (algorithm CCCHCompleteGD) SetMetadata(i int) {

	if algorithm.B.isPointType([]measures.PointType{"3", "4", "5"}) {
		if algorithm.B.IsLoadCurveOrigin(i, []measures.OriginType{measures.STM, measures.TPL, measures.File, measures.Manual, measures.Visual, measures.Filled}) {
			algorithm.B.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 1)
			algorithm.B.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, RealValidMeasure)
			algorithm.B.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralReal)
			algorithm.B.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)
		}
	}

	if algorithm.B.isPointType([]measures.PointType{"1", "2"}) {
		if algorithm.B.IsLoadCurveOrigin(i, []measures.OriginType{measures.STM, measures.TPL, measures.File, measures.Manual, measures.Visual}) {
			if algorithm.B.GetLoadCurveEquipment(i) == measures.Main {
				algorithm.B.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 1)
				algorithm.B.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, FirmMainConfig)
			}
			if algorithm.B.GetLoadCurveEquipment(i) == measures.Redundant {
				algorithm.B.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 2)
				algorithm.B.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, FirmRedundantConfig)
			}
			if algorithm.B.GetLoadCurveEquipment(i) == measures.Receipt {
				algorithm.B.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 3)
				algorithm.B.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, FirmReceiptConfig)
			}
			algorithm.B.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralReal)
			algorithm.B.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)

		}
		if algorithm.B.IsLoadCurveOrigin(i, []measures.OriginType{measures.Filled}) {
			algorithm.B.SetLoadCurveEstimatedCode(i, algorithm.Magnitude, 2)
			algorithm.B.SetLoadCurveEstimatedMethod(i, algorithm.Magnitude, CalculatedBalance)
			algorithm.B.SetLoadCurveGeneralEstimatedMethod(i, algorithm.Magnitude, GeneralCalculated)
			algorithm.B.SetLoadCurveMeasureType(i, algorithm.Magnitude, FirmBalanceMeasure)

		}
	}
}

func (algorithm CCCHCompleteGD) ID() string {
	return algorithm.Name
}

func (algorithm CCCHCompleteGD) Execute(ctx context.Context) error {

	billingLoadCurvesFilled := map[string]struct {
		b *BillingLoadCurve
		i int
	}{}

	loadCurves, err := algorithm.processedMeasureRepository.ListHistoryLoadCurve(ctx, process_measures.QueryHistoryLoadCurve{
		CUPS:          algorithm.B.CUPS,
		StartDate:     algorithm.B.InitDate.AddDate(0, 0, -1),
		EndDate:       algorithm.B.EndDate.AddDate(0, 0, 1),
		Periods:       algorithm.B.Periods,
		WithCriterias: true,
		Type:          measures.QuarterHour,
		Magnitude:     algorithm.Magnitude,
	})

	if err != nil {
		return err
	}

	for index, loadCurveBilling := range algorithm.B.BillingLoadCurve {
		algorithm.SetMetadata(index)
		if loadCurveBilling.Origin == measures.Filled {
			billingLoadCurvesFilled[loadCurveBilling.EndDate.Format("2006-01-02T15")] = struct {
				b *BillingLoadCurve
				i int
			}{b: &algorithm.B.BillingLoadCurve[index], i: index}
		}
	}

	for _, loadCurve := range loadCurves {
		currentDate := loadCurve.EndDate.Format("2006-01-02T15")

		if _, curveFound := billingLoadCurvesFilled[currentDate]; !curveFound {
			continue
		}

		algorithm.SumQuarterMagnitudes(billingLoadCurvesFilled[currentDate].b, &loadCurve)

	}

	return nil
}
