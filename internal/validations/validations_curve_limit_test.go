package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type InputCurveLimit struct {
	measure   MeasureValidatable
	validator ValidatorCurveLimit
}

type WantNewValidatorCurveLimit struct {
	err   bool
	value ValidatorCurveLimit
}

func Test_Unit_Domain_Validations_NewValidatorCurveLimit(t *testing.T) {
	type input struct {
		ValidationData
		measures.Status
	}
	tests := map[string]struct {
		input input

		want WantNewValidatorCurveLimit
	}{
		"Not Valid ValidationData Config Lenght CASE": {
			input: input{},
			want:  WantNewValidatorCurveLimit{err: true, value: ValidatorCurveLimit{}},
		},
		"Not Valid Key Case": {
			input: input{
				ValidationData: ValidationData{
					Type: CurveLimit,
					Keys: []string{"NotValidKey"},
					Config: map[string]string{
						"RandomKey1": "RandomValue",
						"RandomKey2": "RandomValue",
					},
				},
			},

			want: WantNewValidatorCurveLimit{err: true, value: ValidatorCurveLimit{}},
		},
		"Cant Parse type1 in Config CASE": {
			input: input{
				ValidationData: ValidationData{
					Type: CurveLimit,
					Keys: []string{AI},
					Config: map[string]string{
						Type1: "Cant Parse This To float64",
						Type2: "10000",
						Type3: "10000",
						Type4: "10000",
						Type5: "10000",
					},
				},
			},

			want: WantNewValidatorCurveLimit{err: true, value: ValidatorCurveLimit{}},
		},

		"Should NOT Return Error CASE": {
			input: input{
				ValidationData: ValidationData{
					Type: CurveLimit,
					Keys: []string{AI},
					Config: map[string]string{
						Type1: "100000",
						Type2: "10000",
						Type3: "1000",
						Type4: "100",
						Type5: "10",
					},
				},
			},

			want: WantNewValidatorCurveLimit{err: false, value: ValidatorCurveLimit{

				Type1: 100000,
				Type2: 10000,
				Type3: 1000,
				Type4: 100,
				Type5: 10,

				ValidatorBase: ValidatorBase{
					Type: CurveLimit,
					Keys: []string{AI},
					Config: map[string]string{
						Type1: "100000",
						Type2: "10000",
						Type3: "1000",
						Type4: "100",
						Type5: "10",
					},
				},
			}},
		},
	}
	for testName := range tests {
		testCase := tests[testName]

		t.Run(testName, func(t *testing.T) {

			validator, err := NewValidatorCurveLimit(testCase.input.ValidationData, testCase.input.Status)

			assert.Equal(t, testCase.want.value, validator, testName)
			if testCase.want.err {
				assert.NotNil(t, err, testName)
			} else {
				assert.Nil(t, err, testName)
			}
		})
	}
}

func Test_Unit_Domain_CurveLimit_Validator(t *testing.T) {
	tests := map[string]struct {
		input          InputCurveLimit
		want           error
		failValidation bool
	}{
		"should be nil if functions": {
			input: InputCurveLimit{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorCurveLimit{
					ValidatorBase: ValidatorBase{
						Type: CurveLimit,
						Keys: []string{AI},
						Config: map[string]string{
							Type1: "100000",
							Type2: "10000",
							Type3: "1000",
							Type4: "100",
							Type5: "10",
						},
						Action: measures.Invalid,
						Code:   CurveLimit,
					},
					Type1: 100000,
					Type2: 10000,
					Type3: 1000,
					Type4: 100,
					Type5: 10,
				},
			},
			failValidation: false,
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
