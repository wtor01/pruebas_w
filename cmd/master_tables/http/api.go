package main

import (
	"bitbucket.org/sercide/data-ingestion/cmd/master_tables/http/calendar"
	"bitbucket.org/sercide/data-ingestion/cmd/master_tables/http/festive_days"
	"bitbucket.org/sercide/data-ingestion/cmd/master_tables/http/generate_calendars"
	"bitbucket.org/sercide/data-ingestion/cmd/master_tables/http/geographic"
	"bitbucket.org/sercide/data-ingestion/cmd/master_tables/http/seasons"
	"bitbucket.org/sercide/data-ingestion/cmd/master_tables/http/tariff"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"log"
)

func main() {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, s := server.NewHttpServer(context.Background(), "", cnf.Port, cnf.ShutdownTimeout)
	db := postgres.New(cnf)

	authMiddleware := middleware.NewAuth(cnf)

	var redisClient *redis.Client

	if cnf.RedisEnabled {
		redisClient = redis_repos.New(cnf)
		defer redisClient.Close()
	}

	geographic.Register(&s, cnf, authMiddleware, db, redisClient)
	calendar.Register(&s, cnf, authMiddleware, db, redisClient)
	tariff.Register(&s, authMiddleware, db, redisClient)
	seasons.Register(&s, authMiddleware, db, redisClient)
	festive_days.Register(&s, cnf, authMiddleware, db, redisClient)
	generate_calendars.Register(&s, authMiddleware, db, redisClient)

	s.Register(func(router *gin.RouterGroup) {
		// OPENAPI
		server.ServeOpenapi(router, server.SwaggerUIOpts{
			SpecURL: "docs/master_tables/api.yaml",
		})
	})

	if err := s.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
