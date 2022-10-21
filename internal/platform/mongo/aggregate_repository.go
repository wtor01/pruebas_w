package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func NewAggregateRepositoryMongo(client *mongo.Client, database string) *AggregateRepositoryMongo {
	return &AggregateRepositoryMongo{
		client:                client,
		Database:              database,
		CollectionMeterConfig: "view_meter_configs_joined",
		CollectionAggregate:   "aggregations",
	}
}

type AggregateRepositoryMongo struct {
	client                *mongo.Client
	Database              string
	CollectionMeterConfig string
	CollectionAggregate   string
}

func (repository AggregateRepositoryMongo) GenerateAggregation(ctx context.Context, query aggregations.ConfigScheduler) ([]aggregations.Aggregation, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	collection := repository.client.Database(repository.Database).Collection(repository.CollectionMeterConfig)

	match := append(bson.D{
		{"distributor_id", query.DistributorId},
	}, filtersDateInventoryRepositoryMongo(query.Date)...)

	matchStage := bson.D{
		{"$match",
			match,
		},
	}

	groupStageID := bson.D{}

	for _, feature := range query.Features {
		groupStageID = append(groupStageID, bson.D{{feature.ID, "$" + feature.Field}}...)
	}
	groupStage := bson.D{
		{
			"$group",
			bson.D{
				{"_id",
					groupStageID,
				},
				{
					"count",
					bson.D{
						{"$sum", 1},
					},
				},
				{
					"service_points",
					bson.D{
						{
							"$addToSet",
							bson.M{
								"$concat": bson.A{"$service_point.cups", "|||", "$meter.serial_number"},
							},
						},
					},
				},
			},
		},
	}

	features := bson.A{}
	for _, feature := range query.Features {
		features = append(features, bson.M{
			"id":    feature.ID,
			"name":  feature.Name,
			"field": feature.Field,
			"value": "$_id." + feature.ID,
		})
	}
	projectStage := bson.D{
		{"$project",
			bson.D{
				{"_id", 0},
				{"service_points",
					bson.D{
						{"$map",
							bson.D{
								{"input", "$service_points"},
								{"as", "grade"},
								{"in",
									bson.D{
										{"cups",
											bson.D{
												{"$arrayElemAt",
													bson.A{
														bson.D{
															{"$split",
																bson.A{
																	"$$grade",
																	"|||",
																},
															},
														},
														0,
													},
												},
											},
										},
										{"serial_number",
											bson.D{
												{"$arrayElemAt",
													bson.A{
														bson.D{
															{"$split",
																bson.A{
																	"$$grade",
																	"|||",
																},
															},
														},
														1,
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
				{"parameters", features},
			},
		},
	}

	aggregate := append(bson.A{}, matchStage, groupStage, projectStage)

	result := make([]aggregations.Aggregation, 0)

	cursor, err := collection.Aggregate(ctxTimeout, aggregate)

	if err != nil {
		return result, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item aggregations.Aggregation
		if err := cursor.Decode(&item); err != nil {
			return []aggregations.Aggregation{}, err
		}

		if agg, err := aggregations.NewAggregation(item.Parameters, query, item.ServicePoints); err == nil {
			result = append(result, agg)
		}
	}
	if err := cursor.Err(); err != nil {
		return []aggregations.Aggregation{}, err
	}

	return result, err
}

func (repository AggregateRepositoryMongo) SaveAllAggregations(ctx context.Context, aggregationList []aggregations.Aggregation) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	coll := repository.client.Database(repository.Database).Collection(repository.CollectionAggregate)

	var models []mongo.WriteModel

	for _, agg := range aggregationList {
		updateOneModel := mongo.NewUpdateOneModel()
		updateOneModel.SetFilter(bson.D{
			{"_id", agg.Id},
		})
		updateOneModel.SetUpsert(true)
		updateOneModel.SetUpdate(bson.D{
			{"$set", agg},
		})
		models = append(models, updateOneModel)
	}

	_, err := coll.BulkWrite(ctxTimeout, models)

	return err
}
