package features

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations/services"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(s *server.Server, cnf config.Config, authMiddleware *middleware.Auth, db *gorm.DB) {

	repository := postgres.NewAggregationsFeaturesPostgres(db)
	srv := services.NewFeaturesServices(repository)
	controller := Controller{srv: srv}
	s.Register(func(router *gin.RouterGroup) {
		s.Register(func(router *gin.RouterGroup) {
			RegisterHandlerWithMiddlewaresCreateAggregationFeaturesAvailable(
				router,
				controller,
				[]gin.HandlerFunc{
					authMiddleware.HttpSetOAuthUserMiddleware(),
					authMiddleware.HttpMustHaveAdminUserMiddleware(),
				})
		})
		s.Register(func(router *gin.RouterGroup) {
			RegisterHandlerWithMiddlewaresGetAggregationFeaturesAvailable(
				router,
				controller,
				[]gin.HandlerFunc{
					authMiddleware.HttpSetOAuthUserMiddleware(),
					authMiddleware.HttpMustHaveAdminUserMiddleware(),
				})
		})
		s.Register(func(router *gin.RouterGroup) {
			RegisterHandlerWithMiddlewaresDeleteAggregationFeaturesAvailable(
				router,
				controller,
				[]gin.HandlerFunc{
					authMiddleware.HttpSetOAuthUserMiddleware(),
					authMiddleware.HttpMustHaveAdminUserMiddleware(),
				})
		})
		s.Register(func(router *gin.RouterGroup) {
			RegisterHandlerWithMiddlewaresUpdateAggregationFeaturesAvailable(
				router,
				controller,
				[]gin.HandlerFunc{
					authMiddleware.HttpSetOAuthUserMiddleware(),
					authMiddleware.HttpMustHaveAdminUserMiddleware(),
				})
		})
		s.Register(func(router *gin.RouterGroup) {
			RegisterHandlerWithMiddlewaresGetAllAggregationFeaturesAvailable(
				router,
				controller,
				[]gin.HandlerFunc{
					authMiddleware.HttpSetOAuthUserMiddleware(),
					authMiddleware.HttpMustHaveAdminUserMiddleware(),
				})
		})

	})

}

