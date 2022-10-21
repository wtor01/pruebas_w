package aggregations

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations"
	"bitbucket.org/sercide/data-ingestion/internal/aggregations/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	srv *services.AggregationsService
}

func (c Controller) GetAllAggregations(ctx *gin.Context, params GetAllAggregationsParams) {
	dto := aggregations.GetAggregationsDto{
		Offset:              params.Offset,
		Limit:               params.Limit,
		AggregationConfigId: params.AggregationConfigId,
		StartDate:           utils.ToUTC(params.StartDate.Time, c.srv.Location),
		EndDate:             utils.ToUTC(params.EndDate.Time, c.srv.Location),
	}

	result, count, err := c.srv.GetAggregations.Handler(ctx.Request.Context(), dto)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results := utils.MapSlice(result, transformToResponseAggregations)

	paginate := server.GetPaginate(ctx, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
	}, result, int(count))

	ctx.JSON(http.StatusOK, Aggregations{
		Pagination: Pagination{
			Links:  paginate.Links,
			Count:  paginate.Count,
			Limit:  paginate.Limit,
			Offset: paginate.Offset,
			Size:   paginate.Size,
		},
		AggregationsResults: AggregationsResults{Results: &results},
	})

}

func (c Controller) GetAggregation(ctx *gin.Context, aggregationId string) {
	dto := aggregations.GetAggregationDto{
		AggregationConfigId: aggregationId,
	}

	result, err := c.srv.GetAggregation.Handler(ctx, dto)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, transformToResponseAggregation(result))

}
