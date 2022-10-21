package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Unit_Domain_ValidatorQuarterly_New(t *testing.T) {
	tests := map[string]struct {
		want struct {
			err   bool
			value ValidatorQuarterly
		}
		input struct {
			v      ValidationData
			status measures.Status
		}
	}{
		"Should return err with invalid config": {
			want: struct {
				err   bool
				value ValidatorQuarterly
			}{
				err:   true,
				value: ValidatorQuarterly{},
			},
			input: struct {
				v      ValidationData
				status measures.Status
			}{
				v: ValidationData{
					Type:   QuarterlyDate,
					Keys:   []string{StartDate},
					Config: map[string]string{"test": "test"},
				},
				status: measures.Invalid,
			},
		},
		"Should return err with invalid keys": {
			want: struct {
				err   bool
				value ValidatorQuarterly
			}{
				err:   true,
				value: ValidatorQuarterly{},
			},
			input: struct {
				v      ValidationData
				status measures.Status
			}{
				v: ValidationData{
					Type:   QuarterlyDate,
					Keys:   []string{MX},
					Config: map[string]string{},
				},
				status: measures.Invalid,
			},
		},
		"Should return valid ValidatorQuarterly": {
			want: struct {
				err   bool
				value ValidatorQuarterly
			}{
				err: false,
				value: ValidatorQuarterly{
					ValidatorBase: ValidatorBase{
						Type:   QuarterlyDate,
						Keys:   []string{StartDate, EndDate, MeasureDate},
						Config: map[string]string{},
						Action: measures.Invalid,
					},
				},
			},
			input: struct {
				v      ValidationData
				status measures.Status
			}{
				v: ValidationData{
					Type:   QuarterlyDate,
					Keys:   []string{StartDate, EndDate, MeasureDate},
					Config: map[string]string{},
				},
				status: measures.Invalid,
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {

			validator, err := NewValidatorQuarterly(testCase.input.v, testCase.input.status)

			assert.Equal(t, testCase.want.value, validator, testName)
			if testCase.want.err {
				assert.NotNil(t, err, testName)
			} else {
				assert.Nil(t, err, testName)
			}
		})
	}
}

func Test_Unit_Domain_ValidatorQuarterly_Validator(t *testing.T) {
	validDate := time.Date(2022, 05, 05, 00, 30, 00, 000, time.UTC)

	invalidDate := time.Date(2022, 05, 05, 10, 05, 00, 000, time.UTC)

	tests := map[string]struct {
		input          MeasureValidatable
		validator      ValidatorQuarterly
		failValidation bool
	}{
		"should update measure status if invalid date in EndDate to measures.Invalid": {
			input: MeasureValidatable{
				StartDate:   validDate,
				EndDate:     invalidDate,
				ReadingDate: validDate,
				FX:          validDate,
				WhiteListKeys: map[string]struct{}{
					StartDate:   {},
					EndDate:     {},
					MeasureDate: {},
					FX:          {},
				},
			},
			validator: ValidatorQuarterly{
				ValidatorBase: ValidatorBase{
					Type:   QuarterlyDate,
					Keys:   []string{StartDate, EndDate, MeasureDate, FX},
					Config: map[string]string{},
					Action: measures.Invalid,
					Code:   QuarterlyDate,
				},
			},
			failValidation: true,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {

			base := testCase.validator.Validate(testCase.input)
			if testCase.failValidation {
				assert.Equal(t, &testCase.validator.ValidatorBase, base, testCase)
			} else {
				assert.Nil(t, base, testCase)
			}
		})
	}
}
