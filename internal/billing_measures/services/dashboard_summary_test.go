package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_Service_DashboardSummary_Handler(t *testing.T) {
	type input struct {
		dto DashboardSummaryDto
	}

	type output struct {
		GroupSummaryResult billing_measures.FiscalMeasureSummary
		GroupSummaryErr    error
	}

	type want struct {
		err    error
		result billing_measures.FiscalMeasureSummary
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be bbdd err": {
			input: input{dto: NewDashboardSummaryDto(
				"DistributorX",
				"TLG",
				time.Date(2022, 7, 31, 22, 0, 0, 0, time.UTC),
				time.Date(2022, 8, 31, 22, 0, 0, 0, time.UTC),
			)},
			output: output{
				GroupSummaryErr: errors.New("err"),
			},
			want: want{
				err: errors.New("err"),
			},
		},
		"Should be correct TLG": {
			input: input{dto: NewDashboardSummaryDto(
				"DistributorX",
				"TLG",
				time.Date(2022, 7, 31, 22, 0, 0, 0, time.UTC),
				time.Date(2022, 8, 31, 22, 0, 0, 0, time.UTC),
			)},
			output: output{
				GroupSummaryResult: billing_measures.FiscalMeasureSummary{
					MeterType: measures.TLG,
					BalanceType: billing_measures.BalanceTypeSummary{
						TypeSummary: billing_measures.TypeSummary{
							Real:       95,
							Calculated: 2,
							Estimated:  3,
						},
					},
					BalanceOrigin: billing_measures.BalanceOriginSummary{
						Monthly:   100,
						Daily:     0,
						Other:     0,
						NoClosure: 0,
					},
					CurveType: billing_measures.CurveTypeSummary{
						TypeSummary: billing_measures.TypeSummary{
							Real:       78,
							Calculated: 0,
							Estimated:  0,
						},
						Adjusted: 2,
						Outlined: 20,
					},
					CurveStatus: billing_measures.CurveStatusSummary{
						Completed:    90,
						NotCompleted: 5,
						Absent:       5,
					},
					Total: 100,
				},
			},
			want: want{
				err: nil,
				result: billing_measures.FiscalMeasureSummary{
					MeterType: measures.TLG,
					BalanceType: billing_measures.BalanceTypeSummary{
						TypeSummary: billing_measures.TypeSummary{
							Real:       95,
							Calculated: 2,
							Estimated:  3,
						},
					},
					BalanceOrigin: billing_measures.BalanceOriginSummary{
						Monthly:   100,
						Daily:     0,
						Other:     0,
						NoClosure: 0,
					},
					CurveType: billing_measures.CurveTypeSummary{
						TypeSummary: billing_measures.TypeSummary{
							Real:       78,
							Calculated: 0,
							Estimated:  0,
						},
						Adjusted: 2,
						Outlined: 20,
					},
					CurveStatus: billing_measures.CurveStatusSummary{
						Completed:    90,
						NotCompleted: 5,
						Absent:       5,
					},
					Total: 100,
				},
			},
		},
		"Should be correct TLM": {
			input: input{dto: NewDashboardSummaryDto(
				"DistributorX",
				"TLM",
				time.Date(2022, 7, 31, 22, 0, 0, 0, time.UTC),
				time.Date(2022, 8, 31, 22, 0, 0, 0, time.UTC),
			)},
			output: output{
				GroupSummaryResult: billing_measures.FiscalMeasureSummary{
					MeterType: measures.TLM,
					BalanceType: billing_measures.BalanceTypeSummary{
						TypeSummary: billing_measures.TypeSummary{
							Real:       95,
							Calculated: 2,
							Estimated:  3,
						},
					},
					BalanceOrigin: billing_measures.BalanceOriginSummary{
						Monthly:   100,
						Daily:     0,
						Other:     0,
						NoClosure: 0,
					},
					CurveType: billing_measures.CurveTypeSummary{
						TypeSummary: billing_measures.TypeSummary{
							Real:       78,
							Calculated: 0,
							Estimated:  0,
						},
						Adjusted: 2,
						Outlined: 20,
					},
					CurveStatus: billing_measures.CurveStatusSummary{
						Completed:    90,
						NotCompleted: 5,
						Absent:       5,
					},
					Total: 100,
				},
			},
			want: want{
				err: nil,
				result: billing_measures.FiscalMeasureSummary{
					MeterType: measures.TLM,
					BalanceType: billing_measures.BalanceTypeSummary{
						TypeSummary: billing_measures.TypeSummary{
							Real:       95,
							Calculated: 2,
							Estimated:  3,
						},
					},
					BalanceOrigin: billing_measures.BalanceOriginSummary{
						Monthly:   100,
						Daily:     0,
						Other:     0,
						NoClosure: 0,
					},
					CurveType: billing_measures.CurveTypeSummary{
						TypeSummary: billing_measures.TypeSummary{
							Real:       78,
							Calculated: 0,
							Estimated:  0,
						},
						Adjusted: 2,
						Outlined: 20,
					},
					CurveStatus: billing_measures.CurveStatusSummary{
						Completed:    90,
						NotCompleted: 5,
						Absent:       5,
					},
					Total: 100,
				},
			},
		},
		"Should be correct OTHER": {
			input: input{dto: NewDashboardSummaryDto(
				"DistributorX",
				"OTHER",
				time.Date(2022, 7, 31, 22, 0, 0, 0, time.UTC),
				time.Date(2022, 8, 31, 22, 0, 0, 0, time.UTC),
			)},
			output: output{
				GroupSummaryResult: billing_measures.FiscalMeasureSummary{
					MeterType: measures.OTHER,
					BalanceType: billing_measures.BalanceTypeSummary{
						TypeSummary: billing_measures.TypeSummary{
							Real:       95,
							Calculated: 2,
							Estimated:  3,
						},
					},
					BalanceOrigin: billing_measures.BalanceOriginSummary{
						Monthly:   100,
						Daily:     0,
						Other:     0,
						NoClosure: 0,
					},
					CurveType: billing_measures.CurveTypeSummary{
						TypeSummary: billing_measures.TypeSummary{
							Real:       78,
							Calculated: 0,
							Estimated:  0,
						},
						Adjusted: 2,
						Outlined: 20,
					},
					CurveStatus: billing_measures.CurveStatusSummary{
						Completed:    90,
						NotCompleted: 5,
						Absent:       5,
					},
					Total: 100,
				},
			},
			want: want{
				err: nil,
				result: billing_measures.FiscalMeasureSummary{
					MeterType: measures.OTHER,
					BalanceType: billing_measures.BalanceTypeSummary{
						TypeSummary: billing_measures.TypeSummary{
							Real:       95,
							Calculated: 2,
							Estimated:  3,
						},
					},
					BalanceOrigin: billing_measures.BalanceOriginSummary{
						Monthly:   100,
						Daily:     0,
						Other:     0,
						NoClosure: 0,
					},
					CurveType: billing_measures.CurveTypeSummary{
						TypeSummary: billing_measures.TypeSummary{
							Real:       78,
							Calculated: 0,
							Estimated:  0,
						},
						Adjusted: 2,
						Outlined: 20,
					},
					CurveStatus: billing_measures.CurveStatusSummary{
						Completed:    90,
						NotCompleted: 5,
						Absent:       5,
					},
					Total: 100,
				},
			},
		},
	}

	for name, _ := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			dashboardRepository := new(mocks.BillingMeasuresDashboardRepository)
			dashboardRepository.On("GroupFiscalMeasureSummary", mock.Anything, billing_measures.GroupFiscalMeasureSummaryQuery{
				DistributorId: test.input.dto.DistributorId,
				MeterType:     test.input.dto.MeterType,
				StartDate:     test.input.dto.StartDate,
				EndDate:       test.input.dto.EndDate,
			}).Return(test.output.GroupSummaryResult, test.output.GroupSummaryErr)
			srv := NewDashboardSummaryService(dashboardRepository)

			result, err := srv.Handler(ctx, test.input.dto)

			assert.Equal(t, test.want.err, err, name)
			assert.Equal(t, test.want.result, result, name)

		})
	}
}
