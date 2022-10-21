package festive_days

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/festive_days"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/festive_days/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Controller struct {
	srv *services.FestiveDaysServices
}

func NewController(srv *services.FestiveDaysServices) *Controller {
	return &Controller{srv: srv}
}

func (c Controller) ListFestiveDays(ctx *gin.Context, params ListFestiveDaysParams) {

	sortParams := make([]string, 0)

	if params.Sort != nil {
		sortParams = *params.Sort
	}

	query := ""
	if params.Q != nil {
		query = *params.Q
	}

	limit := 10
	if params.Limit != nil {
		limit = *params.Limit
	}

	result, count, err := c.srv.GetListFestiveDays.Handler(ctx, services.GetListFestiveDaysDto{
		Q:      query,
		Limit:  limit,
		Offset: params.Offset,
		Sort:   festive_days.NewSort(sortParams),
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := make([]FestiveDays, 0, cap(result))
	//layout := "02-01-2006"
	for _, d := range result {
		res = append(res, FestiveDaysToResponse(d))
	}

	links := server.GetLinksList(ctx, count, server.Query{
		Limit:  limit,
		Offset: params.Offset,
		Q:      params.Q,
		Sort:   params.Sort,
	})

	ctx.JSON(http.StatusOK, FestiveDaysList{
		Links: struct {
			Next *string `json:"next,omitempty"`
			Prev *string `json:"prev,omitempty"`
			Self string  `json:"self"`
		}{
			Next: links.Next,
			Prev: links.Prev,
			Self: links.Self,
		},
		Limit:   limit,
		Offset:  params.Offset,
		Results: res,
		Count:   count,
		Size:    len(result),
	})
}

func (c Controller) PostFestiveDay(ctx *gin.Context) {

	var req PutFestiveDayJSONRequestBody

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetAuthUser(ctx)

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	description := ""
	if req.Description != "" {
		description = req.Description
	}
	layout := "02-01-2006"
	sd, _ := time.Parse(layout, req.Date)
	err = c.srv.SaveFestiveDay.Handler(ctx, festive_days.FestiveDay{
		Date:         sd,
		Description:  description,
		GeographicId: req.GeographicId,
		CreatedBy:    user.ID,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, "created")
}

func (c Controller) DeleteFestiveDay(ctx *gin.Context, festiveDaysId string) {
	//generated dto and send to service
	err := c.srv.DeleteFestiveDay.Handler(ctx, festiveDaysId)
	//check errors and send
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": "deleted"})
}

func (c Controller) GetFestiveDay(ctx *gin.Context, festiveDaysId string) {

	d, err := c.srv.GetFestiveDayById.Handler(ctx, services.GetFestiveDayByIdDto{
		Id: festiveDaysId,
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, FestiveDaysToResponse(d))
}

func (c Controller) PutFestiveDay(ctx *gin.Context, festiveDaysId string) {

	var req PutFestiveDayJSONRequestBody

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.NewUUID()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	description := ""
	if req.Description != "" {
		description = req.Description
	}
	layout := "02-01-2006"
	sd, _ := time.Parse(layout, req.Date)
	err = c.srv.UpdateFestiveDay.Handler(ctx, festiveDaysId, services.UpdateFestiveDayDto{
		Id:           id.String(),
		Date:         sd,
		Description:  description,
		GeographicId: req.GeographicId,
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, "modified")
}
