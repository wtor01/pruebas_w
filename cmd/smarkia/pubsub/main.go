package main

import (
	smarkia_pubsub "bitbucket.org/sercide/data-ingestion/cmd/smarkia/pubsub/register"
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

	if err != nil {
		log.Fatalln(err)
	}

	db := postgres.New(cnf)
	mongoClient, err := mongo.New(ctx, cnf)

	if err != nil {
		log.Fatalln(err)
	}
	var redisClient *redis.Client

	if cnf.RedisEnabled {
		redisClient = redisRepos.New(cnf)
		defer redisClient.Close()
	}

	if err := smarkia_pubsub.Register(ctx, db, cnf, redisClient, mongoClient); err != nil {
		log.Fatal(err)
	}
}
