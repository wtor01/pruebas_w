package closures

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	tariff_services "bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff/services"
	mongo_r "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"time"
)

func Register(s *server.Server, authMiddleware *middleware.Auth, mongoClient *mongo.Client, cnf config.Config, db *gorm.DB, rd *redis.Client) {

	mongoRepo := mongo_r.NewClosureRepositoryMongo(mongoClient, cnf.MeasureDB)
	inventoryRepo := mongo_r.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)
	tariffRepository := postgres.NewTariffPostgres(db, redis_repos.NewDataCacheRedis(rd))
	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	masterTablesClient := master_tables.NewMasterTables(tariff_services.NewTariffServices(tariffRepository))
	executeServices := services.NewExecuteServices(
		inventoryRepo,
		publisher,
		cnf.ProcessedMeasureTopic,
	)
	srvs := services.NewClosureServices(mongoRepo, inventoryRepo, masterTablesClient, time.Now)
	resumeSrvs := services.NewGetResumeService(mongoRepo, cnf.LocalLocation)

	controller := NewController(srvs, executeServices, resumeSrvs)

	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresGetBillingClosures(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			})
		RegisterHandlerWithMiddlewaresCreateBillingClosures(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			})
		RegisterHandlerWithMiddlewaresUpdateBillingClosures(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			})
		RegisterHandlerWithMiddlewaresExecuteProcessMeasuresServices(

			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Admin),
			})
		RegisterHandlerWithMiddlewaresGetProcessMeasuresResumeByCups(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Admin),
			})
	})

}
