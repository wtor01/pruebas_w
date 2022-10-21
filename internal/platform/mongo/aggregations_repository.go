package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	mongo_utils "bitbucket.org/sercide/data-ingestion/pkg/db/mongo"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func NewAggregationsRepositoryMongo(client *mongo.Client, database, collection string) *AggregationsRepositoryMongo {
	return &AggregationsRepositoryMongo{
		client:     client,
		database:   database,
		collection: collection,
		timeout:    time.Second * 5,
	}
}

type AggregationsRepositoryMongo struct {
	client     *mongo.Client
	database   string
	timeout    time.Duration
	collection string
}

type Count struct {
	count string `bson:"count"`
}

func (repository AggregationsRepositoryMongo) GetAggregations(ctx context.Context, params aggregations.GetAggregationsDto) ([]aggregations.Aggregation, int64, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, repository.timeout)
	defer cancel()
	collection := repository.client.Database(repository.database).Collection(repository.collection)

	filter := bson.D{{
		Key: "$match",
		Value: bson.M{
			"generation_date": bson.M{
				"$gte": params.StartDate,
				"$lte": params.EndDate,
			},
			"type_id": params.AggregationConfigId,
		},
	}}
	sort := bson.D{{
		Key: "$sort",
		Value: bson.M{
			"generation_date": -1,
		},
	}}

	aggregatePipeline := bson.A{}

	aggregateResult := append(aggregatePipeline, filter, sort)
	aggregateCount := append(aggregatePipeline, filter)

	result, count, err := mongo_utils.Paginate[aggregations.Aggregation](ctxTimeout, collection, mongo_utils.PaginateQuery{
		ResultQuery: mongo_utils.AppendAggregate(aggregateResult),
		CountQuery:  mongo_utils.AppendAggregate(aggregateCount),
		Paginate: db.Pagination{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	})

	return result, int64(count), err
}

func (repository AggregationsRepositoryMongo) GetAggregation(ctx context.Context, params aggregations.GetAggregationDto) (aggregations.Aggregation, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, repository.timeout)
	defer cancel()

	collection := repository.client.Database(repository.database).Collection(repository.collection)

	matchFilter := bson.M{}
	matchFilter["_id"] = params.AggregationConfigId

	result := collection.FindOne(ctxTimeout, matchFilter)

	var aggregation aggregations.Aggregation

	err := result.Decode(&aggregation)

	if err != nil {
		return aggregations.Aggregation{}, err
	}

	return aggregation, nil
}
func (repository AggregationsRepositoryMongo) GetPreviousAggregation(ctx context.Context, params aggregations.GetAggregationDto) (aggregations.AggregationPrevious, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, repository.timeout)
	defer cancel()

	collection := repository.client.Database(repository.database).Collection(repository.collection)

	matchFilter := bson.A{
		bson.D{{"$match", bson.D{{"_id", params.AggregationConfigId}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "aggregations"},
					{"let",
						bson.D{
							{"type_id", "$type_id"},
							{"generation_date", "$generation_date"},
							{"id", "$_id"},
						},
					},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match",
									bson.D{
										{"$expr",
											bson.D{
												{"$and",
													bson.A{
														bson.D{
															{"$eq",
																bson.A{
																	"$type_id",
																	"$$type_id",
																},
															},
														},
														bson.D{
															{"$ne",
																bson.A{
																	"$_id",
																	"$$id",
																},
															},
														},
														bson.D{
															{"$lt",
																bson.A{
																	"$generation_date",
																	"$$generation_date",
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
							bson.D{{"$sort", bson.D{{"generation_date", -1}}}},
							bson.D{{"$limit", 1}},
						},
					},
					{"as", "previous_aggregation"},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"previous_aggregation",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$previous_aggregation",
									0,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"current_cups_aggregation",
						bson.D{
							{"$map",
								bson.D{
									{"input", "$service_points"},
									{"as", "sp"},
									{"in", "$$sp.cups"},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"previous_cups_aggregation",
						bson.D{
							{"$map",
								bson.D{
									{"input", "$previous_aggregation.service_points"},
									{"as", "sp"},
									{"in", "$$sp.cups"},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"previous_cups_aggregation_type",
						bson.D{
							{"$map",
								bson.D{
									{"input", "$previous_cups_aggregation"},
									{"as", "sp"},
									{"in",
										bson.D{
											{"$cond",
												bson.D{
													{"if",
														bson.D{
															{"$in",
																bson.A{
																	"$$sp",
																	"$current_cups_aggregation",
																},
															},
														},
													},
													{"then",
														bson.D{
															{"type", "NEUTRAL"},
															{"cups", "$$sp"},
														},
													},
													{"else",
														bson.D{
															{"type", "OUT"},
															{"cups", "$$sp"},
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
		bson.D{
			{"$addFields",
				bson.D{
					{"current_cups_aggregation_type",
						bson.D{
							{"$map",
								bson.D{
									{"input", "$current_cups_aggregation"},
									{"as", "sp"},
									{"in",
										bson.D{
											{"$cond",
												bson.D{
													{"if",
														bson.D{
															{"$in",
																bson.A{
																	"$$sp",
																	"$previous_cups_aggregation",
																},
															},
														},
													},
													{"then",
														bson.D{
															{"type", "NEUTRAL"},
															{"cups", "$$sp"},
														},
													},
													{"else",
														bson.D{
															{"type", "IN"},
															{"cups", "$$sp"},
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
		bson.D{
			{"$unset",
				bson.A{
					"current_cups_aggregation",
					"previous_cups_aggregation",
					"previous_aggregation.service_points",
				},
			},
		},
	}
	var result []aggregations.AggregationPrevious
	cursor, err := collection.Aggregate(ctxTimeout, matchFilter)

	if err != nil {
		return aggregations.AggregationPrevious{}, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return aggregations.AggregationPrevious{}, err
	}

	if err == nil && len(result) == 0 {
		return aggregations.AggregationPrevious{}, errors.New(fmt.Sprintf(""))
	}

	return result[0], nil
}
