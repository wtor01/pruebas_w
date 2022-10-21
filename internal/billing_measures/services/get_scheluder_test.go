package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Service_Billing_Scheduler_Get(t *testing.T) {
	type input struct {
		ctx context.Context
		dto GetSchedulerDTO
	}
	type want struct {
		err    error
		result billing_measures.Scheduler
	}
	type results struct {
		getSchedulerResponse    billing_measures.Scheduler
		getSchedulerResponseErr error
	}
	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"Should return err if empty Id": {
			input: input{
				ctx: context.Background(),
				dto: GetSchedulerDTO{},
			},
			want: want{
				err: billing_measures.ErrSchedulerIdFormat,
			},
		},
		"Should Return Err if Fails GetScheduler": {
			input: input{
				ctx: context.Background(),
				dto: GetSchedulerDTO{
					ID: "54f6f0be-6059-4981-804d-7785f5aaba5e",
				},
			},
			results: results{
				getSchedulerResponse:    billing_measures.Scheduler{},
				getSchedulerResponseErr: errors.New("Get Scheluder Failed"),
			},
			want: want{
				err:    errors.New("Get Scheluder Failed"),
				result: billing_measures.Scheduler{},
			},
		},
		"Should be OK": {
			input: input{
				ctx: context.Background(),
				dto: GetSchedulerDTO{
					ID: "54f6f0be-6059-4981-804d-7785f5aaba5e",
				},
			},
			results: results{
				getSchedulerResponse: billing_measures.Scheduler{
					SchedulerId: "I-want-this-jobID",
				},
				getSchedulerResponseErr: nil,
			},
			want: want{
				err: nil,
				result: billing_measures.Scheduler{
					SchedulerId: "I-want-this-jobID",
				},
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			repo := new(mocks.BillingSchedulerRepository)
			repo.Mock.On("GetScheduler", testCase.input.ctx, testCase.input.dto.ID).Return(testCase.results.getSchedulerResponse, testCase.results.getSchedulerResponseErr)
			srv := NewGetSchedulerByIdService(repo)
			result, err := srv.Handler(testCase.input.ctx, testCase.input.dto)
			assert.Equal(t, testCase.want.result, result, testCase)
			assert.Equal(t, testCase.want.err, err, testCase)
		})
	}
}
