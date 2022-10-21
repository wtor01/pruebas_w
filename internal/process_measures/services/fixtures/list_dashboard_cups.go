package fixtures

import "bitbucket.org/sercide/data-ingestion/internal/measures"

var RESULT_LIST_CUPS = map[string]*measures.DashboardCupsReading{
	"ESXXXX01": {
		Curve: measures.DashboardCupsValues{
			Valid:     541,
			Invalid:   63,
			Supervise: 40,
			None:      100,
			Total:     744,
		},
		Daily: measures.DashboardCupsValues{
			Valid: 31,
			Total: 31,
		},
		Monthly: measures.DashboardCupsValues{
			Valid: 1,
			Total: 1,
		},
	},
	"ESXXXX02": {
		Curve: measures.DashboardCupsValues{
			Valid:     500,
			Invalid:   30,
			Supervise: 40,
			None:      100,
			Total:     670,
		},
		Daily: measures.DashboardCupsValues{
			Valid: 29,
			Total: 29,
		},
		Monthly: measures.DashboardCupsValues{
			Valid: 1,
			Total: 1,
		},
	},
	"ESXXXX03": {
		Curve: measures.DashboardCupsValues{
			Valid:     1000,
			Invalid:   500,
			Supervise: 150,
			Total:     1650,
		},
		Daily: measures.DashboardCupsValues{
			Valid:   24,
			Invalid: 7,
			Total:   31,
		},
		Monthly: measures.DashboardCupsValues{
			Valid: 2,
			Total: 2,
		},
	},
	"ESXXXX04": {
		Curve: measures.DashboardCupsValues{
			Valid:   60000,
			Invalid: 11000,
			Total:   71000,
		},
		Daily: measures.DashboardCupsValues{
			Valid: 31,
			Total: 31,
		},
		Monthly: measures.DashboardCupsValues{
			Valid: 1,
			Total: 1,
		},
	},
	"ESXXXX05": {},
}
