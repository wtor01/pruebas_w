package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	mongo_utils "bitbucket.org/sercide/data-ingestion/pkg/db/mongo"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewGrossMeasureRepositoryDashboardMongo(client *mongo.Client, database, collectionStatsGrossMeasures, collectionStatsGrossMeasuresSerialNumber string) *GrossMeasureRepositoryDashboardMongo {
	return &GrossMeasureRepositoryDashboardMongo{
		client:                                   client,
		Database:                                 database,
		CollectionStatsGrossMeasures:             collectionStatsGrossMeasures,
		CollectionStatsGrossMeasuresSerialNumber: collectionStatsGrossMeasuresSerialNumber,
	}
}

type GrossMeasureRepositoryDashboardMongo struct {
	client                                   *mongo.Client
	Database                                 string
	CollectionStatsGrossMeasures             string
	CollectionStatsGrossMeasuresSerialNumber string
}

func (repository GrossMeasureRepositoryDashboardMongo) ListGrossMeasuresStatisticsSerialNumber(ctx context.Context, q gross_measures.SearchDashboardSerialNumber) (gross_measures.ListGrossMeasuresStatisticsSerialNumberResult, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.CollectionStatsGrossMeasuresSerialNumber)
	filter := bson.D{
		{"distributor_id", q.DistributorId},
		{"month", q.Month},
		{"year", q.Year},
		{"type", q.Type},
	}
	if q.Ghost {
		filter = append(filter, primitive.E{"ghost", q.Ghost})
	}
	matchStage := bson.D{
		{"$match", filter},
	}

	projectStage := bson.D{{"$project", bson.M{
		"distributor_id":     1,
		"month":              1,
		"year":               1,
		"cups":               1,
		"type":               1,
		"serial_number":      1,
		"service_type":       1,
		"service_point_type": 1,
		"dailyStats":         1,
	}}}

	aggregatePipeline := bson.A{}
	aggregateResult := append(aggregatePipeline, matchStage, projectStage)
	aggregateCount := append(aggregatePipeline, matchStage)

	result, count, err := mongo_utils.Paginate[gross_measures.DashboardSerialNumber](ctxTimeout, coll, mongo_utils.PaginateQuery{
		ResultQuery: mongo_utils.AppendAggregate(aggregateResult),
		CountQuery:  mongo_utils.AppendAggregate(aggregateCount),
		Paginate: db.Pagination{
			Limit:  q.Limit,
			Offset: &q.Offset,
		},
	})
	return gross_measures.ListGrossMeasuresStatisticsSerialNumberResult{
		Data:  result,
		Count: count,
	}, err
}
func (repository GrossMeasureRepositoryDashboardMongo) GetStatisticsGlobal(ctx context.Context, q gross_measures.SearchDashboardStats) ([]gross_measures.GrossMeasuresDashboardStatsGlobal, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	coll := repository.client.Database(repository.Database).Collection(repository.CollectionStatsGrossMeasures)
	filter := bson.D{
		{"distributor_id", q.DistributorID},
		{"month", q.Month},
		{"year", q.Year},
		{"type", q.Type},
	}

	opts := options.Find()

	var stats []gross_measures.GrossMeasuresDashboardStatsGlobal

	result, err := coll.Find(ctxTimeout, filter, opts)

	if err != nil {
		return make([]gross_measures.GrossMeasuresDashboardStatsGlobal, 0), err
	}

	err = result.All(ctxTimeout, &stats)

	if err != nil {
		return make([]gross_measures.GrossMeasuresDashboardStatsGlobal, 0), err
	}

	return stats, err
}
