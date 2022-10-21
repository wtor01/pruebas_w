package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/services/fixtures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/pkg/storage"
	storage_mocks "bitbucket.org/sercide/data-ingestion/pkg/storage/mocks"
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Unit_Service_Billing_Consum_Insert(t *testing.T) {

	tests := map[string]struct {
		input             InsertConsumProfileDTO
		want              error
		newReaderErr      error
		storageCreatorErr error
		saveErr           error
		ConsumProfiles    []billing_measures.ConsumProfile
	}{
		"should return err if storageCreator return err": {
			input: InsertConsumProfileDTO{
				File: "fixtures/Datos_coeficientes_consumo.csv",
			},
			want:              errors.New("storageCreatorErr"),
			newReaderErr:      nil,
			storageCreatorErr: errors.New("storageCreatorErr"),
			saveErr:           nil,
		},
		"should return err if newReader return err": {
			input: InsertConsumProfileDTO{
				File: "fixtures/Datos_coeficientes_consumo.csv",
			},
			want:              errors.New("newReader"),
			newReaderErr:      errors.New("newReader"),
			storageCreatorErr: nil,
			saveErr:           nil,
		},
		"should return err if save return err": {
			input: InsertConsumProfileDTO{
				File: "fixtures/Datos_coeficientes_consumo.csv",
			},
			want:              errors.New("error saving ConsumProfile from file"),
			newReaderErr:      nil,
			storageCreatorErr: nil,
			saveErr:           errors.New("saveErr"),
			ConsumProfiles:    fixtures.ConsumProfileInsertFile,
		},
		"should return err if invalid headers file": {
			input: InsertConsumProfileDTO{
				File: "fixtures/Datos_coeficientes_consumo_error_header.csv",
			},
			want:              errors.New("error invalid header file"),
			newReaderErr:      nil,
			storageCreatorErr: nil,
			saveErr:           nil,
			ConsumProfiles:    nil,
		},
		"should call save with corect ConsumProfiles": {
			input: InsertConsumProfileDTO{
				File: "fixtures/Datos_coeficientes_consumo.csv",
			},
			want:              nil,
			newReaderErr:      nil,
			storageCreatorErr: nil,
			saveErr:           nil,
			ConsumProfiles:    fixtures.ConsumProfileInsertFile,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			repo := new(mocks.ConsumProfileRepository)
			m := new(storage_mocks.Storage)

			storageCreator := func(ctx context.Context) (storage.Storage, error) {
				return m, testCase.storageCreatorErr
			}

			ctx := context.Background()

			fileContent, _ := os.ReadFile(testCase.input.File)
			myReader := bytes.NewReader(fileContent)

			m.On("NewReader", ctx, testCase.input.File).Return(myReader, testCase.newReaderErr)

			for _, c := range testCase.ConsumProfiles {
				repo.Mock.On("Save", ctx, c).Return(testCase.saveErr)
			}

			srv := NewInsertConsumProfile(storageCreator, repo)

			err := srv.Handler(ctx, testCase.input)

			assert.Equal(t, testCase.want, err)
		})
	}
}
