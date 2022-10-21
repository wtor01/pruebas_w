package tests

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/tests/fixtures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Unit_Domain_BillingMeasures_HaveClosureMeasure(t *testing.T) {

	type input struct {
		b *billing_measures.BillingMeasure
	}

	testCases := map[string]struct {
		input input
		want  bool
	}{
		"should be true, if its completed": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					BillingBalance: billing_measures.BillingBalance{
						P1: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P2: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P3: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
					},
				},
			},
			want: true,
		},
		"should be false, if its not completed": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					BillingBalance: billing_measures.BillingBalance{
						P1: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P2: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
						P3: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
					},
				},
			},
			want: false,
		},
		"should be true, if its completed example": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4, measures.P5, measures.P6},
					BillingBalance: billing_measures.BillingBalance{
						P1: &billing_measures.BillingBalancePeriod{AI: 2436, BalanceValidationAI: measures.Valid},
						P2: &billing_measures.BillingBalancePeriod{AI: 3142, BalanceValidationAI: measures.Valid},
						P3: &billing_measures.BillingBalancePeriod{AI: 0, BalanceValidationAI: measures.Valid},
						P4: &billing_measures.BillingBalancePeriod{AI: 0, BalanceValidationAI: measures.Valid},
						P5: &billing_measures.BillingBalancePeriod{AI: 7124, BalanceValidationAI: measures.Valid},
						P6: &billing_measures.BillingBalancePeriod{AI: 4, BalanceValidationAI: measures.Valid},
					},
				},
			},
			want: true,
		},
		"should be false, if its not completed example": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3, measures.P4, measures.P5, measures.P6},
					BillingBalance: billing_measures.BillingBalance{
						P1: &billing_measures.BillingBalancePeriod{AI: 2436, BalanceValidationAI: measures.Valid},
						P2: &billing_measures.BillingBalancePeriod{AI: 3142, BalanceValidationAI: measures.Valid},
						P3: &billing_measures.BillingBalancePeriod{AI: 0, BalanceValidationAI: measures.Invalid},
						P4: &billing_measures.BillingBalancePeriod{AI: 0, BalanceValidationAI: measures.Invalid},
						P5: &billing_measures.BillingBalancePeriod{AI: 7124, BalanceValidationAI: measures.Valid},
						P6: &billing_measures.BillingBalancePeriod{AI: 4, BalanceValidationAI: measures.Valid},
					},
				},
			},
			want: false,
		},
	}

	for name, _ := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			haveBillingHistory := billing_measures.NewIsCloseMeasureComplete(test.input.b, measures.AI)

			result := haveBillingHistory.Eval(ctx)
			assert.Equal(t, test.want, result)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsMoreThanOneMissingPeriod(t *testing.T) {

	type input struct {
		b *billing_measures.BillingMeasure
	}

	testCases := map[string]struct {
		input input
		want  bool
	}{
		"should be false, if its not missing periods": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					BillingBalance: billing_measures.BillingBalance{
						P1: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P2: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P3: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
					},
				},
			},
			want: false,
		},
		"should be false, if its only one missing periods": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					BillingBalance: billing_measures.BillingBalance{
						P0: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
						P1: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P2: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
						P3: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
					},
				},
			},
			want: false,
		},
		"should be true, there are two or more missing periods": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					BillingBalance: billing_measures.BillingBalance{
						P1: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P2: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
						P3: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
					},
				},
			},
			want: true,
		},
		"should be true, if its completed example": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					BillingBalance: billing_measures.BillingBalance{
						P1: &billing_measures.BillingBalancePeriod{AI: 12702, BalanceValidationAI: measures.Valid},
						P2: &billing_measures.BillingBalancePeriod{AI: 0, BalanceValidationAI: measures.Invalid},
						P3: &billing_measures.BillingBalancePeriod{AI: 0, BalanceValidationAI: measures.Invalid},
					},
				},
			},
			want: true,
		},
		"should be false, there are not enough missing periods example": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					BillingBalance: billing_measures.BillingBalance{
						P0: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
						P1: &billing_measures.BillingBalancePeriod{AI: 2436, BalanceValidationAI: measures.Valid},
						P2: &billing_measures.BillingBalancePeriod{AI: 3142, BalanceValidationAI: measures.Valid},
						P3: &billing_measures.BillingBalancePeriod{AI: 7124, BalanceValidationAI: measures.Valid},
					},
				},
			},
			want: false,
		},
	}

	for name, _ := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			haveBillingHistory := billing_measures.NewIsMoreThanOneMissingPeriod(test.input.b, measures.AI)

			result := haveBillingHistory.Eval(ctx)
			assert.Equal(t, test.want, result)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsHourly(t *testing.T) {

	type input struct {
		b *billing_measures.BillingMeasure
	}

	testCases := map[string]struct {
		input input
		want  bool
	}{
		"Should be true with measure hourly": {
			input: input{b: &billing_measures.BillingMeasure{Technology: "R", RegisterType: measures.Hourly}},
			want:  true,
		},
		"Should be true with measure both": {
			input: input{b: &billing_measures.BillingMeasure{Technology: "R", RegisterType: measures.Both}},
			want:  true,
		},
		"Should be false with measure quarter": {
			input: input{b: &billing_measures.BillingMeasure{Technology: "R", RegisterType: measures.QuarterHour}},
			want:  false,
		},
		"Should be false with measure None": {
			input: input{b: &billing_measures.BillingMeasure{Technology: "R", RegisterType: measures.NoneType}},
			want:  false,
		},
		"Should be true with technology T": {
			input: input{b: &billing_measures.BillingMeasure{Technology: "T", RegisterType: measures.Hourly}},
			want:  true,
		},
	}

	for name, _ := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			condition := billing_measures.NewIsHourly(test.input.b)
			evalResult := condition.Eval(ctx)
			assert.Equal(t, test.want, evalResult)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_HasCloseHistories(t *testing.T) {
	type input struct {
		b            *billing_measures.BillingMeasure
		magnitude    measures.Magnitude
		ContextNoTlg billing_measures.GraphContext
	}

	type output struct {
		billingHistory  []billing_measures.BillingMeasure
		getHistoriesErr error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   bool
	}{
		"Should be true, have 4": {
			input: input{b: &billing_measures.BillingMeasure{Periods: []measures.PeriodKey{measures.P4, measures.P6, measures.P5}}, magnitude: measures.AI},
			output: output{
				billingHistory: []billing_measures.BillingMeasure{
					{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"},
				},
				getHistoriesErr: nil,
			},
			want: true,
		},
		"Should be true, have 5": {
			input: input{b: &billing_measures.BillingMeasure{Periods: []measures.PeriodKey{measures.P4, measures.P6, measures.P5}}, magnitude: measures.AI},
			output: output{
				billingHistory: []billing_measures.BillingMeasure{
					{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"},
				},
				getHistoriesErr: nil,
			},
			want: true,
		},
		"Should be false, have 3": {
			input: input{b: &billing_measures.BillingMeasure{Periods: []measures.PeriodKey{measures.P4, measures.P6, measures.P5}}, magnitude: measures.AI},
			output: output{
				billingHistory: []billing_measures.BillingMeasure{
					{Id: "1"}, {Id: "2"}, {Id: "3"},
				},
				getHistoriesErr: nil,
			},
			want: false,
		},
		"Should be false, err": {
			input: input{b: &billing_measures.BillingMeasure{Periods: []measures.PeriodKey{measures.P4, measures.P6, measures.P5}}, magnitude: measures.AI},
			output: output{
				billingHistory:  []billing_measures.BillingMeasure{},
				getHistoriesErr: errors.New("err"),
			},
			want: false,
		},
	}

	for name, _ := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			billingRepository := mocks.NewBillingMeasureRepository(t)

			billingRepository.On("GetCloseHistories", ctx, billing_measures.QueryGetCloseHistories{
				CUPS:       test.input.b.CUPS,
				EndDate:    test.input.b.EndDate,
				Periods:    test.input.b.Periods,
				Magnitudes: []measures.Magnitude{test.input.magnitude},
			}).Return(test.output.billingHistory, test.output.getHistoriesErr)

			condition := billing_measures.NewHasClosedHistory(test.input.b, billingRepository, test.input.magnitude, &test.input.ContextNoTlg)
			result := condition.Eval(ctx)

			assert.ElementsMatch(t, condition.ContextNoTlg.ClosedHistory, test.output.billingHistory)

			assert.Equal(t, test.want, result)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_HasAnualHistory(t *testing.T) {

	type input struct {
		b            *billing_measures.BillingMeasure
		magnitude    measures.Magnitude
		ContextNoTlg *billing_measures.GraphContext
	}

	type output struct {
		b              billing_measures.BillingMeasure
		lastHistoryErr error
	}

	type want struct {
		expected bool
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be true, have history with real balance": {
			input: input{
				b:            &billing_measures.BillingMeasure{},
				magnitude:    measures.AE,
				ContextNoTlg: &billing_measures.GraphContext{},
			},
			output: output{
				b: billing_measures.BillingMeasure{
					Id: "1",
					BillingBalance: billing_measures.BillingBalance{
						P0: &billing_measures.BillingBalancePeriod{
							BalanceTypeAE: billing_measures.RealBalance,
						},
					},
				},
				lastHistoryErr: nil,
			},
			want: want{
				expected: true,
			},
		},
		"Should be false, err": {
			input: input{
				b:            &billing_measures.BillingMeasure{},
				magnitude:    measures.AE,
				ContextNoTlg: &billing_measures.GraphContext{},
			},
			output: output{
				b: billing_measures.BillingMeasure{
					Id: "1",
					BillingBalance: billing_measures.BillingBalance{
						P0: &billing_measures.BillingBalancePeriod{
							BalanceTypeAE: billing_measures.RealBalance,
						},
					},
				},
				lastHistoryErr: errors.New("err"),
			},
			want: want{
				expected: false,
			},
		},
	}

	for name, _ := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			billingRepository := mocks.NewBillingMeasureRepository(t)
			billingRepository.On("LastHistory", ctx, billing_measures.QueryLastHistory{
				CUPS:       test.input.b.CUPS,
				InitDate:   test.input.b.InitDate,
				EndDate:    test.input.b.EndDate,
				Periods:    test.input.b.Periods,
				Magnitudes: []measures.Magnitude{test.input.magnitude},
			}).Return(test.output.b, test.output.lastHistoryErr)
			condition := billing_measures.NewHasAnualHistory(test.input.b, billingRepository, test.input.magnitude, test.input.ContextNoTlg)

			result := condition.Eval(ctx)

			assert.Equal(t, test.want.expected, result)
			assert.Equal(t, test.input.ContextNoTlg, condition.ContextNoTlg)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsSimpleHistoric(t *testing.T) {

	type input struct {
		b                 *billing_measures.BillingMeasure
		previousLoadCurve []process_measures.ProcessedLoadCurve
		nextLoadCurve     []process_measures.ProcessedLoadCurve
	}

	testCases := map[string]struct {
		input input
		want  bool
	}{
		"should be false, if missing historic": {
			input: input{
				b: &fixtures.BillingResult_1_Input_averages_conditional,
			},
			want: false,
		},
		"should be false, if last value of historic its filled": {
			input: input{
				b:                 &fixtures.BillingResult_1_Input_averages_conditional,
				previousLoadCurve: fixtures.ProcessMeasures_Input_averages_conditional_previous_filled,
			},
			want: false,
		},
		"should be true, if last value of historic its not filled": {
			input: input{
				b:                 &fixtures.BillingResult_1_Input_averages_conditional,
				previousLoadCurve: fixtures.ProcessMeasures_Input_averages_conditional_previous_not_filled,
			},
			want: true,
		},
		"should be false, if first value of historic its filled": {
			input: input{
				b:             &fixtures.BillingResult_2_Input_averages_conditional,
				nextLoadCurve: fixtures.ProcessMeasures_Input_averages_conditional_next_filled,
			},
			want: false,
		},
		"should be true, if first value of historic its not filled": {
			input: input{
				b:             &fixtures.BillingResult_2_Input_averages_conditional,
				nextLoadCurve: fixtures.ProcessMeasures_Input_averages_conditional_next_not_filled,
			},
			want: true,
		},
		"should be true, valid example task": {
			input: input{
				b: &fixtures.BillingResult_example_Input_averages_conditional_true,
			},
			want: true,
		},
		"should be false, invalid example task": {
			input: input{
				b: &fixtures.BillingResult_example_Input_averages_conditional_false,
			},
			want: false,
		},
	}

	for name, _ := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			processedMeasuresRepository := new(mocks.ProcessedMeasureRepository)

			processedMeasuresRepository.On("ListHistoryLoadCurve", ctx, process_measures.QueryHistoryLoadCurve{
				CUPS:      test.input.b.CUPS,
				StartDate: test.input.b.InitDate.AddDate(0, 0, -1),
				EndDate:   test.input.b.InitDate,
				Count:     1,
				Periods:   test.input.b.GetPeriods(),
				Magnitude: measures.AI,
				IsFuture:  false,
			}).Once().Return(test.input.previousLoadCurve, nil)
			processedMeasuresRepository.On("ListHistoryLoadCurve", ctx, process_measures.QueryHistoryLoadCurve{
				CUPS:      test.input.b.CUPS,
				StartDate: test.input.b.EndDate,
				EndDate:   test.input.b.EndDate.AddDate(0, 0, 1),
				Count:     1,
				Periods:   test.input.b.GetPeriods(),
				Magnitude: measures.AI,
				IsFuture:  true,
			}).Once().Return(test.input.nextLoadCurve, nil)
			haveBillingHistory := billing_measures.NewIsSimpleHistoric(test.input.b, &billing_measures.GraphContext{}, measures.AI, processedMeasuresRepository)

			result := haveBillingHistory.Eval(ctx)
			assert.Equal(t, test.want, result)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsBalanceValid(t *testing.T) {

	type input struct {
		b *billing_measures.BillingMeasure
	}

	testCases := map[string]struct {
		input input
		want  bool
	}{
		"should be true, if its valid": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					BillingBalance: billing_measures.BillingBalance{
						P0: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P1: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
						P2: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
						P3: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
					},
				},
			},
			want: true,
		},
		"should be false, if its not valid": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods: []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					BillingBalance: billing_measures.BillingBalance{
						P0: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
						P1: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P2: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P3: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
					},
				},
			},
			want: false,
		},
	}

	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			haveBillingHistory := billing_measures.NewIsBalanceValid(test.input.b, measures.AI)

			result := haveBillingHistory.Eval(ctx)
			assert.Equal(t, result, test.want)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_HasIterativeHistory_Eval(t *testing.T) {

	type input struct {
		b            *billing_measures.BillingMeasure
		ContextNoTlg *billing_measures.GraphContext
	}

	type output struct {
		listHistoryLoadCurveResult    []process_measures.ProcessedLoadCurve
		listHistoryLoadCurveResultErr error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   bool
	}{
		"should be false, if listHistoryLoadCurve return err": {
			input: input{
				b: &billing_measures.BillingMeasure{
					CUPS:    "CUPS",
					EndDate: time.Now(),
					Periods: []measures.PeriodKey{measures.P1},
				},
				ContextNoTlg: &billing_measures.GraphContext{},
			},
			output: output{
				listHistoryLoadCurveResult:    nil,
				listHistoryLoadCurveResultErr: errors.New(""),
			},
			want: false,
		},
		"should be false, if listHistoryLoadCurve return less than minCount by period": {
			input: input{
				b: &billing_measures.BillingMeasure{
					CUPS:    "CUPS",
					EndDate: time.Now(),
					Periods: []measures.PeriodKey{measures.P1, measures.P2},
				},
				ContextNoTlg: &billing_measures.GraphContext{},
			},
			output: output{
				listHistoryLoadCurveResult:    make([]process_measures.ProcessedLoadCurve, 3),
				listHistoryLoadCurveResultErr: nil,
			},
			want: false,
		}, "should be true, if listHistoryLoadCurve return minCount by period": {
			input: input{
				b: &billing_measures.BillingMeasure{
					CUPS:    "CUPS",
					EndDate: time.Now(),
					Periods: []measures.PeriodKey{measures.P1, measures.P2},
				},
				ContextNoTlg: &billing_measures.GraphContext{},
			},
			output: output{
				listHistoryLoadCurveResult:    make([]process_measures.ProcessedLoadCurve, 12),
				listHistoryLoadCurveResultErr: nil,
			},
			want: true,
		},
	}

	for name, _ := range testCases {
		testCase := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			processedMeasuresRepository := new(mocks.ProcessedMeasureRepository)

			processedMeasuresRepository.On("ListHistoryLoadCurve", ctx, process_measures.QueryHistoryLoadCurve{
				CUPS:          testCase.input.b.CUPS,
				StartDate:     testCase.input.b.InitDate.AddDate(0, -6, 0),
				EndDate:       testCase.input.b.InitDate,
				WithCriterias: true,
				Periods:       testCase.input.b.Periods,
				Magnitude:     measures.AI,
			}).Return(testCase.output.listHistoryLoadCurveResult, testCase.output.listHistoryLoadCurveResultErr)

			processedMeasuresRepository.On("ListHistoryLoadCurve", ctx, process_measures.QueryHistoryLoadCurve{
				CUPS:          testCase.input.b.CUPS,
				StartDate:     testCase.input.b.EndDate,
				EndDate:       testCase.input.b.EndDate.AddDate(0, 6, 0),
				WithCriterias: true,
				IsFuture:      true,
				Periods:       testCase.input.b.Periods,
				Magnitude:     measures.AI,
			}).Return(testCase.output.listHistoryLoadCurveResult, testCase.output.listHistoryLoadCurveResultErr)
			haveBillingHistory := billing_measures.NewHasIterativeHistory(testCase.input.b, processedMeasuresRepository, measures.AI, testCase.input.ContextNoTlg)

			result := haveBillingHistory.Eval(ctx)
			assert.Equal(t, testCase.want, result)
			assert.Equal(t, testCase.input.ContextNoTlg, haveBillingHistory.ContextNoTlg)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsHouseCloseAndCloseAtrAreEmpty(t *testing.T) {

	type input struct {
		b *billing_measures.BillingMeasure
	}

	type want struct {
		expected bool
	}

	testCases := map[string]struct {
		input input
		want  want
	}{
		"Should be true, house closed - close atr empty - point type 3/4": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods:      []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					PointType:    "3",
					Inaccessible: true,
					BillingBalance: billing_measures.BillingBalance{
						P1: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
						P2: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
						P3: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
					},
				},
			},
			want: want{
				expected: true,
			},
		},
		"Should be false, point type 1": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods:      []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					PointType:    "1",
					Inaccessible: true,
					BillingBalance: billing_measures.BillingBalance{
						P1: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P2: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P3: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
					},
				},
			},
		},
		"Should be false, house not close": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods:      []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					PointType:    "3",
					Inaccessible: false,
					BillingBalance: billing_measures.BillingBalance{
						P1: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P2: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P3: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
					},
				},
			},
		},
		"Should be false, close atr are not empty": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Periods:      []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					PointType:    "3",
					Inaccessible: true,
					Magnitudes:   []measures.Magnitude{measures.AI},
					BillingBalance: billing_measures.BillingBalance{
						P1: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Valid},
						P2: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
						P3: &billing_measures.BillingBalancePeriod{BalanceValidationAI: measures.Invalid},
					},
				},
			},
			want: want{
				expected: false,
			},
		},
	}

	for name, _ := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			condition := billing_measures.NewIsHouseCloseAndCloseAtrAreEmpty(test.input.b, measures.AI)

			result := condition.Eval(ctx)

			assert.Equal(t, test.want.expected, result)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsCCCHComplete_Eval(t *testing.T) {

	type input struct {
		b            *billing_measures.BillingMeasure
		ContextNoTlg *billing_measures.GraphContext
	}

	type output struct {
		listHistoryLoadCurveResult    []process_measures.ProcessedLoadCurve
		listHistoryLoadCurveResultErr error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   bool
	}{
		"should be false, if listHistoryLoadCurve return err": {
			input: input{
				b: &billing_measures.BillingMeasure{
					CUPS:         "CUPS",
					EndDate:      time.Now(),
					Periods:      []measures.PeriodKey{measures.P1},
					RegisterType: measures.Hourly,
				},
				ContextNoTlg: &billing_measures.GraphContext{},
			},
			output: output{
				listHistoryLoadCurveResult:    nil,
				listHistoryLoadCurveResultErr: errors.New(""),
			},
			want: false,
		}, "should be false, if load curve of billing measure dont have hourly type": {
			input: input{
				b: &billing_measures.BillingMeasure{
					CUPS:         "CUPS",
					EndDate:      time.Now(),
					Periods:      []measures.PeriodKey{measures.P1},
					RegisterType: measures.QuarterHour,
				},
				ContextNoTlg: &billing_measures.GraphContext{},
			},
			output: output{
				listHistoryLoadCurveResult:    nil,
				listHistoryLoadCurveResultErr: nil,
			},
			want: false,
		}, "should be false, if dont exist quarterly values history": {
			input: input{
				b: &billing_measures.BillingMeasure{
					CUPS:         "CUPS",
					EndDate:      time.Now(),
					Periods:      []measures.PeriodKey{measures.P1},
					RegisterType: measures.Hourly,
				},
				ContextNoTlg: &billing_measures.GraphContext{},
			},
			output: output{
				listHistoryLoadCurveResult:    []process_measures.ProcessedLoadCurve{},
				listHistoryLoadCurveResultErr: nil,
			},
			want: false,
		}, "should be true, if exist quarterly values history": {
			input: input{
				b:            &fixtures.BillingResult_example_Input_ccch_complete_1,
				ContextNoTlg: &billing_measures.GraphContext{},
			},
			output: output{
				listHistoryLoadCurveResult:    fixtures.ProcessedLoadCurve_ccch_example_input_1,
				listHistoryLoadCurveResultErr: nil,
			},
			want: true,
		},
	}

	for name, _ := range testCases {
		testCase := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			processedMeasuresRepository := new(mocks.ProcessedMeasureRepository)

			processedMeasuresRepository.On("ListHistoryLoadCurve", ctx, process_measures.QueryHistoryLoadCurve{
				CUPS:          testCase.input.b.CUPS,
				StartDate:     testCase.input.b.InitDate.AddDate(0, 0, -1),
				EndDate:       testCase.input.b.EndDate.AddDate(0, 0, 1),
				WithCriterias: true,
				Periods:       testCase.input.b.Periods,
				Type:          measures.QuarterHour,
				Magnitude:     measures.AI,
			}).Return(testCase.output.listHistoryLoadCurveResult, testCase.output.listHistoryLoadCurveResultErr)

			evalResult := billing_measures.NewIsCCCHCompleteGD(testCase.input.b, measures.AI, processedMeasuresRepository)

			result := evalResult.Eval(ctx)
			assert.Equal(t, testCase.want, result)
		})
	}
}
