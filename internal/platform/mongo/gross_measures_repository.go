package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"sync"
	"time"
)

func NewGrossMeasureRepositoryMongo(client *mongo.Client, database, collection string, loc *time.Location) *GrossMeasureRepositoryMongo {

	return &GrossMeasureRepositoryMongo{
		client:          client,
		location:        loc,
		Database:        database,
		CollectionCurve: fmt.Sprintf("%s_curve", collection),
		CollectionClose: fmt.Sprintf("%s_close", collection),
	}
}

type GrossMeasureRepositoryMongo struct {
	client          *mongo.Client
	location        *time.Location
	Database        string
	CollectionCurve string
	CollectionClose string
}

func (repository GrossMeasureRepositoryMongo) SaveAllMeasuresClose(ctx context.Context, measures []gross_measures.MeasureCloseWrite) error {

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	coll := repository.client.Database(repository.Database).Collection(repository.CollectionClose)

	var models []mongo.WriteModel

	for _, m := range measures {
		updateOneModel := mongo.NewUpdateOneModel()
		updateOneModel.SetFilter(bson.D{
			{"_id", m.Id},
		})
		updateOneModel.SetUpsert(true)
		updateOneModel.SetUpdate(bson.D{
			{"$set", m},
		})
		models = append(models, updateOneModel)
	}

	_, err := coll.BulkWrite(ctxTimeout, models)

	return err
}

func (repository GrossMeasureRepositoryMongo) SaveAllMeasuresCurve(ctx context.Context, measures []gross_measures.MeasureCurveWrite) error {

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	coll := repository.client.Database(repository.Database).Collection(repository.CollectionCurve)

	var models []mongo.WriteModel

	for _, m := range measures {
		updateOneModel := mongo.NewUpdateOneModel()
		updateOneModel.SetFilter(bson.D{
			{"_id", m.Id},
		})
		updateOneModel.SetUpsert(true)
		updateOneModel.SetUpdate(bson.D{
			{"$set", m},
		})
		models = append(models, updateOneModel)
	}

	_, err := coll.BulkWrite(ctxTimeout, models)

	return err
}

type GetBestMeasuresInDateQuery struct {
	DistributorID string    `json:"distributor_code"`
	MeterID       string    `json:"meter_id"`
	Date          time.Time `json:"date"`
	ReadingType   string    `json:"reading_type"`
}

func (repository GrossMeasureRepositoryMongo) ListDailyCurveMeasures(ctx context.Context, query gross_measures.QueryListForProcessCurve) ([]gross_measures.MeasureCurveWrite, error) {
	collectionCurve := repository.client.Database(repository.Database).Collection(repository.CollectionCurve)
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	startDate := query.Date.Truncate(time.Hour)
	endDate := query.Date.AddDate(0, 0, 1)

	matchMap := bson.M{
		"end_date": bson.M{
			"$gt":  startDate,
			"$lte": endDate,
		},
		"meter_serial_number": query.SerialNumber,
	}

	if query.CurveType != "" {
		matchMap["curve_type"] = query.CurveType
	}

	matchStage := bson.D{
		{"$match",
			matchMap,
		},
	}
	aggregatePipeline := bson.A{}
	projectStage := bson.D{
		{
			"$project",
			bson.M{
				"_id":                 1,
				"start_date":          1,
				"end_date":            1,
				"reading_date":        1,
				"curve_type":          1,
				"type":                1,
				"status":              1,
				"reading_type":        1,
				"contract":            1,
				"register_type":       1,
				"meter_id":            1,
				"meter_serial_number": 1,
				"measure_point_type":  1,
				"concentrator_id":     1,
				"reader_id":           1,
				"file":                1,
				"distributor_id":      1,
				"origin":              1,
				"qualifier":           1,
				"AI":                  1,
				"AE":                  1,
				"R1":                  1,
				"R2":                  1,
				"R3":                  1,
				"R4":                  1,
				"invalidations":       1,
				"origin_sort": bson.M{
					"$switch": bson.M{
						"branches": bson.A{
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$origin", "Manual"},
								},
								"then": 1,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$origin", "STG"},
								},
								"then": 2,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$origin", "STM"},
								},
								"then": 3,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$origin", "TPL"},
								},
								"then": 4,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$origin", "Visual"},
								},
								"then": 5,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$origin", "Autolectura"},
								},
								"then": 6,
							},
						},
						"default": 7,
					},
				},
				"measure_point_type_sort": bson.M{
					"$switch": bson.M{
						"branches": bson.A{
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$measure_point_type", "P"},
								},
								"then": 1,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$measure_point_type", "R"},
								},
								"then": 2,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$measure_point_type", "C"},
								},
								"then": 3,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$measure_point_type", "T"},
								},
								"then": 4,
							},
						},
						"default": 5,
					},
				},
			},
		},
	}

	sortStage := bson.D{
		{"$sort",
			bson.M{
				"end_date":                1,
				"origin_sort":             1,
				"measure_point_type_sort": 1,
			},
		},
	}
	aggregatePipeline = append(aggregatePipeline, matchStage, projectStage, sortStage)

	cursor, err := collectionCurve.Aggregate(ctxTimeout, aggregatePipeline)
	var rows []gross_measures.MeasureCurveWrite

	if err != nil {
		return rows, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}
func (repository GrossMeasureRepositoryMongo) CountGrossMeasuresFromGenerationDate(ctx context.Context, query gross_measures.QueryListForProcessCurveGenerationDate) (int, error) {
	var collection *mongo.Collection
	switch query.ReadingType {
	case measures.Curve:
		collection = repository.client.Database(repository.Database).Collection(repository.CollectionCurve)
	case measures.DailyClosure:
		collection = repository.client.Database(repository.Database).Collection(repository.CollectionClose)
	case measures.BillingClosure:
		collection = repository.client.Database(repository.Database).Collection(repository.CollectionClose)

	}
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	match := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"distributor_id", query.DistributorId},
					{"reading_type", query.ReadingType},
					{"generation_date",
						bson.D{
							{"$gt", query.StartDate},
							{"$lt", query.EndDate},
						},
					},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"meter_serial_number", 1},
					{"day", bson.D{{"$dayOfMonth", "$end_date"}}},
					{"month", bson.D{{"$month", "$end_date"}}},
					{"year", bson.D{{"$year", "$end_date"}}},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id",
						bson.D{
							{"day", "$day"},
							{"month", "$month"},
							{"year", "$year"},
							{"meter_serial_number", "$meter_serial_number"},
						},
					},
				},
			},
		},
		bson.D{{"$count", "meter_serial_number"}},
	}
	cursor, err := collection.Aggregate(ctxTimeout, match)
	if err != nil {
		return 0, err
	}
	res := []struct {
		MeterSerialNumber int `json:"meter_serial_number" bson:"meter_serial_number"`
	}{}
	err = cursor.All(ctx, &res)
	if err != nil {
		return 0, err
	}

	return res[0].MeterSerialNumber, err
}
func (repository GrossMeasureRepositoryMongo) ListGrossMeasuresFromGenerationDate(ctx context.Context, query gross_measures.QueryListForProcessCurveGenerationDate) ([]gross_measures.MeasureCurveMeterSerialNumber, error) {
	var collection *mongo.Collection
	switch query.ReadingType {
	case measures.Curve:
		collection = repository.client.Database(repository.Database).Collection(repository.CollectionCurve)
	case measures.DailyClosure:
		collection = repository.client.Database(repository.Database).Collection(repository.CollectionClose)
	case measures.BillingClosure:
		collection = repository.client.Database(repository.Database).Collection(repository.CollectionClose)

	}
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	aggregatePipeline := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"distributor_id", query.DistributorId},
					{"reading_type", query.ReadingType},
					{"generation_date",
						bson.D{
							{"$gt", query.StartDate},
							{"$lt", query.EndDate},
						},
					},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"meter_serial_number", 1},
					{"day", bson.D{{"$dayOfMonth", "$end_date"}}},
					{"month", bson.D{{"$month", "$end_date"}}},
					{"year", bson.D{{"$year", "$end_date"}}},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id",
						bson.D{
							{"day", "$day"},
							{"month", "$month"},
							{"year", "$year"},
							{"meter_serial_number", "$meter_serial_number"},
						},
					},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"meter_serial_number", "$_id.meter_serial_number"},
					{"day", "$_id.day"},
					{"month", "$_id.month"},
					{"year", "$_id.year"},
				},
			},
		},
		bson.D{
			{"$sort",
				bson.D{
					{"day", 1},
					{"month", 1},
					{"year", 1},
					{"meter_serial_number", 1},
				},
			},
		},
		bson.D{{"$skip", query.Offset}},
		bson.D{{"$limit", query.Limit}},
	}
	cursor, err := collection.Aggregate(ctxTimeout, aggregatePipeline)
	defer cursor.Close(ctx)
	var rows []gross_measures.MeasureCurveMeterSerialNumber

	if err != nil {
		return rows, err
	}

	if err = cursor.All(ctx, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func (repository GrossMeasureRepositoryMongo) ListDailyCloseMeasures(ctx context.Context, query gross_measures.QueryListForProcessClose) ([]gross_measures.MeasureCloseWrite, error) {
	collectionCurve := repository.client.Database(repository.Database).Collection(repository.CollectionClose)
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	startDate := query.Date
	endDate := query.Date.AddDate(0, 0, 1)

	matchStage := bson.D{
		{"$match",
			bson.M{
				"end_date": bson.M{
					"$gt":  startDate.UTC(),
					"$lte": endDate.UTC(),
				},
				"meter_serial_number": query.SerialNumber,
				"reading_type":        query.ReadingType,
			},
		},
	}
	aggregatePipeline := bson.A{}
	projectStage := bson.D{
		{
			"$project",
			bson.M{
				"_id":                 1,
				"start_date":          1,
				"end_date":            1,
				"reading_date":        1,
				"type":                1,
				"status":              1,
				"reading_type":        1,
				"contract":            1,
				"register_type":       1,
				"meter_id":            1,
				"meter_serial_number": 1,
				"measure_point_type":  1,
				"concentrator_id":     1,
				"reader_id":           1,
				"file":                1,
				"distributor_id":      1,
				"origin":              1,
				"qualifier":           1,
				"periods":             1,
				"invalidations":       1,
				"origin_sort": bson.M{
					"$switch": bson.M{
						"branches": bson.A{
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$origin", "Manual"},
								},
								"then": 1,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$origin", "STG"},
								},
								"then": 2,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$origin", "STM"},
								},
								"then": 3,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$origin", "TPL"},
								},
								"then": 4,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$origin", "Visual"},
								},
								"then": 5,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$origin", "Autolectura"},
								},
								"then": 6,
							},
						},
						"default": 7,
					},
				},
				"measure_point_type_sort": bson.M{
					"$switch": bson.M{
						"branches": bson.A{
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$measure_point_type", "P"},
								},
								"then": 1,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$measure_point_type", "R"},
								},
								"then": 2,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$measure_point_type", "C"},
								},
								"then": 3,
							},
							bson.M{
								"case": bson.M{
									"$eq": bson.A{"$measure_point_type", "T"},
								},
								"then": 4,
							},
						},
						"default": 5,
					},
				},
			},
		},
	}

	sortStage := bson.D{
		{"$sort",
			bson.M{
				"end_date":                1,
				"origin_sort":             1,
				"measure_point_type_sort": 1,
			},
		},
	}
	aggregatePipeline = append(aggregatePipeline, matchStage, projectStage, sortStage)

	cursor, err := collectionCurve.Aggregate(ctxTimeout, aggregatePipeline)
	var rows []gross_measures.MeasureCloseWrite

	if err != nil {
		return rows, nil
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &rows); err != nil {
		return nil, err
	}

	return rows, nil
}

func (repository GrossMeasureRepositoryMongo) ListCurveMeasures(ctx context.Context, query gross_measures.QueryListMeasure) ([]gross_measures.MeasureCurveWrite, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	collection := repository.client.Database(repository.Database).Collection(repository.CollectionCurve)

	match := bson.D{
		{
			"meter_serial_number", query.SerialNumber,
		},
		{
			"end_date", bson.M{
				"$gt":  query.StartDate,
				"$lte": query.EndDate,
			},
		},
	}

	cursor, err := collection.Find(ctxTimeout, match)

	if err != nil {
		return []gross_measures.MeasureCurveWrite{}, err
	}

	var grossCurves []gross_measures.MeasureCurveWrite
	if err = cursor.All(ctx, &grossCurves); err != nil {
		return []gross_measures.MeasureCurveWrite{}, err
	}

	return grossCurves, nil
}

func (repository GrossMeasureRepositoryMongo) ListCloseMeasures(ctx context.Context, query gross_measures.QueryListMeasure) ([]gross_measures.MeasureCloseWrite, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	collection := repository.client.Database(repository.Database).Collection(repository.CollectionClose)

	match := bson.D{
		{
			"meter_serial_number", query.SerialNumber,
		},
		{
			"end_date", bson.M{
				"$gt":  query.StartDate,
				"$lte": query.EndDate,
			},
		},
	}

	if query.ReadingType != "" {
		match = append(match, bson.E{
			Key:   "reading_type",
			Value: query.ReadingType,
		})
	}

	cursor, err := collection.Find(ctxTimeout, match)

	if err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}

	var grossClose []gross_measures.MeasureCloseWrite
	if err = cursor.All(ctx, &grossClose); err != nil {
		return []gross_measures.MeasureCloseWrite{}, err
	}

	return grossClose, nil
}

type measureDashboardRow struct {
	Count int `bson:"count"`
	Id    struct {
		Origin      string `bson:"origin"`
		TypeMeasure string `bson:"reading_type"`
		Status      string `bson:"status"`
		Date        string `bson:"date"`
	} `bson:"_id"`
}

func (repository GrossMeasureRepositoryMongo) getGrossMeasureDashboardRow(ctx context.Context, query gross_measures.GetDashboardQuery, rows *[]measureDashboardRow, wg *sync.WaitGroup, collectionName string) {
	collectionCurve := repository.client.Database(repository.Database).Collection(collectionName)
	defer wg.Done()
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	matchStage := bson.D{
		{"$match",
			bson.M{
				"end_date": bson.M{
					"$gte": query.StartDate,
					"$lt":  query.EndDate,
				},
				"distributor_id": query.DistributorId,
			},
		},
	}
	aggregatePipeline := bson.A{}
	groupStage := bson.D{
		{
			"$group",
			bson.M{
				"_id": bson.M{
					"origin":       "$origin",
					"reading_type": "$reading_type",
					"status":       "$status",
					"date":         bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$end_date", "timezone": repository.location.String()}},
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	}
	aggregatePipeline = append(aggregatePipeline, matchStage, groupStage)

	cursor, err := collectionCurve.Aggregate(ctxTimeout, aggregatePipeline)

	if err != nil {
		return
	}

	if err = cursor.All(ctx, rows); err != nil {
		log.Fatal(err)
	}
}

func (repository GrossMeasureRepositoryMongo) GetDashboard(ctx context.Context, query gross_measures.GetDashboardQuery) ([]measures.DashboardMeasureI, error) {
	rowsCurve := make([]measureDashboardRow, 0)
	rowsClose := make([]measureDashboardRow, 0)

	var wg sync.WaitGroup
	wg.Add(2)

	go repository.getGrossMeasureDashboardRow(ctx, query, &rowsCurve, &wg, repository.CollectionCurve)
	go repository.getGrossMeasureDashboardRow(ctx, query, &rowsClose, &wg, repository.CollectionClose)
	wg.Wait()

	rows := make([]measureDashboardRow, 0, cap(rowsCurve)+cap(rowsClose))
	rows = append(rows, rowsCurve...)
	rows = append(rows, rowsClose...)

	resultDashboard := make([]measures.DashboardMeasureI, 0, cap(rows))

	for _, e := range rows {
		t, _ := time.ParseInLocation("2006-01-02", e.Id.Date, repository.location)
		resultDashboard = append(resultDashboard, measures.NewDashboardMeasure(
			e.Count,
			measures.OriginType(e.Id.Origin),
			measures.ReadingType(e.Id.TypeMeasure),
			measures.Status(e.Id.Status),
			t,
		))
	}

	return resultDashboard, nil
}
