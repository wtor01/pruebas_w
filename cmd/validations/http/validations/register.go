package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	mongo_repos "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	redis_repos "bitbucket.org/sercide/data-ingestion/internal/platform/redis"
	"bitbucket.org/sercide/data-ingestion/internal/validations/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func Register(s *server.Server, authMiddleware *middleware.Auth, db *gorm.DB, mongoClient *mongo.Client, rd *redis.Client, cnf config.Config) {
	repository := postgres.NewValidationMeasurePostgresRepository(db, redis_repos.NewDataCacheRedis(rd), cnf.ShutdownTimeout)
	repositoryMeasure := mongo_repos.NewProcessMeasureRepository(
		mongoClient,
		cnf.MeasureDB,
		cnf.CollectionLoadCurve,
		cnf.CollectionDailyLoadCurve,
		cnf.CollectionMonthlyClosure,
		cnf.CollectionDailyClosure,
		cnf.LocalLocation,
	)

	adminServices := services.NewServices(repository, repositoryMeasure)

	controller := NewController(adminServices)

	s.Register(func(router *gin.RouterGroup) {
		// VALIDATION MEASURE
		RegisterHandlerWithMiddlewaresCreateValidationsMeasure(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			},
		)
		RegisterHandlerWithMiddlewaresUpdateValidationsMeasure(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			},
		)
		RegisterHandlerWithMiddlewaresDeleteValidationsMeasure(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			},
		)
		RegisterHandlerWithMiddlewaresListValidationsMeasure(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
			},
		)
		RegisterHandlerWithMiddlewaresGetValidationsMeasure(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
			},
		)

		// VALIDATION MEASURE CONFIG
		RegisterHandlerWithMiddlewaresCreateValidationsMeasureConfig(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresUpdateValidationsMeasureConfig(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator)},
		)
		RegisterHandlerWithMiddlewaresDeleteValidationsMeasureConfig(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator)},
		)
		RegisterHandlerWithMiddlewaresListValidationsMeasureConfig(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresGetValidationsMeasureConfig(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Operator),
			},
		)
		RegisterHandlerWithMiddlewaresPutMeasurementValidation(router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor_id"), auth.Admin),
			})
	})
}
