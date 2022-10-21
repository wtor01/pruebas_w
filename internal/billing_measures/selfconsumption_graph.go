package billing_measures

import (
	"context"
	"errors"
)

func TestPrecond(ctx context.Context) bool {
	return true
}

type SetSurplusType struct {
	name        string
	SurplusType SurplusType
	b           *BillingSelfConsumption
}

func NewSetSurplusType(surplusType SurplusType, b *BillingSelfConsumption) *SetSurplusType {
	return &SetSurplusType{name: "SET_SurplusType", SurplusType: surplusType, b: b}
}

func (s SetSurplusType) Execute(_ context.Context) error {
	s.b.SurplusType = s.SurplusType
	return nil
}

func (s SetSurplusType) ID() string {
	return s.name
}

type SuperVision struct {
	name string
}

func NewSuperVision() *SuperVision {
	return &SuperVision{name: "SuperVision"}
}

func (s SuperVision) Execute(_ context.Context) error {
	return errors.New(string(SupervisionSelfConsumptionStatus))
}

func (s SuperVision) ID() string {
	return s.name
}

func GenerateSelfConsumptionTree(
	c *BillingSelfConsumption,
	coefficientRepository ConsumCoefficientRepository,
) *Graph {
	//ctx := GraphContext{}
	g := NewGraph()
	g.AddGraph([]*Vector{
		{
			// Es Autoconsumo
			node: &Node{
				Id:           "INIT",
				Precondition: NewIsSelfConsumption(),
			},
			children: []*Vector{
				{
					// Es Config A
					node: &Node{
						Id:           "1",
						Precondition: NewIsConfigType(c, []ConfigType{ConfigTypeA}),
					},
					children: []*Vector{

						// C.IND y G.IND
						{
							node: &Node{
								Id:           "1.1",
								Precondition: NewSeveralEval(NewIsIndividualConsumer(c), NewIsIndividualGeneration(c)),
							},
							children: []*Vector{
								{
									// RED INTERNA
									node: &Node{
										Id:           "1.1.1",
										Precondition: NewIsTypeConnection(c, InternalConnection),
									},
									children: []*Vector{
										{

											// SI EXCEDENTES
											node: &Node{
												Id:           "1.1.1.1",
												Precondition: NewAreExcedents(c),
											},
											children: []*Vector{
												{
													// SI COMPENSACION
													node: &Node{
														Id:           "1.1.1.1.1",
														Precondition: NewIsCompensation(c),
														Algorithms: []Algorithm{
															NewSetSurplusType(DcSurplusType, c), NewSelfConsumptionConfA(c),
														},
													},
												},
												{
													// NO COMPENSACION
													node: &Node{
														Id:           "1.1.1.1.2",
														Precondition: NewNegate(NewIsCompensation(c)),
													},
													children: []*Vector{
														{
															// SI SSAA
															node: &Node{
																Id:           "1.1.1.1.2.1",
																Precondition: NewIsSSAA(c),
																Algorithms: []Algorithm{
																	NewSetSurplusType(GdSurplusType, c), NewSelfConsumptionConfA(c),
																},
															},
														},
														{
															// NO SSAA
															node: &Node{
																Id:           "1.1.1.1.2.2",
																Precondition: NewNegate(NewIsSSAA(c)),
																Algorithms: []Algorithm{
																	NewSetSurplusType(GdSurplusType, c), NewSelfConsumptionConfA(c),
																},
															},
														},
													},
												},
											},
										},
										{

											// NO EXCEDENTES
											node: &Node{
												Id:           "1.1.1.2",
												Precondition: NewNegate(NewAreExcedents(c)),
											},
											children: []*Vector{
												{
													// SI COMPENSACION
													node: &Node{
														Id:           "1.1.1.2.1",
														Precondition: NewIsCompensation(c),
														Algorithms: []Algorithm{
															NewSuperVision(),
														},
													},
												},
												{
													// NO COMPENSACION
													node: &Node{
														Id:           "1.1.1.2.2",
														Precondition: NewNegate(NewIsCompensation(c)),
														Algorithms: []Algorithm{
															NewSetSurplusType(NaSurplusType, c), NewSelfConsumptionConfA(c),
														},
													},
												},
											},
										},
									},
								},
								{
									// RED EXTERNA
									node: &Node{
										Id:           "1.1.2",
										Precondition: NewIsTypeConnection(c, ExternalConnection),
										Algorithms: []Algorithm{
											NewSuperVision(),
										},
									},
								},
							},
						},
						{
							// C.COL y G.IND,C.IND y G.COL, C.COL y G.COL
							node: &Node{
								Id:           "1.2",
								Precondition: NewNegate(NewSeveralEval(NewIsIndividualConsumer(c), NewIsIndividualGeneration(c))),
								Algorithms: []Algorithm{
									NewSuperVision(),
								},
							},
						},
					},
				},
				{
					// Es Config B, C
					node: &Node{
						Id:           "2",
						Precondition: NewIsConfigType(c, []ConfigType{ConfigTypeB, ConfigTypeC}),
					},
					children: []*Vector{
						{
							// ES CONSUMO INDIVIDUAL
							node: &Node{
								Id:           "2.1",
								Precondition: NewIsIndividualConsumer(c),
							},
							children: []*Vector{
								{
									// RED INTERNA
									node: &Node{
										Id:           "2.1.1",
										Precondition: NewIsTypeConnection(c, InternalConnection),
									},
									children: []*Vector{
										{
											// SI ES B
											node: &Node{
												Id:           "2.1.1.1",
												Precondition: NewIsConfigType(c, []ConfigType{ConfigTypeB}),
											},
											children: []*Vector{
												{
													// SI EXCEDENTES
													node: &Node{
														Id:           "2.1.1.1.1",
														Precondition: NewAreExcedents(c),
													},
													children: []*Vector{{
														// SI COMPENSACION
														node: &Node{
															Id:           "2.1.1.1.1.1",
															Precondition: NewIsCompensation(c),
															Algorithms: []Algorithm{
																NewSetSurplusType(DcSurplusType, c), NewSelfConsumptionConfB(c),
															},
														},
													}, {
														// NO COMPENSACION
														node: &Node{
															Id:           "2.1.1.1.1.2",
															Precondition: NewNegate(NewIsCompensation(c)),
														},
														children: []*Vector{
															{
																// SI SSAA
																node: &Node{
																	Id:           "2.1.1.1.1.2.1",
																	Precondition: NewIsSSAA(c),
																	Algorithms: []Algorithm{
																		NewSetSurplusType(GdSurplusType, c), NewSelfConsumptionConfB(c),
																	},
																},
															},
															{
																// NO SSAA
																node: &Node{
																	Id:           "2.1.1.1.1.2.2",
																	Precondition: NewNegate(NewIsSSAA(c)),
																	Algorithms: []Algorithm{
																		NewSetSurplusType(GdSurplusType, c), NewSelfConsumptionConfB(c),
																	},
																},
															}},
													}},
												}, {
													// NO EXCEDENTES
													node: &Node{
														Id:           "2.1.1.1.2",
														Precondition: NewNegate(NewAreExcedents(c)),
													},
													children: []*Vector{{
														// SI COMPENSACION
														node: &Node{
															Id:           "2.1.1.1.2.1",
															Precondition: NewIsCompensation(c),
															Algorithms: []Algorithm{
																NewSuperVision(),
															},
														},
													}, {
														// NO COMPENSACION
														node: &Node{
															Id:           "2.1.1.1.2.2",
															Precondition: NewNegate(NewIsCompensation(c)),
															Algorithms: []Algorithm{
																NewSetSurplusType(NaSurplusType, c),
																NewSelfConsumptionConfB(c),
															},
														},
													}},
												}},
										},
										{
											// SI ES C
											node: &Node{
												Id:           "2.1.1.2",
												Precondition: NewIsConfigType(c, []ConfigType{ConfigTypeC}),
											},
											children: []*Vector{
												{
													// SI EXCEDENTES
													node: &Node{
														Id:           "2.1.1.2.1",
														Precondition: NewAreExcedents(c),
													},
													children: []*Vector{
														{
															// SI COMPENSACION
															node: &Node{
																Id:           "2.1.1.2.1.1",
																Precondition: NewIsCompensation(c),
																Algorithms: []Algorithm{
																	NewSetSurplusType(DcSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
																},
															},
														},
														{
															// NO COMPENSACION
															node: &Node{
																Id:           "2.1.1.2.1.2",
																Precondition: NewNegate(NewIsCompensation(c)),
															},
															children: []*Vector{
																{
																	// SI SSAA
																	node: &Node{
																		Id:           "2.1.1.2.1.2.1",
																		Precondition: NewIsSSAA(c),
																		Algorithms: []Algorithm{
																			NewSetSurplusType(GdSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
																		},
																	},
																},
																{
																	// NO SSAA
																	node: &Node{
																		Id:           "2.1.1.2.1.2.2",
																		Precondition: NewNegate(NewIsSSAA(c)),
																		Algorithms: []Algorithm{
																			NewSetSurplusType(GdSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
																		},
																	},
																}},
														}},
												},
												{
													// NO EXCEDENTES
													node: &Node{
														Id:           "2.1.1.2.2",
														Precondition: NewNegate(NewIsCompensation(c)),
													},
													children: []*Vector{{
														// SI COMPENSACION
														node: &Node{
															Id:           "2.1.1.2.2.1",
															Precondition: NewIsCompensation(c),
															Algorithms: []Algorithm{
																NewSuperVision(),
															},
														},
													}, {
														// NO COMPENSACION
														node: &Node{
															Id:           "2.1.1.2.2.2",
															Precondition: NewNegate(NewIsCompensation(c)),
															Algorithms: []Algorithm{
																NewSetSurplusType(NaSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
															},
														},
													}},
												}},
										}},
								},
								{
									// RED EXTERNA
									node: &Node{
										Id:           "2.1.2",
										Precondition: NewIsTypeConnection(c, ExternalConnection),
									},
									children: []*Vector{
										{
											// SI EXCEDENTES
											node: &Node{
												Id:           "2.1.2.1",
												Precondition: NewAreExcedents(c),
											},
											children: []*Vector{
												{
													// SI COMPENSACION
													node: &Node{
														Id:           "2.1.2.1.1",
														Precondition: NewIsCompensation(c),
														Algorithms: []Algorithm{
															NewSetSurplusType(DcSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
														},
													},
												},
												{
													// NO COMPENSACION
													node: &Node{
														Id:           "2.1.2.1.2",
														Precondition: NewNegate(NewIsCompensation(c)),
													},
													children: []*Vector{{
														// SI SSAA
														node: &Node{
															Id:           "2.1.2.1.2.1",
															Precondition: NewIsSSAA(c),
															Algorithms: []Algorithm{
																NewSetSurplusType(GdSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
															},
														},
													}, {
														// NO SSAA
														node: &Node{
															Id:           "2.1.2.1.2.2",
															Precondition: NewNegate(NewIsSSAA(c)),
															Algorithms: []Algorithm{
																NewSuperVision(),
															},
														},
													}},
												}},
										},
										{
											// NO EXCEDENTES
											node: &Node{
												Id:           "2.1.2.2",
												Precondition: NewNegate(NewAreExcedents(c)),
												Algorithms: []Algorithm{
													NewSuperVision(),
												},
											},
										},
									},
								},
							},
						},
						{
							// ES CONSUMO COLECTIVO
							node: &Node{
								Id:           "2.2",
								Precondition: NewNegate(NewIsIndividualConsumer(c)),
							},
							children: []*Vector{
								{
									// RED INTERNA
									node: &Node{
										Id:           "2.2.1",
										Precondition: NewIsTypeConnection(c, InternalConnection),
									},
									children: []*Vector{
										{
											// SI EXCEDENTES
											node: &Node{
												Id:           "2.2.1.1",
												Precondition: NewAreExcedents(c),
											},
											children: []*Vector{
												{
													// SI COMPENSACION
													node: &Node{
														Id:           "2.2.1.1.1",
														Precondition: NewIsCompensation(c),
														Algorithms: []Algorithm{
															NewSetSurplusType(DcSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
														},
													},
												}, {
													// NO COMPENSACION
													node: &Node{
														Id:           "2.2.1.1.2",
														Precondition: NewNegate(NewIsCompensation(c)),
													},
													children: []*Vector{{
														// SI SSAA
														node: &Node{
															Id:           "2.2.1.1.2.1",
															Precondition: NewIsSSAA(c),
															Algorithms: []Algorithm{
																NewSetSurplusType(GdSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
															},
														},
													}, {
														// NO SSAA
														node: &Node{
															Id:           "2.2.1.1.2.2",
															Precondition: NewNegate(NewIsSSAA(c)),
															Algorithms: []Algorithm{
																NewSetSurplusType(GdSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
															},
														},
													}},
												},
											},
										}, {
											// NO EXCEDENTES
											node: &Node{
												Id:           "2.2.1.2",
												Precondition: NewNegate(NewAreExcedents(c)),
											},
											children: []*Vector{
												{
													// SI COMPENSACION
													node: &Node{
														Id:           "2.2.1.2.1",
														Precondition: NewIsCompensation(c),
														Algorithms: []Algorithm{
															NewSetSurplusType(NaSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
														},
													},
												}, {
													// NO COMPENSACION
													node: &Node{
														Id:           "2.2.1.2.2",
														Precondition: NewNegate(NewIsCompensation(c)),
														Algorithms: []Algorithm{
															NewSetSurplusType(NaSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
														},
													},
												},
											},
										},
									},
								},
								{
									// RED EXTERNA
									node: &Node{
										Id:           "2.2.2",
										Precondition: NewIsTypeConnection(c, ExternalConnection),
									},
									children: []*Vector{
										{
											// SI EXCEDENTES
											node: &Node{
												Id:           "2.2.2.1",
												Precondition: NewAreExcedents(c),
											},
											children: []*Vector{
												{
													// SI COMPENSACION
													node: &Node{
														Id:           "2.2.2.1.1",
														Precondition: NewIsCompensation(c),
														Algorithms: []Algorithm{
															NewSetSurplusType(DcSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
														},
													},
												}, {
													// NO COMPENSACION
													node: &Node{
														Id:           "2.2.2.1.2",
														Precondition: NewNegate(NewIsCompensation(c)),
													},
													children: []*Vector{
														{
															// SI SSAA
															node: &Node{
																Id:           "2.2.2.1.2.1",
																Precondition: NewIsSSAA(c),
																Algorithms: []Algorithm{
																	NewSetSurplusType(GdSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
																},
															},
														},
														{
															// NO SSAA
															node: &Node{
																Id:           "node_level_7_order_2",
																Precondition: NewNegate(NewIsSSAA(c)),
																Algorithms: []Algorithm{
																	NewSetSurplusType(GdSurplusType, c), NewSelfConsumptionConfCAndB(c, coefficientRepository),
																},
															},
														},
													},
												},
											},
										}, {
											// NO EXCEDENTES
											node: &Node{
												Id:           "2.2.2.2",
												Precondition: NewNegate(NewAreExcedents(c)),
												Algorithms: []Algorithm{
													NewSuperVision(),
												},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					// No Es Config A, B, C, D, E1
					node: &Node{
						Id:           "3",
						Precondition: NewNegate(NewIsConfigType(c, []ConfigType{ConfigTypeA, ConfigTypeB, ConfigTypeC, ConfigTypeD, ConfigTypeE1})),
						// TODO: ALG SUPERVISION
						Algorithms: []Algorithm{
							NewSuperVision(),
						},
					},
				},
				{
					// Es Config D, E1
					node: &Node{
						Id:           "4",
						Precondition: NewIsConfigType(c, []ConfigType{ConfigTypeD, ConfigTypeE1}),
					},
					children: []*Vector{
						{
							// ES CONSUMO INDIVIDUAL
							node: &Node{
								Id:           "4.1",
								Precondition: NewIsIndividualConsumer(c),
							},
							children: []*Vector{
								{
									// RED INTERNA
									node: &Node{
										Id:           "4.1.1",
										Precondition: NewIsTypeConnection(c, InternalConnection),
									},
									children: []*Vector{
										{
											node: &Node{
												// TYPE D
												Id:           "4.1.1.1",
												Precondition: NewIsConfigType(c, []ConfigType{ConfigTypeD}),
											},
											children: []*Vector{{
												// SI EXCEDENTES
												node: &Node{
													Id:           "4.1.1.1.1",
													Precondition: NewAreExcedents(c),
												},
												children: []*Vector{{
													// SI COMPENSACION
													node: &Node{
														Id:           "4.1.1.1.1.1",
														Precondition: NewIsCompensation(c),
														Algorithms: []Algorithm{
															NewSetSurplusType(DcSurplusType, c), NewSelfConsumptionConfDandE1(c, coefficientRepository),
														},
													},
												}, {
													// NO COMPENSACION
													node: &Node{
														Id:           "4.1.1.1.1.2",
														Precondition: NewNegate(NewIsCompensation(c)),
													},
													children: []*Vector{
														{
															// SI SSAA
															node: &Node{
																Id:           "4.1.1.1.1.2.1",
																Precondition: NewIsSSAA(c),
																Algorithms: []Algorithm{
																	NewSetSurplusType(GdSurplusType, c), NewSelfConsumptionConfDandE1(c, coefficientRepository),
																},
															},
														},
														{
															// NO SSAA
															node: &Node{
																Id:           "4.1.1.1.1.2.2",
																Precondition: NewNegate(NewIsSSAA(c)),
																Algorithms: []Algorithm{
																	NewSetSurplusType(GdSurplusType, c), NewSelfConsumptionConfDandE1(c, coefficientRepository),
																},
															},
														},
													},
												}},
											}, {
												// NO EXCEDENTES
												node: &Node{
													Id:           "4.1.1.1.2",
													Precondition: NewNegate(NewAreExcedents(c)),
												},
												children: []*Vector{{
													// SI COMPENSACION
													node: &Node{
														Id:           "4.1.1.1.2.1",
														Precondition: NewIsCompensation(c),
														Algorithms: []Algorithm{
															NewSuperVision(),
														},
													},
												}, {
													// NO COMPENSACION
													node: &Node{
														Id:           "4.1.1.1.2.2",
														Precondition: NewNegate(NewIsCompensation(c)),
														Algorithms: []Algorithm{
															NewSetSurplusType(NaSurplusType, c),
															NewSelfConsumptionConfDandE1(c, coefficientRepository),
														},
													},
												}},
											}}},
										{
											node: &Node{
												// TYPE E1
												Id:           "4.1.1.2",
												Precondition: NewIsConfigType(c, []ConfigType{ConfigTypeE1}),
											},
											children: []*Vector{
												{
													// SI EXCEDENTES
													node: &Node{
														Id:           "4.1.1.2.1",
														Precondition: NewAreExcedents(c),
													},
													children: []*Vector{
														{
															// SI COMPENSACION
															node: &Node{
																Id:           "4.1.1.2.1.1",
																Precondition: NewIsCompensation(c),
																Algorithms: []Algorithm{
																	NewSetSurplusType(DcSurplusType, c), NewSelfConsumptionConfE1(c, coefficientRepository),
																},
															},
														},
														{
															// NO COMPENSACION
															node: &Node{
																Id:           "4.1.1.2.1.2",
																Precondition: NewNegate(NewIsCompensation(c)),
															},
															children: []*Vector{
																{
																	// SI SSAA
																	node: &Node{
																		Id:           "4.1.1.2.1.2.1",
																		Precondition: NewIsSSAA(c),
																		Algorithms: []Algorithm{
																			NewSetSurplusType(GdSurplusType, c),
																			NewSelfConsumptionConfE1(c, coefficientRepository),
																		},
																	},
																},
																{
																	// NO SSAA
																	node: &Node{
																		Id:           "4.1.1.2.1.2.2",
																		Precondition: NewNegate(NewIsSSAA(c)),
																		Algorithms: []Algorithm{
																			NewSetSurplusType(GdSurplusType, c),
																			NewSelfConsumptionConfE1(c, coefficientRepository),
																		},
																	},
																},
															},
														}},
												}, {
													// NO EXCEDENTES
													node: &Node{
														Id:           "4.1.1.2.2",
														Precondition: NewNegate(NewAreExcedents(c)),
													},
													children: []*Vector{
														{
															// SI COMPENSACION
															node: &Node{
																Id:           "4.1.1.2.2.1",
																Precondition: NewIsCompensation(c),
																Algorithms: []Algorithm{
																	NewSuperVision(),
																},
															},
														},
														{
															// NO COMPENSACION
															node: &Node{
																Id:           "4.1.1.2.2.2",
																Precondition: NewNegate(NewIsCompensation(c)),
																Algorithms: []Algorithm{
																	NewSetSurplusType(NaSurplusType, c), NewSelfConsumptionConfE1(c, coefficientRepository),
																},
															},
														}},
												},
											},
										},
									},
								},
								{
									// RED EXTERNA
									node: &Node{
										Id:           "4.1.2",
										Precondition: NewIsTypeConnection(c, ExternalConnection),
									},
									children: []*Vector{
										{
											// SI EXCEDENTES
											node: &Node{
												Id:           "4.1.2.1",
												Precondition: NewAreExcedents(c),
											},
											children: []*Vector{
												{
													// SI COMPENSACION
													node: &Node{
														Id:           "4.1.2.1.1",
														Precondition: NewIsCompensation(c),
														Algorithms: []Algorithm{
															NewSetSurplusType(DcSurplusType, c),
															NewSelfConsumptionConfDandE1(c, coefficientRepository),
														},
													},
												},
												{
													// NO COMPENSACION
													node: &Node{
														Id:           "4.1.2.1.2",
														Precondition: NewNegate(NewIsCompensation(c)),
													},
													children: []*Vector{
														{
															// SI SSAA
															node: &Node{
																Id:           "4.1.2.1.2.1",
																Precondition: NewIsSSAA(c),
																Algorithms: []Algorithm{
																	NewSetSurplusType(GdSurplusType, c),
																	NewSelfConsumptionConfDandE1(c, coefficientRepository),
																},
															},
														},
														{
															// NO SSAA
															node: &Node{
																Id:           "4.1.2.1.2.2",
																Precondition: NewNegate(NewIsSSAA(c)),
																Algorithms: []Algorithm{
																	NewSuperVision(),
																},
															},
														}},
												},
											},
										}, {
											// NO EXCEDENTES
											node: &Node{
												Id:           "4.1.2.2",
												Precondition: NewNegate(NewAreExcedents(c)),
												Algorithms: []Algorithm{
													NewSuperVision(),
												},
											},
										},
									},
								},
							},
						},
						{

							// ES CONSUMO COLECTIVO
							node: &Node{
								Id:           "4.2",
								Precondition: NewNegate(NewIsIndividualConsumer(c)),
							},
							children: []*Vector{
								{
									// RED INTERNA
									node: &Node{
										Id:           "4.2.1",
										Precondition: NewIsTypeConnection(c, InternalConnection),
									},
									children: []*Vector{
										{
											// SI EXCEDENTES
											node: &Node{
												Id:           "4.2.1.1",
												Precondition: NewAreExcedents(c),
											},
											children: []*Vector{
												{
													// SI COMPENSACION
													node: &Node{
														Id:           "4.2.1.1.1",
														Precondition: NewIsCompensation(c),
														Algorithms: []Algorithm{
															NewSetSurplusType(DcSurplusType, c),
															NewSelfConsumptionConfDandE1(c, coefficientRepository),
														},
													},
												},
												{
													// NO COMPENSACION
													node: &Node{
														Id:           "4.2.1.1.2",
														Precondition: NewNegate(NewIsCompensation(c)),
													},
													children: []*Vector{
														{
															// SI SSAA
															node: &Node{
																Id:           "4.2.1.1.2.1",
																Precondition: NewIsSSAA(c),
																Algorithms: []Algorithm{
																	NewSetSurplusType(GdSurplusType, c),
																	NewSelfConsumptionConfDandE1(c, coefficientRepository),
																},
															},
														},
														{
															// NO SSAA
															node: &Node{
																Id:           "4.2.1.1.2.2",
																Precondition: NewNegate(NewIsSSAA(c)),
																Algorithms: []Algorithm{
																	NewSetSurplusType(GdSurplusType, c),
																	NewSelfConsumptionConfDandE1(c, coefficientRepository),
																},
															},
														},
													},
												},
											},
										},
										{
											// NO EXCEDENTES
											node: &Node{
												Id:           "4.2.1.2",
												Precondition: NewNegate(NewAreExcedents(c)),
											},
											children: []*Vector{
												{
													// SI COMPENSACION
													node: &Node{
														Id:           "4.2.1.2.1",
														Precondition: NewIsCompensation(c),
														Algorithms: []Algorithm{
															NewSetSurplusType(NaSurplusType, c),
															NewSelfConsumptionConfDandE1(c, coefficientRepository),
														},
													},
												}, {
													// NO COMPENSACION
													node: &Node{
														Id:           "4.2.1.2.2",
														Precondition: NewNegate(NewIsCompensation(c)),
														Algorithms: []Algorithm{
															NewSetSurplusType(NaSurplusType, c),
															NewSelfConsumptionConfDandE1(c, coefficientRepository),
														},
													},
												},
											},
										},
									},
								},
								{
									// RED EXTERNAL
									node: &Node{
										Id:           "4.2.2",
										Precondition: NewIsTypeConnection(c, ExternalConnection),
									},
									children: []*Vector{
										{
											// SI EXCEDENTES
											node: &Node{
												Id:           "4.2.2.1",
												Precondition: NewAreExcedents(c),
											},
											children: []*Vector{
												{
													// SI COMPENSACION
													node: &Node{
														Id:           "4.2.2.1.1",
														Precondition: NewIsCompensation(c),
														Algorithms: []Algorithm{
															NewSetSurplusType(DcSurplusType, c),
															NewSelfConsumptionConfDandE1(c, coefficientRepository),
														},
													},
												}, {
													// NO COMPENSACION
													node: &Node{
														Id:           "4.2.2.1.2",
														Precondition: NewNegate(NewIsCompensation(c)),
													},
													children: []*Vector{
														{
															// SI SSAA
															node: &Node{
																Id:           "4.2.2.1.2.1",
																Precondition: NewIsSSAA(c),
																Algorithms: []Algorithm{
																	NewSetSurplusType(GdSurplusType, c),
																	NewSelfConsumptionConfDandE1(c, coefficientRepository),
																},
															},
														},
														{
															// NO SSAA
															node: &Node{
																Id:           "4.2.2.1.2.2",
																Precondition: NewNegate(NewIsSSAA(c)),
																Algorithms: []Algorithm{
																	NewSetSurplusType(GdSurplusType, c),
																	NewSelfConsumptionConfDandE1(c, coefficientRepository),
																},
															},
														}},
												},
											},
										}, {
											// NO EXCEDENTES
											node: &Node{
												Id:           "4.2.2.2",
												Precondition: NewNegate(NewAreExcedents(c)),
												Algorithms: []Algorithm{
													NewSuperVision(),
												},
											},
										},
									},
								}},
						},
					},
				},
			},
		},
	}, nil)

	return g
}
