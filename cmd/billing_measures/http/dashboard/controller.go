package dashboard

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/services"
	"bitbucket.org/sercide/data-ingestion/internal/common/config"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Controller struct {
	srv *services.BillingMeasuresDashboardService
	cnf config.Config
}

func (c Controller) GetBillingMeasuresResumeById(ctx *gin.Context, Id string) {
	result, err := c.srv.GetClosureResumeDashboard.Handler(ctx.Request.Context(), services.GetClosureResumeDashboardDTO{BillingMeasureID: Id})
	if err != nil {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, toFiscalClosureResume(result))
}

func (c Controller) CreateBillingMeasuresMVH(ctx *gin.Context) {
	var req MVH
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := c.srv.ExecuteMvh.Handler(ctx.Request.Context(), services.DtoServiceExecuteMvh{
		Cups:          req.Cups,
		DistributorId: req.DistributorId,
		Date:          req.Time.Time,
		Location:      c.cnf.LocalLocation,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": "success"})
}

func (c Controller) SearchDistributorFiscalBillingMeasures(ctx *gin.Context, params SearchDistributorFiscalBillingMeasuresParams) {
	result, err := c.srv.SearchFiscalBillingMeasures.Handler(ctx.Request.Context(), services.SearchFiscalBillingMeasuresDashboardDTO{
		Cups:          params.Cups,
		DistributorId: params.DistributorId,
		StartDate:     time.Date(params.StartDate.Year(), params.StartDate.Month(), params.StartDate.Day(), 0, 0, 0, 0, c.srv.SearchFiscalBillingMeasures.Location).UTC(),
		EndDate:       time.Date(params.EndDate.Year(), params.EndDate.Month(), params.EndDate.Day()+1, 0, 0, 0, 0, c.srv.SearchFiscalBillingMeasures.Location).UTC(),
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, utils.MapSlice(result, fiscalBillingMeasuresToResponse))

}

func (c Controller) GetBillingMeasureDashboardSummary(ctx *gin.Context, params GetBillingMeasureDashboardSummaryParams) {
	summaryResult, err := c.srv.DashboardSummaryService.Handler(ctx.Request.Context(), services.NewDashboardSummaryDto(
		params.DistributorId,
		string(params.MeterType),
		utils.ToUTC(params.StartDate.Time, c.srv.Location),
		utils.ToUTC(params.EndDate.Time, c.srv.Location),
	))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, toFiscalMeasureSummary(summaryResult))
}

func (sc Controller) GetBillingMeasureTaxMeasurebycups(ctx *gin.Context, params GetBillingMeasureTaxMeasurebycupsParams) {

	result, err := sc.srv.TaxMeasuresByCups.Handler(ctx.Request.Context(), billing_measures.QueryBillingMeasuresTax{
		Offset:        params.Offset,
		Limit:         params.Limit,
		DistributorId: params.DistributorId,
		MeasureType:   string(params.MeasureType),
		StartDate:     time.Date(params.StartDate.Year(), params.StartDate.Month(), params.StartDate.Day(), 0, 0, 0, 0, sc.cnf.LocalLocation).UTC(),
		EndDate:       time.Date(params.EndDate.Year(), params.EndDate.Month(), params.EndDate.Day(), 0, 0, 0, 0, sc.cnf.LocalLocation).UTC(),
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paginate := server.GetPaginate(ctx, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
	}, result.Data, result.Count)

	data := utils.MapSlice(result.Data, func(item billing_measures.BillingMeasuresTax) FiscalMeasureByCups {
		r := FiscalMeasureByCups{
			Cups: &item.Cups,
		}

		if !item.StartDate.IsZero() {
			r.StartDate = &item.StartDate
		}
		if !item.EndDate.IsZero() {
			r.EndDate = &item.EndDate
		}

		atrCalification := FiscalMeasureByCupsAtrCalification(item.ExecutionSummary.BalanceType)
		if atrCalification != "" {
			r.AtrCalification = &atrCalification
		}

		atrType := FiscalMeasureByCupsAtrType(item.ExecutionSummary.BalanceOrigin)
		if atrType != "" {
			r.AtrType = &atrType
		}

		curveCalification := FiscalMeasureByCupsCurveCalification(item.ExecutionSummary.CurveType)
		if curveCalification != "" {
			r.CurveCalification = &curveCalification
		}

		curveStatus := FiscalMeasureByCupsCurveStatus(item.ExecutionSummary.CurveStatus)
		if curveStatus != "" {
			r.CurveStatus = &curveStatus
		}

		return r
	})

	ctx.JSON(http.StatusOK, FiscalMeasureByCupsResponse{
		Pagination: Pagination{
			Links:  paginate.Links,
			Count:  paginate.Count,
			Limit:  paginate.Limit,
			Offset: paginate.Offset,
			Size:   paginate.Size,
		},
		Results: data,
	})
}

func (sc Controller) GetBillingMeasureTaxMeasure(ctx *gin.Context, params GetBillingMeasureTaxMeasureParams) {
}
