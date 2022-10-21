package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	mongo_utils "bitbucket.org/sercide/data-ingestion/pkg/db/mongo"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type BillingMeasuresDashboardRepositoryMongo struct {
	client     *mongo.Client
	database   string
	timeout    time.Duration
	collection string
}

func NewBillingMeasuresDashboardRepositoryMongo(client *mongo.Client, database, collection string) *BillingMeasuresDashboardRepositoryMongo {

	return &BillingMeasuresDashboardRepositoryMongo{
		client:     client,
		database:   database,
		timeout:    time.Second * 5,
		collection: collection,
	}
}

func (repository BillingMeasuresDashboardRepositoryMongo) SearchFiscalBillingMeasures(ctx context.Context, cups string, distributorId string, startDate time.Time, endDate time.Time) ([]billing_measures.BillingMeasure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, repository.timeout)
	defer cancel()

	collection := repository.client.Database(repository.database).Collection(repository.collection)

	filter := bson.D{
		{"cups", cups},
		{"distributor_id", distributorId},
		{"init_date",
			bson.M{
				"$gte": startDate,
			},
		},
		{"end_date",
			bson.M{
				"$lte": endDate,
			},
		},
	}

	opts := options.Find().SetSort(bson.D{{"end_date", 1}})

	var billingMeasures []billing_measures.BillingMeasure

	result, err := collection.Find(ctxTimeout, filter, opts)

	err = result.All(ctxTimeout, &billingMeasures)
	if err != nil {
		return []billing_measures.BillingMeasure{}, err
	}

	return billingMeasures, nil
}

func (repository BillingMeasuresDashboardRepositoryMongo) SearchLastBillingMeasures(ctx context.Context, cups string, distributorId string) (billing_measures.BillingMeasure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, repository.timeout)
	defer cancel()
	collection := repository.client.Database(repository.database).Collection(repository.collection)
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"cups", cups},
					{"distributor_id", distributorId},
				},
			},
		},
	}
	opts := options.FindOne().SetSort(bson.D{{"end_date", -1}}).SetProjection(bson.M{"graph_history": 0})
	var b billing_measures.BillingMeasure

	result := collection.FindOne(ctxTimeout, filter, opts)

	err := result.Decode(&b)

	if err != nil {
		return billing_measures.BillingMeasure{}, err
	}

	return b, err
}

func (repository BillingMeasuresDashboardRepositoryMongo) GroupFiscalMeasureSummary(ctx context.Context, query billing_measures.GroupFiscalMeasureSummaryQuery) (billing_measures.FiscalMeasureSummary, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, repository.timeout)
	defer cancel()

	collection := repository.client.Database(repository.database).Collection(repository.collection)

	matchStage := bson.D{
		{
			"$match", bson.M{
				"distributor_id": query.DistributorId,
				"meter_type":     query.MeterType,
				"end_date": bson.M{
					"$gt":  query.StartDate,
					"$lte": query.EndDate,
				},
			},
		},
	}

	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$distributor_id"},
			{"total", bson.D{{"$sum", 1}}},
			{"balance_real", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.balance_type",
												billing_measures.GeneralReal,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"balance_calculated", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.balance_type",
												billing_measures.GeneralCalculated,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"balance_estimated", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.balance_type",
												billing_measures.GeneralEstimated,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"balance_monthly", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.balance_origin",
												measures.Monthly,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"balance_daily", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.balance_origin",
												measures.Daily,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"balance_other", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.balance_origin",
												measures.Other,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"balance_no_closure", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.balance_origin",
												measures.NoClosure,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"curve_real", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.curve_type",
												billing_measures.GeneralReal,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"curve_adjusted", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.curve_type",
												billing_measures.GeneralAdjusted,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"curve_outlined", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.curve_type",
												billing_measures.GeneralOutlined,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"curve_calculated", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.curve_type",
												billing_measures.GeneralCalculated,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"curve_estimated", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.curve_type",
												billing_measures.GeneralEstimated,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"curve_completed", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.curve_status",
												measures.Complete,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"curve_uncompleted", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.curve_status",
												measures.Incomplete,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
			{"curve_absent", bson.D{
				{"$sum",
					bson.D{
						{"$cond",
							bson.D{
								{"if",
									bson.D{
										{"$eq",
											bson.A{
												"$execution_summary.curve_status",
												measures.Absent,
											},
										},
									},
								},
								{"then", 1},
								{"else", 0},
							},
						},
					},
				},
			}},
		},
		},
	}

	projectStage := bson.D{
		{"$project",
			bson.D{
				{"_id", 0},
				{"total", 1},
				{"balance_type",
					bson.D{
						{"real", "$balance_real"},
						{"calculated", "$balance_calculated"},
						{"estimated", "$balance_estimated"},
					},
				},
				{"curve_type",
					bson.D{
						{"real", "$curve_real"},
						{"adjusted", "$curve_adjusted"},
						{"outlined", "$curve_outlined"},
						{"calculated", "$curve_calculated"},
						{"estimated", "$curve_estimated"},
					},
				},
				{"curve_status",
					bson.D{
						{"completed", "$curve_completed"},
						{"not_completed", "$curve_uncompleted"},
						{"absent", "$curve_absent"},
					},
				},
				{"balance_origin",
					bson.D{
						{"monthly", "$balance_monthly"},
						{"daily", "$balance_daily"},
						{"other", "$balance_other"},
						{"no_closure", "$balance_no_closure"},
					},
				},
			},
		},
	}

	aggregation := append(bson.A{}, matchStage, groupStage, projectStage)

	cursor, err := collection.Aggregate(ctxTimeout, aggregation)

	if err != nil {
		return billing_measures.FiscalMeasureSummary{}, err
	}

	type Result struct {
		Total       int `bson:"total"`
		BalanceType struct {
			Real       int `bson:"real"`
			Calculated int `bson:"calculated"`
			Estimated  int `bson:"estimated"`
		} `bson:"balance_type"`
		BalanceOrigin struct {
			Monthly   int `bson:"monthly"`
			Daily     int `bson:"daily"`
			Other     int `bson:"other"`
			NoClosure int `bson:"no_closure"`
		} `bson:"balance_origin"`
		CurveType struct {
			Real       int `bson:"real"`
			Outlined   int `bson:"outlined"`
			Adjusted   int `bson:"adjusted"`
			Calculated int `bson:"calculated"`
			Estimated  int `bson:"estimated"`
		} `bson:"curve_type"`
		CurveStatus struct {
			Completed    int `bson:"completed"`
			NotCompleted int `bson:"not_completed"`
			Absent       int `bson:"absent"`
		} `bson:"curve_status"`
	}

	var results []Result
	if err = cursor.All(ctx, &results); err != nil {
		return billing_measures.FiscalMeasureSummary{}, err
	}

	if len(results) == 0 {
		return billing_measures.FiscalMeasureSummary{}, errors.New("not found summary")
	}

	result := results[0]
	return billing_measures.FiscalMeasureSummary{
		MeterType: query.MeterType,
		Total:     result.Total,
		CurveType: billing_measures.CurveTypeSummary{
			TypeSummary: billing_measures.TypeSummary{
				Real:       result.CurveType.Real,
				Calculated: result.CurveType.Calculated,
				Estimated:  result.CurveType.Estimated,
			},
			Adjusted: result.CurveType.Adjusted,
			Outlined: result.CurveType.Outlined,
		},
		CurveStatus: billing_measures.CurveStatusSummary{
			Completed:    result.CurveStatus.Completed,
			NotCompleted: result.CurveStatus.NotCompleted,
			Absent:       result.CurveStatus.Absent,
		},
		BalanceType: billing_measures.BalanceTypeSummary{
			TypeSummary: billing_measures.TypeSummary{
				Real:       result.BalanceType.Real,
				Calculated: result.BalanceType.Calculated,
				Estimated:  result.BalanceType.Estimated,
			},
		},
		BalanceOrigin: billing_measures.BalanceOriginSummary{
			Monthly:   result.BalanceOrigin.Monthly,
			Daily:     result.BalanceOrigin.Daily,
			Other:     result.BalanceOrigin.Other,
			NoClosure: result.BalanceOrigin.NoClosure,
		},
	}, nil
}

func (repository BillingMeasuresDashboardRepositoryMongo) SearchBillingMeasureClosureResume(ctx context.Context, billingMeasureID string) (billing_measures.BillingMeasureResumeClosureResponse, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, repository.timeout)
	defer cancel()
	collection := repository.client.Database(repository.database).Collection(repository.collection)
	filter := bson.D{
		{"_id", billingMeasureID},
	}
	opts := options.FindOne().SetSort(bson.D{{"end_date", -1}}).SetProjection(bson.M{
		"periods":                  1,
		"magnitudes":               1,
		"actual_reading_closure":   1,
		"previous_reading_closure": 1,
	})
	var b billing_measures.BillingMeasureResumeClosureResponse
	result := collection.FindOne(ctxTimeout, filter, opts)
	err := result.Decode(&b)
	if err != nil {
		return billing_measures.BillingMeasureResumeClosureResponse{}, err
	}
	return b, err
}
func (repository BillingMeasuresDashboardRepositoryMongo) GetBillingMeasuresTax(ctx context.Context, params billing_measures.QueryBillingMeasuresTax) (billing_measures.BillingMeasuresTaxResult, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	collection := repository.client.Database(repository.database).Collection(repository.collection)

	filter := bson.D{
		{
			"$match", bson.M{
				"meter_type":     params.MeasureType,
				"distributor_id": params.DistributorId,
				"end_date": bson.M{
					"$gt":  params.StartDate,
					"$lte": params.EndDate,
				},
			},
		},
	}
	projectStage := bson.D{{"$project", bson.M{
		"execution_summary": 1,
		"cups":              1,
		"distributor_id":    1,
		"init_date":         1,
		"end_date":          1,
	}}}
	sort := bson.D{{
		Key: "$sort",
		Value: bson.M{
			"end_date": -1,
		},
	}}

	aggregatePipeline := bson.A{}

	aggregateResult := append(aggregatePipeline, filter, sort, projectStage)
	aggregateCount := append(aggregatePipeline, filter)

	result, count, err := mongo_utils.Paginate[billing_measures.BillingMeasuresTax](ctxTimeout, collection, mongo_utils.PaginateQuery{
		ResultQuery: mongo_utils.AppendAggregate(aggregateResult),
		CountQuery:  mongo_utils.AppendAggregate(aggregateCount),
		Paginate: db.Pagination{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	})

	return billing_measures.BillingMeasuresTaxResult{
		Data:  result,
		Count: count,
	}, err
}
