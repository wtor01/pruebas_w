package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Unit_Service_Aggregations_ListFeatures(t *testing.T) {
	type input struct {
		ctx context.Context
		dto ListFeaturesDto
	}
	type output struct {
		listFeaturesResponse    []aggregations.Features
		listFeaturesIntResponse int
		listFeaturesErrResponse error
	}

	tests := map[string]struct {
		input  input
		output output
	}{
		"Should be okay": {
			input: input{
				ctx: context.Background(),
				dto: ListFeaturesDto{},
			},
			output: output{
				listFeaturesResponse: []aggregations.Features{
					{
						ID:    "aeb95d60-adfb-429d-acc5-433d3565227d",
						Name:  "GoodName",
						Field: "GoodField",
					},
					{
						ID:    "aeb95d60-adfb-429d-acc5-433d3565227d",
						Name:  "GoodName",
						Field: "GoodField",
					},
				},
				listFeaturesIntResponse: 2,
				listFeaturesErrResponse: nil,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			repo := new(mocks.AggregationsFeaturesRepository)

			repo.On("ListFeatures", mock.Anything, db.Pagination{
				Limit:  testCase.input.dto.Limit,
				Offset: testCase.input.dto.Offset,
			}).Return(testCase.output.listFeaturesResponse, testCase.output.listFeaturesIntResponse, testCase.output.listFeaturesErrResponse)
			srv := NewListFeatures(repo)

			response, number, err := srv.Handler(testCase.input.ctx, testCase.input.dto)
			assert.Equal(t, testCase.output.listFeaturesResponse, response)
			assert.Equal(t, testCase.output.listFeaturesIntResponse, number)
			assert.Equal(t, testCase.output.listFeaturesErrResponse, err)

		})
	}

}
