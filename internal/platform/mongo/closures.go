package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewClosureRepositoryMongo(client *mongo.Client, database string) *ClosureRepositoryMongo {
	return &ClosureRepositoryMongo{
		client:                  client,
		Database:                database,
		CollectionClosureConfig: "processed_monthly_closure",
		CollectionMeterConfig:   "meter_configs",
	}
}

type ClosureRepositoryMongo struct {
	client                  *mongo.Client
	Database                string
	CollectionClosureConfig string
	CollectionMeterConfig   string
}

// GetResume recupera el cierre anterior y el siguiente
func (repository ClosureRepositoryMongo) GetResume(ctx context.Context, query process_measures.GetResume) (process_measures.ResumesProcessMonthlyClosure, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	MonthlyClosures := process_measures.ResumesProcessMonthlyClosure{
		Previous: &process_measures.ProcessedMonthlyClosure{},
		Next:     &process_measures.ProcessedMonthlyClosure{},
	}
	previousMatch := bson.D{
		{"cups", query.Cups},
		{"distributor_id", query.DistributorId},
		{
			"start_date",
			bson.M{"$lte": query.StartDate},
		},
	}

	nextMatch := bson.D{
		{"cups", query.Cups},
		{"distributor_id", query.DistributorId},
		{
			"end_date",
			bson.M{"$gt": query.EndDate},
		},
	}
	call := func(ctx context.Context, match primitive.D, opts *options.FindOneOptions, pmc chan *process_measures.ProcessedMonthlyClosure, error chan error) {
		defer close(pmc)
		defer close(error)
		result := repository.client.Database(repository.Database).Collection(repository.CollectionClosureConfig).
			FindOne(ctxTimeout, match, opts)
		res := process_measures.ProcessedMonthlyClosure{}
		err := result.Decode(&res)
		if err != nil {
			error <- err
			pmc <- nil
			return
		}
		error <- nil
		pmc <- &res
	}
	errPrevious := make(chan error)
	errNext := make(chan error)

	optsPrevious := options.FindOne().SetSort(bson.D{{"end_date", -1}})
	optsNext := options.FindOne().SetSort(bson.D{{"end_date", 1}})
	previousPMC := make(chan *process_measures.ProcessedMonthlyClosure)
	nextPMC := make(chan *process_measures.ProcessedMonthlyClosure)
	go call(ctxTimeout, previousMatch, optsPrevious, previousPMC, errPrevious)
	go call(ctxTimeout, nextMatch, optsNext, nextPMC, errNext)
	errorPrevious, errorNext := <-errPrevious, <-errNext
	if errorPrevious != nil && errorPrevious != mongo.ErrNoDocuments {
		return process_measures.ResumesProcessMonthlyClosure{}, errors.New("error in Mongo")
	}
	if errorNext != nil && errorNext != mongo.ErrNoDocuments {
		return process_measures.ResumesProcessMonthlyClosure{}, errors.New("error in Mongo")
	}
	MonthlyClosures.Next = <-nextPMC
	MonthlyClosures.Previous = <-previousPMC

	return MonthlyClosures, nil
}

func (repository ClosureRepositoryMongo) GetClosure(ctx context.Context, query process_measures.GetClosure) (process_measures.ProcessedMonthlyClosure, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	var match bson.D
	var opts *options.FindOneOptions
	var config process_measures.ProcessedMonthlyClosure

	if query.Id != "" {
		match = append(bson.D{
			{"distributor_id", query.DistributorId},
			{"_id", query.Id},
		})
	} else {
		match = append(bson.D{
			{"distributor_id", query.DistributorId},
			{"cups", query.CUPS},
			{"start_date",
				bson.M{
					"$gte": query.StartDate,
				},
			},
			{"end_date",
				bson.M{
					"$lte": query.EndDate,
				},
			},
		})
	}
	opts = nil

	result := repository.client.Database(repository.Database).Collection(repository.CollectionClosureConfig).
		FindOne(ctxTimeout, match, opts)

	err := result.Decode(&config)

	if query.Moment == process_measures.Before {

		match = append(bson.D{
			{"distributor_id", query.DistributorId},
			{"end_date",
				bson.M{
					"$lt": config.StartDate,
				},
			},
		})
		opts = options.FindOne().SetSort(bson.D{{"end_date", -1}})

	}
	if query.Moment == process_measures.Next {

		match = append(bson.D{
			{"distributor_id", query.DistributorId},
			{"start_date",
				bson.M{
					"$gt": config.EndDate,
				},
			},
		})
		opts = options.FindOne().SetSort(bson.D{{"start_date", 1}})
	}

	result = repository.client.Database(repository.Database).Collection(repository.CollectionClosureConfig).
		FindOne(ctxTimeout, match, opts)

	err = result.Decode(&config)

	if err != nil {
		return process_measures.ProcessedMonthlyClosure{}, err
	}

	return config, err

}

func (repository ClosureRepositoryMongo) CreateClosure(ctx context.Context, monthly process_measures.ProcessedMonthlyClosure) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.CollectionClosureConfig)

	upsert := true

	_, err := coll.UpdateOne(
		ctxTimeout,
		bson.D{{"_id", monthly.Id}},
		bson.D{{"$set", monthly}},
		&options.UpdateOptions{Upsert: &upsert},
	)

	if err != nil {
		return err
	}
	return err
}

func (repository ClosureRepositoryMongo) UpdateClosure(ctx context.Context, id string, monthly process_measures.ProcessedMonthlyClosure) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	coll := repository.client.Database(repository.Database).Collection(repository.CollectionClosureConfig)

	var model []mongo.WriteModel

	updateOneModel := mongo.NewUpdateOneModel()
	updateOneModel.SetFilter(bson.D{
		{"_id", id},
	})
	updateOneModel.SetUpdate(bson.D{
		{"$set", monthly},
	})

	model = append(model, updateOneModel)

	_, err := coll.BulkWrite(ctxTimeout, model)

	if err != nil {
		return err
	}
	return nil
}

func (repository ClosureRepositoryMongo) GetClosureOne(ctx context.Context, id string) (process_measures.ProcessedMonthlyClosure, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	var match bson.D
	if id != "" {
		match = append(bson.D{
			{"_id", id},
		})
	}

	var config process_measures.ProcessedMonthlyClosure

	result := repository.client.Database(repository.Database).Collection(repository.CollectionClosureConfig).FindOne(ctxTimeout, match)

	err := result.Decode(&config)

	if err != nil {
		return process_measures.ProcessedMonthlyClosure{}, err
	}

	return config, err

}
