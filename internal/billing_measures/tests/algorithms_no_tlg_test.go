package tests

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/tests/fixtures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Unit_Domain_BillingMeasures_CchAverage(t *testing.T) {
	type input struct {
		b            *billing_measures.BillingMeasure
		ContextNoTlg *billing_measures.GraphContext
		magnitude    measures.Magnitude
	}

	type output struct {
		b billing_measures.BillingMeasure
	}

	testCases := map[string]struct {
		input  input
		output output
		want   error
	}{
		"All should be ok, no changes": {
			input: input{
				b:            &fixtures.BillingResult_1_Input_averages,
				ContextNoTlg: &billing_measures.GraphContext{},
				magnitude:    measures.AI,
			},
			output: output{b: fixtures.BillingResult_1_Output_averages},
			want:   nil,
		},
		"All should be ok, one change": {
			input: input{
				b:            &fixtures.BillingResult_2_Input_averages,
				ContextNoTlg: &billing_measures.GraphContext{},
				magnitude:    measures.AI,
			},
			output: output{b: fixtures.BillingResult_2_Output_averages},
			want:   nil,
		},
		"All should be ok, with 2 periods": {
			input: input{
				b:            &fixtures.BillingResult_3_Input_averages,
				ContextNoTlg: &billing_measures.GraphContext{},
				magnitude:    measures.AI,
			},
			output: output{b: fixtures.BillingResult_3_Output_averages},
			want:   nil,
		},
		"All should be ok, with 3 gap": {
			input: input{
				b:            &fixtures.BillingResult_4_Input_averages,
				ContextNoTlg: &billing_measures.GraphContext{},
				magnitude:    measures.AI,
			},
			output: output{b: fixtures.BillingResult_4_Output_averages},
			want:   nil,
		},
		"All should be ok, example task previous value historic": {
			input: input{
				b: &fixtures.BillingResult_Example_Input_averages,
				ContextNoTlg: &billing_measures.GraphContext{
					SimpleHistoric: billing_measures.ContextSimpleHistoric{
						PreviousLoadCurve: fixtures.Previous_ProcessedCurve_Example_Input_averages,
					},
				},
				magnitude: measures.AI,
			},
			output: output{b: fixtures.BillingResult_Example_Output_averages},

			want: nil,
		},
		"All should be ok, example task next value historic": {
			input: input{
				b: &fixtures.BillingResult_Example_Input_averages_example_next,
				ContextNoTlg: &billing_measures.GraphContext{
					SimpleHistoric: billing_measures.ContextSimpleHistoric{
						NextLoadCurve: fixtures.Next_ProcessedCurve_Example_Input_averages,
					},
				},
				magnitude: measures.AI,
			},
			output: output{b: fixtures.BillingResult_Example_Output_averages_example_next},

			want: nil,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {

			ctx := context.Background()
			cchAverage := billing_measures.NewCchAverage(test.input.b, test.input.magnitude, test.input.ContextNoTlg)
			cchAverage.Execute(ctx)
			cchId := cchAverage.ID()
			assert.Equal(t, cchId, "CCH_AVERAGE")
			assert.Equal(t, test.output.b.BillingLoadCurve, cchAverage.B.BillingLoadCurve)
			assert.Equal(t, test.output.b, *cchAverage.B)
		})
	}

}

func Test_Unit_Domain_BillingMeasures_CchPowerUseFactor(t *testing.T) {
	type input struct {
		b         *billing_measures.BillingMeasure
		magnitude measures.Magnitude
	}

	type output struct {
		b billing_measures.BillingMeasure
	}

	testCases := map[string]struct {
		input  input
		output output
		want   error
	}{
		"All should be ok, no changes": {
			input: input{
				magnitude: measures.AI,
				b:         &fixtures.BillingResult_1_Input_Power_Use_Factor,
			},
			output: output{b: fixtures.BillingResult_1_Input_Power_Use_Factor},
			want:   nil,
		},
		"All should be ok, one value changes": {
			input: input{
				magnitude: measures.AI,
				b:         &fixtures.BillingResult_2_Input_Power_Use_Factor,
			},
			output: output{b: fixtures.BillingResult_2_Output_Power_Use_Factor},
			want:   nil,
		},
		"All should be ok, example task": {
			input: input{
				magnitude: measures.AI,
				b:         &fixtures.BillingResult_Example_Input_Power_Use_Factor,
			},
			output: output{b: fixtures.BillingResult_Example_Output_Power_Use_Factor},
			want:   nil,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {

			ctx := context.Background()
			cchPowerUseFactor := billing_measures.NewPowerUseFactorCCH(test.input.b, test.input.magnitude)
			cchPowerUseFactor.Execute(ctx)
			cchId := cchPowerUseFactor.ID()
			assert.Equal(t, cchId, "CCH_POWER_USE_FACTOR")
			assert.Equal(t, test.output.b.BillingLoadCurve, cchPowerUseFactor.B.BillingLoadCurve)
			assert.Equal(t, test.output.b, *cchPowerUseFactor.B)
		})
	}

}

func Test_Unit_Domain_BillingMeasures_BalancePowerUseFactor(t *testing.T) {
	type input struct {
		b         *billing_measures.BillingMeasure
		magnitude measures.Magnitude
		ctx       context.Context
	}

	type output struct {
		b billing_measures.BillingMeasure
	}

	testCases := map[string]struct {
		input  input
		output output
		want   error
	}{
		"All should be ok, no empty": {
			input: input{
				b:         &fixtures.BillingResult_1_Input_Power_Use_Factor_Balance,
				magnitude: measures.AI,
			},
			output: output{b: fixtures.BillingResult_1_Output_Power_Use_Factor_Balance},
			want:   nil,
		},
		"All should be ok, 2 empty": {
			input: input{
				b:         &fixtures.BillingResult_2_Input_Power_Use_Factor_Balance,
				magnitude: measures.AI,
			},
			output: output{b: fixtures.BillingResult_2_Output_Power_Use_Factor_Balance},
			want:   nil,
		},
		"All should be 0 of values of magnitude that its not AI": {
			input: input{
				b:         &fixtures.BillingResult_3_Input_Power_Use_Factor_Balance,
				magnitude: measures.AE,
			},
			output: output{b: fixtures.BillingResult_3_Output_Power_Use_Factor_Balance},
			want:   nil,
		},
	}

	for name, _ := range testCases {
		t.Run(name, func(t *testing.T) {
			test := testCases[name]
			t.Parallel()
			ctx := context.Background()
			cchPowerUseFactor := billing_measures.NewPowerUseFactorBalance(test.input.b, test.input.magnitude)
			cchPowerUseFactor.Execute(ctx)
			cchId := cchPowerUseFactor.ID()
			assert.Equal(t, cchId, "BALANCE_POWER_USE_FACTOR")
			for _, period := range cchPowerUseFactor.B.Periods {
				assert.Equal(t, test.output.b.GetBalancePeriod(period), cchPowerUseFactor.B.GetBalancePeriod(period), period)
			}
			assert.Equal(t, test.output.b.BillingBalance.P0, cchPowerUseFactor.B.BillingBalance.P0, "P0")
			assert.Equal(t, test.output.b, *cchPowerUseFactor.B)
		})
	}

}

func Test_Unit_Domain_BillingMeasures_ClosingWithBalance(t *testing.T) {
	type input struct {
		b         *billing_measures.BillingMeasure
		magnitude measures.Magnitude
	}

	type output struct {
		b billing_measures.BillingMeasure
	}

	testCases := map[string]struct {
		input  input
		output output
		want   error
	}{
		"Should be ok, two invalid": {
			input: input{
				b:         &fixtures.BillingResult_1_Input_Closing_With_Balance,
				magnitude: measures.AI,
			},
			output: output{b: fixtures.BillingResult_1_Output_Closing_With_Balance},
			want:   nil,
		},
		"Should be ok, one invalid": {
			input: input{
				b:         &fixtures.BillingResult_2_Input_Closing_With_Balance,
				magnitude: measures.AI,
			},
			output: output{b: fixtures.BillingResult_2_Output_Closing_With_Balance},
			want:   nil,
		},
		"All should be 0 of values of magnitude that its not AE": {
			input: input{
				b:         &fixtures.BillingResult_3_Input_Closing_With_Balance,
				magnitude: measures.AE,
			},
			output: output{b: fixtures.BillingResult_3_Output_Closing_With_Balance},
			want:   nil,
		},
	}

	for name, _ := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			cchPowerUseFactor := billing_measures.NewClosingWithBalance(test.input.b, test.input.magnitude)
			cchPowerUseFactor.Execute(ctx)
			cchId := cchPowerUseFactor.ID()
			assert.Equal(t, cchId, "BALANCE_CLOSING")
			for _, period := range cchPowerUseFactor.B.Periods {
				assert.Equal(t, test.output.b.GetBalancePeriod(period), cchPowerUseFactor.B.GetBalancePeriod(period), period)
			}
			assert.Equal(t, test.output.b.BillingBalance.P0, cchPowerUseFactor.B.BillingBalance.P0, "P0")
			assert.Equal(t, test.output.b, *cchPowerUseFactor.B)
		})
	}

}

func Test_Unit_Domain_BillingMeasure_ReeBalanceOutline(t *testing.T) {

	var valueProfileTotal1 = 0.0036
	var valueProfileTotal2 = 0.0045
	var valueProfileTotal3 = 0.0087
	var valueProfileTotal4 = 0.0034
	var valueProfileTotal5 = 0.0084
	var valueProfileTotal6 = 0.0036
	var valueProfileTotal7 = 0.0013
	var valueProfileTotal8 = 0.0024

	type input struct {
		b         *billing_measures.BillingMeasure
		magnitude measures.Magnitude
		ctx       context.Context
	}
	type output struct {
		cp                     []billing_measures.ConsumProfile
		loadCurve              []billing_measures.BillingLoadCurve
		algorithmErrorResponse error
	}

	tests := map[string]struct {
		input  input
		output output
	}{
		"Error because consum records not found": {
			input: input{
				b:         &billing_measures.BillingMeasure{},
				magnitude: measures.AI,
				ctx:       context.Background(),
			},
			output: output{
				cp:                     []billing_measures.ConsumProfile{},
				algorithmErrorResponse: billing_measures.ErrorNotFoundConsumProfiles,
			},
		},
		"Should be okay": {
			input: input{
				b: &billing_measures.BillingMeasure{
					Coefficient: billing_measures.CoefficientA,
					BillingBalance: billing_measures.BillingBalance{
						P6: &billing_measures.BillingBalancePeriod{
							AE:              200000,
							EstimatedCodeAE: 1,
						},
					},
					BillingLoadCurve: []billing_measures.BillingLoadCurve{
						{
							EndDate: time.Date(2022, 05, 31, 1, 0, 0, 0, time.UTC),
							Period:  measures.P6,
						},
						{
							EndDate: time.Date(2022, 05, 31, 2, 0, 0, 0, time.UTC),
							Period:  measures.P6,
						},
						{
							EndDate: time.Date(2022, 05, 31, 3, 0, 0, 0, time.UTC),
							Period:  measures.P6,
						},
						{
							EndDate: time.Date(2022, 05, 31, 4, 0, 0, 0, time.UTC),
							Period:  measures.P6,
						},
						{
							EndDate: time.Date(2022, 05, 31, 5, 0, 0, 0, time.UTC),
							Period:  measures.P6,
						},
						{
							EndDate: time.Date(2022, 05, 31, 6, 0, 0, 0, time.UTC),
							Period:  measures.P6,
						},
						{
							EndDate: time.Date(2022, 05, 31, 7, 0, 0, 0, time.UTC),
							Period:  measures.P6,
						},
						{
							EndDate: time.Date(2022, 05, 31, 8, 0, 0, 0, time.UTC),
							Period:  measures.P6,
						},
					},
					Periods: []measures.PeriodKey{measures.P6},
				},
				magnitude: measures.AE,
				ctx:       context.Background(),
			},
			output: output{
				cp: []billing_measures.ConsumProfile{
					{Date: time.Date(2022, 05, 31, 1, 0, 0, 0, time.UTC),
						CoefA: &valueProfileTotal1,
					},
					{Date: time.Date(2022, 05, 31, 2, 0, 0, 0, time.UTC),
						CoefA: &valueProfileTotal2,
					},
					{Date: time.Date(2022, 05, 31, 3, 0, 0, 0, time.UTC),
						CoefA: &valueProfileTotal3,
					},
					{Date: time.Date(2022, 05, 31, 4, 0, 0, 0, time.UTC),
						CoefA: &valueProfileTotal4,
					},
					{Date: time.Date(2022, 05, 31, 5, 0, 0, 0, time.UTC),
						CoefA: &valueProfileTotal5,
					},
					{Date: time.Date(2022, 05, 31, 6, 0, 0, 0, time.UTC),
						CoefA: &valueProfileTotal6,
					},
					{Date: time.Date(2022, 05, 31, 7, 0, 0, 0, time.UTC),
						CoefA: &valueProfileTotal7,
					},
					{Date: time.Date(2022, 05, 31, 8, 0, 0, 0, time.UTC),
						CoefA: &valueProfileTotal8,
					},
				},
				loadCurve: []billing_measures.BillingLoadCurve{
					{
						EndDate:                  time.Date(2022, 05, 31, 1, 0, 0, 0, time.UTC),
						AE:                       20056,
						EstimatedCodeAE:          2,
						EstimatedMethodAE:        billing_measures.Outlined,
						MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
						EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
						Period:                   measures.P6,
					},
					{
						EndDate:                  time.Date(2022, 05, 31, 2, 0, 0, 0, time.UTC),
						AE:                       25070,
						EstimatedCodeAE:          2,
						EstimatedMethodAE:        billing_measures.Outlined,
						EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
						MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
						Period:                   measures.P6,
					}, {
						EndDate:                  time.Date(2022, 05, 31, 3, 0, 0, 0, time.UTC),
						AE:                       48468,
						EstimatedCodeAE:          2,
						EstimatedMethodAE:        billing_measures.Outlined,
						EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
						MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
						Period:                   measures.P6,
					},
					{
						EndDate:                  time.Date(2022, 05, 31, 4, 0, 0, 0, time.UTC),
						AE:                       18942,
						EstimatedCodeAE:          2,
						EstimatedMethodAE:        billing_measures.Outlined,
						EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
						MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
						Period:                   measures.P6,
					},
					{
						EndDate:                  time.Date(2022, 05, 31, 5, 0, 0, 0, time.UTC),
						AE:                       46797,
						EstimatedCodeAE:          2,
						EstimatedMethodAE:        billing_measures.Outlined,
						EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
						MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
						Period:                   measures.P6,
					},
					{
						EndDate:                  time.Date(2022, 05, 31, 6, 0, 0, 0, time.UTC),
						AE:                       20056,
						EstimatedCodeAE:          2,
						EstimatedMethodAE:        billing_measures.Outlined,
						EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
						MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
						Period:                   measures.P6,
					}, {
						EndDate:                  time.Date(2022, 05, 31, 7, 0, 0, 0, time.UTC),
						AE:                       7242,
						EstimatedCodeAE:          2,
						EstimatedMethodAE:        billing_measures.Outlined,
						EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
						MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
						Period:                   measures.P6,
					},
					{
						EndDate:                  time.Date(2022, 05, 31, 8, 0, 0, 0, time.UTC),
						AE:                       13370,
						EstimatedCodeAE:          2,
						EstimatedMethodAE:        billing_measures.Outlined,
						EstimatedGeneralMethodAE: billing_measures.GeneralOutlined,
						MeasureTypeAE:            billing_measures.ProvisionalBalanceMeasure,
						Period:                   measures.P6,
					},
				},
				algorithmErrorResponse: nil,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			repoProfiles := new(mocks.ConsumProfileRepository)

			repoProfiles.On("Search", testCase.input.ctx, billing_measures.QueryConsumProfile{StartDate: testCase.input.b.InitDate, EndDate: testCase.input.b.EndDate}).Return(testCase.output.cp, nil)

			reeBalanceOutline := billing_measures.NewReeBalanceOutline(testCase.input.b, repoProfiles, testCase.input.magnitude)

			response := reeBalanceOutline.Execute(testCase.input.ctx)

			assert.Equal(t, "CCH_REE_BALANCE_OUTLINE", reeBalanceOutline.ID())
			assert.Equal(t, testCase.output.algorithmErrorResponse, response)
			assert.Equal(t, reeBalanceOutline.B.BillingLoadCurve, testCase.output.loadCurve)

		})
	}
}

func Test_Unit_Domain_BillingMeasure_CCHByWindows(t *testing.T) {

	type input struct {
		magnitude measures.Magnitude
		b         *billing_measures.BillingMeasure
		context   *billing_measures.GraphContext
	}

	type output struct {
		loadCurves []billing_measures.BillingLoadCurve
	}

	type want struct {
		err error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be calculated correctly": {
			input: input{
				magnitude: measures.AI,
				b:         &fixtures.BILLING_MEASURE,
				context: &billing_measures.GraphContext{
					IterativeHistory: fixtures.HISTORY_LOAD_CURVE,
				},
			},
			output: output{
				loadCurves: fixtures.BILLING_MEASURE_RESULT.BillingLoadCurve,
			},
		},
	}

	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			algorithm := billing_measures.NewCCHWindows(test.input.b, test.input.magnitude, test.input.context)

			err := algorithm.Execute(ctx)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.output.loadCurves, algorithm.B.BillingLoadCurve)
		})
	}
}

func Test_Unit_Domain_BillingMeasure_CCHWindowsBalanceModulated(t *testing.T) {

	type input struct {
		magnitude measures.Magnitude
		b         *billing_measures.BillingMeasure
		context   *billing_measures.GraphContext
	}

	type output struct {
		loadCurves []billing_measures.BillingLoadCurve
	}

	type want struct {
		err error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be calculated correctly": {
			input: input{
				magnitude: measures.AI,
				b:         &fixtures.BILLING_MEASURE,
				context: &billing_measures.GraphContext{
					IterativeHistory: fixtures.HISTORY_LOAD_CURVE,
				},
			},
			output: output{
				loadCurves: fixtures.BILLING_MEASURE_RESULT_BALANCE_MODULATED.BillingLoadCurve,
			},
		},
	}

	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			algorithm := billing_measures.NewCCHWindowsBalanceModulated(test.input.b, test.input.magnitude, test.input.context)
			err := algorithm.Execute(ctx)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.output.loadCurves, algorithm.B.BillingLoadCurve)
		})
	}
}

func Test_Unit_Domain_BillingMeasure_CCHWindowsCloseModulated(t *testing.T) {

	type input struct {
		magnitude measures.Magnitude
		b         *billing_measures.BillingMeasure
		context   *billing_measures.GraphContext
	}

	type output struct {
		loadCurves []billing_measures.BillingLoadCurve
	}

	type want struct {
		err error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be calculated correctly": {
			input: input{
				magnitude: measures.AI,
				b:         &fixtures.BILLING_MEASURE,
				context: &billing_measures.GraphContext{
					IterativeHistory: fixtures.HISTORY_LOAD_CURVE,
				},
			},
			output: output{
				loadCurves: fixtures.BILLING_MEASURE_RESULT_CLOSE_MODULATED.BillingLoadCurve,
			},
		},
	}

	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			algorithm := billing_measures.NewCCHWindowsCloseModulated(test.input.b, test.input.magnitude, test.input.context)

			err := algorithm.Execute(ctx)

			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.output.loadCurves, algorithm.B.BillingLoadCurve)
		})
	}
}

func Test_Unit_Domain_BillingMeasure_CCCHComplete(t *testing.T) {

	type input struct {
		magnitude measures.Magnitude
		b         *billing_measures.BillingMeasure
		context   *billing_measures.GraphContext
	}

	type output struct {
		listHistoryLoadCurveResult    []process_measures.ProcessedLoadCurve
		b                             *billing_measures.BillingMeasure
		listHistoryLoadCurveResultErr error
	}

	type want struct {
		err error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be calculated correctly": {
			input: input{
				magnitude: measures.AI,
				b:         &fixtures.BillingResult_example_Input_ccch_complete_1,
				context:   &billing_measures.GraphContext{},
			},
			output: output{
				listHistoryLoadCurveResult:    fixtures.ProcessedLoadCurve_ccch_example_input_1,
				listHistoryLoadCurveResultErr: nil,
				b:                             &fixtures.BillingResult_example_Output_ccch_complete_1,
			},
		},
	}

	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			processedMeasuresRepository := new(mocks.ProcessedMeasureRepository)

			processedMeasuresRepository.On("ListHistoryLoadCurve", ctx, process_measures.QueryHistoryLoadCurve{
				CUPS:          test.input.b.CUPS,
				StartDate:     test.input.b.InitDate.AddDate(0, 0, -1),
				EndDate:       test.input.b.EndDate.AddDate(0, 0, 1),
				WithCriterias: true,
				Periods:       test.input.b.Periods,
				Type:          measures.QuarterHour,
				Magnitude:     test.input.magnitude,
			}).Return(test.output.listHistoryLoadCurveResult, test.output.listHistoryLoadCurveResultErr)

			algorithm := billing_measures.NewCCCHCompleteGD(test.input.b, test.input.magnitude, processedMeasuresRepository)
			err := algorithm.Execute(ctx)

			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.output.b.BillingLoadCurve, algorithm.B.BillingLoadCurve)
		})
	}
}
