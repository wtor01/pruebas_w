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

func Test_Unit_Services_Scheduler_Create(t *testing.T) {

	type input struct {
		ctx   context.Context
		dto   CreateSchedulerDTO
		topic string
	}

	type want struct {
		response billing_measures.Scheduler
		err      error
	}

	type results struct {
		searchSchedulerResponse    []billing_measures.Scheduler
		searchSchedulerResponseErr error
		createJobResponse          string
		createJobErr               error
		saveSchedulerErr           error
	}

	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"should return err if invalid params to NewScheduler": {
			input: input{
				ctx: context.Background(),
				dto: CreateSchedulerDTO{},
			},
			want: want{
				response: billing_measures.Scheduler{},
				err:      billing_measures.ErrSchedulerInvalidFormat,
			},
		},
		"should return err if fail SearchScheduler": {
			input: input{
				ctx: context.Background(),
				dto: CreateSchedulerDTO{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ProcessType:   "D-C",
					Scheduler:     "* * 1 * *",
					Name:          "name",
				},
			},
			want: want{
				response: billing_measures.Scheduler{},
				err:      errors.New("err"),
			},
			results: results{
				searchSchedulerResponse:    nil,
				searchSchedulerResponseErr: errors.New("err"),
			},
		},
		"should return err if SearchScheduler return values": {
			input: input{
				ctx: context.Background(),
				dto: CreateSchedulerDTO{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ProcessType:   "D-C",
					Scheduler:     "* * 1 * *",
					Name:          "name",
				},
			},
			want: want{
				response: billing_measures.Scheduler{},
				err:      billing_measures.ErrSchedulerExist,
			},
			results: results{
				searchSchedulerResponse: []billing_measures.Scheduler{
					{
						ID:            "",
						DistributorId: "",
						Name:          "",
						SchedulerId:   "",
						ServiceType:   "",
						PointType:     "",
						MeterType:     nil,
						ProcessType:   "",
						Format:        "",
					},
				},
				searchSchedulerResponseErr: nil,
			},
		},
		"should return err if CreateJob return err": {
			input: input{
				ctx: context.Background(),
				dto: CreateSchedulerDTO{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ProcessType:   "D-C",
					Scheduler:     "* * 1 * *",
					Name:          "name",
				},
			},
			want: want{
				response: billing_measures.Scheduler{},
				err:      errors.New("err"),
			},
			results: results{
				searchSchedulerResponse:    []billing_measures.Scheduler{},
				searchSchedulerResponseErr: nil,
				createJobErr:               errors.New("err"),
				createJobResponse:          "",
			},
		},
		"should return err if SaveScheduler return err and call DeleteJob": {
			input: input{
				ctx: context.Background(),
				dto: CreateSchedulerDTO{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ProcessType:   "D-C",
					Scheduler:     "* * 1 * *",
					Name:          "name",
				},
			},
			want: want{
				response: billing_measures.Scheduler{},
				err:      errors.New("err"),
			},
			results: results{
				searchSchedulerResponse:    []billing_measures.Scheduler{},
				searchSchedulerResponseErr: nil,
				createJobErr:               nil,
				createJobResponse:          "job-id",
				saveSchedulerErr:           errors.New("err"),
			},
		},
		"should return ok": {
			input: input{
				ctx: context.Background(),
				dto: CreateSchedulerDTO{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ProcessType:   "D-C",
					Scheduler:     "* * 1 * *",
					Name:          "name",
				},
			},
			want: want{
				response: billing_measures.Scheduler{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					Name:          "name",
					SchedulerId:   "job-id",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ProcessType:   "D-C",
					Format:        "* * 1 * *",
				},
				err: nil,
			},
			results: results{
				searchSchedulerResponse:    []billing_measures.Scheduler{},
				searchSchedulerResponseErr: nil,
				createJobErr:               nil,
				createJobResponse:          "job-id",
				saveSchedulerErr:           nil,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			repo := new(mocks.BillingSchedulerRepository)
			client := new(mocks2.Client)
			sc, _ := billing_measures.NewScheduler(
				testCase.input.dto.ID,
				testCase.input.dto.Name,
				"",
				testCase.input.dto.DistributorId,
				testCase.input.dto.ServiceType,
				testCase.input.dto.PointType,
				testCase.input.dto.MeterType,
				testCase.input.dto.ProcessType,
				testCase.input.dto.Scheduler,
			)

			clientCreator := func(ctx context.Context) (scheduler.Client, error) {
				sc, _ := billing_measures.NewScheduler(
					testCase.input.dto.ID,
					testCase.input.dto.Name,
					"",
					testCase.input.dto.DistributorId,
					testCase.input.dto.ServiceType,
					testCase.input.dto.PointType,
					testCase.input.dto.MeterType,
					testCase.input.dto.ProcessType,
					testCase.input.dto.Scheduler,
				)
				client.Mock.On("CreateJob", testCase.input.ctx, &sc, testCase.input.topic).Return(testCase.results.createJobResponse, testCase.results.createJobErr)
				client.Mock.On("Close").Return(nil)
				client.Mock.On("DeleteJob", testCase.input.ctx, testCase.results.createJobResponse).Return(nil)
				return client, nil
			}

			repo.Mock.On("SearchScheduler", testCase.input.ctx, billing_measures.SearchScheduler{
				DistributorId: testCase.input.dto.DistributorId,
				ServiceType:   testCase.input.dto.ServiceType,
				PointType:     testCase.input.dto.PointType,
				MeterType:     testCase.input.dto.MeterType,
				ProcessType:   testCase.input.dto.ProcessType,
			}).Return(testCase.results.searchSchedulerResponse, testCase.results.searchSchedulerResponseErr)

			sc.SetSchedulerId(testCase.results.createJobResponse)

			repo.Mock.On("SaveScheduler", testCase.input.ctx, sc).Return(testCase.results.saveSchedulerErr)

			srv := NewCreateScheduler(repo, clientCreator, testCase.input.topic)

			res, err := srv.Handler(testCase.input.ctx, testCase.input.dto)

			assert.Equal(t, testCase.want.response, res, testCase)
			assert.Equal(t, testCase.want.err, err, testCase)
		})
	}
}
