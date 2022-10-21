package main

import (
	"bitbucket.org/sercide/data-ingestion/cmd/aggregations/http/aggregations"
	aggregation_config "bitbucket.org/sercide/data-ingestion/cmd/aggregations/http/config"
	"bitbucket.org/sercide/data-ingestion/cmd/aggregations/http/features"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
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

	authMiddleware := middleware.NewAuth(cnf)

	var redisClient *redis.Client

	if cnf.RedisEnabled {
		redisClient = redis_repos.New(cnf)
		defer redisClient.Close()
	}

	db := postgres.New(cnf)

	mongoClient, err := mongo.New(ctx, cnf)

	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	tp := telemetry.NewTracerProvider("aggregations")

	defer tp.ForceFlush(ctx)

	aggregation_config.Register(&s, authMiddleware, db, cnf)
	aggregations.Register(&s, cnf, authMiddleware, mongoClient)

	features.Register(&s, cnf, authMiddleware, db)

	s.Register(func(router *gin.RouterGroup) {
		server.ServeOpenapi(router, server.SwaggerUIOpts{
			SpecURL: "docs/aggregations/api.yaml",
		})
	})

	if err := s.Run(ctx); err != nil {
		log.Fatal(err)
	}

}
