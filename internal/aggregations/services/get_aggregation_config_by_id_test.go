package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_Service_Aggregation_GetConfigById_Handler(t *testing.T) {
	type input struct {
		configId string
	}

	type output struct {
		GetConfig    aggregations.Config
		GetConfigErr error
	}

	type want struct {
		result aggregations.Config
		err    error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be err": {
			input: input{
				configId: "not Valid",
			},
			output: output{
				GetConfig:    aggregations.Config{},
				GetConfigErr: errors.New("err"),
			},
			want: want{
				result: aggregations.Config{},
				err:    errors.New("err"),
			},
		},
		"Should be find configs": {
			input: input{
				configId: "12345",
			},
			output: output{
				GetConfig: aggregations.Config{
					Id:        "12345",
					Name:      "Test 1",
					Scheduler: "* * * * *",
					StartDate: time.Date(2022, 6, 30, 0, 0, 0, 0, time.UTC),
					Features:  []aggregations.Features{{}},
				},
				GetConfigErr: nil,
			},
			want: want{
				result: aggregations.Config{
					Id:        "12345",
					Name:      "Test 1",
					Scheduler: "* * * * *",
					StartDate: time.Date(2022, 6, 30, 0, 0, 0, 0, time.UTC),
					Features:  []aggregations.Features{{}},
				},
				err: nil,
			},
		},
	}

	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			aggregationRepository := new(mocks.AggregationConfigRepository)
			aggregationRepository.On("GetAggregationConfigById", mock.Anything, test.input.configId).Return(test.output.GetConfig, test.output.GetConfigErr)

			s := NewGetAggregationConfigByIdService(aggregationRepository)
			result, err := s.Handler(context.Background(), test.input.configId)

			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.result, result)
		})
	}
}
