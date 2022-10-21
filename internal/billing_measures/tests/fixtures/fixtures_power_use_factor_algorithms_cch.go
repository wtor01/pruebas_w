package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
)

var BillingResult_1_Input_Power_Use_Factor = billing_measures.BillingMeasure{
	PointType: "2",
	BillingBalance: billing_measures.BillingBalance{P0: &billing_measures.BillingBalancePeriod{
		AI:              200000,
		EstimatedCodeAI: 1,
	},
	}, BillingLoadCurve: []billing_measures.BillingLoadCurve{
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

var BillingResult_2_Input_Power_Use_Factor = billing_measures.BillingMeasure{
	PointType: "2",

	BillingBalance: billing_measures.BillingBalance{P0: &billing_measures.BillingBalancePeriod{
		AI:              200000,
		EstimatedCodeAI: 1,
	},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{

			AI:               30000,
			Origin:           measures.STG,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			Origin:           measures.Filled,
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
		}},
}

var BillingResult_2_Output_Power_Use_Factor = billing_measures.BillingMeasure{
	PointType: "2",
	BillingBalance: billing_measures.BillingBalance{
		P0: &billing_measures.BillingBalancePeriod{
			AI:              200000,
			EstimatedCodeAI: 1,
		},
	}, BillingLoadCurve: []billing_measures.BillingLoadCurve{
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
			Period:                   measures.P1,
			EstimatedCodeAI:          11,
			EstimatedMethodAI:        billing_measures.PowerUseFactor,
			EstimatedGeneralMethodAI: billing_measures.GeneralCalculated,
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
		}},
}

var BillingResult_Example_Input_Power_Use_Factor = billing_measures.BillingMeasure{
	PointType: "2",
	BillingBalance: billing_measures.BillingBalance{P0: &billing_measures.BillingBalancePeriod{
		AI:              200000,
		EstimatedCodeAI: 1,
	},
	}, BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{

			AI:               22000,
			Origin:           measures.STG,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               30000,
			Origin:           measures.STG,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			Origin:           measures.Filled,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			Origin:           measures.Filled,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeR,
		},
		{

			Origin:           measures.Filled,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeR,
		},
		{

			Origin:           measures.Filled,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               36000,
			Origin:           measures.STG,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			AI:               29000,
			Origin:           measures.STG,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
		{

			Origin:           measures.Filled,
			Period:           measures.P1,
			MeasurePointType: measures.MeasurePointTypeP,
		},
	},
}

var BillingResult_Example_Output_Power_Use_Factor = billing_measures.BillingMeasure{
	PointType: "2",
	BillingBalance: billing_measures.BillingBalance{P0: &billing_measures.BillingBalancePeriod{
		AI:              200000,
		EstimatedCodeAI: 1,
	},
	}, BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{

			Origin:                   measures.STG,
			AI:                       22000,
			Period:                   measures.P1,
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
			Period:                   measures.P1,
			EstimatedCodeAI:          11,
			EstimatedMethodAI:        billing_measures.PowerUseFactor,
			EstimatedGeneralMethodAI: billing_measures.GeneralCalculated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.Filled,
			Period:                   measures.P1,
			EstimatedCodeAI:          11,
			EstimatedMethodAI:        billing_measures.PowerUseFactor,
			EstimatedGeneralMethodAI: billing_measures.GeneralCalculated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeR,
		},
		{

			Origin:                   measures.Filled,
			Period:                   measures.P1,
			EstimatedCodeAI:          11,
			EstimatedMethodAI:        billing_measures.PowerUseFactor,
			EstimatedGeneralMethodAI: billing_measures.GeneralCalculated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeR,
		},
		{

			Origin:                   measures.Filled,
			Period:                   measures.P1,
			EstimatedCodeAI:          11,
			EstimatedMethodAI:        billing_measures.PowerUseFactor,
			EstimatedGeneralMethodAI: billing_measures.GeneralCalculated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.STG,
			AI:                       36000,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.STG,
			AI:                       29000,
			Period:                   measures.P1,
			EstimatedCodeAI:          1,
			EstimatedMethodAI:        billing_measures.FirmMainConfig,
			EstimatedGeneralMethodAI: billing_measures.GeneralReal,
			MeasureTypeAI:            billing_measures.FirmBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
		{

			Origin:                   measures.Filled,
			Period:                   measures.P1,
			EstimatedCodeAI:          11,
			EstimatedMethodAI:        billing_measures.PowerUseFactor,
			EstimatedGeneralMethodAI: billing_measures.GeneralCalculated,
			MeasureTypeAI:            billing_measures.ProvisionalBalanceMeasure,
			MeasurePointType:         measures.MeasurePointTypeP,
		},
	},
}
