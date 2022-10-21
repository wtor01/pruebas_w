// Package stats provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package stats

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// get process measure statistics by cups
	// (GET /dashboard/statistics/cups)
	GetProcessMeasureStatisticsByCups(ctx *gin.Context, params GetProcessMeasureStatisticsByCupsParams)
	// get process measure statistics
	// (GET /dashboard/statistics/global)
	GetProcessMeasureStatistics(ctx *gin.Context, params GetProcessMeasureStatisticsParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []gin.HandlerFunc
}

// GetProcessMeasureStatisticsByCups operation middleware
func (siw *ServerInterfaceWrapper) GetProcessMeasureStatisticsByCups(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetProcessMeasureStatisticsByCupsParams

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

	// ------------- Required query parameter "month" -------------
	if paramValue := c.Query("month"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument month is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "month", c.Request.URL.Query(), &params.Month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter month: %s", err)})
		return
	}

	// ------------- Required query parameter "year" -------------
	if paramValue := c.Query("year"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument year is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "year", c.Request.URL.Query(), &params.Year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter year: %s", err)})
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

	siw.Handler.GetProcessMeasureStatisticsByCups(c, params)
}

// GetProcessMeasureStatistics operation middleware
func (siw *ServerInterfaceWrapper) GetProcessMeasureStatistics(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetProcessMeasureStatisticsParams

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

	// ------------- Required query parameter "month" -------------
	if paramValue := c.Query("month"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument month is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "month", c.Request.URL.Query(), &params.Month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter month: %s", err)})
		return
	}

	// ------------- Required query parameter "year" -------------
	if paramValue := c.Query("year"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument year is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "year", c.Request.URL.Query(), &params.Year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter year: %s", err)})
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

	siw.Handler.GetProcessMeasureStatistics(c, params)
}

// RegisterHandlerGetProcessMeasureStatisticsByCups creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetProcessMeasureStatisticsByCups(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetProcessMeasureStatisticsByCups(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetProcessMeasureStatisticsByCups creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetProcessMeasureStatisticsByCups(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetProcessMeasureStatisticsByCups)

	router.GET("/dashboard/statistics/cups", m...)

	return router
}

// RegisterHandlerGetProcessMeasureStatistics creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetProcessMeasureStatistics(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetProcessMeasureStatistics(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetProcessMeasureStatistics creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetProcessMeasureStatistics(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetProcessMeasureStatistics)

	router.GET("/dashboard/statistics/global", m...)

	return router
}
