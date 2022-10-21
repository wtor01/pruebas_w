package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type InputUmbralSummaryCalendar struct {
	measure   MeasureValidatable
	validator ValidatorSummaryCalendar
}

type WantNewValidatorSummaryCalendar struct {
	err   bool
	value ValidatorSummaryCalendar
}

func Test_Unit_Domain_Validations_NewValidatorSummaryCalendar(t *testing.T) {
	type input struct {
		ValidationData
		measures.Status
	}
	tests := map[string]struct {
		input input

		want WantNewValidatorSummaryCalendar
	}{

		"Not Valid Key Case": {
			input: input{
				ValidationData: ValidationData{
					Type:   SummaryCalendar,
					Keys:   []string{"NotValidKey"},
					Config: map[string]string{},
				},
			},

			want: WantNewValidatorSummaryCalendar{err: true, value: ValidatorSummaryCalendar{}},
		},

		"Should NOT Return Error CASE": {
			input: input{
				ValidationData: ValidationData{
					Type:   SummaryCalendar,
					Keys:   []string{AI},
					Config: map[string]string{},
				},
			},

			want: WantNewValidatorSummaryCalendar{err: false, value: ValidatorSummaryCalendar{

				ValidatorBase: ValidatorBase{
					Type:   SummaryCalendar,
					Keys:   []string{AI},
					Config: map[string]string{},
				},
			}},
		},
	}
	for testName := range tests {
		testCase := tests[testName]

		t.Run(testName, func(t *testing.T) {

			validator, err := NewValidatorSummaryCalendar(testCase.input.ValidationData, testCase.input.Status)

			assert.Equal(t, testCase.want.value, validator, testName)
			if testCase.want.err {
				assert.NotNil(t, err, testName)
			} else {
				assert.Nil(t, err, testName)
			}
		})
	}
}

func Test_Unit_Domain_SummaryCalendar_Validator(t *testing.T) {
	tests := map[string]struct {
		input          InputUmbralSummaryCalendar
		want           error
		failValidation bool
	}{
		"should be error P1": {

			input: InputUmbralSummaryCalendar{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P1i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					PeriodsNumbers: []string{"P2"},
				},
				validator: ValidatorSummaryCalendar{
					ValidatorBase: ValidatorBase{
						Type:   SummaryCalendar,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   SummaryCalendar,
					},
				},
			},
			failValidation: true,
		},
		"should be error P2": {

			input: InputUmbralSummaryCalendar{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P1i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2i: Periods{
						Values: measures.Values{AI: 0, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					PeriodsNumbers: []string{"P1"},
				},
				validator: ValidatorSummaryCalendar{
					ValidatorBase: ValidatorBase{
						Type:   SummaryCalendar,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   SummaryCalendar,
					},
				},
			},
			failValidation: true,
		},
		"should be error P3": {

			input: InputUmbralSummaryCalendar{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P1i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3i: Periods{
						Values: measures.Values{AI: 0, AE: 0, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					PeriodsNumbers: []string{"P1", "P2"},
				},
				validator: ValidatorSummaryCalendar{
					ValidatorBase: ValidatorBase{
						Type:   SummaryCalendar,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   SummaryCalendar,
					},
				},
			},
			failValidation: true,
		},
		"should be error P4": {

			input: InputUmbralSummaryCalendar{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P1i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4i: Periods{
						Values: measures.Values{AI: 0, AE: 0, R1: 0, R2: 1, R3: 1, R4: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					PeriodsNumbers: []string{"P1", "P2", "P3"},
				},
				validator: ValidatorSummaryCalendar{
					ValidatorBase: ValidatorBase{
						Type:   SummaryCalendar,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   SummaryCalendar,
					},
				},
			},
			failValidation: true,
		},
		"should be error P5": {

			input: InputUmbralSummaryCalendar{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P1i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5i: Periods{
						Values: measures.Values{AI: 0, AE: 0, R1: 0, R2: 0, R3: 1, R4: 1},
					},
					P6i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					PeriodsNumbers: []string{"P1", "P2", "P3", "P4"},
				},
				validator: ValidatorSummaryCalendar{
					ValidatorBase: ValidatorBase{
						Type:   SummaryCalendar,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   SummaryCalendar,
					},
				},
			},
			failValidation: true,
		},
		"should be error P6": {

			input: InputUmbralSummaryCalendar{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P1i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6i: Periods{
						Values: measures.Values{AI: 0, AE: 0, R1: 0, R2: 0, R3: 0, R4: 1},
					},
					PeriodsNumbers: []string{"P1", "P2", "P3", "P4", "P5"},
				},
				validator: ValidatorSummaryCalendar{
					ValidatorBase: ValidatorBase{
						Type:   SummaryCalendar,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   SummaryCalendar,
					},
				},
			},
			failValidation: true,
		},
		"should be nil if functions": {

			input: InputUmbralSummaryCalendar{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P1i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6i: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					PeriodsNumbers: []string{"P1", "P2", "P3", "P4", "P5", "P6"},
				},
				validator: ValidatorSummaryCalendar{
					ValidatorBase: ValidatorBase{
						Type:   SummaryCalendar,
						Keys:   []string{AI},
						Config: map[string]string{},
						Action: measures.Invalid,
						Code:   SummaryCalendar,
					},
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
