package main

import (
	measures_pubsub "bitbucket.org/sercide/data-ingestion/cmd/gross_measures/pubsub/register"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redisRepos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"context"
	"github.com/go-redis/redis/v9"
	"log"
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

	if err := measures_pubsub.Register(ctx, db, cnf, mongoClient, redisClient); err != nil {
		log.Fatal(err)
	}
}
