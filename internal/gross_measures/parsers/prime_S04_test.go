package parsers

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	S044 "bitbucket.org/sercide/data-ingestion/internal/gross_measures/parsers/fixtures/S04"
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

func Test_Unit_Services_PrimeS04_IsFile(t *testing.T) {
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
				FilePath: "XXX/CIR4621232059_0_S04_0_20211201000120",
			},
			want: false,
		},
		"should be false if invalid path 2": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/distributorID/CIR4621232059_0_S04_0_20211201000120",
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
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_2_S04_0_20211201000120",
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
		"should be false, its Prime02": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/CIR4621232059_0_S02_0_20160831103005",
			},
			want: false,
		},
		"should be false, its Prime02 2": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/CIR4621232059_0_S02_0_20211111165524",
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
		"should be false, its Prime05": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/CIR4621232059_0_S05_0_20170721011003",
			},
			want: false,
		},
		"should be true": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S04_0_20211201000120",
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
			h := NewPrimeS04(parser)
			result := h.IsFile(context.Background(), testCase.input)
			assert.Equal(t, testCase.want, result, testCase)
		})
	}

}

func Test_Unit_Services_PrimeS04_getFileNameFromFilePath(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"should return bucketname string 1": {
			input: "bucketname",
			want:  "bucketname",
		},
		"should return DistribuidorX string 2": {
			input: "bucketname/DistribuidorX",
			want:  "DistribuidorX",
		},
		"should return Input string 3": {
			input: "bucketname/DistribuidorX/Input",
			want:  "Input",
		},

		"should return Prime string 4": {
			input: "bucketname/DistribuidorX/Input/Prime",
			want:  "Prime",
		},
		"should return file name": {
			input: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S04_0_20211201000120",
			want:  "CIR4621232059_0_S04_0_20211201000120",
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
			h := NewPrimeS04(parser)
			result := h.GetFileNameFromFilePath(testCase.input)
			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}

func Test_Unit_Services_PrimeS04_getDistributorFromFilePath(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"should return empty string 1": {
			input: "bucketname",
			want:  "",
		},
		"should return empty string 2": {
			input: "bucketname/DistribuidorX",
			want:  "DistribuidorX",
		},
		"should return empty string 3": {
			input: "bucketname/DistribuidorX/Input",
			want:  "DistribuidorX",
		},

		"should return empty string 4": {
			input: "bucketname/DistribuidorX/Input/Prime",
			want:  "DistribuidorX",
		},
		"should return file name": {
			input: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S04_0_20211201000120",
			want:  "DistribuidorX",
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
			h := NewPrimeS04(parser)
			result := h.GetDistributorFromFilePath(testCase.input)
			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}

func Test_Unit_Services_PrimeS04_getMeasureDateFromFilePath(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"should return empty string 1": {
			input: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_0_20211201000120",
			want:  "",
		},
		"should return empty string 2": {
			input: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S04_23_0_20211201000120",
			want:  "",
		},
		"should return file name": {
			input: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S04_0_20211201000120",
			want:  "20211201000120",
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
			h := NewPrimeS04(parser)
			result := h.getMeasureDateFromFilePath(testCase.input)
			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}

func Test_Unit_Services_PrimeS04_Handle(t *testing.T) {
	tests := map[string]struct {
		input             gross_measures.HandleFileDTO
		want              []gross_measures.MeasureCloseWrite
		err               error
		fixtureFile       string
		storageCreatorErr error
		readAllErr        error
	}{
		"should return error if error in create storage": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S04_0_20211201000120",
			},
			readAllErr:        nil,
			storageCreatorErr: errors.New(""),
			want:              []gross_measures.MeasureCloseWrite{},
			err:               errors.New(""),
			fixtureFile:       "fixtures/S04/CIR4621232059_0_S04_0_20211201000120",
		},
		"should return error if error in readAll file": {
			input: gross_measures.HandleFileDTO{
				FilePath: "",
			},
			storageCreatorErr: nil,
			readAllErr:        errors.New(""),
			want:              []gross_measures.MeasureCloseWrite{},
			err:               errors.New(""),
			fixtureFile:       "",
		},
		"should return error if error in umarshall with empty file": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S04_0_20211201000120",
			},
			storageCreatorErr: nil,
			readAllErr:        nil,
			want:              []gross_measures.MeasureCloseWrite{},
			err:               errors.New("EOF"),
			fixtureFile:       "fixtures/S04/empty",
		},
		"should return empty measures if not cnc": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S04_0_20211201000120",
			},
			storageCreatorErr: nil,
			readAllErr:        nil,
			want:              []gross_measures.MeasureCloseWrite{},
			err:               nil,
			fixtureFile:       "fixtures/S04/empty_Cnc",
		},
		"should return empty measures if not cnt": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S04_0_20211201000120",
			},
			storageCreatorErr: nil,
			readAllErr:        nil,
			want:              []gross_measures.MeasureCloseWrite{},
			err:               nil,
			fixtureFile:       "fixtures/S04/empty_Cnt",
		},
		"should return empty measures if only cnt in error": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S04_0_20211201000120",
			},
			storageCreatorErr: nil,
			readAllErr:        nil,
			want:              []gross_measures.MeasureCloseWrite{},
			err:               nil,
			fixtureFile:       "fixtures/S04/only_errors_Cnt",
		},
		"should return error if date measure in file es incorrect": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S04_0_",
			},
			storageCreatorErr: nil,
			readAllErr:        nil,
			want:              []gross_measures.MeasureCloseWrite{},
			err:               errors.New("invalid date"),
			fixtureFile:       "fixtures/S04/invalid_enddate",
		},
		"should parse CIR4621232059_0_S04_0_20211201000120 file well": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/Prime/CIR4621232059_0_S04_0_20211201000120",
			},
			readAllErr:        nil,
			storageCreatorErr: nil,
			want:              S044.CIR4621232059_0_S04_0_20211201000120_Result,
			err:               nil,
			fixtureFile:       "fixtures/S04/CIR4621232059_0_S04_0_20211201000120",
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
			h := NewPrimeS04(parser)

			result, err := h.Handle(context.Background(), testCase.input)
			orderMeasuresClose(testCase.want)
			orderMeasuresClose(result)

			assert.ElementsMatch(t, testCase.want, result, testCase)
			assert.Equal(t, testCase.err, err, testCase)
		})
	}
}
