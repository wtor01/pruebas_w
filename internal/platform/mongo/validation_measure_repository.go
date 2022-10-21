package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewValidationMeasureRepository(client *mongo.Client, database string, collectionLoadCurve string, collectionDailyLoadCurve string, collectionMonthlyClosure string, collectionDailyClosure string, loc *time.Location) *ValidationMeasureRepository {
	return &ValidationMeasureRepository{
		client:                   client,
		location:                 loc,
		Database:                 database,
		CollectionLoadCurve:      collectionLoadCurve,
		CollectionDailyLoadCurve: collectionDailyLoadCurve,
		CollectionMonthlyClosure: collectionMonthlyClosure,
		CollectionDailyClosure:   collectionDailyClosure,
	}
}

type ValidationMeasureRepository struct {
	client                   *mongo.Client
	location                 *time.Location
	Database                 string
	CollectionLoadCurve      string
	CollectionDailyLoadCurve string
	CollectionMonthlyClosure string
	CollectionDailyClosure   string
}

func (repository ValidationMeasureRepository) GetMonthlyClosureByCup(ctx context.Context, q validations.QueryClosedCupsMeasureOnDate) (validations.ProcessedMonthlyClosure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.CollectionMonthlyClosure)
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"cups", q.CUPS},
					{"validation_status", measures.Valid},
					{"end_date", bson.M{
						"$gte": q.Date.AddDate(0, 0, -3),
						"$lt":  q.Date.AddDate(0, 0, 4),
					}},
				},
			},
		},
	}
	opts := options.FindOne().SetSort(bson.D{{"end_date", -1}})

	var measure validations.ProcessedMonthlyClosure

	result := coll.FindOne(ctxTimeout, filter, opts)

	err := result.Decode(&measure)

	if err != nil {
		return validations.ProcessedMonthlyClosure{}, err
	}

	return measure, err
}
func (repository ValidationMeasureRepository) GetLoadCurveByQuery(ctx context.Context, q validations.QueryCurveCupsMeasureOnDate) ([]validations.ProcessedLoadCurve, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.CollectionLoadCurve)

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"cups", q.CUPS},
					{"end_date", bson.M{
						"$gte": q.StartDate,
						"$lt":  q.EndDate,
					}},
				},
			},
		},
	}
	opts := options.Find().SetSort(bson.D{{"end_date", 1}})
	var b []validations.ProcessedLoadCurve

	result, err := coll.Find(ctxTimeout, filter, opts)

	if err != nil {
		return make([]validations.ProcessedLoadCurve, 0), err
	}

	err = result.All(ctxTimeout, &b)

	if err != nil {
		return make([]validations.ProcessedLoadCurve, 0), err
	}

	return b, err
}
