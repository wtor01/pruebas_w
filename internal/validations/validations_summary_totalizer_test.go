package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type InputSumTot struct {
	measure   MeasureValidatable
	validator ValidatorSummaryTotalizer
}

type WantNewValidatorSummaryTotalizer struct {
	err   bool
	value ValidatorSummaryTotalizer
}

func Test_Unit_Domain_Validations_NewValidatorSummaryTotalizer(t *testing.T) {
	type input struct {
		ValidationData
		measures.Status
	}
	tests := map[string]struct {
		input input
		want  WantNewValidatorSummaryTotalizer
	}{
		"Not Valid length": {
			input: input{
				ValidationData: ValidationData{
					Type: SummaryTotalizer,
					Keys: []string{"NotValidKey"},
					Config: map[string]string{
						"RandomKey1": "RandomValue",
						"RandomKey2": "RandomValue",
					},
				},
			},
			want: WantNewValidatorSummaryTotalizer{err: true, value: ValidatorSummaryTotalizer{}},
		},
		"Not Valid parse": {
			input: input{
				ValidationData: ValidationData{
					Type: SummaryTotalizer,
					Keys: []string{"NotValidKey"},
					Config: map[string]string{
						Tolerance: "no parse",
					},
				},
			},
			want: WantNewValidatorSummaryTotalizer{err: true, value: ValidatorSummaryTotalizer{}},
		},

		"Should NOT Return Error CASE": {
			input: input{
				ValidationData: ValidationData{
					Type: SummaryTotalizer,
					Keys: []string{AI},
					Config: map[string]string{
						Tolerance: "1",
					},
				},
			},
			want: WantNewValidatorSummaryTotalizer{err: false, value: ValidatorSummaryTotalizer{
				Tolerance: 1,
				ValidatorBase: ValidatorBase{
					Type: SummaryTotalizer,
					Keys: []string{AI},
					Config: map[string]string{
						Tolerance: "1",
					},
				},
			}},
		},
	}
	for testName := range tests {
		testCase := tests[testName]

		t.Run(testName, func(t *testing.T) {

			validator, err := NewValidatorSummaryTotalizer(testCase.input.ValidationData, testCase.input.Status)

			assert.Equal(t, testCase.want.value, validator, testName)
			if testCase.want.err {
				assert.NotNil(t, err, testName)
			} else {
				assert.Nil(t, err, testName)
			}
		})
	}
}

func Test_Unit_Domain_SummaryTotalizer_Validator(t *testing.T) {
	tests := map[string]struct {
		input          InputSumTot
		want           error
		failValidation bool
	}{
		"should be nil if functions": {
			input: InputSumTot{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 5, AE: 7, R1: 6, R2: 6, R3: 6, R4: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},

					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorSummaryTotalizer{
					Tolerance: 1,
					ValidatorBase: ValidatorBase{
						Type: SummaryTotalizer,
						Keys: []string{AI},
						Config: map[string]string{
							Tolerance: "1",
						},
						Action: measures.Invalid,
						Code:   SummaryTotalizer,
					},
				},
			},
			failValidation: false,
		},
		"no TLG": {
			input: InputSumTot{
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
				validator: ValidatorSummaryTotalizer{
					ValidatorBase: ValidatorBase{
						Type: SummaryTotalizer,
						Keys: []string{AI},
						Config: map[string]string{
							Tolerance: "1",
						},

						Action: measures.Invalid,
						Code:   SummaryTotalizer,
					},
					Tolerance: 1,
				},
			},
			failValidation: false,
		},
		"no ServiceType": {
			input: InputSumTot{
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
				validator: ValidatorSummaryTotalizer{
					ValidatorBase: ValidatorBase{
						Type: SummaryTotalizer,
						Keys: []string{AI},
						Config: map[string]string{
							Tolerance: "1",
						},
						Action: measures.Invalid,
						Code:   SummaryTotalizer,
					},
					Tolerance: 1,
				},
			},
			failValidation: false,
		},
		"no PointType": {
			input: InputSumTot{
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
				validator: ValidatorSummaryTotalizer{
					ValidatorBase: ValidatorBase{
						Type: SummaryTotalizer,
						Keys: []string{AI},
						Config: map[string]string{
							Tolerance: "1",
						},
						Action: measures.Invalid,
						Code:   SummaryTotalizer,
					},
					Tolerance: 1,
				},
			},
			failValidation: false,
		},
		"should be error AI summatory": {
			input: InputSumTot{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 4, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},

					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorSummaryTotalizer{
					ValidatorBase: ValidatorBase{
						Type: SummaryTotalizer,
						Keys: []string{AI},
						Config: map[string]string{
							Tolerance: "1",
						},
						Action: measures.Invalid,
						Code:   SummaryTotalizer,
					},
					Tolerance: 1,
				},
			},
			failValidation: true,
		},
		"should be error AE summatory": {
			input: InputSumTot{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 7, AE: 8, R1: 6, R2: 6, R3: 6, R4: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},

					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorSummaryTotalizer{
					ValidatorBase: ValidatorBase{
						Type: SummaryTotalizer,
						Keys: []string{AI},
						Config: map[string]string{
							Tolerance: "1",
						},
						Action: measures.Invalid,
						Code:   SummaryTotalizer,
					},
					Tolerance: 1,
				},
			},
			failValidation: true,
		},
		"should be error R1 summatory": {
			input: InputSumTot{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 9, R2: 6, R3: 6, R4: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},

					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorSummaryTotalizer{
					ValidatorBase: ValidatorBase{
						Type: SummaryTotalizer,
						Keys: []string{AI},
						Config: map[string]string{
							Tolerance: "1",
						},
						Action: measures.Invalid,
						Code:   SummaryTotalizer,
					},
					Tolerance: 1,
				},
			},
			failValidation: true,
		},
		"should be error R2 summatory": {
			input: InputSumTot{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 3, R3: 6, R4: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},

					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorSummaryTotalizer{
					ValidatorBase: ValidatorBase{
						Type: SummaryTotalizer,
						Keys: []string{AI},
						Config: map[string]string{
							Tolerance: "1",
						},
						Action: measures.Invalid,
						Code:   SummaryTotalizer,
					},
					Tolerance: 1,
				},
			},
			failValidation: true,
		},
		"should be error R3 summatory": {
			input: InputSumTot{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 4, R4: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},

					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorSummaryTotalizer{
					ValidatorBase: ValidatorBase{
						Type: SummaryTotalizer,
						Keys: []string{AI},
						Config: map[string]string{
							Tolerance: "1",
						},
						Action: measures.Invalid,
						Code:   SummaryTotalizer,
					},
					Tolerance: 1,
				},
			},
			failValidation: true,
		},
		"should be error R4 summatory": {
			input: InputSumTot{
				measure: MeasureValidatable{
					EndDate:        time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate:    time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:             90,
					TechnologyType: "TLG",
					ServiceType:    measures.DcServiceType,
					PointType:      "3",

					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 8},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},

					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
				},
				validator: ValidatorSummaryTotalizer{
					ValidatorBase: ValidatorBase{
						Type: SummaryTotalizer,
						Keys: []string{AI},
						Config: map[string]string{
							Tolerance: "1",
						},
						Action: measures.Invalid,
						Code:   SummaryTotalizer,
					},
					Tolerance: 1,
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
