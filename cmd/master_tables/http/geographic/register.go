package geographic

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/geographic/services"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

// Register generate server
func Register(s *server.Server, cnf config.Config, authMiddleware *middleware.Auth, db *gorm.DB, rd *redis.Client) {
	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	geographicRepository := postgres.NewGeographicPostgres(db, redis_repos.NewDataCacheRedis(rd))

	gssrvs := services.NewGeographicServices(geographicRepository, publisher, cnf.TopicCalendarPeriods)
	controller := NewController(gssrvs)
	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresGetAllGeographicZones(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresDeleteGeographicZone(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresGetGeographicZone(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresInsertGeographicZone(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresModifyGeographicZone(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
	})

}
