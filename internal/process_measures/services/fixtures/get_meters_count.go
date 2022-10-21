package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

var RESULT_GET_METERS_AND_COUNT = measures.GetMetersAndCountResult{
	Data: []measures.GetMetersAndCountData{
		{
			Cups: "ESXXXX01",
			Meters: []measures.GetMetersAndCountMeter{
				{
					CurveType: measures.Hourly,
					EndDate:   time.Date(2022, 05, 25, 22, 0, 0, 0, time.UTC),
					StartDate: time.Date(2021, 01, 1, 23, 0, 0, 0, time.UTC),
				},
				{
					CurveType: measures.Hourly,
					EndDate:   time.Date(2050, 01, 1, 23, 0, 0, 0, time.UTC),
					StartDate: time.Date(2022, 05, 25, 22, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			Cups: "ESXXXX02",
			Meters: []measures.GetMetersAndCountMeter{
				{
					CurveType: measures.Hourly,
					EndDate:   time.Date(2022, 05, 31, 22, 0, 0, 0, time.UTC),
					StartDate: time.Date(2021, 01, 1, 23, 0, 0, 0, time.UTC),
				},
				{
					CurveType: measures.Hourly,
					EndDate:   time.Date(2050, 01, 1, 23, 0, 0, 0, time.UTC),
					StartDate: time.Date(2022, 05, 31, 22, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			Cups: "ESXXXX03",
			Meters: []measures.GetMetersAndCountMeter{
				{
					CurveType: measures.Hourly,
					EndDate:   time.Date(2022, 05, 15, 22, 0, 0, 0, time.UTC),
					StartDate: time.Date(2021, 01, 1, 23, 0, 0, 0, time.UTC),
				},
				{
					CurveType: measures.QuarterHour,
					EndDate:   time.Date(2050, 01, 1, 23, 0, 0, 0, time.UTC),
					StartDate: time.Date(2022, 05, 25, 22, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			Cups: "ESXXXX04",
			Meters: []measures.GetMetersAndCountMeter{
				{
					CurveType: measures.Both,
					EndDate:   time.Date(2023, 05, 1, 22, 0, 0, 0, time.UTC),
					StartDate: time.Date(2021, 01, 1, 23, 0, 0, 0, time.UTC),
				},
				{
					CurveType: measures.QuarterHour,
					EndDate:   time.Date(2050, 01, 1, 23, 0, 0, 0, time.UTC),
					StartDate: time.Date(2023, 05, 1, 22, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			Cups: "ESXXXX05",
			Meters: []measures.GetMetersAndCountMeter{
				{
					CurveType: measures.Both,
					EndDate:   time.Date(2023, 05, 1, 22, 0, 0, 0, time.UTC),
					StartDate: time.Date(2021, 01, 1, 23, 0, 0, 0, time.UTC),
				},
				{
					CurveType: measures.QuarterHour,
					EndDate:   time.Date(2050, 01, 1, 23, 0, 0, 0, time.UTC),
					StartDate: time.Date(2023, 05, 1, 22, 0, 0, 0, time.UTC),
				},
			},
		},
	},
	Count: 10,
}
