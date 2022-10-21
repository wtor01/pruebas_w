package supply_point

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/services"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Controller struct {
	svc      *services.DashboardServices
	location *time.Location
}

func NewController(svc *services.DashboardServices, loc *time.Location) *Controller {
	return &Controller{svc: svc, location: loc}
}

func (c Controller) GetCurveGrossMeasureMeter(ctx *gin.Context, cups string, params GetCurveGrossMeasureMeterParams) {

	dashboardCurves, err := c.svc.DashboardSupplyPointCurvesService.Handler(ctx.Request.Context(), services.NewDashboardSupplyPointCurvesDto(
		cups,
		params.Distributor,
		string(params.CurveType),
		utils.ToUTC(params.Date.Time, c.location)))

	if err != nil {
		ctx.JSON(http.StatusExpectationFailed, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.MapSlice(dashboardCurves, toGrossCurves))
}

func (c Controller) GetGrossMeasureServicePoint(ctx *gin.Context, cups string, params GetGrossMeasureServicePointParams) {
	dashboardSupplyPoint, err := c.svc.DashboardMeasureSupplyPointService.Handler(ctx.Request.Context(),
		services.NewDashboardMeasureSupplyPointDto(
			cups,
			params.Distributor,
			utils.ToUTC(params.StartDate.Time, c.location),
			utils.ToUTC(params.EndDate.Time, c.location),
		))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, toGrossMeasureSupplyPoint(dashboardSupplyPoint))

}
