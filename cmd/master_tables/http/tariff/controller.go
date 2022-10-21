package tariff

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/tariff/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Controller struct {
	tfsrv *services.TariffServices
}

func (c Controller) GetAllTariffs(ctx *gin.Context, params GetAllTariffsParams) {
	// get params
	sortParams := make([]string, 0)
	user, _ := middleware.GetAuthUser(ctx)
	if params.Sort != nil {
		sortParams = *params.Sort
	}
	query := ""
	if params.Q != nil {
		query = *params.Q
	}

	d, count, err := c.tfsrv.GetTariffs.Handler(ctx, services.GetAllDto{
		Q:           query,
		Limit:       params.Limit,
		Offset:      params.Offset,
		Sort:        inventory.NewSort(sortParams),
		CurrentUser: user,
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	//make response
	response := make([]Tariffs, 0)
	links := server.GetLinksList(ctx, count, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
		Q:      params.Q,
		Sort:   params.Sort,
	})
	for _, dis := range d {
		response = append(response, TariffToResponse(dis))

	}
	//response
	ctx.JSON(http.StatusOK, TariffResponseList{
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
		Results: response,
		Count:   count,
		Size:    len(d),
	})
}
func (c Controller) InsertTariffs(ctx *gin.Context) {
	var req InsertTariffsJSONRequestBody
	//Get body and save in a struct
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetAuthUser(ctx)

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//convert body to service Geographic_zones and send it
	err = c.tfsrv.InsertTariff.Handler(ctx, tariff.Tariffs{
		Id:           req.TariffCode,
		CodeOdos:     req.CodeOdos,
		CodeOne:      req.CodeOne,
		Description:  req.Description,
		GeographicId: req.GeographicId,
		Periods:      req.Periods,
		CalendarId:   req.CalendarCode,
		TensionLevel: string(req.TensionLevel),
		CreatedBy:    user.ID,
		CreatedAt:    time.Now(),
		Coef:         string(req.Coef),
	})
	//check if there are an error and send
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": "created"})
}
func (c Controller) ModifyTariffs(ctx *gin.Context, tariffCode string) {
	var req ModifyTariffsJSONRequestBody
	// convert body to ModifyGeographicZoneJSONRequestBody
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetAuthUser(ctx)

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	//generate serrvice geographic struct and send to service
	err = c.tfsrv.ModifyTariff.Handler(ctx, tariffCode, tariff.Tariffs{
		Id:           req.TariffCode,
		Description:  req.Description,
		UpdatedBy:    user.ID,
		UpdatedAt:    time.Now(),
		CodeOne:      req.CodeOne,
		CodeOdos:     req.CodeOdos,
		TensionLevel: string(req.TensionLevel),
		GeographicId: req.GeographicId,
		Periods:      req.Periods,
		CalendarId:   req.CalendarCode,
		Coef:         string(req.Coef),
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": "modified"})
}
func (c Controller) DeleteTariff(ctx *gin.Context, tariffId string) {
	//generated dto and send to service
	err := c.tfsrv.DeleteTariff.Handler(ctx, services.DeleteTariffByCodeDto{
		Id: tariffId,
	})
	//check errors and send
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": "deleted"})
}
func (c Controller) GetTariffs(ctx *gin.Context, tariffId string) {
	//generate dto and send to service
	tf, err := c.tfsrv.GetOneTariff.Handler(ctx, services.GetTariffByCodeDto{
		ID: tariffId,
	})
	//check for errors and send response
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, TariffToResponse(tf))
}

func (c Controller) GetTariffsCalendar(ctx *gin.Context, tariffId string, params GetTariffsCalendarParams) {
	// get params
	sortParams := make([]string, 0)
	user, _ := middleware.GetAuthUser(ctx)
	if params.Sort != nil {
		sortParams = *params.Sort
	}
	query := ""
	if params.Q != nil {
		query = *params.Q
	}

	d, count, err := c.tfsrv.GetAllTariffCalendar.Handler(ctx, tariffId, services.GetAllDto{
		Q:           query,
		Limit:       params.Limit,
		Offset:      params.Offset,
		Sort:        inventory.NewSort(sortParams),
		CurrentUser: user,
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	//make response
	response := make([]TariffCalendar, 0)
	links := server.GetLinksList(ctx, count, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
		Q:      params.Q,
		Sort:   params.Sort,
	})
	for _, dis := range d {
		response = append(response, TCtoResponse(dis))

	}
	//response
	ctx.JSON(http.StatusOK, TariffCalendarResponseList{
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
		Results: response,
		Count:   count,
		Size:    len(d),
	})
}

func TCtoResponse(tc tariff.TariffsCalendar) TariffCalendar {
	layout := "02-01-2006"
	null := "01-01-0001"
	sd := tc.StartDate.Format(layout)
	ed := tc.EndDate.Format(layout)
	if ed == null {
		ed = ""
	}
	return TariffCalendar{
		CalendarCode:   tc.CalendarId,
		TariffCode:     tc.TariffId,
		EndDate:        ed,
		StartDate:      sd,
		GeographicCode: tc.GeoGraphicCode,
	}
}

func NewController(tfsrv *services.TariffServices) *Controller {
	return &Controller{tfsrv: tfsrv}
}
