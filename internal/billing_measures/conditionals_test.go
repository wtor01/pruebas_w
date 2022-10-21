package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Domain_BillingMeasures_IsCchNotCompletedByFill(t *testing.T) {
	tests := map[string]struct {
		b                *BillingMeasure
		expectedResponse bool
	}{
		"Should return false because it is not filled": {
			b: &BillingMeasure{
				BillingLoadCurve: []BillingLoadCurve{
					{Origin: measures.Auto},
					{Origin: measures.STG},
				},
			},
			expectedResponse: true,
		},
		"Should return true because it is filled": {
			b: &BillingMeasure{
				BillingLoadCurve: []BillingLoadCurve{
					{Origin: measures.Auto},
					{Origin: measures.Filled},
				},
			},
			expectedResponse: false,
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			actualResponse := testCase.b.IsCchComplete()
			assert.Equal(t, testCase.expectedResponse, actualResponse)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsChhValid(t *testing.T) {
	tests := map[string]struct {
		b                *BillingMeasure
		expectedResponse bool
	}{
		"Should return false because it is not valid": {
			b: &BillingMeasure{
				BillingLoadCurve: []BillingLoadCurve{
					{Origin: measures.Filled},
					{Origin: measures.Filled},
				},
			},
			expectedResponse: false,
		},
		"Should return true because there are some !Filled": {
			b: &BillingMeasure{
				BillingLoadCurve: []BillingLoadCurve{
					{Origin: measures.Auto},
					{Origin: measures.Filled},
				},
			},
			expectedResponse: true,
		},
		"Should return true because there are not filled": {
			b: &BillingMeasure{
				BillingLoadCurve: []BillingLoadCurve{
					{Origin: measures.Auto},
					{Origin: measures.Auto},
				},
			},
			expectedResponse: true,
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			actualResponse := testCase.b.IsChhValid()
			assert.Equal(t, testCase.expectedResponse, actualResponse)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsThereChhWindows(t *testing.T) {
	tests := map[string]struct {
		b                *BillingMeasure
		expectedResponse bool
	}{
		"There are no windows case": {
			b: &BillingMeasure{
				BillingLoadCurve: []BillingLoadCurve{
					{Origin: measures.Auto},
					{Origin: measures.Filled},
					{Origin: measures.Filled},
					{Origin: measures.Filled},
					{Origin: measures.Auto},
					{Origin: measures.Auto},
					{Origin: measures.Filled},
				},
			},
			expectedResponse: false,
		},
		"There are Windows case": {
			b: &BillingMeasure{
				BillingLoadCurve: []BillingLoadCurve{
					{Origin: measures.Auto},
					{Origin: measures.Filled},
					{Origin: measures.Filled},
					{Origin: measures.Filled},
					{Origin: measures.Filled},
					{Origin: measures.Auto},
				},
			},
			expectedResponse: true,
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			actualResponse := testCase.b.IsThereChhWindows()
			assert.Equal(t, testCase.expectedResponse, actualResponse)
		})
	}
}
