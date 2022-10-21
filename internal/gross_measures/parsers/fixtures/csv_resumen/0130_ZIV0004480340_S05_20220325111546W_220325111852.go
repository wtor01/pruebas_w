package csv_resumen

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

var Result_0130_ZIV0004480340_S05_20220325111546W_220325111852 = []gross_measures.MeasureCloseWrite{
	{
		EndDate:           time.Date(2022, time.March, 23, 23, 0, 0, 0, time.UTC),
		ReadingDate:       time.Date(2022, time.March, 25, 10, 15, 46, 0, time.UTC),
		Type:              measures.Absolute,
		ReadingType:       measures.DailyClosure,
		MeterSerialNumber: "ES0130000001329031HX",
		ConcentratorID:    "",
		File:              "bucketname/DistribuidorX/Input/CSV/0130_ZIV0004480340_S05_20220325111546W_220325111852.csv",
		DistributorID:     "",
		DistributorCDOS:   "DistribuidorX",
		Origin:            measures.File,
		Contract:          "1",
		Periods: []gross_measures.MeasureClosePeriod{
			{
				Period: measures.P0,
				AI:     7.0171e+07,
				AE:     0,
				R1:     0,
				R2:     0,
				R3:     0,
				R4:     1.9786e+07,
			},
			{
				Period: measures.P1,
				AI:     2.7326e+07,
				AE:     0,
				R1:     0,
				R2:     0,
				R3:     0,
				R4:     7.462e+06,
			},
			{
				Period: measures.P2,
				AI:     3.5587e+07,
				AE:     0,
				R1:     0,
				R2:     0,
				R3:     0,
				R4:     1.0191e+07,
			},
			{
				Period: measures.P3,
				AI:     7.257e+06,
				AE:     0,
				R1:     0,
				R2:     0,
				R3:     0,
				R4:     2.132e+06,
			},
			{
				Period: measures.P4,
				AI:     0,
				AE:     0,
				R1:     0,
				R2:     0,
				R3:     0,
				R4:     0,
			},
			{
				Period: measures.P5,
				AI:     0,
				AE:     0,
				R1:     0,
				R2:     0,
				R3:     0,
				R4:     0,
			},
			{
				Period: measures.P6,
				AI:     0,
				AE:     0,
				R1:     0,
				R2:     0,
				R3:     0,
				R4:     0,
			},
		},
	},
}
