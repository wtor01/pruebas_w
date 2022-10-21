package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

var BillingResult_1_Input_Closing_With_Balance = billing_measures.BillingMeasure{
	BillingBalance: billing_measures.BillingBalance{
		EndDate: time.Time{},
		Origin:  measures.STM,
		P0: &billing_measures.BillingBalancePeriod{
			AI: 800000,

			BalanceTypeAI:        "REAL_BY_REMOTE_READ",
			BalanceMeasureTypeAI: "FIRM",
			BalanceOriginAI:      "TLM",
			BalanceValidationAI:  measures.Valid,
			EstimatedCodeAI:      1,
		},
		P6: &billing_measures.BillingBalancePeriod{
			AI: 200000,

			BalanceTypeAI:        "REAL_BY_REMOTE_READ",
			BalanceMeasureTypeAI: "FIRM",
			BalanceOriginAI:      "TLM",
			BalanceValidationAI:  measures.Valid,
			EstimatedCodeAI:      1,
		},
		P5: &billing_measures.BillingBalancePeriod{
			BalanceTypeAI:        "REAL_BY_REMOTE_READ",
			BalanceMeasureTypeAI: "FIRM",
			BalanceOriginAI:      "TLM",
			BalanceValidationAI:  measures.Invalid,
			EstimatedCodeAI:      1,
		},
		P4: &billing_measures.BillingBalancePeriod{

			BalanceTypeAI:        "REAL_BY_REMOTE_READ",
			BalanceMeasureTypeAI: "FIRM",
			BalanceOriginAI:      "TLM",
			BalanceValidationAI:  measures.Invalid,
			EstimatedCodeAI:      1,
		},
	},
	AtrVsCurve: billing_measures.AtrVsCurve{},
	Periods:    []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
	P1Demand:   10,
	P2Demand:   15,
	P3Demand:   20,
	P4Demand:   13,
	P5Demand:   10,
	P6Demand:   6,
}

var BillingResult_1_Output_Closing_With_Balance = billing_measures.BillingMeasure{
	BillingBalance: billing_measures.BillingBalance{
		EndDate: time.Time{},
		Origin:  measures.STM,
		P0: &billing_measures.BillingBalancePeriod{
			AI: 800000,

			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceValidationAI:  measures.Valid,
			EstimatedCodeAI:      1,
		},
		P6: &billing_measures.BillingBalancePeriod{
			AI: 200000,

			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceValidationAI:  measures.Valid,
			EstimatedCodeAI:      1,
		},
		P5: &billing_measures.BillingBalancePeriod{
			AI:                   206897,
			BalanceTypeAI:        billing_measures.EstimatedContractPower,
			BalanceGeneralTypeAI: billing_measures.GeneralEstimated,
			BalanceMeasureTypeAI: billing_measures.ProvisionalBalanceMeasure,
			BalanceOriginAI:      billing_measures.EstimateOrigin,
			BalanceValidationAI:  measures.Invalid,
			EstimatedCodeAI:      6,
		},
		P4: &billing_measures.BillingBalancePeriod{
			AI:                   268966,
			BalanceTypeAI:        billing_measures.EstimatedContractPower,
			BalanceGeneralTypeAI: billing_measures.GeneralEstimated,
			BalanceMeasureTypeAI: billing_measures.ProvisionalBalanceMeasure,
			BalanceOriginAI:      billing_measures.EstimateOrigin,
			BalanceValidationAI:  measures.Invalid,
			EstimatedCodeAI:      6,
		},
	},
	AtrVsCurve: billing_measures.AtrVsCurve{},
	Periods:    []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
	P1Demand:   10,
	P2Demand:   15,
	P3Demand:   20,
	P4Demand:   13,
	P5Demand:   10,
	P6Demand:   6,
}

var BillingResult_2_Input_Closing_With_Balance = billing_measures.BillingMeasure{
	BillingBalance: billing_measures.BillingBalance{
		EndDate: time.Time{},
		Origin:  measures.STM,
		P0: &billing_measures.BillingBalancePeriod{
			AI: 1000000,

			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceValidationAI:  measures.Valid,
			EstimatedCodeAI:      1,
		},
		P6: &billing_measures.BillingBalancePeriod{
			AI: 200000,

			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceValidationAI:  measures.Valid,
			EstimatedCodeAI:      1,
		},
		P5: &billing_measures.BillingBalancePeriod{
			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceValidationAI:  measures.Invalid,
			EstimatedCodeAI:      1,
		},
		P4: &billing_measures.BillingBalancePeriod{
			AI: 500000,

			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceValidationAI:  measures.Valid,
			EstimatedCodeAI:      1,
		},
	},
	AtrVsCurve: billing_measures.AtrVsCurve{},
	Periods:    []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
	P1Demand:   10,
	P2Demand:   15,
	P3Demand:   20,
	P4Demand:   13,
	P5Demand:   10,
	P6Demand:   6,
}

var BillingResult_2_Output_Closing_With_Balance = billing_measures.BillingMeasure{
	BillingBalance: billing_measures.BillingBalance{
		EndDate: time.Time{},
		Origin:  measures.STM,
		P0: &billing_measures.BillingBalancePeriod{
			AI: 1000000,

			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceValidationAI:  measures.Valid,
			EstimatedCodeAI:      1,
		},
		P6: &billing_measures.BillingBalancePeriod{
			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceValidationAI:  measures.Valid,
			EstimatedCodeAI:      1,
			AI:                   200000,
		},
		P5: &billing_measures.BillingBalancePeriod{
			AI:                   103448,
			BalanceTypeAI:        billing_measures.EstimatedContractPower,
			BalanceGeneralTypeAI: billing_measures.GeneralEstimated,
			BalanceMeasureTypeAI: billing_measures.ProvisionalBalanceMeasure,
			BalanceOriginAI:      billing_measures.EstimateOrigin,
			BalanceValidationAI:  measures.Invalid,
			EstimatedCodeAI:      6,
		},
		P4: &billing_measures.BillingBalancePeriod{
			AI:                   500000,
			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceValidationAI:  measures.Valid,
			EstimatedCodeAI:      1,
		},
	},
	AtrVsCurve: billing_measures.AtrVsCurve{},
	Periods:    []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
	P1Demand:   10,
	P2Demand:   15,
	P3Demand:   20,
	P4Demand:   13,
	P5Demand:   10,
	P6Demand:   6,
}

var BillingResult_3_Input_Closing_With_Balance = billing_measures.BillingMeasure{
	BillingBalance: billing_measures.BillingBalance{
		EndDate: time.Time{},
		Origin:  measures.STM,
		P0: &billing_measures.BillingBalancePeriod{
			AE:                  1000000,
			BalanceValidationAE: measures.Valid,
		},
		P6: &billing_measures.BillingBalancePeriod{
			AE: 200000,

			BalanceValidationAE: measures.Valid,
		},
		P5: &billing_measures.BillingBalancePeriod{

			BalanceValidationAE: measures.Invalid,
		},
		P4: &billing_measures.BillingBalancePeriod{
			AE: 500000,

			BalanceValidationAE: measures.Valid,
		},
	},
	AtrVsCurve: billing_measures.AtrVsCurve{},
	Periods:    []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
	P1Demand:   10,
	P2Demand:   15,
	P3Demand:   20,
	P4Demand:   13,
	P5Demand:   10,
	P6Demand:   6,
}

var BillingResult_3_Output_Closing_With_Balance = billing_measures.BillingMeasure{
	BillingBalance: billing_measures.BillingBalance{
		EndDate: time.Time{},
		Origin:  measures.STM,
		P0: &billing_measures.BillingBalancePeriod{
			AE: 1000000,

			BalanceTypeAE:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAE: billing_measures.GeneralReal,
			BalanceMeasureTypeAE: billing_measures.FirmBalanceMeasure,
			BalanceOriginAE:      billing_measures.TlmOrigin,
			BalanceValidationAE:  measures.Valid,
			EstimatedCodeAE:      1,
		},
		P6: &billing_measures.BillingBalancePeriod{
			AE: 200000,

			BalanceTypeAE:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAE: billing_measures.GeneralReal,
			BalanceMeasureTypeAE: billing_measures.FirmBalanceMeasure,
			BalanceOriginAE:      billing_measures.TlmOrigin,
			BalanceValidationAE:  measures.Valid,
			EstimatedCodeAE:      1,
		},
		P5: &billing_measures.BillingBalancePeriod{
			AE:                   0,
			BalanceTypeAE:        billing_measures.EstimatedContractPower,
			BalanceGeneralTypeAE: billing_measures.GeneralEstimated,
			BalanceMeasureTypeAE: billing_measures.ProvisionalBalanceMeasure,
			BalanceOriginAE:      billing_measures.EstimateOrigin,
			BalanceValidationAE:  measures.Invalid,
			EstimatedCodeAE:      6,
		},
		P4: &billing_measures.BillingBalancePeriod{
			AE: 500000,

			BalanceTypeAE:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAE: billing_measures.GeneralReal,
			BalanceMeasureTypeAE: billing_measures.FirmBalanceMeasure,
			BalanceOriginAE:      billing_measures.TlmOrigin,
			BalanceValidationAE:  measures.Valid,
			EstimatedCodeAE:      1,
		},
	},
	AtrVsCurve: billing_measures.AtrVsCurve{},
	Periods:    []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
	P1Demand:   10,
	P2Demand:   15,
	P3Demand:   20,
	P4Demand:   13,
	P5Demand:   10,
	P6Demand:   6,
}
