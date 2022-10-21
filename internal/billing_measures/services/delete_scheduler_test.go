package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	mocks2 "bitbucket.org/sercide/data-ingestion/pkg/scheduler/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Service_Scheduler_Delete(t *testing.T) {

	type input struct {
		ctx context.Context
		dto DeleteSchedulerDTO
	}

	type want struct {
		err error
	}

	type results struct {
		getSchedulerResponse    billing_measures.Scheduler
		getSchedulerResponseErr error
		deleteJobErr            error
		deleteSchedulerErr      error
	}

	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"should return err if empty id": {
			input: input{
				ctx: context.Background(),
				dto: DeleteSchedulerDTO{},
			},
			want: want{
				err: billing_measures.ErrSchedulerIdFormat,
			},
		},
		"should return err if fail GetScheduler": {
			input: input{
				ctx: context.Background(),
				dto: DeleteSchedulerDTO{
					ID: "54f6f0be-6059-4981-804d-7785f5aaba5e",
				},
			},
			want: want{
				err: errors.New("err"),
			},
			results: results{
				getSchedulerResponse:    billing_measures.Scheduler{},
				getSchedulerResponseErr: errors.New("err"),
			},
		},
		"should return err if DeleteJob return err": {
			input: input{
				ctx: context.Background(),
				dto: DeleteSchedulerDTO{
					ID: "54f6f0be-6059-4981-804d-7785f5aaba5e",
				},
			},
			want: want{
				err: errors.New("err"),
			},
			results: results{
				getSchedulerResponse: billing_measures.Scheduler{
					SchedulerId: "job-id",
				},
				getSchedulerResponseErr: nil,
				deleteJobErr:            errors.New("err"),
			},
		},
		"should return err if DeleteScheduler return err": {
			input: input{
				ctx: context.Background(),
				dto: DeleteSchedulerDTO{
					ID: "54f6f0be-6059-4981-804d-7785f5aaba5e",
				},
			},
			want: want{
				err: errors.New("err"),
			},
			results: results{
				getSchedulerResponse: billing_measures.Scheduler{
					SchedulerId: "job-id",
				},
				getSchedulerResponseErr: nil,
				deleteJobErr:            nil,
				deleteSchedulerErr:      errors.New("err"),
			},
		},
		"should return ok": {
			input: input{
				ctx: context.Background(),
				dto: DeleteSchedulerDTO{
					ID: "54f6f0be-6059-4981-804d-7785f5aaba5e",
				},
			},
			want: want{
				err: nil,
			},
			results: results{
				getSchedulerResponse: billing_measures.Scheduler{
					SchedulerId: "job-id",
				},
				getSchedulerResponseErr: nil,
				deleteJobErr:            nil,
				deleteSchedulerErr:      nil,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			repo := new(mocks.BillingSchedulerRepository)
			client := new(mocks2.Client)

			clientCreator := func(ctx context.Context) (scheduler.Client, error) {

				client.Mock.On("DeleteJob", testCase.input.ctx, testCase.results.getSchedulerResponse.SchedulerId).Return(testCase.results.deleteJobErr)
				client.Mock.On("Close").Return(nil)

				return client, nil
			}

			repo.Mock.On("GetScheduler", testCase.input.ctx, testCase.input.dto.ID).Return(testCase.results.getSchedulerResponse, testCase.results.getSchedulerResponseErr)

			repo.Mock.On("DeleteScheduler", testCase.input.ctx, testCase.results.getSchedulerResponse.ID).Return(testCase.results.deleteSchedulerErr)

			srv := NewDeleteScheduler(repo, clientCreator)

			err := srv.Handler(testCase.input.ctx, testCase.input.dto)

			assert.Equal(t, testCase.want.err, err, testCase)
		})
	}
}
