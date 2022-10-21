package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"time"
)

var ProcessMeasures_Input_averages_conditional_previous_filled = []process_measures.ProcessedLoadCurve{{
	EndDate: time.Date(2022, time.May, 8, 23, 0, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: "STM", Period: "P1", MeasurePointType: "P",
}, {
	EndDate: time.Date(2022, time.May, 8, 0, 0, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: "FILLED", Period: "P1", MeasurePointType: "P",
}}

var ProcessMeasures_Input_averages_conditional_previous_not_filled = []process_measures.ProcessedLoadCurve{{
	EndDate: time.Date(2022, time.May, 8, 23, 0, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: "STM", Period: "P1", MeasurePointType: "P",
}, {
	EndDate: time.Date(2022, time.May, 8, 0, 0, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: "STM", Period: "P1", MeasurePointType: "P",
}}

var ProcessMeasures_1_Input_averages_conditional_previous = []process_measures.ProcessedLoadCurve{{
	EndDate: time.Date(2022, time.May, 8, 0, 0, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: "STM", Period: "P1", MeasurePointType: "P",
}}

var ProcessMeasures_Input_averages_conditional_next_filled = []process_measures.ProcessedLoadCurve{{
	EndDate: time.Date(2022, time.May, 8, 0, 0, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: "FILLED", Period: "P1", MeasurePointType: "P",
}, {
	EndDate: time.Date(2022, time.May, 8, 1, 0, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: "STM", Period: "P1", MeasurePointType: "P",
}}

var ProcessMeasures_Input_averages_conditional_next_not_filled = []process_measures.ProcessedLoadCurve{{
	EndDate: time.Date(2022, time.May, 8, 0, 0, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: "STM", Period: "P1", MeasurePointType: "P",
}, {
	EndDate: time.Date(2022, time.May, 8, 1, 0, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: "FILLED", Period: "P1", MeasurePointType: "P",
}}

var BillingResult_1_Input_averages_conditional = billing_measures.BillingMeasure{Id: "1", Periods: []measures.PeriodKey{"P1", "P2", "P3"}, ActualReadingClosure: measures.DailyReadingClosure{CalendarPeriods: measures.DailyReadingClosureCalendar{
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
		EndDate: time.Date(2022, time.May, 31, 9, 0, 0, 0, time.UTC), Origin: "FILLED", Period: "P1", EstimatedCodeAI: 1},
	billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC), AI: 20000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1},
	billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC), AI: 30000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}}}

var BillingResult_2_Input_averages_conditional = billing_measures.BillingMeasure{Id: "1", Periods: []measures.PeriodKey{"P1", "P2", "P3"}, ActualReadingClosure: measures.DailyReadingClosure{CalendarPeriods: measures.DailyReadingClosureCalendar{
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
		EndDate: time.Date(2022, time.May, 31, 9, 0, 0, 0, time.UTC), AI: 20000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1},
	billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC), AI: 30000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC), Origin: "FILLED", Period: "P1", EstimatedCodeAI: 1}}}

var BillingResult_example_Input_averages_conditional_true = billing_measures.BillingMeasure{Id: "1", Periods: []measures.PeriodKey{"P1", "P2", "P3"}, ActualReadingClosure: measures.DailyReadingClosure{CalendarPeriods: measures.DailyReadingClosureCalendar{
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
		EndDate: time.Date(2022, time.May, 31, 1, 0, 0, 0, time.UTC), AI: 20000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1},
	billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 2, 0, 0, 0, time.UTC), AI: 30000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 3, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 4, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 5, 0, 0, 0, time.UTC), Origin: "FILLED", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 6, 0, 0, 0, time.UTC), Origin: "FILLED", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 7, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 8, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}}}

var BillingResult_example_Input_averages_conditional_false = billing_measures.BillingMeasure{Id: "1", Periods: []measures.PeriodKey{"P1", "P2", "P3"}, ActualReadingClosure: measures.DailyReadingClosure{CalendarPeriods: measures.DailyReadingClosureCalendar{
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
		EndDate: time.Date(2022, time.May, 31, 1, 0, 0, 0, time.UTC), AI: 20000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1},
	billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 2, 0, 0, 0, time.UTC), AI: 30000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 3, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 4, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: "STM", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 5, 0, 0, 0, time.UTC), Origin: "FILLED", Period: "P1", EstimatedCodeAI: 1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 6, 0, 0, 0, time.UTC), Origin: "FILLED", Period: "P1", EstimatedCodeAI: 1}}}
