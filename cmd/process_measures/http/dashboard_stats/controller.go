package stats

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller struct {
	srv *services.ProcessMeasuresStatsDashboardServices
}

func (c Controller) GetProcessMeasureStatisticsByCups(ctx *gin.Context, params GetProcessMeasureStatisticsByCupsParams) {
	var offset int
	if params.Offset != nil {
		offset = *params.Offset
	}

	result, count, err := c.srv.GetDashboardMeasuresStatsCups.Handler(ctx.Request.Context(), services.GetDashboardMeasuresStatsCupsDTO{
		DistributorID: params.DistributorId,
		Month:         params.Month,
		Year:          params.Year,
		Type:          string(params.Type),
		Offset:        offset,
		Limit:         params.Limit,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paginate := server.GetPaginate(ctx, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
		Parameters: map[string]string{
			"distributor_id": params.DistributorId,
			"type":           string(params.Type),
			"year":           strconv.Itoa(params.Year),
			"month":          strconv.Itoa(params.Month),
		},
	}, result, count)

	ctx.JSON(http.StatusOK, GetMeasureStatisticsCUPSResponse{
		Pagination: Pagination{
			Links: struct {
				Next *string `json:"next,omitempty"`
				Prev *string `json:"prev,omitempty"`
				Self string  `json:"self"`
			}{
				Next: paginate.Links.Next,
				Prev: paginate.Links.Prev,
				Self: paginate.Links.Self,
			},
			Count:  paginate.Count,
			Limit:  paginate.Limit,
			Offset: paginate.Offset,
			Size:   paginate.Size,
		},
		Results: utils.MapSlice(result, DashboardMeasuresStatsCupsToResponse),
	})
}

func (c Controller) GetProcessMeasureStatistics(ctx *gin.Context, params GetProcessMeasureStatisticsParams) {
	response, err := c.srv.GetDashboardMeasuresStats.Handler(ctx.Request.Context(), services.GetDashboardMeasuresStatsDTO{
		DistributorID: params.DistributorId,
		Month:         params.Month,
		Year:          params.Year,
		Type:          string(params.Type),
	})
	if err != nil {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, utils.MapSlice(response, DashboardMeasuresStatsToResponse))
}
