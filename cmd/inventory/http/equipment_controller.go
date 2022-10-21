package http

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/inventory/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c Controller) ListMeasureEquipment(ctx *gin.Context, distributorId string, params ListMeasureEquipmentParams) {
	sortParams := make([]string, 0)

	if params.Sort != nil {
		sortParams = *params.Sort
	}

	query := ""
	if params.Q != nil {
		query = *params.Q
	}

	ds, count, err := c.srv.GetMeasureEquipmentsService.Handler(ctx, services.GetMeasureEquipmentDto{
		Q:             query,
		Limit:         params.Limit,
		Offset:        params.Offset,
		Sort:          inventory.NewSort(sortParams),
		DistributorID: distributorId,
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res := make([]MeasureEquipment, 0, cap(ds))

	for _, d := range ds {
		res = append(res, measureEquipmentToResponse(d))
	}

	links := server.GetLinksList(ctx, count, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
		Q:      params.Q,
		Sort:   params.Sort,
	})

	ctx.JSON(http.StatusOK, MeasureEquipmentsResponseList{
		Links: struct {
			Next *string `json:"next,omitempty"`
			Prev *string `json:"prev,omitempty"`
			Self string  `json:"self"`
		}{
			Next: links.Next,
			Prev: links.Prev,
			Self: links.Self,
		},
		Limit:   params.Limit,
		Offset:  params.Offset,
		Results: res,
		Count:   count,
		Size:    len(ds),
	})
}

func (c Controller) GetSmarkiaMeasureEquipment(ctx *gin.Context, smarkiaId string) {
	e, err := c.srv.GetMeasureEquipmentBySmarkiaIdService.Handler(ctx, services.GetMeasureEquipmentBySmarkiaDto{
		Id: smarkiaId,
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, measureEquipmentToResponse(e))
}

func (c Controller) GetMeasureEquipment(ctx *gin.Context, distributorId string, id string) {
	e, err := c.srv.GetMeasureEquipmentByIdService.Handler(ctx, services.GetMeasureEquipmentByIdDto{
		Id:            id,
		DistributorId: distributorId,
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, measureEquipmentToResponse(e))
}

func (c Controller) GetMeterConfig(ctx *gin.Context, distributorId string, params GetMeterConfigParams) {
	e, err := c.srv.GetMeterConfigByCupsService.Handler(ctx, services.GetMeterConfigByCupsDto{
		Cups:        *params.Cups,
		Distributor: distributorId,
		Date:        params.Date.AddDate(0, 0, 0),
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, meterConfigToResponse(e))

}
