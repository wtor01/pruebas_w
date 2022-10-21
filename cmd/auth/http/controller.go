package http

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c Controller) AuthorizationUser(ctx *gin.Context) {

	user, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	response := Me{
		Distributors: make([]struct {
			Id   string `json:"id"`
			Name string `json:"name"`
			Role struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"role"`
		}, 0, cap(user.Distributors)),
		Email:   user.Email,
		Id:      user.ID,
		IsAdmin: user.IsAdmin,
		Name:    user.Name,
	}
	for _, p := range user.Distributors {
		response.Distributors = append(response.Distributors, struct {
			Id   string `json:"id"`
			Name string `json:"name"`
			Role struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"role"`
		}{
			Id:   p.ID,
			Name: p.Name,
			Role: struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			}{
				Id:   p.Role.ID,
				Name: p.Role.Name,
			},
		})
	}

	ctx.JSON(http.StatusOK, response)
}
