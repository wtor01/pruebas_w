package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type InputSComZer struct {
	measure   MeasureValidatable
	validator ValidatorZeroConsumption
}

type WantNewValidatorZeroConsumption struct {
	err   bool
	value ValidatorZeroConsumption
}

func Test_Unit_Domain_Validations_NewValidatorZeroConsumption(t *testing.T) {
	type input struct {
		ValidationData
		measures.Status
	}
	tests := map[string]struct {
		input input
		want  WantNewValidatorZeroConsumption
	}{
		"Not Valid length": {
			input: input{
				ValidationData: ValidationData{
					Type: ZeroConsmuption,
					Keys: []string{"NotValidKey"},
					Config: map[string]string{
						"RandomKey1": "RandomValue",
						"RandomKey2": "RandomValue",
					},
				},
			},
			want: WantNewValidatorZeroConsumption{err: true, value: ValidatorZeroConsumption{}},
		},

		"Should NOT Return Error CASE": {
			input: input{
				ValidationData: ValidationData{
					Type:   ZeroConsmuption,
					Keys:   []string{AI},
					Config: map[string]string{},
				},
			},
			want: WantNewValidatorZeroConsumption{err: false, value: ValidatorZeroConsumption{

				ValidatorBase: ValidatorBase{
					Type:   ZeroConsmuption,
					Keys:   []string{AI},
					Config: map[string]string{},
				},
			}},
		},
	}

	for testName := range tests {
		testCase := tests[testName]

		t.Run(testName, func(t *testing.T) {

			validator, err := NewValidatorZeroConsumption(testCase.input.ValidationData, testCase.input.Status)

			assert.Equal(t, testCase.want.value, validator, testName)
			if testCase.want.err {
				assert.NotNil(t, err, testName)
			} else {
				assert.Nil(t, err, testName)
			}
		})
	}

}

func Test_Unit_Domain_ZeroConsumption_Validator(t *testing.T) {
	tests := map[string]struct {
		input          InputSComZer
		want           error
		failValidation bool
	}{

		"no TLG": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					ServiceType:    measures.DcServiceType,
					TechnologyType: "TLM",
					PointType:      "3",

					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: false,
		},
		"no ServiceType": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					ServiceType:    measures.DdServiceType,
					TechnologyType: "TLM",
					PointType:      "3",

					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: false,
		},
		"no PointType": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "2",

					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   SummaryTotalizer,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   SummaryTotalizer,
					},
				},
			},
			failValidation: false,
		},

		"should be nil if functions": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: false,
		},
		"should be error AI p1 period": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 0, AE: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: true,
		},
		"should be error AI p2 period": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 0, AE: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: true,
		},
		"should be error AI p3 period": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 0, AE: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: true,
		},
		"should be error AI p4 period": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 0, AE: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: true,
		},
		"should be error AI p5 period": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 0, AE: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: true,
		},
		"should be error AI p6 period": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 0, AE: 1},
					},
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: true,
		},

		"should be error AE p1 period": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 0},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: true,
		},
		"should be error AE p2 period": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 0},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: true,
		},
		"should be error AE p3 period": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 0},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: true,
		},
		"should be error AE p4 period": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 0},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: true,
		},
		"should be error AE p5 period": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 0},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
					},
				},
			},
			failValidation: true,
		},
		"should be error AE p6 period": {
			input: InputSComZer{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 0},
					},
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorZeroConsumption{
					ValidatorBase: ValidatorBase{
						Type:   ZeroConsmuption,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   ZeroConsmuption,
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
