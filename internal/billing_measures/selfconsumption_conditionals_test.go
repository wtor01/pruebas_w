package billing_measures

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Domain_BillingMeasures_IsConfigType(t *testing.T) {
	tests := map[string]struct {
		b                *BillingSelfConsumption
		types            []ConfigType
		expectedResponse bool
	}{
		"Should return true because it is valid config type": {
			b: &BillingSelfConsumption{
				Config: BillingSelfConsumptionConfig{ConfType: "A"},
			},
			types:            []ConfigType{ConfigTypeA},
			expectedResponse: true,
		},
		"Should return false because is not valid config type": {
			b: &BillingSelfConsumption{
				Config: BillingSelfConsumptionConfig{ConfType: "B"},
			},
			types:            []ConfigType{ConfigTypeA, "D1"},
			expectedResponse: false,
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			eval := NewIsConfigType(testCase.b, testCase.types)
			result := eval.Eval(ctx)
			assert.Equal(t, testCase.expectedResponse, result)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsTypeConnection(t *testing.T) {
	tests := map[string]struct {
		b                *BillingSelfConsumption
		t                ConnectionType
		expectedResponse bool
	}{
		"Should return true because it is valid type connection": {
			b: &BillingSelfConsumption{
				Config: BillingSelfConsumptionConfig{ConnType: "Próxima a través de red"},
			},
			t:                ExternalConnection,
			expectedResponse: true,
		},
		"Should return false because is not valid type connection": {
			b: &BillingSelfConsumption{
				Config: BillingSelfConsumptionConfig{ConnType: "Invent"},
			},
			t:                "Invent123",
			expectedResponse: false,
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			eval := NewIsTypeConnection(testCase.b, ConnectionType(testCase.t))
			result := eval.Eval(ctx)
			assert.Equal(t, testCase.expectedResponse, result)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsIndividualGeneration(t *testing.T) {
	tests := map[string]struct {
		b                *BillingSelfConsumption
		expectedResponse bool
	}{
		"Should return true because it is individual generation": {
			b: &BillingSelfConsumption{
				Points: []BillingSelfConsumptionPoint{{ServicePointType: GdServicePointType}, {ServicePointType: "test"}},
			},
			expectedResponse: true,
		},
		"Should return false because it is not individual generation": {
			b: &BillingSelfConsumption{
				Points: []BillingSelfConsumptionPoint{{ServicePointType: GdServicePointType}, {ServicePointType: GdServicePointType}, {ServicePointType: GdServicePointType}},
			},
			expectedResponse: false,
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			eval := NewIsIndividualGeneration(testCase.b)
			result := eval.Eval(ctx)
			assert.Equal(t, testCase.expectedResponse, result)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsIndividualConsumer(t *testing.T) {
	tests := map[string]struct {
		b                *BillingSelfConsumption
		expectedResponse bool
	}{
		"Should return true because it is individual consumer": {
			b: &BillingSelfConsumption{
				Config: BillingSelfConsumptionConfig{ConsumerType: "Individual"},
			},
			expectedResponse: true,
		},
		"Should return false because it is not individual consumer": {
			b: &BillingSelfConsumption{
				Config: BillingSelfConsumptionConfig{ConsumerType: "test"},
			},
			expectedResponse: false,
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			eval := NewIsIndividualConsumer(testCase.b)
			result := eval.Eval(ctx)
			assert.Equal(t, testCase.expectedResponse, result)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsCompensation(t *testing.T) {
	tests := map[string]struct {
		b                *BillingSelfConsumption
		expectedResponse bool
	}{
		"Should return true because it is compensation": {
			b: &BillingSelfConsumption{
				Config: BillingSelfConsumptionConfig{Compensation: true},
			},
			expectedResponse: true,
		},
		"Should return false because it is not compensation": {
			b: &BillingSelfConsumption{
				Config: BillingSelfConsumptionConfig{Compensation: false},
			},
			expectedResponse: false,
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			eval := NewIsCompensation(testCase.b)
			result := eval.Eval(ctx)
			assert.Equal(t, testCase.expectedResponse, result)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_AreExcedents(t *testing.T) {
	tests := map[string]struct {
		b                *BillingSelfConsumption
		expectedResponse bool
	}{
		"Should return true because are excedents": {
			b: &BillingSelfConsumption{
				Config: BillingSelfConsumptionConfig{Excedents: true},
			},
			expectedResponse: true,
		},
		"Should return false because are not excedents": {
			b: &BillingSelfConsumption{
				Config: BillingSelfConsumptionConfig{Excedents: false},
			},
			expectedResponse: false,
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			eval := NewAreExcedents(testCase.b)
			result := eval.Eval(ctx)
			assert.Equal(t, testCase.expectedResponse, result)
		})
	}
}

func Test_Unit_Domain_BillingMeasures_IsSSAA(t *testing.T) {
	tests := map[string]struct {
		b                *BillingSelfConsumption
		expectedResponse bool
	}{
		"Should return true because are ssaa points": {
			b: &BillingSelfConsumption{
				Points: []BillingSelfConsumptionPoint{{ServicePointType: SsaaServicePointType}, {ServicePointType: "test"}}},
			expectedResponse: true,
		},
		"Should return false because are not ssaa points": {
			b: &BillingSelfConsumption{
				Points: []BillingSelfConsumptionPoint{{ServicePointType: "test"}, {ServicePointType: "test"}}},
			expectedResponse: false,
		},
	}
	for testName, _ := range tests {
		testCase := tests[testName]
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			eval := NewIsSSAA(testCase.b)
			result := eval.Eval(ctx)
			assert.Equal(t, testCase.expectedResponse, result)
		})
	}
}
