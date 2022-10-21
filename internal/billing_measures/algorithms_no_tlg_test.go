package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Domain_BillingMeasures_FlatCastNoTLGAlgorithm(t *testing.T) {
	tests := map[string]struct {
		b         *BillingMeasure
		magnitude measures.Magnitude
		expected  []BillingLoadCurve
		ctx       context.Context
	}{
		"Should Be Okay": {
			ctx:       context.Background(),
			magnitude: measures.AI,
			b: &BillingMeasure{
				PointType: "1",
				Periods:   []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
				BillingBalance: BillingBalance{
					P4: &BillingBalancePeriod{AI: 230000},
					P5: &BillingBalancePeriod{AI: 184000},
					P6: &BillingBalancePeriod{AI: 200000, AE: 200000},
				},
				BillingLoadCurve: []BillingLoadCurve{
					{
						Equipment: measures.Main,
						Origin:    measures.Filled,
						Period:    measures.P6,
					},
					{
						Equipment: measures.Redundant,
						AI:        16000,
						AE:        16000,
						Period:    measures.P6,
					},
					{
						Equipment: measures.Main,
						AI:        24000,
						AE:        24000,
						Period:    measures.P6,
					},
					{
						Equipment: measures.Receipt,
						AI:        25000,
						AE:        25000,
						Period:    measures.P6,
					},
					{
						Equipment: measures.Receipt,
						AI:        19000,
						AE:        19000,
						Period:    measures.P6,
					},
					{
						Equipment: measures.Receipt,
						AI:        27000,
						AE:        27000,
						Period:    measures.P6,
					},
					{
						Equipment: measures.Receipt,
						AI:        29000,
						AE:        29000,
						Period:    measures.P6,
					},
					{
						Equipment: measures.Receipt,
						AI:        30000,
						AE:        30000,
						Period:    measures.P6,
					},
					{
						Equipment: measures.Receipt,
						AI:        22000,
						Period:    measures.P5,
					},
					{
						Equipment: measures.Receipt,
						Origin:    measures.Filled,
						Period:    measures.P4,
					},
					{
						Equipment: measures.Receipt,
						Origin:    measures.Filled,
						Period:    measures.P4,
					},
					{
						Equipment: measures.Receipt,
						Origin:    measures.Filled,
						Period:    measures.P4,
					},
					{
						Equipment: measures.Receipt,
						Origin:    measures.Filled,
						Period:    measures.P4,
					},
					{
						Equipment: measures.Receipt,
						AI:        21000,
						Period:    measures.P4,
					},
					{
						Equipment: measures.Receipt,
						AI:        30000,
						Period:    measures.P5,
					},
					{
						Equipment: measures.Receipt,
						AI:        45000,
						Period:    measures.P5,
					},
					{
						Equipment: measures.Receipt,
						AI:        20000,
						Period:    measures.P5,
					},
					{
						Equipment: measures.Receipt,
						AI:        28000,
						Period:    measures.P5,
					},
					{
						Equipment: measures.Receipt,
						AI:        12000,
						Period:    measures.P4,
					},
					{
						Equipment: measures.Receipt,
						AI:        23000,
						Period:    measures.P4,
					},

					{
						Equipment: measures.Receipt,
						AI:        24000,
						Period:    measures.P4,
					},

					{
						Equipment: measures.Receipt,
						AI:        19000,
						Period:    measures.P4,
					},
					{
						Equipment: measures.Receipt,
						AI:        20000,
						Period:    measures.P5,
					},
					{
						Equipment: measures.Receipt,
						AI:        19000,
						Period:    measures.P5,
					},
				},
			},
			expected: []BillingLoadCurve{
				{
					Origin:                   measures.Filled,
					Equipment:                measures.Main,
					AI:                       30000,
					Period:                   measures.P6,
					EstimatedCodeAI:          8,
					EstimatedMethodAI:        EstimatedByFlatProfile,
					EstimatedGeneralMethodAI: GeneralEstimated,
					MeasureTypeAI:            ProvisionalBalanceMeasure,
				},
				{
					Equipment:                measures.Redundant,
					AI:                       16000,
					AE:                       16000,
					Period:                   measures.P6,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Main,
					AI:                       24000,
					AE:                       24000,
					Period:                   measures.P6,
					EstimatedCodeAI:          1,
					EstimatedMethodAI:        FirmMainConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       25000,
					AE:                       25000,
					Period:                   measures.P6,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       19000,
					AE:                       19000,
					Period:                   measures.P6,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       27000,
					AE:                       27000,
					Period:                   measures.P6,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       29000,
					AE:                       29000,
					Period:                   measures.P6,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       30000,
					AE:                       30000,
					Period:                   measures.P6,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       22000,
					Period:                   measures.P5,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					Origin:                   measures.Filled,
					AI:                       32750,
					Period:                   measures.P4,
					EstimatedCodeAI:          8,
					EstimatedMethodAI:        EstimatedByFlatProfile,
					EstimatedGeneralMethodAI: GeneralEstimated,
					MeasureTypeAI:            ProvisionalBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					Origin:                   measures.Filled,
					AI:                       32750,
					Period:                   measures.P4,
					EstimatedCodeAI:          8,
					EstimatedMethodAI:        EstimatedByFlatProfile,
					EstimatedGeneralMethodAI: GeneralEstimated,
					MeasureTypeAI:            ProvisionalBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					Origin:                   measures.Filled,
					AI:                       32750,
					Period:                   measures.P4,
					EstimatedCodeAI:          8,
					EstimatedMethodAI:        EstimatedByFlatProfile,
					EstimatedGeneralMethodAI: GeneralEstimated,
					MeasureTypeAI:            ProvisionalBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					Origin:                   measures.Filled,
					AI:                       32750,
					Period:                   measures.P4,
					EstimatedCodeAI:          8,
					EstimatedMethodAI:        EstimatedByFlatProfile,
					EstimatedGeneralMethodAI: GeneralEstimated,
					MeasureTypeAI:            ProvisionalBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       21000,
					Period:                   measures.P4,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       30000,
					Period:                   measures.P5,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       45000,
					Period:                   measures.P5,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       20000,
					Period:                   measures.P5,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       28000,
					Period:                   measures.P5,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       12000,
					Period:                   measures.P4,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       23000,
					Period:                   measures.P4,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},

				{
					Equipment:                measures.Receipt,
					AI:                       24000,
					Period:                   measures.P4,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       19000,
					Period:                   measures.P4,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       20000,
					Period:                   measures.P5,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Equipment:                measures.Receipt,
					AI:                       19000,
					Period:                   measures.P5,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			flatCast := NewFlatCastNoTLG(testCase.b, testCase.magnitude)
			flatCast.Execute(context.Background())

			assert.Equal(t, "CCH_FLAT_CAST_NO_TLG", flatCast.ID())
			assert.Equal(t, testCase.expected, flatCast.b.BillingLoadCurve)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_CloseSum(t *testing.T) {

	type input struct {
		b         *BillingMeasure
		magnitude measures.Magnitude
	}

	type output struct {
		b BillingBalance
	}

	type want struct {
		id  string
		err error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be calculated, estimated code 2 Origin STM": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
					BillingBalance: BillingBalance{
						Origin: measures.STM,
						P0:     &BillingBalancePeriod{},
						P4: &BillingBalancePeriod{
							AI: 230000,
						},
						P5: &BillingBalancePeriod{
							AI: 184000,
						},
						P6: &BillingBalancePeriod{
							AI: 200000,
						},
					},
				}},
			output: output{
				b: BillingBalance{
					Origin: measures.STM,
					P0: &BillingBalancePeriod{
						AI:                   614000,
						EstimatedCodeAI:      2,
						BalanceTypeAI:        CalculatedByCloseSum,
						BalanceGeneralTypeAI: GeneralCalculated,
						BalanceOriginAI:      TlmOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P6: &BillingBalancePeriod{
						AI:                   200000,
						EstimatedCodeAI:      1,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      TlmOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P5: &BillingBalancePeriod{
						AI:                   184000,
						EstimatedCodeAI:      1,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      TlmOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P4: &BillingBalancePeriod{
						AI:                   230000,
						EstimatedCodeAI:      1,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      TlmOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
				},
			},
			want: want{id: "BALANCE_CALCULATED_BY_CLOSE_SUM", err: nil},
		},
		"Should be calculated, estimated code 3 Origin TPL": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
					BillingBalance: BillingBalance{
						Origin: measures.TPL,
						P0:     &BillingBalancePeriod{},
						P4: &BillingBalancePeriod{
							AI: 230000,
						},
						P5: &BillingBalancePeriod{
							AI: 184000,
						},
						P6: &BillingBalancePeriod{
							AI: 200000,
						},
					},
				}},
			output: output{
				b: BillingBalance{
					Origin: measures.TPL,
					P0: &BillingBalancePeriod{
						AI:                   614000,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        CalculatedByCloseSum,
						BalanceGeneralTypeAI: GeneralCalculated,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P6: &BillingBalancePeriod{
						AI:                   200000,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        RealByAbsLocalRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P5: &BillingBalancePeriod{
						AI:                   184000,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        RealByAbsLocalRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P4: &BillingBalancePeriod{
						AI:                   230000,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        RealByAbsLocalRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
				},
			},
			want: want{id: "BALANCE_CALCULATED_BY_CLOSE_SUM", err: nil},
		},
		"Should be calculated, estimated code 3 Origin File": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
					BillingBalance: BillingBalance{
						Origin: measures.File,
						P0:     &BillingBalancePeriod{},
						P4: &BillingBalancePeriod{
							AI: 230000,
						},
						P5: &BillingBalancePeriod{
							AI: 184000,
						},
						P6: &BillingBalancePeriod{
							AI: 200000,
						},
					},
				}},
			output: output{
				b: BillingBalance{
					Origin: measures.File,
					P0: &BillingBalancePeriod{
						AI:                   614000,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        CalculatedByCloseSum,
						BalanceGeneralTypeAI: GeneralCalculated,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P6: &BillingBalancePeriod{
						AI:                   200000,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        RealByAbsLocalRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P5: &BillingBalancePeriod{
						AI:                   184000,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        RealByAbsLocalRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P4: &BillingBalancePeriod{
						AI:                   230000,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        RealByAbsLocalRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
				},
			},
			want: want{id: "BALANCE_CALCULATED_BY_CLOSE_SUM", err: nil},
		},
		"Should be calculated, estimated code 4 Origin AutoRead": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
					BillingBalance: BillingBalance{
						Origin: measures.Auto,
						P0:     &BillingBalancePeriod{},
						P4: &BillingBalancePeriod{
							AI: 230000,
						},
						P5: &BillingBalancePeriod{
							AI: 184000,
						},
						P6: &BillingBalancePeriod{
							AI: 200000,
						},
					},
				}},
			output: output{
				b: BillingBalance{
					Origin: measures.Auto,
					P0: &BillingBalancePeriod{
						AI:                   614000,
						EstimatedCodeAI:      4,
						BalanceTypeAI:        CalculatedByCloseSum,
						BalanceGeneralTypeAI: GeneralCalculated,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
					P6: &BillingBalancePeriod{
						AI:                   200000,
						EstimatedCodeAI:      4,
						BalanceTypeAI:        RealByAutoRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
					P5: &BillingBalancePeriod{
						AI:                   184000,
						EstimatedCodeAI:      4,
						BalanceTypeAI:        RealByAutoRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
					P4: &BillingBalancePeriod{
						AI:                   230000,
						EstimatedCodeAI:      4,
						BalanceTypeAI:        RealByAutoRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
				},
			},
			want: want{id: "BALANCE_CALCULATED_BY_CLOSE_SUM", err: nil},
		},
	}

	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			closeSum := NewCloseSum(test.input.b, test.input.magnitude)
			id := closeSum.ID()
			err := closeSum.Execute(ctx)

			assert.Equal(t, test.output.b, closeSum.b.BillingBalance)
			assert.Equal(t, test.want.id, id)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_BalanceCompleteNoTlg(t *testing.T) {

	type input struct {
		b         *BillingMeasure
		magnitude measures.Magnitude
	}

	type output struct {
		b BillingBalance
	}

	type want struct {
		err error
	}

	testCases := map[string]struct {
		input  input
		output output
		want
	}{
		"Should be correct, origin STM": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},

					BillingBalance: BillingBalance{
						Origin: measures.STM,
						P0:     &BillingBalancePeriod{},
						P4:     &BillingBalancePeriod{},
						P5:     &BillingBalancePeriod{},
						P6:     &BillingBalancePeriod{},
					},
				}},
			output: output{
				b: BillingBalance{
					Origin: measures.STM,
					P0: &BillingBalancePeriod{
						EstimatedCodeAI:      1,
						BalanceOriginAI:      TlmOrigin,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P4: &BillingBalancePeriod{
						EstimatedCodeAI:      1,
						BalanceOriginAI:      TlmOrigin,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P5: &BillingBalancePeriod{
						EstimatedCodeAI:      1,
						BalanceOriginAI:      TlmOrigin,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P6: &BillingBalancePeriod{
						EstimatedCodeAI:      1,
						BalanceOriginAI:      TlmOrigin,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
				},
			},
		},
		"Should be correct, origin Manual": {
			input: input{
				magnitude: measures.AE,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
					BillingBalance: BillingBalance{
						Origin: measures.Manual,
						P0:     &BillingBalancePeriod{},
						P4:     &BillingBalancePeriod{},
						P5:     &BillingBalancePeriod{},
						P6:     &BillingBalancePeriod{},
					},
				}},
			output: output{
				b: BillingBalance{
					Origin: measures.Manual,
					P0: &BillingBalancePeriod{
						EstimatedCodeAE:      3,
						BalanceOriginAE:      LocalOrigin,
						BalanceTypeAE:        RealByAbsLocalRead,
						BalanceGeneralTypeAE: GeneralReal,
						BalanceMeasureTypeAE: FirmBalanceMeasure,
					},
					P4: &BillingBalancePeriod{
						EstimatedCodeAE:      3,
						BalanceOriginAE:      LocalOrigin,
						BalanceTypeAE:        RealByAbsLocalRead,
						BalanceGeneralTypeAE: GeneralReal,
						BalanceMeasureTypeAE: FirmBalanceMeasure,
					},
					P5: &BillingBalancePeriod{
						EstimatedCodeAE:      3,
						BalanceOriginAE:      LocalOrigin,
						BalanceTypeAE:        RealByAbsLocalRead,
						BalanceGeneralTypeAE: GeneralReal,
						BalanceMeasureTypeAE: FirmBalanceMeasure,
					},
					P6: &BillingBalancePeriod{
						EstimatedCodeAE:      3,
						BalanceOriginAE:      LocalOrigin,
						BalanceTypeAE:        RealByAbsLocalRead,
						BalanceGeneralTypeAE: GeneralReal,
						BalanceMeasureTypeAE: FirmBalanceMeasure,
					},
				},
			},
		},
		"Should be error, invalid origin": {
			input: input{b: &BillingMeasure{
				Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
				BillingBalance: BillingBalance{
					Origin: measures.STG,
					P0:     &BillingBalancePeriod{},
					P4:     &BillingBalancePeriod{},
					P5:     &BillingBalancePeriod{},
					P6:     &BillingBalancePeriod{},
				},
			}},
			output: output{
				b: BillingBalance{
					Origin: measures.STG,
					P0:     &BillingBalancePeriod{},
					P4:     &BillingBalancePeriod{},
					P5:     &BillingBalancePeriod{},
					P6:     &BillingBalancePeriod{},
				},
			},
			want: want{err: fmt.Errorf("invalid origin %s", measures.STG)},
		},
	}

	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			algorithm := NewBalanceCompleteNoTlg(test.input.b, test.input.magnitude)
			err := algorithm.Execute(ctx)

			assert.Equal(t, test.output.b, algorithm.b.BillingBalance)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_CCHCompleteNoTlg(t *testing.T) {

	type input struct {
		b         *BillingMeasure
		magnitude measures.Magnitude
	}

	type output struct {
		curve []BillingLoadCurve
	}

	type want struct {
		id  string
		err error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be correct, point type 3": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					PointType: "3",
					BillingLoadCurve: []BillingLoadCurve{
						{
							Origin: measures.STM,
						},
						{
							Origin: measures.STM,
						},
					},
				},
			},
			output: output{
				curve: []BillingLoadCurve{
					{
						Origin:                   measures.STM,
						EstimatedCodeAI:          1,
						EstimatedMethodAI:        RealBalance,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
					{
						Origin:                   measures.STM,
						EstimatedCodeAI:          1,
						EstimatedMethodAI:        RealBalance,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
				},
			},
			want: want{
				id:  "CCH_COMPLETE",
				err: nil,
			},
		},
		"Should be correct, point type 4": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					PointType: "4",
					BillingLoadCurve: []BillingLoadCurve{
						{
							Origin: measures.TPL,
						},
						{
							Origin: measures.TPL,
						},
					},
				},
			},
			output: output{
				curve: []BillingLoadCurve{
					{
						Origin:                   measures.TPL,
						EstimatedCodeAI:          1,
						EstimatedMethodAI:        RealBalance,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
					{
						Origin:                   measures.TPL,
						EstimatedCodeAI:          1,
						EstimatedMethodAI:        RealBalance,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
				},
			},
			want: want{
				id:  "CCH_COMPLETE",
				err: nil,
			},
		},
		"Should be correct, point type 5": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					PointType: "5",
					BillingLoadCurve: []BillingLoadCurve{
						{
							Origin: measures.File,
						},
						{
							Origin: measures.File,
						},
					},
				},
			},
			output: output{
				curve: []BillingLoadCurve{
					{
						Origin:                   measures.File,
						EstimatedCodeAI:          1,
						EstimatedMethodAI:        RealBalance,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
					{
						Origin:                   measures.File,
						EstimatedCodeAI:          1,
						EstimatedMethodAI:        RealBalance,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
				},
			},
			want: want{
				id:  "CCH_COMPLETE",
				err: nil,
			},
		},
		"Should be correct, point type 1 with main equipment": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					PointType: "1",
					BillingLoadCurve: []BillingLoadCurve{
						{
							Origin:    measures.Visual,
							Equipment: measures.Main,
						},
						{
							Origin:    measures.Visual,
							Equipment: measures.Main,
						},
					},
				},
			},
			output: output{
				curve: []BillingLoadCurve{
					{
						Origin:                   measures.Visual,
						Equipment:                measures.Main,
						EstimatedCodeAI:          1,
						EstimatedMethodAI:        FirmMainConfig,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
					{
						Origin:                   measures.Visual,
						Equipment:                measures.Main,
						EstimatedCodeAI:          1,
						EstimatedMethodAI:        FirmMainConfig,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
				},
			},
			want: want{
				id:  "CCH_COMPLETE",
				err: nil,
			},
		},
		"Should be correct, point type 1 with redundant equipment": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					PointType: "2",
					BillingLoadCurve: []BillingLoadCurve{
						{
							Origin:    measures.STM,
							Equipment: measures.Redundant,
						},
						{
							Origin:    measures.STM,
							Equipment: measures.Redundant,
						},
					},
				},
			},
			output: output{
				curve: []BillingLoadCurve{
					{
						Origin:                   measures.STM,
						Equipment:                measures.Redundant,
						EstimatedCodeAI:          2,
						EstimatedMethodAI:        FirmRedundantConfig,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
					{
						Origin:                   measures.STM,
						Equipment:                measures.Redundant,
						EstimatedCodeAI:          2,
						EstimatedMethodAI:        FirmRedundantConfig,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
				},
			},
			want: want{
				id:  "CCH_COMPLETE",
				err: nil,
			},
		},
		"Should be correct, point type 2 with receipt equipment": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					PointType: "2",
					BillingLoadCurve: []BillingLoadCurve{
						{
							Origin:    measures.Manual,
							Equipment: measures.Receipt,
						},
						{
							Origin:    measures.Manual,
							Equipment: measures.Receipt,
						},
					},
				},
			},
			output: output{
				curve: []BillingLoadCurve{
					{
						Origin:                   measures.Manual,
						Equipment:                measures.Receipt,
						EstimatedCodeAI:          3,
						EstimatedMethodAI:        FirmReceiptConfig,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
					{
						Origin:                   measures.Manual,
						Equipment:                measures.Receipt,
						EstimatedCodeAI:          3,
						EstimatedMethodAI:        FirmReceiptConfig,
						EstimatedGeneralMethodAI: GeneralReal,
						MeasureTypeAI:            FirmBalanceMeasure,
					},
				},
			},
			want: want{
				id:  "CCH_COMPLETE",
				err: nil,
			},
		},
	}

	for name, _ := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			algorithm := NewCCHCompleteNoTlg(test.input.b, test.input.magnitude)

			id := algorithm.ID()
			err := algorithm.Execute(ctx)

			assert.Equal(t, test.want.id, id)
			assert.Equal(t, test.want.err, err)
			assert.ElementsMatch(t, test.output.curve, algorithm.b.BillingLoadCurve)
			assert.Equal(t, test.output.curve, algorithm.b.BillingLoadCurve)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_FillOneCloseAtr(t *testing.T) {

	type input struct {
		b         *BillingMeasure
		magnitude measures.Magnitude
	}

	type output struct {
		b BillingBalance
	}

	type want struct {
		err error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be calculated P4, STM": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
					BillingBalance: BillingBalance{
						Origin: measures.STM,
						P0: &BillingBalancePeriod{
							AI: 500000,
						},
						P4: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
						P5: &BillingBalancePeriod{
							AI:                  184000,
							BalanceValidationAI: measures.Valid,
						},
						P6: &BillingBalancePeriod{
							AI:                  200000,
							BalanceValidationAI: measures.Valid,
						},
					},
				}},
			output: output{
				b: BillingBalance{
					Origin: measures.STM,
					P0: &BillingBalancePeriod{
						AI:                   500000,
						EstimatedCodeAI:      1,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      TlmOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P4: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Invalid,
						AI:                   116000,
						EstimatedCodeAI:      2,
						BalanceTypeAI:        CalculatedByCloseBalance,
						BalanceGeneralTypeAI: GeneralCalculated,
						BalanceOriginAI:      TlmOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P5: &BillingBalancePeriod{
						AI:                   184000,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      1,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      TlmOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P6: &BillingBalancePeriod{
						AI:                   200000,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      1,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      TlmOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		"Should be calculated P4, TPL": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
					BillingBalance: BillingBalance{
						Origin: measures.TPL,
						P0: &BillingBalancePeriod{
							AI: 500000,
						},
						P4: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
						P5: &BillingBalancePeriod{
							AI:                  184000,
							BalanceValidationAI: measures.Valid,
						},
						P6: &BillingBalancePeriod{
							AI:                  200000,
							BalanceValidationAI: measures.Valid,
						},
					},
				}},
			output: output{
				b: BillingBalance{
					Origin: measures.TPL,
					P0: &BillingBalancePeriod{
						AI:                   500000,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        RealByAbsLocalRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P4: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Invalid,
						AI:                   116000,
						EstimatedCodeAI:      2,
						BalanceTypeAI:        CalculatedByCloseBalance,
						BalanceGeneralTypeAI: GeneralCalculated,
						BalanceOriginAI:      TlmOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P5: &BillingBalancePeriod{
						AI:                   184000,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        RealByAbsLocalRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P6: &BillingBalancePeriod{
						AI:                   200000,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        RealByAbsLocalRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		"Should be calculated P4, Auto": {
			input: input{
				magnitude: measures.AE,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
					BillingBalance: BillingBalance{
						Origin: measures.Auto,
						P0: &BillingBalancePeriod{
							AE: 500000,
						},
						P4: &BillingBalancePeriod{
							BalanceValidationAE: measures.Invalid,
						},
						P5: &BillingBalancePeriod{
							AE:                  184000,
							BalanceValidationAE: measures.Valid,
						},
						P6: &BillingBalancePeriod{
							AE:                  200000,
							BalanceValidationAE: measures.Valid,
						},
					},
				}},
			output: output{
				b: BillingBalance{
					Origin: measures.Auto,
					P0: &BillingBalancePeriod{
						AE:                   500000,
						EstimatedCodeAE:      4,
						BalanceTypeAE:        RealByAutoRead,
						BalanceGeneralTypeAE: GeneralReal,
						BalanceOriginAE:      LocalOrigin,
						BalanceMeasureTypeAE: ProvisionalBalanceMeasure,
					},
					P4: &BillingBalancePeriod{
						BalanceValidationAE:  measures.Invalid,
						AE:                   116000,
						EstimatedCodeAE:      2,
						BalanceTypeAE:        CalculatedByCloseBalance,
						BalanceGeneralTypeAE: GeneralCalculated,
						BalanceOriginAE:      TlmOrigin,
						BalanceMeasureTypeAE: FirmBalanceMeasure,
					},
					P5: &BillingBalancePeriod{
						AE:                   184000,
						BalanceValidationAE:  measures.Valid,
						EstimatedCodeAE:      4,
						BalanceTypeAE:        RealByAutoRead,
						BalanceGeneralTypeAE: GeneralReal,
						BalanceOriginAE:      LocalOrigin,
						BalanceMeasureTypeAE: ProvisionalBalanceMeasure,
					},
					P6: &BillingBalancePeriod{
						AE:                   200000,
						BalanceValidationAE:  measures.Valid,
						EstimatedCodeAE:      4,
						BalanceTypeAE:        RealByAutoRead,
						BalanceGeneralTypeAE: GeneralReal,
						BalanceOriginAE:      LocalOrigin,
						BalanceMeasureTypeAE: ProvisionalBalanceMeasure,
					},
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for name, _ := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			algorithm := NewFillOneCloseAtr(test.input.b, test.input.magnitude)
			err := algorithm.Execute(ctx)

			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.output.b, algorithm.b.BillingBalance)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_SumHoursNoClosureAlgorithm(t *testing.T) {
	type input struct {
		b   *BillingMeasure
		ctx context.Context
	}
	type expected struct {
		err            error
		billingBalance BillingBalance
	}

	tests := map[string]struct {
		input    input
		expected expected
	}{
		"Should return ok": {
			input: input{
				b: &BillingMeasure{
					PointType: "1",

					BillingBalance: BillingBalance{
						Origin: measures.STM,
						P0: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
						P4: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
						P5: &BillingBalancePeriod{
							AI:                  184000,
							BalanceValidationAI: measures.Valid,
						},
						P6: &BillingBalancePeriod{
							AI:                  200000,
							BalanceValidationAI: measures.Valid,
						},
					},
					BillingLoadCurve: []BillingLoadCurve{
						{
							AI:              25600,
							Period:          measures.P4,
							EstimatedCodeAI: 11,
							Origin:          measures.Filled,
						},
						{
							AI:              28000,
							Period:          measures.P4,
							EstimatedCodeAI: 11,
							Origin:          measures.Filled,
						},
						{
							AI:              32750,
							Period:          measures.P4,
							EstimatedCodeAI: 11,
							Origin:          measures.Filled,
						},
						{
							AI:              30400,
							Period:          measures.P4,
							EstimatedCodeAI: 11,
							Origin:          measures.Filled,
						},
						{
							AI:              21000,
							Period:          measures.P4,
							EstimatedCodeAI: 11,
							Origin:          measures.Filled,
						},
						{
							AI:              12000,
							Period:          measures.P4,
							EstimatedCodeAI: 11,
							Origin:          measures.Filled,
						},
						{
							AI:              23000,
							Period:          measures.P4,
							EstimatedCodeAI: 11,
							Origin:          measures.Filled,
						},
						{
							AI:              24000,
							Period:          measures.P4,
							EstimatedCodeAI: 11,
							Origin:          measures.Filled,
						},
						{
							AI:              19000,
							Period:          measures.P4,
							EstimatedCodeAI: 11,
							Origin:          measures.Filled,
						},
						{
							AI:              24500,
							Period:          measures.P6,
							EstimatedCodeAI: 0,
							Origin:          measures.STM,
						},
						{
							AI:              16000,
							Period:          measures.P6,
							EstimatedCodeAI: 0,
							Origin:          measures.STM,
						},
						{
							AI:              24000,
							Period:          measures.P6,
							EstimatedCodeAI: 0,
							Origin:          measures.STM,
						},
						{
							AI:              25000,
							Period:          measures.P6,
							EstimatedCodeAI: 0,
							Origin:          measures.STM,
						},
						{
							AI:              19000,
							Period:          measures.P6,
							EstimatedCodeAI: 0,
							Origin:          measures.STM,
						},
						{
							AI:              27000,
							Period:          measures.P6,
							EstimatedCodeAI: 0,
							Origin:          measures.STM,
						},
						{
							AI:              29000,
							Period:          measures.P6,
							EstimatedCodeAI: 0,
							Origin:          measures.STM,
						},
						{
							AI:              30000,
							Period:          measures.P6,
							EstimatedCodeAI: 0,
							Origin:          measures.STM,
						},
						{
							AI:              22000,
							Period:          measures.P5,
							EstimatedCodeAI: 3,
							Origin:          measures.Auto,
						},
						{
							AI:              30000,
							Period:          measures.P5,
							EstimatedCodeAI: 3,
							Origin:          measures.Auto,
						},
						{
							AI:              45000,
							Period:          measures.P5,
							EstimatedCodeAI: 3,
							Origin:          measures.Auto,
						},
						{
							AI:              20000,
							Period:          measures.P5,
							EstimatedCodeAI: 3,
							Origin:          measures.Auto,
						},
						{
							AI:              28000,
							Period:          measures.P5,
							EstimatedCodeAI: 3,
							Origin:          measures.Auto,
						},
						{
							AI:              20000,
							Period:          measures.P5,
							EstimatedCodeAI: 3,
							Origin:          measures.Auto,
						},
						{
							AI:              19000,
							Period:          measures.P5,
							EstimatedCodeAI: 3,
							Origin:          measures.Auto,
						},
					},
					Periods: []measures.PeriodKey{measures.P4, measures.P6, measures.P5},
				},
				ctx: context.Background(),
			},
			expected: expected{
				err: nil,
				billingBalance: BillingBalance{
					Origin: measures.STM,
					P0: &BillingBalancePeriod{
						AI:                   599750,
						BalanceTypeAI:        CalculatedByCloseSum,
						BalanceGeneralTypeAI: GeneralCalculated,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
						BalanceOriginAI:      EstimateOrigin,
						EstimatedCodeAI:      6,
						BalanceValidationAI:  measures.Invalid,
					},
					P4: &BillingBalancePeriod{
						AI:                   215750,
						BalanceTypeAI:        ObtainedByCurve,
						BalanceGeneralTypeAI: GeneralCalculated,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
						BalanceOriginAI:      EstimateOrigin,
						EstimatedCodeAI:      6,
						BalanceValidationAI:  measures.Invalid,
					},
					P6: &BillingBalancePeriod{
						AI:                   200000,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      1,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
						BalanceOriginAI:      TlmOrigin,
					},
					P5: &BillingBalancePeriod{
						AI:                   184000,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      1,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
						BalanceOriginAI:      TlmOrigin,
					},
				},
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			closeSum := NewSumHoursNoClosure(testCase.input.b, measures.AI)
			id := closeSum.ID()
			err := closeSum.Execute(testCase.input.ctx)
			assert.Equal(t, "BALANCE_SUM_HOURS_NO_CLOSURE", id)
			assert.Equal(t, testCase.expected.err, err)
			assert.Equal(t, testCase.expected.billingBalance, closeSum.b.BillingBalance)

		})
	}
}

func Test_Unit_Domain_BillingMeasure_FlatCastBalanceNoTLG(t *testing.T) {
	tests := map[string]struct {
		b                      *BillingMeasure
		magnitude              measures.Magnitude
		expectedTransformation []BillingLoadCurve
	}{
		"Should Be Okay": {
			magnitude: measures.AI,
			b: &BillingMeasure{
				BillingBalance: BillingBalance{
					P0: &BillingBalancePeriod{
						AI: 614000.0,
					},
				},
				PointType: "1",
				BillingLoadCurve: []BillingLoadCurve{
					{
						Origin: measures.Filled,
						Period: measures.P6,
					},
					{
						AI:        16000.0,
						Period:    measures.P6,
						Equipment: measures.Main,
					},
					{
						AI:        24000.0,
						Period:    measures.P6,
						Equipment: measures.Redundant,
					},
					{
						AI:        25000.0,
						Period:    measures.P6,
						Equipment: measures.Redundant,
					},
					{
						AI:        19000.0,
						Period:    measures.P6,
						Equipment: measures.Redundant,
					},
					{
						AI:        27000.0,
						Period:    measures.P6,
						Equipment: measures.Redundant,
					},
					{
						AI:        29000.0,
						Period:    measures.P6,
						Equipment: measures.Redundant,
					},
					{
						AI:        30000.0,
						Period:    measures.P6,
						Equipment: measures.Redundant,
					},
					{
						AI:        22000.0,
						Period:    measures.P5,
						Equipment: measures.Redundant,
					},
					{
						Origin: measures.Filled,
						Period: measures.P4,
					},
					{
						Origin: measures.Filled,
						Period: measures.P4,
					},
					{
						Origin: measures.Filled,
						Period: measures.P4,
					},
					{
						Origin: measures.Filled,
						Period: measures.P4,
					},
					{
						AI:        21000.0,
						Period:    measures.P4,
						Equipment: measures.Redundant,
					},
					{
						AI:        30000.0,
						Period:    measures.P5,
						Equipment: measures.Redundant,
					},
					{
						AI:        45000.0,
						Period:    measures.P5,
						Equipment: measures.Redundant,
					},
					{
						AI:        20000.0,
						Period:    measures.P5,
						Equipment: measures.Redundant,
					},
					{
						AI:        28000.0,
						Period:    measures.P5,
						Equipment: measures.Redundant,
					},
					{
						AI:        12000.0,
						Period:    measures.P4,
						Equipment: measures.Redundant,
					},
					{
						AI:        23000.0,
						Period:    measures.P4,
						Equipment: measures.Redundant,
					},
					{
						AI:        24000.0,
						Period:    measures.P4,
						Equipment: measures.Redundant,
					},
					{
						AI:        19000.0,
						Period:    measures.P4,
						Equipment: measures.Redundant,
					},
					{
						AI:        20000.0,
						Period:    measures.P5,
						Equipment: measures.Redundant,
					},
					{
						AI:     19000.0,
						Period: measures.P5,
					},
				},
			},
			expectedTransformation: []BillingLoadCurve{
				{
					AI:                       32200.0,
					Origin:                   measures.Filled,
					Period:                   measures.P6,
					MeasureTypeAI:            ProvisionalBalanceMeasure,
					EstimatedCodeAI:          8,
					EstimatedMethodAI:        EstimatedByFlatProfile,
					EstimatedGeneralMethodAI: GeneralEstimated,
				},
				{
					AI:                       16000.0,
					Period:                   measures.P6,
					Equipment:                measures.Main,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          1,
					EstimatedMethodAI:        FirmMainConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       24000.0,
					Period:                   measures.P6,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       25000.0,
					Period:                   measures.P6,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       19000.0,
					Period:                   measures.P6,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       27000.0,
					Period:                   measures.P6,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       29000.0,
					Period:                   measures.P6,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       30000.0,
					Period:                   measures.P6,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       22000.0,
					Period:                   measures.P5,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					Origin:                   measures.Filled,
					AI:                       32200.0,
					Period:                   measures.P4,
					MeasureTypeAI:            ProvisionalBalanceMeasure,
					EstimatedCodeAI:          8,
					EstimatedMethodAI:        EstimatedByFlatProfile,
					EstimatedGeneralMethodAI: GeneralEstimated,
				},
				{
					Origin:                   measures.Filled,
					AI:                       32200.0,
					Period:                   measures.P4,
					MeasureTypeAI:            ProvisionalBalanceMeasure,
					EstimatedCodeAI:          8,
					EstimatedMethodAI:        EstimatedByFlatProfile,
					EstimatedGeneralMethodAI: GeneralEstimated,
				},
				{
					Origin:                   measures.Filled,
					AI:                       32200.0,
					Period:                   measures.P4,
					MeasureTypeAI:            ProvisionalBalanceMeasure,
					EstimatedCodeAI:          8,
					EstimatedMethodAI:        EstimatedByFlatProfile,
					EstimatedGeneralMethodAI: GeneralEstimated,
				},
				{
					Origin:                   measures.Filled,
					AI:                       32200.0,
					Period:                   measures.P4,
					MeasureTypeAI:            ProvisionalBalanceMeasure,
					EstimatedCodeAI:          8,
					EstimatedMethodAI:        EstimatedByFlatProfile,
					EstimatedGeneralMethodAI: GeneralEstimated,
				},
				{
					AI:                       21000.0,
					Period:                   measures.P4,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       30000.0,
					Period:                   measures.P5,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       45000.0,
					Period:                   measures.P5,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       20000.0,
					Period:                   measures.P5,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       28000.0,
					Period:                   measures.P5,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       12000.0,
					Period:                   measures.P4,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       23000.0,
					Period:                   measures.P4,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       24000.0,
					Period:                   measures.P4,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       19000.0,
					Period:                   measures.P4,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       20000.0,
					Period:                   measures.P5,
					Equipment:                measures.Redundant,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
				{
					AI:                       19000.0,
					Period:                   measures.P5,
					MeasureTypeAI:            FirmBalanceMeasure,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
				},
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			flatCast := NewFlatCastBalanceNoTLG(testCase.b, testCase.magnitude)
			flatCast.Execute(context.Background())
			assert.Equal(t, "CCH_FLAT_CAST_BALANCE_NO_TLG", flatCast.ID())
			assert.Equal(t, testCase.expectedTransformation, flatCast.b.BillingLoadCurve)
		})
	}
}

func Test_Unit_Domain_BillingMeasure_BalanceZeroConsumption(t *testing.T) {
	type input struct {
		b         *BillingMeasure
		magnitude measures.Magnitude
	}

	type output struct {
		b BillingBalance
	}

	type want struct {
		err error
		id  string
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be set 0 values in all balance": {
			input: input{
				magnitude: measures.AI,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
					BillingBalance: BillingBalance{
						P0: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
						P6: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
						P5: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
						P4: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
					},
				},
			},
			output: output{
				b: BillingBalance{
					P0: &BillingBalancePeriod{
						AI:                   0,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        EstimateOnlyHistoric,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
						BalanceValidationAI:  measures.Invalid,
					},
					P6: &BillingBalancePeriod{
						AI:                   0,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        EstimateOnlyHistoric,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
						BalanceValidationAI:  measures.Invalid,
					},
					P5: &BillingBalancePeriod{
						AI:                   0,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        EstimateOnlyHistoric,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
						BalanceValidationAI:  measures.Invalid,
					},
					P4: &BillingBalancePeriod{
						AI:                   0,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        EstimateOnlyHistoric,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
						BalanceValidationAI:  measures.Invalid,
					},
				},
			},
			want: want{
				err: nil,
				id:  "BALANCE_ZERO_CONSUMPTION",
			},
		},
	}

	for name := range testCases {
		test := testCases[name]

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			algorithm := NewBalanceZeroConsumption(test.input.b, test.input.magnitude)

			err := algorithm.Execute(ctx)

			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.id, algorithm.ID())
			assert.Equal(t, test.output.b, algorithm.b.BillingBalance)
		})
	}
}

func Test_Unit_Domain_BillingMeasure_ConsumZeroClosedHouseNoTLG(t *testing.T) {
	tests := map[string]struct {
		b                      *BillingMeasure
		expectedTransformation []BillingLoadCurve
		magnitude              measures.Magnitude
	}{
		"Should Be Okay PointType 3": {
			b: &BillingMeasure{
				PointType: "3",
				BillingLoadCurve: []BillingLoadCurve{
					{Origin: measures.Filled},
				},
			},
			expectedTransformation: []BillingLoadCurve{
				{
					AI:                       0.0,
					Origin:                   measures.Filled,
					MeasureTypeAI:            ProvisionalBalanceMeasure,
					EstimatedCodeAI:          5,
					EstimatedMethodAI:        EstimateByHistoricConsumLastYear,
					EstimatedGeneralMethodAI: GeneralEstimated,
				},
			},
			magnitude: measures.AI,
		},
		"Should Be Okay PointType 2": {
			b: &BillingMeasure{
				PointType: "2",
				BillingLoadCurve: []BillingLoadCurve{
					{Origin: measures.Filled},
				},
			},
			expectedTransformation: []BillingLoadCurve{
				{
					AE:                       0.0,
					Origin:                   measures.Filled,
					MeasureTypeAE:            ProvisionalBalanceMeasure,
					EstimatedCodeAE:          9,
					EstimatedMethodAE:        EstimateHistoricMainMeasurePoint,
					EstimatedGeneralMethodAE: GeneralEstimated,
				},
			},
			magnitude: measures.AE,
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			consum := NewConsumZeroClosedHouseNoTLG(testCase.b, testCase.magnitude)
			consum.Execute(context.Background())
			assert.Equal(t, "CCH_CONSUM_ZERO_CLOSED_HOUSE_NO_TLG", consum.ID())
			assert.Equal(t, testCase.expectedTransformation, consum.b.BillingLoadCurve)
		})
	}
}

func Test_Unit_Domain_BillingMeasure_CloseHistoryWithoutBalance(t *testing.T) {

	type input struct {
		magnitude     measures.Magnitude
		b             *BillingMeasure
		context       *GraphContext
		daysOfPeriods map[measures.PeriodKey][]int
	}

	type output struct {
		b BillingBalance
	}

	type want struct {
		id  string
		err error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be calculated correctly balance, STM origin": {
			input: input{
				daysOfPeriods: map[measures.PeriodKey][]int{
					measures.P5: {63, 210, 105, 112, 154},
					measures.P4: {81, 270, 135, 144, 198},
				},
				magnitude: measures.AI,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
					BillingBalance: BillingBalance{
						Origin: measures.STM,
						P0: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
						P6: &BillingBalancePeriod{
							AI:                  200000,
							BalanceValidationAI: measures.Valid,
						},
						P5: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
						P4: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
					},
				},
				context: &GraphContext{
					ClosedHistory: []BillingMeasure{
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 88000,
								},
								P5: &BillingBalancePeriod{
									AI: 130000,
								},
								P6: &BillingBalancePeriod{
									AI: 220000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 286000,
								},
								P5: &BillingBalancePeriod{
									AI: 322000,
								},
								P6: &BillingBalancePeriod{
									AI: 185000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 157000,
								},
								P5: &BillingBalancePeriod{
									AI: 190000,
								},
								P6: &BillingBalancePeriod{
									AI: 200000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 178000,
								},
								P5: &BillingBalancePeriod{
									AI: 210000,
								},
								P6: &BillingBalancePeriod{
									AI: 205000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
					},
				},
			},
			output: output{
				b: BillingBalance{
					Origin: measures.STM,
					P0: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Invalid,
						AI:                   705297,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        CalculatedByCloseSum,
						BalanceGeneralTypeAI: GeneralCalculated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
					P6: &BillingBalancePeriod{
						AI:                   200000,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      1,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      TlmOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P5: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Invalid,
						AI:                   280332,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        EstimateByHistoricLastYear,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
					P4: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Invalid,
						AI:                   224965,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        EstimateByHistoricLastYear,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
				},
			},
			want: want{
				id:  "BALANCE_CLOSE_HISTORIC_WITHOUT_BALANCE",
				err: nil,
			},
		},
		"Should be calculated correctly balance no curve hours 0 values, TPL origin": {
			input: input{
				daysOfPeriods: map[measures.PeriodKey][]int{},
				magnitude:     measures.AI,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
					BillingBalance: BillingBalance{
						Origin: measures.TPL,
						P0: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
						P6: &BillingBalancePeriod{
							AI:                  200000,
							BalanceValidationAI: measures.Valid,
						},
						P5: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
						P4: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
					},
				},
				context: &GraphContext{
					ClosedHistory: []BillingMeasure{
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 88000,
								},
								P5: &BillingBalancePeriod{
									AI: 130000,
								},
								P6: &BillingBalancePeriod{
									AI: 220000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 286000,
								},
								P5: &BillingBalancePeriod{
									AI: 322000,
								},
								P6: &BillingBalancePeriod{
									AI: 185000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 157000,
								},
								P5: &BillingBalancePeriod{
									AI: 190000,
								},
								P6: &BillingBalancePeriod{
									AI: 200000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 178000,
								},
								P5: &BillingBalancePeriod{
									AI: 210000,
								},
								P6: &BillingBalancePeriod{
									AI: 205000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
					},
				},
			},
			output: output{
				b: BillingBalance{
					Origin: measures.TPL,
					P0: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Invalid,
						AI:                   200000,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        CalculatedByCloseSum,
						BalanceGeneralTypeAI: GeneralCalculated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
					P6: &BillingBalancePeriod{
						AI:                   200000,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        RealByAbsLocalRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P5: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Invalid,
						AI:                   0,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        EstimateByHistoricLastYear,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
					P4: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Invalid,
						AI:                   0,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        EstimateByHistoricLastYear,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
				},
			},
			want: want{
				id:  "BALANCE_CLOSE_HISTORIC_WITHOUT_BALANCE",
				err: nil,
			},
		},
	}

	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			//GENERATE CURVE FALSE CURVE
			for period, listOfdays := range test.input.daysOfPeriods {
				for i, days := range listOfdays {
					billingMeasure := test.input.b

					if i < len(listOfdays)-1 {
						billingMeasure = &test.input.context.ClosedHistory[i]
					}

					for i := 0; i < days; i++ {
						billingMeasure.BillingLoadCurve = append(billingMeasure.BillingLoadCurve, BillingLoadCurve{Period: period})
					}
				}
			}

			algorithm := NewCloseHistoryWithoutBalance(test.input.b, test.input.magnitude, test.input.context)
			err := algorithm.Execute(ctx)

			assert.Equal(t, test.output.b, algorithm.b.BillingBalance)
			assert.Equal(t, test.want.id, algorithm.ID())
			assert.Equal(t, test.want.err, err)
		})
	}
}

func Test_Unit_Domain_BillingMeasure_CloseHistoryWithBalance(t *testing.T) {

	type input struct {
		magnitude     measures.Magnitude
		b             *BillingMeasure
		context       *GraphContext
		daysOfPeriods map[measures.PeriodKey][]int
	}

	type output struct {
		b BillingBalance
	}

	type want struct {
		id  string
		err error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be calculated correctly balance, STM origin": {
			input: input{
				daysOfPeriods: map[measures.PeriodKey][]int{
					measures.P5: {63, 210, 105, 112, 154},
					measures.P4: {81, 270, 135, 144, 198},
				},
				magnitude: measures.AI,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
					BillingBalance: BillingBalance{
						Origin: measures.STM,
						P0: &BillingBalancePeriod{
							AI:                  800000,
							BalanceValidationAI: measures.Valid,
						},
						P6: &BillingBalancePeriod{
							AI:                  200000,
							BalanceValidationAI: measures.Valid,
						},
						P5: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
						P4: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
					},
				},
				context: &GraphContext{
					ClosedHistory: []BillingMeasure{
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 88000,
								},
								P5: &BillingBalancePeriod{
									AI: 130000,
								},
								P6: &BillingBalancePeriod{
									AI: 220000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 286000,
								},
								P5: &BillingBalancePeriod{
									AI: 322000,
								},
								P6: &BillingBalancePeriod{
									AI: 185000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 157000,
								},
								P5: &BillingBalancePeriod{
									AI: 190000,
								},
								P6: &BillingBalancePeriod{
									AI: 200000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 178000,
								},
								P5: &BillingBalancePeriod{
									AI: 210000,
								},
								P6: &BillingBalancePeriod{
									AI: 205000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
					},
				},
			},
			output: output{
				b: BillingBalance{
					Origin: measures.STM,
					P0: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Valid,
						AI:                   800000,
						EstimatedCodeAI:      1,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      TlmOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P6: &BillingBalancePeriod{
						AI:                   200000,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      1,
						BalanceTypeAI:        RealByRemoteRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      TlmOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P5: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Invalid,
						AI:                   332872,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        EstimateByHistoricLastYear,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
					P4: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Invalid,
						AI:                   267128,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        EstimateByHistoricLastYear,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
				},
			},
			want: want{
				id:  "BALANCE_CLOSE_HISTORY_WITH_BALANCE",
				err: nil,
			},
		},
		"Should be calculated correctly balance no curve hours 0 values, TPL origin": {
			input: input{
				daysOfPeriods: map[measures.PeriodKey][]int{},
				magnitude:     measures.AI,
				b: &BillingMeasure{
					Periods: []measures.PeriodKey{measures.P6, measures.P5, measures.P4},
					BillingBalance: BillingBalance{
						Origin: measures.TPL,
						P0: &BillingBalancePeriod{
							AI:                  800000,
							BalanceValidationAI: measures.Valid,
						},
						P6: &BillingBalancePeriod{
							AI:                  200000,
							BalanceValidationAI: measures.Valid,
						},
						P5: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
						P4: &BillingBalancePeriod{
							BalanceValidationAI: measures.Invalid,
						},
					},
				},
				context: &GraphContext{
					ClosedHistory: []BillingMeasure{
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 88000,
								},
								P5: &BillingBalancePeriod{
									AI: 130000,
								},
								P6: &BillingBalancePeriod{
									AI: 220000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 286000,
								},
								P5: &BillingBalancePeriod{
									AI: 322000,
								},
								P6: &BillingBalancePeriod{
									AI: 185000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 157000,
								},
								P5: &BillingBalancePeriod{
									AI: 190000,
								},
								P6: &BillingBalancePeriod{
									AI: 200000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
						{
							BillingBalance: BillingBalance{
								P4: &BillingBalancePeriod{
									AI: 178000,
								},
								P5: &BillingBalancePeriod{
									AI: 210000,
								},
								P6: &BillingBalancePeriod{
									AI: 205000,
								},
							},

							BillingLoadCurve: []BillingLoadCurve{},
						},
					},
				},
			},
			output: output{
				b: BillingBalance{
					Origin: measures.TPL,
					P0: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Valid,
						AI:                   800000,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        RealByAbsLocalRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P6: &BillingBalancePeriod{
						AI:                   200000,
						BalanceValidationAI:  measures.Valid,
						EstimatedCodeAI:      3,
						BalanceTypeAI:        RealByAbsLocalRead,
						BalanceGeneralTypeAI: GeneralReal,
						BalanceOriginAI:      LocalOrigin,
						BalanceMeasureTypeAI: FirmBalanceMeasure,
					},
					P5: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Invalid,
						AI:                   0,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        EstimateByHistoricLastYear,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
					P4: &BillingBalancePeriod{
						BalanceValidationAI:  measures.Invalid,
						AI:                   0,
						EstimatedCodeAI:      5,
						BalanceTypeAI:        EstimateByHistoricLastYear,
						BalanceGeneralTypeAI: GeneralEstimated,
						BalanceOriginAI:      EstimateOrigin,
						BalanceMeasureTypeAI: ProvisionalBalanceMeasure,
					},
				},
			},
			want: want{
				id:  "BALANCE_CLOSE_HISTORY_WITH_BALANCE",
				err: nil,
			},
		},
	}

	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			//GENERATE CURVE FALSE CURVE
			for period, listOfdays := range test.input.daysOfPeriods {
				for i, days := range listOfdays {
					billingMeasure := test.input.b

					if i < len(listOfdays)-1 {
						billingMeasure = &test.input.context.ClosedHistory[i]
					}

					for i := 0; i < days; i++ {
						billingMeasure.BillingLoadCurve = append(billingMeasure.BillingLoadCurve, BillingLoadCurve{Period: period})
					}
				}
			}

			algorithm := NewClosingHistoryWithBalance(test.input.b, test.input.magnitude, test.input.context)
			err := algorithm.Execute(ctx)

			assert.Equal(t, test.output.b, algorithm.B.BillingBalance)
			assert.Equal(t, test.want.id, algorithm.ID())
			assert.Equal(t, test.want.err, err)
		})
	}
}

func Test_Unit_Domain_BillingMeasure_CCHPenalty(t *testing.T) {
	type input struct {
		ctx       context.Context
		b         *BillingMeasure
		magnitude measures.Magnitude
		period    measures.PeriodKey
	}
	type output struct {
		blc []BillingLoadCurve
	}
	tests := map[string]struct {
		input  input
		output output
	}{
		"Point type 3, 4 or 5 case": {
			input: input{
				magnitude: measures.AI,
				ctx:       context.Background(),
				b: &BillingMeasure{
					PointType: "3",
					BillingLoadCurve: []BillingLoadCurve{
						{
							Origin: measures.STM,
						},
						{
							Origin: measures.Filled,
						},
					},
				},
			},
			output: output{blc: []BillingLoadCurve{
				{
					Origin:                   measures.STM,
					EstimatedCodeAI:          1,
					EstimatedMethodAI:        RealValidMeasure,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Origin:                   measures.Filled,
					EstimatedCodeAI:          5,
					EstimatedMethodAI:        EstimateByHistoricLastYear,
					EstimatedGeneralMethodAI: GeneralEstimated,
					MeasureTypeAI:            ProvisionalBalanceMeasure,
				},
			}},
		},
		"Point type 1 or 2 case": {
			input: input{
				ctx: context.Background(),
				b: &BillingMeasure{
					PointType: "1",
					BillingLoadCurve: []BillingLoadCurve{
						{
							Origin:    measures.STM,
							Equipment: measures.Main,
						},
						{
							Origin:    measures.STM,
							Equipment: measures.Redundant,
						},
						{
							Origin:    measures.STM,
							Equipment: measures.Receipt,
						},
						{
							Origin: measures.Filled,
						},
					},
				},
				magnitude: measures.AI,
			},
			output: output{blc: []BillingLoadCurve{
				{
					Origin:                   measures.STM,
					Equipment:                measures.Main,
					EstimatedCodeAI:          1,
					EstimatedMethodAI:        FirmMainConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Origin:                   measures.STM,
					Equipment:                measures.Redundant,
					EstimatedCodeAI:          2,
					EstimatedMethodAI:        FirmRedundantConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Origin:                   measures.STM,
					Equipment:                measures.Receipt,
					EstimatedCodeAI:          3,
					EstimatedMethodAI:        FirmReceiptConfig,
					EstimatedGeneralMethodAI: GeneralReal,
					MeasureTypeAI:            FirmBalanceMeasure,
				},
				{
					Origin:                   measures.Filled,
					EstimatedCodeAI:          22,
					EstimatedMethodAI:        EstimateWhosePenaltiesForClientsTypeOneAndTwo,
					EstimatedGeneralMethodAI: GeneralEstimated,
					MeasureTypeAI:            ProvisionalBalanceMeasure,
				},
			}},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			cchPenalty := NewCCHPenalty(testCase.input.b, testCase.input.magnitude)
			cchPenalty.Execute(testCase.input.ctx)
			assert.Equal(t, "CCH_PENALTY", cchPenalty.ID())
			assert.Equal(t, testCase.output.blc, cchPenalty.B.BillingLoadCurve)
		})
	}
}
