package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	clientMocks "bitbucket.org/sercide/data-ingestion/internal/common/clients/clients_mocks"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
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

func TestProcessCurve(t *testing.T) {

	generationDate := func() time.Time {
		return time.Date(2022, 01, 1, 0, 0, 0, 0, time.UTC)
	}

	type input struct {
		dto                 measures.ProcessMeasurePayload
		ctx                 context.Context
		processedCurve      []process_measures.ProcessedLoadCurve
		processedDailyCurve process_measures.ProcessedDailyLoadCurve
	}

	type output struct {
		listMeasureCurve       []gross_measures.MeasureCurveWrite
		listMeasureCurveErr    error
		saveAllCurveClosureErr error
		saveDailyLoadCurveErr  error
	}

	type wants struct {
		err error
	}

	var testCases = map[string]struct {
		input  input
		output output
		wants  wants
	}{
		"should return err if listMeasure fail": {
			input: input{
				dto:                 measures.ProcessMeasurePayload{},
				ctx:                 context.Background(),
				processedCurve:      []process_measures.ProcessedLoadCurve{},
				processedDailyCurve: process_measures.ProcessedDailyLoadCurve{},
			},
			output: output{
				listMeasureCurveErr:    errors.New("err"),
				listMeasureCurve:       []gross_measures.MeasureCurveWrite{},
				saveAllCurveClosureErr: nil,
				saveDailyLoadCurveErr:  nil,
			},
			wants: wants{
				err: errors.New("err"),
			},
		},
		"should err if save load curve fail": {
			input: input{
				dto:                 fixtures.DTO_PROCESSED_CURVE,
				ctx:                 context.Background(),
				processedCurve:      fixtures.RESULT_PROCESSED_LOAD_CURVE,
				processedDailyCurve: fixtures.RESULT_PROCESSED_DAILY_LOAD_CURVE,
			},
			output: output{
				listMeasureCurve:       fixtures.MEASURES_CURVE,
				listMeasureCurveErr:    nil,
				saveAllCurveClosureErr: errors.New("err"),
				saveDailyLoadCurveErr:  nil,
			},
			wants: wants{err: errors.New("err")},
		},
		"should err if save daily load curve fail": {
			input: input{
				dto:                 fixtures.DTO_PROCESSED_CURVE,
				ctx:                 context.Background(),
				processedCurve:      fixtures.RESULT_PROCESSED_LOAD_CURVE,
				processedDailyCurve: fixtures.RESULT_PROCESSED_DAILY_LOAD_CURVE,
			},
			output: output{
				listMeasureCurve:       fixtures.MEASURES_CURVE,
				listMeasureCurveErr:    nil,
				saveAllCurveClosureErr: nil,
				saveDailyLoadCurveErr:  errors.New("err"),
			},
			wants: wants{err: errors.New("err")},
		},
		"should transform ok": {
			input: input{
				dto:                 fixtures.DTO_PROCESSED_CURVE,
				ctx:                 context.Background(),
				processedCurve:      fixtures.RESULT_PROCESSED_LOAD_CURVE,
				processedDailyCurve: fixtures.RESULT_PROCESSED_DAILY_LOAD_CURVE,
			},
			output: output{
				listMeasureCurve:       fixtures.MEASURES_CURVE,
				listMeasureCurveErr:    nil,
				saveAllCurveClosureErr: nil,
			},
			wants: wants{
				err: nil,
			},
		},
		"should fill curve quarter ok": {
			input: input{
				dto:                 fixtures.DTO_PROCESSED_CURVE_QUARTER,
				ctx:                 context.Background(),
				processedCurve:      fixtures.RESULT_PROCESSED_LOAD_CURVE_QUARTER_FILLED,
				processedDailyCurve: fixtures.RESULT_PROCESSED_DAILY_LOAD_CURVE_QUARTERLY,
			},
			output: output{
				listMeasureCurve:       []gross_measures.MeasureCurveWrite{},
				listMeasureCurveErr:    nil,
				saveAllCurveClosureErr: nil,
			},
			wants: wants{
				err: nil,
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			grossRepository := new(mocks.GrossMeasureRepository)

			processMeasureRepository := new(mocks.ProcessedMeasureRepository)
			calendarPeriodRepository := new(mocks.CalendarPeriodRepository)
			validationMongo := new(mocks.ValidationMongoRepository)
			seasonsRepository := new(mocks.RepositorySeasons)
			validation := new(clientMocks.Validation)
			masterTablesClient := new(clientMocks.MasterTables)

			grossRepository.On("ListDailyCurveMeasures", test.input.ctx, gross_measures.QueryListForProcessCurve{
				SerialNumber: test.input.dto.MeterConfig.SerialNumber(),
				Date:         test.input.dto.Date,
				CurveType:    test.input.dto.CurveType,
			}).Return(test.output.listMeasureCurve, test.output.listMeasureCurveErr)

			validation.On("GetValidationConfigList", mock.Anything, test.input.dto.MeterConfig.DistributorID, string(measures.Process)).Return(
				make([]clients.ValidationConfig, 0), nil)

			masterTablesClient.On("GetTariff", mock.Anything, mock.Anything).Return(
				clients.Tariffs{}, nil)

			calendarPeriodRepository.On("GetCalendarPeriod", mock.Anything, mock.Anything).Return(fixtures.CALENDAR_PERIOD, nil)

			for i := range test.input.processedCurve {
				test.input.processedCurve[i].GenerateID()
			}

			processMeasureRepository.On("SaveAllProcessedLoadCurve", test.input.ctx, mock.MatchedBy(func(s []process_measures.ProcessedLoadCurve) bool {
				return assert.ElementsMatch(t, test.input.processedCurve, s)
			})).Return(test.output.saveAllCurveClosureErr)

			test.input.processedDailyCurve.GenerateID()

			processMeasureRepository.On("SaveProcessedDailyLoadCurve", test.input.ctx, test.input.processedDailyCurve).Return(test.output.saveDailyLoadCurveErr)
			seasonsRepository.On("GetDayTypeByMonth", test.input.ctx, int(test.input.dto.Date.Month()), fixtures.CALENDAR_PERIOD.IsFestiveDay()).Return(seasons.DayTypes{IsFestive: fixtures.CALENDAR_PERIOD.IsFestiveDay(), Month: int(test.input.dto.Date.Month())}, nil)
			loc, _ := time.LoadLocation("Europe/Madrid")
			srv := NewProcessCurve(
				grossRepository,
				processMeasureRepository,
				calendarPeriodRepository,
				seasonsRepository,
				loc,
				validation,
				masterTablesClient,
				validationMongo,
			)

			srv.generatorDate = generationDate
			err := srv.Handle(test.input.ctx, test.input.dto)

			assert.Equal(t, test.wants.err, err)
		})
	}
}
