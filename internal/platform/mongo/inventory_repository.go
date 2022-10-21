package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/db"
	mongo_utils "bitbucket.org/sercide/data-ingestion/pkg/db/mongo"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewInventoryRepositoryMongo(client *mongo.Client, database string) *InventoryRepositoryMongo {
	return &InventoryRepositoryMongo{
		client:                client,
		Database:              database,
		CollectionMeterConfig: "view_meter_configs_joined",
	}
}

type InventoryRepositoryMongo struct {
	client                *mongo.Client
	Database              string
	CollectionMeterConfig string
}

func (repository InventoryRepositoryMongo) ListMeterConfigByCups(ctx context.Context, query measures.ListMeterConfigByCups) ([]measures.MeterConfig, error) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "ListMeterConfigByDate")
	defer span.End()

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	opts := options.Find()

	match := bson.D{
		{"distributor_id", query.DistributorId},
		{"cups", query.CUPS},
		{"start_date",
			bson.M{
				"$lte": query.EndDate.UTC(),
			},
		},
		{"contractual_situations.init_date",
			bson.M{
				"$lte": query.EndDate.UTC(),
			},
		},
		{"end_date",
			bson.M{
				"$gte": query.StartDate.UTC(),
			},
		},
		{"contractual_situations.end_date",
			bson.M{
				"$gte": query.StartDate.UTC(),
			},
		},
	}

	cursor, err := repository.client.Database(repository.Database).Collection(repository.CollectionMeterConfig).Find(ctxTimeout, match, opts)

	var results []measures.MeterConfig

	if err != nil {
		return results, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return []measures.MeterConfig{}, err
	}

	return results, nil
}

func filtersDateInventoryRepositoryMongo(date time.Time) bson.D {
	return bson.D{
		{"start_date",
			bson.M{
				"$lte": date.UTC(),
			},
		},
		{"contractual_situations.init_date",
			bson.M{
				"$lte": date.UTC(),
			},
		},
		{"end_date",
			bson.M{
				"$gt": date.UTC(),
			},
		},
		{"contractual_situations.end_date",
			bson.M{
				"$gt": date.UTC(),
			},
		},
	}
}

func (repository InventoryRepositoryMongo) GetMeterConfigByMeter(ctx context.Context, query measures.GetMeterConfigByMeterQuery) (measures.MeterConfig, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	match := append(bson.D{
		{"meter.serial_number", query.MeterSerialNumber},
	}, filtersDateInventoryRepositoryMongo(query.Date)...)

	var config measures.MeterConfig

	result := repository.client.Database(repository.Database).Collection(repository.CollectionMeterConfig).FindOne(ctxTimeout, match)

	err := result.Decode(&config)

	if err != nil {
		return measures.MeterConfig{}, err
	}

	return config, err
}

func (repository InventoryRepositoryMongo) ListMeterConfigByDate(ctx context.Context, query measures.ListMeterConfigByDateQuery) ([]measures.MeterConfig, error) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "ListMeterConfigByDate")
	defer span.End()

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	opts := options.Find()
	opts.SetLimit(int64(query.Limit)).SetSkip(int64(query.Offset))

	match := append(bson.D{
		{"distributor_id", query.DistributorID},
		{"service_point.service_type", query.ServiceType},
		{"service_point.point_type", query.PointType},
		{"type", bson.M{"$in": query.MeterType}},
	}, filtersDateInventoryRepositoryMongo(query.Date)...)

	cursor, err := repository.client.Database(repository.Database).Collection(repository.CollectionMeterConfig).Find(ctxTimeout, match, opts)

	var results []measures.MeterConfig

	if err != nil {
		return results, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return []measures.MeterConfig{}, err
	}

	return results, nil
}

func (repository InventoryRepositoryMongo) CountMeterConfigByDate(ctx context.Context, query measures.ListMeterConfigByDateQuery) (int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	match := append(bson.D{
		{"distributor_id", query.DistributorID},
		{"service_point.service_type", query.ServiceType},
		{"service_point.point_type", query.PointType},
		{"type", bson.M{"$in": query.MeterType}},
	}, filtersDateInventoryRepositoryMongo(query.Date)...)

	cursor, err := repository.client.Database(repository.Database).Collection(repository.CollectionMeterConfig).CountDocuments(ctxTimeout, match)

	return int(cursor), err

}

func (repository InventoryRepositoryMongo) GetMeterConfigByCups(ctx context.Context, query measures.GetMeterConfigByCupsQuery) (measures.MeterConfig, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	match := append(bson.D{
		{"distributor_id", query.Distributor},
		{"service_point.cups", query.CUPS},
	},
		filtersDateInventoryRepositoryMongo(query.Time)...)

	var config measures.MeterConfig

	result := repository.client.Database(repository.Database).Collection(repository.CollectionMeterConfig).FindOne(ctxTimeout, match)

	err := result.Decode(&config)

	if err != nil {
		return measures.MeterConfig{}, err
	}

	return config, err
}

func (repository InventoryRepositoryMongo) GetMeterConfigByCupsAPI(ctx context.Context, query measures.GetMeterConfigByCupsQuery) (measures.MeterConfig, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	match := append(bson.D{
		{"distributor_id", query.Distributor},
		{"cups", query.CUPS},
	},
		filtersDateInventoryRepositoryMongo(query.Time)...)

	var config measures.MeterConfig

	result := repository.client.Database(repository.Database).Collection(repository.CollectionMeterConfig).FindOne(ctxTimeout, match)

	err := result.Decode(&config)

	if err != nil {
		return measures.MeterConfig{}, err
	}

	return config, err
}

func (repository InventoryRepositoryMongo) GroupMetersByType(ctx context.Context, query measures.GroupMetersByTypeQuery) (map[measures.MeterType]map[measures.RegisterType]measures.MeasureCount, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	database := repository.client.Database(repository.Database)
	collection := database.Collection(repository.CollectionMeterConfig)

	matchStage := bson.D{
		{
			"$match",
			bson.D{
				{"distributor_id", query.DistributorId},
				{"start_date", bson.D{{"$lt", query.EndDate}}},
				{"end_date", bson.D{{"$gt", query.StartDate}}},
				{"contractual_situations.init_date", bson.D{{"$lt", query.EndDate}}},
				{"contractual_situations.end_date", bson.D{{"$gt", query.StartDate}}},
				{"type", bson.M{
					"$ne": nil,
				}},
			},
		},
	}

	groupStage := bson.D{{
		"$group",
		bson.D{
			{"_id",
				bson.D{
					{"type", "$type"},
					{"curve_type", "$curve_type"},
				},
			},
			{"count", bson.D{
				{"$sum", 1},
			}},
			{"serial_numbers", bson.D{
				{"$addToSet", "$serial_number"},
			}},
		},
	}}

	projectStage := bson.D{{
		"$project",
		bson.M{
			"_id":            0,
			"type":           "$_id.type",
			"curve_type":     "$_id.curve_type",
			"serial_numbers": 1,
			"count":          1,
		},
	}}

	type Result struct {
		Type          measures.MeterType    `bson:"type"`
		CurveType     measures.RegisterType `bson:"curve_type"`
		SerialNumbers []string              `bson:"serial_numbers"`
		TotalConfigs  int                   `bson:"count"`
	}

	aggregate := append(bson.A{}, matchStage, groupStage, projectStage)
	cursor, err := collection.Aggregate(ctxTimeout, aggregate)
	defer cursor.Close(ctx)

	if err != nil {
		return map[measures.MeterType]map[measures.RegisterType]measures.MeasureCount{}, err
	}

	mappedResult := make(map[measures.MeterType]map[measures.RegisterType]measures.MeasureCount)

	for cursor.Next(ctx) {
		var r Result
		err = cursor.Decode(&r)
		if err != nil {
			return map[measures.MeterType]map[measures.RegisterType]measures.MeasureCount{}, err
		}

		if mappedResult[r.Type] == nil {
			mappedResult[r.Type] = make(map[measures.RegisterType]measures.MeasureCount)
		}
		mappedResult[r.Type][r.CurveType] = measures.MeasureCount{
			Count: len(r.SerialNumbers),
			Total: r.TotalConfigs,
		}
	}

	return mappedResult, nil
}

func (repository InventoryRepositoryMongo) GetMetersAndCountByDistributorId(ctx context.Context, query measures.GetMetersAndCountByDistributorIdQuery) (measures.GetMetersAndCountResult, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	database := repository.client.Database(repository.Database)
	collection := database.Collection(repository.CollectionMeterConfig)
	matchStage := bson.D{
		{"$match",
			bson.D{
				{"distributor_id", query.DistributorId},
				{"type", query.Type},
				{"start_date", bson.D{{"$lt", query.EndDate}}},
				{"end_date", bson.D{{"$gt", query.StartDate}}},
				{"contractual_situations.init_date", bson.D{{"$lt", query.EndDate}}},
				{"contractual_situations.end_date", bson.D{{"$gt", query.StartDate}}},
			},
		},
	}

	sortStage := bson.D{
		{
			"$sort", bson.D{
				{
					"end_date", 1,
				},
			},
		},
	}

	groupStage := bson.D{
		{
			"$group",
			bson.D{
				{"_id", "$cups"},
				{"meters",
					bson.D{
						{"$push",
							bson.D{
								{"curve_type", "$curve_type"},
								{"end_date", "$end_date"},
								{"start_date", "$start_date"},
							},
						},
					},
				},
			},
		},
	}

	aggregateResult := append(bson.A{}, matchStage, sortStage, groupStage)
	aggregateCount := append(bson.A{}, matchStage, bson.D{
		{
			"$group",
			bson.D{
				{"_id", "$cups"},
			},
		},
	})

	offset := int(query.Offset)

	result, count, err := mongo_utils.Paginate[measures.GetMetersAndCountData](ctxTimeout, collection, mongo_utils.PaginateQuery{
		ResultQuery: mongo_utils.AppendAggregate(aggregateResult),
		CountQuery:  mongo_utils.AppendAggregate(aggregateCount),
		Paginate: db.Pagination{
			Limit:  int(query.Limit),
			Offset: &offset,
		},
	})

	if err != nil {
		return measures.GetMetersAndCountResult{}, err
	}

	return measures.GetMetersAndCountResult{
		Data:  result,
		Count: count,
	}, nil

}
