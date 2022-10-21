package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_PkgUtils_Maps_Merge(t *testing.T) {
	tests := map[string]struct {
		want  map[string]string
		input struct {
			old  map[string]string
			newM map[string]string
		}
	}{
		"Should merge maps": {
			want: map[string]string{"a": "1", "b": "2"},
			input: struct {
				old  map[string]string
				newM map[string]string
			}{
				old:  map[string]string{"a": "1"},
				newM: map[string]string{"b": "2"},
			},
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {

			Merge(testCase.input.old, testCase.input.newM)

			assert.Equal(t, testCase.want, testCase.input.old, testName)

		})
	}
}
