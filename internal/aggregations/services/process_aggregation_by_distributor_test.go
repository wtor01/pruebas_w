package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Unit_Services_Aggregations_ProcessAggregationByDistributorService(t *testing.T) {
	type expected struct {
		generateAggregationResult []aggregations.Aggregation
		generateAggregationErr    error

		saveAllAggregationsErr error
	}

	tests := map[string]struct {
		input    aggregations.ConfigScheduler
		expected expected
		want     error
	}{
		"should return err if GenerateAggregation fail": {
			input: aggregations.ConfigScheduler{},
			expected: expected{
				generateAggregationResult: []aggregations.Aggregation{},
				generateAggregationErr:    errors.New(""),
			},
			want: errors.New(""),
		},
		"should not save if GenerateAggregation return 0 aggregations": {
			input: aggregations.ConfigScheduler{},
			expected: expected{
				generateAggregationResult: []aggregations.Aggregation{},
			},
			want: nil,
		},
		"should return err if SaveAllAggregations fail": {
			input: aggregations.ConfigScheduler{},
			expected: expected{
				generateAggregationResult: []aggregations.Aggregation{{
					Id: "test",
				}},
				generateAggregationErr: nil,
				saveAllAggregationsErr: errors.New(""),
			},
			want: errors.New(""),
		},
		"should return ok if SaveAllAggregations ok": {
			input: aggregations.ConfigScheduler{},
			expected: expected{
				generateAggregationResult: []aggregations.Aggregation{{
					Id: "test",
				}},
				generateAggregationErr: nil,
				saveAllAggregationsErr: nil,
			},
			want: nil,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {

			aggregationRepository := new(mocks.AggregationRepository)
			aggregationRepository.On(
				"GenerateAggregation",
				mock.Anything,
				testCase.input,
			).Return(
				testCase.expected.generateAggregationResult,
				testCase.expected.generateAggregationErr,
			)
			aggregationRepository.On(
				"SaveAllAggregations",
				mock.Anything,
				testCase.expected.generateAggregationResult,
			).Return(testCase.expected.saveAllAggregationsErr)

			svc := NewProcessAggregationByDistributorService(aggregationRepository)
			result := svc.Handler(context.Background(), testCase.input)
			
			calls := 0

			if len(testCase.expected.generateAggregationResult) != 0 {
				calls = 1
			}

			aggregationRepository.AssertNumberOfCalls(t, "SaveAllAggregations", calls)

			assert.Equal(t, testCase.want, result)
		})
	}
}
