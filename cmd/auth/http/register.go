package http

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
)

func Register(s *server.Server, authMiddleware *middleware.Auth) {

	controller := NewController()

	s.Register(func(router *gin.RouterGroup) {
		RegisterHandlerWithMiddlewaresAuthorizationUser(router, controller, []gin.HandlerFunc{authMiddleware.HttpSetOAuthUserMiddleware()})

	})
}
