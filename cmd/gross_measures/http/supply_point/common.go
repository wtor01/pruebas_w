package supply_point

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
)

func toGrossMeasureSupplyPoint(dashboardMeasure gross_measures.DashboardMeasureSupplyPoint) GrossMeasureServicePoint {
	return GrossMeasureServicePoint{
		CalendarCurve:          utils.MapSlice(dashboardMeasure.CalendarCurve, toGrossCalendarStatus),
		CalendarDailyClosure:   utils.MapSlice(dashboardMeasure.CalendarDailyClosure, toGrossCalendarStatus),
		CalendarMonthlyClosure: utils.MapSlice(dashboardMeasure.CalendarMonthlyClosure, toGrossCalendarStatus),
		ListDailyClosures:      utils.MapSlice(dashboardMeasure.ListDailyClosure, toGrossListDaily),
		ListMonthlyClosures:    utils.MapSlice(dashboardMeasure.ListMonthlyClosure, toGrossListMonthly),
		Cups:                   dashboardMeasure.Cups,
		MagnitudeEnergy:        string(dashboardMeasure.MagnitudeEnergy),
		Magnitudes: utils.MapSlice(dashboardMeasure.Magnitudes, func(item measures.Magnitude) GrossMeasureServicePointMagnitudes {
			return GrossMeasureServicePointMagnitudes(item)
		}),
		Periods: utils.MapSlice(dashboardMeasure.Periods, func(item measures.PeriodKey) GrossMeasureServicePointPeriods {
			return GrossMeasureServicePointPeriods(item)
		}),
		SerialNumber: dashboardMeasure.SerialNumber,
		Type:         GrossMeasureServicePointType(dashboardMeasure.MeterType),
	}
}

func toGrossCalendarStatus(calendar gross_measures.CalendarStatus) ServicePointCalendarStatus {
	return ServicePointCalendarStatus{
		Date:   calendar.Date,
		Status: GrossMeasureValidationStatus(calendar.Status),
	}
}

func toGrossMeasureValues(values measures.Values) GrossMeasureValues {
	return GrossMeasureValues{
		AE: values.AE,
		AI: values.AI,
		R1: values.R1,
		R2: values.R2,
		R3: values.R3,
		R4: values.R4,
	}
}

func toGrossMeasureMonthlyValues(monthlyValues measures.ValuesMonthly) GrossMeasureMonthlyValues {
	return GrossMeasureMonthlyValues{
		GrossMeasureValues: toGrossMeasureValues(monthlyValues.Values),
		AEi:                monthlyValues.AEi,
		AIi:                monthlyValues.AIi,
		R1i:                monthlyValues.R1i,
		R2i:                monthlyValues.R2i,
		R3i:                monthlyValues.R3i,
		R4i:                monthlyValues.R4i,
	}
}

func toGrossListDaily(dailyCloses gross_measures.ListDailyClose) ServicePointDailyValues {
	return ServicePointDailyValues{
		EndDate: dailyCloses.EndDate,
		Status:  GrossMeasureValidationStatus(dailyCloses.Status),
		File:    dailyCloses.Origin,
		P0:      toGrossMeasureValues(dailyCloses.P0),
		P1:      toGrossMeasureValues(dailyCloses.P1),
		P2:      toGrossMeasureValues(dailyCloses.P2),
		P3:      toGrossMeasureValues(dailyCloses.P3),
		P4:      toGrossMeasureValues(dailyCloses.P4),
		P5:      toGrossMeasureValues(dailyCloses.P5),
		P6:      toGrossMeasureValues(dailyCloses.P6),
	}
}

func toGrossListMonthly(monthlyCloses gross_measures.ListMonthlyClose) ServicePointMonthlyValues {
	return ServicePointMonthlyValues{
		InitDate: monthlyCloses.InitDate,
		EndDate:  monthlyCloses.EndDate,
		Status:   GrossMeasureValidationStatus(monthlyCloses.Status),
		File:     monthlyCloses.Origin,
		P0:       toGrossMeasureMonthlyValues(monthlyCloses.P0),
		P1:       toGrossMeasureMonthlyValues(monthlyCloses.P1),
		P2:       toGrossMeasureMonthlyValues(monthlyCloses.P2),
		P3:       toGrossMeasureMonthlyValues(monthlyCloses.P3),
		P4:       toGrossMeasureMonthlyValues(monthlyCloses.P4),
		P5:       toGrossMeasureMonthlyValues(monthlyCloses.P5),
		P6:       toGrossMeasureMonthlyValues(monthlyCloses.P6),
	}
}

func toGrossCurves(dashboardCurves gross_measures.DashboardSupplyPointCurve) CurveGrossMeasureMeter {
	return CurveGrossMeasureMeter{
		ServicePointCalendarStatus: ServicePointCalendarStatus{
			Date:   dashboardCurves.Hour,
			Status: GrossMeasureValidationStatus(dashboardCurves.Status),
		},
		File:   dashboardCurves.File,
		Values: toGrossMeasureValues(dashboardCurves.Values),
	}
}
