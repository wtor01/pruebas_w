package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	mocks2 "bitbucket.org/sercide/data-ingestion/pkg/scheduler/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Services_Scheduler_Update(t *testing.T) {

	type input struct {
		ctx   context.Context
		dto   UpdateSchedulerDTO
		topic string
	}

	type want struct {
		response process_measures.Scheduler
		err      error
	}

	type results struct {
		getSchedulerResponse       process_measures.Scheduler
		getSchedulerResponseErr    error
		searchSchedulerResponse    []process_measures.Scheduler
		searchSchedulerResponseErr error
		updateSchedulerErr         error
		saveSchedulerErr           error
	}

	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"should return err if GetScheduler fail": {
			input: input{
				ctx: context.Background(),
				dto: UpdateSchedulerDTO{},
			},
			want: want{
				response: process_measures.Scheduler{},
				err:      errors.New("err"),
			},
			results: results{
				getSchedulerResponse:    process_measures.Scheduler{},
				getSchedulerResponseErr: errors.New("err"),
			},
		},
		"should return err if invalid params to Update": {
			input: input{
				ctx: context.Background(),
				dto: UpdateSchedulerDTO{
					Scheduler: "fail",
				},
			},
			want: want{
				response: process_measures.Scheduler{},
				err:      process_measures.ErrSchedulerInvalidFormat,
			},
			results: results{
				getSchedulerResponse:    process_measures.Scheduler{},
				getSchedulerResponseErr: nil,
			},
		},
		"should return err if SearchScheduler return values and is not the same id": {
			input: input{
				ctx: context.Background(),
				dto: UpdateSchedulerDTO{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ReadingType:   "curve",
					Scheduler:     "* * 1 * *",
					Description:   "description",
				},
			},
			want: want{
				response: process_measures.Scheduler{},
				err:      process_measures.ErrSchedulerExist,
			},
			results: results{
				getSchedulerResponse: process_measures.Scheduler{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					Name:          "name",
					Description:   "description",
					SchedulerId:   "job-id",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ReadingType:   "curve",
					Format:        "* * 1 * *",
				},
				getSchedulerResponseErr: nil,
				searchSchedulerResponse: []process_measures.Scheduler{
					{
						ID:            "4603aebd-afaa-4ba0-9252-4ab32fbcc9b5",
						DistributorId: "",
						Name:          "",
						Description:   "",
						SchedulerId:   "",
						ServiceType:   "",
						PointType:     "",
						MeterType:     nil,
						ReadingType:   "",
						Format:        "",
					},
				},
				searchSchedulerResponseErr: nil,
			},
		},
		"should continue if SearchScheduler return values but is the same id": {
			input: input{
				ctx: context.Background(),
				dto: UpdateSchedulerDTO{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ReadingType:   "curve",
					Scheduler:     "* * 1 * *",
					Description:   "description",
				},
			},
			want: want{
				response: process_measures.Scheduler{},
				err:      errors.New("err"),
			},
			results: results{
				getSchedulerResponse: process_measures.Scheduler{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					Name:          "name",
					Description:   "description",
					SchedulerId:   "job-id",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ReadingType:   "curve",
					Format:        "* * 1 * *",
				},
				getSchedulerResponseErr: nil,
				searchSchedulerResponse: []process_measures.Scheduler{
					{
						ID:            "cbfc1a75-542b-4667-882b-f19863163866",
						DistributorId: "",
						Name:          "",
						Description:   "",
						SchedulerId:   "",
						ServiceType:   "",
						PointType:     "",
						MeterType:     nil,
						ReadingType:   "",
						Format:        "",
					},
				},
				searchSchedulerResponseErr: nil,
				updateSchedulerErr:         errors.New("err"),
			},
		},
		"should return error if  UpdateJob fail": {
			input: input{
				ctx: context.Background(),
				dto: UpdateSchedulerDTO{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ReadingType:   "curve",
					Scheduler:     "* * 1 * *",
					Description:   "description",
				},
			},
			want: want{
				response: process_measures.Scheduler{},
				err:      errors.New("err"),
			},
			results: results{
				getSchedulerResponse: process_measures.Scheduler{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					Name:          "name",
					Description:   "description",
					SchedulerId:   "job-id",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ReadingType:   "curve",
					Format:        "* * 1 * *",
				},
				getSchedulerResponseErr:    nil,
				searchSchedulerResponse:    []process_measures.Scheduler{},
				searchSchedulerResponseErr: nil,
				updateSchedulerErr:         errors.New("err"),
			},
		},
		"should return error if  SaveScheduler fail and call to UpdateJob with old values": {
			input: input{
				ctx: context.Background(),
				dto: UpdateSchedulerDTO{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ReadingType:   "curve",
					Scheduler:     "* * 4 * *",
					Description:   "description",
				},
			},
			want: want{
				response: process_measures.Scheduler{},
				err:      errors.New("err"),
			},
			results: results{
				getSchedulerResponse: process_measures.Scheduler{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					Name:          "name",
					Description:   "description",
					SchedulerId:   "job-id",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ReadingType:   "curve",
					Format:        "* * 1 * *",
				},
				getSchedulerResponseErr:    nil,
				searchSchedulerResponse:    []process_measures.Scheduler{},
				searchSchedulerResponseErr: nil,
				updateSchedulerErr:         nil,
				saveSchedulerErr:           errors.New("err"),
			},
		},
		"should return ok": {
			input: input{
				ctx: context.Background(),
				dto: UpdateSchedulerDTO{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ReadingType:   "curve",
					Scheduler:     "* * 4 * *",
					Description:   "description",
				},
			},
			want: want{
				response: process_measures.Scheduler{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					Name:          "name",
					Description:   "description",
					SchedulerId:   "job-id",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ReadingType:   "curve",
					Format:        "* * 4 * *",
				},
				err: nil,
			},
			results: results{
				getSchedulerResponse: process_measures.Scheduler{
					ID:            "cbfc1a75-542b-4667-882b-f19863163866",
					DistributorId: "cbfc1a75-542b-4667-882b-f19863163866",
					Name:          "name",
					Description:   "description",
					SchedulerId:   "job-id",
					ServiceType:   "G-D",
					PointType:     "1",
					MeterType:     []string{"TLG"},
					ReadingType:   "curve",
					Format:        "* * 1 * *",
				},
				getSchedulerResponseErr:    nil,
				searchSchedulerResponse:    []process_measures.Scheduler{},
				searchSchedulerResponseErr: nil,
				updateSchedulerErr:         nil,
				saveSchedulerErr:           nil,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			repo := new(mocks.SchedulerRepository)
			client := new(mocks2.Client)

			schedulerCloned := testCase.results.getSchedulerResponse.Clone()
			_ = schedulerCloned.Update(
				testCase.input.dto.Description,
				testCase.input.dto.DistributorId,
				testCase.input.dto.ServiceType,
				testCase.input.dto.PointType,
				testCase.input.dto.MeterType,
				testCase.input.dto.ReadingType,
				testCase.input.dto.Scheduler,
			)

			clientCreator := func(ctx context.Context) (scheduler.Client, error) {
				client.Mock.On("UpdateJob", testCase.input.ctx, &schedulerCloned, testCase.input.topic).Return(testCase.results.updateSchedulerErr).Once()
				client.Mock.On("UpdateJob", testCase.input.ctx, &testCase.results.getSchedulerResponse, testCase.input.topic).Return(testCase.results.updateSchedulerErr)

				client.Mock.On("Close").Return(nil)
				return client, nil
			}

			repo.Mock.On("GetScheduler", testCase.input.ctx, testCase.input.dto.ID).Return(testCase.results.getSchedulerResponse, testCase.results.getSchedulerResponseErr)

			repo.Mock.On("SearchScheduler", testCase.input.ctx, process_measures.SearchScheduler{
				DistributorId: testCase.input.dto.DistributorId,
				ServiceType:   testCase.input.dto.ServiceType,
				PointType:     testCase.input.dto.PointType,
				MeterType:     testCase.input.dto.MeterType,
				ReadingType:   testCase.input.dto.ReadingType,
			}).Return(testCase.results.searchSchedulerResponse, testCase.results.searchSchedulerResponseErr)

			repo.Mock.On("SaveScheduler", testCase.input.ctx, schedulerCloned).Return(testCase.results.saveSchedulerErr)

			srv := NewUpdateScheduler(repo, clientCreator, testCase.input.topic)

			res, err := srv.Handler(testCase.input.ctx, testCase.input.dto)

			if testCase.results.saveSchedulerErr != nil {
				client.Mock.AssertNumberOfCalls(t, "UpdateJob", 2)
			}

			assert.Equal(t, testCase.want.response, res, testCase)
			assert.Equal(t, testCase.want.err, err, testCase)
		})
	}
}
