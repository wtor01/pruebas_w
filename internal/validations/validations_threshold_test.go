package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type InputUmbral struct {
	measure   MeasureValidatable
	validator ValidatorThreshold
}

type WantNewValidator struct {
	err   bool
	value ValidatorThreshold
}

func Test_Unit_Domain_Validations_NewValidatorThreshold(t *testing.T) {
	type input struct {
		ValidationData
		measures.Status
	}
	tests := map[string]struct {
		input input
		want  WantNewValidator
	}{
		"Not Valid ValidationData Config Lenght CASE": {
			input: input{},
			want:  WantNewValidator{err: true, value: ValidatorThreshold{}},
		},
		"Not Valid Key Case": {
			input: input{
				ValidationData: ValidationData{
					Type: Threshold,
					Keys: []string{"NotValidKey"},
					Config: map[string]string{
						"RandomKey1": "RandomValue",
						"RandomKey2": "RandomValue",
					},
				},
			},
			want: WantNewValidator{err: true, value: ValidatorThreshold{}},
		},
		"Cant Parse MaxKey in Config CASE": {
			input: input{
				ValidationData: ValidationData{
					Type: Threshold,
					Keys: []string{AI},
					Config: map[string]string{
						MaxKey: "Cant Parse This To float64",
						MinKey: "64.0",
					},
				},
			},
			want: WantNewValidator{err: true, value: ValidatorThreshold{}},
		},
		"Cant Parse Min in Config CASE": {
			input: input{
				ValidationData: ValidationData{
					Type: Threshold,
					Keys: []string{AI},
					Config: map[string]string{
						MaxKey: "65.0",
						MinKey: "Cant Parse This To float64",
					},
				},
			},
			want: WantNewValidator{err: true, value: ValidatorThreshold{}},
		},
		"Should Return Ranged Error CASE": {
			input: input{
				ValidationData: ValidationData{
					Type: Threshold,
					Keys: []string{AI},
					Config: map[string]string{
						MaxKey: "65.0",
						MinKey: "80.0",
					},
				},
			},
			want: WantNewValidator{err: true, value: ValidatorThreshold{}},
		},
		"Should NOT Return Error CASE": {
			input: input{
				ValidationData: ValidationData{
					Type: Threshold,
					Keys: []string{AI},
					Config: map[string]string{
						MinKey: "65.0",
						MaxKey: "80.0",
					},
				},
			},
			want: WantNewValidator{err: false, value: ValidatorThreshold{
				Min: 65,
				Max: 80,
				ValidatorBase: ValidatorBase{
					Type: Threshold,
					Keys: []string{AI},
					Config: map[string]string{
						MinKey: "65.0",
						MaxKey: "80.0",
					},
				},
			}},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]

		t.Run(testName, func(t *testing.T) {

			validator, err := NewValidatorThreshold(testCase.input.ValidationData, testCase.input.Status)

			assert.Equal(t, testCase.want.value, validator, testName)
			if testCase.want.err {
				assert.NotNil(t, err, testName)
			} else {
				assert.Nil(t, err, testName)
			}
		})
	}
}

func Test_Unit_Domain_Threshold_Validator(t *testing.T) {
	tests := map[string]struct {
		input          InputUmbral
		want           error
		failValidation bool
	}{
		"should be nil if AI in threshold": {
			input: InputUmbral{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          9,
					WhiteListKeys: map[string]struct{}{
						AI: {},
						AE: {},
						R1: {},
						R2: {},
						R3: {},
						R4: {},
						MX: {},
						E:  {},
					},
				},
				validator: ValidatorThreshold{
					ValidatorBase: ValidatorBase{
						Type: QuarterlyDate,
						Keys: []string{StartDate, EndDate, MeasureDate, FX},
						Config: map[string]string{
							MaxKey: "10",
							MinKey: "0",
						},
						Action: measures.Invalid,
						Code:   QuarterlyDate,
					},
					Max: 10,
					Min: 0,
				},
			},
			failValidation: false,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {

			base := testCase.input.validator.Validate(testCase.input.measure)
			if testCase.failValidation {
				assert.Equal(t, &testCase.input.validator.ValidatorBase, base, testCase)
			} else {
				assert.Nil(t, base, testCase)
			}
		})
	}
}
