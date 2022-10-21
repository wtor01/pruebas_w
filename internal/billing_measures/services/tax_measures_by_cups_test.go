package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Unit_Service_TaxMeasuresByCups(t *testing.T) {
	type input struct {
		ctx context.Context
		dto billing_measures.QueryBillingMeasuresTax
	}
	type want struct {
		err    error
		count  int
		result billing_measures.BillingMeasuresTaxResult
	}
	type results struct {
		taxMeasuresResponse      billing_measures.BillingMeasuresTaxResult
		taxMeasuresResponseCount int
		taxMeasuresResponseErr   error
	}
	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"Should return err if fail TaxMeasures": {
			input: input{
				ctx: context.Background(),
				dto: billing_measures.QueryBillingMeasuresTax{},
			},
			want: want{
				err:    errors.New("TaxMeasures Have Failed"),
				count:  0,
				result: billing_measures.BillingMeasuresTaxResult{},
			},
			results: results{
				taxMeasuresResponse:      billing_measures.BillingMeasuresTaxResult{},
				taxMeasuresResponseCount: 0,
				taxMeasuresResponseErr:   errors.New("TaxMeasures Have Failed"),
			},
		},
		"Should return ok": {
			input: input{
				ctx: context.Background(),
				dto: billing_measures.QueryBillingMeasuresTax{},
			},
			want: want{
				err:   nil,
				count: 1,
				result: billing_measures.BillingMeasuresTaxResult{

					Data: []billing_measures.BillingMeasuresTax{{
						Cups:             "21321312",
						DistributorId:    "1231312",
						StartDate:        time.Time{},
						EndDate:          time.Time{},
						ExecutionSummary: billing_measures.ExecutionSummary{}},
					},
					Count: 1,
				},
			},
			results: results{
				taxMeasuresResponse: billing_measures.BillingMeasuresTaxResult{

					Data: []billing_measures.BillingMeasuresTax{{
						Cups:             "21321312",
						DistributorId:    "1231312",
						StartDate:        time.Time{},
						EndDate:          time.Time{},
						ExecutionSummary: billing_measures.ExecutionSummary{}},
					},
					Count: 1,
				},

				taxMeasuresResponseCount: 1,
				taxMeasuresResponseErr:   nil,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			repo := new(mocks.BillingMeasuresDashboardRepository)

			repo.Mock.On("GetBillingMeasuresTax", testCase.input.ctx, billing_measures.QueryBillingMeasuresTax{
				Offset:        testCase.input.dto.Offset,
				Limit:         testCase.input.dto.Limit,
				DistributorId: testCase.input.dto.DistributorId,
				MeasureType:   testCase.input.dto.MeasureType,
				StartDate:     testCase.input.dto.StartDate,
				EndDate:       testCase.input.dto.EndDate,
			}).Return(testCase.results.taxMeasuresResponse, testCase.results.taxMeasuresResponseErr)

			srv := NewTaxMeasuresByCups(repo)

			result, err := srv.Handler(testCase.input.ctx, testCase.input.dto)
			assert.Equal(t, testCase.want.result, result, testCase)
			assert.Equal(t, testCase.want.count, result.Count, testCase)
			assert.Equal(t, testCase.want.err, err, testCase)
		})
	}
}
