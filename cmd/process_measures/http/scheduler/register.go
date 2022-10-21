package scheduler

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(s *server.Server, cnf config.Config, authMiddleware *middleware.Auth, db *gorm.DB) {

	repository := postgres.NewProcessMeasureSchedulerPostgres(db)

	srv := services.NewSchedulerServices(repository, scheduler.NewGcpScheduler(cnf.ProjectID, cnf.Location, cnf.TimeZone), cnf.ProcessedMeasureTopic)

	controller := Controller{srv: srv}

	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresCreateProcessMeasuresScheduler(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresDeleteProcessMeasuresScheduler(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresListProcessMeasuresScheduler(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresPatchProcessMeasuresScheduler(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresGetProcessMeasureSchedulerById(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
	})
}
