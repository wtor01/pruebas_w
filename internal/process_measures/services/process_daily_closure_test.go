package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	clientMocks "bitbucket.org/sercide/data-ingestion/internal/common/clients/clients_mocks"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services/fixtures"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestProcessDailyClosure(t *testing.T) {

	generationDate := func() time.Time {
		return time.Time{}
	}

	type input struct {
		dto               measures.ProcessMeasurePayload
		ctx               context.Context
		processedMeasures process_measures.ProcessedDailyClosure
	}

	type output struct {
		listMeasureClose       []gross_measures.MeasureCloseWrite
		listMeasureCloseErr    error
		saveAllDailyClosureErr error
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
				ctx: context.Background(),
			},
			output: output{
				listMeasureCloseErr:    errors.New("err"),
				listMeasureClose:       []gross_measures.MeasureCloseWrite{},
				saveAllDailyClosureErr: nil,
			},
			wants: wants{
				err: errors.New("err"),
			},
		},
		"should not return err if list gross measures is empty": {
			input: input{
				dto: measures.ProcessMeasurePayload{},
				ctx: context.Background(),
			},
			output: output{
				listMeasureCloseErr:    nil,
				listMeasureClose:       []gross_measures.MeasureCloseWrite{},
				saveAllDailyClosureErr: nil,
			},
			wants: wants{
				err: nil,
			},
		},
		"should return err if is invalid measure": {
			input: input{
				dto: measures.ProcessMeasurePayload{},
				ctx: context.Background(),
			},
			output: output{
				listMeasureClose:       []gross_measures.MeasureCloseWrite{fixtures.MEASURES_CLOSE_DAILY[1]},
				listMeasureCloseErr:    nil,
				saveAllDailyClosureErr: nil,
			},
			wants: wants{err: errors.New("invalid measure")},
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
				ctx:               context.Background(),
				processedMeasures: fixtures.RESULT_PROCESSED_DAILY_CLOSURE,
			},
			output: output{
				listMeasureClose:       fixtures.MEASURES_CLOSE_DAILY,
				listMeasureCloseErr:    nil,
				saveAllDailyClosureErr: nil,
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
				ctx:               context.Background(),
				processedMeasures: fixtures.RESULT_PROCESSED_DAILY_CLOSURE,
			},
			output: output{
				listMeasureClose:       fixtures.MEASURES_CLOSE_DAILY,
				listMeasureCloseErr:    nil,
				saveAllDailyClosureErr: errors.New("err"),
			},
			wants: wants{err: errors.New("err")},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			grossRepository := new(mocks.GrossMeasureRepository)
			processMeasureRepository := new(mocks.ProcessedMeasureRepository)
			calendarPeriodRepository := new(mocks.CalendarPeriodRepository)
			validationMongo := new(mocks.ValidationMongoRepository)
			loc, _ := time.LoadLocation("Europe/Madrid")
			masterTablesClient := new(clientMocks.MasterTables)

			masterTablesClient.On("GetTariff", mock.Anything, mock.Anything).Return(
				clients.Tariffs{}, nil)
			validation := new(clientMocks.Validation)

			grossRepository.On("ListDailyCloseMeasures", test.input.ctx, gross_measures.QueryListForProcessClose{
				ReadingType:  measures.DailyClosure,
				SerialNumber: test.input.dto.MeterConfig.SerialNumber(),
				Date:         test.input.dto.Date,
			}).Return(test.output.listMeasureClose, test.output.listMeasureCloseErr)
			validation.On("GetValidationConfigList", mock.Anything, test.input.dto.MeterConfig.DistributorID, string(measures.Process)).Return(
				make([]clients.ValidationConfig, 0), nil)
			calendarPeriodRepository.On("GetCalendarPeriod", mock.Anything, mock.Anything).Return(fixtures.CALENDAR_PERIOD, nil)

			test.input.processedMeasures.GenerateID()

			processMeasureRepository.On("SaveDailyClosure", test.input.ctx, test.input.processedMeasures).Return(test.output.saveAllDailyClosureErr)
			processMeasureRepository.On(
				"GetProcessedDailyClosureByCup",
				test.input.ctx,
				process_measures.QueryClosedCupsMeasureOnDate{CUPS: test.input.dto.MeterConfig.Cups(), Date: test.input.dto.Date.AddDate(0, 0, 1)},
			).Return(
				process_measures.ProcessedDailyClosure{},
				nil,
			)

			srv := NewProcessDailyClosure(
				grossRepository,
				processMeasureRepository,
				validation,
				calendarPeriodRepository,
				loc,
				masterTablesClient,
				validationMongo,
			)

			processMeasureRepository.On(
				"GetProcessedDailyClosureByCup",
				test.input.ctx,
				process_measures.QueryClosedCupsMeasureOnDate{CUPS: test.input.dto.MeterConfig.Cups(), Date: test.input.dto.Date.AddDate(0, 0, 1)},
			).Return(
				process_measures.ProcessedDailyClosure{},
				nil,
			)

			srv.generatorDate = generationDate
			err := srv.Handle(test.input.ctx, test.input.dto)

			assert.Equal(t, test.wants.err, err)
		},
		)
	}
}
