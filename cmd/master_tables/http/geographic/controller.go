package geographic

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/geographic"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/geographic/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Controller struct {
	gssrv *services.GeographicServices
}

// GetAllGeographicZones controller for Get all Geographic Zones
func (c2 Controller) GetAllGeographicZones(c *gin.Context, params GetAllGeographicZonesParams) {
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

	d, count, err := c2.gssrv.GetGeographicZones.Handler(c, services.GetAllDto{
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
	response := make([]GeographicZoneWithId, 0)
	links := server.GetLinksList(c, count, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
		Q:      params.Q,
		Sort:   params.Sort,
	})
	for _, dis := range d {
		response = append(response, GeographicToResponse(dis))

	}
	//response
	c.JSON(http.StatusOK, GeographicResponseList{
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

//InsertGeographicZone controller for update Geographic Zone
func (c2 Controller) InsertGeographicZone(c *gin.Context) {
	var req InsertGeographicZoneJSONRequestBody
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
	err = c2.gssrv.InsertGeographicZone.Handler(c, geographic.GeographicZones{
		Description: req.Description,
		Code:        req.Code,
		CreatedBy:   user.ID,
		CreatedAt:   time.Now(),
	})
	//check if there are an error and send
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "created"})

}

//GetGeographicZone controller for get geographic zone
func (c2 Controller) GetGeographicZone(c *gin.Context, geographicId string) {
	//generate dto and send to service
	gz, err := c2.gssrv.GetOneGeographicsZone.Handler(c, services.GetGeographicByIdDto{
		Id: geographicId,
	})
	//check for errors and send response
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, GeographicToResponse(gz))
}

//ModifyGeographicZone controller for modify geographic zone
func (c2 Controller) ModifyGeographicZone(c *gin.Context, geographicId string) {
	var req ModifyGeographicZoneJSONRequestBody
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
	err = c2.gssrv.ModifyGeographicZone.Handler(c, geographicId, geographic.GeographicZones{
		Description: req.Description,
		Code:        req.Code,
		UpdatedBy:   user.ID,
		UpdatedAt:   time.Now(),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "modified"})
}

//DeleteGeographicZone controller for geographic zone
func (c2 Controller) DeleteGeographicZone(c *gin.Context, geographicId string) {
	//generated dto and send to service
	err := c2.gssrv.DeleteGeographicsZone.Handler(c, services.DeleteGeographicByIdDto{
		Id: geographicId,
	})
	//check errors and send
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "deleted"})
}

func NewController(gssrv *services.GeographicServices) *Controller {
	return &Controller{gssrv: gssrv}
}
