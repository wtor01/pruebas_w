package scheduler

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/services"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/pkg/scheduler"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(s *server.Server, cnf config.Config, authMiddleware *middleware.Auth, db *gorm.DB) {
	repository := postgres.NewProcessBillingSchedulerPostgres(db)
	srv := services.NewSchedulerServices(repository, scheduler.NewGcpScheduler(cnf.ProjectID, cnf.Location, cnf.TimeZone), cnf.TopicBillingMeasures)
	controller := Controller{srv: srv}
	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresListBillingMeasuresScheduler(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresGetBillingMeasuresSchedulerById(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresDeleteBillingMeasuresScheduler(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresCreateBillingMeasuresScheduler(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})

		RegisterHandlerWithMiddlewaresPatchBillingMeasuresScheduler(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
	})
}
