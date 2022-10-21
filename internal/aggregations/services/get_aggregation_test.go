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

func Test_Unit_Measure_Services_GetAggregation_Handler(t *testing.T) {

	loc, _ := time.LoadLocation("Europe/Madrid")

	tests := map[string]struct {
		want                      []aggregations.Aggregation
		wantErr                   error
		input                     aggregations.GetAggregationsDto
		repoGetDashboardResultErr error
		countResult               int64
		clientResult              []aggregations.Aggregation
		clientResultErr           error
	}{
		"should return err if GetDashboard fail": {
			want:                      []aggregations.Aggregation{},
			wantErr:                   errors.New(""),
			input:                     aggregations.GetAggregationsDto{},
			countResult:               0,
			repoGetDashboardResultErr: errors.New(""),
			clientResult:              []aggregations.Aggregation{},
			clientResultErr:           nil,
		},
		"measureShouldBe": {
			want:                      []aggregations.Aggregation{},
			wantErr:                   nil,
			repoGetDashboardResultErr: nil,
			input:                     aggregations.GetAggregationsDto{},
			countResult:               0,
			clientResult:              []aggregations.Aggregation{},
			clientResultErr:           nil,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			dashboardRepository := new(mocks.AggregationMongoRepository)

			dashboardRepository.On("GetAggregations", mock.Anything, aggregations.GetAggregationsDto{
				Offset:              nil,
				Limit:               10,
				AggregationConfigId: "22",
				StartDate:           time.Time{},
				EndDate:             time.Time{},
			}).Return(testCase.clientResult, testCase.countResult, testCase.repoGetDashboardResultErr)

			h := NewAggregationsService(dashboardRepository, loc)

			result, _, err := h.GetAggregations.Handler(ctx, aggregations.GetAggregationsDto{
				Offset:              nil,
				Limit:               10,
				AggregationConfigId: "22",
				StartDate:           time.Time{},
				EndDate:             time.Time{},
			})

			assert.Equal(t, testCase.want, result, testCase)
			assert.Equal(t, testCase.wantErr, err, testCase)
		})
	}
}
