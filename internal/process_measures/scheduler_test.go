package process_measures

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Domain_Scheduler_isValidSchedulerFormat(t *testing.T) {
	tests := map[string]struct {
		input string
		want  bool
	}{
		"1": {
			input: "/1 * * * *",
			want:  false,
		},
		"2": {
			input: "1.1 * * * *",
			want:  false,
		},
		"3": {
			input: "1 ? * * *",
			want:  false,
		},
		"4": {
			input: "1-2 * * * *",
			want:  false,
		},
		"5": {
			input: "1",
			want:  false,
		},
		"6": {
			input: "1,1 * * * *",
			want:  false,
		},
		"7": {
			input: "60 * * * *",
			want:  false,
		},
		"8": {
			input: "-1 * * * *",
			want:  false,
		},
		"9": {
			input: "52 24 * * *",
			want:  false,
		},
		"10": {
			input: "52 -1 * * *",
			want:  false,
		},
		"11": {
			input: "52 23 32 * *",
			want:  false,
		},
		"12": {
			input: "52 23 0 * *",
			want:  false,
		},
		"13": {
			input: "52 23 1 0 *",
			want:  false,
		},
		"14": {
			input: "52 23 1 13 *",
			want:  false,
		},
		"15": {
			input: "52 23 1 12 -1",
			want:  false,
		},
		"16": {
			input: "52 23 1 12 7",
			want:  false,
		},
		"17": {
			input: "52 23 1 12 6",
			want:  true,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, testCase.want, isValidSchedulerFormat(testCase.input), testCase)
		})
	}
}

func Test_Unit_Domain_Scheduler_isValidSchedulerName(t *testing.T) {
	tests := map[string]struct {
		input string
		want  bool
	}{
		"fail 1": {
			input: "name!",
			want:  false,
		},
		"fail 2": {
			input: "failaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa123",
			want:  false,
		},
		"fail 3": {
			input: "",
			want:  false,
		},
		"ok 1": {
			input: "name-des-2",
			want:  true,
		},
		"ok 2": {
			input: "name_des-2",
			want:  true,
		},
		"ok 3": {
			input: "name_des-2_",
			want:  true,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, testCase.want, isValidSchedulerName(testCase.input), testCase)
		})
	}
}

func Test_Unit_Domain_Scheduler_isValidSchedulerDescription(t *testing.T) {
	tests := map[string]struct {
		input string
		want  bool
	}{

		"fail 1": {
			input: "failaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa123",
			want:  false,
		},
		"ok 1": {
			input: "description test",
			want:  true,
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, testCase.want, isValidSchedulerDescription(testCase.input), testCase)
		})
	}
}
