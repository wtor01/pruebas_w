package festive_days

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/festive_days/services"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"

	//"bitbucket.org/sercide/data-ingestion/internal/master_tables"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"gorm.io/gorm"
)

func Register(s *server.Server, cnf config.Config, authMiddleware *middleware.Auth, db *gorm.DB, rd *redis.Client) {
	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	repository := postgres.NewFestiveDaysPostgres(db, redis_repos.NewDataCacheRedis(rd))
	service := services.New(repository, publisher, cnf.TopicCalendarPeriods)
	controller := NewController(service)

	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresListFestiveDays(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			})
		RegisterHandlerWithMiddlewaresPostFestiveDay(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			})
		RegisterHandlerWithMiddlewaresGetFestiveDay(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			})
		RegisterHandlerWithMiddlewaresPutFestiveDay(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			})
		RegisterHandlerWithMiddlewaresDeleteFestiveDay(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			})
	})

}
