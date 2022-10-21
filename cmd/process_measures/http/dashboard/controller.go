package dashboard

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"fmt"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	dashboardService *services.DashboardServices
}

func NewController(dashboardService *services.DashboardServices) *Controller {
	return &Controller{dashboardService: dashboardService}
}

func (c Controller) GetProcessMeasureDashboard(ctx *gin.Context, params GetProcessMeasureDashboardParams) {
	response, err := c.dashboardService.GetDashboard.Handle(ctx.Request.Context(), services.NewGetDashboardMeasuresDTO(
		params.DistributorId,
		utils.ToUTC(params.StartDate.Time, c.dashboardService.Location),
		utils.ToUTC(params.EndDate.Time, c.dashboardService.Location),
	))

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	daily := make([]struct {
		Date        openapi_types.Date            `json:"date"`
		Others      DashboardProcessMeasureOthers `json:"others"`
		Telegestion DashboardProcessMeasureTLG    `json:"telegestion"`
		Telemedida  DashboardProcessMeasureTLM    `json:"telemedida"`
	}, 0, cap(response.Daily))

	for _, d := range response.Daily {
		daily = append(daily, struct {
			Date        openapi_types.Date            `json:"date"`
			Others      DashboardProcessMeasureOthers `json:"others"`
			Telegestion DashboardProcessMeasureTLG    `json:"telegestion"`
			Telemedida  DashboardProcessMeasureTLM    `json:"telemedida"`
		}{
			Date: openapi_types.Date{d.Date},
			Others: DashboardProcessMeasureOthers{
				Closing: DashboardProcessMeasureData{
					Invalid:          d.Others.Closing.Invalid,
					MeasuresShouldBe: d.Others.Closing.MeasuresShouldBe,
					Supervise:        d.Others.Closing.Supervise,
					Valid:            d.Others.Closing.Valid,
				},
			},
			Telegestion: DashboardProcessMeasureTLG{
				Closing: DashboardProcessMeasureData{
					Invalid:          d.Telegestion.Closing.Invalid,
					MeasuresShouldBe: d.Telegestion.Closing.MeasuresShouldBe,
					Supervise:        d.Telegestion.Closing.Supervise,
					Valid:            d.Telegestion.Closing.Valid,
				},
				Curva: DashboardProcessMeasureData{
					Invalid:          d.Telegestion.Curva.Invalid,
					MeasuresShouldBe: d.Telegestion.Curva.MeasuresShouldBe,
					Supervise:        d.Telegestion.Curva.Supervise,
					Valid:            d.Telegestion.Curva.Valid,
				},
				Resumen: DashboardProcessMeasureData{
					Invalid:          d.Telegestion.Resumen.Invalid,
					MeasuresShouldBe: d.Telegestion.Resumen.MeasuresShouldBe,
					Supervise:        d.Telegestion.Resumen.Supervise,
					Valid:            d.Telegestion.Resumen.Valid,
				},
			},
			Telemedida: DashboardProcessMeasureTLM{
				Closing: DashboardProcessMeasureData{
					Invalid:          d.Telemedida.Closing.Invalid,
					MeasuresShouldBe: d.Telemedida.Closing.MeasuresShouldBe,
					Supervise:        d.Telemedida.Closing.Supervise,
					Valid:            d.Telemedida.Closing.Valid,
				},
				Curva: DashboardProcessMeasureData{
					Invalid:          d.Telemedida.Curva.Invalid,
					MeasuresShouldBe: d.Telemedida.Curva.MeasuresShouldBe,
					Supervise:        d.Telemedida.Curva.Supervise,
					Valid:            d.Telemedida.Curva.Valid,
				},
			},
		},
		)
	}

	ctx.JSON(http.StatusOK, DashboardProcessMeasure{
		Totals: struct {
			Others      DashboardProcessMeasureOthers `json:"others"`
			Telegestion DashboardProcessMeasureTLG    `json:"telegestion"`
			Telemedida  DashboardProcessMeasureTLM    `json:"telemedida"`
		}{
			Others: DashboardProcessMeasureOthers{
				Closing: DashboardProcessMeasureData{
					Invalid:          response.Totals.Others.Closing.Invalid,
					MeasuresShouldBe: response.Totals.Others.Closing.MeasuresShouldBe,
					Supervise:        response.Totals.Others.Closing.Supervise,
					Valid:            response.Totals.Others.Closing.Valid,
				},
			},
			Telegestion: DashboardProcessMeasureTLG{
				Closing: DashboardProcessMeasureData{
					Invalid:          response.Totals.Telegestion.Closing.Invalid,
					MeasuresShouldBe: response.Totals.Telegestion.Closing.MeasuresShouldBe,
					Supervise:        response.Totals.Telegestion.Closing.Supervise,
					Valid:            response.Totals.Telegestion.Closing.Valid,
				},
				Curva: DashboardProcessMeasureData{
					Invalid:          response.Totals.Telegestion.Curva.Invalid,
					MeasuresShouldBe: response.Totals.Telegestion.Curva.MeasuresShouldBe,
					Supervise:        response.Totals.Telegestion.Curva.Supervise,
					Valid:            response.Totals.Telegestion.Curva.Valid,
				},
				Resumen: DashboardProcessMeasureData{
					Invalid:          response.Totals.Telegestion.Resumen.Invalid,
					MeasuresShouldBe: response.Totals.Telegestion.Resumen.MeasuresShouldBe,
					Supervise:        response.Totals.Telegestion.Resumen.Supervise,
					Valid:            response.Totals.Telegestion.Resumen.Valid,
				},
			},
			Telemedida: DashboardProcessMeasureTLM{
				Closing: DashboardProcessMeasureData{
					Invalid:          response.Totals.Telemedida.Closing.Invalid,
					MeasuresShouldBe: response.Totals.Telemedida.Closing.MeasuresShouldBe,
					Supervise:        response.Totals.Telemedida.Closing.Supervise,
					Valid:            response.Totals.Telemedida.Closing.Valid,
				},
				Curva: DashboardProcessMeasureData{
					Invalid:          response.Totals.Telemedida.Curva.Invalid,
					MeasuresShouldBe: response.Totals.Telemedida.Curva.MeasuresShouldBe,
					Supervise:        response.Totals.Telemedida.Curva.Supervise,
					Valid:            response.Totals.Telemedida.Curva.Valid,
				},
			},
		},
		Daily: daily,
	})
}

func (c Controller) GetCurveProcessServicePoint(ctx *gin.Context, params GetCurveProcessServicePointParams) {
	response, err := c.dashboardService.SearchServicePointDashboardCurves.Handler(ctx.Request.Context(), services.NewSearchServicePointProcessMeasuresDashboardCurvesDTO(
		params.Distributor,
		params.Cups,
		utils.ToUTC(params.StartDate.Time, c.dashboardService.Location),
		utils.ToUTC(params.EndDate.Time, c.dashboardService.Location),
		measures.MeasureCurveReadingType(params.CurveType),
	))

	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, utils.MapSlice(response, func(item process_measures.ServicePointDashboardCurve) CurveProcessServicePoint {
		return CurveProcessServicePoint{
			Date:   item.Date.UTC().Format("2006-01-02T15:04:05-0700"),
			Status: string(item.Status),
			Values: CurveValues{
				AE: item.Values.AE,
				AI: item.Values.AI,
				R1: item.Values.R1,
				R2: item.Values.R2,
				R3: item.Values.R3,
				R4: item.Values.R4,
			},
		}
	}))
}

func (c Controller) GetDashboardProcessServicePoint(ctx *gin.Context, params GetDashboardProcessServicePointParams) {
	response, err := c.dashboardService.SearchServicePointProcessMeasuresDashboard.Handler(ctx.Request.Context(), services.NewSearchServicePointProcessMeasuresDashboardDTO(
		params.Distributor,
		params.Cups,
		utils.ToUTC(params.StartDate.Time, c.dashboardService.Location),
		utils.ToUTC(params.EndDate.Time, c.dashboardService.Location),
	))

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, GetDashboardProcessServicePointToResponse(response))
}

func (c Controller) GetProcessMeasureDashboardList(ctx *gin.Context, params GetProcessMeasureDashboardListParams) {

	var offset int64
	if params.Offset != nil {
		offset = int64(*params.Offset)
	}

	response, err := c.dashboardService.ListDashboardCups.Handle(ctx.Request.Context(), services.NewListDashboardCupsDTO(
		params.DistributorId,
		int64(params.Limit),
		offset,
		string(params.Type),
		utils.ToUTC(params.StartDate.Time, c.dashboardService.Location),
		utils.ToUTC(params.EndDate.Time, c.dashboardService.Location),
	))

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	linkList := server.GetLinksList(ctx, response.Total, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
		Q:      nil,
		Sort:   nil,
		Parameters: map[string]string{
			"distributor_id": params.DistributorId,
			"type":           string(params.Type),
			"start_date":     params.StartDate.String(),
			"end_date":       params.EndDate.String(),
		},
		Host: ctx.Request.Host,
	})

	results := GetProcessMeasureDashboardListToResponse(response)

	pagination := Pagination{
		Links: struct {
			Next *string `json:"next,omitempty"`
			Prev *string `json:"prev,omitempty"`
			Self string  `json:"self"`
		}{
			Next: linkList.Next,
			Prev: linkList.Prev,
			Self: linkList.Self,
		},
		Count:  response.Total,
		Limit:  params.Limit,
		Offset: params.Offset,
		Size:   len(results),
	}
	ctx.JSON(http.StatusOK, DashboardCupsProcessMeasure{
		Pagination: pagination,
		Results:    results,
	})
}
