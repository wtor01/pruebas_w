package tests

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/tests/fixtures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Domain_Algorithms_CchPartialEstimation(t *testing.T) {
	type input struct {
		b      *billing_measures.BillingMeasure
		period []measures.PeriodKey
		ctx    context.Context
	}

	type output struct {
		b  billing_measures.BillingMeasure
		cp []billing_measures.ConsumProfile
	}

	testCases := map[string]struct {
		input  input
		output output
		want   error
	}{
		"Should return error for consum records not found": {
			input: input{
				ctx:    context.Background(),
				b:      &billing_measures.BillingMeasure{},
				period: []measures.PeriodKey{measures.P0},
			},
			output: output{b: billing_measures.BillingMeasure{},
				cp: []billing_measures.ConsumProfile{}},
			want: billing_measures.ErrorNotFoundConsumProfiles,
		},
		"All should be ok, no changes": {
			input: input{
				ctx:    context.Background(),
				b:      &fixtures.BillingResult_1_Input_partial,
				period: []measures.PeriodKey{measures.P0},
			},
			output: output{b: fixtures.BillingResult_1_Output_partial,
				cp: fixtures.Fixture_consume_profiles_partial},
			want: nil,
		},
		"Change for one value": {
			input: input{
				ctx:    context.Background(),
				b:      &fixtures.BillingResult_2_Input_partial,
				period: []measures.PeriodKey{measures.P0},
			},
			output: output{b: fixtures.BillingResult_2_Output_partial,
				cp: fixtures.Fixture_consume_profiles_partial},
			want: nil,
		},
		"Change for many values and periods": {
			input: input{
				ctx:    context.Background(),
				b:      &fixtures.BillingResult_3_Input_partial,
				period: []measures.PeriodKey{measures.P0, measures.P1},
			},
			output: output{b: fixtures.BillingResult_3_Output_partial,
				cp: fixtures.Fixture_consume_profiles_partial},
			want: nil,
		},
		"Example task": {
			input: input{
				ctx:    context.Background(),
				b:      &fixtures.BillingResult_1_Input_partial_example,
				period: []measures.PeriodKey{measures.P4},
			},
			output: output{b: fixtures.BillingResult_1_Output_partial_example,
				cp: fixtures.Fixture_consume_profiles_example_partial},
			want: nil,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			repoProfiles := new(mocks.ConsumProfileRepository)

			repoProfiles.On("Search", test.input.ctx, billing_measures.QueryConsumProfile{StartDate: test.input.b.InitDate, EndDate: test.input.b.EndDate}).Return(test.output.cp, nil)

			cchPartial := billing_measures.NewCchPartialEstimation(test.input.b, measures.P0, repoProfiles, measures.AE)
			var response error
			for _, period := range test.input.period {
				cchPartial.B = test.input.b
				cchPartial.Period = period
				response = cchPartial.Execute(ctx)
			}
			cchId := cchPartial.ID()
			assert.Equal(t, cchId, "CCH_PARTIAL_ESTIMATE")
			assert.Equal(t, test.want, response)
			assert.Equal(t, test.output.b.BillingLoadCurve, cchPartial.B.BillingLoadCurve)
			assert.Equal(t, test.output.b, *cchPartial.B)
		})
	}

}

func Test_Unit_Domain_Algorithms_CchTotalEstimation(t *testing.T) {
	type input struct {
		b      *billing_measures.BillingMeasure
		period []measures.PeriodKey
		ctx    context.Context
	}

	type output struct {
		b  billing_measures.BillingMeasure
		cp []billing_measures.ConsumProfile
	}

	testCases := map[string]struct {
		input  input
		output output
		want   error
	}{
		"Should return error if consum records are not found": {
			input: input{
				ctx:    context.Background(),
				b:      &billing_measures.BillingMeasure{},
				period: []measures.PeriodKey{measures.P0},
			},
			output: output{b: billing_measures.BillingMeasure{},
				cp: []billing_measures.ConsumProfile{}},
			want: billing_measures.ErrorNotFoundConsumProfiles,
		},
		"All should be ok": {
			input: input{
				ctx:    context.Background(),
				b:      &fixtures.BillingResult_1_Input_total,
				period: []measures.PeriodKey{measures.P0},
			},
			output: output{b: fixtures.BillingResult_1_Output_total,
				cp: fixtures.Fixture_consume_profiles_total},
			want: nil,
		},
		"Example task test": {
			input: input{
				ctx:    context.Background(),
				b:      &fixtures.BillingResult_1_Input_total_example,
				period: []measures.PeriodKey{measures.P0},
			},
			output: output{b: fixtures.BillingResult_1_Output_total_example,
				cp: fixtures.Fixture_consume_profiles_example_total},
			want: nil,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {

			repoProfiles := new(mocks.ConsumProfileRepository)

			repoProfiles.On("Search", test.input.ctx, billing_measures.QueryConsumProfile{StartDate: test.input.b.InitDate, EndDate: test.input.b.EndDate}).Return(test.output.cp, nil)

			cchPartial := billing_measures.NewCchTotalEstimation(test.input.b, measures.P0, repoProfiles, measures.AE)
			var response error
			for _, period := range test.input.period {
				cchPartial.B = test.input.b
				cchPartial.Period = period
				response = cchPartial.Execute(context.Background())
			}
			cchId := cchPartial.ID()
			assert.Equal(t, cchId, "CCH_TOTAL_ESTIMATE")
			assert.Equal(t, test.output.b.BillingLoadCurve, cchPartial.B.BillingLoadCurve)
			assert.Equal(t, test.want, response)
			assert.Equal(t, test.output.b, *cchPartial.B)
		})
	}

}
