package main

import (
	"bitbucket.org/sercide/data-ingestion/cmd/inventory/http"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
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
	mongoClient, err := mongo.New(ctx, cnf)
	if err != nil {
		log.Fatalln(err)
	}
	authMiddleware := middleware.NewAuth(cnf)

	var redisClient *redis.Client

	if cnf.RedisEnabled {
		redisClient = redis_repos.New(cnf)
		defer redisClient.Close()
	}
	http.Register(&s, authMiddleware, db, redisClient, mongoClient, cnf)

	s.Register(func(router *gin.RouterGroup) {
		// OPENAPI
		server.ServeOpenapi(router, server.SwaggerUIOpts{
			SpecURL: "docs/inventory.yaml",
		})
	})

	if err := s.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
