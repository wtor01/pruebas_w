package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	clientMocks "bitbucket.org/sercide/data-ingestion/internal/common/clients/clients_mocks"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services/fixtures"
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	event_mocks "bitbucket.org/sercide/data-ingestion/pkg/event/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestProcessMonthlyHandle(t *testing.T) {
	generationDate := func() time.Time {
		return time.Time{}
	}

	type input struct {
		dto               measures.ProcessMeasurePayload
		processedMeasures []process_measures.ProcessedMonthlyClosure
	}

	type output struct {
		listMeasureClose         []gross_measures.MeasureCloseWrite
		listMeasureCloseErr      error
		saveAllMonthlyClosureErr error
		publisherErr             error
		publishErr               error
	}

	type wants struct {
		err error
	}

	testCases := map[string]struct {
		input  input
		output output
		wants  wants
	}{
		"should return err if listMeasure fail": {
			input: input{
				dto: measures.ProcessMeasurePayload{},
			},
			output: output{
				listMeasureCloseErr:      errors.New("invalid measures"),
				listMeasureClose:         []gross_measures.MeasureCloseWrite{},
				saveAllMonthlyClosureErr: nil,
			},
			wants: wants{
				err: errors.New("invalid measures"),
			},
		},
		"should transform ok": {
			input: input{
				dto: measures.ProcessMeasurePayload{
					MeterConfig: measures.MeterConfig{
						DistributorID:   "DistributorX",
						DistributorCode: "0130",
						Meter: measures.Meter{
							SerialNumber: "123456",
						},
						ServicePoint: measures.ServicePoint{
							Cups: "XXXXXXXXXXXX",
						},
						MeasurePoint: measures.MeasurePoint{
							Type: measures.MeasurePointTypeR,
						},
						ContractualSituations: measures.ContractualSituations{},
					},
				},
				processedMeasures: []process_measures.ProcessedMonthlyClosure{fixtures.RESULT_PROCESSED_MONTHLY_CLOSURE},
			},
			output: output{
				listMeasureClose:         fixtures.MEASURES_CLOSE_MONTHLY,
				listMeasureCloseErr:      nil,
				saveAllMonthlyClosureErr: nil,
			},
			wants: wants{
				err: nil,
			},
		},
		"should err if save fail": {
			input: input{
				dto: measures.ProcessMeasurePayload{
					MeterConfig: measures.MeterConfig{
						DistributorID:   "DistributorX",
						DistributorCode: "0130",
						Meter: measures.Meter{
							SerialNumber: "123456",
						},
						ServicePoint: measures.ServicePoint{
							Cups: "XXXXXXXXXXXX",
						},
						MeasurePoint: measures.MeasurePoint{
							Type: measures.MeasurePointTypeR,
						},
						ContractualSituations: measures.ContractualSituations{},
					},
				},
				processedMeasures: []process_measures.ProcessedMonthlyClosure{fixtures.RESULT_PROCESSED_MONTHLY_CLOSURE},
			},
			output: output{
				listMeasureClose:         fixtures.MEASURES_CLOSE_MONTHLY,
				listMeasureCloseErr:      nil,
				saveAllMonthlyClosureErr: errors.New("err"),
			},
			wants: wants{err: errors.New("err")},
		},
		"should publish ok": {
			input: input{
				dto: measures.ProcessMeasurePayload{
					MeterConfig: measures.MeterConfig{
						DistributorID:   "DistributorX",
						DistributorCode: "0130",
						Meter: measures.Meter{
							SerialNumber: "123456",
						},
						ServicePoint: measures.ServicePoint{
							Cups: "XXXXXXXXXXXX",
						},
						MeasurePoint: measures.MeasurePoint{
							Type: measures.MeasurePointTypeR,
						},
						ContractualSituations: measures.ContractualSituations{},
					},
				},
				processedMeasures: []process_measures.ProcessedMonthlyClosure{fixtures.RESULT_PROCESSED_MONTHLY_CLOSURE},
			},
			output: output{
				listMeasureClose:         fixtures.MEASURES_CLOSE_MONTHLY,
				listMeasureCloseErr:      nil,
				saveAllMonthlyClosureErr: nil,
				publisherErr:             nil,
				publishErr:               nil,
			},
			wants: wants{
				err: nil,
			},
		},
		"should publish ok without incremental measures": {
			input: input{
				dto: measures.ProcessMeasurePayload{
					MeterConfig: measures.MeterConfig{
						DistributorID:   "DistributorX",
						DistributorCode: "0130",
						Meter: measures.Meter{
							SerialNumber: "123456",
						},
						ServicePoint: measures.ServicePoint{
							Cups: "XXXXXXXXXXXX",
						},
						MeasurePoint: measures.MeasurePoint{
							Type: measures.MeasurePointTypeR,
						},
						ContractualSituations: measures.ContractualSituations{},
					},
				},
				processedMeasures: []process_measures.ProcessedMonthlyClosure{fixtures.RESULT_PROCESSED_MONTHLY_CLOSURE_WITHOUT_INCREMENT},
			},
			output: output{
				listMeasureClose:         []gross_measures.MeasureCloseWrite{fixtures.MEASURES_CLOSE_MONTHLY[0]},
				listMeasureCloseErr:      nil,
				saveAllMonthlyClosureErr: nil,
				publisherErr:             nil,
				publishErr:               nil,
			},
			wants: wants{
				err: nil,
			},
		},
		"should publish ok with only incremental measures": {
			input: input{
				dto: measures.ProcessMeasurePayload{
					MeterConfig: measures.MeterConfig{
						DistributorID:   "DistributorX",
						DistributorCode: "0130",
						Meter: measures.Meter{
							SerialNumber: "123456",
						},
						ServicePoint: measures.ServicePoint{
							Cups: "XXXXXXXXXXXX",
						},
						MeasurePoint: measures.MeasurePoint{
							Type: measures.MeasurePointTypeR,
						},
						ContractualSituations: measures.ContractualSituations{},
					},
				},
				processedMeasures: []process_measures.ProcessedMonthlyClosure{fixtures.RESULT_PROCESSED_MONTHLY_CLOSURE_WITHOUT_ABS},
			},
			output: output{
				listMeasureClose:         []gross_measures.MeasureCloseWrite{fixtures.MEASURES_CLOSE_MONTHLY[2]},
				listMeasureCloseErr:      nil,
				saveAllMonthlyClosureErr: nil,
				publisherErr:             nil,
				publishErr:               nil,
			},
			wants: wants{
				err: nil,
			},
		},
		"should save 2 close with different hour": {
			input: input{
				dto: measures.ProcessMeasurePayload{
					MeterConfig: measures.MeterConfig{
						DistributorID:   "DistributorX",
						DistributorCode: "0130",
						Meter: measures.Meter{
							SerialNumber: "123456",
						},
						ServicePoint: measures.ServicePoint{
							Cups: "XXXXXXXXXXXX",
						},
						MeasurePoint: measures.MeasurePoint{
							Type: measures.MeasurePointTypeR,
						},
						ContractualSituations: measures.ContractualSituations{},
					},
				},
				processedMeasures: []process_measures.ProcessedMonthlyClosure{
					fixtures.RESULT_PROCESSED_MONTHLY_CLOSURE,
					fixtures.RESULT_PROCESSED_2_MONTHLY_CLOSURE,
				},
			},
			output: output{
				listMeasureClose:         fixtures.MEASURES_CLOSE_MONTHLY_2_HOURS,
				listMeasureCloseErr:      nil,
				saveAllMonthlyClosureErr: nil,
				publisherErr:             nil,
				publishErr:               nil,
			},
			wants: wants{
				err: nil,
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			grossRepository := new(mocks.GrossMeasureRepository)
			validation := new(clientMocks.Validation)
			validationMongo := new(mocks.ValidationMongoRepository)
			cal := new(mocks.RepositoryCalendar)
			processMeasureRepository := new(mocks.ProcessedMeasureRepository)
			calendarPeriodRepository := new(mocks.CalendarPeriodRepository)
			loc, _ := time.LoadLocation("Europe/Madrid")
			masterTablesClient := new(clientMocks.MasterTables)
			pub := new(event_mocks.Publisher)

			for i, write := range test.output.listMeasureClose {
				write.GenerateID()
				test.output.listMeasureClose[i] = write
			}

			pub.Mock.On("Publish", mock.Anything, "Test", mock.Anything, mock.Anything).Return(test.output.publishErr)
			publisherCreator := func(ctx context.Context) (event.Publisher, error) {
				return pub, test.output.publisherErr
			}

			masterTablesClient.On("GetTariff", mock.Anything, mock.Anything).Return(
				clients.Tariffs{}, nil)
			grossRepository.On("ListDailyCloseMeasures", mock.Anything, gross_measures.QueryListForProcessClose{
				ReadingType:  measures.BillingClosure,
				SerialNumber: test.input.dto.MeterConfig.SerialNumber(),
				Date:         test.input.dto.Date,
			}).Return(test.output.listMeasureClose, test.output.listMeasureCloseErr)

			for _, measure := range test.input.processedMeasures {
				measure.GenerateID()
				measure.SetLastDailyClose(process_measures.ProcessedMonthlyClosure{})
				measure.Periods = []measures.PeriodKey{
					measures.P1,
					measures.P2,
					measures.P3,
				}
				processMeasureRepository.On("SaveMonthlyClosure", mock.Anything, measure).Return(test.output.saveAllMonthlyClosureErr).Once()
			}

			processMeasureRepository.On(
				"GetMonthlyClosureByCup",
				mock.Anything,
				process_measures.QueryClosedCupsMeasureOnDate{CUPS: test.input.dto.MeterConfig.Cups(), Date: test.input.dto.Date.AddDate(0, -1, 1)},
			).Return(
				process_measures.ProcessedMonthlyClosure{},
				nil,
			)
			validation.On("GetValidationConfigList", mock.Anything, test.input.dto.MeterConfig.DistributorID, string(measures.Process)).Return(
				make([]clients.ValidationConfig, 0), nil)
			cal.On("GetPeriodsByDate", mock.Anything, test.input.dto.MeterConfig.CalendarID, test.input.dto.Date).Return(make([]calendar.PeriodCalendar, 0), nil)
			calendarPeriodRepository.On("GetCalendarPeriod", mock.Anything, mock.Anything).Return(fixtures.CALENDAR_PERIOD, nil)

			srv := NewProcessMonthlyClosure(
				grossRepository,
				processMeasureRepository,
				cal,
				validation,
				calendarPeriodRepository,
				loc,
				masterTablesClient,
				validationMongo,
				publisherCreator,
				"Test",
			)
			srv.generatorDate = generationDate
			err := srv.Handle(context.Background(), test.input.dto)

			assert.Equal(t, test.wants.err, err)
		})
	}
}
