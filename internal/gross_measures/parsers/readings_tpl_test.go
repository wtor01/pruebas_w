package parsers

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/parsers/fixtures/readings_tpl"
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

func Test_Unit_Services_ReadingsTpl_IsFile(t *testing.T) {
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
		"should be false if empty path TPL": {
			input: gross_measures.HandleFileDTO{
				FilePath: "TPL",
			},
			want: false,
		},
		"should be false if invalid path": {
			input: gross_measures.HandleFileDTO{
				FilePath: "XXX/20220302083136246_PONTEAREAS_EXT_P02_2022.zip",
			},
			want: false,
		},
		"should be false if invalid path 2": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/distributorID/iivm_20220302083136246_PONTEAREAS_EXT_P02_2022.zip",
			},
			want: false,
		},
		"should be false if invalid path iicc": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/distributorID/Input/TPL/iicc_20220302083136246_PONTEAREAS_EXT_P02_2022.zip",
			},
			want: false,
		},
		"should be false if invalid path iivm": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/distributorID/Input/TPL/iivm_20220302083136246_PONTEAREAS_EXT_P02_2022.zip",
			},
			want: false,
		},
		"should be false if invalid path 5": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/distributorID/ss",
			},
			want: false,
		},
		"should be false if invalid format filename": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/TPL/20220302083136246_PONTEAREAS_EXT_P02_2022.XX",
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
		"should be true": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/TPL/20220302083136246_PONTEAREAS_EXT_P02_2022.zip",
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
			h := NewReadingsTpl(parser)
			result := h.IsFile(context.Background(), testCase.input)
			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}

func Test_Unit_Services_ReadingsTpl_Handle(t *testing.T) {
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
				FilePath: "bucketname/DistribuidorX/Input/TPL/20220302083136246_PONTEAREAS_EXT_P02_2022.zip",
			},
			readAllErr:        nil,
			storageCreatorErr: errors.New(""),
			want:              []gross_measures.MeasureCloseWrite{},
			err:               errors.New(""),
			fixtureFile:       "fixtures/readings_tpl/20220302083136246_PONTEAREAS_EXT_P02_2022.zip",
		},
		"should return error if error in NewReader file": {
			input: gross_measures.HandleFileDTO{
				FilePath: "",
			},
			storageCreatorErr: nil,
			readAllErr:        errors.New(""),
			want:              []gross_measures.MeasureCloseWrite{},
			err:               errors.New(""),
			fixtureFile:       "",
		},
		"should parse 20220302083136246_PONTEAREAS_EXT_P02_2022.zip file well": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/TPL/20220302083136246_PONTEAREAS_EXT_P02_2022.zip",
			},
			readAllErr:        nil,
			storageCreatorErr: nil,
			want:              readings_tpl.Result_20220302083136246_PONTEAREAS_EXT_P02_2022,
			err:               nil,
			fixtureFile:       "fixtures/readings_tpl/20220302083136246_PONTEAREAS_EXT_P02_2022.zip",
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
			h := NewReadingsTpl(parser)

			result, err := h.Handle(context.Background(), testCase.input)
			if testCase.input.FilePath != "" {
				_, errTmpFile := os.Stat(h.GetTmpFilePath(testCase.input.FilePath))
				_, errZipTmpFile := os.Stat(h.GetZipTmpFilePath(testCase.input.FilePath))
				assert.Equal(t, true, errors.Is(errTmpFile, os.ErrNotExist))
				assert.Equal(t, true, errors.Is(errZipTmpFile, os.ErrNotExist))

			}
			orderMeasuresClose(testCase.want)
			orderMeasuresClose(result)
			assert.ElementsMatch(t, testCase.want, result, testCase)

			if testCase.err != nil {
				assert.Error(t, err, testCase)
			} else {
				assert.Equal(t, testCase.err, err, testCase)
			}
		})
	}
}

func Test_Unit_Services_ReadingsTpl_setRowValuesToMeasure(t *testing.T) {
	tests := map[string]struct {
		inputRow     []string
		inputMeasure gross_measures.MeasureClosePeriod
		want         gross_measures.MeasureClosePeriod
		err          error
	}{
		"should parse AbsA+": {
			inputRow:     []string{"0", "AbsA+", "P1 TD6", "149111"},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{AI: 149111000},
			err:          nil,
		},
		"should not parse AbsA+": {
			inputRow:     []string{"0", "AbsA+", "P1 TD6", ""},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{},
			err:          errors.New(""),
		},
		"should parse AbsA-": {
			inputRow:     []string{"0", "AbsA-", "P1 TD6", "149111"},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{AE: 149111000},
			err:          nil,
		},
		"should not parse AbsA-": {
			inputRow:     []string{"0", "AbsA-", "P1 TD6", ""},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{},
			err:          errors.New(""),
		},
		"should parse AbsRi+": {
			inputRow:     []string{"0", "AbsRi+", "P1 TD6", "149111"},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{R1: 149111000},
			err:          nil,
		},
		"should not parse AbsRi+": {
			inputRow:     []string{"0", "AbsRi+", "P1 TD6", ""},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{},
			err:          errors.New(""),
		},
		"should parse AbsRc+": {
			inputRow:     []string{"0", "AbsRc+", "P1 TD6", "149111"},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{R2: 149111000},
			err:          nil,
		},
		"should not parse AbsRc+": {
			inputRow:     []string{"0", "AbsRc+", "P1 TD6", ""},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{},
			err:          errors.New(""),
		},
		"should parse AbsRi-": {
			inputRow:     []string{"0", "AbsRi-", "P1 TD6", "149111"},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{R3: 149111000},
			err:          nil,
		},
		"should not parse AbsRi-": {
			inputRow:     []string{"0", "AbsRi-", "P1 TD6", ""},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{},
			err:          errors.New(""),
		},
		"should parse AbsRc-": {
			inputRow:     []string{"0", "AbsRc-", "P1 TD6", "149111"},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{R4: 149111000},
			err:          nil,
		},
		"should not parse AbsRc-": {
			inputRow:     []string{"0", "AbsRc-", "P1 TD6", ""},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{},
			err:          errors.New(""),
		},
		"should parse MaxA": {
			inputRow:     []string{"0", "MaxA", "P1 TD6", "149111"},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{MX: 149111000},
			err:          nil,
		},
		"should not parse MaxA": {
			inputRow:     []string{"0", "MaxA", "P1 TD6", ""},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{},
			err:          errors.New(""),
		},
		"should parse ExcA": {
			inputRow:     []string{"0", "ExcA", "P1 TD6", "149111"},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{E: 149111000},
			err:          nil,
		},
		"should not parse ExcA": {
			inputRow:     []string{"0", "ExcA", "P1 TD6", ""},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{},
			err:          errors.New(""),
		},
		"should not parse invalid magnitude": {
			inputRow:     []string{"0", "test", "P1 TD6", ""},
			inputMeasure: gross_measures.MeasureClosePeriod{},
			want:         gross_measures.MeasureClosePeriod{},
			err:          nil,
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
			h := NewReadingsTpl(parser)

			err := h.setRowValuesToMeasure(&testCase.inputMeasure, testCase.inputRow[1], testCase.inputRow[3])

			assert.Equal(t, testCase.want, testCase.inputMeasure, testCase)
			if testCase.err != nil {
				assert.Error(t, err, testCase)
			} else {
				assert.Equal(t, testCase.err, err, testCase)
			}
		})
	}
}
