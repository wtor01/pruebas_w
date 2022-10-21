package dashboard

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/services"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Controller struct {
	dashboardService *services.DashboardService
}

func NewController(svc *services.DashboardService) *Controller {
	return &Controller{svc}
}

func (c Controller) GetMeasureDashboard(ctx *gin.Context, params GetMeasureDashboardParams) {
	response, err := c.dashboardService.Handle(ctx, services.NewDashboardServiceDTO(
		params.DistributorId,
		time.Date(params.StartDate.Year(), params.StartDate.Month(), params.StartDate.Day(), 0, 0, 0, 0, c.dashboardService.Location),
		time.Date(params.EndDate.Year(), params.EndDate.Month(), params.EndDate.Day(), 0, 0, 0, 0, c.dashboardService.Location),
	))

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	daily := make([]struct {
		Date        openapi_types.Date `json:"date"`
		Others      Others             `json:"others"`
		Telegestion Telegestion        `json:"telegestion"`
		Telemedida  Telemedida         `json:"telemedida"`
	}, 0, cap(response.Daily))

	for _, d := range response.Daily {
		daily = append(daily, struct {
			Date        openapi_types.Date `json:"date"`
			Others      Others             `json:"others"`
			Telegestion Telegestion        `json:"telegestion"`
			Telemedida  Telemedida         `json:"telemedida"`
		}{
			Date: openapi_types.Date{d.Date},
			Others: Others{
				Closing: Data{
					Invalid:          d.Others.Closing.Invalid,
					MeasuresShouldBe: d.Others.Closing.MeasuresShouldBe,
					Supervise:        d.Others.Closing.Supervise,
					Valid:            d.Others.Closing.Valid,
				},
			},
			Telegestion: Telegestion{
				Closing: Data{
					Invalid:          d.Telegestion.Closing.Invalid,
					MeasuresShouldBe: d.Telegestion.Closing.MeasuresShouldBe,
					Supervise:        d.Telegestion.Closing.Supervise,
					Valid:            d.Telegestion.Closing.Valid,
				},
				Curva: Data{
					Invalid:          d.Telegestion.Curva.Invalid,
					MeasuresShouldBe: d.Telegestion.Curva.MeasuresShouldBe,
					Supervise:        d.Telegestion.Curva.Supervise,
					Valid:            d.Telegestion.Curva.Valid,
				},
				Resumen: Data{
					Invalid:          d.Telegestion.Resumen.Invalid,
					MeasuresShouldBe: d.Telegestion.Resumen.MeasuresShouldBe,
					Supervise:        d.Telegestion.Resumen.Supervise,
					Valid:            d.Telegestion.Resumen.Valid,
				},
			},
			Telemedida: Telemedida{
				Closing: Data{
					Invalid:          d.Telemedida.Closing.Invalid,
					MeasuresShouldBe: d.Telemedida.Closing.MeasuresShouldBe,
					Supervise:        d.Telemedida.Closing.Supervise,
					Valid:            d.Telemedida.Closing.Valid,
				},
				Curva: Data{
					Invalid:          d.Telemedida.Curva.Invalid,
					MeasuresShouldBe: d.Telemedida.Curva.MeasuresShouldBe,
					Supervise:        d.Telemedida.Curva.Supervise,
					Valid:            d.Telemedida.Curva.Valid,
				},
			},
		},
		)
	}

	ctx.JSON(http.StatusOK, DashboardMeasure{
		Totals: struct {
			Others      Others      `json:"others"`
			Telegestion Telegestion `json:"telegestion"`
			Telemedida  Telemedida  `json:"telemedida"`
		}{
			Others: Others{
				Closing: Data{
					Invalid:          response.Totals.Others.Closing.Invalid,
					MeasuresShouldBe: response.Totals.Others.Closing.MeasuresShouldBe,
					Supervise:        response.Totals.Others.Closing.Supervise,
					Valid:            response.Totals.Others.Closing.Valid,
				},
			},
			Telegestion: Telegestion{
				Closing: Data{
					Invalid:          response.Totals.Telegestion.Closing.Invalid,
					MeasuresShouldBe: response.Totals.Telegestion.Closing.MeasuresShouldBe,
					Supervise:        response.Totals.Telegestion.Closing.Supervise,
					Valid:            response.Totals.Telegestion.Closing.Valid,
				},
				Curva: Data{
					Invalid:          response.Totals.Telegestion.Curva.Invalid,
					MeasuresShouldBe: response.Totals.Telegestion.Curva.MeasuresShouldBe,
					Supervise:        response.Totals.Telegestion.Curva.Supervise,
					Valid:            response.Totals.Telegestion.Curva.Valid,
				},
				Resumen: Data{
					Invalid:          response.Totals.Telegestion.Resumen.Invalid,
					MeasuresShouldBe: response.Totals.Telegestion.Resumen.MeasuresShouldBe,
					Supervise:        response.Totals.Telegestion.Resumen.Supervise,
					Valid:            response.Totals.Telegestion.Resumen.Valid,
				},
			},
			Telemedida: Telemedida{
				Closing: Data{
					Invalid:          response.Totals.Telemedida.Closing.Invalid,
					MeasuresShouldBe: response.Totals.Telemedida.Closing.MeasuresShouldBe,
					Supervise:        response.Totals.Telemedida.Closing.Supervise,
					Valid:            response.Totals.Telemedida.Closing.Valid,
				},
				Curva: Data{
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
