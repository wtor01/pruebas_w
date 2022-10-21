package dashboard

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"strings"
)

func fiscalBillingMeasuresToResponse(fBillingMeasuresDashboard billing_measures.FiscalBillingMeasuresDashboard) FiscalBillingMeasures {

	return FiscalBillingMeasures{
		Balance: BalancePeriods{
			Method: BalancePeriodsMethod(fBillingMeasuresDashboard.Balance.Method),
			Origin: BalancePeriodsOrigin(fBillingMeasuresDashboard.Balance.Origin),
			P0:     toFiscalBalance(fBillingMeasuresDashboard.Balance.P0),
			P1:     toFiscalBalance(fBillingMeasuresDashboard.Balance.P1),
			P2:     toFiscalBalance(fBillingMeasuresDashboard.Balance.P2),
			P3:     toFiscalBalance(fBillingMeasuresDashboard.Balance.P3),
			P4:     toFiscalBalance(fBillingMeasuresDashboard.Balance.P4),
			P5:     toFiscalBalance(fBillingMeasuresDashboard.Balance.P5),
			P6:     toFiscalBalance(fBillingMeasuresDashboard.Balance.P6),
		},
		MeterType:          string(fBillingMeasuresDashboard.Type),
		CalendarCurve:      toFiscalCalendarCurve(fBillingMeasuresDashboard.CalendarCurve),
		Cups:               fBillingMeasuresDashboard.Cups,
		Curve:              toFiscalCurve(fBillingMeasuresDashboard.Curve),
		EndDate:            fBillingMeasuresDashboard.EndDate,
		LastMvDate:         fBillingMeasuresDashboard.LastMvDate,
		Magnitudes:         toFiscalMagnitudes(fBillingMeasuresDashboard.Magnitudes),
		Periods:            toFiscalPeriods(fBillingMeasuresDashboard.Periods),
		PrincipalMagnitude: FiscalBillingMeasuresPrincipalMagnitude(fBillingMeasuresDashboard.PrincipalMagnitude),
		StartDate:          fBillingMeasuresDashboard.StartDate,
		GraphHistory:       toFiscalGraph(fBillingMeasuresDashboard.GraphHistory),
		Status:             FiscalBillingMeasuresStatus(fBillingMeasuresDashboard.Status),
		Summary:            toFiscalSummary(fBillingMeasuresDashboard.Summary),
		ExecutionSummary:   toFiscalExecutionSummary(fBillingMeasuresDashboard.ExecutionSummary),
		Id:                 fBillingMeasuresDashboard.Id,
	}
}
func toFiscalSummary(fiscalBillingMeasuresSummary billing_measures.Summary) SummaryStructure {
	toResponse := SummaryStructure{
		Adjusted: ItemSummaryStructure{
			P1:    fiscalBillingMeasuresSummary.Adjusted.P1,
			P2:    fiscalBillingMeasuresSummary.Adjusted.P2,
			P3:    fiscalBillingMeasuresSummary.Adjusted.P3,
			P4:    fiscalBillingMeasuresSummary.Adjusted.P4,
			P5:    fiscalBillingMeasuresSummary.Adjusted.P5,
			P6:    fiscalBillingMeasuresSummary.Adjusted.P6,
			Total: fiscalBillingMeasuresSummary.Adjusted.Total,
		},
		Calculated: ItemSummaryStructure{
			P1:    fiscalBillingMeasuresSummary.Calculated.P1,
			P2:    fiscalBillingMeasuresSummary.Calculated.P2,
			P3:    fiscalBillingMeasuresSummary.Calculated.P3,
			P4:    fiscalBillingMeasuresSummary.Calculated.P4,
			P5:    fiscalBillingMeasuresSummary.Calculated.P5,
			P6:    fiscalBillingMeasuresSummary.Calculated.P6,
			Total: fiscalBillingMeasuresSummary.Calculated.Total,
		},
		Consum: struct {
			P1    float64 `json:"p1"`
			P2    float64 `json:"p2"`
			P3    float64 `json:"p3"`
			P4    float64 `json:"p4"`
			P5    float64 `json:"p5"`
			P6    float64 `json:"p6"`
			Total float64 `json:"total"`
		}{
			P1:    fiscalBillingMeasuresSummary.Consum.P1,
			P2:    fiscalBillingMeasuresSummary.Consum.P2,
			P3:    fiscalBillingMeasuresSummary.Consum.P3,
			P4:    fiscalBillingMeasuresSummary.Consum.P4,
			P5:    fiscalBillingMeasuresSummary.Consum.P5,
			P6:    fiscalBillingMeasuresSummary.Consum.P6,
			Total: fiscalBillingMeasuresSummary.Consum.Total,
		},
		Estimated: ItemSummaryStructure{
			P1:    fiscalBillingMeasuresSummary.Estimated.P1,
			P2:    fiscalBillingMeasuresSummary.Estimated.P2,
			P3:    fiscalBillingMeasuresSummary.Estimated.P3,
			P4:    fiscalBillingMeasuresSummary.Estimated.P4,
			P5:    fiscalBillingMeasuresSummary.Estimated.P5,
			P6:    fiscalBillingMeasuresSummary.Estimated.P6,
			Total: fiscalBillingMeasuresSummary.Estimated.Total,
		},
		Outlined: ItemSummaryStructure{
			P1:    fiscalBillingMeasuresSummary.Outlined.P1,
			P2:    fiscalBillingMeasuresSummary.Outlined.P2,
			P3:    fiscalBillingMeasuresSummary.Outlined.P3,
			P4:    fiscalBillingMeasuresSummary.Outlined.P4,
			P5:    fiscalBillingMeasuresSummary.Outlined.P5,
			P6:    fiscalBillingMeasuresSummary.Outlined.P6,
			Total: fiscalBillingMeasuresSummary.Outlined.Total,
		},
		Real: ItemSummaryStructure{
			P1:    fiscalBillingMeasuresSummary.Real.P1,
			P2:    fiscalBillingMeasuresSummary.Real.P2,
			P3:    fiscalBillingMeasuresSummary.Real.P3,
			P4:    fiscalBillingMeasuresSummary.Real.P4,
			P5:    fiscalBillingMeasuresSummary.Real.P5,
			P6:    fiscalBillingMeasuresSummary.Real.P6,
			Total: fiscalBillingMeasuresSummary.Real.Total,
		},
	}
	return toResponse
}

func toFiscalPeriods(fiscalBillingMeasuresPeriods []measures.PeriodKey) []FiscalBillingMeasuresPeriods {
	toResponse := make([]FiscalBillingMeasuresPeriods, 0)
	for _, magnitude := range fiscalBillingMeasuresPeriods {
		toResponse = append(toResponse, FiscalBillingMeasuresPeriods(magnitude))
	}
	return toResponse
}

func toFiscalMagnitudes(fiscalBillingMeasuresMagnitudes []measures.Magnitude) []FiscalBillingMeasuresMagnitudes {
	toResponse := make([]FiscalBillingMeasuresMagnitudes, 0)
	for _, magnitude := range fiscalBillingMeasuresMagnitudes {
		toResponse = append(toResponse, FiscalBillingMeasuresMagnitudes(magnitude))
	}
	return toResponse
}

func toFiscalCurve(fiscalBillingMeasuresCurves []billing_measures.Curve) Curve {
	toResponse := make([]ItemCurve, 0)
	for _, curve := range fiscalBillingMeasuresCurves {
		curveLine := make([]DailyCurve, 0)
		for _, curveStats := range curve.Values {
			curveLine = append(curveLine, DailyCurve{
				Ae:     curveStats.AE,
				AeAuto: curveStats.AeAuto,
				Ai:     curveStats.AI,
				AiAuto: curveStats.AiAuto,
				Date:   curveStats.Hour,
				R1:     curveStats.R1,
				R2:     curveStats.R2,
				R3:     curveStats.R3,
				R4:     curveStats.R4,
				Status: DailyCurveStatus(curveStats.Status),
			})
		}
		toResponse = append(toResponse, ItemCurve{
			Date:   curve.Date,
			Values: curveLine,
		})
	}
	return toResponse
}

func toFiscalCalendarCurve(fiscalBillingMeasuresCalendarCurves []billing_measures.CalendarCurve) CalendarCurve {
	toResponse := make([]CalendarCurveItem, 0)
	for _, calCurve := range fiscalBillingMeasuresCalendarCurves {
		toResponse = append(toResponse, CalendarCurveItem{
			Date:   calCurve.Date,
			Status: CalendarCurveItemStatus(calCurve.Status),
		})
	}
	return toResponse
}

func toFiscalBalance(fiscalBillingMeasuresBalance *billing_measures.BalancePeriod) *PeriodFeatures {
	if fiscalBillingMeasuresBalance == nil {
		return &PeriodFeatures{}
	}
	return &PeriodFeatures{

		Ae:            fiscalBillingMeasuresBalance.AE,
		Ai:            fiscalBillingMeasuresBalance.AI,
		BalanceTypeAe: string(fiscalBillingMeasuresBalance.BalanceTypeAE),
		BalanceTypeAi: string(fiscalBillingMeasuresBalance.BalanceTypeAI),
		BalanceTypeR1: string(fiscalBillingMeasuresBalance.BalanceTypeR1),
		BalanceTypeR2: string(fiscalBillingMeasuresBalance.BalanceTypeR2),
		BalanceTypeR3: string(fiscalBillingMeasuresBalance.BalanceTypeR3),
		BalanceTypeR4: string(fiscalBillingMeasuresBalance.BalanceTypeR4),
		R1:            fiscalBillingMeasuresBalance.R1,
		R2:            fiscalBillingMeasuresBalance.R2,
		R3:            fiscalBillingMeasuresBalance.R3,
		R4:            fiscalBillingMeasuresBalance.R4,
	}
}

func toFiscalGraph(billingGraphHistory map[string]*billing_measures.Graph) []FiscalMeasureGraph {
	fiscalGraph := make([]FiscalMeasureGraph, 0, len(billingGraphHistory))
	for key, graph := range billingGraphHistory {
		keys := strings.Split(key, "_")
		period := string(measures.P0)
		magnitude := keys[0]
		if len(keys) > 1 {
			period = keys[0]
			magnitude = keys[1]
		}

		graphInfo := FiscalMeasureGraph{
			FinishedAt: graph.FinishedAt,
			Period:     FiscalMeasureGraphPeriod(period),
			Magnitude:  FiscalMeasureGraphMagnitude(magnitude),
			StartedAt:  graph.StartedAt,
		}

		cchAlgorithm, _ := utils.FindSlice(graph.Algorithms, func(item string) bool {
			return strings.HasPrefix(item, "CCH")
		})

		if cchAlgorithm != nil {
			graphInfo.CchAlgorithm = *cchAlgorithm
		}

		balanceAlgorithm, _ := utils.FindSlice(graph.Algorithms, func(item string) bool {
			return strings.HasPrefix(item, "BALANCE")
		})

		if balanceAlgorithm != nil {
			graphInfo.BalanceAlgorithm = *balanceAlgorithm
		}

		fiscalGraph = append(fiscalGraph, graphInfo)
	}

	return fiscalGraph
}

func toFiscalExecutionSummary(executionSummary billing_measures.ExecutionSummary) ExecutionSummary {
	return ExecutionSummary{
		BalanceOrigin: ExecutionSummaryBalanceOrigin(executionSummary.BalanceOrigin),
		BalanceType:   MethodType(executionSummary.BalanceType),
		CurveStatus:   ExecutionSummaryCurveStatus(executionSummary.CurveStatus),
		CurveType:     MethodType(executionSummary.CurveType),
	}
}

func toFiscalClosureResume(resume billing_measures.BillingMeasureDashboardResumeClosure) BillingMeasuresResume {
	return BillingMeasuresResume{
		ActualReadingClosure:   toReadingClosureResume(resume.ActualReadingClose),
		PreviousReadingClosure: toReadingClosureResume(resume.PreviousReadingClose),
	}
}
func toReadingClosureResume(readingClosure billing_measures.ReadingClosureResume) ReadingClosure {
	return ReadingClosure{
		MeasureDevice: readingClosure.MeterSerialNumber,
		MeasureType:   readingClosure.ClosureType,
		Origin:        readingClosure.Origin,
		StartDate:     readingClosure.InitDate,
		EndDate:       readingClosure.EndDate,
		P0:            toMagnitudePeriodFeatures(readingClosure.P0),
		P1:            toMagnitudePeriodFeatures(readingClosure.P1),
		P2:            toMagnitudePeriodFeatures(readingClosure.P2),
		P3:            toMagnitudePeriodFeatures(readingClosure.P3),
		P4:            toMagnitudePeriodFeatures(readingClosure.P4),
		P5:            toMagnitudePeriodFeatures(readingClosure.P5),
		P6:            toMagnitudePeriodFeatures(readingClosure.P6),
	}
}

func toMagnitudePeriodFeatures(periodByMagnitudes billing_measures.PeriodsByMagnitude) MagnitudePeriodFeatures {
	return MagnitudePeriodFeatures{
		Ae: MagnitudeFeatures{
			Consum:  periodByMagnitudes.AE.Consum,
			Reading: periodByMagnitudes.AE.Reading,
		},
		Ai: MagnitudeFeatures{
			Consum:  periodByMagnitudes.AI.Consum,
			Reading: periodByMagnitudes.AI.Reading,
		},
		R1: MagnitudeFeatures{
			Consum:  periodByMagnitudes.R1.Consum,
			Reading: periodByMagnitudes.R1.Reading,
		},
		R2: MagnitudeFeatures{
			Consum:  periodByMagnitudes.R2.Consum,
			Reading: periodByMagnitudes.R3.Reading,
		},
		R3: MagnitudeFeatures{
			Consum:  periodByMagnitudes.R3.Consum,
			Reading: periodByMagnitudes.R3.Reading,
		},
		R4: MagnitudeFeatures{
			Consum:  periodByMagnitudes.R4.Consum,
			Reading: periodByMagnitudes.R4.Reading,
		},
	}
}

func toFiscalMeasureSummary(summary billing_measures.FiscalMeasureSummary) FiscalMeasureSummary {
	fiscalSummary := FiscalMeasureSummary{
		MeterType: MeterType(summary.MeterType),
		BalanceOrigin: BalanceOriginSummary{
			Daily:     summary.BalanceOrigin.Daily,
			Monthly:   summary.BalanceOrigin.Monthly,
			NoClosure: summary.BalanceOrigin.NoClosure,
			Other:     summary.BalanceOrigin.Other,
		},
		BalanceType: BalanceTypeSummary{
			Calculated: summary.BalanceType.Calculated,
			Estimated:  summary.BalanceType.Estimated,
			Real:       summary.BalanceType.Real,
		},
		CurveStatus: CurveStatusSummary{
			Absent:       summary.CurveStatus.Absent,
			Completed:    summary.CurveStatus.Completed,
			NotCompleted: summary.CurveStatus.NotCompleted,
		},
		CurveType: CurveTypeSummary{
			Adjusted: summary.CurveType.Adjusted,
			Outlined: summary.CurveType.Outlined,
			Real:     summary.CurveType.Real,
		},
		Total: summary.Total,
	}

	if summary.MeterType != measures.TLG {
		fiscalSummary.CurveType.Calculated = &summary.CurveType.Calculated
		fiscalSummary.CurveType.Estimated = &summary.CurveType.Estimated
	}

	return fiscalSummary
}
