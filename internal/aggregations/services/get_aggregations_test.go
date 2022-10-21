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

func Test_Unit_Measure_Services_GetAggregations_Handler(t *testing.T) {

	loc, _ := time.LoadLocation("Europe/Madrid")

	tests := map[string]struct {
		want                      aggregations.AggregationPrevious
		wantErr                   error
		input                     aggregations.GetAggregationDto
		repoGetDashboardResultErr error
		clientResult              aggregations.AggregationPrevious
		clientResultErr           error
	}{
		"should return err if GetDashboard fail": {
			want:                      aggregations.AggregationPrevious{},
			wantErr:                   errors.New(""),
			input:                     aggregations.GetAggregationDto{},
			repoGetDashboardResultErr: errors.New(""),
			clientResult:              aggregations.AggregationPrevious{},
			clientResultErr:           nil,
		},
		"measureShouldBe": {
			want:                      aggregations.AggregationPrevious{},
			wantErr:                   nil,
			repoGetDashboardResultErr: nil,
			input:                     aggregations.GetAggregationDto{},
			clientResult:              aggregations.AggregationPrevious{},
			clientResultErr:           nil,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			dashboardRepository := new(mocks.AggregationMongoRepository)

			dashboardRepository.On("GetPreviousAggregation", mock.Anything, aggregations.GetAggregationDto{AggregationConfigId: testCase.input.AggregationConfigId}).Return(testCase.clientResult, testCase.repoGetDashboardResultErr)

			h := NewAggregationsService(dashboardRepository, loc)
			result, err := h.GetAggregation.Handler(ctx, testCase.input)

			assert.Equal(t, testCase.want, result, testCase)
			assert.Equal(t, testCase.wantErr, err, testCase)
		})
	}
}
