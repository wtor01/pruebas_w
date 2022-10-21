package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type InputCloseMeter struct {
	measure   MeasureValidatable
	validator ValidatorCloseMeter
}

type WantNewValidatorCloseMeter struct {
	err   bool
	value ValidatorCloseMeter
}

func Test_Unit_Domain_Validations_NewValidatorCloseMeter(t *testing.T) {
	type input struct {
		ValidationData
		measures.Status
	}
	tests := map[string]struct {
		input input

		want WantNewValidatorCloseMeter
	}{

		"Should NOT Return Error CASE": {
			input: input{
				ValidationData: ValidationData{
					Type:   SummaryCalendar,
					Keys:   []string{AI},
					Config: map[string]string{},
				},
			},

			want: WantNewValidatorCloseMeter{err: false, value: ValidatorCloseMeter{

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

			validator, err := NewValidatorCloseMeter(testCase.input.ValidationData, testCase.input.Status)

			assert.Equal(t, testCase.want.value, validator, testName)
			if testCase.want.err {
				assert.NotNil(t, err, testName)
			} else {
				assert.Nil(t, err, testName)
			}
		})
	}
}

func Test_Unit_Domain_CloseMeter_Validator(t *testing.T) {
	tests := map[string]struct {
		input          InputCloseMeter
		want           error
		failValidation bool
	}{
		"should be error P0": {

			input: InputCloseMeter{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P0: Periods{
						Values: measures.Values{AI: 0, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
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
					P0LM: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
					},
					P1LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
				},
				validator: ValidatorCloseMeter{
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
		"should be error P1": {

			input: InputCloseMeter{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P0: Periods{
						Values: measures.Values{AI: 12, AE: 12, R1: 12, R2: 12, R3: 12, R4: 12},
					},
					P1: Periods{
						Values: measures.Values{AI: 0, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
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
					P0LM: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
					},
					P1LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
				},
				validator: ValidatorCloseMeter{
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

			input: InputCloseMeter{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 0},
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
					P0LM: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
					},
					P1LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
				},
				validator: ValidatorCloseMeter{
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

			input: InputCloseMeter{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
					},
					P1: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 0, R4: 1},
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
					P0LM: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
					},
					P1LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
				},
				validator: ValidatorCloseMeter{
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

			input: InputCloseMeter{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
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
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 0, R3: 1, R4: 1},
					},
					P5: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P0LM: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
					},
					P1LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
				},
				validator: ValidatorCloseMeter{
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

			input: InputCloseMeter{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
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
						Values: measures.Values{AI: 1, AE: 0, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P0LM: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
					},
					P1LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
				},
				validator: ValidatorCloseMeter{
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

			input: InputCloseMeter{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
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
						Values: measures.Values{AI: 1, AE: 1, R1: 0, R2: 1, R3: 1, R4: 1},
					},
					P0LM: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
					},
					P1LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
				},
				validator: ValidatorCloseMeter{
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

			input: InputCloseMeter{
				measure: MeasureValidatable{
					EndDate:     time.Date(2016, 8, 30, 10, 0, 0, 0, time.UTC),
					ReadingDate: time.Date(2016, 8, 31, 8, 30, 05, 0, time.UTC),
					AI:          90,
					ServiceType: measures.DcServiceType,
					PointType:   "1",
					WhiteListKeys: map[string]struct{}{
						AI: {},
					},
					P0: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
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
					P0LM: Periods{
						Values: measures.Values{AI: 6, AE: 6, R1: 6, R2: 6, R3: 6, R4: 6},
					},
					P1LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P2LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P3LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P4LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P5LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
					P6LM: Periods{
						Values: measures.Values{AI: 1, AE: 1, R1: 1, R2: 1, R3: 1, R4: 1},
					},
				},
				validator: ValidatorCloseMeter{
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
