package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/services/fixtures"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	clients_mocks "bitbucket.org/sercide/data-ingestion/internal/common/clients/clients_mocks"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	event_mocks "bitbucket.org/sercide/data-ingestion/pkg/event/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_Services_BillingMeasures_Process_mvh_by_cup_Handle(t *testing.T) {

	noLastBillingMeasure := billing_measures.NoLastBillingMeasure
	type want struct {
		resultLast                         billing_measures.BillingMeasure
		resultLastErr                      error
		Status                             billing_measures.Status
		DescriptionStatus                  billing_measures.DescriptionStatus
		GetPrevious                        billing_measures.BillingMeasure
		GetPreviousErr                     error
		result                             error
		result1Save                        error
		ListProcessedLoadCurveResult       []process_measures.ProcessedLoadCurve
		ListProcessedLoadCurveErr          error
		GetDailyReadingClosureByCupsResult measures.DailyReadingClosure
		GetDailyReadingClosureByCupsErr    error
	}

	tests := map[string]struct {
		input measures.ProcessMeasurePayload
		want  want
	}{
		"should return err if repo.Last return err": {
			input: measures.ProcessMeasurePayload{
				MeterConfig: measures.MeterConfig{
					Type: measures.TLG,
					ServicePoint: measures.ServicePoint{
						Cups: "1",
					},
				},
				Date: time.Now(),
			},
			want: want{
				resultLast:        billing_measures.BillingMeasure{},
				Status:            billing_measures.Supervision,
				DescriptionStatus: billing_measures.NoLastBillingMeasure,
				resultLastErr:     errors.New("err last"),
				result:            errors.New("err last"),
			},
		},
		"should return err if repo.Save return err": {
			input: measures.ProcessMeasurePayload{
				MeterConfig: measures.MeterConfig{
					Type: measures.TLG,
					ServicePoint: measures.ServicePoint{
						Cups: "1",
					},
				},
				Date: time.Now(),
			},
			want: want{
				Status: billing_measures.Calculating,
				resultLast: billing_measures.BillingMeasure{
					DistributorCode: "0130",
					DistributorID:   "id",
					CUPS:            "CUPS",
					EndDate:         time.Now().AddDate(0, -1, 0),
					InitDate:        time.Now().AddDate(0, -2, 0),
					MeterType:       measures.TLG,
				},
				resultLastErr: nil,
				result:        errors.New("err save"),
				result1Save:   errors.New("err save"),
			},
		},
		"should return err if client.ListProcessedLoadCurve return err": {
			input: measures.ProcessMeasurePayload{
				MeterConfig: measures.MeterConfig{
					Type: measures.TLG,
					ServicePoint: measures.ServicePoint{
						Cups: "1",
					},
				},
				Date: time.Now(),
			},
			want: want{
				Status: billing_measures.Calculating,
				resultLast: billing_measures.BillingMeasure{
					DescriptionStatus: &noLastBillingMeasure,
					DistributorCode:   "0130",
					DistributorID:     "id",
					CUPS:              "CUPS",
					EndDate:           time.Now().AddDate(0, -1, 0),
					InitDate:          time.Now().AddDate(0, -2, 0),
					MeterType:         measures.TLG,
				},
				resultLastErr:                nil,
				result:                       errors.New("err client.ListProcessedLoadCurve"),
				result1Save:                  nil,
				ListProcessedLoadCurveResult: []process_measures.ProcessedLoadCurve{},
				ListProcessedLoadCurveErr:    errors.New("err client.ListProcessedLoadCurve"),
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			b := billing_measures.NewBillingMeasure(
				testCase.input.MeterConfig.Cups(),
				testCase.want.resultLast.EndDate,
				testCase.input.Date.AddDate(0, 0, 1),
				testCase.want.resultLast.DistributorCode,
				testCase.want.resultLast.DistributorID,
				[]measures.PeriodKey{measures.P1, measures.P2, measures.P3},
				[]measures.Magnitude{},
				testCase.input.MeterConfig.Type,
			)
			calendarPeriodRepository := new(mocks.CalendarPeriodRepository)
			masterTablesClient := new(clients_mocks.MasterTables)

			loc, _ := time.LoadLocation("Europe/Madrid")
			b.PointType = testCase.input.MeterConfig.PointType()
			b.RegisterType = testCase.input.MeterConfig.CurveType
			repo := new(mocks.BillingMeasureRepository)
			repoProfiles := new(mocks.ConsumProfileRepository)
			repo.On("Last", mock.Anything, billing_measures.QueryLast{
				CUPS: testCase.input.MeterConfig.Cups(),
				Date: testCase.input.Date,
			}).Return(testCase.want.resultLast, testCase.want.resultLastErr)
			repo.On("GetPrevious", mock.Anything, billing_measures.GetPrevious{
				CUPS:     b.CUPS,
				InitDate: b.InitDate,
				EndDate:  b.EndDate,
			}).Return(testCase.want.GetPrevious, testCase.want.GetPreviousErr)

			repo.On("Save", mock.Anything, mock.MatchedBy(func(compare billing_measures.BillingMeasure) bool {
				b.GenerationDate = compare.GenerationDate
				if testCase.want.DescriptionStatus != "" {
					b.DescriptionStatus = &testCase.want.DescriptionStatus
				}
				b.Status = testCase.want.Status
				return assert.Equal(t, b, compare)
			})).Return(testCase.want.result1Save)

			processedMeasuresRepository := new(mocks.ProcessedMeasureRepository)
			processedMeasuresRepository.On("ProcessedLoadCurveByCups", mock.Anything, process_measures.QueryProcessedLoadCurveByCups{
				CUPS:      testCase.input.MeterConfig.Cups(),
				StartDate: testCase.want.resultLast.EndDate,
				EndDate:   testCase.input.Date.AddDate(0, 0, 1),
				CurveType: measures.HourlyMeasureCurveReadingType,
				Status:    measures.Valid,
			}).Return(testCase.want.ListProcessedLoadCurveResult, testCase.want.ListProcessedLoadCurveErr)
			processedMeasuresRepository.On("GetProcessedDailyClosureByCup", mock.Anything, process_measures.QueryClosedCupsMeasureOnDate{
				CUPS: testCase.input.MeterConfig.Cups(),
				Date: testCase.input.Date.AddDate(0, 0, 1),
			}).Return(testCase.want.GetDailyReadingClosureByCupsResult, testCase.want.GetDailyReadingClosureByCupsErr)
			masterTablesClient.On("GetTariff", mock.Anything, mock.Anything).Return(
				clients.Tariffs{}, nil)

			inventoryClient := new(clients_mocks.Inventory)
			calendarPeriodRepository.On("GetCalendarPeriod", mock.Anything, mock.Anything).Return(fixtures.CALENDAR_PERIOD, nil)

			p := new(event_mocks.Publisher)

			publisherCreator := func(ctx context.Context) (event.Publisher, error) {
				return p, nil
			}

			svc := NewProcessMvhByCup(repo, processedMeasuresRepository, inventoryClient, repoProfiles, calendarPeriodRepository, loc, masterTablesClient, publisherCreator, "topic")
			assert.Equal(t, testCase.want.result, svc.Handler(ctx, testCase.input), testName)
		})
	}
}
func Test_Unit_Services_BillingMeasures_Process_mvh_by_cup_fillEmptyCurves(t *testing.T) {

	type want struct {
		result                         []process_measures.ProcessedLoadCurve
		resultErr                      error
		ProcessedLoadCurveByCupsResult []process_measures.ProcessedLoadCurve
		ProcessedLoadCurveByCupsErr    error
	}

	type input struct {
		cups               string
		startDate, endDate time.Time
		processedLoadCurve []process_measures.ProcessedLoadCurve
	}

	tests := map[string]struct {
		input input
		want  want
	}{
		"should not change ProcessedLoadCurve if not empty": {
			input: input{
				cups:      "cups",
				startDate: time.Date(2022, 01, 01, 1, 0, 0, 0, time.UTC),
				endDate:   time.Date(2022, 01, 01, 2, 0, 0, 0, time.UTC),
				processedLoadCurve: []process_measures.ProcessedLoadCurve{
					{
						EndDate: time.Date(2022, 01, 01, 1, 0, 0, 0, time.UTC),
					},
					{
						EndDate: time.Date(2022, 01, 01, 2, 0, 0, 0, time.UTC),
					},
				},
			},
			want: want{
				resultErr: nil,
				result: []process_measures.ProcessedLoadCurve{
					{
						EndDate: time.Date(2022, 01, 01, 1, 0, 0, 0, time.UTC),
					},
					{
						EndDate: time.Date(2022, 01, 01, 2, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		"should return err if fetch ProcessedLoadCurveByCups return err": {
			input: input{
				cups:      "cups",
				startDate: time.Date(2022, 01, 01, 1, 0, 0, 0, time.UTC),
				endDate:   time.Date(2022, 01, 01, 2, 0, 0, 0, time.UTC),
				processedLoadCurve: []process_measures.ProcessedLoadCurve{
					{
						EndDate: time.Date(2022, 01, 01, 1, 0, 0, 0, time.UTC),
					},
					{
						EndDate: time.Date(2022, 01, 01, 2, 0, 0, 0, time.UTC),
						Origin:  measures.Filled,
					},
				},
			},
			want: want{
				resultErr:                   errors.New("err"),
				ProcessedLoadCurveByCupsErr: errors.New("err"),
			},
		},
		"should not fill processedLoadCurve hourly with processedLoadCurve quarter": {
			input: input{
				cups:      "cups",
				startDate: time.Date(2022, 01, 01, 1, 0, 0, 0, time.UTC),
				endDate:   time.Date(2022, 01, 01, 2, 0, 0, 0, time.UTC),
				processedLoadCurve: []process_measures.ProcessedLoadCurve{
					{
						EndDate: time.Date(2022, 01, 01, 1, 0, 0, 0, time.UTC),
					},
					{
						EndDate: time.Date(2022, 01, 01, 2, 0, 0, 0, time.UTC),
						Origin:  measures.Filled,
					},
				},
			},
			want: want{
				result: []process_measures.ProcessedLoadCurve{
					{
						EndDate: time.Date(2022, 01, 01, 1, 0, 0, 0, time.UTC),
					},
					{
						EndDate: time.Date(2022, 01, 01, 2, 0, 0, 0, time.UTC),
						Origin:  measures.Filled,
					},
				},
			},
		},
		"should fill processedLoadCurve hourly with processedLoadCurve quarter": {
			input: input{
				cups:      "cups",
				startDate: time.Date(2022, 01, 01, 1, 0, 0, 0, time.UTC),
				endDate:   time.Date(2022, 01, 01, 2, 0, 0, 0, time.UTC),
				processedLoadCurve: []process_measures.ProcessedLoadCurve{
					{
						EndDate: time.Date(2022, 01, 01, 1, 0, 0, 0, time.UTC),
					},
					{
						EndDate: time.Date(2022, 01, 01, 2, 0, 0, 0, time.UTC),
						Origin:  measures.Filled,
					},
				},
			},
			want: want{
				result: []process_measures.ProcessedLoadCurve{
					{
						EndDate: time.Date(2022, 01, 01, 1, 0, 0, 0, time.UTC),
					},
					{
						EndDate: time.Date(2022, 01, 01, 2, 0, 0, 0, time.UTC),
						AI:      10,
						AE:      20,
						R1:      30,
						R2:      40,
						R3:      50,
						R4:      60,
						Origin:  measures.CalculatedWithQuarter,
					},
				},
				ProcessedLoadCurveByCupsResult: []process_measures.ProcessedLoadCurve{
					{
						EndDate: time.Date(2022, 01, 01, 1, 15, 0, 0, time.UTC),
						AI:      2.5,
						AE:      5,
						R1:      7,
						R2:      10,
						R3:      12.5,
						R4:      15,
						Origin:  measures.CalculatedWithQuarter,
					},
					{
						EndDate: time.Date(2022, 01, 01, 1, 30, 0, 0, time.UTC),
						AI:      2.5,
						AE:      5,
						R1:      7,
						R2:      10,
						R3:      12.5,
						R4:      15,
						Origin:  measures.CalculatedWithQuarter,
					},
					{
						EndDate: time.Date(2022, 01, 01, 1, 45, 0, 0, time.UTC),
						AI:      2.5,
						AE:      5,
						R1:      7,
						R2:      10,
						R3:      12.5,
						R4:      15,
						Origin:  measures.CalculatedWithQuarter,
					},
					{
						EndDate: time.Date(2022, 01, 01, 2, 0, 0, 0, time.UTC),
						AI:      2.5,
						AE:      5,
						R1:      9,
						R2:      10,
						R3:      12.5,
						R4:      15,
						Origin:  measures.CalculatedWithQuarter,
					},
				},
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {

			loc, _ := time.LoadLocation("Europe/Madrid")

			processedMeasuresRepository := new(mocks.ProcessedMeasureRepository)
			processedMeasuresRepository.On("ProcessedLoadCurveByCups", mock.Anything, process_measures.QueryProcessedLoadCurveByCups{
				CUPS:      testCase.input.cups,
				StartDate: testCase.input.startDate,
				EndDate:   testCase.input.endDate,
				CurveType: measures.QuarterMeasureCurveReadingType,
				Status:    measures.Valid,
			}).Return(testCase.want.ProcessedLoadCurveByCupsResult, testCase.want.ProcessedLoadCurveByCupsErr)

			svc := NewProcessMvhByCup(
				nil,
				processedMeasuresRepository,
				nil,
				nil,
				nil,
				loc,
				nil,
				nil,
				"",
			)

			result, err := svc.fillEmptyCurves(
				context.Background(),
				testCase.input.cups,
				testCase.input.startDate,
				testCase.input.endDate,
				testCase.input.processedLoadCurve,
			)

			assert.ElementsMatch(t, testCase.want.result, result)
			assert.Equal(t, testCase.want.resultErr, err)
		})
	}
}
