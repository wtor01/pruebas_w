package parsers

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/parsers/fixtures/csv_cierre"
	"bitbucket.org/sercide/data-ingestion/pkg/storage"
	storage_mocks "bitbucket.org/sercide/data-ingestion/pkg/storage/mocks"
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
	"time"
)

func Test_Unit_Services_CierreCsv_IsFile(t *testing.T) {
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
				FilePath: "XXX/0348_CIERRE_G-D_INC_2112010700_211101000020211201150430",
			},
			want: false,
		},
		"should be false if invalid path 2": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/distributorID/0348_CIERRE_G-D_INC_2112010700_211101000020211201150430",
			},
			want: false,
		},
		"should be false if invalid path 3": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/distributorID/ss",
			},
			want: false,
		},
		"should be false if invalid lench filename, its big": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/0348_CIERRE_G-D_INC_s_2112010700_211101000020211201150430.csv",
			},
			want: false,
		},
		"should be false if invalid lench filename, its short": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/0348_CIERRE_2112010700_211101000020211201150430.csv",
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
		"should be true for big path": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/0348_CIERRE_G-D_INC_2112010700_211101000020211201150430.csv",
			},
			want: true,
		},
		"should be true for short path": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/0130_CIERRE_INC_2201010200_2112010000.csv",
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
			h := NewCloseCsv(parser)
			result := h.IsFile(context.Background(), testCase.input)
			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}

func Test_Unit_Services_CierreCsv_Handle(t *testing.T) {
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
				FilePath: "bucketname/DistribuidorX/Input/CSV/0130_CIERRE_INC_2201010200_2112010000",
			},
			readAllErr:        nil,
			storageCreatorErr: errors.New(""),
			want:              []gross_measures.MeasureCloseWrite{},
			err:               errors.New(""),
			fixtureFile:       "fixtures/csv_cierre/0130_CIERRE_INC_2201010200_2112010000",
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
		"should parse 0348_CIERRE_G-D_INC_2112010700_211101000020211201150430.csv file well": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/CSV/0130_CIERRE_INC_2201010200_2112010000.csv",
			},
			readAllErr:        nil,
			storageCreatorErr: nil,
			want:              csv_cierre.Result_0130_CIERRE_INC_2201010200_2112010000,
			err:               nil,
			fixtureFile:       "fixtures/csv_cierre/0130_CIERRE_INC_2201010200_2112010000.csv",
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
			myReader := bytes.NewReader(fileContent)

			m.On("NewReader", mock.Anything, mock.Anything).Return(myReader, testCase.readAllErr)
			m.On("Close").Return(nil)

			parser := NewParser(storageCreator, loc)
			h := NewCloseCsv(parser)

			result, err := h.Handle(context.Background(), testCase.input)
			orderMeasuresClose(testCase.want)
			orderMeasuresClose(result)
			assert.ElementsMatch(t, testCase.want, result, testCase)
			assert.Equal(t, testCase.err, err, testCase)

		})
	}
}
