package dashboard_stats

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
)

func DashboardMeasuresStatsToResponse(pmds gross_measures.GrossMeasuresDashboardStatsGlobal) GrossMeasureDashboard {
	var measureType GrossMeasureDashboardType

	switch pmds.Type {
	case "OTHER":
		measureType = GrossMeasureDashboardTypeOTHER
	case "TLG":
		measureType = GrossMeasureDashboardTypeTLG
	case "TLM":
		measureType = GrossMeasureDashboardTypeTLM

	}

	dailyStats := utils.MapSlice(pmds.DailyStats, func(item gross_measures.DashboardSingleStats) GrossMeasureDashboardDaily {
		return GrossMeasureDashboardDaily{
			DailyClosure:           item.DailyClosure,
			Day:                    &item.Day,
			ExpectedDailyClosure:   item.ExpectedDailyClosure,
			ExpectedHourlyCurve:    item.ExpectedHourlyCurve,
			ExpectedMonthlyClosure: item.ExpectedMonthlyClosure,
			ExpectedQuarterlyCurve: item.ExpectedQuarterlyCurve,
			HourlyCurve:            item.HourlyCurve,
			MonthlyClosure:         item.MonthlyClosure,
			QuarterlyCurve:         item.QuarterlyCurve,
		}
	})

	return GrossMeasureDashboard{
		GrossMeasuresDailyStats: &dailyStats,
		DistributorId:           pmds.DistributorID,
		Month:                   pmds.Month,
		Type:                    measureType,
		Year:                    pmds.Year,
	}
}

func GrossMeasureStatisticsByCupsToResponse(dashboardSerialNumberObj gross_measures.DashboardSerialNumber) GrossMeasureDashboardCUPS {
	var measureType GrossMeasureDashboardCUPSType
	switch dashboardSerialNumberObj.Type {
	case string(measures.OTHER):
		measureType = GrossMeasureDashboardCUPSTypeOTHER
	case string(measures.TLG):
		measureType = GrossMeasureDashboardCUPSTypeTLG
	case string(measures.TLM):
		measureType = GrossMeasureDashboardCUPSTypeTLM
	}

	dailyStats := utils.MapSlice(dashboardSerialNumberObj.DailyStats, func(item gross_measures.DashboardSingleStats) GrossMeasureDashboardDaily {
		return GrossMeasureDashboardDaily{
			DailyClosure:           item.DailyClosure,
			Day:                    &item.Day,
			ExpectedDailyClosure:   item.ExpectedDailyClosure,
			ExpectedHourlyCurve:    item.ExpectedHourlyCurve,
			ExpectedMonthlyClosure: item.ExpectedMonthlyClosure,
			ExpectedQuarterlyCurve: item.ExpectedQuarterlyCurve,
			HourlyCurve:            item.HourlyCurve,
			MonthlyClosure:         item.MonthlyClosure,
			QuarterlyCurve:         item.QuarterlyCurve,
		}
	})

	return GrossMeasureDashboardCUPS{
		Cups:                    dashboardSerialNumberObj.Cups,
		DistributorId:           dashboardSerialNumberObj.DistributorId,
		GrossMeasuresDailyStats: &dailyStats,
		Month:                   dashboardSerialNumberObj.Month,
		SerialNumber:            dashboardSerialNumberObj.SerialNumber,
		ServicePointType:        dashboardSerialNumberObj.ServicePointType,
		ServiceType:             dashboardSerialNumberObj.ServiceType,
		Type:                    measureType,
		Year:                    dashboardSerialNumberObj.Year,
	}
}
