package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"time"
)

var BillingResult_1_Input_Power_Use_Factor_Balance = billing_measures.BillingMeasure{
	BillingBalance: billing_measures.BillingBalance{
		EndDate: time.Time{},
		Origin:  measures.STM,
		P0: &billing_measures.BillingBalancePeriod{
			AI: 50,
			AE: 50,
			R1: 50,
			R2: 50,
			R3: 50,
			R4: 50,
		},
		P1: &billing_measures.BillingBalancePeriod{
			AI: 50,
			AE: 50,
			R1: 50,
			R2: 50,
			R3: 50,
			R4: 50,
		},
		P2: &billing_measures.BillingBalancePeriod{
			AI: 50,
			AE: 50,
			R1: 50,
			R2: 50,
			R3: 50,
			R4: 50,
		},
		P3: &billing_measures.BillingBalancePeriod{
			AI: 50,
			AE: 50,
			R1: 50,
			R2: 50,
			R3: 50,
			R4: 50,
		},
		P4: &billing_measures.BillingBalancePeriod{
			AI: 50,
			AE: 50,
			R1: 50,
			R2: 50,
			R3: 50,
			R4: 50,
		},
		P5: nil,
		P6: nil,
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC), AI: 30000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: measures.STG, Period: measures.P1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC), AI: 40000, AE: 40000, R1: 40000, R2: 40000, R3: 40000, R4: 40000, Origin: measures.STG, Period: measures.P2}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STG, Period: measures.P3}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STG, Period: measures.P4}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: measures.STG, Period: measures.P2}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 15, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STG, Period: measures.P3}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 16, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STG, Period: measures.P4}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 17, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: measures.STG, Period: measures.P1}},
	AtrVsCurve: billing_measures.AtrVsCurve{},
	Periods:    []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4},
	P1Demand:   10,
	P2Demand:   15,
	P3Demand:   20,
	P4Demand:   30,
	P5Demand:   30,
	P6Demand:   30,
}
var BillingResult_1_Output_Power_Use_Factor_Balance = billing_measures.BillingMeasure{
	BillingBalance: billing_measures.BillingBalance{
		EndDate: time.Time{},
		Origin:  measures.STM,
		P0: &billing_measures.BillingBalancePeriod{
			AI: 200,
			AE: 50,
			R1: 50,
			R2: 50,
			R3: 50,
			R4: 50,

			BalanceTypeAI:        billing_measures.CalculatedByCloseSum,
			BalanceGeneralTypeAI: billing_measures.GeneralCalculated,
			BalanceMeasureTypeAI: billing_measures.ProvisionalBalanceMeasure,
			BalanceOriginAI:      billing_measures.EstimateOrigin,
			EstimatedCodeAI:      6,
		},
		P1: &billing_measures.BillingBalancePeriod{
			AI: 50,
			AE: 50,
			R1: 50,
			R2: 50,
			R3: 50,
			R4: 50,

			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceValidationAI:  "",
			EstimatedCodeAI:      1,
		},
		P2: &billing_measures.BillingBalancePeriod{
			AI: 50,
			AE: 50,
			R1: 50,
			R2: 50,
			R3: 50,
			R4: 50,

			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			EstimatedCodeAI:      1,
		},
		P3: &billing_measures.BillingBalancePeriod{
			AI: 50,
			AE: 50,
			R1: 50,
			R2: 50,
			R3: 50,
			R4: 50,

			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceValidationAI:  "",
			EstimatedCodeAI:      1,
		},
		P4: &billing_measures.BillingBalancePeriod{
			AI:                   50,
			AE:                   50,
			R1:                   50,
			R2:                   50,
			R3:                   50,
			R4:                   50,
			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceValidationAI:  "",
			EstimatedCodeAI:      1,
		},
		P5: nil,
		P6: nil,
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC), AI: 30000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: measures.STG, Period: measures.P1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC), AI: 40000, AE: 40000, R1: 40000, R2: 40000, R3: 40000, R4: 40000, Origin: measures.STG, Period: measures.P2}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STG, Period: measures.P3}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STG, Period: measures.P4}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: measures.STG, Period: measures.P2}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 15, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STG, Period: measures.P3}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 16, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STG, Period: measures.P4}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 17, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: measures.STG, Period: measures.P1}},
	AtrVsCurve: billing_measures.AtrVsCurve{},
	Periods:    []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4},
	P1Demand:   10,
	P2Demand:   15,
	P3Demand:   20,
	P4Demand:   30,
	P5Demand:   30,
	P6Demand:   30,
}

var BillingResult_2_Input_Power_Use_Factor_Balance = billing_measures.BillingMeasure{
	BillingBalance: billing_measures.BillingBalance{
		EndDate: time.Time{},
		Origin:  measures.STM,
		P0:      &billing_measures.BillingBalancePeriod{},
		P1: &billing_measures.BillingBalancePeriod{
			AI: 50,
			AE: 50,
			R1: 50,
			R2: 50,
			R3: 50,
			R4: 50,
		},
		P2: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
		P3: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
		P4: &billing_measures.BillingBalancePeriod{
			AI: 35,
			AE: 35,
			R1: 35,
			R2: 35,
			R3: 35,
			R4: 35,
		},
		P5: nil,
		P6: nil,
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC), AI: 30000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: measures.STG, Period: measures.P1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC), AI: 40000, AE: 40000, R1: 40000, R2: 40000, R3: 40000, R4: 40000, Origin: measures.STG, Period: measures.P2}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STG, Period: measures.P3}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STG, Period: measures.P4}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: measures.STG, Period: measures.P2}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 15, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STG, Period: measures.P3}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 16, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STG, Period: measures.P4}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 17, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: measures.STG, Period: measures.P1}},
	AtrVsCurve: billing_measures.AtrVsCurve{},
	Periods:    []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4},
	P1Demand:   10,
	P2Demand:   15,
	P3Demand:   20,
	P4Demand:   30,
	P5Demand:   30,
	P6Demand:   30,
}

var BillingResult_2_Output_Power_Use_Factor_Balance = billing_measures.BillingMeasure{
	BillingBalance: billing_measures.BillingBalance{
		EndDate: time.Time{},
		Origin:  measures.STM,
		P0: &billing_measures.BillingBalancePeriod{AE: 0,
			AI:                   108,
			R1:                   0,
			R2:                   0,
			R3:                   0,
			R4:                   0,
			BalanceTypeAI:        billing_measures.CalculatedByCloseSum,
			BalanceGeneralTypeAI: billing_measures.GeneralCalculated,
			BalanceMeasureTypeAI: billing_measures.ProvisionalBalanceMeasure,
			BalanceOriginAI:      billing_measures.EstimateOrigin,
			BalanceValidationAI:  "",
			EstimatedCodeAI:      6,
		},
		P1: &billing_measures.BillingBalancePeriod{
			AI:                   50,
			AE:                   50,
			R1:                   50,
			R2:                   50,
			R3:                   50,
			R4:                   50,
			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceValidationAI:  "",
			EstimatedCodeAI:      1,
		},
		P2: &billing_measures.BillingBalancePeriod{
			AE:                   0,
			AI:                   10,
			R1:                   0,
			R2:                   0,
			R3:                   0,
			R4:                   0,
			BalanceTypeAI:        billing_measures.PowerUseFactor,
			BalanceGeneralTypeAI: billing_measures.GeneralEstimated,
			BalanceMeasureTypeAI: billing_measures.ProvisionalBalanceMeasure,
			BalanceOriginAI:      billing_measures.EstimateOrigin,
			EstimatedCodeAI:      6,
			BalanceValidationAI:  measures.Invalid,
		},
		P3: &billing_measures.BillingBalancePeriod{
			AE:                   0,
			AI:                   13,
			R1:                   0,
			R2:                   0,
			R3:                   0,
			R4:                   0,
			BalanceTypeAI:        billing_measures.PowerUseFactor,
			BalanceGeneralTypeAI: billing_measures.GeneralEstimated,
			BalanceMeasureTypeAI: billing_measures.ProvisionalBalanceMeasure,
			BalanceOriginAI:      billing_measures.EstimateOrigin,
			BalanceValidationAI:  measures.Invalid,
			EstimatedCodeAI:      6,
		},
		P4: &billing_measures.BillingBalancePeriod{
			AI:                   35,
			AE:                   35,
			R1:                   35,
			R2:                   35,
			R3:                   35,
			R4:                   35,
			BalanceTypeAI:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAI: billing_measures.GeneralReal,
			BalanceOriginAI:      billing_measures.TlmOrigin,
			BalanceMeasureTypeAI: billing_measures.FirmBalanceMeasure,
			EstimatedCodeAI:      1,
		},
		P5: nil,
		P6: nil,
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC), AI: 30000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: measures.STG, Period: measures.P1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC), AI: 40000, AE: 40000, R1: 40000, R2: 40000, R3: 40000, R4: 40000, Origin: measures.STG, Period: measures.P2}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STG, Period: measures.P3}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STG, Period: measures.P4}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: measures.STG, Period: measures.P2}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 15, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STG, Period: measures.P3}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 16, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STG, Period: measures.P4}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 17, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: measures.STG, Period: measures.P1}},
	AtrVsCurve:   billing_measures.AtrVsCurve{},
	Periods:      []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4},
	TariffID:     "",
	CalendarCode: "",
	P1Demand:     10,
	P2Demand:     15,
	P3Demand:     20,
	P4Demand:     30,
	P5Demand:     30,
	P6Demand:     30,
}

var BillingResult_3_Input_Power_Use_Factor_Balance = billing_measures.BillingMeasure{
	BillingBalance: billing_measures.BillingBalance{
		EndDate: time.Time{},
		Origin:  measures.STM,
		P0: &billing_measures.BillingBalancePeriod{
			BalanceValidationAE: measures.Valid,
		},
		P1: &billing_measures.BillingBalancePeriod{
			AI:                  50,
			AE:                  50,
			R1:                  50,
			R2:                  50,
			R3:                  50,
			R4:                  50,
			BalanceValidationAE: measures.Valid,
		},
		P2: &billing_measures.BillingBalancePeriod{BalanceValidationAE: measures.Invalid},
		P3: &billing_measures.BillingBalancePeriod{BalanceValidationAE: measures.Invalid},
		P4: &billing_measures.BillingBalancePeriod{
			AI:                  35,
			AE:                  35,
			R1:                  35,
			R2:                  35,
			R3:                  35,
			R4:                  35,
			BalanceValidationAE: measures.Valid,
		},
		P5: nil,
		P6: nil,
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC), AI: 30000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: measures.STG, Period: measures.P1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC), AI: 40000, AE: 40000, R1: 40000, R2: 40000, R3: 40000, R4: 40000, Origin: measures.STG, Period: measures.P2}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STG, Period: measures.P3}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STG, Period: measures.P4}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: measures.STG, Period: measures.P2}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 15, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STG, Period: measures.P3}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 16, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STG, Period: measures.P4}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 17, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: measures.STG, Period: measures.P1}},
	AtrVsCurve: billing_measures.AtrVsCurve{},
	Periods:    []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4},
	P1Demand:   10,
	P2Demand:   15,
	P3Demand:   20,
	P4Demand:   30,
	P5Demand:   30,
	P6Demand:   30,
}

var BillingResult_3_Output_Power_Use_Factor_Balance = billing_measures.BillingMeasure{
	BillingBalance: billing_measures.BillingBalance{
		EndDate: time.Time{},
		Origin:  measures.STM,
		P0: &billing_measures.BillingBalancePeriod{
			AE:                   0,
			AI:                   0,
			R1:                   0,
			R2:                   0,
			R3:                   0,
			R4:                   0,
			BalanceTypeAE:        billing_measures.CalculatedByCloseSum,
			BalanceGeneralTypeAE: billing_measures.GeneralCalculated,
			BalanceMeasureTypeAE: billing_measures.ProvisionalBalanceMeasure,
			BalanceOriginAE:      billing_measures.EstimateOrigin,
			EstimatedCodeAE:      6,
			BalanceValidationAE:  measures.Valid,
		},
		P1: &billing_measures.BillingBalancePeriod{
			AI:                   50,
			AE:                   50,
			R1:                   50,
			R2:                   50,
			R3:                   50,
			R4:                   50,
			BalanceTypeAE:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAE: billing_measures.GeneralReal,
			BalanceMeasureTypeAE: billing_measures.FirmBalanceMeasure,
			BalanceOriginAE:      billing_measures.TlmOrigin,
			BalanceValidationAE:  measures.Valid,
			EstimatedCodeAE:      1,
		},
		P2: &billing_measures.BillingBalancePeriod{
			AE:                   0,
			AI:                   0,
			R1:                   0,
			R2:                   0,
			R3:                   0,
			R4:                   0,
			BalanceTypeAE:        billing_measures.PowerUseFactor,
			BalanceGeneralTypeAE: billing_measures.GeneralEstimated,
			BalanceMeasureTypeAE: billing_measures.ProvisionalBalanceMeasure,
			BalanceOriginAE:      billing_measures.EstimateOrigin,
			EstimatedCodeAE:      6,
			BalanceValidationAE:  measures.Invalid,
		},
		P3: &billing_measures.BillingBalancePeriod{
			AE:                   0,
			AI:                   0,
			R1:                   0,
			R2:                   0,
			R3:                   0,
			R4:                   0,
			BalanceTypeAE:        billing_measures.PowerUseFactor,
			BalanceGeneralTypeAE: billing_measures.GeneralEstimated,
			BalanceMeasureTypeAE: billing_measures.ProvisionalBalanceMeasure,
			BalanceOriginAE:      billing_measures.EstimateOrigin,
			EstimatedCodeAE:      6,
			BalanceValidationAE:  measures.Invalid,
		},
		P4: &billing_measures.BillingBalancePeriod{
			AI:                   35,
			AE:                   35,
			R1:                   35,
			R2:                   35,
			R3:                   35,
			R4:                   35,
			BalanceTypeAE:        billing_measures.RealByRemoteRead,
			BalanceGeneralTypeAE: billing_measures.GeneralReal,
			BalanceOriginAE:      billing_measures.TlmOrigin,
			BalanceMeasureTypeAE: billing_measures.FirmBalanceMeasure,
			EstimatedCodeAE:      1,
			BalanceValidationAE:  measures.Valid,
		},
		P5: nil,
		P6: nil,
	},
	BillingLoadCurve: []billing_measures.BillingLoadCurve{billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 10, 0, 0, 0, time.UTC), AI: 30000, AE: 30000, R1: 30000, R2: 30000, R3: 30000, R4: 30000, Origin: measures.STG, Period: measures.P1}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 11, 0, 0, 0, time.UTC), AI: 40000, AE: 40000, R1: 40000, R2: 40000, R3: 40000, R4: 40000, Origin: measures.STG, Period: measures.P2}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 12, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STG, Period: measures.P3}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 13, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STG, Period: measures.P4}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 14, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: measures.STG, Period: measures.P2}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 15, 0, 0, 0, time.UTC), AI: 50000, AE: 50000, R1: 50000, R2: 50000, R3: 50000, R4: 50000, Origin: measures.STG, Period: measures.P3}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 16, 0, 0, 0, time.UTC), AI: 60000, AE: 60000, R1: 60000, R2: 60000, R3: 60000, R4: 60000, Origin: measures.STG, Period: measures.P4}, billing_measures.BillingLoadCurve{
		EndDate: time.Date(2022, time.May, 31, 17, 0, 0, 0, time.UTC), AI: 72000, AE: 72000, R1: 72000, R2: 72000, R3: 72000, R4: 72000, Origin: measures.STG, Period: measures.P1}},
	AtrVsCurve:   billing_measures.AtrVsCurve{},
	Periods:      []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4},
	TariffID:     "",
	CalendarCode: "",
	P1Demand:     10,
	P2Demand:     15,
	P3Demand:     20,
	P4Demand:     30,
	P5Demand:     30,
	P6Demand:     30,
}
