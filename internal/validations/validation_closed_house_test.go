package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
)

type InputsClosedHouse struct {
	measure   MeasureValidatable
	validator ValidatorClosedHouse
}

type WantNewValidatorClosedHouse struct {
	err   bool
	value ValidatorClosedHouse
}

func Test_Unit_Domain_Validations_NewValidatorClosedHouse(t *testing.T) {
	type input struct {
		ValidationData
		measures.Status
	}
	tests := map[string]struct {
		input input
		want  WantNewValidatorClosedHouse
	}{
		"Not Valid length": {
			input: input{
				ValidationData: ValidationData{
					Type: ClosedHouse,
					Keys: []string{"NotValidKey"},
					Config: map[string]string{
						"RandomKey1": "RandomValue",
						"RandomKey2": "RandomValue",
					},
				},
			},
			want: WantNewValidatorClosedHouse{err: true, value: ValidatorClosedHouse{}},
		},

		"Should NOT Return Error CASE": {
			input: input{
				ValidationData: ValidationData{
					Type:   ClosedHouse,
					Keys:   []string{AI, AE, R1, R2, R3, R4},
					Config: map[string]string{},
				},
			},
			want: WantNewValidatorClosedHouse{err: false, value: ValidatorClosedHouse{

				ValidatorBase: ValidatorBase{
					Type:   ClosedHouse,
					Keys:   []string{AI, AE, R1, R2, R3, R4},
					Config: map[string]string{},
				},
			}},
		},
	}

	for testName := range tests {
		testCase := tests[testName]

		t.Run(testName, func(t *testing.T) {

			validator, err := NewValidatorClosedHouse(testCase.input.ValidationData, testCase.input.Status)

			assert.Equal(t, testCase.want.value, validator, testName)
			if testCase.want.err {
				assert.NotNil(t, err, testName)
			} else {
				assert.Nil(t, err, testName)
			}
		})
	}

}

func Test_Unit_Domain_ClosedHouse_Validator(t *testing.T) {
	tests := map[string]struct {
		input          InputsClosedHouse
		want           error
		failValidation bool
	}{

		"Is not tecnhology validate": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "TLG",
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: false,
		},

		"There isn't power": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{AI: 0, AE: 0, R1: 0, R2: 0, R3: 0, R4: 0},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: false,
		},
		"There is power in P1 AI": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{AI: 2},
					},
					P1: Periods{
						Values: measures.Values{AI: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P2 AI": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{AI: 2},
					},
					P1: Periods{
						Values: measures.Values{AI: 0},
					},
					P2: Periods{
						Values: measures.Values{AI: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P3 AI": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{AI: 2},
					},
					P1: Periods{
						Values: measures.Values{AI: 0},
					},
					P2: Periods{
						Values: measures.Values{AI: 0},
					},
					P3: Periods{
						Values: measures.Values{AI: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P4 AI": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{AI: 2},
					},
					P1: Periods{
						Values: measures.Values{AI: 0},
					},
					P2: Periods{
						Values: measures.Values{AI: 0},
					},
					P3: Periods{
						Values: measures.Values{AI: 0},
					},
					P4: Periods{
						Values: measures.Values{AI: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P5 AI": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{AI: 2},
					},
					P1: Periods{
						Values: measures.Values{AI: 0},
					},
					P2: Periods{
						Values: measures.Values{AI: 0},
					},
					P3: Periods{
						Values: measures.Values{AI: 0},
					},
					P4: Periods{
						Values: measures.Values{AI: 0},
					},
					P5: Periods{
						Values: measures.Values{AI: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P6 AI": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{AI: 2},
					},
					P1: Periods{
						Values: measures.Values{AI: 0},
					},
					P2: Periods{
						Values: measures.Values{AI: 0},
					},
					P3: Periods{
						Values: measures.Values{AI: 0},
					},
					P4: Periods{
						Values: measures.Values{AI: 0},
					},
					P5: Periods{
						Values: measures.Values{AI: 0},
					},
					P6: Periods{
						Values: measures.Values{AI: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},

		"There is power in P1 AE": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{AE: 2},
					},
					P1: Periods{
						Values: measures.Values{AE: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P2 AE": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{AE: 2},
					},
					P1: Periods{
						Values: measures.Values{AE: 0},
					},
					P2: Periods{
						Values: measures.Values{AE: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P3 AE": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{AE: 2},
					},
					P1: Periods{
						Values: measures.Values{AE: 0},
					},
					P2: Periods{
						Values: measures.Values{AE: 0},
					},
					P3: Periods{
						Values: measures.Values{AE: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P4 AE": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{AE: 2},
					},
					P1: Periods{
						Values: measures.Values{AE: 0},
					},
					P2: Periods{
						Values: measures.Values{AE: 0},
					},
					P3: Periods{
						Values: measures.Values{AE: 0},
					},
					P4: Periods{
						Values: measures.Values{AE: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P5 AE": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{AE: 2},
					},
					P1: Periods{
						Values: measures.Values{AE: 0},
					},
					P2: Periods{
						Values: measures.Values{AE: 0},
					},
					P3: Periods{
						Values: measures.Values{AE: 0},
					},
					P4: Periods{
						Values: measures.Values{AE: 0},
					},
					P5: Periods{
						Values: measures.Values{AE: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P6 AE": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{AE: 2},
					},
					P1: Periods{
						Values: measures.Values{AE: 0},
					},
					P2: Periods{
						Values: measures.Values{AE: 0},
					},
					P3: Periods{
						Values: measures.Values{AE: 0},
					},
					P4: Periods{
						Values: measures.Values{AE: 0},
					},
					P5: Periods{
						Values: measures.Values{AE: 0},
					},
					P6: Periods{
						Values: measures.Values{AE: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},

		"There is power in P1 R1": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R1: 2},
					},
					P1: Periods{
						Values: measures.Values{R1: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P2 R1": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R1: 2},
					},
					P1: Periods{
						Values: measures.Values{R1: 0},
					},
					P2: Periods{
						Values: measures.Values{R1: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P3 R1": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R1: 2},
					},
					P1: Periods{
						Values: measures.Values{R1: 0},
					},
					P2: Periods{
						Values: measures.Values{R1: 0},
					},
					P3: Periods{
						Values: measures.Values{R1: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P4 R1": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R1: 2},
					},
					P1: Periods{
						Values: measures.Values{R1: 0},
					},
					P2: Periods{
						Values: measures.Values{R1: 0},
					},
					P3: Periods{
						Values: measures.Values{R1: 0},
					},
					P4: Periods{
						Values: measures.Values{R1: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P5 R1": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R1: 2},
					},
					P1: Periods{
						Values: measures.Values{R1: 0},
					},
					P2: Periods{
						Values: measures.Values{R1: 0},
					},
					P3: Periods{
						Values: measures.Values{R1: 0},
					},
					P4: Periods{
						Values: measures.Values{R1: 0},
					},
					P5: Periods{
						Values: measures.Values{R1: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P6 R1": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R1: 2},
					},
					P1: Periods{
						Values: measures.Values{R1: 0},
					},
					P2: Periods{
						Values: measures.Values{R1: 0},
					},
					P3: Periods{
						Values: measures.Values{R1: 0},
					},
					P4: Periods{
						Values: measures.Values{R1: 0},
					},
					P5: Periods{
						Values: measures.Values{R1: 0},
					},
					P6: Periods{
						Values: measures.Values{R1: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},

		"There is power in P1 R2": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R2: 2},
					},
					P1: Periods{
						Values: measures.Values{R2: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P2 R2": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R2: 2},
					},
					P1: Periods{
						Values: measures.Values{R2: 0},
					},
					P2: Periods{
						Values: measures.Values{R2: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P3 R2": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R2: 2},
					},
					P1: Periods{
						Values: measures.Values{R2: 0},
					},
					P2: Periods{
						Values: measures.Values{R2: 0},
					},
					P3: Periods{
						Values: measures.Values{R2: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P4 R2": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R2: 2},
					},
					P1: Periods{
						Values: measures.Values{R2: 0},
					},
					P2: Periods{
						Values: measures.Values{R2: 0},
					},
					P3: Periods{
						Values: measures.Values{R2: 0},
					},
					P4: Periods{
						Values: measures.Values{R2: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P5 R2": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R2: 2},
					},
					P1: Periods{
						Values: measures.Values{R2: 0},
					},
					P2: Periods{
						Values: measures.Values{R2: 0},
					},
					P3: Periods{
						Values: measures.Values{R2: 0},
					},
					P4: Periods{
						Values: measures.Values{R2: 0},
					},
					P5: Periods{
						Values: measures.Values{R2: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P6 R2": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R2: 2},
					},
					P1: Periods{
						Values: measures.Values{R2: 0},
					},
					P2: Periods{
						Values: measures.Values{R2: 0},
					},
					P3: Periods{
						Values: measures.Values{R2: 0},
					},
					P4: Periods{
						Values: measures.Values{R2: 0},
					},
					P5: Periods{
						Values: measures.Values{R2: 0},
					},
					P6: Periods{
						Values: measures.Values{R2: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},

		"There is power in P1 R3": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R3: 2},
					},
					P1: Periods{
						Values: measures.Values{R3: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P2 R3": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R3: 2},
					},
					P1: Periods{
						Values: measures.Values{R3: 0},
					},
					P2: Periods{
						Values: measures.Values{R3: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P3 R3": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R3: 2},
					},
					P1: Periods{
						Values: measures.Values{R3: 0},
					},
					P2: Periods{
						Values: measures.Values{R3: 0},
					},
					P3: Periods{
						Values: measures.Values{R3: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P4 R3": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R3: 2},
					},
					P1: Periods{
						Values: measures.Values{R3: 0},
					},
					P2: Periods{
						Values: measures.Values{R3: 0},
					},
					P3: Periods{
						Values: measures.Values{R3: 0},
					},
					P4: Periods{
						Values: measures.Values{R3: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P5 R3": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R3: 2},
					},
					P1: Periods{
						Values: measures.Values{R3: 0},
					},
					P2: Periods{
						Values: measures.Values{R3: 0},
					},
					P3: Periods{
						Values: measures.Values{R3: 0},
					},
					P4: Periods{
						Values: measures.Values{R3: 0},
					},
					P5: Periods{
						Values: measures.Values{R3: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P6 R3": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R3: 2},
					},
					P1: Periods{
						Values: measures.Values{R3: 0},
					},
					P2: Periods{
						Values: measures.Values{R3: 0},
					},
					P3: Periods{
						Values: measures.Values{R3: 0},
					},
					P4: Periods{
						Values: measures.Values{R3: 0},
					},
					P5: Periods{
						Values: measures.Values{R3: 0},
					},
					P6: Periods{
						Values: measures.Values{R3: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},

		"There is power in P1 R4": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R4: 2},
					},
					P1: Periods{
						Values: measures.Values{R4: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P2 R4": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R4: 2},
					},
					P1: Periods{
						Values: measures.Values{R4: 0},
					},
					P2: Periods{
						Values: measures.Values{R4: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P3 R4": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R4: 2},
					},
					P1: Periods{
						Values: measures.Values{R4: 0},
					},
					P2: Periods{
						Values: measures.Values{R4: 0},
					},
					P3: Periods{
						Values: measures.Values{R4: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P4 R4": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R4: 2},
					},
					P1: Periods{
						Values: measures.Values{R4: 0},
					},
					P2: Periods{
						Values: measures.Values{R4: 0},
					},
					P3: Periods{
						Values: measures.Values{R4: 0},
					},
					P4: Periods{
						Values: measures.Values{R4: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P5 R4": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R4: 2},
					},
					P1: Periods{
						Values: measures.Values{R4: 0},
					},
					P2: Periods{
						Values: measures.Values{R4: 0},
					},
					P3: Periods{
						Values: measures.Values{R4: 0},
					},
					P4: Periods{
						Values: measures.Values{R4: 0},
					},
					P5: Periods{
						Values: measures.Values{R4: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
					},
				},
			},
			failValidation: true,
		},
		"There is power in P6 R4": {
			input: InputsClosedHouse{
				measure: MeasureValidatable{
					TechnologyType: "INACCESIBLE",
					P0: Periods{
						Values: measures.Values{R4: 2},
					},
					P1: Periods{
						Values: measures.Values{R4: 0},
					},
					P2: Periods{
						Values: measures.Values{R4: 0},
					},
					P3: Periods{
						Values: measures.Values{R4: 0},
					},
					P4: Periods{
						Values: measures.Values{R4: 0},
					},
					P5: Periods{
						Values: measures.Values{R4: 0},
					},
					P6: Periods{
						Values: measures.Values{R4: 2},
					},
				},
				validator: ValidatorClosedHouse{
					ValidatorBase: ValidatorBase{
						Type:   ClosedHouse,
						Keys:   []string{AI, AE, R1, R2, R3, R4},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ClosedHouse,
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
