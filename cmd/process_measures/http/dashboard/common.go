package dashboard

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
)

func (c *DailyClosure) setPeriodDashboardServicePoint(period measures.PeriodKey, values *measures.Values) {

	if values == nil {
		return
	}
	v := CurveValues{
		AE: values.AE,
		AI: values.AI,
		R1: values.R1,
		R2: values.R2,
		R3: values.R3,
		R4: values.R4,
	}

	switch period {
	case measures.P0:
		c.P0 = &v
	case measures.P1:
		c.P1 = &v
	case measures.P2:
		c.P2 = &v
	case measures.P3:
		c.P3 = &v
	case measures.P4:
		c.P4 = &v
	case measures.P5:
		c.P5 = &v
	case measures.P6:
		c.P6 = &v
	}
}
func (c *MonthlyClosure) setPeriodDashboardServicePoint(period measures.PeriodKey, values *measures.ValuesMonthly) {
	if values == nil {
		return
	}
	v := MonthlyValues{
		AE:  values.AE,
		AI:  values.AI,
		R1:  values.R1,
		R2:  values.R2,
		R3:  values.R3,
		R4:  values.R4,
		AEi: values.AEi,
		AIi: values.AIi,
		R1i: values.R1i,
		R2i: values.R2i,
		R3i: values.R3i,
		R4i: values.R4i,
	}

	switch period {
	case measures.P0:
		c.P0 = &v
	case measures.P1:
		c.P1 = &v
	case measures.P2:
		c.P2 = &v
	case measures.P3:
		c.P3 = &v
	case measures.P4:
		c.P4 = &v
	case measures.P5:
		c.P5 = &v
	case measures.P6:
		c.P6 = &v
	}
}
func (c *Curve) setPeriodDashboardServicePoint(period measures.PeriodKey, values *measures.Values) {

	if values == nil {
		return
	}
	v := CurveValues{
		AE: values.AE,
		AI: values.AI,
		R1: values.R1,
		R2: values.R2,
		R3: values.R3,
		R4: values.R4,
	}

	switch period {
	case measures.P1:
		c.P1 = &v
	case measures.P2:
		c.P2 = &v
	case measures.P3:
		c.P3 = &v
	case measures.P4:
		c.P4 = &v
	case measures.P5:
		c.P5 = &v
	case measures.P6:
		c.P6 = &v
	}
}

func GetDashboardProcessServicePointToResponse(data process_measures.ServicePointDashboardWithType) ServicePointProcessDashboard {
	ms := utils.MapSlice(data.ServicePoints, func(item process_measures.ServicePointDashboard) ServicePointDashboard {
		d := ServicePointDashboard{
			Curve:           Curve{Status: string(item.Curve.Status)},
			DailyClosure:    DailyClosure{Status: string(item.DailyClosure.Status)},
			MonthlyClosure:  MonthlyClosure{Status: string(item.MonthlyClosure.Status)},
			Date:            item.Date,
			MagnitudeEnergy: string(item.MagnitudeEnergy),
			Magnitudes: utils.MapSlice(item.Magnitudes, func(m measures.Magnitude) ServicePointDashboardMagnitudes {
				return ServicePointDashboardMagnitudes(m)
			}),
			Periods: utils.MapSlice(item.Periods, func(m measures.PeriodKey) ServicePointDashboardPeriods {
				return ServicePointDashboardPeriods(m)
			}),
		}

		if d.MonthlyClosure.Status != "" {
			d.MonthlyClosure.Id = item.MonthlyClosure.ID
			d.MonthlyClosure.InitDate = item.MonthlyClosure.InitDate
			d.MonthlyClosure.EndDate = item.MonthlyClosure.EndDate
		}

		for _, period := range append(item.Periods, measures.P0) {
			d.Curve.setPeriodDashboardServicePoint(period, item.Curve.GetPeriodMeasure(period))
			d.DailyClosure.setPeriodDashboardServicePoint(period, item.DailyClosure.GetPeriodMeasure(period))
			d.MonthlyClosure.setPeriodDashboardServicePoint(period, item.MonthlyClosure.GetPeriodMeasure(period))
		}

		return d

	})

	sppd := ServicePointProcessDashboard{}
	days := append(ServicePointDashboardResponseDays{}, ms...)
	sppd.Days = &days
	meterType := ServicePointDashboardResponseType(data.Type)
	sppd.Type = &meterType
	return sppd
}

func GetProcessMeasureDashboardListToResponse(data measures.DashboardListCups) []DashboardCups {
	cupsList := make([]DashboardCups, 0, cap(data.Cups))

	for _, cup := range data.Cups {
		cupsList = append(cupsList, DashboardCups{
			Cups: cup.Cups,
			Curve: DashboardCupsValues{
				Invalid:   cup.Values.Curve.Invalid,
				None:      cup.Values.Curve.None,
				ShouldBe:  cup.Values.Curve.ShouldBe,
				Supervise: cup.Values.Curve.Supervise,
				Total:     cup.Values.Curve.Total,
				Valid:     cup.Values.Curve.Valid,
			},
			Daily: DashboardCupsValues{
				Invalid:   cup.Values.Daily.Invalid,
				None:      cup.Values.Daily.None,
				ShouldBe:  cup.Values.Daily.ShouldBe,
				Supervise: cup.Values.Daily.Supervise,
				Total:     cup.Values.Daily.Total,
				Valid:     cup.Values.Daily.Valid,
			},
			Monthly: DashboardCupsValues{
				Invalid:   cup.Values.Monthly.Invalid,
				None:      cup.Values.Monthly.None,
				ShouldBe:  cup.Values.Monthly.ShouldBe,
				Supervise: cup.Values.Monthly.Supervise,
				Total:     cup.Values.Monthly.Total,
				Valid:     cup.Values.Monthly.Valid,
			},
		})
	}

	return cupsList
}
