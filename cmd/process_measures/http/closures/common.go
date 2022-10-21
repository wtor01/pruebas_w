package closures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"time"
)

// getPeriod obtain Periods struct
func getPeriod(calendar *process_measures.ProcessedMonthlyClosurePeriod, period measures.PeriodKey) []Periods {
	cal := make([]Periods, 0)

	for _, magnitude := range measures.ValidMagnitudes {
		magnitudeValue, magnitudeIncremental := calendar.GetMagnitude(magnitude)
		if magnitudeValue != 0 {
			cal = append(cal, Periods{
				Absolute:    &magnitudeValue,
				Incremental: &magnitudeIncremental,
				Magnitud:    PeriodsMagnitud(magnitude),
				Period:      PeriodsPeriod(period),
			})
		}
	}
	return cal
}

// getPeriod obtain the periods(P0, P1, ...) for processedMonthlyClosure
func getPeriods(closure process_measures.ProcessedMonthlyClosure) []Periods {
	periods := make([]Periods, 0)

	for _, period := range append(measures.ValidPeriodsCurve, measures.P0) {
		periodCalendar := closure.CalendarPeriods.GetPeriodValues(period)
		if periodCalendar != nil {
			periods = append(periods, getPeriod(periodCalendar, period)...)
		}
	}

	return periods
}

// closureToResponse pass processedMonthlyClosure struct to monthlyclosure
func closureToResponse(closure process_measures.ProcessedMonthlyClosure) MonthlyClosure {

	origin := string(closure.Origin)
	measureType := string(closure.MeasurePointType)
	periods := getPeriods(closure)

	return MonthlyClosure{
		Id: closure.Id,
		MonthlyClosureBase: MonthlyClosureBase{
			InitDate: closure.StartDate.String(),
			EndDate:  closure.EndDate.String(),
			Periods:  &periods,
			//status
			Origen:            &origin,
			MeasureType:       &measureType,
			MeterSerialNumber: &closure.MeterSerialNumber,
			Cups:              &closure.CUPS,
		},
	}
}

// getPeriodUpdate return ProcessedMonthlyClosurePeriod for ProcessedMonthlyClosureCalendar
func getMonthlyClosure(periods Periods) *process_measures.ProcessedMonthlyClosurePeriod {
	p := &process_measures.ProcessedMonthlyClosurePeriod{}
	switch periods.Magnitud {
	case PeriodsMagnitud(measures.AI):
		p.AI = *periods.Absolute
		p.AIi = *periods.Incremental
		break
	case PeriodsMagnitud(measures.AE):
		p.AE = *periods.Absolute
		p.AEi = *periods.Incremental
		break
	case PeriodsMagnitud(measures.R1):
		p.R1 = *periods.Absolute
		p.R1i = *periods.Incremental
		break
	case PeriodsMagnitud(measures.R2):
		p.R2 = *periods.Absolute
		p.R2i = *periods.Incremental
		break
	case PeriodsMagnitud(measures.R3):
		p.R3 = *periods.Absolute
		p.R3i = *periods.Incremental
		break
	case PeriodsMagnitud(measures.R4):
		p.R4 = *periods.Absolute
		p.R4i = *periods.Incremental
		break
	}
	return p
}

// getPeriodUpdate return processedMonthlyClosureCalendar for ProcessedMonthlyClosure
func getPeriodUpdate(periods []Periods) process_measures.ProcessedMonthlyClosureCalendar {

	cal := process_measures.ProcessedMonthlyClosureCalendar{
		P0: &process_measures.ProcessedMonthlyClosurePeriod{},
		P1: &process_measures.ProcessedMonthlyClosurePeriod{},
		P2: &process_measures.ProcessedMonthlyClosurePeriod{},
		P3: &process_measures.ProcessedMonthlyClosurePeriod{},
		P4: &process_measures.ProcessedMonthlyClosurePeriod{},
		P5: &process_measures.ProcessedMonthlyClosurePeriod{},
		P6: &process_measures.ProcessedMonthlyClosurePeriod{},
	}

	for i := 0; i < len(periods); i++ {
		if periods[i].Period == PeriodsPeriodP0 {
			cal.P0 = getMonthlyClosure(periods[i])
		}
		if periods[i].Period == PeriodsPeriodP1 {
			cal.P1 = getMonthlyClosure(periods[i])
		}
		if periods[i].Period == PeriodsPeriodP2 {
			cal.P2 = getMonthlyClosure(periods[i])
		}
		if periods[i].Period == PeriodsPeriodP3 {
			cal.P3 = getMonthlyClosure(periods[i])
		}
		if periods[i].Period == PeriodsPeriodP4 {
			cal.P4 = getMonthlyClosure(periods[i])
		}
		if periods[i].Period == PeriodsPeriodP5 {
			cal.P5 = getMonthlyClosure(periods[i])
		}
		if periods[i].Period == PeriodsPeriodP6 {
			cal.P6 = getMonthlyClosure(periods[i])
		}
	}

	return cal
}

// closureToUpdateOrCreate pass monthlyclosure struct to processedMonthlyClosure
func closureToUpdateOrCreate(month MonthlyClosure) process_measures.ProcessedMonthlyClosure {

	period := getPeriodUpdate(*month.Periods)
	layout := "2006-01-02 15:04:05 +0000 UTC"

	start, _ := time.Parse(layout, month.MonthlyClosureBase.InitDate)
	end, _ := time.Parse(layout, month.MonthlyClosureBase.EndDate)

	return process_measures.ProcessedMonthlyClosure{
		Id:                month.Id,
		CUPS:              *month.Cups,
		StartDate:         start,
		EndDate:           end,
		ValidationStatus:  measures.Status(month.Status),
		Origin:            measures.OriginType(*month.Origen),
		MeasurePointType:  measures.MeasurePointType(*month.MeasureType),
		MeterSerialNumber: *month.MeterSerialNumber,
		CalendarPeriods:   period,
	}
}

// resumeToResponse convierte los dos ProcessedMonthlyClosure en el dto que se va a devolver
func resumeToResponse(closure process_measures.ResumesProcessMonthlyClosure) ProcessMeasuresResume {
	return ProcessMeasuresResume{
		NextReadingClosure:     resumeToReadingClosure(closure.Next),
		PreviousReadingClosure: resumeToReadingClosure(closure.Previous),
	}
}
func resumeToReadingClosure(pmc *process_measures.ProcessedMonthlyClosure) *ReadingClosure {
	if pmc == nil {
		return &ReadingClosure{}
	}

	origin := string(pmc.Origin)
	measureType := string(pmc.MeasurePointType)
	sd := openapi_types.Date{pmc.StartDate}
	ed := openapi_types.Date{pmc.EndDate}
	return &ReadingClosure{
		MeasureDevice: &pmc.MeterSerialNumber,
		MeasureType:   &measureType,
		Origin:        &origin,
		StartDate:     &sd,
		EndDate:       &ed,
		P0:            resumeToMagnitudePeriodFeatures(pmc.CalendarPeriods.P0),
		P1:            resumeToMagnitudePeriodFeatures(pmc.CalendarPeriods.P1),
		P2:            resumeToMagnitudePeriodFeatures(pmc.CalendarPeriods.P2),
		P3:            resumeToMagnitudePeriodFeatures(pmc.CalendarPeriods.P3),
		P4:            resumeToMagnitudePeriodFeatures(pmc.CalendarPeriods.P4),
		P5:            resumeToMagnitudePeriodFeatures(pmc.CalendarPeriods.P5),
		P6:            resumeToMagnitudePeriodFeatures(pmc.CalendarPeriods.P6),
	}
}
func resumeToMagnitudePeriodFeatures(p *process_measures.ProcessedMonthlyClosurePeriod) *MagnitudePeriodFeatures {
	return &MagnitudePeriodFeatures{
		Ae: resumeToMagnitudeFeatures(p.AE, p.AEi),
		Ai: resumeToMagnitudeFeatures(p.AI, p.AIi),
		R1: resumeToMagnitudeFeatures(p.R1, p.R1i),
		R2: resumeToMagnitudeFeatures(p.R2, p.R2i),
		R3: resumeToMagnitudeFeatures(p.R4, p.R3i),
		R4: resumeToMagnitudeFeatures(p.R2, p.R4i),
	}
}
func resumeToMagnitudeFeatures(consum float64, reading float64) *MagnitudeFeatures {
	cns := float32(consum)
	rdng := float32(reading)
	return &MagnitudeFeatures{
		Consum:        &cns,
		MaxDemand:     nil,
		MaxDemandDate: nil,
		Reading:       &rdng,
	}
}
