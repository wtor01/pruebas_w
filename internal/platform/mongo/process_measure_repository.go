package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

func NewProcessMeasureRepository(client *mongo.Client, database string, collectionLoadCurve string, collectionDailyLoadCurve string, collectionMonthlyClosure string, collectionDailyClosure string, loc *time.Location) *ProcessMeasureRepository {
	return &ProcessMeasureRepository{
		client:                   client,
		location:                 loc,
		Database:                 database,
		CollectionLoadCurve:      collectionLoadCurve,
		CollectionDailyLoadCurve: collectionDailyLoadCurve,
		CollectionMonthlyClosure: collectionMonthlyClosure,
		CollectionDailyClosure:   collectionDailyClosure,
	}
}

type ProcessMeasureRepository struct {
	client                   *mongo.Client
	location                 *time.Location
	Database                 string
	CollectionLoadCurve      string
	CollectionDailyLoadCurve string
	CollectionMonthlyClosure string
	CollectionDailyClosure   string
}

func (repository ProcessMeasureRepository) GetMonthlyClosureByID(ctx context.Context, id string) (process_measures.ProcessedMonthlyClosure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.CollectionMonthlyClosure)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return process_measures.ProcessedMonthlyClosure{}, err
	}
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"_id", objectID},
				},
			},
		},
	}

	var b process_measures.ProcessedMonthlyClosure

	result := coll.FindOne(ctxTimeout, filter)

	if err != nil {
		return process_measures.ProcessedMonthlyClosure{}, err
	}

	err = result.Decode(&b)

	if err != nil {
		return process_measures.ProcessedMonthlyClosure{}, err
	}

	return b, err
}

func (repository ProcessMeasureRepository) GetDailyClosureByID(ctx context.Context, id string) (process_measures.ProcessedDailyClosure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.CollectionDailyClosure)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return process_measures.ProcessedDailyClosure{}, err
	}
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"_id", objectID},
				},
			},
		},
	}

	var b process_measures.ProcessedDailyClosure

	result := coll.FindOne(ctxTimeout, filter)

	if err != nil {
		return process_measures.ProcessedDailyClosure{}, err
	}

	err = result.Decode(&b)

	if err != nil {
		return process_measures.ProcessedDailyClosure{}, err
	}

	return b, err
}

func (repository ProcessMeasureRepository) GetLoadCurveByID(ctx context.Context, id string) (process_measures.ProcessedLoadCurve, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.CollectionLoadCurve)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return process_measures.ProcessedLoadCurve{}, err
	}
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"_id", objectID},
				},
			},
		},
	}

	var b process_measures.ProcessedLoadCurve

	result := coll.FindOne(ctxTimeout, filter)

	if err != nil {
		return process_measures.ProcessedLoadCurve{}, err
	}

	err = result.Decode(&b)

	if err != nil {
		return process_measures.ProcessedLoadCurve{}, err
	}

	return b, err
}

func (repository ProcessMeasureRepository) ProcessedLoadCurveByCups(ctx context.Context, q process_measures.QueryProcessedLoadCurveByCups) ([]process_measures.ProcessedLoadCurve, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.CollectionLoadCurve)
	filter := bson.D{
		{"cups", q.CUPS},
		{
			"end_date", bson.M{
				"$gt":  q.StartDate,
				"$lte": q.EndDate,
			},
		},
	}

	if q.Status != "" {
		filter = append(filter, bson.D{{"validation_status", q.Status}}...)
	}

	if q.CurveType != "" {
		filter = append(filter, bson.D{{"curve_type", q.CurveType}}...)
	}

	opts := options.Find().SetSort(bson.D{{"end_date", 1}})

	var b []process_measures.ProcessedLoadCurve

	result, err := coll.Find(ctxTimeout, filter, opts)

	if err != nil {
		return make([]process_measures.ProcessedLoadCurve, 0), err
	}

	err = result.All(ctxTimeout, &b)

	if err != nil {
		return make([]process_measures.ProcessedLoadCurve, 0), err
	}

	return b, err
}

func (repository ProcessMeasureRepository) SaveProcessedDailyLoadCurve(ctx context.Context, measure process_measures.ProcessedDailyLoadCurve) error {

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	coll := repository.client.Database(repository.Database).Collection(repository.CollectionDailyLoadCurve)

	upsert := true
	_, err := coll.UpdateOne(
		ctxTimeout,
		bson.D{{"_id", measure.Id}},
		bson.D{{"$set", measure}},
		&options.UpdateOptions{Upsert: &upsert},
	)

	return err
}

func (repository ProcessMeasureRepository) SaveMonthlyClosure(ctx context.Context, measure process_measures.ProcessedMonthlyClosure) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	coll := repository.client.Database(repository.Database).Collection(repository.CollectionMonthlyClosure)

	upsert := true
	_, err := coll.UpdateOne(
		ctxTimeout,
		bson.D{{"_id", measure.Id}},
		bson.D{{"$set", measure}},
		&options.UpdateOptions{Upsert: &upsert},
	)

	return err
}

func (repository ProcessMeasureRepository) SaveDailyClosure(ctx context.Context, measure process_measures.ProcessedDailyClosure) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	coll := repository.client.Database(repository.Database).Collection(repository.CollectionDailyClosure)

	upsert := true
	_, err := coll.UpdateOne(
		ctxTimeout,
		bson.D{{"_id", measure.Id}},
		bson.D{{"$set", measure}},
		&options.UpdateOptions{Upsert: &upsert},
	)

	return err
}

func (repository ProcessMeasureRepository) SaveAllProcessedLoadCurve(ctx context.Context, measures []process_measures.ProcessedLoadCurve) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	coll := repository.client.Database(repository.Database).Collection(repository.CollectionLoadCurve)

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

func (repository ProcessMeasureRepository) getProcessMeasureDashboardRow(ctx context.Context, query process_measures.GetDashboardQuery, rows *[]measureDashboardRow, wg *sync.WaitGroup, collectionName string) {
	collection := repository.client.Database(repository.Database).Collection(collectionName)
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
				"origin":         bson.M{"$ne": measures.Filled},
			},
		},
	}
	aggregatePipeline := bson.A{}
	readingTypes := map[string]measures.ReadingType{
		repository.CollectionMonthlyClosure: measures.BillingClosure,
		repository.CollectionDailyClosure:   measures.DailyClosure,
		repository.CollectionLoadCurve:      measures.Curve,
	}

	readingType := readingTypes[collectionName]

	groupStage := bson.D{
		{
			"$group",
			bson.M{
				"_id": bson.M{
					"origin":       "$origin",
					"reading_type": readingType,
					"status":       "$validation_status",
					"date":         bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$end_date", "timezone": repository.location.String()}},
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	}
	aggregatePipeline = append(aggregatePipeline, matchStage, groupStage)

	cursor, err := collection.Aggregate(ctxTimeout, aggregatePipeline)

	if err != nil {
		return
	}

	if err = cursor.All(ctx, rows); err != nil {
		log.Fatal(err)
	}
}

func (repository ProcessMeasureRepository) GetDashboard(ctx context.Context, query process_measures.GetDashboardQuery) ([]measures.DashboardMeasureI, error) {
	rowsCurve := make([]measureDashboardRow, 0)
	rowsDailyClose := make([]measureDashboardRow, 0)
	rowsBillingClose := make([]measureDashboardRow, 0)

	var wg sync.WaitGroup
	wg.Add(3)

	go repository.getProcessMeasureDashboardRow(ctx, query, &rowsCurve, &wg, repository.CollectionLoadCurve)
	go repository.getProcessMeasureDashboardRow(ctx, query, &rowsBillingClose, &wg, repository.CollectionMonthlyClosure)
	go repository.getProcessMeasureDashboardRow(ctx, query, &rowsDailyClose, &wg, repository.CollectionDailyClosure)
	wg.Wait()

	rows := make([]measureDashboardRow, 0, cap(rowsCurve)+cap(rowsBillingClose)+cap(rowsDailyClose))
	rows = append(rows, rowsCurve...)
	rows = append(rows, rowsBillingClose...)
	rows = append(rows, rowsDailyClose...)

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

func (repository ProcessMeasureRepository) GetMonthlyClosureByCup(ctx context.Context, q process_measures.QueryClosedCupsMeasureOnDate) (process_measures.ProcessedMonthlyClosure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.CollectionMonthlyClosure)
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"cups", q.CUPS},
					{"validation_status", measures.Valid},
					{"end_date", q.Date.UTC()},
				},
			},
		},
	}
	opts := options.FindOne()

	var measure process_measures.ProcessedMonthlyClosure

	result := coll.FindOne(ctxTimeout, filter, opts)

	err := result.Decode(&measure)

	if err != nil {
		return process_measures.ProcessedMonthlyClosure{}, err
	}

	return measure, err
}
func (repository ProcessMeasureRepository) GetMonthlyClosureMeasuresByCup(ctx context.Context, q process_measures.QueryMonthlyClosedMeasures) ([]process_measures.ProcessedMonthlyClosure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.CollectionMonthlyClosure)
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"cups", q.CUPS},
					{"end_date", bson.M{
						"$gt":  q.StartDate,
						"$lte": q.EndDate,
					}},
				},
			},
		},
	}
	opts := options.Find().SetSort(bson.D{{"end_date", 1}})

	monthlyMeasures := make([]process_measures.ProcessedMonthlyClosure, 0)

	result, err := coll.Find(ctxTimeout, filter, opts)

	if err != nil {
		return make([]process_measures.ProcessedMonthlyClosure, 0), err
	}

	err = result.All(ctxTimeout, &monthlyMeasures)

	if err != nil {
		return make([]process_measures.ProcessedMonthlyClosure, 0), err
	}

	return monthlyMeasures, err
}

func (repository ProcessMeasureRepository) GetProcessedDailyClosureByCup(ctx context.Context, q process_measures.QueryClosedCupsMeasureOnDate) (process_measures.ProcessedDailyClosure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.CollectionDailyClosure)
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"cups", q.CUPS},
					{"validation_status", measures.Valid},
					{"end_date", q.Date.UTC()},
				},
			},
		},
	}
	opts := options.FindOne()

	var measure process_measures.ProcessedDailyClosure

	result := coll.FindOne(ctxTimeout, filter, opts)

	err := result.Decode(&measure)

	if err != nil {
		return process_measures.ProcessedDailyClosure{}, err
	}

	return measure, err
}
func (repository ProcessMeasureRepository) ProcessedDailyClosureByCups(ctx context.Context, q process_measures.QueryHistoryDailyClosureByCups) ([]process_measures.ProcessedDailyClosure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.CollectionDailyClosure)
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"cups", q.CUPS},
					{
						"end_date", bson.M{
							"$gt":  q.StartDate,
							"$lte": q.EndDate,
						},
					},
				},
			},
		},
	}
	opts := options.Find().SetSort(bson.D{{"end_date", 1}})

	var b []process_measures.ProcessedDailyClosure

	result, err := coll.Find(ctxTimeout, filter, opts)

	if err != nil {
		return make([]process_measures.ProcessedDailyClosure, 0), err
	}

	err = result.All(ctxTimeout, &b)

	if err != nil {
		return make([]process_measures.ProcessedDailyClosure, 0), err
	}

	return b, err
}

func (repository ProcessMeasureRepository) getProcessMeasureCupsList(ctx context.Context, wg *sync.WaitGroup, query process_measures.ListCupsQuery, collectionName string, dashboardMeasures map[string]*measures.DashboardCupsReading) {
	defer wg.Done()
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	database := repository.client.Database(repository.Database)
	collection := database.Collection(collectionName)

	readingTypes := map[string]measures.ReadingType{
		repository.CollectionMonthlyClosure: measures.BillingClosure,
		repository.CollectionDailyClosure:   measures.DailyClosure,
		repository.CollectionLoadCurve:      measures.Curve,
	}

	readingType := readingTypes[collectionName]

	matchStage := bson.D{
		{
			"$match",
			bson.M{
				"cups": bson.M{
					"$in": query.Cups,
				},
				"end_date": bson.M{
					"$gt": query.StartDate,
					"$lt": query.EndDate,
				},
				"origin": bson.M{
					"$ne": measures.Filled,
				},
			},
		},
	}

	groupByCupsStatus := bson.D{
		{
			"$group",
			bson.M{
				"_id": bson.M{
					"cups":   "$cups",
					"status": "$validation_status",
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	groupByCups := bson.D{
		{"$group",
			bson.M{
				"_id": bson.M{
					"cups": "$_id.cups",
				},
				"total_count": bson.M{
					"$sum": "$count",
				},
				"status": bson.D{
					{
						"$push", bson.M{
							"status": "$_id.status",
							"count":  "$count",
						},
					},
				},
			},
		},
	}

	projectStage := bson.D{
		{
			"$project", bson.M{
				"_id":         0,
				"cups":        "$_id.cups",
				"total_count": 1,
				"status":      1,
			},
		},
	}

	agreggations := append(bson.A{}, matchStage, groupByCupsStatus, groupByCups, projectStage)

	cursor, err := collection.Aggregate(ctxTimeout, agreggations)

	if err != nil {
		return
	}

	type Result struct {
		Cups        string `bson:"cups"`
		TotalCount  int    `bson:"total_count"`
		StatusCount []struct {
			Status measures.Status `bson:"status"`
			Count  int             `bson:"count"`
		} `bson:"status"`
	}

	var result Result

	for cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return
		}
		for _, status := range result.StatusCount {
			dashboardMeasures[result.Cups].SetReading(readingType, status.Status, status.Count)
		}
		dashboardMeasures[result.Cups].SetTotal(readingType, result.TotalCount)

	}
}

func (repository ProcessMeasureRepository) GetCupsMeasures(ctx context.Context, q process_measures.ListCupsQuery) (map[string]*measures.DashboardCupsReading, error) {

	dashboardMeasures := map[string]*measures.DashboardCupsReading{}

	for _, cup := range q.Cups {
		dashboardMeasures[cup] = &measures.DashboardCupsReading{}
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go repository.getProcessMeasureCupsList(ctx, &wg, q, repository.CollectionLoadCurve, dashboardMeasures)
	go repository.getProcessMeasureCupsList(ctx, &wg, q, repository.CollectionDailyClosure, dashboardMeasures)
	go repository.getProcessMeasureCupsList(ctx, &wg, q, repository.CollectionMonthlyClosure, dashboardMeasures)
	wg.Wait()

	return dashboardMeasures, nil
}

func (repository ProcessMeasureRepository) ListHistoryLoadCurve(ctx context.Context, q process_measures.QueryHistoryLoadCurve) ([]process_measures.ProcessedLoadCurve, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.CollectionLoadCurve)

	sortStage := bson.D{
		{"$sort",
			bson.M{
				"diffDate": 1,
			},
		},
	}

	rowsByPeriod := make(map[measures.PeriodKey][]process_measures.ProcessedLoadCurve)
	errorByPeriod := make(map[measures.PeriodKey]error)
	wg := sync.WaitGroup{}

	for _, period := range q.Periods {
		rowsByPeriod[period] = make([]process_measures.ProcessedLoadCurve, 0, 18)
		wg.Add(1)
		go func(rowsByPeriod map[measures.PeriodKey][]process_measures.ProcessedLoadCurve, period measures.PeriodKey) {
			defer wg.Done()

			substractDate := q.EndDate

			if q.IsFuture {
				substractDate = q.StartDate
			}

			fields := append(bson.D{},
				bson.E{
					Key: "diffDate",
					Value: bson.M{
						"$abs": bson.M{
							"$subtract": bson.A{"$end_date", substractDate},
						},
					}})

			filters := bson.M{
				"cups": q.CUPS,
				"origin": bson.M{
					"$ne": measures.Filled,
				},
				"magnitudes": q.Magnitude,
				"period":     period,
			}

			timeFilter := bson.M{}

			if !q.StartDate.IsZero() {
				timeFilter["$gt"] = q.StartDate.In(repository.location).UTC()
			}

			if !q.EndDate.IsZero() {
				timeFilter["$lt"] = q.EndDate.In(repository.location).UTC()
			}

			if len(timeFilter) != 0 {
				filters["end_date"] = timeFilter
			}
			if q.Type != "" {
				filters["curve_type"] = q.Type
			}

			matchStage := bson.D{
				{"$match",
					filters,
				},
			}

			addFieldStage := bson.D{
				{
					"$addFields",
					fields,
				},
			}

			aggregatePipeline := append(bson.A{}, matchStage, addFieldStage, sortStage)

			if q.WithCriterias {
				groupByCriteriaStage := bson.D{
					{"$group",
						bson.D{
							{"_id",
								bson.D{
									{"criteria1", bson.M{
										"$concat": bson.A{"$period", "$season_id", "$day_type_id"},
									}},
									{"criteria2", bson.M{
										"$concat": bson.A{"$period", "$season_id"},
									}},
									{"criteria3", bson.M{
										"$concat": bson.A{"$period"},
									}},
								},
							},
							{"curves", bson.D{
								{"$topN", bson.M{
									"output": "$$ROOT",
									"sortBy": bson.M{
										"diffDate": 1,
									},
									"n": 6,
								}},
							}},
						},
					},
				}

				unwindStage := bson.D{{"$unwind", bson.D{{"path", "$curves"}}}}

				replaceRootStage := bson.D{{"$replaceRoot", bson.D{{"newRoot", "$curves"}}}}

				aggregatePipeline = append(aggregatePipeline, groupByCriteriaStage, unwindStage, replaceRootStage)
			}

			if q.Count != 0 {
				limitStage := bson.D{
					{"$limit",
						q.Count,
					},
				}
				aggregatePipeline = append(aggregatePipeline, limitStage)
			}

			cursor, err := coll.Aggregate(ctxTimeout, aggregatePipeline)

			if err != nil {
				errorByPeriod[period] = err
				return
			}

			defer cursor.Close(ctx)

			rows := rowsByPeriod[period]

			if err = cursor.All(ctx, &rows); err != nil {
				errorByPeriod[period] = err
			}
			rowsByPeriod[period] = rows
		}(rowsByPeriod, period)
	}
	wg.Wait()
	arr := make([]process_measures.ProcessedLoadCurve, 0, len(q.Periods)*18)

	for _, v := range rowsByPeriod {
		if len(v) > 0 {
			arr = append(arr, v...)
		}
	}
	var err error

	for _, e := range errorByPeriod {
		if e != nil {
			err = e
		}
	}

	return arr, err
}
