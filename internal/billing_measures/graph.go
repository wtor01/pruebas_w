package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"time"
)

type NodeStatus = string

const (
	NodeStatusFail    = "FAIL"
	NodeStatusSuccess = "SUCCESS"
)

type NodeKey = string

type Graph struct {
	Dict       map[NodeKey]*Node     `bson:"dict" json:"dict,omitempty"`
	From       map[NodeKey][]NodeKey `bson:"-" json:"from,omitempty"`
	To         map[NodeKey][]NodeKey `bson:"-" json:"to,omitempty"`
	Algorithms []string              `bson:"algorithms"`
	StartedAt  time.Time             `bson:"started_at" json:"started_at"`
	FinishedAt time.Time             `bson:"finished_at" json:"finished_at"`
}

func NewGraph() *Graph {
	return &Graph{
		Dict: make(map[NodeKey]*Node),
		From: make(map[NodeKey][]NodeKey),
		To:   make(map[NodeKey][]NodeKey),
	}
}

func (g Graph) firstNodes() []NodeKey {

	allNodes := utils.CopyMap(g.Dict)

	for _, nodeKeys := range g.From {
		for _, k := range nodeKeys {
			delete(allNodes, k)
		}
	}

	return utils.MapKeys(allNodes)
}

func (g *Graph) AddVector(n1 *Node, n2 *Node) *Graph {
	from := g.From[n1.Id]

	if from == nil {
		from = make([]NodeKey, 0)
		g.From[n1.Id] = from
	}

	g.Dict[n1.Id] = n1

	if n2 == nil {
		return g
	}

	g.Dict[n2.Id] = n2

	fromSet := utils.NewSet(from)
	fromSet.Add(n2.Id)

	g.From[n1.Id] = fromSet.Slice()

	to := g.To[n2.Id]

	if to == nil {
		to = make([]NodeKey, 0)
	}

	toSet := utils.NewSet(to)
	toSet.Add(n1.Id)

	g.To[n2.Id] = toSet.Slice()

	return g
}

type Vector struct {
	node     *Node
	children []*Vector
}

func (g *Graph) AddGraph(vectors []*Vector, parent *Node) {

	for _, v := range vectors {
		if parent != nil {
			g.AddVector(parent, v.node)
		}
		g.AddVector(v.node, nil)
		if len(v.children) != 0 {
			g.AddGraph(v.children, v.node)
		}
	}
}

type Node struct {
	Id            NodeKey
	Precondition  Condition   `bson:"-"`
	Algorithms    []Algorithm `bson:"-"`
	Done          bool
	StartedAt     time.Time
	FinishedAt    time.Time
	Status        NodeStatus
	FailAlgorithm string
}

func (g Graph) execute(ctx context.Context, nodeIndex []NodeKey) *Node {
	if len(nodeIndex) == 0 {
		return nil
	}
	node := g.Dict[nodeIndex[0]]
	node.StartedAt = time.Now().UTC()

	if !node.Precondition.Eval(ctx) {
		node.FinishedAt = time.Now().UTC()
		return g.execute(ctx, nodeIndex[1:])
	}

	node.FinishedAt = time.Now().UTC()
	node.Done = true

	to, ok := g.From[node.Id]

	if !ok || len(to) == 0 {
		return node
	}

	return g.execute(ctx, to)
}

func (g *Graph) Execute(ctx context.Context) error {
	g.StartedAt = time.Now().UTC()

	defer func(g *Graph, now func() time.Time) {
		g.FinishedAt = now().UTC()
	}(g, time.Now)

	var node *Node

	if node = g.execute(ctx, g.firstNodes()); node == nil {
		return nil
	}
	if node == nil {
		return nil
	}

	for _, fn := range node.Algorithms {
		if err := fn.Execute(ctx); err != nil {
			node.Status = NodeStatusFail
			node.FailAlgorithm = fn.ID()
			node.FinishedAt = time.Now().UTC()
			return err
		}
	}

	node.Status = NodeStatusSuccess
	node.FinishedAt = time.Now().UTC()

	return nil
}

func GenerateTree(b *BillingMeasure, period measures.PeriodKey, magnitude measures.Magnitude, repository BillingMeasureRepository, repositoryProfiles ConsumProfileRepository) *Graph {
	ctx := GraphContext{}
	g := NewGraph()
	g.AddGraph([]*Vector{
		{
			node: &Node{ // ES TLG
				Id:           "node_level_1_order_1",
				Precondition: NewIsTlg(b),
			},
			children: []*Vector{
				{
					// ES TIPO 3 4 5
					node: &Node{
						Id:           "node_level_2_order_1",
						Precondition: NewIsPointType([]measures.PointType{"3", "4", "5"}, b),
					},
					children: []*Vector{
						{
							// CON MEDIDAS HORARIAS
							node: &Node{
								Id:           "node_level_3_order_1",
								Precondition: NewIsRegisterType(measures.Hourly, b),
							},
							children: []*Vector{
								{
									// SALDO ATR VALIDO PARA PERIODO X
									node: &Node{
										Id:           "node_level_4_order_1",
										Precondition: NewIsValidBillingAtrByPeriod(period, magnitude, b),
									},
									children: []*Vector{
										{
											// CCH NO COMPLETA PARA PERIODO X

											node: &Node{
												Id:           "node_level_5_order_1",
												Precondition: NewNegate(NewIsCurvePeriodCompletedByPeriod(period, b)),
											},
											children: []*Vector{
												{
													// CCH Sin Horas para Px
													node: &Node{
														Id:           "node_level_6_order_1",
														Precondition: NewNegate(NewAreSomeCurveMeasureForPeriod(period, b)),
														Algorithms: []Algorithm{
															NewBalanceCompleted(b, period, magnitude),
															NewCchTotalEstimation(b, period, repositoryProfiles, magnitude),
														},
													},
												},
												{
													// CCH Con alguna hora Válida paa Px
													node: &Node{
														Id:           "node_level_6_order_2",
														Precondition: NewAreSomeCurveMeasureForPeriod(period, b),
													},
													children: []*Vector{
														{
															// ISaldo ATR - CCH < 1KW PARA PERIODO X
															node: &Node{
																Id:           "node_level_7_order_1",
																Precondition: NewIsAtrVsCurveValidByPeriod(period, b, magnitude),
																Algorithms: []Algorithm{
																	NewBalanceCompleted(b, period, magnitude),
																	NewCchPartialAdjustment(b, period, magnitude),
																},
															},
															children: nil,
														},
														{
															// ISaldo ATR - CCH > 1KW PARA PERIODO X
															node: &Node{
																Id:           "node_level_7_order_2",
																Precondition: NewNegate(NewIsAtrVsCurveValidByPeriod(period, b, magnitude)),
															},
															children: []*Vector{
																{
																	// Saldo ATR > CCH
																	node: &Node{
																		Id:           "node_level_8_order_1",
																		Precondition: NewIsBalanceGreaterCch(b, period, magnitude),
																		Algorithms: []Algorithm{
																			NewBalanceCompleted(b, period, magnitude),
																			NewCchPartialEstimation(b, period, repositoryProfiles, magnitude),
																		},
																	},
																},
																{
																	// Saldo ATR < CCH
																	node: &Node{
																		Id:           "node_level_8_order_2",
																		Precondition: NewNegate(NewIsBalanceGreaterCch(b, period, magnitude)),
																		Algorithms: []Algorithm{
																			NewBalanceCompleted(b, period, magnitude),
																			NewCchPartialAdjustment(b, period, magnitude),
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
											// CCH COMPLETA PARA PERIODO X

											node: &Node{
												Id:           "node_level_5_order_2",
												Precondition: NewIsCurvePeriodCompletedByPeriod(period, b),
											},
											children: []*Vector{
												{
													// ISaldo ATR - CCH > 1KW PARA PERIODO X
													node: &Node{
														Id:           "node_level_6_order_3",
														Precondition: NewNegate(NewIsAtrVsCurveValidByPeriod(period, b, magnitude)),
														Algorithms: []Algorithm{
															NewBalanceCompleted(b, period, magnitude),
															NewCchCompleteAdjustment(b, period, magnitude),
														},
													},
													children: nil,
												},
												{
													// ISaldo ATR - CCH < 1KW PARA PERIODO X
													node: &Node{
														Id:           "node_level_6_order_4",
														Precondition: NewIsAtrVsCurveValidByPeriod(period, b, magnitude),
														Algorithms: []Algorithm{
															NewBalanceCompleted(b, period, magnitude),
															NewCchCompleted(b, period, magnitude),
														},
													},
													children: nil,
												},
											},
										},
									},
								},
								{
									// SALDO ATR NO VALIDO PARA PERIODO X
									node: &Node{
										Id:           "node_level_4_order_2",
										Precondition: NewNegate(NewIsValidBillingAtrByPeriod(period, magnitude, b)),
									},
									children: []*Vector{
										{
											//CURVA COMPLETA PARA EL PERIODO X
											node: &Node{
												Id:           "node_level_5_order_3",
												Precondition: NewIsCurvePeriodCompletedByPeriod(period, b),
												Algorithms: []Algorithm{
													NewCchCompleted(b, period, magnitude),
													NewBalanceCalculatedByCch(b, period, magnitude),
												},
											},
										},
										{
											//CURVA INCOMPLETA PARA EL PERIDO X
											node: &Node{
												Id:           "node_level_5_order_4",
												Precondition: NewNegate(NewIsCurvePeriodCompletedByPeriod(period, b)),
											},
											children: []*Vector{
												{
													//NO EXISTE HISTORICO
													node: &Node{
														Id:           "node_level_6_order_5",
														Precondition: NewNegate(NewHasBillingHistory(b, repository, period, magnitude, &ctx)),
													},
													children: []*Vector{
														{
															//CURVA CON ALGUNA HORA PARA PERIODO X
															node: &Node{
																Id:           "node_level_7_order_3",
																Precondition: NewAreSomeCurveMeasureForPeriod(period, b),
																Algorithms: []Algorithm{
																	NewEstimatedBalanceByPowerDemand(b, period, magnitude),
																	NewCchPartialEstimation(b, period, repositoryProfiles, magnitude),
																},
															},
														},
														{
															//CURVA SIN HORAS PARA PERIODO X
															node: &Node{
																Id:           "node_level_7_order_4",
																Precondition: NewNegate(NewAreSomeCurveMeasureForPeriod(period, b)),
																Algorithms: []Algorithm{
																	NewEstimatedBalanceByPowerDemand(b, period, magnitude),
																	NewCchTotalEstimation(b, period, repositoryProfiles, magnitude),
																},
															},
														},
													},
												},
												{
													//EXISTE HISTORICO
													node: &Node{
														Id:           "node_level_6_order_6",
														Precondition: NewHasBillingHistory(b, repository, period, magnitude, &ctx),
													},
													children: []*Vector{
														{
															//CURVA CON ALGUNA HORA PARA PERIODO X
															node: &Node{
																Id:           "node_level_7_order_5",
																Precondition: NewAreSomeCurveMeasureForPeriod(period, b),
																Algorithms: []Algorithm{
																	NewEstimatedHistoryTlg(b, period, magnitude, &ctx),
																	NewCchPartialEstimation(b, period, repositoryProfiles, magnitude),
																},
															},
														},
														{
															//CURVA SIN HORAS PARA PERIODO X
															node: &Node{
																Id:           "node_level_7_order_6",
																Precondition: NewNegate(NewAreSomeCurveMeasureForPeriod(period, b)),
																Algorithms: []Algorithm{
																	NewEstimatedHistoryTlg(b, period, magnitude, &ctx),
																	NewCchTotalEstimation(b, period, repositoryProfiles, magnitude),
																},
															},
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
							// CON MEDIDAS NO HORARIAS
							node: &Node{
								Id:           "node_level_3_order_2",
								Precondition: NewIsRegisterType(measures.QuarterHour, b),
							},
							children: []*Vector{
								{
									// Saldo ATR Válido
									node: &Node{
										Id:           "node_level_4_order_3",
										Precondition: NewIsValidBillingAtrByPeriod(period, magnitude, b),
										Algorithms: []Algorithm{
											NewBalanceCompleted(b, period, magnitude),
											NewCchTotalEstimation(b, period, repositoryProfiles, magnitude),
										},
									},
									children: nil,
								},
								{
									// Saldo ATR No Válido
									node: &Node{
										Id:           "node_level_4_order_4",
										Precondition: NewNegate(NewIsValidBillingAtrByPeriod(period, magnitude, b)),
									},
									children: []*Vector{
										{
											//Existe Historico
											node: &Node{
												Id:           "node_level_5_order_5",
												Precondition: NewHasBillingHistory(b, repository, period, magnitude, &ctx),
												Algorithms: []Algorithm{
													NewEstimatedHistoryTlg(b, period, magnitude, &ctx),
													NewCchTotalEstimation(b, period, repositoryProfiles, magnitude),
												},
											},
											children: nil,
										},
										{
											//No Existe Historico
											node: &Node{
												Id:           "node_level_5_order_6",
												Precondition: NewNegate(NewHasBillingHistory(b, repository, period, magnitude, &ctx)),
												Algorithms: []Algorithm{
													NewEstimatedBalanceByPowerDemand(b, period, magnitude),
													NewCchTotalEstimation(b, period, repositoryProfiles, magnitude),
												},
											},
											children: nil,
										},
									},
								},
							},
						},
					},
				},
				{
					// ES TIPO 1 2
					node: &Node{
						Id:           "node_level_2_order_2",
						Precondition: NewIsPointType([]measures.PointType{"1", "2"}, b),
					},
					children: []*Vector{},
				},
			},
		},
	}, nil)

	return g
}

type GraphContext struct {
	LastHistory                 BillingMeasure
	IsLastHistoryRequested      bool
	ClosedHistory               []BillingMeasure
	IsClosedHistoryRequested    bool
	AnualHistory                BillingMeasure
	IsAnualHistoryRequested     bool
	IterativeHistory            []process_measures.ProcessedLoadCurve
	IsIterativeHistoryRequested bool
	SimpleHistoric              ContextSimpleHistoric
	IsSimpleHistoricRequested   bool
}

type ContextSimpleHistoric struct {
	PreviousLoadCurve []process_measures.ProcessedLoadCurve
	NextLoadCurve     []process_measures.ProcessedLoadCurve
}

func GenerateDcTreeNoTlg(
	b *BillingMeasure,
	processedMeasureRepository process_measures.ProcessedMeasureRepository,
	magnitude measures.Magnitude,
	repository BillingMeasureRepository,
	repositoryProfiles ConsumProfileRepository,
) *Graph {
	ctx := GraphContext{}
	g := NewGraph()
	g.AddGraph([]*Vector{
		{
			// No Es TLG
			node: &Node{
				Id:           "node_level_1_order_1",
				Precondition: NewNegate(NewIsTlg(b)),
			},
			children: []*Vector{
				{
					// Es Punto 1,2,3,4
					node: &Node{
						Id:           "node_level_2_order_1",
						Precondition: NewIsPointType([]measures.PointType{"1", "2", "3", "4"}, b),
					},
					children: []*Vector{
						{
							// Sin Medida Horaria
							node: &Node{
								Id:           "node_level_3_order_1",
								Precondition: NewNegate(NewIsHourly(b)),
							},
							children: []*Vector{
								{
									// Tipo de Punto 4
									node: &Node{
										Id:           "node_level_4_order_1",
										Precondition: NewIsPointType([]measures.PointType{"4"}, b),
									},
									children: []*Vector{
										{
											// Cierre ATR Completo
											node: &Node{
												Id:           "node_level_5_order_1",
												Precondition: NewIsCloseMeasureComplete(b, magnitude),
											},
											children: []*Vector{
												{
													// Saldo ATR Válido
													node: &Node{
														Id:           "node_level_6_order_1",
														Precondition: NewIsBalanceValid(b, magnitude),
														Algorithms: []Algorithm{
															NewBalanceCompleteNoTlg(b, magnitude),
															NewReeBalanceOutline(b, repositoryProfiles, magnitude),
														},
													},
												},
												{
													// Saldo ATR Inválido
													node: &Node{
														Id:           "node_level_6_order_2",
														Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
														Algorithms: []Algorithm{
															NewCloseSum(b, magnitude),
															NewReeBalanceOutline(b, repositoryProfiles, magnitude),
														},
													},
												},
											},
										},
										{
											// Cierre ATR Incompleto
											node: &Node{
												Id:           "node_level_5_order_2",
												Precondition: NewNegate(NewIsCloseMeasureComplete(b, magnitude)),
											},
											children: []*Vector{
												{
													// No Falta mas de un periodo
													node: &Node{
														Id:           "node_level_6_order_3",
														Precondition: NewNegate(NewIsMoreThanOneMissingPeriod(b, magnitude)),
													},
													children: []*Vector{
														{
															// Saldo ATR Válido
															node: &Node{
																Id:           "node_level_7_order_1",
																Precondition: NewIsBalanceValid(b, magnitude),
																Algorithms: []Algorithm{
																	NewFillOneCloseAtr(b, magnitude),
																	NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																},
															},
														},
														{
															// Saldo ATR Invalido
															node: &Node{
																Id:           "node_level_7_order_2",
																Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
															},
															children: []*Vector{
																{
																	// Existe Historico Cierres
																	node: &Node{
																		Id:           "node_level_8_order_1",
																		Precondition: NewHasClosedHistory(b, repository, magnitude, &ctx),
																		Algorithms: []Algorithm{
																			NewCloseHistoryWithoutBalance(b, magnitude, &ctx),
																			NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																		},
																	},
																},
																{
																	// No Existe Historico Cierres
																	node: &Node{
																		Id:           "node_level_8_order_2",
																		Precondition: NewNegate(NewHasClosedHistory(b, repository, magnitude, &ctx)),
																		Algorithms: []Algorithm{
																			NewPowerUseFactorBalance(b, magnitude),
																			NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																		},
																	},
																},
															},
														},
													},
												},
												{
													// Falta mas de un periodo
													node: &Node{
														Id:           "node_level_6_order_4",
														Precondition: NewIsMoreThanOneMissingPeriod(b, magnitude),
													},
													children: []*Vector{
														{
															// Faltan todos los periodos, el totalizador y es casa cerrada
															node: &Node{
																Id:           "node_level_7_order_3",
																Precondition: NewIsHouseCloseAndALlPeriodsAtrAreEmpty(b, magnitude),
																Algorithms: []Algorithm{
																	NewPowerUseFactorBalance(b, magnitude),
																	NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																},
															},
														},
														{
															// No Faltan todos los periodos, el totalizador y no es casa cerrada
															node: &Node{
																Id:           "node_level_7_order_4",
																Precondition: NewNegate(NewIsHouseCloseAndALlPeriodsAtrAreEmpty(b, magnitude)),
															},
															children: []*Vector{
																{
																	// No Existe Historico Cierres
																	node: &Node{
																		Id:           "node_level_8_order_3",
																		Precondition: NewNegate(NewHasClosedHistory(b, repository, magnitude, &ctx)),
																	},
																	children: []*Vector{
																		{
																			// Saldo ATR VÁLIDO
																			node: &Node{
																				Id:           "node_level_9_order_1",
																				Precondition: NewIsBalanceValid(b, magnitude),
																				Algorithms: []Algorithm{
																					NewClosingWithBalance(b, magnitude),
																					NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																				},
																			},
																		},
																		{
																			// Saldo ATR INVALIDO
																			node: &Node{
																				Id:           "node_level_9_order_2",
																				Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
																				Algorithms: []Algorithm{
																					NewPowerUseFactorBalance(b, magnitude),
																					NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																				},
																			},
																		},
																	},
																},
																{
																	// Existe Historico Cierres
																	node: &Node{
																		Id:           "node_level_8_order_4",
																		Precondition: NewHasClosedHistory(b, repository, magnitude, &ctx),
																	},
																	children: []*Vector{
																		{
																			// Saldo ATR VÁLIDO
																			node: &Node{
																				Id:           "node_level_9_order_3",
																				Precondition: NewIsBalanceValid(b, magnitude),
																				Algorithms: []Algorithm{
																					NewClosingHistoryWithBalance(b, magnitude, &ctx),
																					NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																				},
																			},
																		},
																		{
																			// Saldo ATR INVALIDO
																			node: &Node{
																				Id:           "node_level_9_order_4",
																				Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
																				Algorithms: []Algorithm{
																					NewCloseHistoryWithoutBalance(b, magnitude, &ctx),
																					NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																				},
																			},
																		},
																	},
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
									// Tipo de Punto 1,2,3
									node: &Node{
										Id:           "node_level_4_order_2",
										Precondition: NewIsPointType([]measures.PointType{"1", "2", "3"}, b),
										Algorithms: []Algorithm{
											Log{"BALANCE_SUPERVISION"},
											Log{"CCH_SUPERVISION"},
										},
									},
								},
							},
						},
						{
							// Con Medida Horaria
							node: &Node{
								Id:           "node_level_3_order_2",
								Precondition: NewIsHourly(b),
							},
							children: []*Vector{
								{
									// CCH Completa
									node: &Node{
										Id:           "node_level_4_order_3",
										Precondition: NewIsChhCompleted(b),
									},
									children: []*Vector{
										{
											// SALDO ATR VÁLIDO
											node: &Node{
												Id:           "node_level_5_order_3",
												Precondition: NewIsBalanceValid(b, magnitude),
												Algorithms: []Algorithm{
													NewCCHCompleteNoTlg(b, magnitude),
													NewBalanceCompleteNoTlg(b, magnitude),
												},
											},
										},
										{
											// SALDO ATR INVALIDO
											node: &Node{
												Id:           "node_level_5_order_4",
												Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
											},
											children: []*Vector{
												{
													// Cierre ATR Completo
													node: &Node{
														Id:           "node_level_6_order_5",
														Precondition: NewIsCloseMeasureComplete(b, magnitude),
														Algorithms: []Algorithm{
															NewCCHCompleteNoTlg(b, magnitude),
															NewCloseSum(b, magnitude),
														},
													},
												},
												{
													// Cierre ATR Incompleto
													node: &Node{
														Id:           "node_level_6_order_6",
														Precondition: NewNegate(NewIsCloseMeasureComplete(b, magnitude)),
														Algorithms: []Algorithm{
															NewCCHCompleteNoTlg(b, magnitude),
															NewSumHoursNoClosure(b, magnitude),
														},
													},
												},
											},
										},
									},
								},
								{
									// CCH Incompleta
									node: &Node{
										Id:           "node_level_4_order_4",
										Precondition: NewNegate(NewIsChhCompleted(b)),
									},
									children: []*Vector{
										{
											// CCH Invalida
											node: &Node{
												Id:           "node_level_5_order_5",
												Precondition: NewNegate(NewIsChhValid(b)),
											},
											children: []*Vector{
												{
													// Punto 3 o 4, cierres ATR y totalizador completamente vacio y casa cerrada
													node: &Node{
														Id:           "node_level_6_order_7",
														Precondition: NewIsHouseCloseAndCloseAtrAreEmpty(b, magnitude),
														Algorithms: []Algorithm{
															NewConsumZeroClosedHouseNoTLG(b, magnitude),
															NewBalanceZeroConsumption(b, magnitude),
														},
													},
												},
												{
													// !Punto 3 o 4, cierres ATR y totalizador completamente vacio y casa cerrada
													node: &Node{
														Id:           "node_level_6_order_8",
														Precondition: NewNegate(NewIsHouseCloseAndCloseAtrAreEmpty(b, magnitude)),
													},
													children: []*Vector{
														{
															// Existe Historico Iterativo
															node: &Node{
																Id:           "node_level_7_order_5",
																Precondition: NewHasIterativeHistory(b, processedMeasureRepository, magnitude, &ctx),
															},
															children: []*Vector{
																{
																	// Cierre ATR Completo
																	node: &Node{
																		Id:           "node_level_8_order_5",
																		Precondition: NewIsCloseMeasureComplete(b, magnitude),
																	},
																	children: []*Vector{
																		{
																			// Saldo ATR Totalizado Invalido
																			node: &Node{
																				Id:           "node_level_9_order_5",
																				Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
																				Algorithms: []Algorithm{
																					NewCCHWindowsCloseModulated(b, magnitude, &ctx),
																					NewCloseSum(b, magnitude),
																				},
																			},
																		},
																		{
																			// Saldo ATR Totalizado Válido
																			node: &Node{
																				Id:           "node_level_9_order_6",
																				Precondition: NewIsBalanceValid(b, magnitude),
																				Algorithms: []Algorithm{
																					NewCCHWindowsCloseModulated(b, magnitude, &ctx),
																					NewBalanceCompleteNoTlg(b, magnitude),
																				},
																			},
																		},
																	},
																},
																{
																	// Cierre ATR Incompleto
																	node: &Node{
																		Id:           "node_level_8_order_6",
																		Precondition: NewNegate(NewIsCloseMeasureComplete(b, magnitude)),
																	},
																	children: []*Vector{
																		{
																			// Saldo ATR Invalido
																			node: &Node{
																				Id:           "node_level_9_order_7",
																				Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
																				Algorithms: []Algorithm{
																					NewCCHWindows(b, magnitude, &ctx),
																					NewSumHoursNoClosure(b, magnitude),
																				},
																			},
																		},
																		{
																			// Saldo ATR Válido
																			node: &Node{
																				Id:           "node_level_9_order_8",
																				Precondition: NewIsBalanceValid(b, magnitude),
																				Algorithms: []Algorithm{
																					NewCCHWindowsBalanceModulated(b, magnitude, &ctx),
																					NewSumHoursNoClosure(b, magnitude),
																				},
																			},
																		},
																	},
																},
															},
														},
														{
															// No Existe Historico Iterativo
															node: &Node{
																Id:           "node_level_7_order_6",
																Precondition: NewNegate(NewHasIterativeHistory(b, processedMeasureRepository, magnitude, &ctx)),
																Algorithms: []Algorithm{
																	NewPowerUseFactorCCH(b, magnitude),
																	NewSumHoursNoClosure(b, magnitude),
																},
															},
														},
													},
												},
											},
										},
										{
											// CCH Válida
											node: &Node{
												Id:           "node_level_5_order_6",
												Precondition: NewIsChhValid(b),
											},
											children: []*Vector{
												{
													// Hay Ventanas
													node: &Node{
														Id:           "node_level_6_order_9",
														Precondition: NewIsThereChhWindows(b),
													},
													children: []*Vector{
														{
															// No Hay Historico iterativo
															node: &Node{
																Id:           "node_level_7_order_7",
																Precondition: NewNegate(NewHasIterativeHistory(b, processedMeasureRepository, magnitude, &ctx)),
																Algorithms: []Algorithm{
																	NewPowerUseFactorCCH(b, magnitude),
																	NewSumHoursNoClosure(b, magnitude),
																},
															},
														},
														{
															// Hay Historico iterativo
															node: &Node{
																Id:           "node_level_7_order_8",
																Precondition: NewHasIterativeHistory(b, processedMeasureRepository, magnitude, &ctx),
															},
															children: []*Vector{
																{
																	// Cierre ATR Completo
																	node: &Node{
																		Id:           "node_level_8_order_7",
																		Precondition: NewIsCloseMeasureComplete(b, magnitude),
																	},
																	children: []*Vector{
																		{
																			// Saldo ATR Válido
																			node: &Node{
																				Id:           "node_level_9_order_9",
																				Precondition: NewIsBalanceValid(b, magnitude),
																				Algorithms: []Algorithm{
																					NewCCHWindowsCloseModulated(b, magnitude, &ctx),
																					NewBalanceCompleteNoTlg(b, magnitude),
																				},
																			},
																		},
																		{
																			// Saldo ATR Invalido
																			node: &Node{
																				Id:           "node_level_9_order_10",
																				Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
																				Algorithms: []Algorithm{
																					NewCCHWindowsCloseModulated(b, magnitude, &ctx),
																					NewCloseSum(b, magnitude),
																				},
																			},
																		},
																	},
																},
																{
																	// Cierre ATR Incompleto
																	node: &Node{
																		Id:           "node_level_8_order_8",
																		Precondition: NewNegate(NewIsCloseMeasureComplete(b, magnitude)),
																	},
																	children: []*Vector{
																		{
																			// Saldo ATR Válido
																			node: &Node{
																				Id:           "node_level_9_order_11",
																				Precondition: NewIsBalanceValid(b, magnitude),
																				Algorithms: []Algorithm{
																					NewCCHWindowsBalanceModulated(b, magnitude, &ctx),
																					NewSumHoursNoClosure(b, magnitude),
																				},
																			},
																		},
																		{
																			// Saldo ATR Invalido
																			node: &Node{
																				Id:           "node_level_9_order_12",
																				Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
																				Algorithms: []Algorithm{
																					NewCCHWindows(b, magnitude, &ctx),
																					NewSumHoursNoClosure(b, magnitude),
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
													// No Hay Ventanas
													node: &Node{
														Id:           "node_level_6_order_10",
														Precondition: NewNegate(NewIsThereChhWindows(b)),
													},
													children: []*Vector{
														{
															// Cierre ATR Completo
															node: &Node{
																Id:           "node_level_7_order_9",
																Precondition: NewIsCloseMeasureComplete(b, magnitude),
															},
															children: []*Vector{
																{
																	// Saldo ATR Válido
																	node: &Node{
																		Id:           "node_level_8_order_9",
																		Precondition: NewIsBalanceValid(b, magnitude),
																		Algorithms: []Algorithm{
																			NewFlatCastNoTLG(b, magnitude),
																			NewBalanceCompleteNoTlg(b, magnitude),
																		},
																	},
																},
																{
																	// Saldo ATR Invalido
																	node: &Node{
																		Id:           "node_level_8_order_10",
																		Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
																		Algorithms: []Algorithm{
																			NewFlatCastNoTLG(b, magnitude),
																			NewCloseSum(b, magnitude),
																		},
																	},
																},
															},
														},
														{
															// Cierre ATR Incompleto
															node: &Node{
																Id:           "node_level_7_order_10",
																Precondition: NewNegate(NewIsCloseMeasureComplete(b, magnitude)),
															},
															children: []*Vector{
																{
																	// Saldo ATR Válido
																	node: &Node{
																		Id:           "node_level_8_order_11",
																		Precondition: NewIsBalanceValid(b, magnitude),
																		Algorithms: []Algorithm{
																			NewFlatCastBalanceNoTLG(b, magnitude),
																			NewSumHoursNoClosure(b, magnitude),
																		},
																	},
																},
																{
																	// Saldo ATR Invalido
																	node: &Node{
																		Id:           "node_level_8_order_12",
																		Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
																	},
																	children: []*Vector{
																		{
																			//El Hueco esta en dias centrales del periodo de facturacion
																			node: &Node{
																				Id:           "node_level_9_order_13",
																				Precondition: NewIsEmptyCentralHoursCch(b),
																				Algorithms: []Algorithm{
																					NewCchAverage(b, magnitude, &ctx),
																					NewSumHoursNoClosure(b, magnitude),
																				},
																			},
																		},
																		{
																			//El Hueco no esta en dias centrales del periodo de facturacion
																			node: &Node{
																				Id:           "node_level_9_order_14",
																				Precondition: NewNegate(NewIsEmptyCentralHoursCch(b)),
																			},
																			children: []*Vector{
																				{
																					// Hay Historico
																					node: &Node{
																						Id:           "node_level_10_order_1",
																						Precondition: NewIsSimpleHistoric(b, &ctx, magnitude, processedMeasureRepository),
																						Algorithms: []Algorithm{
																							NewCchAverage(b, magnitude, &ctx),
																							NewSumHoursNoClosure(b, magnitude),
																						},
																					},
																				},
																				{
																					// No Hay Historico
																					node: &Node{
																						Id:           "node_level_10_order_2",
																						Precondition: NewNegate(NewIsSimpleHistoric(b, &ctx, magnitude, processedMeasureRepository)),
																						Algorithms: []Algorithm{
																							NewPowerUseFactorCCH(b, magnitude),
																							NewSumHoursNoClosure(b, magnitude),
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
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
					// Es Punto 5
					node: &Node{
						Id:           "node_level_2_order_2",
						Precondition: NewIsPointType([]measures.PointType{"5"}, b),
					},
					children: []*Vector{
						{
							// Con Medida Horaria
							node: &Node{
								Id:           "node_level_3_order_3",
								Precondition: NewIsHourly(b),
								Algorithms: []Algorithm{
									Log{"BALANCE_SUPERVISION"},
									Log{"CCH_SUPERVISION"},
								},
							},
						},
						{
							// Sin Medida Horaria
							node: &Node{
								Id:           "node_level_3_order_4",
								Precondition: NewNegate(NewIsHourly(b)),
							},
							children: []*Vector{
								{
									// Cierre ATR Completo
									node: &Node{
										Id:           "node_level_4_order_5",
										Precondition: NewIsCloseMeasureComplete(b, magnitude),
									},
									children: []*Vector{
										{
											// Saldo ATR Válido
											node: &Node{
												Id:           "node_level_5_order_7",
												Precondition: NewIsBalanceValid(b, magnitude),
												Algorithms: []Algorithm{
													NewBalanceCompleteNoTlg(b, magnitude),
													NewReeBalanceOutline(b, repositoryProfiles, magnitude),
												},
											},
										},
										{
											// Saldo ATR Invalido
											node: &Node{
												Id:           "node_level_5_order_8",
												Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
												Algorithms: []Algorithm{
													NewCloseSum(b, magnitude),
													NewReeBalanceOutline(b, repositoryProfiles, magnitude),
												},
											},
										},
									},
								},
								{
									// Cierre ATR Incompleto
									node: &Node{
										Id:           "node_level_4_order_6",
										Precondition: NewNegate(NewIsCloseMeasureComplete(b, magnitude)),
									},
									children: []*Vector{
										{
											//Faltan todos los periodos de cierre el saldo totalizador y es casa cerrada
											node: &Node{
												Id:           "node_level_5_order_9",
												Precondition: NewIsHouseCloseAndALlPeriodsAtrAreEmpty(b, magnitude),
												Algorithms: []Algorithm{
													NewBalanceZeroConsumption(b, magnitude),
													NewConsumZeroClosedHouseNoTLG(b, magnitude),
												},
											},
										},
										{
											//Cualquier otro caso
											node: &Node{
												Id:           "node_level_5_order_10",
												Precondition: NewNegate(NewIsHouseCloseAndALlPeriodsAtrAreEmpty(b, magnitude)),
											},
											children: []*Vector{
												{
													// No Existe Historico Anual
													node: &Node{
														Id:           "node_level_6_order_11",
														Precondition: NewNegate(NewHasAnualHistory(b, repository, magnitude, &ctx)),
														Algorithms: []Algorithm{
															NewPowerUseFactorBalance(b, magnitude),
															NewReeBalanceOutline(b, repositoryProfiles, magnitude),
														},
													},
												},
												{
													// Existe Historico Anual
													node: &Node{
														Id:           "node_level_6_order_12",
														Precondition: NewHasAnualHistory(b, repository, magnitude, &ctx),
													},
													children: []*Vector{
														{
															// Hay cambio de potencia
															node: &Node{
																Id:           "node_level_7_order_11",
																Precondition: Simple{true}, // TODO: FALTA CONDICIONAL
																Algorithms: []Algorithm{
																	NewPowerUseFactorBalance(b, magnitude),
																	NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																},
															},
														},
														{
															// No Hay cambio de potencia
															node: &Node{
																Id:           "node_level_7_order_12",
																Precondition: NewNegate(Simple{false}), // TODO: FALTA CONDICIONAL
																Algorithms: []Algorithm{
																	NewEstimatedHistoryNoTlg(b, magnitude, &ctx),
																	NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil)

	return g
}

func GenerateGDTreeNoTlg(
	b *BillingMeasure,
	magnitude measures.Magnitude,
	repositoryProfiles ConsumProfileRepository,
	repository BillingMeasureRepository,
	processedMeasureRepository process_measures.ProcessedMeasureRepository,
) *Graph {
	ctx := GraphContext{}
	g := NewGraph()
	g.AddGraph([]*Vector{
		{
			//No es TLG
			node: &Node{
				Id:           "node_level_1_order_1",
				Precondition: NewNegate(NewIsTlg(b)),
			},
			children: []*Vector{
				{
					//Es Punto 1,2,3,4,5
					node: &Node{
						Id:           "node_level_2_order_1",
						Precondition: NewIsPointType([]measures.PointType{"1", "2", "3", "4", "5"}, b),
					},
					children: []*Vector{
						{
							//Sin Medida Horaria
							node: &Node{
								Id:           "node_level_3_order_1",
								Precondition: NewNegate(NewIsHourly(b)),
							},
							children: []*Vector{
								{
									//Tipo de Punto 4 o 5
									node: &Node{
										Id:           "node_level_4_order_1",
										Precondition: NewIsPointType([]measures.PointType{"4", "5"}, b),
									},
									children: []*Vector{
										{
											//Cierre ART Completo
											node: &Node{
												Id:           "node_level_5_order_1",
												Precondition: NewIsCloseMeasureComplete(b, magnitude),
											},
											children: []*Vector{
												{
													//Saldo ATR Válido
													node: &Node{
														Id:           "node_level_6_order_1",
														Precondition: NewIsBalanceValid(b, magnitude),
														Algorithms: []Algorithm{
															NewBalanceCompleteNoTlg(b, magnitude),
															NewReeBalanceOutline(b, repositoryProfiles, magnitude),
														},
													},
												},
												{
													//Saldo ATR Inválido
													node: &Node{
														Id:           "node_level_6_order_2",
														Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
														Algorithms: []Algorithm{
															NewCloseSum(b, magnitude),
															NewReeBalanceOutline(b, repositoryProfiles, magnitude),
														},
													},
												},
											},
										},
										{
											//Cierre ATR Incompleto
											node: &Node{
												Id:           "node_level_5_order_2",
												Precondition: NewNegate(NewIsCloseMeasureComplete(b, magnitude)),
											},
											children: []*Vector{
												{
													//No Falta Más De Un Periodo
													node: &Node{
														Id:           "node_level_6_order_3",
														Precondition: NewNegate(NewIsMoreThanOneMissingPeriod(b, magnitude)),
													},
													children: []*Vector{
														{
															//Saldo ATR Válido
															node: &Node{
																Id:           "node_level_7_order_1",
																Precondition: NewIsBalanceValid(b, magnitude),
																Algorithms: []Algorithm{
																	NewFillOneCloseAtr(b, magnitude),
																	NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																},
															},
														},
														{
															//Saldo ART Inválido
															node: &Node{
																Id:           "node_level_7_order_2",
																Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
															},
															children: []*Vector{
																{
																	//Existe Histórico Cierres
																	node: &Node{
																		Id:           "node_level_8_order_1",
																		Precondition: NewHasClosedHistory(b, repository, magnitude, &ctx),
																		Algorithms: []Algorithm{
																			NewCloseHistoryWithoutBalance(b, magnitude, &ctx),
																			NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																		},
																	},
																},
																{
																	//No Existe Histórico de cierre
																	node: &Node{
																		Id:           "node_level_8_order_2",
																		Precondition: NewNegate(NewHasClosedHistory(b, repository, magnitude, &ctx)),
																		Algorithms: []Algorithm{
																			NewPowerUseFactorBalance(b, magnitude),
																			NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																		},
																	},
																},
															},
														},
													},
												},
												{
													//Falta más De Un Periodo
													node: &Node{
														Id:           "node_level_6_order_4",
														Precondition: NewIsMoreThanOneMissingPeriod(b, magnitude),
													},
													children: []*Vector{
														{
															//No Existe Histórico De Cierres
															node: &Node{
																Id:           "node_level_7_order_3",
																Precondition: NewNegate(NewHasClosedHistory(b, repository, magnitude, &ctx)),
															},
															children: []*Vector{
																{
																	//Saldo ART Válido
																	node: &Node{
																		Id:           "node_level_8_order_3",
																		Precondition: NewIsBalanceValid(b, magnitude),
																		Algorithms: []Algorithm{
																			//CierresSinHistóricoConSaldo<-Algoritmo
																			NewClosingWithBalance(b, magnitude),
																			NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																		},
																	},
																},
																{
																	//Saldo ATR Inválido
																	node: &Node{
																		Id:           "node_level_8_order_4",
																		Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
																		Algorithms: []Algorithm{
																			NewPowerUseFactorBalance(b, magnitude),
																			NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																		},
																	},
																},
															},
														},
														{
															//Existe Histórico de Cierres
															node: &Node{
																Id:           "node_level_7_order_4",
																Precondition: NewHasClosedHistory(b, repository, magnitude, &ctx),
															},
															children: []*Vector{
																{
																	//Saldo ATR Válido
																	node: &Node{
																		Id:           "node_level_8_order_5",
																		Precondition: NewIsBalanceValid(b, magnitude),
																		Algorithms: []Algorithm{
																			NewClosingHistoryWithBalance(b, magnitude, &ctx),
																			NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																		},
																	},
																},
																{
																	//Saldo ATR Inválido
																	node: &Node{
																		Id:           "node_level_8_order_6",
																		Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
																		Algorithms: []Algorithm{
																			NewCloseHistoryWithoutBalance(b, magnitude, &ctx),
																			NewReeBalanceOutline(b, repositoryProfiles, magnitude),
																		},
																	},
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
									//Tipo De Punto 1,2 o 3
									node: &Node{
										Id:           "node_level_4_order_2",
										Precondition: NewIsPointType([]measures.PointType{"1", "2", "3"}, b),
										Algorithms: []Algorithm{
											Log{"BALANCE_SUPERVISION"},
											Log{"CCH_SUPERVISION"},
										},
									},
								},
							},
						},
						{
							//Con Medida Horaria
							node: &Node{
								Id:           "node_level_3_order_2",
								Precondition: NewIsHourly(b),
							},
							children: []*Vector{
								{
									//Saldo ATR Válido
									node: &Node{
										Id:           "node_level_4_order_3",
										Precondition: NewIsBalanceValid(b, magnitude),
									},
									children: []*Vector{
										{
											//CCH Completa
											node: &Node{
												Id:           "node_level_5_order_3",
												Precondition: NewIsChhCompleted(b),
												Algorithms: []Algorithm{
													NewCCHCompleteNoTlg(b, magnitude),
													NewBalanceCompleteNoTlg(b, magnitude),
												},
											},
										},
										{
											//CCH Incompleta
											node: &Node{
												Id:           "node_level_5_order_4",
												Precondition: NewNegate(NewIsChhCompleted(b)),
											},
											children: []*Vector{
												{
													//CCH Inválida
													node: &Node{
														Id:           "node_level_6_order_5",
														Precondition: NewNegate(NewIsChhValid(b)),
													},
													children: []*Vector{
														{
															//Existe Histórico Iterativo
															node: &Node{
																Id:           "node_level_7_order_5",
																Precondition: NewHasIterativeHistory(b, processedMeasureRepository, magnitude, &ctx),
																Algorithms: []Algorithm{
																	NewBalanceCompleteNoTlg(b, magnitude),
																	NewCCHWindowsBalanceModulated(b, magnitude, &ctx),
																},
															},
														},
														{
															//No Existe Histórico Iterativo
															node: &Node{
																Id:           "node_level_7_order_6",
																Precondition: NewNegate(NewHasIterativeHistory(b, processedMeasureRepository, magnitude, &ctx)),
																Algorithms: []Algorithm{
																	NewBalanceCompleteNoTlg(b, magnitude),
																	NewFlatCastBalanceNoTLG(b, magnitude),
																},
															},
														},
													},
												},
												{
													//CCH Válida
													node: &Node{
														Id:           "node_level_6_order_6",
														Precondition: NewIsChhValid(b),
													},
													children: []*Vector{
														{
															//Hay Ventanas
															node: &Node{
																Id:           "node_level_7_order_7",
																Precondition: NewIsThereChhWindows(b),
															},
															children: []*Vector{
																{
																	//Existe Histórico Iterativo
																	node: &Node{
																		Id:           "node_level_8_order_7",
																		Precondition: NewHasIterativeHistory(b, processedMeasureRepository, magnitude, &ctx),
																		Algorithms: []Algorithm{
																			NewBalanceCompleteNoTlg(b, magnitude),
																			NewCCHWindowsBalanceModulated(b, magnitude, &ctx),
																		},
																	},
																},
																{
																	//No Existe Histórico Iterativo
																	node: &Node{
																		Id:           "node_level_8_order_8",
																		Precondition: NewNegate(NewHasIterativeHistory(b, processedMeasureRepository, magnitude, &ctx)),
																		Algorithms: []Algorithm{
																			NewBalanceCompleteNoTlg(b, magnitude),
																			NewFlatCastNoTLG(b, magnitude),
																		},
																	},
																},
															},
														},
														{
															//No Hay Ventanas
															node: &Node{
																Id:           "node_level_7_order_8",
																Precondition: NewNegate(NewIsThereChhWindows(b)),
															},
															children: []*Vector{
																{
																	//El Hueco Está En Las Horas Centrales De La CCH
																	node: &Node{
																		Id:           "node_level_8_order_9",
																		Precondition: NewIsEmptyCentralHoursCch(b),
																		Algorithms: []Algorithm{
																			NewCCHCompleteNoTlg(b, magnitude),
																			NewCchAverage(b, magnitude, &ctx),
																		},
																	},
																},
																{
																	//El Hueco No Está En Las Horas Centrales De La CCH
																	node: &Node{
																		Id:           "node_level_8_order_10",
																		Precondition: NewNegate(NewIsEmptyCentralHoursCch(b)),
																	},
																	children: []*Vector{
																		{
																			//Existe Histórico Simple
																			node: &Node{
																				Id:           "node_level_9_order_1",
																				Precondition: NewIsSimpleHistoric(b, &ctx, magnitude, processedMeasureRepository),
																				Algorithms: []Algorithm{
																					NewCCHCompleteNoTlg(b, magnitude),
																					NewCchAverage(b, magnitude, &ctx),
																				},
																			},
																		},
																		{
																			//No Existe Histórico Simple
																			node: &Node{
																				Id:           "node_level_9_order_2",
																				Precondition: NewNegate(NewIsSimpleHistoric(b, &ctx, magnitude, processedMeasureRepository)),
																			},
																			children: []*Vector{
																				{
																					//Existe Histórico Iterativo
																					node: &Node{
																						Id:           "node_level_10_order_1",
																						Precondition: NewHasIterativeHistory(b, processedMeasureRepository, magnitude, &ctx),
																						Algorithms: []Algorithm{
																							NewCCHCompleteNoTlg(b, magnitude),
																							NewCCHWindowsBalanceModulated(b, magnitude, &ctx),
																						},
																					},
																				},
																				{
																					//No Existe Histórico Iterativo
																					node: &Node{
																						Id:           "node_level_10_order_2",
																						Precondition: NewNegate(NewHasIterativeHistory(b, processedMeasureRepository, magnitude, &ctx)),
																						Algorithms: []Algorithm{
																							NewCCHCompleteNoTlg(b, magnitude),
																							NewFlatCastNoTLG(b, magnitude),
																						},
																					},
																				},
																			},
																		},
																	},
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
									//Saldo ATR Inválido
									node: &Node{
										Id:           "node_level_4_order_4",
										Precondition: NewNegate(NewIsBalanceValid(b, magnitude)),
									},
									children: []*Vector{
										{
											//CCH Completa
											node: &Node{
												Id:           "node_level_5_order_5",
												Precondition: NewIsChhCompleted(b),
												Algorithms: []Algorithm{
													NewCCHCompleteNoTlg(b, magnitude),
													Log{"CCH_LOG_MV_SALDO=Calculado_por_CCH"},
													// TODO FALTA ALGORITMO VERSION NO-TLG NewBalanceCalculatedByCch()
												},
											},
										},
										{
											//CCH Incompleta
											node: &Node{
												Id:           "node_level_5_order_6",
												Precondition: NewNegate(NewIsChhCompleted(b)),
											},
											children: []*Vector{
												{
													//CCH Inválida
													node: &Node{
														Id:           "node_level_6_order_7",
														Precondition: NewNegate(NewIsChhValid(b)),
														Algorithms: []Algorithm{
															NewCCHPenalty(b, magnitude),
															NewSumHoursNoClosure(b, magnitude),
														},
													},
												},
												{
													//CCH Válida
													node: &Node{
														Id:           "node_level_6_order_8",
														Precondition: NewIsChhValid(b),
													},
													children: []*Vector{
														{
															//Hay Ventanas
															node: &Node{
																Id:           "node_level_7_order_9",
																Precondition: NewIsThereChhWindows(b),
																Algorithms: []Algorithm{
																	NewCCHPenalty(b, magnitude),
																	NewSumHoursNoClosure(b, magnitude),
																},
															},
														},
														{
															//No Hay Ventanas
															node: &Node{
																Id:           "node_level_7_order_10",
																Precondition: NewNegate(NewIsThereChhWindows(b)),
															},
															children: []*Vector{
																{
																	//Sí Proceso Aplicado Sobre CCH Que Dispone de CCCH Con Registros Para Las Horas Que Faltan
																	node: &Node{
																		Id:           "node_level_8_order_11",
																		Precondition: NewIsCCCHCompleteGD(b, magnitude, processedMeasureRepository),
																		Algorithms: []Algorithm{
																			NewCCCHCompleteGD(b, magnitude, processedMeasureRepository),
																			NewSumHoursNoClosure(b, magnitude),
																		},
																	},
																},
																{
																	//No Proceso Aplicado Sobre CCH Que Dispone de CCCH Con Registros Para Las Horas Que Faltan
																	node: &Node{
																		Id:           "node_level_8_order_12",
																		Precondition: NewNegate(NewIsCCCHCompleteGD(b, magnitude, processedMeasureRepository)),
																	},
																	children: []*Vector{
																		{
																			//El Hueco Está En Las Horas Centrales De La CCH
																			node: &Node{
																				Id:           "node_level_9_order_3",
																				Precondition: NewIsEmptyCentralHoursCch(b),
																				Algorithms: []Algorithm{
																					NewCchAverage(b, magnitude, &ctx),
																					NewSumHoursNoClosure(b, magnitude),
																				},
																			},
																		},
																		{
																			//El Hueco No Está En Las Horas Centrales De La CCH
																			node: &Node{
																				Id:           "node_level_9_order_4",
																				Precondition: NewNegate(NewIsEmptyCentralHoursCch(b)),
																			},
																			children: []*Vector{
																				{
																					//Existe Histórico Simple
																					node: &Node{
																						Id:           "node_level_10_order_3",
																						Precondition: NewIsSimpleHistoric(b, &ctx, magnitude, processedMeasureRepository),
																						Algorithms: []Algorithm{
																							NewCchAverage(b, magnitude, &ctx),
																							NewSumHoursNoClosure(b, magnitude),
																						},
																					},
																				},
																				{
																					//No Existe Histórico Simple
																					node: &Node{
																						Id:           "node_level_10_order_4",
																						Precondition: NewNegate(NewIsSimpleHistoric(b, &ctx, magnitude, processedMeasureRepository)),
																						Algorithms: []Algorithm{
																							NewCCHPenalty(b, magnitude),
																							NewSumHoursNoClosure(b, magnitude),
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil)
	return g
}
