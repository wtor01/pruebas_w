package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

var valueProfileTotal = 0.000135829136
var valueProfileExampleTotal10 = 0.000143329677
var valueProfileExampleTotal11 = 0.000150676079
var valueProfileExampleTotal12 = 0.000152627050
var valueProfileExampleTotal13 = 0.000154928318
var valueProfileExampleTotal14 = 0.000148656018
var valueProfileExampleTotal19 = 0.000122299751
var valueProfileExampleTotal20 = 0.000116301350
var valueProfileExampleTotal21 = 0.000109264826
var valueProfileExampleTotal22 = 0.000109298395

var Fixture_consume_profiles_total = []billing_measures.ConsumProfile{{
	Date:    time.Date(2022, 05, 31, 10, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileTotal,
	CoefB:   &valueProfileTotal,
	CoefC:   &valueProfileTotal,
	CoefD:   &valueProfileTotal,
}, {
	Date:    time.Date(2022, 05, 31, 11, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileTotal,
	CoefB:   &valueProfileTotal,
	CoefC:   &valueProfileTotal,
	CoefD:   &valueProfileTotal,
}, {
	Date:    time.Date(2022, 05, 31, 12, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileTotal,
	CoefB:   &valueProfileTotal,
	CoefC:   &valueProfileTotal,
	CoefD:   &valueProfileTotal,
}, {
	Date:    time.Date(2022, 05, 31, 13, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileTotal,
	CoefB:   &valueProfileTotal,
	CoefC:   &valueProfileTotal,
	CoefD:   &valueProfileTotal,
}, {
	Date:    time.Date(2022, 05, 31, 14, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileTotal,
	CoefB:   &valueProfileTotal,
	CoefC:   &valueProfileTotal,
	CoefD:   &valueProfileTotal,
}}
var Fixture_consume_profiles_total2 = []billing_measures.ConsumProfile{}

var Fixture_consume_profiles_example_total = []billing_measures.ConsumProfile{{
	Date:    time.Date(2022, 05, 18, 10, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileExampleTotal10,
	CoefB:   &valueProfileExampleTotal10,
	CoefC:   &valueProfileExampleTotal10,
	CoefD:   &valueProfileExampleTotal10,
}, {
	Date:    time.Date(2022, 05, 18, 11, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileExampleTotal11,
	CoefB:   &valueProfileExampleTotal11,
	CoefC:   &valueProfileExampleTotal11,
	CoefD:   &valueProfileExampleTotal11,
}, {
	Date:    time.Date(2022, 05, 18, 12, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileExampleTotal12,
	CoefB:   &valueProfileExampleTotal12,
	CoefC:   &valueProfileExampleTotal12,
	CoefD:   &valueProfileExampleTotal12,
}, {
	Date:    time.Date(2022, 05, 18, 13, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileExampleTotal13,
	CoefB:   &valueProfileExampleTotal13,
	CoefC:   &valueProfileExampleTotal13,
	CoefD:   &valueProfileExampleTotal13,
}, {
	Date:    time.Date(2022, 05, 18, 14, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileExampleTotal14,
	CoefB:   &valueProfileExampleTotal14,
	CoefC:   &valueProfileExampleTotal14,
	CoefD:   &valueProfileExampleTotal14,
}, {
	Date:    time.Date(2022, 05, 18, 19, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileExampleTotal19,
	CoefB:   &valueProfileExampleTotal19,
	CoefC:   &valueProfileExampleTotal19,
	CoefD:   &valueProfileExampleTotal19,
}, {
	Date:    time.Date(2022, 05, 18, 20, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileExampleTotal20,
	CoefB:   &valueProfileExampleTotal20,
	CoefC:   &valueProfileExampleTotal20,
	CoefD:   &valueProfileExampleTotal20,
}, {
	Date:    time.Date(2022, 05, 18, 21, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileExampleTotal21,
	CoefB:   &valueProfileExampleTotal21,
	CoefC:   &valueProfileExampleTotal21,
	CoefD:   &valueProfileExampleTotal21,
}, {
	Date:    time.Date(2022, 05, 18, 22, 0, 0, 0, time.UTC),
	Version: 0,
	Type:    "",
	CoefA:   &valueProfileExampleTotal22,
	CoefB:   &valueProfileExampleTotal22,
	CoefC:   &valueProfileExampleTotal22,
	CoefD:   &valueProfileExampleTotal22,
}}

var BillingResult_1_Input_total = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientA,
	BillingBalance: billing_measures.BillingBalance{
		P0: &billing_measures.BillingBalancePeriod{
			AE:              150000,
			EstimatedCodeAE: 1,
		},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			EndDate: time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
	},
}

var BillingResult_1_Output_total = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientA,
	BillingBalance: billing_measures.BillingBalance{
		P0: &billing_measures.BillingBalancePeriod{
			AE:              150000,
			EstimatedCodeAE: 1,
		},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			EndDate:                  time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       30000,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       30000,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       30000,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       30000,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       30000,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
	},
}

var BillingResult_1_Input_total2 = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientA,
	BillingBalance: billing_measures.BillingBalance{
		P0: &billing_measures.BillingBalancePeriod{
			AE:              150000,
			EstimatedCodeAE: 1,
		},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			EndDate: time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},

		{
			EndDate: time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
	},
}

var BillingResult_1_Output_total2 = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientA,
	BillingBalance: billing_measures.BillingBalance{P0: &billing_measures.BillingBalancePeriod{
		AE:              150000,
		EstimatedCodeAE: 1,
	},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			EndDate:                  time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},

		{
			EndDate:                  time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
	},
}

var BillingResult_1_Input_total_example = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientA,
	BillingBalance: billing_measures.BillingBalance{
		P0: &billing_measures.BillingBalancePeriod{
			AE:              42589,
			EstimatedCodeAE: 1,
		},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			EndDate: time.Date(2022, time.May, 18, 10, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 18, 11, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 18, 12, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 18, 13, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 18, 14, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 18, 19, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		}, {
			EndDate: time.Date(2022, time.May, 18, 20, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 18, 21, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
		{
			EndDate: time.Date(2022, time.May, 18, 22, 0, 0, 0, time.UTC),
			Origin:  measures.Filled,
			Period:  measures.P0,
		},
	},
}

var BillingResult_1_Output_total_example = billing_measures.BillingMeasure{
	Coefficient: billing_measures.CoefficientA,
	BillingBalance: billing_measures.BillingBalance{
		P0: &billing_measures.BillingBalancePeriod{
			AE:              42589,
			EstimatedCodeAE: 1,
		},
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{
		{
			EndDate:                  time.Date(2022, time.May, 18, 10, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       5056,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 11, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       5315,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 12, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       5384,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 13, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       5465,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 14, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       5244,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 19, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       4314,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 20, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       4102,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 21, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       3854,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
		{
			EndDate:                  time.Date(2022, time.May, 18, 22, 0, 0, 0, time.UTC),
			Origin:                   measures.Filled,
			AE:                       3855,
			Period:                   measures.P0,
			EstimatedCodeAE:          2,
			EstimatedMethodAE:        billing_measures.ProfileMeasure,
			EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
			MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
		},
	},
}
