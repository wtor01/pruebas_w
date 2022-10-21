package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Unit_Domain_ValidatorHourly_New(t *testing.T) {
	tests := map[string]struct {
		want struct {
			err   bool
			value ValidatorHourly
		}
		input struct {
			v      ValidationData
			status measures.Status
		}
	}{
		"Should return err with invalid config": {
			want: struct {
				err   bool
				value ValidatorHourly
			}{
				err:   true,
				value: ValidatorHourly{},
			},
			input: struct {
				v      ValidationData
				status measures.Status
			}{
				v: ValidationData{
					Type:   HourDate,
					Keys:   []string{StartDate},
					Config: map[string]string{"test": "test"},
				},
				status: measures.Invalid,
			},
		},
		"Should return err with invalid keys": {
			want: struct {
				err   bool
				value ValidatorHourly
			}{
				err:   true,
				value: ValidatorHourly{},
			},
			input: struct {
				v      ValidationData
				status measures.Status
			}{
				v: ValidationData{
					Type:   HourDate,
					Keys:   []string{AI},
					Config: map[string]string{},
				},
				status: measures.Invalid,
			},
		},
		"Should return valid ValidatorHour": {
			want: struct {
				err   bool
				value ValidatorHourly
			}{
				err: false,
				value: ValidatorHourly{
					ValidatorBase: ValidatorBase{
						Type:   HourDate,
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
					Type:   HourDate,
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

			validator, err := NewValidatorHourly(testCase.input.v, testCase.input.status)

			assert.Equal(t, testCase.want.value, validator, testName)
			if testCase.want.err {
				assert.NotNil(t, err, testName)
			} else {
				assert.Nil(t, err, testName)
			}
		})
	}
}

func Test_Unit_Domain_ValidatorHour_Validator(t *testing.T) {

	today := time.Now()
	invalidDate := time.Date(today.Year(), today.Month(), today.Day(), today.Hour(), 2, 0, 0, time.Local)
	validDate := time.Date(today.Year(), today.Month(), today.Day(), today.Hour(), 0, 0, 0, time.Local)

	tests := map[string]struct {
		input     MeasureValidatable
		validator ValidatorHourly
	}{
		"should update measure status if invalid date in EndDate to measures.Invalid": {
			input: MeasureValidatable{
				StartDate:   validDate,
				EndDate:     invalidDate,
				ReadingDate: validDate,
				FX:          validDate,
			},
			validator: ValidatorHourly{
				ValidatorBase: ValidatorBase{
					Type:   HourDate,
					Keys:   []string{StartDate, EndDate, MeasureDate, FX},
					Config: map[string]string{},
					Action: measures.Invalid,
					Code:   HourDate,
				},
			},
		},
		"should update measure status if invalid date in EndDate to measures.Alert": {
			input: MeasureValidatable{
				StartDate:   validDate,
				EndDate:     invalidDate,
				ReadingDate: validDate,
				FX:          validDate,
			},
			validator: ValidatorHourly{
				ValidatorBase: ValidatorBase{
					Type:   HourDate,
					Keys:   []string{StartDate, EndDate, MeasureDate, FX},
					Config: map[string]string{},
					Action: measures.Alert,
					Code:   HourDate,
				},
			},
		},
		"should update measure status if invalid date in StartDate to measures.Supervise": {
			input: MeasureValidatable{
				StartDate:   invalidDate,
				EndDate:     validDate,
				ReadingDate: validDate,
				FX:          validDate,
			},
			validator: ValidatorHourly{
				ValidatorBase: ValidatorBase{
					Type:   HourDate,
					Keys:   []string{StartDate, EndDate, MeasureDate, FX},
					Config: map[string]string{},
					Action: measures.Supervise,
					Code:   HourDate,
				},
			},
		},
		"should update measure status if invalid date in ReadingDate to measures.Invalid": {
			input: MeasureValidatable{
				StartDate:   validDate,
				EndDate:     validDate,
				ReadingDate: invalidDate,
				FX:          validDate,
			},
			validator: ValidatorHourly{
				ValidatorBase: ValidatorBase{
					Type:   HourDate,
					Keys:   []string{StartDate, EndDate, MeasureDate, FX},
					Config: map[string]string{},
					Action: measures.Invalid,
					Code:   HourDate,
				},
			},
		},
		"should update measure status if invalid date in FX to measures.Invalid": {
			input: MeasureValidatable{
				StartDate:   validDate,
				EndDate:     validDate,
				ReadingDate: validDate,
				FX:          invalidDate,
			},
			validator: ValidatorHourly{
				ValidatorBase: ValidatorBase{
					Type:   HourDate,
					Keys:   []string{StartDate, EndDate, MeasureDate, FX},
					Config: map[string]string{},
					Action: measures.Invalid,
					Code:   HourDate,
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
