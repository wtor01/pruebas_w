package validations

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth/platform/middleware"
	validations2 "bitbucket.org/sercide/data-ingestion/internal/validations"
	validations3 "bitbucket.org/sercide/data-ingestion/internal/validations/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type Controller struct {
	svr *validations3.Services
}

func NewController(svr *validations3.Services) *Controller {
	return &Controller{svr: svr}
}

func (c Controller) ListValidationsMeasure(ctx *gin.Context, params ListValidationsMeasureParams) {

	result, count, err := c.svr.ListValidationMeasureService.Handle(ctx, validations3.ListValidationMeasureDto{
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

	ctx.JSON(http.StatusOK, ListValidationMeasures{
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
		Results: utils.MapSlice[validations2.ValidationMeasure, ValidationMeasure](result, validationMeasureToResponse),
	})

}

func (c Controller) CreateValidationsMeasure(ctx *gin.Context) {
	var req CreateValidationMeasure

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	description := ""

	if req.Description != nil {
		description = *req.Description
	}

	user, err := middleware.GetAuthUser(ctx)

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	v, err := c.svr.CreateValidationMeasureService.Handle(ctx, validations3.CreateValidationMeasureDto{
		Id:          uuid.NewV4().String(),
		UserID:      user.ID,
		Name:        req.Name,
		Action:      string(req.Action),
		Enabled:     req.Enabled,
		MeasureType: string(req.MeasureType),
		Type:        string(req.Type),
		Code:        req.Code,
		Message:     req.Message,
		Description: description,
		Params: validations2.Params{
			Type: string(req.Params.Type),
			Validations: utils.MapSlice[Validation, validations2.Validation](req.Params.Validations, func(item Validation) validations2.Validation {
				config := make(map[string]string)

				if item.Config != nil {
					config = item.Config.AdditionalProperties
				}
				return validations2.Validation{
					Id:   item.Id,
					Name: item.Name,
					Type: string(item.Type),
					Keys: utils.MapSlice[ValidationKeys, string](item.Keys, func(item ValidationKeys) string {
						return string(item)
					}),
					Required: item.Required,
					Config:   config,
				}
			}),
		},
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, validationMeasureToResponse(v))
}

func (c Controller) DeleteValidationsMeasure(ctx *gin.Context, validationId string) {
	err := c.svr.DeleteValidationMeasureByIdService.Handle(ctx, validations3.DeleteValidationMeasureByIdDto{
		ID: validationId,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c Controller) GetValidationsMeasure(ctx *gin.Context, validationId string) {
	result, err := c.svr.GetValidationMeasureByIdService.Handle(ctx, validations3.GetValidationMeasureByIdDto{
		ID: validationId,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, validationMeasureToResponse(result))

}

func (c Controller) UpdateValidationsMeasure(ctx *gin.Context, validationId string) {
	var req UpdateValidationsMeasureJSONRequestBody

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	description := ""

	if req.Description != nil {
		description = *req.Description
	}

	user, err := middleware.GetAuthUser(ctx)

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	v, err := c.svr.UpdateValidationMeasureService.Handle(ctx, validations3.UpdateValidationMeasureDto{
		Id:          validationId,
		UserID:      user.ID,
		Name:        req.Name,
		Action:      string(req.Action),
		Enabled:     req.Enabled,
		MeasureType: string(req.MeasureType),
		Type:        string(req.Type),
		Code:        req.Code,
		Message:     req.Message,
		Description: description,
		Params: validations2.Params{
			Type: string(req.Params.Type),
			Validations: utils.MapSlice[Validation, validations2.Validation](req.Params.Validations, func(item Validation) validations2.Validation {
				config := make(map[string]string)

				if item.Config != nil {
					config = item.Config.AdditionalProperties
				}
				return validations2.Validation{
					Id:   item.Id,
					Type: string(item.Type),
					Keys: utils.MapSlice[ValidationKeys, string](item.Keys, func(item ValidationKeys) string {
						return string(item)
					}),
					Required: item.Required,
					Config:   config,
					Name:     item.Name,
				}
			}),
		},
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, validationMeasureToResponse(v))
}

func (c Controller) ListValidationsMeasureConfig(ctx *gin.Context, distributorId string, params ListValidationsMeasureConfigParams) {

	t := ""

	if params.Type != nil {
		t = string(*params.Type)
	}

	result, err, _ := c.svr.ListValidationMeasureConfigService.Handle(ctx, validations3.ListValidationMeasureConfigDto{
		Type:          t,
		DistributorID: distributorId,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paginate := server.GetPaginate(ctx, server.Query{
		Limit: len(result),
	}, result, len(result))

	ctx.JSON(http.StatusOK, ListValidationMeasureConfig{
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
		Results: utils.MapSlice(result, validationMeasureConfigToResponse),
	})
}

func (c Controller) CreateValidationsMeasureConfig(ctx *gin.Context, distributorId string) {
	var req ValidationMeasureConfigCreate
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetAuthUser(ctx)

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	extraConfig := make([]ExtraConfig, 0)

	if req.ExtraConfig != nil {
		extraConfig = *req.ExtraConfig
	}

	config, err := c.svr.CreateValidationMeasureConfigService.Handle(ctx, validations3.CreateValidationMeasureConfigDto{
		ID:                  uuid.NewV4().String(),
		UserID:              user.ID,
		ValidationMeasureID: req.ValidationMeasureId,
		DistributorID:       distributorId,
		Action:              string(req.Action),
		Enabled:             req.Enabled,
		Params: utils.MapSlice(extraConfig, func(item ExtraConfig) validations2.Config {

			return validations2.Config{
				ID:    item.Id,
				Extra: item.AdditionalProperties,
			}
		}),
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, validationMeasureConfigToResponse(config))
}

func (c Controller) DeleteValidationsMeasureConfig(ctx *gin.Context, distributorId string, configurationId string) {
	err := c.svr.DeleteValidationMeasureConfigService.Handle(ctx, validations3.DeleteValidationMeasureConfigDto{
		ID:            configurationId,
		DistributorID: distributorId,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c Controller) GetValidationsMeasureConfig(ctx *gin.Context, distributorId string, configurationId string) {

	config, err := c.svr.GetValidationMeasureConfigService.Handle(ctx, validations3.GetValidationMeasureConfigDto{
		ID:            configurationId,
		DistributorID: distributorId,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, validationMeasureConfigToResponse(config))
}

func (c Controller) UpdateValidationsMeasureConfig(ctx *gin.Context, distributorId string, configurationId string) {
	var req ValidationMeasureConfigCreate
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetAuthUser(ctx)

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	extraConfig := make([]ExtraConfig, 0)

	if req.ExtraConfig != nil {
		extraConfig = *req.ExtraConfig
	}

	config, err := c.svr.CreateValidationMeasureConfigService.Handle(ctx, validations3.CreateValidationMeasureConfigDto{
		ID:                  configurationId,
		UserID:              user.ID,
		ValidationMeasureID: req.ValidationMeasureId,
		DistributorID:       distributorId,
		Action:              string(req.Action),
		Enabled:             req.Enabled,
		Params: utils.MapSlice(extraConfig, func(item ExtraConfig) validations2.Config {

			return validations2.Config{
				ID:    item.Id,
				Extra: item.AdditionalProperties,
			}
		}),
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, validationMeasureConfigToResponse(config))
}

func (c Controller) PutMeasurementValidation(ctx *gin.Context, measureType PutMeasurementValidationParamsMeasureType) {
	var req PutMeasurementValidationJSONRequestBody

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := c.svr.PutMeasureValidation.Handle(ctx, validations3.PutValidationMeasureDTO{
		MeasureType:      string(measureType),
		Status:           string(req.Status),
		ID:               req.ID,
		InvalidationCode: req.InvalidationCode,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
