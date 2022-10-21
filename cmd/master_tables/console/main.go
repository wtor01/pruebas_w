package main

import (
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/services"
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

	ctx := context.Background()

	db := postgres.New(cnf)
	var redisClient *redis.Client

	if cnf.RedisEnabled {
		redisClient = redis_repos.New(cnf)
		defer redisClient.Close()
	}

	calendarRepository := postgres.NewCalendarPostgres(db, redis_repos.NewDataCacheRedis(redisClient))
	calendarPeriodRepository := redis_repos.NewProcessMeasureFestiveDays(redisClient)
	festiveDaysRepository := postgres.NewFestiveDaysPostgres(db, redis_repos.NewDataCacheRedis(redisClient))

	sr := services.NewCalendarPeriodsPubSubServices(calendarRepository, calendarPeriodRepository, festiveDaysRepository)

	err = sr.CalendarPeriodGenerateService.Handler(ctx)
	if err != nil {
		fmt.Println(err)
	}

	err = sr.FestiveDaysGenerateService.Handler(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
