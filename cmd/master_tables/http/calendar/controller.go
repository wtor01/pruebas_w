package calendar

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/calendar/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

type Controller struct {
	csrv *services.CalendarServices
}

func (c2 Controller) GetAllCalendars(c *gin.Context, params GetAllCalendarsParams) {
	// get params
	sortParams := make([]string, 0)
	user, _ := middleware.GetAuthUser(c)
	if params.Sort != nil {
		sortParams = *params.Sort
	}
	query := ""
	if params.Q != nil {
		query = *params.Q
	}

	d, count, err := c2.csrv.GetCalendars.Handler(c, services.GetAllDto{
		Q:           query,
		Limit:       params.Limit,
		Offset:      params.Offset,
		Sort:        inventory.NewSort(sortParams),
		CurrentUser: user,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	//make response
	response := make([]CalendarWithId, 0)
	links := server.GetLinksList(c, count, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
		Q:      params.Q,
		Sort:   params.Sort,
	})
	for _, dis := range d {
		response = append(response, CalendarToResponse(dis))

	}
	//response
	c.JSON(http.StatusOK, CalendarResponseList{
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

func (c2 Controller) InsertCalendars(c *gin.Context) {
	var req InsertCalendarsJSONRequestBody
	//Get body and save in a struct
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetAuthUser(c)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//convert body to service Geographic_zones and send it
	err = c2.csrv.InsertCalendar.Handler(c, calendar.Calendar{
		Description:    req.Description,
		Id:             req.Code,
		GeographicCode: req.GeographicCode,
		CreatedBy:      user.ID,
		CreatedAt:      time.Now(),
	})
	//check if there are an error and send
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "created"})
}
func (c2 Controller) GetCalendar(c *gin.Context, calendarId string) {
	//generate dto and send to service
	gz, err := c2.csrv.GetOneCalendar.Handler(c, services.GetCalendarByIdDto{
		Id: calendarId,
	})
	//check for errors and send response
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, CalendarToResponse(gz))
}
func (c2 Controller) ModifyCalendar(c *gin.Context, calendarId string) {
	var req ModifyCalendarJSONRequestBody
	// convert body to ModifyGeographicZoneJSONRequestBody
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetAuthUser(c)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//generate serrvice geographic struct and send to service
	err = c2.csrv.PutCalendar.Handler(c, calendarId, calendar.Calendar{
		Description:    req.Description,
		GeographicCode: req.GeographicCode,
		Periods:        req.Periods,
		UpdatedBy:      user.ID,
		UpdatedAt:      time.Now(),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "modified"})
}
func (c2 Controller) DeleteCalendar(c *gin.Context, calendarId string) {
	//generated dto and send to service
	err := c2.csrv.DeleteCalendar.Handler(c, services.DeleteCalendarByIdDto{
		Id: calendarId,
	})
	//check errors and send
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "deleted"})
}

func (c2 Controller) GetAllPeriodsCalendars(c *gin.Context, calendarId string, params GetAllPeriodsCalendarsParams) {
	// get params
	sortParams := make([]string, 0)
	user, _ := middleware.GetAuthUser(c)
	if params.Sort != nil {
		sortParams = *params.Sort
	}
	query := ""
	if params.Q != nil {
		query = *params.Q
	}

	d, count, err := c2.csrv.GetPeriodCalendar.Handler(c, calendarId, services.GetAllDto{
		Q:           query,
		Limit:       params.Limit,
		Offset:      params.Offset,
		Sort:        inventory.NewSort(sortParams),
		CurrentUser: user,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	//make response
	response := make([]CalendarPeriod, 0)
	links := server.GetLinksList(c, count, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
		Q:      params.Q,
		Sort:   params.Sort,
	})
	for _, dis := range d {
		response = append(response, PCtoResponse(dis))

	}
	//response
	c.JSON(http.StatusOK, PeriodResponseList{
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
func (c2 Controller) GetPeriod(c *gin.Context, periodId string) {
	//generate dto and send to service
	gz, err := c2.csrv.GetPeriodById.Handler(c, services.GetPeriodByIdDto{
		Id: uuid.Must(uuid.FromString(periodId)),
	})
	//check for errors and send response
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, PCtoResponse(gz))
}
func (c2 Controller) InsertPeriods(c *gin.Context, calendarId string) {
	var req InsertPeriodsJSONRequestBody
	//Get body and save in a struct
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetAuthUser(c)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//convert body to service Geographic_zones and send it
	err = c2.csrv.InsertPeriod.Handler(c, calendarId, calendar.PeriodCalendar{
		Description:  req.Description,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		StartHour:    req.StartHour,
		EndHour:      req.EndHour,
		Year:         req.Year,
		Energy:       req.Energy,
		Power:        req.Power,
		DayType:      calendar.DayType(req.DayType),
		UpdatedAt:    time.Now(),
		CreatedBy:    user.ID,
		PeriodNumber: string(req.PeriodNumber),
	})
	//check if there are an error and send
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "created"})
}
func (c2 Controller) ModifyPeriodCalendar(c *gin.Context, periodCode string) {
	var req ModifyPeriodCalendarJSONRequestBody
	// convert body to ModifyGeographicZoneJSONRequestBody
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetAuthUser(c)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//generate serrvice geographic struct and send to service
	err = c2.csrv.PutPeriodCalendar.Handler(c, periodCode, calendar.PeriodCalendar{
		Description:  req.Description,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		StartHour:    req.StartHour,
		EndHour:      req.EndHour,
		Year:         req.Year,
		DayType:      calendar.DayType(req.DayType),
		UpdatedAt:    time.Now(),
		UpdatedBy:    user.ID,
		Energy:       req.Energy,
		Power:        req.Power,
		PeriodNumber: string(req.PeriodNumber),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "modified"})
}
func (c2 Controller) DeletePeriod(c *gin.Context, periodCode string) {
	//generated dto and send to service
	err := c2.csrv.DeletePeriod.Handler(c, services.DeletePeriodByCodeDto{
		Code: periodCode,
	})
	//check errors and send
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "deleted"})
}

func NewController(csrv *services.CalendarServices) *Controller {
	return &Controller{csrv: csrv}
}
