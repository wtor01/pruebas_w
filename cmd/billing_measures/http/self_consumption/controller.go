package self_consumption

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	srv *services.SelfConsumptionServices
}

// SearchActivesSelfConsumptionUnitConfigByDistributorId Function to list the self-consumption units of a distributor active on a certain date
func (sc Controller) SearchActivesSelfConsumptionUnitConfigByDistributorId(ctx *gin.Context, distributorId string, params SearchActivesSelfConsumptionUnitConfigByDistributorIdParams) {

	var offset int

	if params.Offset == nil {
		offset = 0
	} else {
		offset = *params.Offset
	}

	result, count, err := sc.srv.GetSelfConsumptionActiveByDistributor.Handler(ctx.Request.Context(), billing_measures.GetSelfConsumptionByDistributortDto{
		DistributorId: distributorId,
		Date:          params.Date.Time,
		Limit:         params.Limit,
		Offset:        offset,
	})

	res := make([]SelfConsumptionUnitConfig, 0, cap(result))

	for _, d := range result {
		res = append(res, selfConsumptionToResponse(d))
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paginate := server.GetPaginate(ctx, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
	}, result, count)

	ctx.JSON(http.StatusOK, ListSelfConsumptionUnitConfigs{
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
		Results: res,
	})

}

// SearchSelfConsumptionUnitConfig Function to obtain the self-consumption units to which a cup belongs for a given date
func (sc Controller) SearchSelfConsumptionUnitConfig(ctx *gin.Context, distributorId string, params SearchSelfConsumptionUnitConfigParams) {

	result, err := sc.srv.GetSelfConsumptionByCup.Handler(ctx.Request.Context(), billing_measures.GetSelfConsumptionByCUP{
		DistributorId: distributorId,
		CUP:           params.Cups,
		Date:          params.Date.Time,
	})

	var response []SelfConsumptionUnitConfig

	response = append(response, selfConsumptionToResponse(result))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (sc Controller) GetSelfConsumptionByCau(ctx *gin.Context, distributorId string, cau string, params GetSelfConsumptionByCauParams) {
	billingConsumption, err := sc.srv.GetBillingSelfConsumptionByCauService.Handler(ctx.Request.Context(),
		services.NewGetBillingSelfConsumptionDto(
			cau,
			distributorId,
			utils.ToUTC(params.StartDate.Time, sc.srv.Location),
			utils.ToUTC(params.EndDate.Time, sc.srv.Location),
		))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.MapSlice(billingConsumption, billingSelfConsumptionToResponse))
}
