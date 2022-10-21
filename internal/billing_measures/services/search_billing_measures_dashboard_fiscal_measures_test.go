package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_Services_Billing_Measures_Dashboard_Search(t *testing.T) {
	initDate, _ := time.Parse("2006-01-02 15:04", "2022-04-30 22:00")
	endDate, _ := time.Parse("2006-01-02 15:04", "2022-05-31 22:00")
	type input struct {
		ctx context.Context
		dto SearchFiscalBillingMeasuresDashboardDTO
	}
	type expected struct {
		expectedErr    error
		expectedResult []billing_measures.FiscalBillingMeasuresDashboard
	}
	type output struct {
		searchFiscalBillingMeasuresResponse    []billing_measures.BillingMeasure
		searchLastBillingMeasuresResponse      billing_measures.BillingMeasure
		searchFiscalBillingMeasuresResponseErr error
	}

	tests := map[string]struct {
		input    input
		expected expected
		output   output
	}{
		"Should return okay": {
			input: input{
				ctx: context.Background(),
				dto: SearchFiscalBillingMeasuresDashboardDTO{
					Id:            nil,
					Cups:          "ES0130000000357054DJ",
					DistributorId: "DistributorID",
					StartDate:     initDate,
					EndDate:       endDate,
				},
			},
			expected: expected{
				expectedErr: nil,
				expectedResult: []billing_measures.FiscalBillingMeasuresDashboard{
					{
						Id:                 "",
						Status:             billing_measures.Calculated,
						LastMvDate:         endDate.Format("02-01-2006 15:04"),
						Cups:               "ES0130000000357054DJ",
						StartDate:          initDate.Format("02-01-2006"),
						EndDate:            endDate.Format("02-01-2006"),
						PrincipalMagnitude: measures.AI,
						Periods:            []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
						Magnitudes:         []measures.Magnitude{measures.AI},
						Summary: billing_measures.Summary{
							Real: billing_measures.SummaryItem{
								Total: 50,
								P1:    100,
								P2:    0,
								P3:    50,
							},
							Adjusted: billing_measures.SummaryItem{
								Total: 41.67,
								P1:    0,
								P2:    100,
								P3:    25,
							},
							Outlined: billing_measures.SummaryItem{
								Total: 8.33,
								P1:    0,
								P2:    0,
								P3:    25,
							},
							Calculated: billing_measures.SummaryItem{},
							Estimated:  billing_measures.SummaryItem{},
							Consum: billing_measures.SummaryItem{
								Total: 180,
								P1:    20,
								P2:    40,
								P3:    120,
							},
						},
						CalendarCurve: []billing_measures.CalendarCurve{
							{
								Date:   "2022-05-01",
								Status: billing_measures.GeneralAdjusted,
							},
						},
						Balance: billing_measures.Balance{
							Method: billing_measures.GeneralOutlined,
							P0: &billing_measures.BalancePeriod{
								AI:            180,
								BalanceTypeAI: billing_measures.GeneralOutlined,
							},
							P1: &billing_measures.BalancePeriod{
								AI:            20,
								BalanceTypeAI: billing_measures.GeneralReal,
							},
							P2: &billing_measures.BalancePeriod{
								AI:            40,
								BalanceTypeAI: billing_measures.GeneralAdjusted,
							},
							P3: &billing_measures.BalancePeriod{
								AI:            120,
								BalanceTypeAI: billing_measures.GeneralOutlined,
							},
							P4: nil,
							P5: nil,
							P6: nil,
						},
						Curve: []billing_measures.Curve{
							{
								Date: "2022-05-01",
								Values: []billing_measures.Values{
									{
										Hour:   "00:00",
										Status: billing_measures.GeneralReal,
										AI:     10,
									},
									{
										Hour:   "01:00",
										Status: billing_measures.GeneralReal,
										AI:     10,
									},
									{
										Hour:   "02:00",
										Status: billing_measures.GeneralAdjusted,
										AI:     20,
									},
									{
										Hour:   "03:00",
										Status: billing_measures.GeneralAdjusted,
										AI:     20,
									},
									{
										Hour:   "04:00",
										Status: billing_measures.GeneralOutlined,
										AI:     30,
									},
									{
										Hour:   "05:00",
										Status: billing_measures.GeneralAdjusted,
										AI:     30,
									},
									{
										Hour:   "06:00",
										Status: billing_measures.GeneralReal,
										AI:     30,
									},
									{
										Hour:   "07:00",
										Status: billing_measures.GeneralReal,
										AI:     30,
									},
								},
							},
						},
						ExecutionSummary: billing_measures.ExecutionSummary{
							BalanceType:   billing_measures.GeneralAdjusted,
							CurveType:     billing_measures.GeneralAdjusted,
							CurveStatus:   measures.Absent,
							BalanceOrigin: measures.Monthly,
						},
					},
				},
			},
			output: output{
				searchFiscalBillingMeasuresResponse: []billing_measures.BillingMeasure{
					{
						CUPS:     "ES0130000000357054DJ",
						EndDate:  endDate,
						InitDate: initDate,
						ActualReadingClosure: measures.DailyReadingClosure{
							ServiceType: "D-C",
							Origin:      measures.STG,
							ClosureType: measures.Monthly,
							EndDate:     time.Date(2022, 1, 31, 22, 0, 0, 0, time.UTC),
						},
						Status: billing_measures.Calculated,
						BillingBalance: billing_measures.BillingBalance{
							P0: &billing_measures.BillingBalancePeriod{
								AI:                   180,
								BalanceTypeAI:        billing_measures.Outlined,
								BalanceGeneralTypeAI: billing_measures.GeneralOutlined,
							},
							P1: &billing_measures.BillingBalancePeriod{
								AI:                   20,
								BalanceTypeAI:        billing_measures.RealBalance,
								BalanceGeneralTypeAI: billing_measures.GeneralReal,
							},
							P2: &billing_measures.BillingBalancePeriod{
								AI:                   40,
								BalanceTypeAI:        billing_measures.Adjustment,
								BalanceGeneralTypeAI: billing_measures.GeneralAdjusted,
							},
							P3: &billing_measures.BillingBalancePeriod{
								AI:                   120,
								BalanceTypeAI:        billing_measures.Outlined,
								BalanceGeneralTypeAI: billing_measures.GeneralOutlined,
							},
							P4: nil,
							P5: nil,
							P6: nil,
						},
						BillingLoadCurve: []billing_measures.BillingLoadCurve{
							{
								EndDate:                  time.Date(2022, 4, 30, 22, 0, 0, 0, time.UTC),
								Period:                   measures.P1,
								AI:                       10,
								EstimatedGeneralMethodAI: billing_measures.GeneralReal,
							},
							{
								EndDate:                  time.Date(2022, 4, 30, 23, 0, 0, 0, time.UTC),
								Period:                   measures.P1,
								AI:                       10,
								EstimatedGeneralMethodAI: billing_measures.GeneralReal,
							},
							{
								EndDate:                  time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC),
								Period:                   measures.P2,
								AI:                       20,
								EstimatedGeneralMethodAI: billing_measures.GeneralAdjusted,
							},
							{
								EndDate:                  time.Date(2022, 5, 1, 1, 0, 0, 0, time.UTC),
								Period:                   measures.P2,
								AI:                       20,
								EstimatedGeneralMethodAI: billing_measures.GeneralAdjusted,
							},
							{
								EndDate:                  time.Date(2022, 5, 1, 2, 0, 0, 0, time.UTC),
								Period:                   measures.P3,
								AI:                       30,
								EstimatedGeneralMethodAI: billing_measures.GeneralOutlined,
							},
							{
								EndDate:                  time.Date(2022, 5, 1, 3, 0, 0, 0, time.UTC),
								Period:                   measures.P3,
								AI:                       30,
								EstimatedGeneralMethodAI: billing_measures.GeneralAdjusted,
							},
							{
								EndDate:                  time.Date(2022, 5, 1, 4, 0, 0, 0, time.UTC),
								Period:                   measures.P3,
								AI:                       30,
								EstimatedGeneralMethodAI: billing_measures.GeneralReal,
							},
							{
								EndDate:                  time.Date(2022, 5, 1, 5, 0, 0, 0, time.UTC),
								Period:                   measures.P3,
								AI:                       30,
								EstimatedGeneralMethodAI: billing_measures.GeneralReal,
							},
						},
						Periods: []measures.PeriodKey{
							measures.P1,
							measures.P2,
							measures.P3,
						},
						Magnitudes: []measures.Magnitude{measures.AI},
						ExecutionSummary: billing_measures.ExecutionSummary{
							BalanceType:   billing_measures.GeneralAdjusted,
							CurveType:     billing_measures.GeneralAdjusted,
							CurveStatus:   measures.Absent,
							BalanceOrigin: measures.Monthly,
						},
					},
				},
				searchLastBillingMeasuresResponse: billing_measures.BillingMeasure{
					EndDate: endDate,
				},
				searchFiscalBillingMeasuresResponseErr: nil,
			},
		},
	}

	for testsName, _ := range tests {
		testCase := tests[testsName]
		t.Run(testsName, func(t *testing.T) {
			t.Parallel()
			repo := new(mocks.BillingMeasuresDashboardRepository)
			repo.On("SearchFiscalBillingMeasures", mock.Anything, testCase.input.dto.Cups, testCase.input.dto.DistributorId, testCase.input.dto.StartDate, testCase.input.dto.EndDate).Return(testCase.output.searchFiscalBillingMeasuresResponse, testCase.output.searchFiscalBillingMeasuresResponseErr)

			repo.On("SearchLastBillingMeasures", mock.Anything, testCase.input.dto.Cups, testCase.input.dto.DistributorId).Return(testCase.output.searchLastBillingMeasuresResponse, testCase.output.searchFiscalBillingMeasuresResponseErr)

			location, _ := time.LoadLocation("Europe/Madrid")
			srv := NewSearchFiscalBillingMeasuresDashboard(repo, location)
			result, err := srv.Handler(testCase.input.ctx, testCase.input.dto)
			assert.Equal(t, testCase.expected.expectedResult, result)
			assert.Equal(t, testCase.expected.expectedErr, err)

		})

	}
}

//TODO:GENERAR MOCKERY
