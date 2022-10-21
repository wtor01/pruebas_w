package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

var valueProfilePartial = 0.000135829136

var Fixture_consume_profiles_partial = []billing_measures.ConsumProfile{
	{
		Date:    time.Date(2022, 05, 31, 10, 0, 0, 0, time.UTC),
		Version: 0,
		Type:    "",
		CoefA:   &valueProfilePartial,
		CoefB:   &valueProfilePartial,
		CoefC:   &valueProfilePartial,
		CoefD:   &valueProfilePartial,
	}, {
		Date:    time.Date(2022, 05, 31, 11, 0, 0, 0, time.UTC),
		Version: 0,
		Type:    "",
		CoefA:   &valueProfilePartial,
		CoefB:   &valueProfilePartial,
		CoefC:   &valueProfilePartial,
		CoefD:   &valueProfilePartial,
	}, {
		Date:    time.Date(2022, 05, 31, 12, 0, 0, 0, time.UTC),
		Version: 0,
		Type:    "",
		CoefA:   &valueProfilePartial,
		CoefB:   &valueProfilePartial,
		CoefC:   &valueProfilePartial,
		CoefD:   &valueProfilePartial,
	}, {
		Date:    time.Date(2022, 05, 31, 13, 0, 0, 0, time.UTC),
		Version: 0,
		Type:    "",
		CoefA:   &valueProfilePartial,
		CoefB:   &valueProfilePartial,
		CoefC:   &valueProfilePartial,
		CoefD:   &valueProfilePartial,
	}, {
		Date:    time.Date(2022, 05, 31, 14, 0, 0, 0, time.UTC),
		Version: 0,
		Type:    "",
		CoefA:   &valueProfilePartial,
		CoefB:   &valueProfilePartial,
		CoefC:   &valueProfilePartial,
		CoefD:   &valueProfilePartial,
	}}

var valueProfileExamplePartial13 = 0.000154928318
var valueProfileExamplePartial14 = 0.000148656018
var valueProfileExamplePartial20 = 0.000116301350
var valueProfileExamplePartial21 = 0.000109264826

var Fixture_consume_profiles_example_partial = []billing_measures.ConsumProfile{
	{
		Date:    time.Date(2022, 05, 18, 13, 0, 0, 0, time.UTC),
		Version: 0,
		Type:    "",
		CoefA:   &valueProfileExamplePartial13,
		CoefB:   &valueProfileExamplePartial13,
		CoefC:   &valueProfileExamplePartial13,
		CoefD:   &valueProfileExamplePartial13,
	}, {
		Date:    time.Date(2022, 05, 18, 14, 0, 0, 0, time.UTC),
		Version: 0,
		Type:    "",
		CoefA:   &valueProfileExamplePartial14,
		CoefB:   &valueProfileExamplePartial14,
		CoefC:   &valueProfileExamplePartial14,
		CoefD:   &valueProfileExamplePartial14,
	}, {
		Date:    time.Date(2022, 05, 18, 20, 0, 0, 0, time.UTC),
		Version: 0,
		Type:    "",
		CoefA:   &valueProfileExamplePartial20,
		CoefB:   &valueProfileExamplePartial20,
		CoefC:   &valueProfileExamplePartial20,
		CoefD:   &valueProfileExamplePartial20,
	}, {
		Date:    time.Date(2022, 05, 18, 21, 0, 0, 0, time.UTC),
		Version: 0,
		Type:    "",
		CoefA:   &valueProfileExamplePartial21,
		CoefB:   &valueProfileExamplePartial21,
		CoefC:   &valueProfileExamplePartial21,
		CoefD:   &valueProfileExamplePartial21,
	}}

var BillingResult_1_Input_partial = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientA,
	BillingBalance: billing_measures.BillingBalance{
		P0: &billing_measures.BillingBalancePeriod{
			AE:              150000,
			EstimatedCodeAE: 1,
		},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			Origin:                   measures.STG,
			AE:                       15200,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			Origin:                   measures.STG,
			AE:                       0,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			Origin:                   measures.STG,
			AE:                       35200,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			Origin:                   measures.STG,
			AE:                       46020,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			Origin:                   measures.STG,
			AE:                       0,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
	},
}

var BillingResult_1_Output_partial = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientA,
	BillingBalance: billing_measures.BillingBalance{
		P0: &billing_measures.BillingBalancePeriod{
			AE:              150000,
			EstimatedCodeAE: 1,
		},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			Origin:                   measures.STG,
			AE:                       15200,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			Origin:                   measures.STG,
			AE:                       0,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			Origin:                   measures.STG,
			AE:                       35200,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			Origin:                   measures.STG,
			AE:                       46020,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			Origin:                   measures.STG,
			AE:                       0,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
	},
}

var BillingResult_2_Input_partial = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientA,
	BillingBalance: billing_measures.BillingBalance{
		P0: &billing_measures.BillingBalancePeriod{
			AE:              150000,
			EstimatedCodeAE: 1,
		},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			Origin:          measures.STG,
			AE:              15200,
			Period:          measures.P0,
			EstimatedCodeAE: 1,
		},
		{
			Origin:          measures.STG,
			AE:              0,
			Period:          measures.P0,
			EstimatedCodeAE: 1,
		},
		{
			Origin:          measures.STG,
			AE:              35200,
			Period:          measures.P0,
			EstimatedCodeAE: 1,
		},
		{
			Origin:          measures.STG,
			AE:              46020,
			Period:          measures.P0,
			EstimatedCodeAE: 1,
		},
		{
			EndDate:         time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC),
			Origin:          measures.Filled,
			AE:              0,
			Period:          measures.P0,
			EstimatedCodeAE: 3,
		},
	},
}

var BillingResult_2_Output_partial = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientA,
	BillingBalance: billing_measures.BillingBalance{
		P0: &billing_measures.BillingBalancePeriod{
			AE:              150000,
			EstimatedCodeAE: 1,
		},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			Origin:                   measures.STG,
			AE:                       15200,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			Origin:                   measures.STG,
			AE:                       0,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			Origin:                   measures.STG,
			AE:                       35200,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			Origin:                   measures.STG,
			AE:                       46020,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       53580,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
	},
}

var BillingResult_3_Input_partial = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientB,
	BillingBalance: billing_measures.BillingBalance{
		P0: &billing_measures.BillingBalancePeriod{
			AE:              150000,
			EstimatedCodeAE: 1,
		},
		P1: &billing_measures.BillingBalancePeriod{
			AE:              13000,
			EstimatedCodeAE: 1,
		},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			EndDate: time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC), Origin: measures.STG,
			AE:              15200,
			Period:          measures.P0,
			EstimatedCodeAE: 1,
		},
		{
			EndDate: time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC), Origin: measures.STG,
			AE:              0,
			Period:          measures.P0,
			EstimatedCodeAE: 1,
		},
		{
			EndDate: time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC), Origin: measures.STG,
			AE:              35200,
			Period:          measures.P0,
			EstimatedCodeAE: 1,
		},
		{
			EndDate: time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC), Origin: measures.Filled,
			AE:              46020,
			Period:          measures.P1,
			EstimatedCodeAE: 3,
		},
		{
			EndDate: time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC), Origin: measures.Filled,
			AE:              0,
			Period:          measures.P0,
			EstimatedCodeAE: 3,
		},
	},
}

var BillingResult_3_Output_partial = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientB,
	BillingBalance: billing_measures.BillingBalance{
		P0: &billing_measures.BillingBalancePeriod{
			AE:              150000,
			EstimatedCodeAE: 1,
		},
		P1: &billing_measures.BillingBalancePeriod{
			AE:              13000,
			EstimatedCodeAE: 1,
		},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			EndDate:                  time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC),
			Origin:                   measures.STG,
			AE:                       15200,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC),
			Origin:                   measures.STG,
			AE:                       0,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC),
			Origin:                   measures.STG,
			AE:                       35200,
			Period:                   measures.P0,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       13000,
			Period:                   measures.P1,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       99600,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
	},
}

var BillingResult_1_Input_partial_example = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientA,
	BillingBalance: billing_measures.BillingBalance{
		P4: &billing_measures.BillingBalancePeriod{
			AE:              42589,
			EstimatedCodeAE: 1,
		},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			EndDate:         time.Date(2022, time.May, 18, 10, 0, 0, 0, time.UTC),
			Origin:          measures.STG,
			AE:              4532,
			Period:          measures.P4,
			EstimatedCodeAE: 1,
		},
		{
			EndDate:         time.Date(2022, time.May, 18, 11, 0, 0, 0, time.UTC),
			Origin:          measures.STG,
			AE:              4356,
			Period:          measures.P4,
			EstimatedCodeAE: 1,
		},
		{
			EndDate:         time.Date(2022, time.May, 18, 12, 0, 0, 0, time.UTC),
			Origin:          measures.STG,
			AE:              5237,
			Period:          measures.P4,
			EstimatedCodeAE: 1},
		{
			EndDate:         time.Date(2022, time.May, 18, 13, 0, 0, 0, time.UTC),
			Origin:          measures.Filled,
			Period:          measures.P4,
			EstimatedCodeAE: 1,
		},
		{
			EndDate:         time.Date(2022, time.May, 18, 14, 0, 0, 0, time.UTC),
			Origin:          measures.Filled,
			Period:          measures.P4,
			EstimatedCodeAE: 1,
		},
		{
			EndDate:         time.Date(2022, time.May, 18, 19, 0, 0, 0, time.UTC),
			Origin:          measures.STG,
			AE:              5468,
			Period:          measures.P4,
			EstimatedCodeAE: 1,
		},
		{
			EndDate:         time.Date(2022, time.May, 18, 20, 0, 0, 0, time.UTC),
			Origin:          measures.Filled,
			Period:          measures.P4,
			EstimatedCodeAE: 1,
		},
		{
			EndDate:         time.Date(2022, time.May, 18, 21, 0, 0, 0, time.UTC),
			Origin:          measures.Filled,
			Period:          measures.P4,
			EstimatedCodeAE: 1,
		},
		{
			EndDate:         time.Date(2022, time.May, 18, 22, 0, 0, 0, time.UTC),
			Origin:          measures.STG,
			AE:              4589,
			Period:          measures.P4,
			EstimatedCodeAE: 1,
		},
	},
}

var BillingResult_1_Output_partial_example = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientA,
	BillingBalance: billing_measures.BillingBalance{
		P4: &billing_measures.BillingBalancePeriod{
			AE:              42589,
			EstimatedCodeAE: 1,
		},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			EndDate:                  time.Date(2022, time.May, 18, 10, 0, 0, 0, time.UTC),
			Origin:                   measures.STG,
			AE:                       4532,
			Period:                   measures.P4,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 11, 0, 0, 0, time.UTC),
			Origin:                   measures.STG,
			AE:                       4356,
			Period:                   measures.P4,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 12, 0, 0, 0, time.UTC),
			Origin:                   measures.STG,
			AE:                       5237,
			Period:                   measures.P4,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 13, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       5389,
			Period:                   measures.P4,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 14, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       5171,
			Period:                   measures.P4,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 19, 0, 0, 0, time.UTC),
			Origin:                   measures.STG,
			AE:                       5468,
			Period:                   measures.P4,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 20, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       4046,
			Period:                   measures.P4,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 21, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       3801,
			Period:                   measures.P4,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 22, 0, 0, 0, time.UTC),
			Origin:                   measures.STG,
			AE:                       4589,
			Period:                   measures.P4,
			EstimatedCodeAE:          1,
			EstimatedMethodAE:        billing_measures.RealBalance,
			EstimatedGeneralMethodAE: billing_measures.GeneralReal,
			MeasureTypeAE:            billing_measures.FirmBalanceMeasure,
		},
	},
}
