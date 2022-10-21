package http

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/inventory/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c Controller) ListDistributors(ctx *gin.Context, params ListDistributorsParams) {

	sortParams := make([]string, 0)

	user, _ := middleware.GetAuthUser(ctx)

	if params.Sort != nil {
		sortParams = *params.Sort
	}

	query := ""
	if params.Q != nil {
		query = *params.Q
	}

	ds, count, err := c.srv.GetDistributorsService.Handler(ctx, services.GetDistributorsDto{
		Q:           query,
		Limit:       params.Limit,
		Offset:      params.Offset,
		Sort:        inventory.NewSort(sortParams),
		CurrentUser: user,
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res := make([]Distributor, 0, cap(ds))

	for _, d := range ds {
		res = append(res, distributorToResponse(d))
	}

	links := server.GetLinksList(ctx, count, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
		Q:      params.Q,
		Sort:   params.Sort,
	})

	ctx.JSON(http.StatusOK, DistributorsResponseList{
		Links: struct {
			Next *string `json:"next,omitempty"`
			Prev *string `json:"prev,omitempty"`
			Self string  `json:"self"`
		}{
			Next: links.Next,
			Prev: links.Prev,
			Self: links.Self,
		},
		Limit:   params.Limit,
		Offset:  params.Offset,
		Results: res,
		Count:   count,
		Size:    len(ds),
	})
}

func (c Controller) GetSmarkiaDistributor(ctx *gin.Context, smarkiaId string) {
	d, err := c.srv.GetDistributorBySmarkiaIdService.Handler(ctx, services.GetDistributorBySmarkiaDto{
		Id: smarkiaId,
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, distributorToResponse(d))
}

func (c Controller) GetDistributor(ctx *gin.Context, id string) {
	d, err := c.srv.GetDistributorByIdService.Handler(ctx, services.GetDistributorBySmarkiaDto{
		Id: id,
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, distributorToResponse(d))
}
