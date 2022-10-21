package closures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type Controller struct {
	closureService  *services.ClosureServices
	getClosure      *services.GetClosure
	executeServices *services.ExecuteServices
	getResumes      *services.GetResume
}

func NewController(closureService *services.ClosureServices, executeServices *services.ExecuteServices, getResumes *services.GetResume) *Controller {
	return &Controller{closureService: closureService, executeServices: executeServices, getResumes: getResumes}
}
func (c Controller) GetProcessMeasuresResumeByCups(ctx *gin.Context, params GetProcessMeasuresResumeByCupsParams) {

	resumes, err := c.getResumes.Handler(ctx, services.ListResumeDto{
		DistributorId: params.DistributorId,
		Cups:          params.Cups,
		StartDate:     params.StartDate.Time,
		EndDate:       params.EndDate.Time,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}

	ctx.JSON(http.StatusOK, resumeToResponse(resumes))
}

func (c Controller) ExecuteProcessMeasuresServices(ctx *gin.Context) {
	var req ExecuteProcessMeasure
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := c.executeServices.Handler(ctx.Request.Context(), services.DtoServiceExecute{
		Cups:          req.Cups,
		DistributorId: req.DistributorId,
		StartDate:     req.StartDate.Time,
		EndDate:       req.EndDate.Time,
		ReadingType:   measures.ReadingType(req.Type),
	})
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusOK)
}
func (c Controller) GetBillingClosures(ctx *gin.Context, distributorId string, params GetBillingClosuresParams) {
	var response process_measures.ProcessedMonthlyClosure
	moment := string(*params.Moment)
	var err error
	if params.Id != nil {
		response, err = c.closureService.GetClosure.Handler(ctx, services.ListClosureDto{
			DistributorId: distributorId,
			Id:            *params.Id,
			Moment:        moment,
		})
	} else {
		response, err = c.closureService.GetClosure.Handler(ctx, services.ListClosureDto{
			DistributorId: distributorId,
			Cups:          *params.Cups,
			StartDate:     params.StartDate.Time,
			EndDate:       params.EndDate.Time,
			Moment:        moment,
		})

	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, closureToResponse(response))
}

func (c Controller) CreateBillingClosures(ctx *gin.Context, distributorId string) {
	var req MonthlyClosure

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := uuid.NewUUID()
	result := closureToUpdateOrCreate(req)
	result.Id = id.String()
	result.DistributorID = distributorId
	err := c.closureService.CreateClosure.Handler(ctx, services.CreateClosureDto{
		Monthly: result,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": "created"})

}

func (c Controller) UpdateBillingClosures(ctx *gin.Context, distributorId string, id string) {

	var req MonthlyClosure

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := closureToUpdateOrCreate(req)
	result.Id = id
	result.DistributorID = distributorId

	err := c.closureService.UpdateClosure.Handler(ctx, services.UpdateClosureDto{
		Monthly: result,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusCreated, gin.H{"success": "modified"})

}
