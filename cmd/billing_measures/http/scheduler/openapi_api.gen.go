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
	// List billing measures scheduler
	// (GET /billing-measures/scheduler)
	ListBillingMeasuresScheduler(ctx *gin.Context, params ListBillingMeasuresSchedulerParams)
	// Create billing measures scheduler
	// (POST /billing-measures/scheduler)
	CreateBillingMeasuresScheduler(ctx *gin.Context)
	// Delete billing measures scheduler
	// (DELETE /billing-measures/scheduler/{id})
	DeleteBillingMeasuresScheduler(ctx *gin.Context, id string)
	// Get billing measures scheduler by id
	// (GET /billing-measures/scheduler/{id})
	GetBillingMeasuresSchedulerById(ctx *gin.Context, id string)
	// Patch billing measures scheduler
	// (PATCH /billing-measures/scheduler/{id})
	PatchBillingMeasuresScheduler(ctx *gin.Context, id string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []gin.HandlerFunc
}

// ListBillingMeasuresScheduler operation middleware
func (siw *ServerInterfaceWrapper) ListBillingMeasuresScheduler(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListBillingMeasuresSchedulerParams

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

	siw.Handler.ListBillingMeasuresScheduler(c, params)
}

// CreateBillingMeasuresScheduler operation middleware
func (siw *ServerInterfaceWrapper) CreateBillingMeasuresScheduler(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.CreateBillingMeasuresScheduler(c)
}

// DeleteBillingMeasuresScheduler operation middleware
func (siw *ServerInterfaceWrapper) DeleteBillingMeasuresScheduler(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter id: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.DeleteBillingMeasuresScheduler(c, id)
}

// GetBillingMeasuresSchedulerById operation middleware
func (siw *ServerInterfaceWrapper) GetBillingMeasuresSchedulerById(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter id: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.GetBillingMeasuresSchedulerById(c, id)
}

// PatchBillingMeasuresScheduler operation middleware
func (siw *ServerInterfaceWrapper) PatchBillingMeasuresScheduler(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter id: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.PatchBillingMeasuresScheduler(c, id)
}

// RegisterHandlerListBillingMeasuresScheduler creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerListBillingMeasuresScheduler(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresListBillingMeasuresScheduler(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresListBillingMeasuresScheduler creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresListBillingMeasuresScheduler(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.ListBillingMeasuresScheduler)

	router.GET("/billing-measures/scheduler", m...)

	return router
}

// RegisterHandlerCreateBillingMeasuresScheduler creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerCreateBillingMeasuresScheduler(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresCreateBillingMeasuresScheduler(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresCreateBillingMeasuresScheduler creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresCreateBillingMeasuresScheduler(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.CreateBillingMeasuresScheduler)

	router.POST("/billing-measures/scheduler", m...)

	return router
}

// RegisterHandlerDeleteBillingMeasuresScheduler creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerDeleteBillingMeasuresScheduler(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresDeleteBillingMeasuresScheduler(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresDeleteBillingMeasuresScheduler creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresDeleteBillingMeasuresScheduler(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.DeleteBillingMeasuresScheduler)

	router.DELETE("/billing-measures/scheduler/:id", m...)

	return router
}

// RegisterHandlerGetBillingMeasuresSchedulerById creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetBillingMeasuresSchedulerById(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetBillingMeasuresSchedulerById(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetBillingMeasuresSchedulerById creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetBillingMeasuresSchedulerById(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetBillingMeasuresSchedulerById)

	router.GET("/billing-measures/scheduler/:id", m...)

	return router
}

// RegisterHandlerPatchBillingMeasuresScheduler creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerPatchBillingMeasuresScheduler(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresPatchBillingMeasuresScheduler(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresPatchBillingMeasuresScheduler creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresPatchBillingMeasuresScheduler(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.PatchBillingMeasuresScheduler)

	router.PATCH("/billing-measures/scheduler/:id", m...)

	return router
}
