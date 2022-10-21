package middleware

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/auth/services"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/internal/platform/firebase"
	"bitbucket.org/sercide/data-ingestion/internal/platform/postgres"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
)

const UserKey = "user"

const (
	AuthorizationHeader string = "Authorization"
)

type Auth struct {
	meService *services.MeService
}

func NewAuth(cnf config.Config) *Auth {
	firebaseClient, err := firebase.NewOAuthClient()
	if err != nil {
		panic(err)
	}
	db := postgres.New(cnf)
	authRepository := postgres.NewAuthRepository(db)

	meService := services.NewMeService(authRepository, firebaseClient)

	return &Auth{meService: meService}
}

var regexToken = regexp.MustCompile(`(?i)BEARER`)

func GetAuthUser(ctx *gin.Context) (auth.User, error) {
	userAuth, ok := ctx.Get(UserKey)
	if !ok {
		return auth.User{}, errors.New("")
	}

	u, ok := userAuth.(auth.User)

	if !ok {
		return auth.User{}, errors.New("")
	}

	return u, nil
}

func tokenFromHeader(authorizationHeader string) string {
	authorizationHeader = regexToken.ReplaceAllString(authorizationHeader, "")
	strings.Trim(authorizationHeader, " ")

	return strings.Trim(authorizationHeader, " ")
}

func (a Auth) HttpSetOAuthUserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token := tokenFromHeader(ctx.Request.Header.Get(AuthorizationHeader))

		u, err := a.meService.Handler(ctx, services.MeDto{
			Token: token,
		})

		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		ctx.Set(UserKey, u)
		ctx.Next()
	}
}

func (a Auth) HttpMustHaveUserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		_, err := GetAuthUser(ctx)

		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		ctx.Next()
	}
}

func (a Auth) HttpMustHaveAdminUserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		u, err := GetAuthUser(ctx)

		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if !u.IsAdmin {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}

type HttpExtractorDistributor = func(ctx *gin.Context) string

func HttpExtractDistributorFromQuery(distributorKeyQuery string) HttpExtractorDistributor {
	return func(ctx *gin.Context) string {
		return ctx.Query(distributorKeyQuery)
	}
}

func HttpExtractDistributorFromParam(distributorKeyQuery string) HttpExtractorDistributor {
	return func(ctx *gin.Context) string {
		return ctx.Param(distributorKeyQuery)
	}
}

func HttpOAuthUserPermissionMiddleware(distributorExtractor HttpExtractorDistributor, roleName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		user, err := GetAuthUser(ctx)

		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		distributorID := distributorExtractor(ctx)

		if !user.IsAdmin && !user.HasPermissionInDistributor(distributorID, roleName) {
			_ = ctx.AbortWithError(http.StatusForbidden, errors.New("forbidden"))
			return
		}

		ctx.Next()
	}
}
