package generate_calendars

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/services"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

func Register(s *server.Server, authMiddleware *middleware.Auth, db *gorm.DB, rd *redis.Client) {
	calendarRepository := postgres.NewCalendarPostgres(db, redis_repos.NewDataCacheRedis(rd))
	festiveDaysRepository := postgres.NewFestiveDaysPostgres(db, redis_repos.NewDataCacheRedis(rd))
	calendarPeriodRepository := redis_repos.NewProcessMeasureFestiveDays(rd)

	srv := services.NewCalendarPeriodsPubSubServices(calendarRepository, calendarPeriodRepository, festiveDaysRepository)

	controller := NewController(srv)

	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresGenerateCalendarsPeriods(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			},
		)
		RegisterHandlerWithMiddlewaresGenerateFestiveDays(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			},
		)
	})

}
