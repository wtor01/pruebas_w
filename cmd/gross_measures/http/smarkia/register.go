package smarkia

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	services "bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia/services"
	inventory_services "bitbucket.org/sercide/data-ingestion/internal/inventory/services"
	mongo_repos "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/internal/platform/smarkia_api"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func Register(s *server.Server, authMiddleware *middleware.Auth, db *gorm.DB, cnf config.Config, mongoClient *mongo.Client, rd *redis.Client) {

	smarkiaApi := smarkia_api.NewApi(cnf.SmarkiaToken, cnf.SmarkiaHost, cnf.LocalLocation)
	inventoryRepository := postgres.NewInventoryPostgres(db, redis_repos.NewDataCacheRedis(rd))
	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	inventoryRepositoryMongo := mongo_repos.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)

	recoverSmarkiaMeasuresService := services.NewRecoverSmarkiaMeasures(
		smarkiaApi,
		inventory.NewInventory(inventory_services.New(inventoryRepository, inventoryRepositoryMongo), redis_repos.NewDataCacheRedis(rd)),
		cnf.LocalLocation,
		publisher,
		cnf.TopicMeasures,
		inventoryRepositoryMongo,
	)

	controller := NewController(recoverSmarkiaMeasuresService)

	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresRecoverSmarkiaMeasures(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
	})
}
