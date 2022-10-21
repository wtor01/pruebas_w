package inventory

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Services_Distributor_Validator(t *testing.T) {
	tests := map[string]struct {
		input Distributor
		want  error
	}{
		"should be error if wrong name": {
			input: Distributor{
				Name:      "",
				R1:        "320800",
				SmarkiaId: "21312313512",
			},
			want: errors.New("fields cannot be empty"),
		},

		"should be error if wrong R1": {
			input: Distributor{
				Name:      "122",
				R1:        "",
				SmarkiaId: "21312313512",
			},
			want: errors.New("fields cannot be empty"),
		},
		"should be error if wrong Smarkia": {
			input: Distributor{
				Name:      "122",
				R1:        "fasd21",
				SmarkiaId: "",
			},
			want: errors.New("fields cannot be empty"),
		},
		"should be not error": {
			input: Distributor{
				Name:      "Ejemplo",
				R1:        "fasd21",
				SmarkiaId: "3213",
			},
			want: nil,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			result := ValidateDistributorForm(testCase.input)
			assert.Equal(t, testCase.want, result, testCase)
		})
	}
}
