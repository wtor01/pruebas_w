package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Service_Scheduler_Get(t *testing.T) {
	type input struct {
		ctx context.Context
		dto GetSchedulerDTO
	}
	type want struct {
		err    error
		result process_measures.Scheduler
	}
	type results struct {
		getSchedulerResponse    process_measures.Scheduler
		getSchedulerResponseErr error
	}
	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"should return err if empty id": {
			input: input{
				ctx: context.Background(),
				dto: GetSchedulerDTO{},
			},
			want: want{
				err: process_measures.ErrSchedulerIdFormat,
			},
		},
		"should return err if fail GetScheduler": {
			input: input{
				ctx: context.Background(),
				dto: GetSchedulerDTO{
					ID: "54f6f0be-6059-4981-804d-7785f5aaba5e",
				},
			},
			want: want{
				err: errors.New("err"),
			},
			results: results{
				getSchedulerResponse:    process_measures.Scheduler{},
				getSchedulerResponseErr: errors.New("err"),
			},
		},
		"should return ok": {
			input: input{
				ctx: context.Background(),
				dto: GetSchedulerDTO{
					ID: "54f6f0be-6059-4981-804d-7785f5aaba5e",
				},
			},
			want: want{
				err: nil,
				result: process_measures.Scheduler{
					SchedulerId: "job-id",
				},
			},
			results: results{
				getSchedulerResponse: process_measures.Scheduler{
					SchedulerId: "job-id",
				},
				getSchedulerResponseErr: nil,
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			repo := new(mocks.SchedulerRepository)
			repo.Mock.On("GetScheduler", testCase.input.ctx, testCase.input.dto.ID).Return(testCase.results.getSchedulerResponse, testCase.results.getSchedulerResponseErr)

			srv := NewGetSchedulerByIdService(repo)
			result, err := srv.Handler(testCase.input.ctx, testCase.input.dto)
			assert.Equal(t, testCase.want.result, result, testCase)
			assert.Equal(t, testCase.want.err, err, testCase)

		})
	}

}
