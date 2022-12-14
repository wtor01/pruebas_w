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
	// measure dashboard
	// (GET /dashboard/measures)
	GetMeasureDashboard(ctx *gin.Context, params GetMeasureDashboardParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []gin.HandlerFunc
}

// GetMeasureDashboard operation middleware
func (siw *ServerInterfaceWrapper) GetMeasureDashboard(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetMeasureDashboardParams

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

	siw.Handler.GetMeasureDashboard(c, params)
}

// RegisterHandlerGetMeasureDashboard creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetMeasureDashboard(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetMeasureDashboard(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetMeasureDashboard creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetMeasureDashboard(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetMeasureDashboard)

	router.GET("/dashboard/measures", m...)

	return router
}
