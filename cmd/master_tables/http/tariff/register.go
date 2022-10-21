package tariff

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff/services"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

// Register generate server
func Register(s *server.Server, authMiddleware *middleware.Auth, db *gorm.DB, rd *redis.Client) {
	geographicRepository := postgres.NewTariffPostgres(db, redis_repos.NewDataCacheRedis(rd))

	tfsrv := services.NewTariffServices(geographicRepository)
	controller := NewController(tfsrv)
	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresGetAllTariffs(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresDeleteTariff(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresGetTariffs(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresInsertTariffs(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresModifyTariffs(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresGetTariffsCalendar(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
	})

}
