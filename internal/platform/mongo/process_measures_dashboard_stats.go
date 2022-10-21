package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ProcessMeasuresDashboardStatsRepositoryMongo struct {
	client                               *mongo.Client
	database                             string
	timeout                              time.Duration
	collectionStatsProcessMeasures       string
	collectionStatsProcessMeasuresByCups string
}

func NewProcessMeasuresDashboardStatsRepositoryMongo(client *mongo.Client, database, collectionStatsProcessMeasures, collectionStatsProcessMeasuresByCups string) *ProcessMeasuresDashboardStatsRepositoryMongo {

	return &ProcessMeasuresDashboardStatsRepositoryMongo{
		client:                               client,
		database:                             database,
		timeout:                              time.Second * 5,
		collectionStatsProcessMeasures:       collectionStatsProcessMeasures,
		collectionStatsProcessMeasuresByCups: collectionStatsProcessMeasuresByCups,
	}
}

func (repository ProcessMeasuresDashboardStatsRepositoryMongo) GetStatisticsGlobal(ctx context.Context, q process_measures.SearchDashboardStats) ([]process_measures.ProcessMeasureDashboardStatsGlobal, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.database).Collection(repository.collectionStatsProcessMeasures)
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"distributor_id", q.DistributorID},
					{"month", q.Month},
					{"year", q.Year},
					{"type", q.Type},
				},
			},
		},
	}
	opts := options.Find().SetSort(bson.D{{"year", 1}, {"month", 1}})

	var stats []process_measures.ProcessMeasureDashboardStatsGlobal

	result, err := coll.Find(ctxTimeout, filter, opts)

	if err != nil {
		return make([]process_measures.ProcessMeasureDashboardStatsGlobal, 0), err
	}

	err = result.All(ctxTimeout, &stats)

	if err != nil {
		return make([]process_measures.ProcessMeasureDashboardStatsGlobal, 0), err
	}

	return stats, err
}

func (repository ProcessMeasuresDashboardStatsRepositoryMongo) GetStatisticsCups(ctx context.Context, q process_measures.SearchDashboardStats) ([]process_measures.ProcessMeasureDashboardStatsGlobal, int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.database).Collection(repository.collectionStatsProcessMeasuresByCups)
	filter := bson.D{
		{"distributor_id", q.DistributorID},
		{"month", q.Month},
		{"year", q.Year},
		{"type", q.Type},
	}
	matchStage := bson.D{
		{"$match", filter},
	}
	sortStage := bson.D{
		{
			"$sort", bson.D{
				{
					"year", 1,
				},
				{
					"month", 1,
				},
			},
		},
	}
	facetStage := bson.D{
		{"$facet",
			bson.M{
				"data": bson.A{
					bson.D{{"$skip", q.Offset}},
					bson.D{{"$limit", q.Limit}},
				},
				"total": bson.A{
					bson.D{{"$count", "total"}},
				},
			},
		},
	}
	projectStage := bson.D{
		{"$project",
			bson.D{
				{"_id", 0},
				{"data", 1},
				{"total", bson.D{{"$first", "$total.total"}}},
			},
		},
	}

	type Row struct {
		Data  []process_measures.ProcessMeasureDashboardStatsGlobal `bson:"data"`
		Count int                                                   `bson:"total"`
	}

	aggregate := append(bson.A{}, matchStage, sortStage, facetStage, projectStage)

	cursor, err := coll.Aggregate(ctxTimeout, aggregate)
	defer cursor.Close(ctx)

	if err != nil {
		return make([]process_measures.ProcessMeasureDashboardStatsGlobal, 0), 0, err
	}
	var row []Row
	err = cursor.All(ctxTimeout, &row)
	if err != nil {
		return make([]process_measures.ProcessMeasureDashboardStatsGlobal, 0), 0, err
	}

	return row[0].Data, row[0].Count, err
}
