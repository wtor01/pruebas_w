package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"log"
)

func main() {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	db := postgres.New(cnf)

	var redisClient *redis.Client

	if cnf.RedisEnabled {
		redisClient = redis_repos.New(cnf)
		defer redisClient.Close()
	}

	inventoryRepository := postgres.NewInventoryPostgres(db, redis_repos.NewDataCacheRedis(redisClient))

	result, count, err := inventoryRepository.ListDistributors(context.Background(), inventory.Search{
		Values: map[string]string{
			"is_smarkia_active": "1",
		},
		Limit: 200,
	})

	fmt.Print(result, count, err)
}
