package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
)

func NewBillingMeasureRepositoryMongo(client *mongo.Client, database, collection string) *BillingMeasureRepositoryMongo {

	return &BillingMeasureRepositoryMongo{
		client:     client,
		Database:   database,
		Collection: collection,
	}
}

type BillingMeasureRepositoryMongo struct {
	client     *mongo.Client
	Database   string
	Collection string
}

func (repository BillingMeasureRepositoryMongo) Last(ctx context.Context, q billing_measures.QueryLast) (billing_measures.BillingMeasure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	lastDayPreviousMonth := q.Date.AddDate(0, 0, -q.Date.Day())

	coll := repository.client.Database(repository.Database).Collection(repository.Collection)
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"cups", q.CUPS},
					{"status", bson.M{
						"$in": bson.A{billing_measures.ReadyToBill, billing_measures.Billed},
					}},
					{"end_date", bson.D{
						{"$lte", q.Date},
						{"$gte", lastDayPreviousMonth},
					}},
				},
			},
		},
	}
	opts := options.FindOne().SetSort(bson.D{{"end_date", -1}})

	var b billing_measures.BillingMeasure

	result := coll.FindOne(ctxTimeout, filter, opts)

	err := result.Decode(&b)

	if err != nil {
		return billing_measures.BillingMeasure{}, err
	}

	return b, err
}

func (repository BillingMeasureRepositoryMongo) GetPrevious(ctx context.Context, query billing_measures.GetPrevious) (billing_measures.BillingMeasure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.Collection)
	filter := bson.D{
		{"cups", query.CUPS},
		{"init_date", query.InitDate},
		{"end_date", query.EndDate},
	}
	opts := options.FindOne().SetSort(bson.D{{"generation_date", -1}})

	var b billing_measures.BillingMeasure

	result := coll.FindOne(ctxTimeout, filter, opts)

	err := result.Decode(&b)

	if err != nil {
		return billing_measures.BillingMeasure{}, err
	}

	return b, err
}

func (repository BillingMeasureRepositoryMongo) Save(ctx context.Context, measure billing_measures.BillingMeasure) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.Collection)
	filter := bson.D{
		{"_id", measure.Id},
	}

	upsert := true
	_, err := coll.UpdateOne(ctxTimeout, filter, bson.D{
		{"$set", measure}}, &options.UpdateOptions{
		Upsert: &upsert,
	})

	return err
}

func (repository BillingMeasureRepositoryMongo) SaveAll(ctx context.Context, measures []billing_measures.BillingMeasure) error {

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		g, ctxErrGroup := errgroup.WithContext(sessCtx)
		for _, measure := range measures {
			measure := measure
			g.Go(func() error {
				return repository.Save(ctxErrGroup, measure)
			})

		}
		if err := g.Wait(); err != nil {
			return nil, err

		}
		return nil, nil

	}

	session, err := repository.client.StartSession()

	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, callback)

	if err != nil {
		return err
	}

	return nil
}

func (repository BillingMeasureRepositoryMongo) LastHistory(ctx context.Context, q billing_measures.QueryLastHistory) (billing_measures.BillingMeasure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	collection := repository.client.Database(repository.Database).Collection(repository.Collection)
	initDate := q.InitDate.AddDate(-1, 0, 0)
	endDate := q.EndDate.AddDate(-1, 0, 0)
	filtersKeys := bson.D{
		{"cups", q.CUPS},
		{"init_date",
			bson.M{
				"$eq": initDate,
			},
		},
		{"end_date",
			bson.M{
				"$eq": endDate,
			},
		},
		{
			"status", bson.M{
				"$in": bson.A{billing_measures.ReadyToBill, billing_measures.Billed},
			},
		},
	}

	filtersKeys = append(filtersKeys, repository.getFilterBillingBalanceReal(q.Periods, q.Magnitudes)...)

	filter := bson.D{{"$and", bson.A{filtersKeys}}}

	var billingMeasure billing_measures.BillingMeasure

	result := collection.FindOne(ctxTimeout, filter)

	err := result.Decode(&billingMeasure)

	if err != nil {
		return billing_measures.BillingMeasure{}, err
	}

	return billingMeasure, nil
}

func (repository BillingMeasureRepositoryMongo) getFilterBillingBalanceReal(periods []measures.PeriodKey, magnitudes []measures.Magnitude) bson.D {
	filtersKeys := bson.D{}
	for _, period := range periods {
		periodName := strings.ToLower(string(period))
		for _, magnitude := range magnitudes {
			magnitudeName := strings.ToLower(string(magnitude))
			filterName := fmt.Sprintf("billing_balance.%s.balance_type_%s", periodName, magnitudeName)
			filtersKeys = append(filtersKeys, bson.E{Key: filterName, Value: bson.M{
				"$in": billing_measures.RealBalanceTypes,
			}})
		}
	}

	return filtersKeys
}

func (repository BillingMeasureRepositoryMongo) GetCloseHistories(ctx context.Context, q billing_measures.QueryGetCloseHistories) ([]billing_measures.BillingMeasure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	collection := repository.client.Database(repository.Database).Collection(repository.Collection)
	filter := bson.D{
		{"cups", q.CUPS},
		{"status", bson.M{
			"$in": bson.A{billing_measures.ReadyToBill, billing_measures.Billed},
		}}}

	filter = append(filter, repository.getFilterBillingBalanceReal(q.Periods, q.Magnitudes)...)

	aggregatePipeline := bson.A{}
	filterStage := bson.D{
		{"$and",
			bson.A{
				filter,
			},
		},
	}
	addFieldStage := bson.D{
		{
			"$addFields",
			bson.D{
				{
					"diffDate", bson.M{
						"$abs": bson.M{
							"$subtract": bson.A{"$end_date", q.EndDate},
						},
					},
				},
			},
		},
	}

	sortStage := bson.D{
		{
			"$sort",
			bson.M{
				"diffDate": 1,
			},
		},
	}

	limitStage := bson.D{
		{
			"$limit",
			4,
		},
	}

	aggregate := append(aggregatePipeline, filterStage, addFieldStage, sortStage, limitStage)
	var billingMeasure []billing_measures.BillingMeasure

	cursor, err := collection.Aggregate(ctxTimeout, aggregate)

	if err != nil && len(billingMeasure) == 0 {
		return []billing_measures.BillingMeasure{}, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &billingMeasure); err != nil {
		return nil, err
	}

	return billingMeasure, nil
}

func (repository BillingMeasureRepositoryMongo) Find(ctx context.Context, id string) (billing_measures.BillingMeasure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	collection := repository.client.Database(repository.Database).Collection(repository.Collection)

	filter := bson.D{
		{"_id", id},
	}

	var billingMeasure billing_measures.BillingMeasure

	result := collection.FindOne(ctxTimeout, filter)

	err := result.Decode(&billingMeasure)

	if err != nil {
		return billing_measures.BillingMeasure{}, err
	}

	return billingMeasure, nil
}
