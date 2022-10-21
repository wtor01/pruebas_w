package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Unit_Domain_ValidatorDaily_New(t *testing.T) {
	tests := map[string]struct {
		want struct {
			err   bool
			value ValidatorDaily
		}
		input struct {
			v      ValidationData
			status measures.Status
		}
	}{
		"Should return err with invalid config": {
			want: struct {
				err   bool
				value ValidatorDaily
			}{
				err:   true,
				value: ValidatorDaily{},
			},
			input: struct {
				v      ValidationData
				status measures.Status
			}{
				v: ValidationData{
					Type:   DailyDate,
					Keys:   []string{StartDate},
					Config: map[string]string{"test": "test"},
				},
				status: measures.Invalid,
			},
		},
		"Should return err with invalid keys": {
			want: struct {
				err   bool
				value ValidatorDaily
			}{
				err:   true,
				value: ValidatorDaily{},
			},
			input: struct {
				v      ValidationData
				status measures.Status
			}{
				v: ValidationData{
					Type:   DailyDate,
					Keys:   []string{MX},
					Config: map[string]string{},
				},
				status: measures.Invalid,
			},
		},
		"Should return valid ValidatorDaily": {
			want: struct {
				err   bool
				value ValidatorDaily
			}{
				err: false,
				value: ValidatorDaily{
					ValidatorBase: ValidatorBase{
						Type:   DailyDate,
						Keys:   []string{StartDate, EndDate, MeasureDate, FX},
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
					Type:   DailyDate,
					Keys:   []string{StartDate, EndDate, MeasureDate, FX},
					Config: map[string]string{},
				},
				status: measures.Invalid,
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {

			validator, err := NewValidatorDaily(testCase.input.v, testCase.input.status)

			assert.Equal(t, testCase.want.value, validator, testName)
			if testCase.want.err {
				assert.NotNil(t, err, testName)
			} else {
				assert.Nil(t, err, testName)
			}
		})
	}
}

func Test_Unit_Domain_ValidatorDaily_Validator(t *testing.T) {
	validDate := time.Date(2022, 05, 05, 00, 00, 00, 000, time.UTC)

	invalidDate := time.Date(2022, 05, 05, 00, 45, 00, 000, time.UTC)

	tests := map[string]struct {
		input      MeasureValidatable
		validator  ValidatorDaily
		returnBase bool
	}{
		"should return validator base if invalid date in EndDate": {
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
			validator: ValidatorDaily{
				ValidatorBase: ValidatorBase{
					Type:   DailyDate,
					Keys:   []string{StartDate, EndDate, MeasureDate, FX},
					Config: map[string]string{},
					Action: measures.Invalid,
					Code:   DailyDate,
				},
			},
			returnBase: true,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {

			err := testCase.validator.Validate(testCase.input)
			if testCase.returnBase {
				assert.Equal(t, &testCase.validator.ValidatorBase, err)
			}
		})
	}
}
