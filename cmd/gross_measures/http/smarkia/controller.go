package smarkia

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	recoverSmarkiaMeasuresService *services.RecoverSmarkiaMeasures
}

func NewController(recoverSmarkiaMeasuresService *services.RecoverSmarkiaMeasures) *Controller {
	return &Controller{recoverSmarkiaMeasuresService: recoverSmarkiaMeasuresService}
}

func (c Controller) RecoverSmarkiaMeasures(ctx *gin.Context) {
	var req RecoverSmarkiaMeasures
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.recoverSmarkiaMeasuresService.Handler(ctx.Request.Context(), services.RecoverSmarkiaMeasuresDTO{
		CUPS:          req.Cups,
		DistributorID: req.DistributorId,
		Date:          req.Date.Time,
		ProcessName:   smarkia.ProcessName(req.ProcessName),
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}
