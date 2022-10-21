package stats

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
)

func DashboardMeasuresStatsToResponse(pmds process_measures.ProcessMeasureDashboardStatsGlobal) MeasureStatistics {
	var measureType MeasureStatisticsType

	switch pmds.Type {
	case "OTHER":
		measureType = MeasureStatisticsTypeOTHER
	case "TLG":
		measureType = MeasureStatisticsTypeTLG
	case "TLM":
		measureType = MeasureStatisticsTypeTLM

	}

	dailyStats := utils.MapSlice(pmds.DailyStats, func(item process_measures.DailyStats) MeasureStatisticsDaily {
		return MeasureStatisticsDaily{
			DailyClosure: MeasureStatisticsDailyResults{
				Absent:     &item.DailyClosure.Absent,
				Complete:   &item.DailyClosure.Complete,
				Incomplete: &item.DailyClosure.Incomplete,
				Inv:        &item.DailyClosure.Inv,
				Sup:        &item.DailyClosure.Sup,
				Total:      &item.DailyClosure.Total,
				Val:        &item.DailyClosure.Val,
			},
			Date:                   item.Date,
			ExpectedCurve:          item.ExpectedCurve,
			ExpectedDailyClosure:   item.ExpectedDailyClosure,
			ExpectedMonthlyClosure: item.ExpectedMonthlyClosure,
			LoadCurve: MeasureStatisticsDailyResults{
				Absent:     &item.LoadCurve.Absent,
				Complete:   &item.LoadCurve.Complete,
				Incomplete: &item.LoadCurve.Incomplete,
				Inv:        &item.LoadCurve.Inv,
				Sup:        &item.LoadCurve.Sup,
				Total:      &item.LoadCurve.Total,
				Val:        &item.LoadCurve.Val,
			},
			MonthlyClosure: MeasureStatisticsDailyResults{
				Absent:     &item.MonthlyClosure.Absent,
				Complete:   &item.MonthlyClosure.Complete,
				Incomplete: &item.MonthlyClosure.Incomplete,
				Inv:        &item.MonthlyClosure.Inv,
				Sup:        &item.MonthlyClosure.Sup,
				Total:      &item.MonthlyClosure.Total,
				Val:        &item.MonthlyClosure.Val,
			},
		}
	})

	return MeasureStatistics{
		DailyStats:    &dailyStats,
		DistributorId: pmds.DistributorID,
		Month:         pmds.Month,
		Type:          measureType,
		Year:          pmds.Year,
	}
}

func DashboardMeasuresStatsCupsToResponse(pmds process_measures.ProcessMeasureDashboardStatsGlobal) MeasureStatisticsCUPS {

	var measureType MeasureStatisticsCUPSType

	switch pmds.Type {
	case "OTHER":
		measureType = MeasureStatisticsCUPSTypeOTHER
	case "TLG":
		measureType = MeasureStatisticsCUPSTypeTLG
	case "TLM":
		measureType = MeasureStatisticsCUPSTypeTLM

	}

	dailyStats := utils.MapSlice(pmds.DailyStats, func(item process_measures.DailyStats) MeasureStatisticsDaily {
		return MeasureStatisticsDaily{
			DailyClosure: MeasureStatisticsDailyResults{
				Absent:     &item.DailyClosure.Absent,
				Complete:   &item.DailyClosure.Complete,
				Incomplete: &item.DailyClosure.Incomplete,
				Inv:        &item.DailyClosure.Inv,
				Sup:        &item.DailyClosure.Sup,
				Total:      &item.DailyClosure.Total,
				Val:        &item.DailyClosure.Val,
			},
			Date:                   item.Date,
			ExpectedCurve:          item.ExpectedCurve,
			ExpectedDailyClosure:   item.ExpectedDailyClosure,
			ExpectedMonthlyClosure: item.ExpectedMonthlyClosure,
			LoadCurve: MeasureStatisticsDailyResults{
				Absent:     &item.LoadCurve.Absent,
				Complete:   &item.LoadCurve.Complete,
				Incomplete: &item.LoadCurve.Incomplete,
				Inv:        &item.LoadCurve.Inv,
				Sup:        &item.LoadCurve.Sup,
				Total:      &item.LoadCurve.Total,
				Val:        &item.LoadCurve.Val,
			},
			MonthlyClosure: MeasureStatisticsDailyResults{
				Absent:     &item.MonthlyClosure.Absent,
				Complete:   &item.MonthlyClosure.Complete,
				Incomplete: &item.MonthlyClosure.Incomplete,
				Inv:        &item.MonthlyClosure.Inv,
				Sup:        &item.MonthlyClosure.Sup,
				Total:      &item.MonthlyClosure.Total,
				Val:        &item.MonthlyClosure.Val,
			},
		}
	})

	return MeasureStatisticsCUPS{
		Cups:          pmds.Cups,
		DailyStats:    &dailyStats,
		DistributorId: pmds.DistributorID,
		Month:         pmds.Month,
		Type:          measureType,
		Year:          pmds.Year,
	}

}
