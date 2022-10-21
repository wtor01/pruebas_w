package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Unit_Service_Aggregations_GetFeatures(t *testing.T) {
	type input struct {
		ctx context.Context
		dto GetFeaturesDTO
	}
	type output struct {
		getFeaturesResponse      aggregations.Features
		getFeaturesErrorResponse error
	}

	tests := map[string]struct {
		input  input
		output output
	}{
		"Should be okay": {
			input: input{
				ctx: context.Background(),
				dto: GetFeaturesDTO{
					ID: "aeb95d60-adfb-429d-acc5-433d3565227d",
				},
			},
			output: output{
				getFeaturesResponse: aggregations.Features{
					ID:    "aeb95d60-adfb-429d-acc5-433d3565227d",
					Name:  "GoodName",
					Field: "GoodField",
				},
				getFeaturesErrorResponse: nil,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			repo := new(mocks.AggregationsFeaturesRepository)
			repo.On("GetFeatures", mock.Anything, testCase.input.dto.ID).Return(testCase.output.getFeaturesResponse, testCase.output.getFeaturesErrorResponse)

			srv := NewGetFeaturesService(repo)
			response, err := srv.Handler(testCase.input.ctx, testCase.input.dto)
			assert.Equal(t, testCase.output.getFeaturesResponse, response)
			assert.Equal(t, testCase.output.getFeaturesErrorResponse, err)

		})
	}

}
