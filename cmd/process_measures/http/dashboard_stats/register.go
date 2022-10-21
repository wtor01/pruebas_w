package stats

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	mongo_r "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(s *server.Server, cnf config.Config, authMiddleware *middleware.Auth, mongoClient *mongo.Client) {
	repository := mongo_r.NewProcessMeasuresDashboardStatsRepositoryMongo(mongoClient, cnf.MeasureDB, "stats_processed_measure", "stats_processed_measure_by_cups")
	srv := services.NewProcessMeasuresStatsDashboardServices(repository)
	controller := Controller{
		srv: srv,
	}
	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresGetProcessMeasureStatistics(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor_id"), auth.Admin),
			})
		RegisterHandlerWithMiddlewaresGetProcessMeasureStatisticsByCups(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor_id"), auth.Admin),
			})
	})

}
