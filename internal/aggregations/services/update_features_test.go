package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Unit_Service_Aggregations_UpdateFeatures(t *testing.T) {
	type input struct {
		ctx context.Context
		dto UpdateFeaturesDTO
	}
	type output struct {
		searchFeaturesResponse      []aggregations.Features
		searchFeaturesErrorResponse error
		saveFeaturesErrResponse     error
		getFeaturesResponse         aggregations.Features
		getFeaturesErrResponse      error
	}
	type expected struct {
		featuresExpected    aggregations.Features
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
				dto: UpdateFeaturesDTO{
					ID:    "aeb95d60-adfb-429d-acc5-433d3565227d",
					Name:  "Change*",
					Field: "Change*",
				},
			},
			output: output{
				getFeaturesResponse: aggregations.Features{
					ID:    "aeb95d60-adfb-429d-acc5-433d3565227d",
					Name:  "GoodName",
					Field: "GoodField",
				},
				getFeaturesErrResponse: nil,
				searchFeaturesResponse: []aggregations.Features{
					{
						ID:    "aeb95d60-adfb-429d-acc5-433d3565227d",
						Name:  "GoodName",
						Field: "GoodField",
					},
				},
				searchFeaturesErrorResponse: nil,
				saveFeaturesErrResponse:     nil,
			},
			expected: expected{
				featuresExpected: aggregations.Features{
					ID:    "aeb95d60-adfb-429d-acc5-433d3565227d",
					Name:  "Change*",
					Field: "Change*",
				},
				featuresErrExpected: nil,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			repo := new(mocks.AggregationsFeaturesRepository)
			repo.On("GetFeatures", mock.Anything, testCase.input.dto.ID).Return(testCase.output.getFeaturesResponse, testCase.output.getFeaturesErrResponse)

			repo.On("SearchFeatures", mock.Anything, aggregations.SearchFeatures{
				Name:  testCase.input.dto.Name,
				Field: testCase.input.dto.Field,
			}).Return(testCase.output.searchFeaturesResponse, testCase.output.searchFeaturesErrorResponse)
			
			repo.On("SaveFeatures", mock.Anything, aggregations.Features{
				ID:    testCase.input.dto.ID,
				Name:  testCase.input.dto.Name,
				Field: testCase.input.dto.Field,
			}).Return(testCase.output.saveFeaturesErrResponse)

			srv := NewUpdateFeaturesService(repo)

			response, err := srv.Handler(testCase.input.ctx, testCase.input.dto)

			assert.Equal(t, testCase.expected.featuresExpected, response)
			assert.Equal(t, testCase.expected.featuresErrExpected, err)

		})
	}

}
