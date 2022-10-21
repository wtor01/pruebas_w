package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Unit_Measure_Services_DashboardStats_Handle(t *testing.T) {

	tests := map[string]struct {
		want                      []gross_measures.GrossMeasuresDashboardStatsGlobal
		wantErr                   error
		input                     GetDashboardMeasuresStatsDTO
		repoGetDashboardResult    []gross_measures.GrossMeasuresDashboardStatsGlobal
		repoGetDashboardResultErr error
	}{
		"should return err if GetDashboardStats fail": {
			want:                      []gross_measures.GrossMeasuresDashboardStatsGlobal{},
			wantErr:                   errors.New("invalid input format"),
			input:                     GetDashboardMeasuresStatsDTO{},
			repoGetDashboardResult:    nil,
			repoGetDashboardResultErr: errors.New(""),
		},
		"should be ok": {
			want: []gross_measures.GrossMeasuresDashboardStatsGlobal{{
				DistributorID: "32131",
				Month:         9,
				Year:          2022,
				Type:          "TLG",
				DailyStats: []gross_measures.DashboardSingleStats{{
					Day:                    0,
					HourlyCurve:            0,
					QuarterlyCurve:         0,
					DailyClosure:           0,
					MonthlyClosure:         0,
					ExpectedHourlyCurve:    0,
					ExpectedQuarterlyCurve: 0,
					ExpectedDailyClosure:   0,
					ExpectedMonthlyClosure: 0,
				}},
			}},
			wantErr: nil,
			repoGetDashboardResult: []gross_measures.GrossMeasuresDashboardStatsGlobal{{
				DistributorID: "32131",
				Month:         9,
				Year:          2022,
				Type:          "TLG",
				DailyStats: []gross_measures.DashboardSingleStats{{
					Day:                    0,
					HourlyCurve:            0,
					QuarterlyCurve:         0,
					DailyClosure:           0,
					MonthlyClosure:         0,
					ExpectedHourlyCurve:    0,
					ExpectedQuarterlyCurve: 0,
					ExpectedDailyClosure:   0,
					ExpectedMonthlyClosure: 0,
				}},
			}},
			repoGetDashboardResultErr: nil,
			input: GetDashboardMeasuresStatsDTO{
				DistributorID: "32131",
				Month:         9,
				Year:          2022,
				Type:          "TLG",
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			dashboardRepository := new(mocks.GrossMeasuresDashboardStatsRepository)

			dashboardRepository.On("GetStatisticsGlobal", mock.Anything, gross_measures.SearchDashboardStats{
				DistributorID: testCase.input.DistributorID,
				Month:         testCase.input.Month,
				Year:          testCase.input.Year,
				Type:          testCase.input.Type,
			}).Return(testCase.repoGetDashboardResult, testCase.repoGetDashboardResultErr)

			h := NewGetDashboardMeasuresStats(dashboardRepository)

			result, err := h.Handler(ctx, testCase.input)

			assert.Equal(t, testCase.want, result, testCase)
			assert.Equal(t, testCase.wantErr, err, testCase)
		})
	}
}
