package fixtures

import "bitbucket.org/sercide/data-ingestion/internal/measures"

var RESULT_DASHBOARD_CUPS_LIST = measures.DashboardListCups{
	Cups: []measures.DashboardCups{
		{
			Cups: "ESXXXX01",
			Values: measures.DashboardCupsReading{
				Curve: measures.DashboardCupsValues{
					Valid:     541,
					Invalid:   63,
					Supervise: 40,
					None:      100,
					Total:     744,
					ShouldBe:  744,
				},
				Daily: measures.DashboardCupsValues{
					Valid:    31,
					Total:    31,
					ShouldBe: 31,
				},
				Monthly: measures.DashboardCupsValues{
					Valid:    1,
					Total:    1,
					ShouldBe: 2,
				},
			},
		},
		{
			Cups: "ESXXXX02",
			Values: measures.DashboardCupsReading{
				Curve: measures.DashboardCupsValues{
					Valid:     500,
					Invalid:   30,
					Supervise: 40,
					None:      100,
					Total:     670,
					ShouldBe:  744,
				},
				Daily: measures.DashboardCupsValues{
					Valid:    29,
					Total:    29,
					ShouldBe: 31,
				},
				Monthly: measures.DashboardCupsValues{
					Valid:    1,
					Total:    1,
					ShouldBe: 1,
				},
			},
		},
		{
			Cups: "ESXXXX03",
			Values: measures.DashboardCupsReading{
				Curve: measures.DashboardCupsValues{
					Valid:     1000,
					Invalid:   500,
					Supervise: 150,
					Total:     1650,
					ShouldBe:  1896,
				},
				Daily: measures.DashboardCupsValues{
					Valid:    24,
					Invalid:  7,
					Total:    31,
					ShouldBe: 31,
				},
				Monthly: measures.DashboardCupsValues{
					Valid:    2,
					Total:    2,
					ShouldBe: 2,
				},
			},
		},
		{
			Cups: "ESXXXX04",
			Values: measures.DashboardCupsReading{
				Curve: measures.DashboardCupsValues{
					Valid:    60000,
					Invalid:  11000,
					Total:    71000,
					ShouldBe: 71424,
				},
				Daily: measures.DashboardCupsValues{
					Valid:    31,
					Total:    31,
					ShouldBe: 31,
				},
				Monthly: measures.DashboardCupsValues{
					Valid:    1,
					Total:    1,
					ShouldBe: 1,
				},
			},
		},
		{
			Cups: "ESXXXX05",
			Values: measures.DashboardCupsReading{
				Curve: measures.DashboardCupsValues{
					ShouldBe: 71424,
				},
				Daily: measures.DashboardCupsValues{
					ShouldBe: 31,
				},
				Monthly: measures.DashboardCupsValues{
					ShouldBe: 1,
				},
			},
		},
	},
	Total: 10,
}
