package closing_tpl

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

var Result_iivm_420796_1555_2_20210803 = []gross_measures.MeasureCloseWrite{
	{
		StartDate:         time.Date(2021, time.July, 31, 22, 0, 0, 0, time.UTC),
		EndDate:           time.Date(2021, time.June, 30, 22, 0, 0, 0, time.UTC),
		ReadingDate:       time.Date(2021, time.August, 2, 22, 0, 0, 0, time.UTC),
		Type:              measures.Absolute,
		Contract:          "1",
		ReadingType:       measures.BillingClosure,
		MeterSerialNumber: "420796",
		ConcentratorID:    "",
		File:              "bucketname/DistribuidorX/Input/TPL/iivm_420796_1555_2_20210803.tpl",
		DistributorID:     "",
		DistributorCDOS:   "DistribuidorX",
		Origin:            measures.TPL,
		Periods: []gross_measures.MeasureClosePeriod{
			{
				Period: measures.P0,
				AI:     3.41404e+06,
				AE:     3.41404e+06,
				R1:     1.03735e+06,
				R2:     2.118087e+06,
				R3:     1.03735e+06,
				R4:     2.118087e+06,
				MX:     77,
				FX:     time.Date(2021, time.July, 19, 10, 45, 0, 0, time.UTC),
				E:      0,
			},
			{
				Period: measures.P1,
				AI:     614088,
				AE:     614088,
				R1:     178918,
				R2:     360186,
				R3:     178918,
				R4:     360186,
				MX:     77,
				FX:     time.Date(2021, time.July, 19, 10, 45, 0, 0, time.UTC),
				E:      0,
			},
		},
		Qualifier: "",
	},
	{
		StartDate:         time.Date(2021, time.July, 31, 22, 0, 0, 0, time.UTC),
		EndDate:           time.Date(2021, time.June, 30, 22, 0, 0, 0, time.UTC),
		ReadingDate:       time.Date(2021, time.August, 2, 22, 0, 0, 0, time.UTC),
		Type:              measures.Incremental,
		Contract:          "1",
		ReadingType:       measures.BillingClosure,
		MeterSerialNumber: "420796",
		ConcentratorID:    "",
		File:              "bucketname/DistribuidorX/Input/TPL/iivm_420796_1555_2_20210803.tpl",
		DistributorID:     "",
		DistributorCDOS:   "DistribuidorX",
		Origin:            measures.TPL,
		Qualifier:         "",
		Periods: []gross_measures.MeasureClosePeriod{
			{
				Period: measures.P0,
				AI:     17703,
				AE:     0,
				R1:     0,
				R2:     58109,
				R3:     0,
				R4:     0,
				MX:     77,
				FX:     time.Date(2021, time.July, 19, 10, 45, 0, 0, time.UTC),
				E:      0,
			},
			{
				Period: measures.P1,
				AI:     6206,
				AE:     0,
				R1:     0,
				R2:     13320,
				R3:     0,
				R4:     0,
				MX:     77,
				FX:     time.Date(2021, time.July, 19, 10, 45, 0, 0, time.UTC),
				E:      0,
			},
		},
	},
}
