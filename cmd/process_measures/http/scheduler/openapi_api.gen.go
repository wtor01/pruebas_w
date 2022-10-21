// Package scheduler provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package scheduler

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List process measures scheduler
	// (GET /process-measures/scheduler)
	ListProcessMeasuresScheduler(ctx *gin.Context, params ListProcessMeasuresSchedulerParams)
	// Create process measures scheduler
	// (POST /process-measures/scheduler)
	CreateProcessMeasuresScheduler(ctx *gin.Context)
	// Delete process measures scheduler
	// (DELETE /process-measures/scheduler/{id})
	DeleteProcessMeasuresScheduler(ctx *gin.Context, id string)
	// Get process measure scheduler by id
	// (GET /process-measures/scheduler/{id})
	GetProcessMeasureSchedulerById(ctx *gin.Context, id string)
	// Patch process measures scheduler
	// (PATCH /process-measures/scheduler/{id})
	PatchProcessMeasuresScheduler(ctx *gin.Context, id string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []gin.HandlerFunc
}

// ListProcessMeasuresScheduler operation middleware
func (siw *ServerInterfaceWrapper) ListProcessMeasuresScheduler(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListProcessMeasuresSchedulerParams

	// ------------- Optional query parameter "offset" -------------
	if paramValue := c.Query("offset"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "offset", c.Request.URL.Query(), &params.Offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter offset: %s", err)})
		return
	}

	// ------------- Required query parameter "limit" -------------
	if paramValue := c.Query("limit"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument limit is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "limit", c.Request.URL.Query(), &params.Limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter limit: %s", err)})
		return
	}

	siw.Handler.ListProcessMeasuresScheduler(c, params)
}

// CreateProcessMeasuresScheduler operation middleware
func (siw *ServerInterfaceWrapper) CreateProcessMeasuresScheduler(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.CreateProcessMeasuresScheduler(c)
}

// DeleteProcessMeasuresScheduler operation middleware
func (siw *ServerInterfaceWrapper) DeleteProcessMeasuresScheduler(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter id: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.DeleteProcessMeasuresScheduler(c, id)
}

// GetProcessMeasureSchedulerById operation middleware
func (siw *ServerInterfaceWrapper) GetProcessMeasureSchedulerById(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter id: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.GetProcessMeasureSchedulerById(c, id)
}

// PatchProcessMeasuresScheduler operation middleware
func (siw *ServerInterfaceWrapper) PatchProcessMeasuresScheduler(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter id: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.PatchProcessMeasuresScheduler(c, id)
}

// RegisterHandlerListProcessMeasuresScheduler creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerListProcessMeasuresScheduler(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresListProcessMeasuresScheduler(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresListProcessMeasuresScheduler creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresListProcessMeasuresScheduler(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.ListProcessMeasuresScheduler)

	router.GET("/process-measures/scheduler", m...)

	return router
}

// RegisterHandlerCreateProcessMeasuresScheduler creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerCreateProcessMeasuresScheduler(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresCreateProcessMeasuresScheduler(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresCreateProcessMeasuresScheduler creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresCreateProcessMeasuresScheduler(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.CreateProcessMeasuresScheduler)

	router.POST("/process-measures/scheduler", m...)

	return router
}

// RegisterHandlerDeleteProcessMeasuresScheduler creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerDeleteProcessMeasuresScheduler(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresDeleteProcessMeasuresScheduler(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresDeleteProcessMeasuresScheduler creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresDeleteProcessMeasuresScheduler(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.DeleteProcessMeasuresScheduler)

	router.DELETE("/process-measures/scheduler/:id", m...)

	return router
}

// RegisterHandlerGetProcessMeasureSchedulerById creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetProcessMeasureSchedulerById(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetProcessMeasureSchedulerById(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetProcessMeasureSchedulerById creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetProcessMeasureSchedulerById(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetProcessMeasureSchedulerById)

	router.GET("/process-measures/scheduler/:id", m...)

	return router
}

// RegisterHandlerPatchProcessMeasuresScheduler creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerPatchProcessMeasuresScheduler(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresPatchProcessMeasuresScheduler(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresPatchProcessMeasuresScheduler creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresPatchProcessMeasuresScheduler(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.PatchProcessMeasuresScheduler)

	router.PATCH("/process-measures/scheduler/:id", m...)

	return router
}
