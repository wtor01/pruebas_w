package config

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/internal/aggregations/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	aggregationConfigService *services.AggregationConfigServices
}

func (c Controller) GetAllAggregationsConfig(ctx *gin.Context, params GetAllAggregationsConfigParams) {
	var query string
	if params.Q != nil {
		query = *params.Q
	}

	result, count, err := c.aggregationConfigService.GetAggregationConfigsService.Handler(ctx.Request.Context(),
		services.NewGetAggregationConfigsServiceDto(
			query,
			params.Limit,
			params.Offset,
		))

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	linksList := server.GetLinksList(ctx, count, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
		Q:      params.Q,
		Host:   ctx.Request.Host,
	})

	ctx.JSON(http.StatusOK, GetAggregationsConfig{
		Pagination: Pagination{
			Links:  linksList,
			Count:  count,
			Limit:  params.Limit,
			Offset: params.Offset,
			Size:   len(result),
		},
		Results: utils.MapSlice(result, AggregationConfigToResponse),
	})

}

func (c Controller) CreateAggregationConfig(ctx *gin.Context) {
	var body CreateAggregationConfigJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	aggregationConfigDto := services.NewCreateConfigDto(body.Name, body.Scheduler, body.Description, body.StartDate, body.EndDate, utils.MapSlice(body.Features, func(item AggregationFeature) aggregations.ConfigFeatureDto {
		return aggregations.NewConfigFeatureDto(item.Id, item.Name, item.Field)
	}))

	result, err := c.aggregationConfigService.CreateAggregationConfigService.Handler(ctx.Request.Context(), aggregationConfigDto)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusCreated, AggregationConfigToResponse(result))
}

func (c Controller) DeleteAggregationConfig(ctx *gin.Context, aggregationConfigId string) {
	err := c.aggregationConfigService.DeleteAggregationConfigService.Handler(ctx.Request.Context(), aggregationConfigId)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, aggregationConfigId)
}

func (c Controller) GetAggregationConfig(ctx *gin.Context, aggregationConfigId string) {
	result, err := c.aggregationConfigService.GetAggregationConfigByIdService.Handler(ctx.Request.Context(), aggregationConfigId)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, AggregationConfigToResponse(result))
}

func (c Controller) UpdateAggregationConfig(ctx *gin.Context, aggregationConfigId string) {
	var body UpdateAggregationConfigJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	aggregationConfigDto := services.NewUpdateConfigDto(body.Scheduler, body.Description, body.StartDate, body.EndDate, utils.MapSlice(body.Features, func(item AggregationFeature) aggregations.ConfigFeatureDto {
		return aggregations.NewConfigFeatureDto(item.Id, item.Name, item.Field)
	}))

	result, err := c.aggregationConfigService.UpdateAggregationConfigService.Handler(ctx.Request.Context(), aggregationConfigId, aggregationConfigDto)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, AggregationConfigToResponse(result))
}

func NewController(s *services.AggregationConfigServices) *Controller {
	return &Controller{
		aggregationConfigService: s,
	}
}
