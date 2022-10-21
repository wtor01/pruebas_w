package aggregations

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations/services"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	mongo_r "bitbucket.org/sercide/data-ingestion/internal/platform/mongo"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(s *server.Server, cnf config.Config, authMiddleware *middleware.Auth, mongoClient *mongo.Client) {
	repository := mongo_r.NewAggregationsRepositoryMongo(mongoClient, cnf.MeasureDB, "aggregations")
	srv := services.NewAggregationsService(repository, cnf.LocalLocation)

	controller := Controller{srv: srv}

	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresGetAllAggregations(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
	})
	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresGetAggregation(
			router,
			controller,
			[]gin.HandlerFunc{
				authMiddleware.HttpSetOAuthUserMiddleware(),
				authMiddleware.HttpMustHaveAdminUserMiddleware(),
			})
	})

}
