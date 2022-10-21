package scheduler

import (
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services"
	"bitbucket.org/sercide/data-ingestion/pkg/server"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type Controller struct {
	srv *services.SchedulerServices
}

func (c Controller) GetProcessMeasureSchedulerById(ctx *gin.Context, id string) {
	result, err := c.srv.GetSchedulerById.Handler(ctx, services.GetSchedulerDTO{
		ID: id,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, schedulerToResponse(result))
}

func (c Controller) PatchProcessMeasuresScheduler(ctx *gin.Context, id string) {
	var req PatchProcessMeasureScheduler

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	distributorID := ""

	if req.DistributorId != nil {
		distributorID = *req.DistributorId
	}

	s, err := c.srv.UpdateScheduler.Handler(ctx, services.UpdateSchedulerDTO{
		ID:            id,
		DistributorId: distributorID,
		ServiceType:   string(req.ServiceType),
		PointType:     string(req.PointType),
		MeterType: utils.MapSlice(req.MeterType, func(item ProcessMeasureSchedulerUpdatableMeterType) string {
			return string(item)
		}),
		ReadingType: string(req.ReadingType),
		Scheduler:   req.Scheduler,
		Description: req.Description,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, s)
}

func (c Controller) ListProcessMeasuresScheduler(ctx *gin.Context, params ListProcessMeasuresSchedulerParams) {
	result, count, err := c.srv.ListScheduler.Handler(ctx, services.ListSchedulerDto{
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

	ctx.JSON(http.StatusOK, ListProcessMeasuresScheduler{
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
		Results: utils.MapSlice(result, schedulerToResponse),
	})
}

func (c Controller) DeleteProcessMeasuresScheduler(ctx *gin.Context, id string) {
	err := c.srv.DeleteScheduler.Handler(ctx, services.DeleteSchedulerDTO{
		ID: id,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c Controller) CreateProcessMeasuresScheduler(ctx *gin.Context) {
	var req ProcessMeasureSchedulerBase

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.NewUUID()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	distributorID := ""

	if req.DistributorId != nil {
		distributorID = *req.DistributorId
	}

	s, err := c.srv.CreateScheduler.Handler(ctx, services.CreateSchedulerDTO{
		ID:            id.String(),
		DistributorId: distributorID,
		ServiceType:   string(req.ServiceType),
		PointType:     string(req.PointType),
		MeterType: utils.MapSlice(req.MeterType, func(item ProcessMeasureSchedulerUpdatableMeterType) string {
			return string(item)
		}),
		ReadingType: string(req.ReadingType),
		Scheduler:   req.Scheduler,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, s)
}
