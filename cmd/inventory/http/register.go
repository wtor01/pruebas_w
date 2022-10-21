package http

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/inventory/services"
	mongo_r "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func Register(s *server.Server, authMiddleware *middleware.Auth, db *gorm.DB, rd *redis.Client, mongoClient *mongo.Client, cnf config.Config) {
	inventoryRepository := postgres.NewInventoryPostgres(db, redis_repos.NewDataCacheRedis(rd))
	mongoRepo := mongo_r.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)
	srvs := services.New(inventoryRepository, mongoRepo)

	controller := NewController(srvs)

	s.Register(func(router *gin.RouterGroup) {
		// EQUIPMENTS ENDPOINTS
		RegisterHandlerWithMiddlewaresListMeasureEquipment(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			})
		RegisterHandlerWithMiddlewaresGetMeasureEquipment(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			})
		RegisterHandlerWithMiddlewaresGetSmarkiaMeasureEquipment(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresGetMeterConfig(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})

		// DISTRIBUTORS ENDPOINTS
		RegisterHandlerWithMiddlewaresListDistributors(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresGetSmarkiaDistributor(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresGetDistributor(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("id"), auth.Operator),
			})
	})
}
