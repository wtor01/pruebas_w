package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	mocks2 "bitbucket.org/sercide/data-ingestion/pkg/scheduler/mocks"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_Service_Aggregation_UpdateConfig_Handler(t *testing.T) {
	type input struct {
		configId string
		dto      UpdateConfigDto
	}

	type output struct {
		GetConfig           aggregations.Config
		GetConfigErr        error
		UpdateJobErr        error
		UpdateConfigErr     error
		UpdateConfig        aggregations.Config
		GetFeaturesByIdsErr error
		GetFeaturesByIds    []aggregations.Features
	}

	type want struct {
		err    error
		result aggregations.Config
	}

	invalidEndDate := time.Date(2020, 5, 30, 0, 0, 0, 0, time.UTC)

	loc, _ := time.LoadLocation("Europe/Madrid")
	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"should be err get config": {
			input: input{
				configId: "1234",
				dto: NewUpdateConfigDto(
					"",
					nil,
					time.Time{},
					nil,
					[]aggregations.ConfigFeatureDto{{}},
				),
			},
			output: output{
				GetConfig:    aggregations.Config{},
				GetConfigErr: errors.New("err get config"),
			},
			want: want{
				err:    errors.New("err get config"),
				result: aggregations.Config{},
			},
		},
		"should be validation scheduler format err": {
			input: input{
				configId: "1234",
				dto: NewUpdateConfigDto(
					"",
					nil,
					time.Time{},
					nil,
					[]aggregations.ConfigFeatureDto{{}},
				),
			},
			output: output{
				GetConfig: aggregations.Config{
					Id:          "1234",
					Name:        "Test",
					Scheduler:   "* * * * *",
					Description: "",
					StartDate:   time.Time{},
					EndDate:     time.Time{},
					Features:    []aggregations.Features{{}},
				},
			},
			want: want{
				err:    scheduler.ErrSchedulerInvalidFormat,
				result: aggregations.Config{},
			},
		},
		"should be validation date err": {
			input: input{
				configId: "1234",
				dto: NewUpdateConfigDto(
					"* * * * *",
					nil,
					time.Time{},
					&invalidEndDate,
					[]aggregations.ConfigFeatureDto{{}},
				),
			},
			output: output{
				GetConfig: aggregations.Config{
					Id:          "1234",
					Name:        "Test",
					Scheduler:   "* * * * *",
					Description: "",
					StartDate:   time.Time{},
					EndDate:     time.Time{},
					Features:    []aggregations.Features{{}},
				},
			},
			want: want{
				err:    scheduler.ErrInvalidScheduler,
				result: aggregations.Config{},
			},
		},
		"should be validation features err": {
			input: input{
				configId: "1234",
				dto: NewUpdateConfigDto(
					"* * * * *",
					nil,
					time.Time{},
					nil,
					[]aggregations.ConfigFeatureDto{},
				),
			},
			output: output{

				GetConfig: aggregations.Config{
					Id:          "1234",
					Name:        "Test",
					Scheduler:   "* * * * *",
					Description: "",
					StartDate:   time.Time{},
					EndDate:     time.Time{},
					Features:    []aggregations.Features{{}},
				},
			},
			want: want{
				err:    scheduler.ErrInvalidScheduler,
				result: aggregations.Config{},
			},
		},
		"should be update job err": {
			input: input{
				configId: "1234",
				dto: NewUpdateConfigDto(
					"* * * * *",
					nil,
					time.Time{},
					nil,
					[]aggregations.ConfigFeatureDto{{}},
				),
			},
			output: output{
				GetFeaturesByIds: []aggregations.Features{
					{
						ID:    "goodID",
						Name:  "goodName",
						Field: "goodField",
					},
				},
				GetFeaturesByIdsErr: nil,
				GetConfig: aggregations.Config{
					Id:          "1234",
					Name:        "Test",
					Scheduler:   "* * * * *",
					Description: "",
					StartDate:   time.Time{},
					EndDate:     time.Time{},
					Features:    []aggregations.Features{{}},
				},
				UpdateJobErr: errors.New("err update job"),
			},
			want: want{
				err:    errors.New("err update job"),
				result: aggregations.Config{},
			},
		},
		"should be update config err": {
			input: input{
				configId: "1234",
				dto: NewUpdateConfigDto(
					"* * * * *",
					nil,
					time.Time{},
					nil,
					[]aggregations.ConfigFeatureDto{{}},
				),
			},
			output: output{
				GetFeaturesByIds: []aggregations.Features{
					{
						ID:    "goodID",
						Name:  "goodName",
						Field: "goodField",
					},
				},
				GetFeaturesByIdsErr: nil,
				GetConfig: aggregations.Config{
					Id:          "1234",
					Name:        "Test",
					Scheduler:   "* * * * *",
					Description: "",
					StartDate:   time.Time{},
					EndDate:     time.Time{},
					Features:    []aggregations.Features{{}},
				},
				UpdateConfigErr: errors.New("err update config"),
			},
			want: want{
				err:    errors.New("err update config"),
				result: aggregations.Config{},
			},
		},
		"should be update config": {
			input: input{
				configId: "1234",
				dto: NewUpdateConfigDto(
					"* * * 1 2",
					nil,
					time.Date(2022, 5, 30, 0, 0, 0, 0, time.UTC),
					nil,
					[]aggregations.ConfigFeatureDto{{}, {}},
				),
			},
			output: output{
				GetConfig: aggregations.Config{
					Id:          "1234",
					Name:        "Test",
					Scheduler:   "* * * * *",
					SchedulerId: "1234",
					Description: "",
					StartDate:   time.Time{},
					EndDate:     time.Time{},
					Features:    []aggregations.Features{{}},
				},
				GetFeaturesByIds: []aggregations.Features{
					{
						ID:    "goodID",
						Name:  "goodName",
						Field: "goodField",
					},
				},
				GetFeaturesByIdsErr: nil,
				UpdateConfig: aggregations.Config{
					Id:          "1234",
					Name:        "Test",
					Scheduler:   "* * * 1 2",
					SchedulerId: "1234",
					StartDate:   time.Date(2022, 5, 30, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Time{},
					Description: "",
					Features: []aggregations.Features{{
						ID:    "goodID",
						Name:  "goodName",
						Field: "goodField",
					}},
				},
			},
			want: want{
				result: aggregations.Config{
					Id:          "1234",
					Name:        "Test",
					Scheduler:   "* * * 1 2",
					SchedulerId: "1234",
					StartDate:   time.Date(2022, 5, 30, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Time{},
					Description: "",
					Features: []aggregations.Features{
						{
							ID:    "goodID",
							Name:  "goodName",
							Field: "goodField",
						},
					},
				},
			},
		},
	}

	for name := range testCases {
		test := testCases[name]

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			featuresRepo := new(mocks.AggregationsFeaturesRepository)

			featuresRepo.On("GetFeaturesByIds", mock.Anything, utils.MapSlice(test.input.dto.Features, func(item aggregations.ConfigFeatureDto) string {
				return item.Id
			})).Return(test.output.GetFeaturesByIds, test.output.GetFeaturesByIdsErr)

			aggregationRepository := new(mocks.AggregationConfigRepository)
			schedulerClientMock := new(mocks2.Client)

			updatedConfig := test.output.GetConfig.Clone()
			_ = updatedConfig.Update(
				test.input.dto.Scheduler,
				test.input.dto.Description,
				test.input.dto.StartDate,
				test.input.dto.EndDate,
				test.output.GetFeaturesByIds,
			)

			aggregationRepository.On("GetAggregationConfigById", mock.Anything, test.input.configId).Return(test.output.GetConfig, test.output.GetConfigErr)
			aggregationRepository.On("SaveAggregationConfig", mock.Anything, updatedConfig).Return(test.output.UpdateConfig, test.output.UpdateConfigErr)

			schedulerClient := func(ctx context.Context) (scheduler.Client, error) {
				schedulerClientMock.Mock.On("UpdateJob", mock.Anything, updatedConfig, "topic").Return(test.output.UpdateJobErr).Once()
				schedulerClientMock.Mock.On("UpdateJob", mock.Anything, test.output.GetConfig, "topic").Return(test.output.UpdateJobErr)
				schedulerClientMock.Mock.On("Close").Return(nil)
				return schedulerClientMock, nil
			}

			s := NewUpdateAggregationConfigService(aggregationRepository, featuresRepo, schedulerClient, "topic", loc)

			result, err := s.Handler(context.Background(), test.input.configId, test.input.dto)

			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.result, result)

			if test.output.UpdateConfigErr != nil {
				schedulerClientMock.AssertNumberOfCalls(t, "UpdateJob", 2)
			}
		})
	}
}
