package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

//CCH-PerfiladoRee-Saldo
var toTransformCchOutlinedBalance = []billing_measures.BillingLoadCurve{
	//P1
	{
		EndDate: time.Date(2022, 05, 31, 1, 0, 0, 0, time.UTC),
		Period:  measures.P1,
	},
	{
		EndDate: time.Date(2022, 05, 31, 2, 0, 0, 0, time.UTC),
		Period:  measures.P1,
	},
	{
		EndDate: time.Date(2022, 05, 31, 3, 0, 0, 0, time.UTC),
		Period:  measures.P1,
	},
	//P2
	{
		EndDate: time.Date(2022, 05, 31, 4, 0, 0, 0, time.UTC),
		Period:  measures.P2,
	},
	{
		EndDate: time.Date(2022, 05, 31, 5, 0, 0, 0, time.UTC),
		Period:  measures.P2,
	},
	{
		EndDate: time.Date(2022, 05, 31, 6, 0, 0, 0, time.UTC),
		Period:  measures.P2,
	},
	//P3
	{
		EndDate: time.Date(2022, 05, 31, 7, 0, 0, 0, time.UTC),
		Period:  measures.P2,
	},
	{
		EndDate: time.Date(2022, 05, 31, 8, 0, 0, 0, time.UTC),
		Period:  measures.P2,
	},
	{
		EndDate: time.Date(2022, 05, 31, 9, 0, 0, 0, time.UTC),
		Period:  measures.P2,
	},
	{
		EndDate: time.Date(2022, 05, 31, 10, 0, 0, 0, time.UTC),
		Period:  measures.P2,
	},
}

//CCH-CchComplete
var toTransformCchComplete = []billing_measures.BillingLoadCurve{
	{
		Origin:    measures.STM,
		Equipment: measures.Redundant,
	},
	{
		Origin:    measures.STM,
		Equipment: measures.Redundant,
	},
	{
		Origin:    measures.STM,
		Equipment: measures.Redundant,
	},
	{
		Origin:    measures.STM,
		Equipment: measures.Redundant,
	},
	{
		Origin:    measures.STM,
		Equipment: measures.Redundant,
	},
	{
		Origin:    measures.STM,
		Equipment: measures.Redundant,
	},
	{
		Origin:    measures.STM,
		Equipment: measures.Redundant,
	},
	{
		Origin:    measures.STM,
		Equipment: measures.Redundant,
	},
	{
		Origin:    measures.Manual,
		Equipment: measures.Receipt,
	},
	{
		Origin:    measures.Manual,
		Equipment: measures.Receipt,
	},
}
