package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"time"
)

var BillingResult_1_Input_averages = billing_measures.BillingMeasure{
	PointType: "2",
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{

			AI:               30000,
			Origin:           measures.STG,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               40000,
			Origin:           measures.STG,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               50000,
			Origin:           measures.STG,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               60000,
			Origin:           measures.STG,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               72000,
			Origin:           measures.STG,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
	},
}

var BillingResult_1_Output_averages = billing_measures.BillingMeasure{
	PointType: "2",
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{

			AI:                       30000,
			Origin:                   measures.STG,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			AI:                       40000,
			Origin:                   measures.STG,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			AI:                       50000,
			Origin:                   measures.STG,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			AI:                       60000,
			Origin:                   measures.STG,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			AI:                       72000,
			Origin:                   measures.STG,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
	},
}

var BillingResult_2_Input_averages = billing_measures.BillingMeasure{
	PointType: "2",
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{

			AI:               30000,
			Origin:           measures.STG,
			Period:           measures.P1,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			Origin:           measures.Filled,
			Period:           measures.P1,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               50000,
			Origin:           measures.STG,
			Period:           measures.P1,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               60000,
			Origin:           measures.STG,
			Period:           measures.P1,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               72000,
			Origin:           measures.STG,
			Period:           measures.P1,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
	},
}

var BillingResult_2_Output_averages = billing_measures.BillingMeasure{
	PointType: "2",
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{

			Origin:                   measures.STG,
			AI:                       30000,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.Filled,
			AI:                       40000,
			Period:                   measures.P1,
			EstimatedCodeAI:          9,
			EstimatedMethodAI:        billing_measures.EstimateOnlyHistoric,
			EstimatedGeneralMethodAI: billing_measures.GeneralEstimated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.STG,
			AI:                       50000,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.STG,
			AI:                       60000,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.STG,
			AI:                       72000,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
	},
}

var BillingResult_3_Input_averages = billing_measures.BillingMeasure{
	PointType: "2",
	BillingLoadCurve: []billing_measures.BillingLoadCurve{

		{

			AI:               20000,
			Origin:           measures.STG,
			Period:           measures.P2,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},

		{

			AI:               30000,
			Origin:           measures.STG,
			Period:           measures.P1,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			Origin:           measures.Filled,
			Period:           measures.P2,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               50000,
			Origin:           measures.STG,
			Period:           measures.P1,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               60000,
			Origin:           measures.STG,
			Period:           measures.P2,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               72000,
			Origin:           measures.STG,
			Period:           measures.P1,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
	},
}

var BillingResult_3_Output_averages = billing_measures.BillingMeasure{
	PointType: "2",
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{

			Origin:                   measures.STG,
			AI:                       20000,
			Period:                   measures.P2,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP},
		{

			Origin:                   measures.STG,
			AI:                       30000,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP},
		{

			Origin:                   measures.Filled,
			AI:                       40000,
			Period:                   measures.P2,
			EstimatedCodeAI:          9,
			EstimatedMethodAI:        billing_measures.EstimateOnlyHistoric,
			EstimatedGeneralMethodAI: billing_measures.GeneralEstimated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP},
		{

			Origin:                   measures.STG,
			AI:                       50000,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP},
		{

			Origin:                   measures.STG,
			AI:                       60000,
			Period:                   measures.P2,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP},
		{

			Origin:                   measures.STG,
			AI:                       72000,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP},
	},
}

var BillingResult_4_Input_averages = billing_measures.BillingMeasure{
	PointType: "2",
	BillingLoadCurve: []billing_measures.BillingLoadCurve{

		{
			AI:               20000,
			Origin:           measures.STG,
			Period:           measures.P2,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{
			AI:               30000,
			Origin:           measures.STG,
			Period:           measures.P1,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{
			Origin:           measures.Filled,
			Period:           measures.P2,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP},

		{

			Origin:           measures.Filled,
			Period:           measures.P2,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP},

		{
			Origin:           measures.Filled,
			Period:           measures.P2,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP},

		{
			AI:               50000,
			Origin:           measures.STG,
			Period:           measures.P1,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{
			AI:               60000,
			Origin:           measures.STG,
			Period:           measures.P2,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{
			AI:               72000,
			Origin:           measures.STG,
			Period:           measures.P1,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
	},
}

var BillingResult_4_Output_averages = billing_measures.BillingMeasure{
	PointType: "2",
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{

			Origin:                   measures.STG,
			AI:                       20000,
			Period:                   measures.P2,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.STG,
			AI:                       30000,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.Filled,
			AI:                       40000,
			Period:                   measures.P2,
			EstimatedCodeAI:          9,
			EstimatedMethodAI:        billing_measures.EstimateOnlyHistoric,
			EstimatedGeneralMethodAI: billing_measures.GeneralEstimated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.Filled,
			AI:                       40000,
			Period:                   measures.P2,
			EstimatedCodeAI:          9,
			EstimatedMethodAI:        billing_measures.EstimateOnlyHistoric,
			EstimatedGeneralMethodAI: billing_measures.GeneralEstimated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.Filled,
			AI:                       40000,
			Period:                   measures.P2,
			EstimatedCodeAI:          9,
			EstimatedMethodAI:        billing_measures.EstimateOnlyHistoric,
			EstimatedGeneralMethodAI: billing_measures.GeneralEstimated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.STG,
			AI:                       50000,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.STG,
			AI:                       60000,
			Period:                   measures.P2,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.STG,
			AI:                       72000,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
	},
}

var Previous_ProcessedCurve_Example_Input_averages = []process_measures.ProcessedLoadCurve{{
	EndDate: time.Date(2022, time.May, 8, 0, 0, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: measures.STG, Period: measures.P6, MeasurePointType: "P",
}}

var Next_ProcessedCurve_Example_Input_averages = []process_measures.ProcessedLoadCurve{{
	EndDate: time.Date(2022, time.May, 10, 0, 0, 0, 0, time.UTC), AI: 33000, AE: 33000, R1: 33000, R2: 33000, R3: 33000, R4: 33000, Origin: measures.STG, Period: measures.P6, MeasurePointType: "P",
}}
var BillingResult_Example_Input_averages = billing_measures.BillingMeasure{
	PointType: "2",
	BillingLoadCurve: []billing_measures.BillingLoadCurve{

		{

			Origin:           measures.Filled,
			Period:           measures.P6,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			Origin:           measures.Filled,
			Period:           measures.P6,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               24000,
			Origin:           measures.STG,
			Period:           measures.P6,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               25000,
			Origin:           measures.STG,
			Period:           measures.P6,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               19000,
			Origin:           measures.STG,
			Period:           measures.P6,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
	},
}

var BillingResult_Example_Output_averages = billing_measures.BillingMeasure{
	PointType: "2",
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{

			Origin:                   measures.Filled,
			AI:                       28500,
			Period:                   measures.P6,
			EstimatedCodeAI:          9,
			EstimatedMethodAI:        billing_measures.EstimateOnlyHistoric,
			EstimatedGeneralMethodAI: billing_measures.GeneralEstimated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP},
		{

			Origin:                   measures.Filled,
			AI:                       28500,
			Period:                   measures.P6,
			EstimatedCodeAI:          9,
			EstimatedMethodAI:        billing_measures.EstimateOnlyHistoric,
			EstimatedGeneralMethodAI: billing_measures.GeneralEstimated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP},
		{

			Origin:                   measures.STG,
			AI:                       24000,
			Period:                   measures.P6,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP},
		{

			Origin:                   measures.STG,
			AI:                       25000,
			Period:                   measures.P6,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP},
		{

			Origin:                   measures.STG,
			AI:                       19000,
			Period:                   measures.P6,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP},
	},
}

var BillingResult_Example_Input_averages_example_next = billing_measures.BillingMeasure{
	PointType: "2",
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			AI:               24000,
			Origin:           measures.STG,
			Period:           measures.P6,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{
			AI:               25000,
			Origin:           measures.STG,
			Period:           measures.P6,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{
			AI:               19000,
			Origin:           measures.STG,
			Period:           measures.P6,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			Origin:           measures.Filled,
			Period:           measures.P6,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			Origin:           measures.Filled,
			Period:           measures.P6,
			EstimatedCodeAI:  1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
	},
}

var BillingResult_Example_Output_averages_example_next = billing_measures.BillingMeasure{
	PointType: "2",
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			Origin:                   measures.STG,
			AI:                       24000,
			Period:                   measures.P6,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{
			Origin:                   measures.STG,
			AI:                       25000,
			Period:                   measures.P6,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{
			Origin:                   measures.STG,
			AI:                       19000,
			Period:                   measures.P6,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{
			Origin:                   measures.Filled,
			AI:                       26000,
			Period:                   measures.P6,
			EstimatedCodeAI:          9,
			EstimatedMethodAI:        billing_measures.EstimateOnlyHistoric,
			EstimatedGeneralMethodAI: billing_measures.GeneralEstimated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{
			Origin:                   measures.Filled,
			AI:                       26000,
			Period:                   measures.P6,
			EstimatedCodeAI:          9,
			EstimatedMethodAI:        billing_measures.EstimateOnlyHistoric,
			EstimatedGeneralMethodAI: billing_measures.GeneralEstimated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
	},
}
