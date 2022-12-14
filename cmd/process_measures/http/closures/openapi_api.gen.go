// Package closures provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package closures

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get billing closures
	// (GET /distributor/{distributor_id}/billing_closure)
	GetBillingClosures(ctx *gin.Context, distributorId string, params GetBillingClosuresParams)
	// Create billing closures
	// (POST /distributor/{distributor_id}/billing_closure)
	CreateBillingClosures(ctx *gin.Context, distributorId string)
	// Update billing closures
	// (PUT /distributor/{distributor_id}/billing_closure/{id})
	UpdateBillingClosures(ctx *gin.Context, distributorId string, id string)
	// Get process measures resume by cups
	// (GET /process-measures/create-close/resume)
	GetProcessMeasuresResumeByCups(ctx *gin.Context, params GetProcessMeasuresResumeByCupsParams)
	// Create measures services
	// (POST /process-measures/execute)
	ExecuteProcessMeasuresServices(ctx *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []gin.HandlerFunc
}

// GetBillingClosures operation middleware
func (siw *ServerInterfaceWrapper) GetBillingClosures(c *gin.Context) {

	var err error

	// ------------- Path parameter "distributor_id" -------------
	var distributorId string

	err = runtime.BindStyledParameter("simple", false, "distributor_id", c.Param("distributor_id"), &distributorId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter distributor_id: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetBillingClosuresParams

	// ------------- Optional query parameter "id" -------------
	if paramValue := c.Query("id"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "id", c.Request.URL.Query(), &params.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter id: %s", err)})
		return
	}

	// ------------- Optional query parameter "moment" -------------
	if paramValue := c.Query("moment"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "moment", c.Request.URL.Query(), &params.Moment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter moment: %s", err)})
		return
	}

	// ------------- Optional query parameter "cups" -------------
	if paramValue := c.Query("cups"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "cups", c.Request.URL.Query(), &params.Cups)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter cups: %s", err)})
		return
	}

	// ------------- Optional query parameter "start_date" -------------
	if paramValue := c.Query("start_date"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "start_date", c.Request.URL.Query(), &params.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter start_date: %s", err)})
		return
	}

	// ------------- Optional query parameter "end_date" -------------
	if paramValue := c.Query("end_date"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "end_date", c.Request.URL.Query(), &params.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter end_date: %s", err)})
		return
	}

	siw.Handler.GetBillingClosures(c, distributorId, params)
}

// CreateBillingClosures operation middleware
func (siw *ServerInterfaceWrapper) CreateBillingClosures(c *gin.Context) {

	var err error

	// ------------- Path parameter "distributor_id" -------------
	var distributorId string

	err = runtime.BindStyledParameter("simple", false, "distributor_id", c.Param("distributor_id"), &distributorId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter distributor_id: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.CreateBillingClosures(c, distributorId)
}

// UpdateBillingClosures operation middleware
func (siw *ServerInterfaceWrapper) UpdateBillingClosures(c *gin.Context) {

	var err error

	// ------------- Path parameter "distributor_id" -------------
	var distributorId string

	err = runtime.BindStyledParameter("simple", false, "distributor_id", c.Param("distributor_id"), &distributorId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter distributor_id: %s", err)})
		return
	}

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter id: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.UpdateBillingClosures(c, distributorId, id)
}

// GetProcessMeasuresResumeByCups operation middleware
func (siw *ServerInterfaceWrapper) GetProcessMeasuresResumeByCups(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetProcessMeasuresResumeByCupsParams

	// ------------- Required query parameter "cups" -------------
	if paramValue := c.Query("cups"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument cups is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "cups", c.Request.URL.Query(), &params.Cups)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter cups: %s", err)})
		return
	}

	// ------------- Required query parameter "distributor_id" -------------
	if paramValue := c.Query("distributor_id"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument distributor_id is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "distributor_id", c.Request.URL.Query(), &params.DistributorId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter distributor_id: %s", err)})
		return
	}

	// ------------- Required query parameter "start_date" -------------
	if paramValue := c.Query("start_date"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument start_date is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "start_date", c.Request.URL.Query(), &params.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter start_date: %s", err)})
		return
	}

	// ------------- Required query parameter "end_date" -------------
	if paramValue := c.Query("end_date"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument end_date is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "end_date", c.Request.URL.Query(), &params.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter end_date: %s", err)})
		return
	}

	siw.Handler.GetProcessMeasuresResumeByCups(c, params)
}

// ExecuteProcessMeasuresServices operation middleware
func (siw *ServerInterfaceWrapper) ExecuteProcessMeasuresServices(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.ExecuteProcessMeasuresServices(c)
}

// RegisterHandlerGetBillingClosures creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetBillingClosures(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetBillingClosures(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetBillingClosures creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetBillingClosures(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetBillingClosures)

	router.GET("/distributor/:distributor_id/billing_closure", m...)

	return router
}

// RegisterHandlerCreateBillingClosures creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerCreateBillingClosures(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresCreateBillingClosures(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresCreateBillingClosures creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresCreateBillingClosures(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.CreateBillingClosures)

	router.POST("/distributor/:distributor_id/billing_closure", m...)

	return router
}

// RegisterHandlerUpdateBillingClosures creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerUpdateBillingClosures(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresUpdateBillingClosures(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresUpdateBillingClosures creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresUpdateBillingClosures(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.UpdateBillingClosures)

	router.PUT("/distributor/:distributor_id/billing_closure/:id", m...)

	return router
}

// RegisterHandlerGetProcessMeasuresResumeByCups creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetProcessMeasuresResumeByCups(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetProcessMeasuresResumeByCups(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetProcessMeasuresResumeByCups creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetProcessMeasuresResumeByCups(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetProcessMeasuresResumeByCups)

	router.GET("/process-measures/create-close/resume", m...)

	return router
}

// RegisterHandlerExecuteProcessMeasuresServices creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerExecuteProcessMeasuresServices(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresExecuteProcessMeasuresServices(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresExecuteProcessMeasuresServices creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresExecuteProcessMeasuresServices(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.ExecuteProcessMeasuresServices)

	router.POST("/process-measures/execute", m...)

	return router
}
