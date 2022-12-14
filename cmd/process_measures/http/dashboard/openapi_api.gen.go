// Package dashboard provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package dashboard

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// process measure dashboard
	// (GET /process-measures/dashboard/cups)
	GetProcessMeasureDashboardList(ctx *gin.Context, params GetProcessMeasureDashboardListParams)
	// Search service point  for process measures
	// (GET /process-measures/dashboard/curve-process-measures)
	GetCurveProcessServicePoint(ctx *gin.Context, params GetCurveProcessServicePointParams)
	// process measure dashboard
	// (GET /process-measures/dashboard/measures)
	GetProcessMeasureDashboard(ctx *gin.Context, params GetProcessMeasureDashboardParams)
	// Search service point  for process measures
	// (GET /process-measures/dashboard/service-point-process-measures)
	GetDashboardProcessServicePoint(ctx *gin.Context, params GetDashboardProcessServicePointParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []gin.HandlerFunc
}

// GetProcessMeasureDashboardList operation middleware
func (siw *ServerInterfaceWrapper) GetProcessMeasureDashboardList(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetProcessMeasureDashboardListParams

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

	// ------------- Optional query parameter "offset" -------------
	if paramValue := c.Query("offset"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "offset", c.Request.URL.Query(), &params.Offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter offset: %s", err)})
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

	// ------------- Required query parameter "type" -------------
	if paramValue := c.Query("type"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument type is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "type", c.Request.URL.Query(), &params.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter type: %s", err)})
		return
	}

	siw.Handler.GetProcessMeasureDashboardList(c, params)
}

// GetCurveProcessServicePoint operation middleware
func (siw *ServerInterfaceWrapper) GetCurveProcessServicePoint(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetCurveProcessServicePointParams

	// ------------- Required query parameter "distributor" -------------
	if paramValue := c.Query("distributor"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument distributor is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "distributor", c.Request.URL.Query(), &params.Distributor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter distributor: %s", err)})
		return
	}

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

	// ------------- Required query parameter "curve_type" -------------
	if paramValue := c.Query("curve_type"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument curve_type is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "curve_type", c.Request.URL.Query(), &params.CurveType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter curve_type: %s", err)})
		return
	}

	siw.Handler.GetCurveProcessServicePoint(c, params)
}

// GetProcessMeasureDashboard operation middleware
func (siw *ServerInterfaceWrapper) GetProcessMeasureDashboard(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetProcessMeasureDashboardParams

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

	siw.Handler.GetProcessMeasureDashboard(c, params)
}

// GetDashboardProcessServicePoint operation middleware
func (siw *ServerInterfaceWrapper) GetDashboardProcessServicePoint(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetDashboardProcessServicePointParams

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

	// ------------- Required query parameter "distributor" -------------
	if paramValue := c.Query("distributor"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument distributor is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "distributor", c.Request.URL.Query(), &params.Distributor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter distributor: %s", err)})
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

	siw.Handler.GetDashboardProcessServicePoint(c, params)
}

// RegisterHandlerGetProcessMeasureDashboardList creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetProcessMeasureDashboardList(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetProcessMeasureDashboardList(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetProcessMeasureDashboardList creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetProcessMeasureDashboardList(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetProcessMeasureDashboardList)

	router.GET("/process-measures/dashboard/cups", m...)

	return router
}

// RegisterHandlerGetCurveProcessServicePoint creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetCurveProcessServicePoint(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetCurveProcessServicePoint(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetCurveProcessServicePoint creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetCurveProcessServicePoint(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetCurveProcessServicePoint)

	router.GET("/process-measures/dashboard/curve-process-measures", m...)

	return router
}

// RegisterHandlerGetProcessMeasureDashboard creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetProcessMeasureDashboard(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetProcessMeasureDashboard(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetProcessMeasureDashboard creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetProcessMeasureDashboard(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetProcessMeasureDashboard)

	router.GET("/process-measures/dashboard/measures", m...)

	return router
}

// RegisterHandlerGetDashboardProcessServicePoint creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetDashboardProcessServicePoint(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetDashboardProcessServicePoint(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetDashboardProcessServicePoint creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetDashboardProcessServicePoint(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetDashboardProcessServicePoint)

	router.GET("/process-measures/dashboard/service-point-process-measures", m...)

	return router
}
