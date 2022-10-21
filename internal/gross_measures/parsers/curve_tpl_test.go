package parsers

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/parsers/fixtures/curve_tpl"
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

func Test_Unit_Services_CurveTpl_IsFile(t *testing.T) {
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
				FilePath: "XXX/iicc_594001_1555_2_20210803.tpl",
			},
			want: false,
		},
		"should be false if invalid path 2": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/distributorID/iicc_594001_1555_2_20210803.tpl",
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
				FilePath: "bucketname/DistribuidorX/Input/TPL/iicc_594001_1555_2_20210803.XX",
			},
			want: false,
		},
		"should be true": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/TPL/iicc_594001_1555_2_20210803.tpl",
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
			h := NewCurveTpl(parser)
			result := h.IsFile(context.Background(), testCase.input)
			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}

func Test_Unit_Services_CurveTpl_Handle(t *testing.T) {
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
				FilePath: "bucketname/DistribuidorX/Input/TPL/iicc_594001_1555_2_20210803.tpl",
			},
			readAllErr:        nil,
			storageCreatorErr: errors.New(""),
			want:              []gross_measures.MeasureCurveWrite{},
			err:               errors.New(""),
			fixtureFile:       "fixtures/curve_tpl/iicc_594001_1555_2_20210803.tpl",
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
		"should parse iicc_594001_1555_2_20210803.tpl file well": {
			input: gross_measures.HandleFileDTO{
				FilePath: "bucketname/DistribuidorX/Input/TPL/iicc_594001_1555_2_20210803.tpl",
			},
			readAllErr:        nil,
			storageCreatorErr: nil,
			want:              curve_tpl.Result_iicc_594001_1555_2_20210803,
			err:               nil,
			fixtureFile:       "fixtures/curve_tpl/iicc_594001_1555_2_20210803.tpl",
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
			h := NewCurveTpl(parser)

			result, err := h.Handle(context.Background(), testCase.input)

			assert.Equal(t, testCase.want, result, testCase)
			assert.Equal(t, testCase.err, err, testCase)
		})
	}
}
