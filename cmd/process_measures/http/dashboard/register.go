package dashboard

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	mongo_repos "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func Register(s *server.Server, authMiddleware *middleware.Auth, db *gorm.DB, cnf config.Config, mongoClient *mongo.Client, rd *redis.Client) {
	inventoryMeasuresRepository := mongo_repos.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)
	processRepository := mongo_repos.NewProcessMeasureRepository(
		mongoClient,
		cnf.MeasureDB,
		cnf.CollectionLoadCurve,
		cnf.CollectionDailyLoadCurve,
		cnf.CollectionMonthlyClosure,
		cnf.CollectionDailyClosure,
		cnf.LocalLocation,
	)

	svcs := services.NewDashboardServices(
		processRepository,
		inventoryMeasuresRepository,
		redis_repos.NewProcessMeasureFestiveDays(rd),
		cnf.LocalLocation,
	)

	controller := NewController(svcs)

	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresGetProcessMeasureDashboard(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor_id"), auth.Admin),
			})

		RegisterHandlerWithMiddlewaresGetCurveProcessServicePoint(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor"), auth.Admin),
			})

		RegisterHandlerWithMiddlewaresGetDashboardProcessServicePoint(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor"), auth.Admin),
			})
		RegisterHandlerWithMiddlewaresGetProcessMeasureDashboardList(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Admin),
			})
	})
}
