package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type InputsExcessiveConsumption struct {
	measure   MeasureValidatable
	validator ValidatorExcessiveConsumption
}

type WantNewValidatorExcessiveConsumption struct {
	err   bool
	value ValidatorExcessiveConsumption
}

func Test_Unit_Domain_Validations_NewValidatorExcessiveConsumption(t *testing.T) {
	type input struct {
		ValidationData
		measures.Status
	}
	tests := map[string]struct {
		input input
		want  WantNewValidatorExcessiveConsumption
	}{
		"Not Valid length": {
			input: input{
				ValidationData: ValidationData{
					Type: ExcessiveConsumption,
					Keys: []string{"NotValidKey"},
					Config: map[string]string{
						"RandomKey1": "RandomValue",
						"RandomKey2": "RandomValue",
					},
				},
			},
			want: WantNewValidatorExcessiveConsumption{err: true, value: ValidatorExcessiveConsumption{}},
		},

		"Should NOT Return Error CASE": {
			input: input{
				ValidationData: ValidationData{
					Type:   ExcessiveConsumption,
					Keys:   []string{AI, AE},
					Config: map[string]string{},
				},
			},
			want: WantNewValidatorExcessiveConsumption{err: false, value: ValidatorExcessiveConsumption{

				ValidatorBase: ValidatorBase{
					Type:   ExcessiveConsumption,
					Keys:   []string{AI, AE},
					Config: map[string]string{},
				},
			}},
		},
	}
	for testName := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {

			validator, err := NewValidatorExcessiveConsumption(testCase.input.ValidationData, testCase.input.Status)

			assert.Equal(t, testCase.want.value, validator, testName)
			if testCase.want.err {
				assert.NotNil(t, err, testName)
			} else {
				assert.Nil(t, err, testName)
			}
		})
	}
}

func Test_Unit_Domain_ExcessiveConsumption_Validator(t *testing.T) {
	tests := map[string]struct {
		input          InputsExcessiveConsumption
		want           error
		failValidation bool
	}{

		"should be nil in functions": {
			input: InputsExcessiveConsumption{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					Period:      measures.P1,
					P1: Periods{
						Values: measures.Values{AI: 10, AE: 10},
					},
					P1Demand:     10,
					RegisterType: measures.Hourly,
				},
				validator: ValidatorExcessiveConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ExcessiveConsumption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ExcessiveConsumption,
					},
				},
			},
			failValidation: false,
		},

		"should be valid in P1 AI Hourly": {
			input: InputsExcessiveConsumption{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					Period:      measures.P1,
					P1: Periods{
						Values: measures.Values{AI: 10, AE: 10},
					},
					P1Demand:     10,
					RegisterType: measures.Hourly,
				},
				validator: ValidatorExcessiveConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ExcessiveConsumption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ExcessiveConsumption,
					},
				},
			},
			failValidation: false,
		},
		"should be valid in P2 AE Hourly": {
			input: InputsExcessiveConsumption{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					Period:      measures.P2,
					P2: Periods{
						Values: measures.Values{AI: 10, AE: 10},
					},
					P2Demand:     11,
					RegisterType: measures.Hourly,
				},
				validator: ValidatorExcessiveConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ExcessiveConsumption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ExcessiveConsumption,
					},
				},
			},
			failValidation: false,
		},
		"should be valid in P3 AI QuarterHour": {
			input: InputsExcessiveConsumption{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					Period:      measures.P3,
					P3: Periods{
						Values: measures.Values{AI: 2, AE: 2},
					},
					P3Demand:     9,
					RegisterType: measures.QuarterHour,
				},
				validator: ValidatorExcessiveConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ExcessiveConsumption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ExcessiveConsumption,
					},
				},
			},
			failValidation: false,
		},
		"should be valid in P4 AI QuarterHour": {
			input: InputsExcessiveConsumption{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					Period:      measures.P4,
					P4: Periods{
						Values: measures.Values{AI: 2, AE: 2},
					},
					P4Demand:     9,
					RegisterType: measures.QuarterHour,
				},
				validator: ValidatorExcessiveConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ExcessiveConsumption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ExcessiveConsumption,
					},
				},
			},
			failValidation: false,
		},

		"should be invalid in P5 AI Hourly": {
			input: InputsExcessiveConsumption{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					Period:      measures.P5,
					P5: Periods{
						Values: measures.Values{AI: 10, AE: 10},
					},
					P5Demand:     9,
					RegisterType: measures.Hourly,
				},
				validator: ValidatorExcessiveConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ExcessiveConsumption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ExcessiveConsumption,
					},
				},
			},
			failValidation: true,
		},
		"should be invalid in P6 AE Hourly": {
			input: InputsExcessiveConsumption{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					Period:      measures.P6,
					P6: Periods{
						Values: measures.Values{AI: 10, AE: 10},
					},
					P6Demand:     9,
					RegisterType: measures.Hourly,
				},
				validator: ValidatorExcessiveConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ExcessiveConsumption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ExcessiveConsumption,
					},
				},
			},
			failValidation: true,
		},
	}
	for testName := range tests {
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
