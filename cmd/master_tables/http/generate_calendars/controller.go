package generate_calendars

import (
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	srv *services.CalendarPeriodsPubSubServices
}

func NewController(s *services.CalendarPeriodsPubSubServices) *Controller {
	return &Controller{
		srv: s,
	}
}

func (c Controller) GenerateCalendarsPeriods(ctx *gin.Context) {
	err := c.srv.CalendarPeriodGenerateService.Handler(ctx.Request.Context())

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusCreated, struct {
		msg string
	}{
		msg: "OK",
	})
	return
}

func (c Controller) GenerateFestiveDays(ctx *gin.Context) {
	err := c.srv.FestiveDaysGenerateService.Handler(ctx.Request.Context())

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusCreated, struct {
		msg string
	}{
		msg: "OK",
	})
	return
}
