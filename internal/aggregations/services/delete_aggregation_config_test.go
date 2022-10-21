package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	mocks2 "bitbucket.org/sercide/data-ingestion/pkg/scheduler/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_Service_Aggregation_DeleteConfig_Handler(t *testing.T) {
	type input struct {
		id string
	}

	type output struct {
		GetConfig       aggregations.Config
		GetConfigErr    error
		DeleteJobErr    error
		DeleteConfigErr error
	}

	type want struct {
		err error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be delete config": {
			input: input{
				id: "1234",
			},
			output: output{
				GetConfig: aggregations.Config{
					Id:          "1234",
					Name:        "Test 2",
					Scheduler:   "* 1 * 2 *",
					SchedulerId: "SchedulerId",
					StartDate:   time.Date(2022, 7, 31, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2022, 8, 31, 0, 0, 0, 0, time.UTC),
					Features: []aggregations.Features{
						{
							ID:    "123",
							Name:  "DistributorCode",
							Field: "distributor_code",
						},
					},
				},
				GetConfigErr: nil,
			},
			want: want{
				err: nil,
			},
		},
		"Should be fail delete job": {
			input: input{
				id: "1234",
			},
			output: output{
				GetConfig: aggregations.Config{
					Id:          "1234",
					Name:        "Test 2",
					Scheduler:   "* 1 * 2 *",
					SchedulerId: "SchedulerId",
					StartDate:   time.Date(2022, 7, 31, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2022, 8, 31, 0, 0, 0, 0, time.UTC),
					Features: []aggregations.Features{
						{
							ID:    "123",
							Name:  "DistributorCode",
							Field: "distributor_code",
						},
					},
				},
				DeleteJobErr: errors.New("err delete job"),
			},
			want: want{
				err: errors.New("err delete job"),
			},
		},
		"Should be fail delete config": {
			input: input{
				id: "1234",
			},
			output: output{
				GetConfig: aggregations.Config{
					Id:          "1234",
					Name:        "Test 2",
					Scheduler:   "* 1 * 2 *",
					SchedulerId: "SchedulerId",
					StartDate:   time.Date(2022, 7, 31, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2022, 8, 31, 0, 0, 0, 0, time.UTC),
					Features: []aggregations.Features{
						{
							ID:    "123",
							Name:  "DistributorCode",
							Field: "distributor_code",
						},
					},
				},
				DeleteConfigErr: errors.New("err delete config"),
			},
			want: want{
				err: errors.New("err delete config"),
			},
		},
		"Should be fail get config": {
			input: input{
				id: "1234",
			},
			output: output{
				GetConfig:    aggregations.Config{},
				GetConfigErr: errors.New("err get config"),
			},
			want: want{
				err: errors.New("err get config"),
			},
		},
	}

	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			aggregationRepository := new(mocks.AggregationConfigRepository)
			aggregationRepository.On("GetAggregationConfigById", mock.Anything, test.input.id).Return(test.output.GetConfig, test.output.GetConfigErr)
			aggregationRepository.On("DeleteAggregationConfig", mock.Anything, test.output.GetConfig.Id).Return(test.output.DeleteConfigErr)

			schedulerClientMock := new(mocks2.Client)

			schedulerClient := func(ctx context.Context) (scheduler.Client, error) {
				schedulerClientMock.Mock.On("Close").Return(nil)
				schedulerClientMock.Mock.On("DeleteJob", mock.Anything, test.output.GetConfig.SchedulerId).Return(test.output.DeleteJobErr)
				schedulerClientMock.Mock.On("CreateJob", mock.Anything, test.output.GetConfig, "topic").Return("1234", nil)
				return schedulerClientMock, nil
			}

			s := NewDeleteAggregationConfigService(aggregationRepository, schedulerClient, "topic")
			err := s.Handler(context.Background(), test.input.id)

			assert.Equal(t, test.want.err, err)
		})
	}
}
