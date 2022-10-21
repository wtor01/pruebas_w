package seasons

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons/services"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

func Register(s *server.Server, authMiddleware *middleware.Auth, db *gorm.DB, rd *redis.Client) {
	seasonsRepository := postgres.NewSeasonPostgres(db, redis_repos.NewDataCacheRedis(rd))

	gssrvs := services.NewSeasonsServices(seasonsRepository)
	controller := NewController(gssrvs)
	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresGetAllSeasons(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresGetSeason(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresInsertSeason(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresModifySeason(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresDeleteSeason(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresGetAllDayTypes(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresGetDayType(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresInsertDayType(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresModifyDayType(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresDeleteDayType(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
	})

}
