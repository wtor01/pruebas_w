package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
)

type InputQualifier struct {
	measure            MeasureValidatable
	ValidatorQualifier ValidatorQualifier
}

func Test_Unit_Domain_ValidatorQualifier_New(t *testing.T) {
	tests := map[string]struct {
		want struct {
			err   bool
			value ValidatorQualifier
		}
		input struct {
			v      ValidationData
			status measures.Status
		}
	}{
		"Should return err with invalid config": {
			want: struct {
				err   bool
				value ValidatorQualifier
			}{
				err:   true,
				value: ValidatorQualifier{},
			},
			input: struct {
				v      ValidationData
				status measures.Status
			}{
				v: ValidationData{
					Type:   Qualifier,
					Keys:   []string{StartDate},
					Config: map[string]string{"test": "test"},
				},
				status: measures.Invalid,
			},
		},
		"Should return err with invalid keys": {
			want: struct {
				err   bool
				value ValidatorQualifier
			}{
				err:   true,
				value: ValidatorQualifier{},
			},
			input: struct {
				v      ValidationData
				status measures.Status
			}{
				v: ValidationData{
					Type:   Qualifier,
					Keys:   []string{AI},
					Config: map[string]string{},
				},
				status: measures.Invalid,
			},
		},
		"Should return valid Qualifier": {
			want: struct {
				err   bool
				value ValidatorQualifier
			}{
				err: false,
				value: ValidatorQualifier{
					ValidatorBase: ValidatorBase{
						Type:   Qualifier,
						Keys:   []string{Qualifier},
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
					Type:   Qualifier,
					Keys:   []string{Qualifier},
					Config: map[string]string{},
				},
				status: measures.Invalid,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {

			validator, err := NewValidatorQualifier(testCase.input.v, testCase.input.status)

			assert.Equal(t, testCase.want.value, validator, testName)

			if testCase.want.err {
				assert.NotNil(t, err, testName)
			} else {
				assert.Nil(t, err, testName)
			}
		})
	}
}

func Test_Unit_Services_Qualifier_Validator(t *testing.T) {
	tests := map[string]struct {
		input         InputQualifier
		updateMeasure bool
	}{
		"should be nil if quality 0 (00000000)": {
			input: InputQualifier{
				measure: MeasureValidatable{
					Qualifier: "0",
					WhiteListKeys: map[string]struct{}{
						Qualifier: {},
					},
				},

				ValidatorQualifier: ValidatorQualifier{
					ValidatorBase: NewValidatorBase(ValidationData{
						Type: QualifierValidation,
						Keys: []string{"qualifier"},
					}, measures.Invalid),
				},
			},

			updateMeasure: false,
		},
		"should be nil if quality 16 (00010000)": {
			input: InputQualifier{
				measure: MeasureValidatable{
					Qualifier: "16",
					WhiteListKeys: map[string]struct{}{
						Qualifier: {},
					},
				},
				ValidatorQualifier: ValidatorQualifier{
					ValidatorBase: NewValidatorBase(ValidationData{
						Type: QualifierValidation,
						Keys: []string{"qualifier"},
					}, measures.Invalid),
				},
			},
			updateMeasure: false,
		},
		"should be updated measure to new status if quality 96 (01100000)": {
			input: InputQualifier{
				measure: MeasureValidatable{
					Qualifier: "01100000",
					WhiteListKeys: map[string]struct{}{
						Qualifier: {},
					},
				},
				ValidatorQualifier: ValidatorQualifier{
					ValidatorBase: NewValidatorBase(ValidationData{
						Type: QualifierValidation,
						Keys: []string{"qualifier"},
					}, measures.Invalid),
				},
			},
			updateMeasure: true,
		},
		"should be updated measure to new status if quality 164 (10100100)": {
			input: InputQualifier{
				measure: MeasureValidatable{
					Qualifier: "164",
					WhiteListKeys: map[string]struct{}{
						Qualifier: {},
					},
				},
				ValidatorQualifier: ValidatorQualifier{
					ValidatorBase: NewValidatorBase(ValidationData{
						Type: QualifierValidation,
						Keys: []string{"qualifier"},
					}, measures.Invalid),
				},
			},
			updateMeasure: true,
		},

		"should be nil if quality 00010000": {
			input: InputQualifier{
				measure: MeasureValidatable{
					Qualifier: "00010000",
					WhiteListKeys: map[string]struct{}{
						Qualifier: {},
					},
				},
				ValidatorQualifier: ValidatorQualifier{
					ValidatorBase: NewValidatorBase(ValidationData{
						Type: QualifierValidation,
						Keys: []string{"qualifier"},
					}, measures.Invalid),
				},
			},
			updateMeasure: false,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			base := testCase.input.ValidatorQualifier.Validate(testCase.input.measure)
			if testCase.updateMeasure {
				assert.Equal(t, &testCase.input.ValidatorQualifier.ValidatorBase, base, testCase)
			} else {
				assert.Nil(t, base, testCase)
			}
		})
	}
}
