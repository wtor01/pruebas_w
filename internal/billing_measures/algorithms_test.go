package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Unit_Domain_Algorithms_CchPartialAdjustment(t *testing.T) {
	type input struct {
		b         *BillingMeasure
		period    measures.PeriodKey
		magnitude measures.Magnitude
	}

	type output struct {
		b         []BillingLoadCurve
		outputErr error
	}

	testCases := map[string]struct {
		input  input
		output output
	}{
		"Adjust for period 1": {
			input: input{
				b: &BillingMeasure{
					BillingBalance: BillingBalance{
						P1: &BillingBalancePeriod{
							AI: 10000,
							AE: 20000,
							R1: 30000,
							R2: 40000,
							R3: 50000,
							R4: 60000,
						},
					},
					BillingLoadCurve: []BillingLoadCurve{
						{
							AI:     10000,
							AE:     20000,
							R1:     30000,
							R2:     40000,
							R3:     50000,
							R4:     60000,
							Period: measures.P1,
						},
						{
							AI:     10000,
							AE:     10000,
							R1:     10000,
							R2:     10000,
							R3:     10000,
							R4:     10000,
							Period: measures.P1,
						},
						{
							AI:     444444444,
							AE:     444444444,
							R1:     444444444,
							R2:     444444444,
							R3:     444444444,
							R4:     444444444,
							Period: measures.P2,
						},
					},
				},
				period:    measures.P1,
				magnitude: measures.AE,
			},
			output: output{
				b: []BillingLoadCurve{
					{
						AI:                       10000,
						AE:                       13333,
						R1:                       30000,
						R2:                       40000,
						R3:                       50000,
						R4:                       60000,
						Period:                   measures.P1,
						MeasureTypeAE:            FirmBalanceMeasure,
						EstimatedCodeAE:          3,
						EstimatedMethodAE:        Adjustment,
						EstimatedGeneralMethodAE: GeneralAdjusted,
					},
					{
						AI:                       10000,
						AE:                       6667,
						R1:                       10000,
						R2:                       10000,
						R3:                       10000,
						R4:                       10000,
						Period:                   measures.P1,
						MeasureTypeAE:            FirmBalanceMeasure,
						EstimatedCodeAE:          3,
						EstimatedMethodAE:        Adjustment,
						EstimatedGeneralMethodAE: GeneralAdjusted,
					},
					{
						AI:     444444444,
						AE:     444444444,
						R1:     444444444,
						R2:     444444444,
						R3:     444444444,
						R4:     444444444,
						Period: measures.P2,
					},
				},
				outputErr: nil,
			},
		},

		"Adjust for period 2 and balance code is 2": {
			input: input{
				b: &BillingMeasure{
					BillingBalance: BillingBalance{
						P2: &BillingBalancePeriod{
							AI:              10000,
							AE:              20000,
							R1:              30000,
							R2:              40000,
							R3:              50000,
							R4:              60000,
							EstimatedCodeAE: 4,
						},
					},
					BillingLoadCurve: []BillingLoadCurve{
						{
							AI:     10000,
							AE:     20000,
							R1:     30000,
							R2:     40000,
							R3:     50000,
							R4:     60000,
							Period: measures.P2,
						},
						{
							AI:     10000,
							AE:     20000,
							R1:     30000,
							R2:     40000,
							R3:     50000,
							R4:     60000,
							Period: measures.P2,
						},
					},
				},
				period:    measures.P2,
				magnitude: measures.AE,
			},
			output: output{
				b: []BillingLoadCurve{
					{
						AI:                       10000,
						AE:                       10000,
						R1:                       30000,
						R2:                       40000,
						R3:                       50000,
						R4:                       60000,
						Period:                   measures.P2,
						EstimatedCodeAE:          3,
						EstimatedMethodAE:        Adjustment,
						MeasureTypeAE:            ProvisionalBalanceMeasure,
						EstimatedGeneralMethodAE: GeneralAdjusted,
					},
					{
						AI:                       10000,
						AE:                       10000,
						R1:                       30000,
						R2:                       40000,
						R3:                       50000,
						R4:                       60000,
						Period:                   measures.P2,
						EstimatedCodeAE:          3,
						EstimatedMethodAE:        Adjustment,
						EstimatedGeneralMethodAE: GeneralAdjusted,
						MeasureTypeAE:            ProvisionalBalanceMeasure,
					},
				},
				outputErr: nil,
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			cchPartialAdjustment := NewCchPartialAdjustment(test.input.b, test.input.period, test.input.magnitude)
			ExecuteResponse := cchPartialAdjustment.Execute(context.Background())
			cchId := cchPartialAdjustment.ID()
			assert.Equal(t, test.output.outputErr, ExecuteResponse)
			assert.Equal(t, cchId, "CCH_PARTIAL_ADJUST")
			assert.Equal(t, test.output.b, cchPartialAdjustment.b.BillingLoadCurve)
		})
	}

}

func Test_Unit_Domain_Algorithms_CchCompleted(t *testing.T) {
	type input struct {
		b         *BillingMeasure
		period    measures.PeriodKey
		magnitude measures.Magnitude
	}

	type output struct {
		b []BillingLoadCurve
	}

	tests := map[string]struct {
		input  input
		output output
	}{
		"Execute for period measures.P1 and magnitude measures.AI": {
			input: input{
				b: &BillingMeasure{
					BillingLoadCurve: []BillingLoadCurve{
						{
							Period:                   measures.P1,
							EstimatedCodeAE:          2,
							EstimatedMethodAE:        EstimateBalance,
							EstimatedGeneralMethodAE: GeneralEstimated,
							MeasureTypeAE:            ProvisionalBalanceMeasure,
						},
						{
							Period: measures.P1,
						},
						{
							Period: measures.P1,
						},
						{
							Period: measures.P2,
						},
					},
				},
				period:    measures.P1,
				magnitude: measures.AI,
			},
			output: output{
				b: []BillingLoadCurve{
					{
						Period:                   measures.P1,
						EstimatedCodeAI:          1,
						EstimatedMethodAI:        RealBalance,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
						EstimatedCodeAE:          2,
						EstimatedMethodAE:        EstimateBalance,
						EstimatedGeneralMethodAE: GeneralEstimated,
						MeasureTypeAE:            ProvisionalBalanceMeasure,
					},
					{
						Period:                   measures.P1,
						EstimatedCodeAI:          1,
						EstimatedMethodAI:        RealBalance,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
					{
						Period:                   measures.P1,
						EstimatedCodeAI:          1,
						EstimatedMethodAI:        RealBalance,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
					{
						Period: measures.P2,
					},
				},
			},
		},
	}

	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			cchCompleted := NewCchCompleted(testCase.input.b, testCase.input.period, testCase.input.magnitude)
			cchID := cchCompleted.ID()
			cchCompleted.Execute(ctx)
			assert.Equal(t, testCase.output.b, cchCompleted.b.BillingLoadCurve)
			assert.Equal(t, "CCH_COMPLETED", cchID)
		})
	}
}

func Test_Unit_Domain_Algorithms_BalanceCalculatedByCch(t *testing.T) {
	type input struct {
		b         *BillingMeasure
		period    measures.PeriodKey
		magnitude measures.Magnitude
	}

	type output struct {
		b         BillingBalance
		outputErr error
	}

	tests := map[string]struct {
		input  input
		output output
	}{
		"Should Return Error If PeriodKey is not Valid": {
			input: input{
				b:      &BillingMeasure{},
				period: measures.PeriodKey("Random key"),
			},
			output: output{
				outputErr: ErrorPeriodKeyDoesNotMatch,
			},
		},
		"Should work for period 1": {
			input: input{
				b: &BillingMeasure{
					BillingBalance: BillingBalance{
						P0: &BillingBalancePeriod{},
						P1: &BillingBalancePeriod{
							AI: 50,
						},
					},
					BillingLoadCurve: []BillingLoadCurve{
						{
							AI:     10,
							AE:     20,
							R1:     30,
							R2:     40,
							R3:     50,
							R4:     60,
							Period: measures.P1,
						},
						{
							AI:     10,
							AE:     20,
							R1:     30,
							R2:     40,
							R3:     50,
							R4:     60,
							Period: measures.P1,
						},
						{
							AI:     10,
							AE:     20,
							R1:     30,
							R2:     40,
							R3:     50,
							R4:     60,
							Period: measures.P1,
						},
						{
							AI:     10,
							AE:     20,
							R1:     30,
							R2:     40,
							R3:     50,
							R4:     60,
							Period: measures.P1,
						},
						{
							AI:     10,
							AE:     20,
							R1:     30,
							R2:     40,
							R3:     50,
							R4:     60,
							Period: measures.P2,
						},
					},
				},
				period:    measures.P1,
				magnitude: measures.AE,
			},
			output: output{
				b: BillingBalance{
					P0: &BillingBalancePeriod{BalanceGeneralTypeAE: GeneralCalculated},
					P1: &BillingBalancePeriod{
						AE:                   80,
						AI:                   50,
						BalanceTypeAE:        CalculatedBalance,
						BalanceGeneralTypeAE: GeneralCalculated,
						BalanceMeasureTypeAE: FirmBalanceMeasure,
						BalanceOriginAE:      TlgOrigin,
						EstimatedCodeAE:      2,
					},
				},
				outputErr: nil,
			},
		},
	}

	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			cchCompletedEx := NewBalanceCalculatedByCch(testCase.input.b, testCase.input.period, testCase.input.magnitude)
			cchID := cchCompletedEx.ID()
			executeResponse := cchCompletedEx.Execute(context.Background())
			assert.Equal(t, testCase.output.outputErr, executeResponse)
			assert.Equal(t, testCase.output.b, cchCompletedEx.b.BillingBalance)
			assert.Equal(t, "BALANCE_CALCULATED_BY_CCH", cchID)
		})
	}

}

func Test_Unit_Domain_Algorithms_CchCompleteAdjustment(t *testing.T) {
	type input struct {
		b         *BillingMeasure
		period    measures.PeriodKey
		magnitude measures.Magnitude
	}

	type output struct {
		b         []BillingLoadCurve
		outputErr error
	}

	testCases := map[string]struct {
		input  input
		output output
	}{
		"Adjust for period 1": {
			input: input{
				b: &BillingMeasure{
					BillingBalance: BillingBalance{
						P1: &BillingBalancePeriod{
							AI:              10000,
							AE:              20000,
							R1:              30000,
							R2:              40000,
							R3:              50000,
							R4:              60000,
							EstimatedCodeAE: 1,
						},
						P2: &BillingBalancePeriod{
							AI: 44444,
							AE: 44444,
							R1: 44444,
							R2: 44444,
							R3: 44444,
							R4: 44444,
						},
					},
					BillingLoadCurve: []BillingLoadCurve{
						{
							AI:     10000,
							AE:     20000,
							R1:     30000,
							R2:     40000,
							R3:     50000,
							R4:     60000,
							Period: measures.P1,
						},
						{
							AI:     10000,
							AE:     10000,
							R1:     10000,
							R2:     10000,
							R3:     10000,
							R4:     10000,
							Period: measures.P1,
						},
						{
							AI:     444444444,
							AE:     444444444,
							R1:     444444444,
							R2:     444444444,
							R3:     444444444,
							R4:     444444444,
							Period: measures.P2,
						},
					},
				},
				period:    measures.P1,
				magnitude: measures.AE,
			},
			output: output{
				b: []BillingLoadCurve{
					{
						AI:                       10000,
						AE:                       13333,
						R1:                       30000,
						R2:                       40000,
						R3:                       50000,
						R4:                       60000,
						Period:                   measures.P1,
						MeasureTypeAE:            FirmBalanceMeasure,
						EstimatedCodeAE:          3,
						EstimatedMethodAE:        Adjustment,
						EstimatedGeneralMethodAE: GeneralAdjusted,
					},
					{
						AI:                       10000,
						AE:                       6667,
						R1:                       10000,
						R2:                       10000,
						R3:                       10000,
						R4:                       10000,
						Period:                   measures.P1,
						MeasureTypeAE:            FirmBalanceMeasure,
						EstimatedCodeAE:          3,
						EstimatedMethodAE:        Adjustment,
						EstimatedGeneralMethodAE: GeneralAdjusted,
					},
					{
						AI:     444444444,
						AE:     444444444,
						R1:     444444444,
						R2:     444444444,
						R3:     444444444,
						R4:     444444444,
						Period: measures.P2,
					},
				},
				outputErr: nil,
			},
		},

		"Adjust for period 2 and balance code is 2": {
			input: input{
				b: &BillingMeasure{
					BillingBalance: BillingBalance{
						P2: &BillingBalancePeriod{
							AI:              10000,
							AE:              20000,
							R1:              30000,
							R2:              40000,
							R3:              50000,
							R4:              60000,
							EstimatedCodeAE: 4,
						},
					},
					BillingLoadCurve: []BillingLoadCurve{
						{
							AI:     10000,
							AE:     20000,
							R1:     30000,
							R2:     40000,
							R3:     50000,
							R4:     60000,
							Period: measures.P2,
						},
						{
							AI:     10000,
							AE:     20000,
							R1:     30000,
							R2:     40000,
							R3:     50000,
							R4:     60000,
							Period: measures.P2,
						},
					},
				},
				period:    measures.P2,
				magnitude: measures.AE,
			},
			output: output{
				b: []BillingLoadCurve{
					{
						AI:                       10000,
						AE:                       10000,
						R1:                       30000,
						R2:                       40000,
						R3:                       50000,
						R4:                       60000,
						Period:                   measures.P2,
						EstimatedCodeAE:          3,
						EstimatedMethodAE:        Adjustment,
						MeasureTypeAE:            ProvisionalBalanceMeasure,
						EstimatedGeneralMethodAE: GeneralAdjusted,
					},
					{
						AI:                       10000,
						AE:                       10000,
						R1:                       30000,
						R2:                       40000,
						R3:                       50000,
						R4:                       60000,
						Period:                   measures.P2,
						EstimatedCodeAE:          3,
						EstimatedMethodAE:        Adjustment,
						EstimatedGeneralMethodAE: GeneralAdjusted,
						MeasureTypeAE:            ProvisionalBalanceMeasure,
					},
				},

				outputErr: nil,
			},
		},

		"Should Return Invalid Key Period": {
			input: input{
				b:      &BillingMeasure{},
				period: measures.PeriodKey("Random Key"),
			},
			output: output{
				outputErr: ErrorPeriodKeyDoesNotMatch,
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			cchPartialCompleteAdjustment := NewCchCompleteAdjustment(test.input.b, test.input.period, test.input.magnitude)
			ExecuteResponse := cchPartialCompleteAdjustment.Execute(context.Background())
			cchId := cchPartialCompleteAdjustment.ID()
			assert.Equal(t, test.output.outputErr, ExecuteResponse)
			assert.Equal(t, cchId, "CCH_COMPLETE_ADJUST")
			assert.Equal(t, cchPartialCompleteAdjustment.b.BillingLoadCurve, test.output.b)
		})
	}

}

func Test_Unit_Domain_Algorithms_Balance_Completed_Execute(t *testing.T) {

	type input struct {
		getBillingMeasure func() BillingMeasure
		period            measures.PeriodKey
		magnitude         measures.Magnitude
	}

	type want struct {
		billingBalance BillingBalance
		err            error
	}

	tests := map[string]struct {
		input input
		want  want
	}{
		"should change to BillingBalance for P0 in case ActualReadingClosure.Origin = measures.STM": {
			input: input{
				getBillingMeasure: func() BillingMeasure {
					b := NewBillingMeasure(
						"ES0130000000357054DJ",
						time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
						time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
						"0130",
						"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
						[]measures.PeriodKey{},
						[]measures.Magnitude{},
						measures.TLG,
					)
					b.ActualReadingClosure.Origin = measures.STM
					b.BillingBalance.P0.BalanceValidationAI = measures.Valid

					return b
				},
				period:    measures.P0,
				magnitude: measures.AI,
			},
			want: want{
				billingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						BalanceTypeAI:        RealBalance,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
						BalanceOriginAI:      RemoteBalanceOrigin,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      1,
					},
				},
				err: nil,
			},
		},
		"should change to BillingBalance for P1 in case ActualReadingClosure.Origin = measures.STG": {
			input: input{
				getBillingMeasure: func() BillingMeasure {
					b := NewBillingMeasure(
						"ES0130000000357054DJ",
						time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
						time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
						"0130",
						"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
						[]measures.PeriodKey{measures.P1},
						[]measures.Magnitude{},
						measures.TLG,
					)
					b.ActualReadingClosure.Origin = measures.STG
					b.BillingBalance.P0 = nil
					b.BillingBalance.P1.BalanceValidationAI = measures.Valid

					return b
				},
				period:    measures.P1,
				magnitude: measures.AI,
			},
			want: want{
				billingBalance: BillingBalance{
					P1: &BillingBalancePeriod{
						BalanceTypeAI:        RealBalance,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
						BalanceOriginAI:      RemoteBalanceOrigin,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      1,
					},
				},
				err: nil,
			},
		},
		"should change to BillingBalance for P2 in case ActualReadingClosure.Origin = measures.TPL": {
			input: input{
				getBillingMeasure: func() BillingMeasure {
					b := NewBillingMeasure(
						"ES0130000000357054DJ",
						time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
						time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
						"0130",
						"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
						[]measures.PeriodKey{measures.P2},
						[]measures.Magnitude{},
						measures.TLG,
					)
					b.ActualReadingClosure.Origin = measures.TPL
					b.BillingBalance.P0 = nil
					b.BillingBalance.P2.BalanceValidationAI = measures.Valid

					return b
				},
				period:    measures.P2,
				magnitude: measures.AI,
			},
			want: want{
				billingBalance: BillingBalance{
					P2: &BillingBalancePeriod{
						BalanceTypeAI:        RealBalance,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
						BalanceOriginAI:      RemoteBalanceOrigin,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      3,
					},
				},
				err: nil,
			},
		},
		"should change to BillingBalance for P3 in case ActualReadingClosure.Origin = measures.Auto": {
			input: input{
				getBillingMeasure: func() BillingMeasure {
					b := NewBillingMeasure(
						"ES0130000000357054DJ",
						time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
						time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
						"0130",
						"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
						[]measures.PeriodKey{measures.P3},
						[]measures.Magnitude{},
						measures.TLG,
					)
					b.ActualReadingClosure.Origin = measures.Auto
					b.BillingBalance.P0 = nil
					b.BillingBalance.P3.BalanceValidationAI = measures.Valid

					return b
				},
				period:    measures.P3,
				magnitude: measures.AI,
			},
			want: want{
				billingBalance: BillingBalance{
					P3: &BillingBalancePeriod{
						BalanceTypeAI:        EstimateBalance,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
						BalanceOriginAI:      AutoBalanceOrigin,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      4,
					},
				},
				err: nil,
			},
		},
		"should change to BillingBalance for P4 in case ActualReadingClosure.Origin = measures.Manual": {
			input: input{
				getBillingMeasure: func() BillingMeasure {
					b := NewBillingMeasure(
						"ES0130000000357054DJ",
						time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
						time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
						"0130",
						"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
						[]measures.PeriodKey{measures.P4},
						[]measures.Magnitude{},
						measures.TLG,
					)
					b.ActualReadingClosure.Origin = measures.Manual
					b.BillingBalance.P0 = nil
					b.BillingBalance.P4.BalanceValidationAI = measures.Valid

					return b
				},
				period:    measures.P4,
				magnitude: measures.AI,
			},
			want: want{
				billingBalance: BillingBalance{
					P4: &BillingBalancePeriod{
						BalanceTypeAI:        RealBalance,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
						BalanceOriginAI:      RemoteBalanceOrigin,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      3,
					},
				},
				err: nil,
			},
		},
		"should change to BillingBalance for P5 in case ActualReadingClosure.Origin = measures.Visual": {
			input: input{
				getBillingMeasure: func() BillingMeasure {
					b := NewBillingMeasure(
						"ES0130000000357054DJ",
						time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
						time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
						"0130",
						"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
						[]measures.PeriodKey{measures.P5},
						[]measures.Magnitude{},
						measures.TLG,
					)
					b.ActualReadingClosure.Origin = measures.Visual
					b.BillingBalance.P0 = nil
					b.BillingBalance.P5.BalanceValidationAI = measures.Valid

					return b
				},
				period:    measures.P5,
				magnitude: measures.AI,
			},
			want: want{
				billingBalance: BillingBalance{
					P5: &BillingBalancePeriod{
						BalanceTypeAI:        RealBalance,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
						BalanceOriginAI:      RemoteBalanceOrigin,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      3,
					},
				},
				err: nil,
			},
		},
		"should change to BillingBalance for P6 in case ActualReadingClosure.Origin = measures.STM": {
			input: input{
				getBillingMeasure: func() BillingMeasure {
					b := NewBillingMeasure(
						"ES0130000000357054DJ",
						time.Date(2022, 5, 1, 0, 0, 0, 0, loc).UTC(),
						time.Date(2022, 6, 1, 0, 0, 0, 0, loc).UTC(),
						"0130",
						"5140cda8-1daa-4f52-b85d-bf29d45ca62a",
						[]measures.PeriodKey{measures.P6},
						[]measures.Magnitude{},
						measures.TLG,
					)
					b.ActualReadingClosure.Origin = measures.STM
					b.BillingBalance.P0 = nil
					b.BillingBalance.P6.BalanceValidationAI = measures.Valid

					return b
				},
				period:    measures.P6,
				magnitude: measures.AI,
			},
			want: want{
				billingBalance: BillingBalance{
					P6: &BillingBalancePeriod{
						BalanceTypeAI:        RealBalance,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
						BalanceOriginAI:      RemoteBalanceOrigin,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      1,
					},
				},
				err: nil,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			billingMeasure := testCase.input.getBillingMeasure()
			algorithm := NewBalanceCompleted(&billingMeasure, testCase.input.period, testCase.input.magnitude)
			assert.Equal(t, testCase.want.err, algorithm.Execute(ctx), testName)
			assert.Equal(t, testCase.want.billingBalance, billingMeasure.BillingBalance, testName)
		})
	}

}

func Test_Unit_Domain_Algorithms_EstimatedHistoryTlg(t *testing.T) {
	type input struct {
		b          *BillingMeasure
		period     measures.PeriodKey
		magnitude  measures.Magnitude
		contextTlg GraphContext
	}

	type output struct {
		billingBalance BillingBalance
		id             string
	}

	testCases := map[string]struct {
		input  input
		output output
		want   error
	}{
		"Should be calculate balance correct (Period Hours Are Equal)": {
			input: input{
				period:    measures.P0,
				magnitude: measures.AI,
				contextTlg: GraphContext{
					LastHistory: BillingMeasure{
						RegisterType: measures.Hourly,
						BillingBalance: BillingBalance{
							P0: &BillingBalancePeriod{
								AI: 20410,
							},
						},
						BillingLoadCurve: []BillingLoadCurve{
							{
								AI:     10000,
								Period: measures.P0,
							},
							{
								AI:     10000,
								Period: measures.P0,
							},
							{
								AI:     10000,
								Period: measures.P0,
							},
							{
								AI:     10000,
								Period: measures.P0,
							},
							{
								AI:     10000,
								Period: measures.P0,
							},
						},
					},
				},
				b: &BillingMeasure{
					RegisterType: measures.Hourly,
					BillingBalance: BillingBalance{
						P0: &BillingBalancePeriod{},
					},
					BillingLoadCurve: []BillingLoadCurve{
						{
							AI:     10000,
							Period: measures.P0,
						},
						{
							AI:     10000,
							Period: measures.P0,
						},
						{
							AI:     10000,
							Period: measures.P0,
						},
						{
							AI:     10000,
							Period: measures.P0,
						},
						{
							AI:     10000,
							Period: measures.P0,
						},
					},
				},
			},
			output: output{
				billingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI:                   20410,
						BalanceTypeAI:        EstimateHistoryProfile,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
						BalanceOriginAI:      EstimateOrigin,
						EstimatedCodeAI:      5,
					},
				},
				id: "BALANCE_ESTIMATED_HISTORY_TLG",
			},
			want: nil,
		},
		"Should be calculate balance correct (Period Hours Not Equal)": {
			input: input{
				period:    measures.P0,
				magnitude: measures.AE,
				contextTlg: GraphContext{
					LastHistory: BillingMeasure{
						RegisterType: measures.Hourly,
						BillingBalance: BillingBalance{
							P0: &BillingBalancePeriod{
								AE: 23100,
							},
						},
						BillingLoadCurve: []BillingLoadCurve{
							{
								AE:     10000,
								Period: measures.P0,
							},
							{
								AE:     10000,
								Period: measures.P0,
							},
							{
								AE:     10000,
								Period: measures.P0,
							},
							{
								AE:     10000,
								Period: measures.P0,
							},
							{
								AE:     10000,
								Period: measures.P0,
							},
						},
					},
				},
				b: &BillingMeasure{
					RegisterType: measures.Hourly,
					BillingBalance: BillingBalance{
						P0: &BillingBalancePeriod{},
					},
					BillingLoadCurve: []BillingLoadCurve{
						{
							AE:     10000,
							Period: measures.P0,
						},
						{
							AE:     10000,
							Period: measures.P0,
						},
						{
							AE:     10000,
							Period: measures.P0,
						},
					},
				},
			},
			output: output{
				billingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AE:                   13860,
						BalanceTypeAE:        EstimateHistoryProfile,
						BalanceGeneralTypeAE: GeneralEstimated,
						BalanceMeasureTypeAE: ProvisionalBalanceMeasure,
						BalanceOriginAE:      EstimateOrigin,
						EstimatedCodeAE:      5,
					},
				},
				id: "BALANCE_ESTIMATED_HISTORY_TLG",
			},
			want: nil,
		},
		"Should be calculate balance correct (Period in QuarterHour Hours Are Equal)": {
			input: input{
				period:    measures.P0,
				magnitude: measures.R1,
				contextTlg: GraphContext{
					LastHistory: BillingMeasure{
						RegisterType: measures.QuarterHour,
						BillingBalance: BillingBalance{
							P0: &BillingBalancePeriod{
								R1: 21100,
							},
						},
						BillingLoadCurve: []BillingLoadCurve{
							{
								R1:     10000,
								Period: measures.P0,
							},
							{
								R1:     10000,
								Period: measures.P0,
							},
							{
								R1:     10000,
								Period: measures.P0,
							},
							{
								R1:     10000,
								Period: measures.P0,
							},
							{
								R1:     10000,
								Period: measures.P0,
							},
						},
					},
				},
				b: &BillingMeasure{
					RegisterType: measures.QuarterHour,
					BillingBalance: BillingBalance{
						P0: &BillingBalancePeriod{},
					},
					BillingLoadCurve: []BillingLoadCurve{
						{
							R1:     10000,
							Period: measures.P0,
						},
						{
							R1:     10000,
							Period: measures.P0,
						},
						{
							R1:     10000,
							Period: measures.P0,
						},
						{
							R1:     10000,
							Period: measures.P0,
						},
						{
							R1:     10000,
							Period: measures.P0,
						},
					},
				},
			},
			output: output{
				billingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						R1:                   21100,
						BalanceTypeR1:        EstimateHistoryProfile,
						BalanceGeneralTypeR1: GeneralEstimated,
						BalanceMeasureTypeR1: ProvisionalBalanceMeasure,
						BalanceOriginR1:      EstimateOrigin,
						EstimatedCodeR1:      5,
					},
				},
				id: "BALANCE_ESTIMATED_HISTORY_TLG",
			},
			want: nil,
		},
		"Should be calculate balance correct (Period in QuarterHour Hours Not Equal)": {
			input: input{
				period:    measures.P0,
				magnitude: measures.R4,
				contextTlg: GraphContext{
					LastHistory: BillingMeasure{
						RegisterType: measures.QuarterHour,
						BillingBalance: BillingBalance{
							P0: &BillingBalancePeriod{
								R4: 13301,
							},
						},
						BillingLoadCurve: []BillingLoadCurve{
							{
								R4:     10000,
								Period: measures.P0,
							},
							{
								R4:     10000,
								Period: measures.P0,
							},
							{
								R4:     10000,
								Period: measures.P0,
							},
							{
								R4:     10000,
								Period: measures.P0,
							},
						},
					},
				},
				b: &BillingMeasure{
					RegisterType: measures.QuarterHour,
					BillingBalance: BillingBalance{
						P0: &BillingBalancePeriod{},
					},
					BillingLoadCurve: []BillingLoadCurve{
						{
							R4:     10000,
							Period: measures.P0,
						},
						{
							R4:     10000,
							Period: measures.P0,
						},
						{
							R4:     10000,
							Period: measures.P0,
						},
					},
				},
			},
			output: output{
				billingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						R4:                   9976,
						BalanceTypeR4:        EstimateHistoryProfile,
						BalanceGeneralTypeR4: GeneralEstimated,
						BalanceMeasureTypeR4: ProvisionalBalanceMeasure,
						BalanceOriginR4:      EstimateOrigin,
						EstimatedCodeR4:      5,
					},
				},
				id: "BALANCE_ESTIMATED_HISTORY_TLG",
			},
			want: nil,
		},
		"Should be error actual period is nil": {
			input: input{
				period:    measures.P0,
				magnitude: measures.R2,
				contextTlg: GraphContext{
					LastHistory: BillingMeasure{
						RegisterType: measures.Hourly,
						BillingBalance: BillingBalance{
							P0: &BillingBalancePeriod{
								AI: 20410,
								AE: 23100,
								R1: 21100,
								R2: 21340,
								R3: 23110,
								R4: 13301,
							},
						},
					},
				},
				b: &BillingMeasure{
					RegisterType: measures.Hourly,
				},
			},
			output: output{
				id: "BALANCE_ESTIMATED_HISTORY_TLG",
			},
			want: errors.New("period cannot be nil"),
		},
		"Should be error history period is nil": {
			input: input{
				period:    measures.P0,
				magnitude: measures.R3,
				contextTlg: GraphContext{
					LastHistory: BillingMeasure{
						RegisterType: measures.Hourly,
					},
				},
				b: &BillingMeasure{
					RegisterType: measures.Hourly,
					BillingBalance: BillingBalance{
						P0: &BillingBalancePeriod{},
					},
				},
			},
			output: output{
				id: "BALANCE_ESTIMATED_HISTORY_TLG",
			},
			want: errors.New("period cannot be nil"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			estimatedHistoryTlg := NewEstimatedHistoryTlg(test.input.b, test.input.period, test.input.magnitude, &test.input.contextTlg)

			id := estimatedHistoryTlg.ID()
			result := estimatedHistoryTlg.Execute(context.Background())

			if result != nil {
				assert.Equal(t, result, test.want)
			} else {
				assert.Equal(t, test.output.id, id)
				assert.Equal(t, *test.output.billingBalance.P0, *test.input.b.BillingBalance.P0)
			}
		})
	}
}

func Test_Unit_Domain_Algorithms_EstimatedBalanceByPowerDemand(t *testing.T) {
	type input struct {
		inputB         *BillingMeasure
		inputPeriodKey measures.PeriodKey
		magnitude      measures.Magnitude
	}
	type expected struct {
		expectedB     BillingBalance
		expectedError error
	}
	tests := map[string]struct {
		input    input
		expected expected
	}{
		"Should return err because of the bad key input": {
			input: input{
				inputB:         &BillingMeasure{},
				inputPeriodKey: measures.PeriodKey("RandomKey"),
			},
			expected: expected{
				expectedError: ErrorPeriodKeyDoesNotMatch,
			},
		},
		"For powerDemanded <=10 with Period 1 magnitude measures.AI": {
			input: input{
				inputB: &BillingMeasure{
					BillingBalance: BillingBalance{
						P1: &BillingBalancePeriod{},
					},
					BillingLoadCurve: []BillingLoadCurve{
						{Period: measures.P1},
						{Period: measures.P1},
						{Period: measures.P1},
						{Period: measures.P1},
						{Period: measures.P1},
						{Period: measures.P1},
					},
					P1Demand: 2,
				},
				inputPeriodKey: measures.P1,
				magnitude:      measures.AI,
			},
			expected: expected{
				expectedB: BillingBalance{
					P1: &BillingBalancePeriod{
						AI:                   250,
						BalanceTypeAI:        EstimatedContractPower,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
						BalanceOriginAI:      EstimateOrigin,
						EstimatedCodeAI:      6,
					},
				},
				expectedError: nil,
			},
		},
		"For powerDemanded <10 with Period 3 magnitude measures.AI": {
			input: input{
				inputB: &BillingMeasure{
					BillingBalance: BillingBalance{
						P3: &BillingBalancePeriod{},
					},
					BillingLoadCurve: []BillingLoadCurve{
						{Period: measures.P3},
						{Period: measures.P3},
						{Period: measures.P3},
						{Period: measures.P3},
						{Period: measures.P3},
						{Period: measures.P3},
					},
					P3Demand: 10,
				},
				inputPeriodKey: measures.P3,
				magnitude:      measures.AI,
			},
			expected: expected{
				expectedB: BillingBalance{
					P3: &BillingBalancePeriod{
						AI:                   6750,
						BalanceTypeAI:        EstimatedContractPower,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
						BalanceOriginAI:      EstimateOrigin,
						EstimatedCodeAI:      6,
					},
				},
				expectedError: nil,
			},
		},
		"For powerDemanded >10 with Period 3 magnitude measures.AI": {
			input: input{
				inputB: &BillingMeasure{
					BillingBalance: BillingBalance{
						P3: &BillingBalancePeriod{},
					},
					BillingLoadCurve: []BillingLoadCurve{
						{Period: measures.P3},
						{Period: measures.P3},
						{Period: measures.P3},
						{Period: measures.P3},
						{Period: measures.P3},
						{Period: measures.P3},
					},
					P3Demand: 20,
				},
				inputPeriodKey: measures.P3,
				magnitude:      measures.AI,
			},
			expected: expected{
				expectedB: BillingBalance{
					P3: &BillingBalancePeriod{
						AI:                   39600,
						BalanceTypeAI:        EstimatedContractPower,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
						BalanceOriginAI:      EstimateOrigin,
						EstimatedCodeAI:      6,
					},
				},
				expectedError: nil,
			},
		},
		"For magnitude != measures.AI": {
			input: input{
				inputB: &BillingMeasure{
					BillingBalance: BillingBalance{
						P4: &BillingBalancePeriod{
							AE: 0,
						},
					},
					P4Demand: 11,
				},
				inputPeriodKey: measures.P4,
				magnitude:      measures.AE,
			},
			expected: expected{
				expectedB: BillingBalance{
					P4: &BillingBalancePeriod{
						AE:                   0,
						BalanceTypeAE:        EstimatedContractPower,
						BalanceGeneralTypeAE: GeneralEstimated,
						BalanceMeasureTypeAE: ProvisionalBalanceMeasure,
						BalanceOriginAE:      EstimateOrigin,
						EstimatedCodeAE:      6,
					},
				},
				expectedError: nil,
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			exeResult := NewEstimatedBalanceByPowerDemand(testCase.input.inputB, testCase.input.inputPeriodKey, testCase.input.magnitude)
			exeErr := exeResult.Execute(ctx)
			assert.Equal(t, testCase.expected.expectedError, exeErr)
			assert.Equal(t, testCase.expected.expectedB, exeResult.b.BillingBalance)
		})
	}
}
