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

func Test_Unit_Service_Aggregation_GetAllConfigs_Handler(t *testing.T) {
	type input struct {
		dto GetAggregationConfigsServiceDto
	}

	type output struct {
		GetConfigs    []aggregations.Config
		Count         int
		GetConfigsErr error
	}

	type want struct {
		result []aggregations.Config
		count  int
		err    error
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be err": {
			input: input{
				dto: NewGetAggregationConfigsServiceDto("", 10, nil),
			},
			output: output{
				GetConfigs:    []aggregations.Config{},
				Count:         0,
				GetConfigsErr: errors.New("err"),
			},
			want: want{
				result: []aggregations.Config{},
				count:  0,
				err:    errors.New("err"),
			},
		},
		"Should be find configs": {
			input: input{
				dto: NewGetAggregationConfigsServiceDto("", 10, nil),
			},
			output: output{
				GetConfigs: []aggregations.Config{
					{
						Id:        "12345",
						Name:      "Test 1",
						Scheduler: "* * * * *",
						StartDate: time.Date(2022, 6, 30, 0, 0, 0, 0, time.UTC),
						Features:  []aggregations.Features{},
					},
					{
						Id:        "1234",
						Name:      "Test 2",
						Scheduler: "* 1 * 2 *",
						StartDate: time.Date(2022, 7, 31, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2022, 8, 31, 0, 0, 0, 0, time.UTC),
						Features: []aggregations.Features{
							{
								ID:    "123",
								Name:  "DistributorCode",
								Field: "distributor_code",
							},
						},
					},
				},
				Count:         2,
				GetConfigsErr: nil,
			},
			want: want{
				result: []aggregations.Config{
					{
						Id:        "12345",
						Name:      "Test 1",
						Scheduler: "* * * * *",
						StartDate: time.Date(2022, 6, 30, 0, 0, 0, 0, time.UTC),
						Features:  []aggregations.Features{},
					},
					{
						Id:        "1234",
						Name:      "Test 2",
						Scheduler: "* 1 * 2 *",
						StartDate: time.Date(2022, 7, 31, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2022, 8, 31, 0, 0, 0, 0, time.UTC),
						Features: []aggregations.Features{
							{
								ID:    "123",
								Name:  "DistributorCode",
								Field: "distributor_code",
							},
						},
					},
				},
				count: 2,
				err:   nil,
			},
		},
	}

	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			aggregationRepository := new(mocks.AggregationConfigRepository)
			aggregationRepository.On("GetAggregationConfigs", mock.Anything, aggregations.GetConfigsQuery{
				Query:  test.input.dto.Q,
				Limit:  test.input.dto.Limit,
				Offset: test.input.dto.Offset,
			}).Return(test.output.GetConfigs, test.output.Count, test.output.GetConfigsErr)

			s := NewGetAggregationConfigService(aggregationRepository)
			result, count, err := s.Handler(context.Background(), test.input.dto)

			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.result, result)
			assert.Equal(t, test.want.count, count)
		})
	}
}
