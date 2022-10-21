package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Unit_Domain_ValidatorFuture_New(t *testing.T) {
	tests := map[string]struct {
		want struct {
			err   bool
			value ValidatorFuture
		}
		input struct {
			v      ValidationData
			status measures.Status
		}
	}{
		"Should return err with invalid config": {
			want: struct {
				err   bool
				value ValidatorFuture
			}{
				err:   true,
				value: ValidatorFuture{},
			},
			input: struct {
				v      ValidationData
				status measures.Status
			}{
				v: ValidationData{
					Type:   FutureDate,
					Keys:   []string{StartDate},
					Config: map[string]string{"test": "test"},
				},
				status: measures.Invalid,
			},
		},
		"Should return err with invalid keys": {
			want: struct {
				err   bool
				value ValidatorFuture
			}{
				err:   true,
				value: ValidatorFuture{},
			},
			input: struct {
				v      ValidationData
				status measures.Status
			}{
				v: ValidationData{
					Type:   FutureDate,
					Keys:   []string{AI},
					Config: map[string]string{},
				},
				status: measures.Invalid,
			},
		},
		"Should return valid ValidatorFuture": {
			want: struct {
				err   bool
				value ValidatorFuture
			}{
				err: false,
				value: ValidatorFuture{
					ValidatorBase: ValidatorBase{
						Type:   FutureDate,
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
					Type:   FutureDate,
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

			validator, err := NewValidatorFuture(testCase.input.v, testCase.input.status)

			assert.Equal(t, testCase.want.value, validator, testName)
			if testCase.want.err {
				assert.NotNil(t, err, testName)
			} else {
				assert.Nil(t, err, testName)
			}
		})
	}
}

func Test_Unit_Domain_ValidatorFuture_Validator(t *testing.T) {

	futureDate := time.Now()
	futureDate = futureDate.Add(time.Hour + 1)

	oldDate := time.Now()
	oldDate = oldDate.Add(-time.Second * 1)

	tests := map[string]struct {
		input     MeasureValidatable
		validator ValidatorFuture
	}{
		"should update measure status if invalid date in EndDate to measures.Invalid": {
			input: MeasureValidatable{
				StartDate:   oldDate,
				EndDate:     futureDate,
				ReadingDate: oldDate,

				FX: oldDate,
			},
			validator: ValidatorFuture{
				ValidatorBase: ValidatorBase{
					Type:   FutureDate,
					Keys:   []string{StartDate, EndDate, MeasureDate, FX},
					Config: map[string]string{},
					Action: measures.Invalid,
					Code:   FutureDate,
				},
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {

			err := testCase.validator.Validate(testCase.input)
			assert.Nil(t, err, testName)
		})
	}
}
