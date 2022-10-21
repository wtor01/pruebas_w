package supply_point

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/clients/internal_clients/master_tables"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/services"
	tariff_services "bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff/services"
	mongo_repos "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func Register(s *server.Server, authMiddleware *middleware.Auth, db *gorm.DB, cnf config.Config, mongoClient *mongo.Client, rd *redis.Client) {
	measureRepository := mongo_repos.NewGrossMeasureRepositoryMongo(mongoClient, cnf.MeasureDB, cnf.MeasureCollection, cnf.LocalLocation)
	inventoryMeasuresRepository := mongo_repos.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)
	calendarRepository := redis_repos.NewProcessMeasureFestiveDays(rd)
	tariffRepository := postgres.NewTariffPostgres(db, redis_repos.NewDataCacheRedis(rd))
	masterTablesClient := master_tables.NewMasterTables(tariff_services.NewTariffServices(tariffRepository))

	svcs := services.NewDashboardServices(
		measureRepository,
		inventoryMeasuresRepository,
		calendarRepository,
		masterTablesClient,
		cnf.LocalLocation,
	)

	controller := NewController(svcs, cnf.LocalLocation)

	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresGetCurveGrossMeasureMeter(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor_id"), auth.Admin),
			})
		RegisterHandlerWithMiddlewaresGetGrossMeasureServicePoint(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor_id"), auth.Admin),
			})
	})
}
