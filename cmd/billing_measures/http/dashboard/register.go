package dashboard

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/services"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	mongo_r "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/internal/platform/pubsub"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(s *server.Server, cnf config.Config, authMiddleware *middleware.Auth, mongoClient *mongo.Client) {
	publisher := pubsub.PublisherCreatorGcp(cnf.ProjectID)
	repository := mongo_r.NewBillingMeasuresDashboardRepositoryMongo(mongoClient, cnf.MeasureDB, "billing_measure")

	inventoryRepository := mongo_r.NewInventoryRepositoryMongo(mongoClient, cnf.MeasureDB)
	srv := services.NewBillingMeasuresDashboardService(repository, cnf.LocalLocation, inventoryRepository, publisher, cnf.TopicBillingMeasures)
	controller := Controller{
		srv: srv,
		cnf: cnf,
	}

	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresSearchDistributorFiscalBillingMeasures(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor_id"), auth.Admin),
			})
		RegisterHandlerWithMiddlewaresCreateBillingMeasuresMVH(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor_id"), auth.Admin),
			})
		RegisterHandlerWithMiddlewaresGetBillingMeasuresResumeById(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor_id"), auth.Admin),
			})
		RegisterHandlerWithMiddlewaresGetBillingMeasureDashboardSummary(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor_id"), auth.Admin),
			})

		RegisterHandlerWithMiddlewaresGetBillingMeasureTaxMeasurebycups(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				middleware.HttpOAuthUserPermissionMiddleware(middleware.HttpExtractDistributorFromQuery("distributor_id"), auth.Admin),
			})

	})
}
