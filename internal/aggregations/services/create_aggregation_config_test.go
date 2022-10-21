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

func Test_Unit_Service_Aggregation_CreateConfig_Handler(t *testing.T) {
	type input struct {
		dto              CreateConfigDto
		GetFeaturesByIds []string
	}

	type output struct {
		CreateJobId         string
		CreateJobErr        error
		CreateConfigErr     error
		CreateConfig        aggregations.Config
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
		"should be error by id": {
			input: input{
				dto: NewCreateConfigDto(
					"Test A",
					"* * * * *",
					nil,
					time.Time{},
					nil,
					[]aggregations.ConfigFeatureDto{{}},
				),
			},
			want: want{
				err:    errors.New("err get features by id"),
				result: aggregations.Config{},
			},
			output: output{
				GetFeaturesByIdsErr: errors.New("err get features by id"),
				GetFeaturesByIds:    []aggregations.Features{},
			},
		},
		"should be validation name err": {
			input: input{
				dto: NewCreateConfigDto(
					"Test A",
					"* * * * *",
					nil,
					time.Time{},
					nil,
					[]aggregations.ConfigFeatureDto{{}},
				),
			},
			want: want{
				err:    scheduler.ErrSchedulerInvalidName,
				result: aggregations.Config{},
			},
		},
		"should be validation scheduler format err": {
			input: input{
				dto: NewCreateConfigDto(
					"Test-A",
					"",
					nil,
					time.Time{},
					nil,
					[]aggregations.ConfigFeatureDto{{}},
				),
			},
			want: want{
				err:    scheduler.ErrSchedulerInvalidFormat,
				result: aggregations.Config{},
			},
		},
		"should be validation date err": {
			input: input{
				dto: NewCreateConfigDto(
					"Test-A",
					"* * * * *",
					nil,
					time.Time{},
					&invalidEndDate,
					[]aggregations.ConfigFeatureDto{{}},
				),
			},
			want: want{
				err:    scheduler.ErrInvalidScheduler,
				result: aggregations.Config{},
			},
		},
		"should be validation features err": {
			input: input{
				dto: NewCreateConfigDto(
					"Test-A",
					"* * * * *",
					nil,
					time.Time{},
					nil,
					[]aggregations.ConfigFeatureDto{},
				),
			},
			want: want{
				err:    scheduler.ErrInvalidScheduler,
				result: aggregations.Config{},
			},
		},
		"should be create job err": {
			input: input{
				dto: NewCreateConfigDto(
					"Test-A",
					"* 1 * * *",
					nil,
					time.Time{},
					nil,
					[]aggregations.ConfigFeatureDto{{}},
				),
			},
			output: output{
				CreateJobErr:        errors.New("err create job"),
				CreateConfigErr:     nil,
				CreateConfig:        aggregations.Config{},
				GetFeaturesByIdsErr: nil,
				GetFeaturesByIds: []aggregations.Features{
					{
						ID:    "goodID",
						Name:  "goodName",
						Field: "goodField",
					},
				},
			},
			want: want{
				err:    errors.New("err create job"),
				result: aggregations.Config{},
			},
		},
		"should be create config err": {
			input: input{
				dto: NewCreateConfigDto(
					"Test-A",
					"* 1 * * *",
					nil,
					time.Time{},
					nil,
					[]aggregations.ConfigFeatureDto{{}},
				),
			},
			output: output{
				CreateJobId:         "1234",
				CreateConfigErr:     errors.New("err create config"),
				CreateConfig:        aggregations.Config{},
				GetFeaturesByIdsErr: nil,
				GetFeaturesByIds: []aggregations.Features{
					{
						ID:    "goodID",
						Name:  "goodName",
						Field: "goodField",
					},
				},
			},
			want: want{
				err:    errors.New("err create config"),
				result: aggregations.Config{},
			},
		},
		"should be create config": {
			input: input{
				GetFeaturesByIds: []string{"GoodId"},
				dto: NewCreateConfigDto(
					"Test-A",
					"* 1 * * *",
					nil,
					time.Time{},
					nil,
					[]aggregations.ConfigFeatureDto{{}},
				),
			},
			output: output{
				CreateJobId:     "1234",
				CreateConfigErr: nil,
				CreateConfig: aggregations.Config{
					Name:        "Test-A",
					Scheduler:   "* 1 * * *",
					SchedulerId: "1234",
					StartDate:   time.Time{},
					EndDate:     time.Time{},
					Description: "",
					Features: []aggregations.Features{{
						ID:    "goodID",
						Name:  "goodName",
						Field: "goodField",
					},
					},
				},
				GetFeaturesByIdsErr: nil,
				GetFeaturesByIds: []aggregations.Features{
					{
						ID:    "goodID",
						Name:  "goodName",
						Field: "goodField",
					},
				},
			},
			want: want{
				err: nil,
				result: aggregations.Config{
					Name:        "Test-A",
					Scheduler:   "* 1 * * *",
					SchedulerId: "1234",
					StartDate:   time.Time{},
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
			config, _ := aggregations.NewConfig(
				test.input.dto.Name,
				test.input.dto.Description,
				test.input.dto.Scheduler,
				test.input.dto.StartDate,
				test.input.dto.EndDate,
				test.output.GetFeaturesByIds)

			featuresRepo := new(mocks.AggregationsFeaturesRepository)

			featuresRepo.On("GetFeaturesByIds", mock.Anything, utils.MapSlice(test.input.dto.Features, func(item aggregations.ConfigFeatureDto) string {
				return item.Id
			})).Return(test.output.GetFeaturesByIds, test.output.GetFeaturesByIdsErr)

			aggregationRepository := new(mocks.AggregationConfigRepository)
			aggregationRepository.On("SaveAggregationConfig", mock.Anything, mock.MatchedBy(func(compare aggregations.Config) bool {
				config.Id = compare.Id
				_, err := utils.Parse(config.Id)
				return err == nil
			})).Return(test.output.CreateConfig, test.output.CreateConfigErr)

			schedulerClientMock := new(mocks2.Client)

			schedulerClient := func(ctx context.Context) (scheduler.Client, error) {
				schedulerClientMock.Mock.On("Close").Return(nil)
				schedulerClientMock.Mock.On("DeleteJob", mock.Anything, test.output.CreateJobId).Return(nil)
				schedulerClientMock.Mock.On("CreateJob", mock.Anything, mock.MatchedBy(func(compare aggregations.Config) bool {
					config.Id = compare.Id
					_, err := utils.Parse(config.Id)
					return err == nil
				}), "topic").Return(test.output.CreateJobId, test.output.CreateJobErr)
				return schedulerClientMock, nil
			}

			s := NewCreateAggregationConfigService(aggregationRepository, featuresRepo, schedulerClient, "topic", loc)

			result, err := s.Handler(context.Background(), test.input.dto)

			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.result, result)

		})
	}
}
