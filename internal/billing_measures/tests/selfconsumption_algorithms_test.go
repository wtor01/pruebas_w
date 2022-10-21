package tests

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/tests/fixtures"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func createFloat(x float64) *float64 {
	return &x
}
func Test_Unit_Domain_Algorithms_SelfConsumptionConfA(t *testing.T) {
	type input struct {
		b *billing_measures.BillingSelfConsumption
	}

	type output struct {
		b billing_measures.BillingSelfConsumption
	}

	tests := map[string]struct {
		input  input
		output output
	}{
		"Execute 1": {
			input: input{
				b: &billing_measures.BillingSelfConsumption{
					Curve: []billing_measures.BillingSelfConsumptionCurve{{
						Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
							AI: 10,
							AE: 20,
						}},
					}},
				},
			},
			output: output{
				b: billing_measures.BillingSelfConsumption{
					Curve: []billing_measures.BillingSelfConsumptionCurve{{
						Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
							AI: 10,
							AE: 20,
							BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
								EHEX: createFloat(10),
								EHCR: createFloat(0),
							},
						}},
						BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
							EHEX: createFloat(10),
							EHCR: createFloat(0),
						},
					}},
				},
			},
		},
		"Execute 2": {
			input: input{
				b: &billing_measures.BillingSelfConsumption{
					Curve: []billing_measures.BillingSelfConsumptionCurve{{
						Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
							AI: 40,
							AE: 30,
						}},
					},
						{
							Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
								AI: 10,
								AE: 20,
							}},
						},
						{
							Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
								AI: 400,
								AE: 30,
							}},
						},
						{
							Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
								AI: 40,
								AE: 10,
							}},
						}},
				},
			},
			output: output{
				b: billing_measures.BillingSelfConsumption{
					Curve: []billing_measures.BillingSelfConsumptionCurve{{
						Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
							AI: 40,
							AE: 30,
							BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
								EHCR: createFloat(10),
								EHEX: createFloat(0),
							},
						}},
						BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
							EHCR: createFloat(10),
							EHEX: createFloat(0),
						},
					},
						{
							Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
								AI: 10,
								AE: 20,
								BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
									EHCR: createFloat(0),
									EHEX: createFloat(10),
								},
							}},
							BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
								EHCR: createFloat(0),
								EHEX: createFloat(10),
							},
						},
						{
							Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
								AI: 400,
								AE: 30,
								BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
									EHCR: createFloat(370),
									EHEX: createFloat(0),
								},
							}},
							BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
								EHCR: createFloat(370),
								EHEX: createFloat(0),
							},
						},
						{
							Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
								AI: 40,
								AE: 10,
								BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
									EHCR: createFloat(30),
									EHEX: createFloat(0),
								},
							}},
							BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
								EHCR: createFloat(30),
								EHEX: createFloat(0),
							},
						}},
				},
			},
		},

		"Example real": {
			input: input{
				b: &fixtures.BillingSelfConsumption_fixtures_A_real_example_1_input,
			},
			output: output{b: fixtures.BillingSelfConsumption_fixtures_A_real_example_1_output},
		},
	}

	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			algConf := billing_measures.NewSelfConsumptionConfA(testCase.input.b)
			cchID := algConf.ID()
			algConf.Execute(ctx)
			assert.Equal(t, &testCase.output.b, algConf.B)
			assert.Equal(t, "CONFIGURATION_A", cchID)
		})
	}
}
func Test_Unit_Domain_Algorithms_SelfConsumptionConfB(t *testing.T) {
	type input struct {
		b *billing_measures.BillingSelfConsumption
	}

	type output struct {
		b billing_measures.BillingSelfConsumption
	}

	tests := map[string]struct {
		input  input
		output output
	}{
		"Execute 1": {
			input: input{
				b: &billing_measures.BillingSelfConsumption{
					Points: []billing_measures.BillingSelfConsumptionPoint{
						{
							ServicePointType: billing_measures.GdServicePointType,
							CUPS:             "1234",
						},
					},
					Curve: []billing_measures.BillingSelfConsumptionCurve{
						{
							Points: []billing_measures.BillingSelfConsumptionCurvePoint{
								{
									CUPS: "1234",
									AI:   10,
									AE:   20,
								},
							},
						},
					},
				},
			},
			output: output{
				b: billing_measures.BillingSelfConsumption{
					Points: []billing_measures.BillingSelfConsumptionPoint{
						{
							ServicePointType: billing_measures.GdServicePointType,
							CUPS:             "1234",
						},
					},
					Curve: []billing_measures.BillingSelfConsumptionCurve{
						{
							Points: []billing_measures.BillingSelfConsumptionCurvePoint{
								{
									CUPS: "1234",
									AI:   10,
									AE:   20,
									BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
										EHGN: createFloat(10),
										EHCR: nil,
										EHEX: createFloat(0),
										EHAU: nil,
										EHSA: createFloat(0),
										EHDC: nil,
									},
								},
							},
							BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
								EHGN: createFloat(10),
								EHCR: nil,
								EHEX: nil,
								EHAU: nil,
								EHSA: createFloat(0),
								EHDC: nil,
							},
						},
					},
				},
			},
		},
	}

	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			algConf := billing_measures.NewSelfConsumptionConfB(testCase.input.b)
			cchID := algConf.ID()
			algConf.Execute(ctx)
			assert.Equal(t, &testCase.output.b, algConf.B)
			assert.Equal(t, "CONFIGURATION_B", cchID)
		})
	}
}

func Test_Unit_Domain_Algorithms_SelfConsumptionConfCAndB(t *testing.T) {
	type input struct {
		b *billing_measures.BillingSelfConsumption
	}

	type output struct {
		b billing_measures.BillingSelfConsumption
	}

	tests := map[string]struct {
		input  input
		output output
	}{
		"Execute 1": {
			input: input{
				b: &billing_measures.BillingSelfConsumption{
					Points: []billing_measures.BillingSelfConsumptionPoint{{CUPS: "1234", ServicePointType: billing_measures.ConsumoServicePointType}, {CUPS: "4321", ServicePointType: billing_measures.GdServicePointType}},
					Curve: []billing_measures.BillingSelfConsumptionCurve{{

						Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
							CUPS: "1234",
							AI:   10,
							AE:   20,
						}, {
							CUPS: "4321",
							AI:   30,
							AE:   40,
						}},
					}},
				},
			},
			output: output{
				b: fixtures.BillingSelfConsumption_fixtures_C_B,
			},
		},
		//"Execute real example": {
		//	input:  input{b: &fixtures.BillingSelfConsumption_fixtures_C_B_real_example_1_input},
		//	output: output{b: fixtures.BillingSelfConsumption_fixtures_C_B_real_example_1_output},
		//},
	}

	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			clientConsumCoefficient := new(mocks.ConsumCoefficientRepository)

			clientConsumCoefficient.On("Search", mock.Anything, billing_measures.QueryConsumCoefficient{}).Return(0.3, nil)
			algConf := billing_measures.NewSelfConsumptionConfCAndB(testCase.input.b, clientConsumCoefficient)
			cchID := algConf.ID()
			algConf.Execute(ctx)
			assert.Equal(t, &testCase.output.b, algConf.B)
			assert.Equal(t, "CONFIGURATION_C_B", cchID)
		})
	}
}

func Test_Unit_Domain_Algorithms_SelfConsumptionConfDandE1(t *testing.T) {
	type input struct {
		b *billing_measures.BillingSelfConsumption
	}

	type output struct {
		b billing_measures.BillingSelfConsumption
	}

	tests := map[string]struct {
		input  input
		output output
	}{
		"Execute 1": {
			input: input{
				b: &billing_measures.BillingSelfConsumption{
					Points: []billing_measures.BillingSelfConsumptionPoint{{CUPS: "1234", ServicePointType: billing_measures.ConsumoServicePointType}, {CUPS: "4321", ServicePointType: billing_measures.GdServicePointType}},
					Curve: []billing_measures.BillingSelfConsumptionCurve{{

						Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
							CUPS: "1234",
							AI:   10,
							AE:   20,
						}, {
							CUPS: "4321",
							AI:   30,
							AE:   40,
						}},
					}},
				},
			},
			output: output{
				b: fixtures.BillingSelfConsumption_fixtures_D_E1,
			},
		},
	}

	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			clientConsumCoefficient := new(mocks.ConsumCoefficientRepository)

			clientConsumCoefficient.On("Search", mock.Anything, billing_measures.QueryConsumCoefficient{}).Return(0.3, nil)
			algConf := billing_measures.NewSelfConsumptionConfDandE1(testCase.input.b, clientConsumCoefficient)
			cchID := algConf.ID()
			algConf.Execute(ctx)
			assert.Equal(t, &testCase.output.b, algConf.B)
			assert.Equal(t, "CONFIGURATION_D_E1", cchID)
		})
	}
}

func Test_Unit_Domain_Algorithms_SelfConsumptionConfE1(t *testing.T) {
	type input struct {
		b *billing_measures.BillingSelfConsumption
	}

	type output struct {
		b billing_measures.BillingSelfConsumption
	}

	tests := map[string]struct {
		input  input
		output output
	}{
		"Execute 1": {
			input: input{
				b: &billing_measures.BillingSelfConsumption{
					Points: []billing_measures.BillingSelfConsumptionPoint{{CUPS: "1234", ServicePointType: billing_measures.ConsumoServicePointType}, {CUPS: "4321", ServicePointType: billing_measures.GdServicePointType}},
					Curve: []billing_measures.BillingSelfConsumptionCurve{{

						Points: []billing_measures.BillingSelfConsumptionCurvePoint{{
							CUPS: "1234",
							AI:   10,
							AE:   20,
						}, {
							CUPS: "4321",
							AI:   30,
							AE:   40,
						}},
					}},
				},
			},
			output: output{
				b: fixtures.BillingSelfConsumption_fixtures_E1,
			},
		},
	}

	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()
			clientConsumCoefficient := new(mocks.ConsumCoefficientRepository)

			clientConsumCoefficient.On("Search", mock.Anything, billing_measures.QueryConsumCoefficient{}).Return(0.3, nil)
			algConf := billing_measures.NewSelfConsumptionConfE1(testCase.input.b, clientConsumCoefficient)
			cchID := algConf.ID()
			algConf.Execute(ctx)
			assert.Equal(t, &testCase.output.b, algConf.B)
			assert.Equal(t, "CONFIGURATION_E1", cchID)
		})
	}
}
