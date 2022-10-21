package parsers

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	S022 "bitbucket.org/sercide/data-ingestion/internal/gross_measures/parsers/fixtures/S02"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/storage"
	storage_mocks "bitbucket.org/sercide/data-ingestion/pkg/storage/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
	"time"
)

func Test_Unit_Services_PrimeS02_IsFile(t *testing.T) {
	tests := map[string]struct {
		input gross_measures.HandleFileDTO
		want  bool
	}{
		"should be false if empty path": {
			input: gross_measures.HandleFileDTO{
				FilePath: "",
			},
			want: false,
		},
		"should be false if invalid path": {
			input: gross_measures.HandleFileDTO{
				FilePath: "XXX/CIR4621232059_0_S02_0_20160831103005",
			},
			want: false,
		},
		"should be false if invalid path 2": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/distributorID/CIR4621232059_0_S04_0_20160831103005",
			},
			want: false,
		},
		"should be false if invalid path 3": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/distributorID/ss",
			},
			want: false,
		},
		"should be false if invalid format filename": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_2_S02_0_20160831103005",
			},
			want: false,
		},
		"should be false, its Csv_curva": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/0166_CCCH_2201020400_220101000020220102050735.csv",
			},
			want: false,
		},
		"should be false, its Csv_resumen": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/0130_ZIV0004480340_S05_20220325111546W_220325111852.csv",
			},
			want: false,
		},
		"should be false, its TPL": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/20220302083136246_PONTEAREAS_EXT_P02_2022.zip",
			},
			want: false,
		},
		"should be false, its Csv_cierre 1": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/0348_CIERRE_G-D_INC_2112010700_211101000020211201150430.csv",
			},
			want: false,
		},
		"should be false,  its Csv_cierre 2": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/0130_CIERRE_INC_2201010200_2112010000.csv",
			},
			want: false,
		},
		"should be false, its Prime04": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/CIR4621232059_0_S04_0_20211201000120",
			},
			want: false,
		},
		"should be false, its Prime05": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/CIR4621232059_0_S05_0_20170721011003",
			},
			want: false,
		},
		"should be true": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_20160831103005",
			},
			want: true,
		},
	}

	loc, _ := time.LoadLocation("Europe/Madrid")

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			m := new(storage_mocks.Storage)
			storageCreator := func(ctx context.Context) (storage.Storage, error) {
				return m, nil
			}
			parser := NewParser(storageCreator, loc)
			h := NewPrimeS02(parser)
			result := h.IsFile(context.Background(), testCase.input)
			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}

func Test_Unit_Services_PrimeS02_getFileNameFromFilePath(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"should return empty string 1": {
			input: "bucketname",
			want:  "bucketname",
		},
		"should return empty string 2": {
			input: "bucketname/DistribuidorX",
			want:  "DistribuidorX",
		},
		"should return empty string 3": {
			input: "bucketname/DistribuidorX/Input",
			want:  "Input",
		},

		"should return empty string 4": {
			input: "bucketname/DistribuidorX/Input/Prime",
			want:  "Prime",
		},
		"should return file name": {
			input: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_20211111165524",
			want:  "CIR4621232059_0_S02_0_20211111165524",
		},
	}
	loc, _ := time.LoadLocation("Europe/Madrid")
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			m := new(storage_mocks.Storage)
			storageCreator := func(ctx context.Context) (storage.Storage, error) {
				return m, nil
			}
			parser := NewParser(storageCreator, loc)
			h := NewPrimeS02(parser)
			result := h.GetFileNameFromFilePath(testCase.input)
			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}

func Test_Unit_Services_PrimeS02_getDistributorFromFilePath(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"should return empty string 1": {
			want:  "",
			input: "bucketname",
		},
		"should return file name": {
			want:  "DistribuidorX",
			input: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_20211111165524",
		},
	}

	loc, _ := time.LoadLocation("Europe/Madrid")
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			m := new(storage_mocks.Storage)
			storageCreator := func(ctx context.Context) (storage.Storage, error) {
				return m, nil
			}
			parser := NewParser(storageCreator, loc)
			h := NewPrimeS02(parser)
			result := h.GetDistributorFromFilePath(testCase.input)
			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}

func Test_Unit_Services_PrimeS02_getMeasureDateFromFilePath(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"should return empty string 1": {
			want:  "",
			input: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_0_20211111165524",
		},
		"should return empty string 2": {
			want:  "",
			input: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_23_0_20211111165524",
		},
		"should return file name": {
			want:  "20211111165524",
			input: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_20211111165524",
		},
	}

	loc, _ := time.LoadLocation("Europe/Madrid")
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			m := new(storage_mocks.Storage)
			storageCreator := func(ctx context.Context) (storage.Storage, error) {
				return m, nil
			}
			parser := NewParser(storageCreator, loc)
			h := NewPrimeS02(parser)
			result := h.getMeasureDateFromFilePath(testCase.input)
			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}

func Test_Unit_Services_PrimeS02_Handle(t *testing.T) {
	tests := map[string]struct {
		input             gross_measures.HandleFileDTO
		want              []gross_measures.MeasureCurveWrite
		err               error
		fixtureFile       string
		storageCreatorErr error
		readAllErr        error
	}{
		"should return error if error in create storage": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_20160831103005",
			},
			readAllErr:        nil,
			storageCreatorErr: errors.New(""),
			want:              []gross_measures.MeasureCurveWrite{},
			err:               errors.New(""),
			fixtureFile:       "fixtures/S02/CIR4621232059_0_S02_0_20160831103005",
		},
		"should return error if error in readAll file": {
			input: gross_measures.HandleFileDTO{
				FilePath: "",
			},
			storageCreatorErr: nil,
			readAllErr:        errors.New(""),
			want:              []gross_measures.MeasureCurveWrite{},
			err:               errors.New(""),
			fixtureFile:       "",
		},
		"should return error if error in umarshall with empty file": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_20160831103005",
			},
			storageCreatorErr: nil,
			readAllErr:        nil,
			want:              []gross_measures.MeasureCurveWrite{},
			err:               errors.New("EOF"),
			fixtureFile:       "fixtures/S02/empty",
		},
		"should return empty measures if not cnc": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_20160831103005",
			},
			storageCreatorErr: nil,
			readAllErr:        nil,
			want:              []gross_measures.MeasureCurveWrite{},
			err:               nil,
			fixtureFile:       "fixtures/S02/empty_Cnc",
		},
		"should return empty measures if not cnt": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_20160831103005",
			},
			storageCreatorErr: nil,
			readAllErr:        nil,
			want:              []gross_measures.MeasureCurveWrite{},
			err:               nil,
			fixtureFile:       "fixtures/S02/empty_Cnt",
		},
		"should return empty measures if only cnt in error": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_20160831103005",
			},
			storageCreatorErr: nil,
			readAllErr:        nil,
			want:              []gross_measures.MeasureCurveWrite{},
			err:               nil,
			fixtureFile:       "fixtures/S02/only_errors_Cnt",
		},
		"should return measures with valid date": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_20160831103005",
			},
			storageCreatorErr: nil,
			readAllErr:        nil,
			want: []gross_measures.MeasureCurveWrite{
				{
					EndDate:           time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:       time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					Type:              measures.Incremental,
					ReadingType:       measures.Curve,
					CurveType:         measures.HourlyMeasureCurveReadingType,
					MeterSerialNumber: "SAG0125150903",
					ConcentratorID:    "CIR4621232059",
					File:              "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_20160831103005",
					DistributorID:     "",
					DistributorCDOS:   "DistribuidorX",
					Origin:            measures.STG,
					Qualifier:         "0",
					AI:                0,
					AE:                0,
					R1:                0,
					R2:                0,
					R3:                0,
					R4:                0,
				},
			},
			err:         nil,
			fixtureFile: "fixtures/S02/invalid_enddate",
		},
		"should return error if date measure in file es incorrect": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_",
			},
			storageCreatorErr: nil,
			readAllErr:        nil,
			want:              []gross_measures.MeasureCurveWrite{},
			err:               errors.New("invalid date"),
			fixtureFile:       "fixtures/S02/invalid_enddate",
		},
		"should parse CIR4621232059_0_S02_0_20160831103005 file well": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_20160831103005",
			},
			readAllErr:        nil,
			storageCreatorErr: nil,
			want:              S022.CIR4621232059_0_S02_0_20160831103005_Result,
			err:               nil,
			fixtureFile:       "fixtures/S02/CIR4621232059_0_S02_0_20160831103005",
		},
		"should parse CIR4621232059_0_S02_0_20211111165524 file well": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S02_0_20211111165524",
			},
			readAllErr:        nil,
			storageCreatorErr: nil,
			want:              S022.CIR4621232059_0_S02_0_20211111165524_Result,
			err:               nil,
			fixtureFile:       "fixtures/S02/CIR4621232059_0_S02_0_20211111165524",
		},
	}

	loc, _ := time.LoadLocation("Europe/Madrid")
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			m := new(storage_mocks.Storage)

			storageCreator := func(ctx context.Context) (storage.Storage, error) {
				return m, testCase.storageCreatorErr
			}

			fileContent, err := os.ReadFile(testCase.fixtureFile)

			m.On("ReadAll", mock.Anything, mock.Anything).Return(fileContent, testCase.readAllErr)
			m.On("Close").Return(nil)

			parser := NewParser(storageCreator, loc)
			h := NewPrimeS02(parser)

			result, err := h.Handle(context.Background(), testCase.input)

			assert.Equal(t, testCase.want, result, testCase)
			assert.Equal(t, testCase.err, err, testCase)
		})
	}
}
