package config

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations/services"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(s *server.Server, authMiddleware *middleware.Auth, db *gorm.DB, cnf config.Config) {

	aggregationRepository := postgres.NewAggregationsConfigRepository(db)
	featuresRepository := postgres.NewAggregationsFeaturesPostgres(db)
	schedulerCreater := scheduler.NewGcpScheduler(cnf.ProjectID, cnf.Location, cnf.TimeZone)
	svcs := services.NewAggregationConfigServices(aggregationRepository, featuresRepository, schedulerCreater, cnf.TopicAggregations, cnf.LocalLocation)

	controller := NewController(svcs)

	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresGetAllAggregationsConfig(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresGetAggregationConfig(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresCreateAggregationConfig(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresUpdateAggregationConfig(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresDeleteAggregationConfig(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
	})
}
