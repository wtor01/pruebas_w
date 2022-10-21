package main

import (
	"bitbucket.org/sercide/data-ingestion/cmd/auth/http"
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"context"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	logger := log.Default()

	cnf, err := config.LoadConfig(logger)

	if err != nil {
		logger.Fatal(err)
	}

	ctx, s := server.NewHttpServer(context.Background(), "", cnf.Port, cnf.ShutdownTimeout)

	authMiddleware := middleware.NewAuth(cnf)

	http.Register(&s, authMiddleware)

	s.Register(func(router *gin.RouterGroup) {

		server.ServeOpenapi(router, server.SwaggerUIOpts{
			SpecURL: "docs/auth.yaml",
		})
	})

	if err := s.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
