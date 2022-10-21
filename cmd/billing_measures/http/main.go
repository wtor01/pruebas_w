package main

import (
	"bitbucket.org/sercide/data-ingestion/cmd/billing_measures/http/dashboard"
	"bitbucket.org/sercide/data-ingestion/cmd/billing_measures/http/scheduler"
	"bitbucket.org/sercide/data-ingestion/cmd/billing_measures/http/self_consumption"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redisRepos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
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

	tp := telemetry.NewTracerProvider("billing_measure-api")

	defer tp.ForceFlush(ctx)

	var redisClient *redis.Client

	if cnf.RedisEnabled {
		redisClient = redisRepos.New(cnf)
		defer redisClient.Close()
	}

	scheduler.Register(&s, cnf, authMiddleware, db)
	dashboard.Register(&s, cnf, authMiddleware, mongoClient)
	self_consumption.Register(&s, authMiddleware, mongoClient, cnf)

	s.Register(func(router *gin.RouterGroup) {
		server.ServeOpenapi(router, server.SwaggerUIOpts{
			SpecURL: "docs/billing_measures/api.yaml",
		})
	})

	if err := s.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
