package main

import (
	aggregations_pubsub "bitbucket.org/sercide/data-ingestion/cmd/aggregations/pubsub/register"
	billing_measures_pubsub "bitbucket.org/sercide/data-ingestion/cmd/billing_measures/pubsub/register"
	gross_measures_pubsub "bitbucket.org/sercide/data-ingestion/cmd/gross_measures/pubsub/register"
	calendar_redis_pubsub "bitbucket.org/sercide/data-ingestion/cmd/master_tables/pubsub/register"
	process_measures_pubsub "bitbucket.org/sercide/data-ingestion/cmd/process_measures/pubsub/register"
	re_process_measure_pubsub "bitbucket.org/sercide/data-ingestion/cmd/re_process_measures/pubsub/register"
	smarkia_pubsub "bitbucket.org/sercide/data-ingestion/cmd/smarkia/pubsub/register"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redisRepos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"context"
	"github.com/go-redis/redis/v9"
	"log"
	"os"
)

func main() {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}
	ctx := context.Background()

	mongoClient, err := mongo.New(ctx, cnf)

	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := postgres.New(cnf)

	var redisClient *redis.Client

	if cnf.RedisEnabled {
		redisClient = redisRepos.New(cnf)
		defer redisClient.Close()
	}

	switch os.Getenv("PUBSUB_PROCESS") {
	case "GROSS_MEASURES":
		{
			if err := gross_measures_pubsub.Register(ctx, db, cnf, mongoClient, redisClient); err != nil {
				log.Fatal(err)
			}
		}
	case "PROCESS_MEASURES":
		{
			if err := process_measures_pubsub.Register(ctx, cnf, mongoClient, db, redisClient); err != nil {
				log.Fatal(err)
			}
		}
	case "RE_PROCESS_MEASURES":
		{
			if err := re_process_measure_pubsub.Register(ctx, cnf, mongoClient, db, redisClient); err != nil {
				log.Fatal(err)
			}
		}
	case "BILLING_MEASURES":
		{
			if err := billing_measures_pubsub.Register(ctx, cnf, mongoClient, db, redisClient); err != nil {
				log.Fatal(err)
			}
		}
	case "SMARKIA":
		{
			if err := smarkia_pubsub.Register(ctx, db, cnf, redisClient, mongoClient); err != nil {
				log.Fatal(err)
			}
		}
	case "AGGREGATIONS":
		{
			if err := aggregations_pubsub.Register(ctx, cnf, mongoClient, db, redisClient); err != nil {
				log.Fatal(err)
			}
		}
	case "CALENDARS_REDIS":
		{
			if err := calendar_redis_pubsub.Register(ctx, cnf, db, redisClient); err != nil {
				log.Fatal(err)
			}
		}
	}
}
