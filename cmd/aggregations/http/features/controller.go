package features

import (
	"bitbucket.org/sercide/data-ingestion/internal/aggregations/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type Controller struct {
	srv *services.FeaturesServices
}

func (c Controller) GetAllAggregationFeaturesAvailable(ctx *gin.Context, params GetAllAggregationFeaturesAvailableParams) {
	result, count, err := c.srv.ListFeatures.Handler(ctx.Request.Context(), services.ListFeaturesDto{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	paginate := server.GetPaginate(ctx, server.Query{
		Limit:  params.Limit,
		Offset: params.Offset,
	}, result, count)
	ctx.JSON(http.StatusOK, AggregationFeatures{
		Links: struct {
			Next *string `json:"next,omitempty"`
			Prev *string `json:"prev,omitempty"`
			Self string  `json:"self"`
		}{
			Next: paginate.Links.Next,
			Prev: paginate.Links.Prev,
			Self: paginate.Links.Self,
		},
		Count:   paginate.Count,
		Limit:   paginate.Limit,
		Offset:  paginate.Offset,
		Results: utils.MapSlice(result, featuresToResponse),
		Size:    paginate.Size,
	})

}

func (c Controller) CreateAggregationFeaturesAvailable(ctx *gin.Context) {
	var req AggregationFeaturesBase
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	id, err := uuid.NewUUID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	s, err := c.srv.CreateFeatures.Handler(ctx.Request.Context(), services.CreateFeaturesDTO{
		ID:    id.String(),
		Name:  req.Name,
		Field: req.Field,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, AggregationFeature{
		Id: s.ID,
		AggregationFeaturesBase: AggregationFeaturesBase{
			Field: s.Field,
			Name:  s.Name,
		},
	})
}

func (c Controller) DeleteAggregationFeaturesAvailable(ctx *gin.Context, featureId string) {
	err := c.srv.DeleteFeatures.Handler(ctx.Request.Context(), services.DeleteFeaturesDTO{ID: featureId})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}

func (c Controller) GetAggregationFeaturesAvailable(ctx *gin.Context, featureId string) {
	result, err := c.srv.GetFeatures.Handler(ctx.Request.Context(), services.GetFeaturesDTO{
		ID: featureId,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, featuresToResponse(result))
}

func (c Controller) UpdateAggregationFeaturesAvailable(ctx *gin.Context, featureId string) {
	var req UpdateAggregationFeature

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	features, err := c.srv.UpdateFeatures.Handler(ctx.Request.Context(), services.UpdateFeaturesDTO{
		ID:    featureId,
		Field: req.Field,
		Name:  req.Name,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, featuresToResponse(features))
}
