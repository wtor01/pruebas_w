package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	clients_mocks "bitbucket.org/sercide/data-ingestion/internal/common/clients/clients_mocks"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Unit_Services_Closure_Create(t *testing.T) {

	type input struct {
		ctx     context.Context
		monthly CreateClosureDto
	}

	type want struct {
		err error
	}

	type results struct {
		resultMeterConfig      measures.MeterConfig
		resultMeterConfigErr   error
		resulTariff            clients.Tariffs
		resulTariffErr         error
		resultCreateClosureErr error
	}

	tests := map[string]struct {
		input   input
		want    want
		results results
	}{
		"should be error in meter config": {
			input: input{
				ctx:     context.Background(),
				monthly: CreateClosureDto{},
			},
			results: results{
				resultMeterConfig:      measures.MeterConfig{},
				resultMeterConfigErr:   errors.New("error"),
				resulTariff:            clients.Tariffs{},
				resulTariffErr:         nil,
				resultCreateClosureErr: nil,
			},
			want: want{
				err: errors.New("error"),
			},
		},
		"should be error in get tariff": {
			input: input{
				ctx:     context.Background(),
				monthly: CreateClosureDto{},
			},
			results: results{
				resultMeterConfig:      measures.MeterConfig{},
				resultMeterConfigErr:   nil,
				resulTariff:            clients.Tariffs{},
				resulTariffErr:         errors.New("error"),
				resultCreateClosureErr: nil,
			},
			want: want{
				err: errors.New("error"),
			},
		},
		"should be error": {
			input: input{
				ctx:     context.Background(),
				monthly: CreateClosureDto{},
			},
			results: results{
				resultMeterConfig:      measures.MeterConfig{},
				resultMeterConfigErr:   nil,
				resulTariff:            clients.Tariffs{},
				resulTariffErr:         nil,
				resultCreateClosureErr: errors.New("error"),
			},
			want: want{
				err: errors.New("error"),
			},
		},
		"should correct": {
			input: input{
				ctx:     context.Background(),
				monthly: CreateClosureDto{},
			},
			results: results{
				resultMeterConfig:      measures.MeterConfig{},
				resultMeterConfigErr:   nil,
				resulTariff:            clients.Tariffs{},
				resulTariffErr:         nil,
				resultCreateClosureErr: nil,
			},
			want: want{},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			repoProcessMeasure := new(mocks.ProcessMeasureClosureRepository)
			repoInventory := new(mocks.InventoryRepository)
			repoMasterTables := new(clients_mocks.MasterTables)

			repoInventory.Mock.On("GetMeterConfigByCupsAPI", testCase.input.ctx, measures.GetMeterConfigByCupsQuery{
				CUPS:        testCase.input.monthly.Monthly.CUPS,
				Time:        testCase.input.monthly.Monthly.EndDate,
				Distributor: testCase.input.monthly.Monthly.DistributorID,
			}).Return(testCase.results.resultMeterConfig, testCase.results.resultMeterConfigErr)

			repoMasterTables.Mock.On("GetTariff", testCase.input.ctx, clients.GetTariffDto{
				ID: testCase.results.resultMeterConfig.TariffID(),
			}).Return(testCase.results.resulTariff, testCase.results.resulTariffErr)

			processMeasureMonthlyClosure := process_measures.NewProcessedMonthlyClosure(testCase.results.resultMeterConfig, time.Time{})
			processMeasureMonthlyClosure.ContractNumber = string(testCase.results.resultMeterConfig.PriorityContract)
			processMeasureMonthlyClosure.StartDate = testCase.input.monthly.Monthly.StartDate
			processMeasureMonthlyClosure.EndDate = testCase.input.monthly.Monthly.EndDate
			processMeasureMonthlyClosure.GenerationDate = time.Date(2022, time.September, 9, 13, 57, 9, 685583200, time.Local)
			processMeasureMonthlyClosure.ReadingDate = time.Date(2022, time.September, 9, 13, 57, 9, 685583200, time.Local)
			processMeasureMonthlyClosure.Origin = "MANUAL"
			processMeasureMonthlyClosure.CalendarPeriods = testCase.input.monthly.Monthly.CalendarPeriods
			processMeasureMonthlyClosure.Coefficient = testCase.results.resulTariff.Coef
			processMeasureMonthlyClosure.Id = testCase.input.monthly.Monthly.Id

			repoProcessMeasure.Mock.On("CreateClosure", testCase.input.ctx,
				*processMeasureMonthlyClosure,
			).Return(testCase.results.resultCreateClosureErr)

			srv := NewCreateClosure(repoProcessMeasure, repoInventory, repoMasterTables, func() time.Time { return time.Date(2022, time.September, 9, 13, 57, 9, 685583200, time.Local) })

			err := srv.Handler(testCase.input.ctx, testCase.input.monthly)

			assert.Equal(t, testCase.want.err, err, testCase)

		})
	}

}
