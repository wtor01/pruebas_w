package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/clients"
	clients_mocks "bitbucket.org/sercide/data-ingestion/internal/common/clients/clients_mocks"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Unit_Service_GrossMeasures_GetDashboardSupplyPoint_Handler(t *testing.T) {
	type input struct {
		dto DashboardMeasureSupplyPointDto
	}

	type output struct {
		MeterConfig       measures.MeterConfig
		MeterConfigErr    error
		Tariff            clients.Tariffs
		TariffErr         error
		CalendarPeriod    measures.CalendarPeriod
		CalendarPeriodErr error
		CurveMeasures     []gross_measures.MeasureCurveWrite
		CurveMeasuresErr  error
		CloseMeasures     []gross_measures.MeasureCloseWrite
		CloseMeasuresErr  error
	}

	type want struct {
		result gross_measures.DashboardMeasureSupplyPoint
		err    error
	}

	dto := NewDashboardMeasureSupplyPointDto(
		"CUPS",
		"DistributorX",
		time.Date(2022, 6, 30, 22, 0, 0, 0, time.UTC),
		time.Date(2022, 7, 31, 22, 0, 0, 0, time.UTC),
	)

	festiveHours := measures.CalendarPeriodHours{
		Hour1:  "P3",
		Hour2:  "P3",
		Hour3:  "P3",
		Hour4:  "P3",
		Hour5:  "P3",
		Hour6:  "P3",
		Hour7:  "P3",
		Hour8:  "P3",
		Hour9:  "P3",
		Hour10: "P3",
		Hour11: "P3",
		Hour12: "P3",
		Hour13: "P3",
		Hour14: "P3",
		Hour15: "P3",
		Hour16: "P3",
		Hour17: "P3",
		Hour18: "P3",
		Hour19: "P3",
		Hour20: "P3",
		Hour21: "P3",
		Hour22: "P3",
		Hour23: "P3",
		Hour00: "P3",
	}
	periodHours := measures.CalendarPeriodHours{
		Hour1:  "P1",
		Hour2:  "P1",
		Hour3:  "P2",
		Hour4:  "P2",
		Hour5:  "P2",
		Hour6:  "P2",
		Hour7:  "P2",
		Hour8:  "P2",
		Hour9:  "P2",
		Hour10: "P3",
		Hour11: "P3",
		Hour12: "P3",
		Hour13: "P3",
		Hour14: "P3",
		Hour15: "P3",
		Hour16: "P3",
		Hour17: "P3",
		Hour18: "P3",
		Hour19: "P3",
		Hour20: "P1",
		Hour21: "P1",
		Hour22: "P1",
		Hour23: "P1",
		Hour00: "P1",
	}

	testCases := map[string]struct {
		input  input
		output output
		want   want
	}{
		"Should be meter config err": {
			input: input{dto: dto},
			output: output{
				MeterConfig:    measures.MeterConfig{},
				MeterConfigErr: errors.New("err meter config"),
			},
			want: want{
				result: gross_measures.DashboardMeasureSupplyPoint{},
				err:    errors.New("err meter config"),
			},
		},
		"Should be tariff err": {
			input: input{dto: dto},
			output: output{
				MeterConfig: measures.MeterConfig{},
				Tariff:      clients.Tariffs{},
				TariffErr:   errors.New("err tariff"),
			},
			want: want{
				result: gross_measures.DashboardMeasureSupplyPoint{},
				err:    errors.New("err tariff"),
			},
		},
		"Should be calendar period err": {
			input: input{dto: dto},
			output: output{
				MeterConfig:       measures.MeterConfig{},
				Tariff:            clients.Tariffs{},
				CalendarPeriod:    measures.CalendarPeriod{},
				CalendarPeriodErr: errors.New("err calendar period"),
			},
			want: want{
				result: gross_measures.DashboardMeasureSupplyPoint{},
				err:    errors.New("err calendar period"),
			},
		},
		"Should be gross measures close err": {
			input: input{dto: dto},
			output: output{
				MeterConfig:      measures.MeterConfig{},
				Tariff:           clients.Tariffs{},
				CalendarPeriod:   measures.CalendarPeriod{},
				CloseMeasures:    []gross_measures.MeasureCloseWrite{},
				CloseMeasuresErr: errors.New("err close measures"),
				CurveMeasures:    []gross_measures.MeasureCurveWrite{},
			},
			want: want{
				result: gross_measures.DashboardMeasureSupplyPoint{},
				err:    errors.New("err close measures"),
			},
		},
		"Should be gross measures curve err": {
			input: input{dto: dto},
			output: output{
				MeterConfig:      measures.MeterConfig{},
				Tariff:           clients.Tariffs{},
				CalendarPeriod:   measures.CalendarPeriod{},
				CloseMeasures:    []gross_measures.MeasureCloseWrite{},
				CurveMeasures:    []gross_measures.MeasureCurveWrite{},
				CurveMeasuresErr: errors.New("err curve measures"),
			},
			want: want{
				result: gross_measures.DashboardMeasureSupplyPoint{},
				err:    errors.New("err curve measures"),
			},
		},
		"Should be ok": {
			input: input{dto: dto},
			output: output{
				MeterConfig: measures.MeterConfig{
					ID:            "1",
					StartDate:     time.Date(2022, 1, 1, 22, 0, 0, 0, time.UTC),
					EndDate:       time.Date(2050, 1, 1, 22, 0, 0, 0, time.UTC),
					DistributorID: "DistributorX",
					Type:          measures.TLG,
					AI:            1,
					AE:            1,
					Meter: measures.Meter{
						SerialNumber: "MeterX",
					},
					ServicePoint: measures.ServicePoint{
						Cups:        "CUPS",
						ServiceType: measures.DcServiceType,
					},
					MeasurePoint: measures.MeasurePoint{},
					ContractualSituations: measures.ContractualSituations{
						TariffID: "2.0TD",
					},
				},
				Tariff: clients.Tariffs{
					Id:           "2.0TD",
					GeographicId: "ES",
					CalendarId:   "TD3",
				},
				CalendarPeriod: measures.CalendarPeriod{
					CalendarCode:        "TD3",
					GeographicID:        "ES",
					CalendarPeriodHours: periodHours,
					Festive:             festiveHours,
					StartDate:           time.Time{},
					EndDate:             time.Time{},
					Day:                 dto.EndDate,
					IsFestive:           false,
				},
				CloseMeasures: []gross_measures.MeasureCloseWrite{
					{
						Id:                "2",
						EndDate:           time.Date(2022, 7, 1, 22, 0, 0, 0, time.UTC),
						ReadingDate:       time.Time{},
						GenerationDate:    time.Time{},
						Type:              measures.Absolute,
						Status:            measures.Valid,
						ReadingType:       measures.DailyClosure,
						Contract:          "",
						MeterSerialNumber: "MeterX",
						ConcentratorID:    "",
						File:              "Test/File_2",
						DistributorID:     "DistributorX",
						DistributorCDOS:   "0130",
						Origin:            measures.STG,
						Qualifier:         "",
						Periods: []gross_measures.MeasureClosePeriod{
							{
								Period: "P0",
								AI:     1000000,
								AE:     2300000,
							},
							{
								Period: "P1",
								AI:     10000,
								AE:     23000,
							},
							{
								Period: "P2",
								AI:     10000,
								AE:     23000,
							},
							{
								Period: "P3",
								AI:     5000,
								AE:     2500,
							},
							{
								Period: "P4",
							},
							{
								Period: "P5",
							},
							{
								Period: "P6",
							},
						},
						Invalidations: nil,
					},
					{
						Id:                "3",
						EndDate:           time.Date(2022, 7, 2, 22, 0, 0, 0, time.UTC),
						ReadingDate:       time.Time{},
						GenerationDate:    time.Time{},
						Type:              measures.Absolute,
						Status:            measures.Valid,
						ReadingType:       measures.DailyClosure,
						Contract:          "",
						MeterSerialNumber: "MeterX",
						ConcentratorID:    "",
						File:              "Test/File_3",
						DistributorID:     "DistributorX",
						DistributorCDOS:   "0130",
						Origin:            measures.STG,
						Qualifier:         "",
						Periods: []gross_measures.MeasureClosePeriod{
							{
								Period: "P0",
								AI:     1000000,
								AE:     2300000,
							},
							{
								Period: "P1",
								AI:     10000,
								AE:     23000,
							},
							{
								Period: "P2",
								AI:     10000,
								AE:     23000,
							},
							{
								Period: "P3",
								AI:     5000,
								AE:     2500,
							},
							{
								Period: "P4",
							},
							{
								Period: "P5",
							},
							{
								Period: "P6",
							},
						},
						Invalidations: nil,
					},
					{
						Id:                "4",
						EndDate:           time.Date(2022, 7, 3, 22, 0, 0, 0, time.UTC),
						ReadingDate:       time.Time{},
						GenerationDate:    time.Time{},
						Type:              measures.Absolute,
						Status:            measures.Valid,
						ReadingType:       measures.DailyClosure,
						Contract:          "",
						MeterSerialNumber: "MeterX",
						ConcentratorID:    "",
						File:              "Test/File_4",
						DistributorID:     "DistributorX",
						DistributorCDOS:   "0130",
						Origin:            measures.STG,
						Qualifier:         "",
						Periods: []gross_measures.MeasureClosePeriod{
							{
								Period: "P0",
								AI:     1000000,
								AE:     2300000,
							},
							{
								Period: "P1",
								AI:     10000,
								AE:     23000,
							},
							{
								Period: "P2",
								AI:     10000,
								AE:     23000,
							},
							{
								Period: "P3",
								AI:     5000,
								AE:     2500,
							},
							{
								Period: "P4",
							},
							{
								Period: "P5",
							},
							{
								Period: "P6",
							},
						},
						Invalidations: nil,
					},
					{
						Id:                "5",
						EndDate:           time.Date(2022, 7, 3, 22, 0, 0, 0, time.UTC),
						ReadingDate:       time.Time{},
						GenerationDate:    time.Time{},
						Type:              measures.Absolute,
						Status:            measures.Invalid,
						ReadingType:       measures.DailyClosure,
						Contract:          "",
						MeterSerialNumber: "MeterX",
						ConcentratorID:    "",
						File:              "Test/File_4_1",
						DistributorID:     "DistributorX",
						DistributorCDOS:   "0130",
						Origin:            measures.STG,
						Qualifier:         "",
						Periods: []gross_measures.MeasureClosePeriod{
							{
								Period: "P0",
								AI:     1000000,
								AE:     2300000,
							},
							{
								Period: "P1",
								AI:     10000,
								AE:     23000,
							},
							{
								Period: "P2",
								AI:     10000,
								AE:     23000,
							},
							{
								Period: "P3",
								AI:     5000,
								AE:     2500,
							},
							{
								Period: "P4",
							},
							{
								Period: "P5",
							},
							{
								Period: "P6",
							},
						},
						Invalidations: nil,
					},
					{
						Id:                "6",
						StartDate:         time.Date(2022, 6, 30, 22, 0, 0, 0, time.UTC),
						EndDate:           time.Date(2022, 7, 31, 22, 0, 0, 0, time.UTC),
						ReadingDate:       time.Time{},
						GenerationDate:    time.Time{},
						Type:              measures.Absolute,
						Status:            measures.Valid,
						ReadingType:       measures.BillingClosure,
						Contract:          "",
						MeterSerialNumber: "MeterX",
						ConcentratorID:    "",
						File:              "Test/File_5",
						DistributorID:     "DistributorX",
						DistributorCDOS:   "0130",
						Origin:            measures.STG,
						Qualifier:         "",
						Periods: []gross_measures.MeasureClosePeriod{
							{
								Period: "P0",
								AI:     1000000,
								AE:     2300000,
							},
							{
								Period: "P1",
								AI:     10000,
								AE:     23000,
							},
							{
								Period: "P2",
								AI:     10000,
								AE:     23000,
							},
							{
								Period: "P3",
								AI:     5000,
								AE:     2500,
							},
							{
								Period: "P4",
							},
							{
								Period: "P5",
							},
							{
								Period: "P6",
							},
						},
						Invalidations: nil,
					},
					{
						Id:                "6",
						StartDate:         time.Date(2022, 6, 30, 22, 0, 0, 0, time.UTC),
						EndDate:           time.Date(2022, 7, 31, 22, 0, 0, 0, time.UTC),
						ReadingDate:       time.Time{},
						GenerationDate:    time.Time{},
						Type:              measures.Incremental,
						Status:            measures.Valid,
						ReadingType:       measures.BillingClosure,
						Contract:          "",
						MeterSerialNumber: "MeterX",
						ConcentratorID:    "",
						File:              "Test/File_5",
						DistributorID:     "DistributorX",
						DistributorCDOS:   "0130",
						Origin:            measures.STG,
						Qualifier:         "",
						Periods: []gross_measures.MeasureClosePeriod{
							{
								Period: "P0",
								AI:     1000000,
								AE:     2300000,
							},
							{
								Period: "P1",
								AI:     10000,
								AE:     23000,
							},
							{
								Period: "P2",
								AI:     10000,
								AE:     23000,
							},
							{
								Period: "P3",
								AI:     5000,
								AE:     2500,
							},
							{
								Period: "P4",
							},
							{
								Period: "P5",
							},
							{
								Period: "P6",
							},
						},
						Invalidations: nil,
					},
				},
				CurveMeasures: []gross_measures.MeasureCurveWrite{
					{
						Id:                "1",
						EndDate:           time.Date(2022, 6, 30, 23, 0, 0, 0, time.UTC),
						Type:              measures.Absolute,
						Status:            measures.Invalid,
						CurveType:         measures.HourlyMeasureCurveReadingType,
						MeterSerialNumber: "MeterX",
						File:              "Test/File_1",
						DistributorID:     "DistributorX",
						DistributorCDOS:   "0130",
						Origin:            measures.STG,
					},
					{
						Id:                "2",
						EndDate:           time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC),
						Type:              measures.Absolute,
						Status:            measures.Valid,
						CurveType:         measures.HourlyMeasureCurveReadingType,
						MeterSerialNumber: "MeterX",
						File:              "Test/File_1",
						DistributorID:     "DistributorX",
						DistributorCDOS:   "0130",
						Origin:            measures.STG,
					},
					{
						Id:                "3",
						EndDate:           time.Date(2022, 7, 1, 1, 0, 0, 0, time.UTC),
						Type:              measures.Absolute,
						Status:            measures.Valid,
						CurveType:         measures.HourlyMeasureCurveReadingType,
						MeterSerialNumber: "MeterX",
						File:              "Test/File_1",
						DistributorID:     "DistributorX",
						DistributorCDOS:   "0130",
						Origin:            measures.STG,
					},
					{
						Id:                "4",
						EndDate:           time.Date(2022, 7, 1, 23, 0, 0, 0, time.UTC),
						Type:              measures.Absolute,
						Status:            measures.Valid,
						CurveType:         measures.HourlyMeasureCurveReadingType,
						MeterSerialNumber: "MeterX",
						File:              "Test/File_2",
						DistributorID:     "DistributorX",
						DistributorCDOS:   "0130",
						Origin:            measures.STG,
					},
				},
			},
			want: want{
				result: gross_measures.DashboardMeasureSupplyPoint{
					Cups:            dto.Cups,
					SerialNumber:    "MeterX",
					Magnitudes:      []measures.Magnitude{measures.AI, measures.AE},
					MagnitudeEnergy: measures.AI,
					MeterType:       measures.TLG,
					Periods:         []measures.PeriodKey{measures.P1, measures.P2, measures.P3},
					CalendarDailyClosure: []gross_measures.CalendarStatus{
						{
							Date:   "2022-07-01",
							Status: measures.Valid,
						},
						{
							Date:   "2022-07-02",
							Status: measures.Valid,
						},
						{
							Date:   "2022-07-03",
							Status: measures.Invalid,
						},
					},
					CalendarMonthlyClosure: []gross_measures.CalendarStatus{
						{
							Date:   "2022-07-31",
							Status: measures.Valid,
						},
					},
					CalendarCurve: []gross_measures.CalendarStatus{
						{
							Date:   "2022-07-01",
							Status: measures.Invalid,
						},
						{
							Date:   "2022-07-02",
							Status: measures.Valid,
						},
					},
					ListDailyClosure: []gross_measures.ListDailyClose{
						{

							EndDate: "2022-07-01",
							Status:  measures.Valid,

							Origin: "File_2",
							P0: measures.Values{
								AI: 1000000,
								AE: 2300000,
							},
							P1: measures.Values{
								AI: 10000,
								AE: 23000,
							},
							P2: measures.Values{
								AI: 10000,
								AE: 23000,
							},
							P3: measures.Values{
								AI: 5000,
								AE: 2500,
							},
						},
						{

							EndDate: "2022-07-02",
							Status:  measures.Valid,

							Origin: "File_3",
							P0: measures.Values{
								AI: 1000000,
								AE: 2300000,
							},
							P1: measures.Values{
								AI: 10000,
								AE: 23000,
							},
							P2: measures.Values{
								AI: 10000,
								AE: 23000,
							},
							P3: measures.Values{
								AI: 5000,
								AE: 2500,
							},
						},
						{
							EndDate: "2022-07-03",
							Status:  measures.Valid,

							Origin: "File_4",
							P0: measures.Values{
								AI: 1000000,
								AE: 2300000,
							},
							P1: measures.Values{
								AI: 10000,
								AE: 23000,
							},
							P2: measures.Values{
								AI: 10000,
								AE: 23000,
							},
							P3: measures.Values{
								AI: 5000,
								AE: 2500,
							},
						},
						{

							EndDate: "2022-07-03",
							Status:  measures.Invalid,

							Origin: "File_4_1",
							P0: measures.Values{
								AI: 1000000,
								AE: 2300000,
							},
							P1: measures.Values{
								AI: 10000,
								AE: 23000,
							},
							P2: measures.Values{
								AI: 10000,
								AE: 23000,
							},
							P3: measures.Values{
								AI: 5000,
								AE: 2500,
							},
						},
					},
					ListMonthlyClosure: []gross_measures.ListMonthlyClose{
						{

							EndDate:  "2022-07-31",
							Status:   measures.Valid,
							InitDate: "2022-06-30",
							Origin:   "File_5",
							P0: measures.ValuesMonthly{
								Values: measures.Values{
									AI: 1000000,
									AE: 2300000,
								},
								AIi: 1000000,
								AEi: 2300000,
							},
							P1: measures.ValuesMonthly{
								Values: measures.Values{
									AI: 10000,
									AE: 23000,
								},
								AIi: 10000,
								AEi: 23000,
							},
							P2: measures.ValuesMonthly{
								Values: measures.Values{
									AI: 10000,
									AE: 23000,
								},
								AIi: 10000,
								AEi: 23000,
							},
							P3: measures.ValuesMonthly{
								Values: measures.Values{
									AI: 5000,
									AE: 2500,
								},
								AIi: 5000,
								AEi: 2500,
							},
						},
					},
				},
			},
		},
	}

	loc, _ := time.LoadLocation("Europe/Madrid")
	for name, _ := range testCases {
		test := testCases[name]

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			inventoryRepository := new(mocks.InventoryRepository)
			inventoryRepository.On("GetMeterConfigByCups", mock.Anything, measures.GetMeterConfigByCupsQuery{
				CUPS:        test.input.dto.Cups,
				Time:        test.input.dto.EndDate,
				Distributor: test.input.dto.DistributorId,
			}).Return(test.output.MeterConfig, test.output.MeterConfigErr)

			grossMeasureRepository := new(mocks.GrossMeasureRepository)
			grossMeasureRepository.On("ListCloseMeasures", mock.Anything, gross_measures.QueryListMeasure{
				SerialNumber: test.output.MeterConfig.SerialNumber(),
				StartDate:    test.input.dto.StartDate,
				EndDate:      test.input.dto.EndDate,
			}).Return(test.output.CloseMeasures, test.output.CloseMeasuresErr)

			grossMeasureRepository.On("ListCurveMeasures", mock.Anything, gross_measures.QueryListMeasure{
				SerialNumber: test.output.MeterConfig.SerialNumber(),
				StartDate:    test.input.dto.StartDate,
				EndDate:      test.input.dto.EndDate,
			}).Return(test.output.CurveMeasures, test.output.CurveMeasuresErr)

			tariffClient := new(clients_mocks.MasterTables)
			tariffClient.On("GetTariff", mock.Anything, clients.GetTariffDto{
				ID: test.output.MeterConfig.TariffID(),
			}).Return(test.output.Tariff, test.output.TariffErr)

			calendarPeriodRepository := new(mocks.CalendarPeriodRepository)
			calendarPeriodRepository.On("GetCalendarPeriod", mock.Anything, measures.SearchCalendarPeriod{
				Day:          test.input.dto.EndDate,
				GeographicID: test.output.Tariff.GeographicId,
				CalendarCode: test.output.Tariff.CalendarId,
				Location:     loc,
			}).Return(test.output.CalendarPeriod, test.output.CalendarPeriodErr)

			srv := NewDashboardMeasureSupplyPointService(grossMeasureRepository, inventoryRepository, calendarPeriodRepository, tariffClient, loc)

			dashboardMeasure, err := srv.Handler(context.Background(), test.input.dto)

			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.result, dashboardMeasure)
		})
	}
}
