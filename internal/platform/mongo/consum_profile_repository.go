package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewConsumProfileRepository(client *mongo.Client, database string) *ConsumProfileRepositoryMongo {

	return &ConsumProfileRepositoryMongo{
		client:     client,
		database:   database,
		timeout:    time.Second * 5,
		collection: "consum_profiles",
	}
}

type ConsumProfileRepositoryMongo struct {
	client     *mongo.Client
	timeout    time.Duration
	collection string
	database   string
}

func (repository ConsumProfileRepositoryMongo) Save(ctx context.Context, profile billing_measures.ConsumProfile) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, repository.timeout)
	defer cancel()
	coll := repository.client.Database(repository.database).Collection(repository.collection)

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"date", profile.Date},
					{"version", profile.Version},
					{"type", profile.Type},
				},
			},
		},
	}

	upsert := true
	_, err := coll.UpdateOne(ctxTimeout, filter, bson.D{
		{"$set", profile}}, &options.UpdateOptions{
		Upsert: &upsert,
	})

	return err
}

func (repository ConsumProfileRepositoryMongo) Search(ctx context.Context, q billing_measures.QueryConsumProfile) ([]billing_measures.ConsumProfile, error) {
	matchStage := bson.D{
		{"$match",
			bson.M{
				"date": bson.M{
					"$gt":  q.StartDate,
					"$lte": q.EndDate,
				},
			},
		},
	}
	aggregatePipeline := bson.A{}
	projectStage := bson.D{
		{
			"$project",
			bson.M{
				"_id":     1,
				"date":    1,
				"version": 1,
				"type":    1,
				"coef_a":  1,
				"coef_b":  1,
				"coef_c":  1,
				"coef_d":  1,
				"type_sort": bson.M{
					"$switch": bson.M{
						"branches": bson.A{
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$type", billing_measures.FinalTypeConsumProfile},
								},
								"then": 1,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$type", billing_measures.ProvisionalTypeConsumProfile},
								},
								"then": 2,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$type", billing_measures.InitialTypeConsumProfile},
								},
								"then": 3,
							},
						},
						"default": 4,
					},
				},
			},
		},
	}

	sortStage := bson.D{
		{"$sort",
			bson.M{
				"date":      1,
				"type_sort": 1,
			},
		},
	}
	groupStage := bson.D{
		{
			"$group",
			bson.M{
				"_id": bson.M{
					"date": "$date",
				},
				"date": bson.M{
					"$first": "$date",
				},
				"version": bson.M{
					"$first": "$version",
				},
				"type": bson.M{
					"$first": "$type",
				},
				"coef_a": bson.M{
					"$first": "$coef_a",
				},
				"coef_b": bson.M{
					"$first": "$coef_b",
				},
				"coef_c": bson.M{
					"$first": "$coef_c",
				},
				"coef_d": bson.M{
					"$first": "$coef_d",
				},
			},
		},
	}

	sortLast := bson.D{
		{"$sort",
			bson.M{
				"date": 1,
			},
		},
	}
	aggregatePipeline = append(aggregatePipeline, matchStage, projectStage, sortStage, groupStage, sortLast)
	coll := repository.client.Database(repository.database).Collection(repository.collection)
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	cursor, err := coll.Aggregate(ctxTimeout, aggregatePipeline)

	var rows []billing_measures.ConsumProfile

	if err != nil {
		return rows, nil
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &rows); err != nil {
		return nil, err
	}

	return rows, nil
}
