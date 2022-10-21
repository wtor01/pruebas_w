package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"context"
	"errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Measure_Services_CreateValidationMeasureService_Handle(t *testing.T) {
	tests := map[string]struct {
		want       validations.ValidationMeasure
		wantErr    bool
		input      CreateValidationMeasureDto
		saveReturn error
	}{
		"should return error if validations fail": {
			wantErr: true,
			input:   CreateValidationMeasureDto{},
		},
		"should return error and empty ValidationRule if SaveValidationMeasure fail": {
			wantErr: true,
			input: CreateValidationMeasureDto{
				Id:          uuid.NewV4().String(),
				UserID:      uuid.NewV4().String(),
				Name:        "Umbral",
				Action:      string(measures.Invalid),
				Enabled:     false,
				MeasureType: string(measures.Incremental),
				Type:        string(validations.Immediate),
				Code:        "V001",
				Message:     "fail",
				Description: "description",
				Params: validations.Params{
					Type: validations.Simple,
					Validations: []validations.Validation{
						{
							Id:       uuid.NewV4().String(),
							Type:     validations.Threshold,
							Keys:     []string{validations.AI},
							Required: false,
							Config: map[string]string{
								"max": "90",
								"min": "10",
							},
						},
					},
				},
			},
			saveReturn: errors.New(""),
		},
		"should return ok if SaveValidationMeasure success": {
			wantErr: false,
			input: CreateValidationMeasureDto{
				Id:          uuid.NewV4().String(),
				UserID:      uuid.NewV4().String(),
				Name:        "Umbral",
				Action:      string(measures.Invalid),
				Enabled:     false,
				MeasureType: string(measures.Incremental),
				Type:        validations.Immediate,
				Code:        "V001",
				Message:     "fail",
				Description: "description",
				Params: validations.Params{
					Type: validations.Simple,
					Validations: []validations.Validation{
						{
							Id:       uuid.NewV4().String(),
							Type:     validations.Threshold,
							Keys:     []string{validations.AI},
							Required: false,
							Config: map[string]string{
								"max": "90",
								"min": "10",
							},
						},
					},
				},
			},
			saveReturn: nil,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			repository := new(mocks.AdminRepository)
			validation, _ := validations.NewValidationMeasure(
				testCase.input.Id,
				testCase.input.UserID,
				testCase.input.Name,
				testCase.input.Action,
				testCase.input.Enabled,
				testCase.input.MeasureType,
				testCase.input.Type,
				testCase.input.Code,
				testCase.input.Message,
				testCase.input.Description,
				testCase.input.Params,
			)
			repository.On("SaveValidationMeasure", ctx, validation).Return(testCase.saveReturn)

			service := NewCreateValidationMeasureService(repository)

			result, err := service.Handle(ctx, testCase.input)

			if testCase.wantErr {
				assert.Equal(t, validations.ValidationMeasure{}, result)
				assert.Error(t, err)
			} else {
				assert.Equal(t, validation, result)
				assert.Nil(t, err)
			}
		})
	}
}
