package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"time"
)

var ProcessedLoadCurve_ccch_example_input_1 = []process_measures.ProcessedLoadCurve{{

	EndDate: time.Date(2022, time.May, 31, 5, 0, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: measures.STM, Period: measures.P1, MeasurePointType: measures.MeasurePointTypeP,
}, {
	EndDate: time.Date(2022, time.May, 31, 5, 15, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: measures.STM, Period: measures.P1, MeasurePointType: measures.MeasurePointTypeP,
}, {
	EndDate: time.Date(2022, time.May, 31, 5, 30, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: measures.STM, Period: measures.P1, MeasurePointType: measures.MeasurePointTypeP,
}, {
	EndDate: time.Date(2022, time.May, 31, 5, 45, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: measures.STM, Period: measures.P1, MeasurePointType: measures.MeasurePointTypeP,
}, {

	EndDate: time.Date(2022, time.May, 31, 6, 0, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: measures.STM, Period: measures.P1, MeasurePointType: measures.MeasurePointTypeP,
}, {
	EndDate: time.Date(2022, time.May, 31, 6, 15, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: measures.STM, Period: measures.P1, MeasurePointType: measures.MeasurePointTypeP,
}, {
	EndDate: time.Date(2022, time.May, 31, 6, 30, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: measures.STM, Period: measures.P1, MeasurePointType: measures.MeasurePointTypeP,
}, {
	EndDate: time.Date(2022, time.May, 31, 6, 45, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: measures.STM, Period: measures.P1, MeasurePointType: measures.MeasurePointTypeP,
}}
var BillingResult_example_Input_ccch_complete_1 = billing_measures.BillingMeasure{Id: "1", PointType: "3", RegisterType: measures.Hourly, Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3}, ActualReadingClosure: measures.DailyReadingClosure{CalendarPeriods: measures.DailyReadingClosureCalendar{
	P1: &measures.DailyReadingClosureCalendarPeriod{
		Filled: false,
	},
	P2: &measures.DailyReadingClosureCalendarPeriod{
		Filled: false,
	},
	P3: &measures.DailyReadingClosureCalendarPeriod{
		Filled: false,
	},
}}, BillingLoadCurve: []billing_measures.BillingLoadCurve{
	billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 1, 0, 0, 0, time.UTC), AI: 20000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: measures.STM, Period: measures.P1, EstimatedCodeAI: 1},
	billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 2, 0, 0, 0, time.UTC), AI: 30000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: measures.STM, Period: measures.P1, EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 3, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STM, Period: measures.P1, EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 4, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STM, Period: measures.P1, EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 5, 0, 0, 0, time.UTC), Origin: measures.Filled, Period: measures.P1, EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 6, 0, 0, 0, time.UTC), Origin: measures.Filled, Period: measures.P1, EstimatedCodeAI: 1}}}

var BillingResult_example_Output_ccch_complete_1 = billing_measures.BillingMeasure{Id: "1", RegisterType: measures.Hourly, Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3}, ActualReadingClosure: measures.DailyReadingClosure{CalendarPeriods: measures.DailyReadingClosureCalendar{
	P1: &measures.DailyReadingClosureCalendarPeriod{
		Filled: false,
	},
	P2: &measures.DailyReadingClosureCalendarPeriod{
		Filled: false,
	},
	P3: &measures.DailyReadingClosureCalendarPeriod{
		Filled: false,
	},
}}, BillingLoadCurve: []billing_measures.BillingLoadCurve{
	billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 1, 0, 0, 0, time.UTC), AI: 20000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: measures.STM, Period: measures.P1, EstimatedCodeAI: 1, EstimatedMethodAI: "REAL_VALID_MEASURE", EstimatedGeneralMethodAI: "REAL", MeasureTypeAI: "FIRM"},
	billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 2, 0, 0, 0, time.UTC), AI: 30000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: measures.STM, Period: measures.P1, EstimatedCodeAI: 1, EstimatedMethodAI: "REAL_VALID_MEASURE", EstimatedGeneralMethodAI: "REAL", MeasureTypeAI: "FIRM"}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 3, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STM, Period: measures.P1, EstimatedCodeAI: 1, EstimatedMethodAI: "REAL_VALID_MEASURE", EstimatedGeneralMethodAI: "REAL", MeasureTypeAI: "FIRM"}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 4, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STM, Period: measures.P1, EstimatedCodeAI: 1, EstimatedMethodAI: "REAL_VALID_MEASURE", EstimatedGeneralMethodAI: "REAL", MeasureTypeAI: "FIRM"}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 5, 0, 0, 0, time.UTC), AI: 132000, Origin: measures.Filled, Period: measures.P1, EstimatedCodeAI: 1, EstimatedMethodAI: "REAL_VALID_MEASURE", EstimatedGeneralMethodAI: "REAL", MeasureTypeAI: "FIRM"}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 6, 0, 0, 0, time.UTC), AI: 132000, Origin: measures.Filled, Period: measures.P1, EstimatedCodeAI: 1, EstimatedMethodAI: "REAL_VALID_MEASURE", EstimatedGeneralMethodAI: "REAL", MeasureTypeAI: "FIRM"}}}
