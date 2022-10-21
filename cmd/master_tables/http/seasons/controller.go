package seasons

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/seasons/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

type Controller struct {
	gssrv *services.SeasonsServices
}

func NewController(gssrv *services.SeasonsServices) *Controller {
	return &Controller{gssrv: gssrv}
}

func (c2 Controller) GetAllSeasons(c *gin.Context, params GetAllSeasonsParams) {
	sortParams := make([]string, 0)
	user, _ := middleware.GetAuthUser(c)
	if params.Sort != nil {
		sortParams = *params.Sort
	}
	query := ""
	if params.Q != nil {
		query = *params.Q
	}

	d, count, err := c2.gssrv.GetSeasons.Handler(c, services.GetAllDto{
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

	response := make([]Seasons, 0)
	links := server.GetLinksList(c, count, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
		Q:      params.Q,
		Sort:   params.Sort,
	})

	for _, dis := range d {
		response = append(response, SeasonsToResponse(dis))
	}

	c.JSON(http.StatusOK, SeasonsResponseList{
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
func (c2 Controller) GetSeason(c *gin.Context, seasonId string) {
	gz, err := c2.gssrv.GetSeasonById.Handler(c, services.GetSeasonByIdDto{
		Id: seasonId,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, SeasonsToResponse(gz))
}
func (c2 Controller) InsertSeason(c *gin.Context) {
	var req InsertSeasonJSONRequestBody

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := middleware.GetAuthUser(c)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	err = c2.gssrv.InsertSeason.Handler(c, seasons.Seasons{
		Name:           req.Name,
		Description:    req.Description,
		GeographicCode: req.GeographicCode,
		CreatedBy:      user.ID,
		CreatedAt:      time.Now(),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "created"})

}
func (c2 Controller) ModifySeason(c *gin.Context, seasonId string) {
	var req ModifySeasonJSONRequestBody
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := middleware.GetAuthUser(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	seasonUUID := uuid.Must(uuid.FromString(seasonId))
	err = c2.gssrv.ModifySeason.Handler(c, seasonUUID, seasons.Seasons{
		Name:           req.Name,
		Description:    req.Description,
		GeographicCode: req.GeographicCode,
		CreatedBy:      user.ID,
		CreatedAt:      time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": "modified"})
}
func (c2 Controller) DeleteSeason(c *gin.Context, seasonId string) {
	err := c2.gssrv.DeleteSeason.Handler(c, services.DeleteSeasonByIdDto{
		Id: uuid.Must(uuid.FromString(seasonId)),
	})
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "deleted"})
}

func (c2 Controller) GetAllDayTypes(c *gin.Context, seasonId string, params GetAllDayTypesParams) {
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
	d, count, err := c2.gssrv.GetDayTypes.Handler(c, seasonId, services.GetAllDto{
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
	response := make([]DayTypes, 0)
	links := server.GetLinksList(c, count, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
		Q:      params.Q,
		Sort:   params.Sort,
	})
	for _, dis := range d {
		response = append(response, DayTypesToResponse(dis))
	}
	c.JSON(http.StatusOK, DayTypesResponseList{
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
func (c2 Controller) GetDayType(c *gin.Context, dayTypeId string) {
	dt, err := c2.gssrv.GetDayTypeById.Handler(c, services.GetDayTypeByIdDto{
		Id: dayTypeId,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, DayTypesToResponse(dt))
}
func (c2 Controller) InsertDayType(c *gin.Context, seasonId string) {
	var req InsertDayTypeJSONRequestBody
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := middleware.GetAuthUser(c)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = c2.gssrv.InsertDayType.Handler(c, seasonId, seasons.DayTypes{
		Name:      req.Name,
		Month:     req.Month,
		IsFestive: req.IsFestive,
		CreatedBy: user.ID,
		CreatedAt: time.Now(),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": "created"})
}
func (c2 Controller) ModifyDayType(c *gin.Context, seasonsId string) {
	var req ModifyDayTypeJSONRequestBody
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetAuthUser(c)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	seasonUUID := uuid.Must(uuid.FromString(seasonsId))
	err = c2.gssrv.ModifyDayType.Handler(c, seasonUUID, seasons.DayTypes{
		Name:      req.Name,
		Month:     req.Month,
		IsFestive: req.IsFestive,
		UpdatedAt: time.Now(),
		UpdatedBy: user.ID,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "modified"})

}
func (c2 Controller) DeleteDayType(c *gin.Context, seasonId string) {
	err := c2.gssrv.DeleteDayType.Handler(c, services.DeleteSeasonByIdDto{
		Id: uuid.Must(uuid.FromString(seasonId)),
	})
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "deleted"})
}
