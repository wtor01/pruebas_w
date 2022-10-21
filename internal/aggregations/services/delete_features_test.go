package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Unit_Service_Aggregations_DeleteFeatures(t *testing.T) {
	type input struct {
		ctx context.Context
		dto DeleteFeaturesDTO
	}
	type output struct {
		getFeaturesResponse       aggregations.Features
		getFeaturesErrorResponse  error
		deleteFeaturesErrResponse error
	}
	type expected struct {
		featuresErrExpected error
	}

	tests := map[string]struct {
		input    input
		output   output
		expected expected
	}{
		"Should be okay": {
			input: input{
				ctx: context.Background(),
				dto: DeleteFeaturesDTO{
					ID: "aeb95d60-adfb-429d-acc5-433d3565227d",
				},
			},
			output: output{
				getFeaturesResponse: aggregations.Features{
					ID:    "aeb95d60-adfb-429d-acc5-433d3565227d",
					Name:  "GoodName",
					Field: "GoodField",
				},
				getFeaturesErrorResponse:  nil,
				deleteFeaturesErrResponse: nil,
			},
			expected: expected{
				featuresErrExpected: nil,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			repo := new(mocks.AggregationsFeaturesRepository)
			repo.On("GetFeatures", mock.Anything, testCase.input.dto.ID).Return(testCase.output.getFeaturesResponse, testCase.output.getFeaturesErrorResponse)

			repo.On("DeleteFeatures", mock.Anything, testCase.input.dto.ID).Return(testCase.output.deleteFeaturesErrResponse)
			srv := NewDeleteFeaturesService(repo)
			response := srv.Handler(testCase.input.ctx, testCase.input.dto)
			assert.Equal(t, testCase.expected.featuresErrExpected, response)

		})
	}

}
