package self_consumption

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/services"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	mongo_r "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(s *server.Server, authMiddleware *middleware.Auth, mongoClient *mongo.Client, cnf config.Config) {

	repository := mongo_r.NewSelfConsumptionRepository(mongoClient, cnf.MeasureDB, cnf.LocalLocation)
	billingSelfConsumptionRepository := mongo_r.NewBillingSelfConsumptionRepository(mongoClient, cnf.MeasureDB, cnf.LocalLocation)
	srv := services.NewSelfConsumptionServices(repository, billingSelfConsumptionRepository, cnf.LocalLocation)
	controller := Controller{srv: srv}
	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresSearchActivesSelfConsumptionUnitConfigByDistributorId(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresSearchSelfConsumptionUnitConfig(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
		RegisterHandlerWithMiddlewaresGetSelfConsumptionByCau(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromParam("distributor_id"), auth.Admin),
			})
	})

}
