package dashboard_stats

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/services"
	"strconv"

	"bitbucket.org/sercide/data-ingestion/pkg/server"

	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	dashboardService *services.GrossMeasuresStatsDashboardServices
}

func (c Controller) GetProcessMeasureStatisticsByCups(ctx *gin.Context, params GetProcessMeasureStatisticsByCupsParams) {
	var offset int
	if params.Offset != nil {
		offset = *params.Offset
	}
	result, err := c.dashboardService.ListDashboardMeasuresStatsSerialNumber.Handler(ctx.Request.Context(), services.ListDashboardMeasuresStatsSerialNumberDTO{
		DistributorId: params.DistributorId,
		Month:         params.Month,
		Year:          params.Year,
		Type:          string(params.Type),
		Ghost:         params.Ghost,
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
			"month":          strconv.Itoa(params.Month),
			"year":           strconv.Itoa(params.Year),
			"type":           string(params.Type),
			"ghost":          strconv.FormatBool(params.Ghost),
		},
	}, result.Data, result.Count)

	ctx.JSON(http.StatusOK, GetGrossMeasureDashboardCUPSResponse{
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
		Results: utils.MapSlice(result.Data, GrossMeasureStatisticsByCupsToResponse),
	})
}

func (c Controller) GetDashboardGrossMeasure(ctx *gin.Context, params GetDashboardGrossMeasureParams) {
	response, err := c.dashboardService.GetDashboardMeasuresStats.Handler(ctx.Request.Context(), services.GetDashboardMeasuresStatsDTO{
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
