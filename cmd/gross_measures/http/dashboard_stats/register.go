package dashboard_stats

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/services"
	mongo_r "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(s *server.Server, cnf config.Config, authMiddleware *middleware.Auth, mongoClient *mongo.Client) {

	repository := mongo_r.NewGrossMeasureRepositoryDashboardMongo(mongoClient, cnf.MeasureDB, "stats_gross_measure", "stats_gross_measure_by_serial_number")
	srv := services.NewGrossMeasuresStatsDashboardServices(repository)
	controller := Controller{
		dashboardService: srv,
	}

	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresGetDashboardGrossMeasure(
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
