// Package config provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package config

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all aggregations config
	// (GET /admin/aggregations/config/)
	GetAllAggregationsConfig(ctx *gin.Context, params GetAllAggregationsConfigParams)
	// Create aggregation config
	// (POST /admin/aggregations/config/)
	CreateAggregationConfig(ctx *gin.Context)
	// Delete aggregation config
	// (DELETE /admin/aggregations/config/{aggregation_config_id})
	DeleteAggregationConfig(ctx *gin.Context, aggregationConfigId string)
	// Get aggregation config
	// (GET /admin/aggregations/config/{aggregation_config_id})
	GetAggregationConfig(ctx *gin.Context, aggregationConfigId string)
	// Update aggregation config
	// (PUT /admin/aggregations/config/{aggregation_config_id})
	UpdateAggregationConfig(ctx *gin.Context, aggregationConfigId string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []gin.HandlerFunc
}

// GetAllAggregationsConfig operation middleware
func (siw *ServerInterfaceWrapper) GetAllAggregationsConfig(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAllAggregationsConfigParams

	// ------------- Optional query parameter "q" -------------
	if paramValue := c.Query("q"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "q", c.Request.URL.Query(), &params.Q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter q: %s", err)})
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

	siw.Handler.GetAllAggregationsConfig(c, params)
}

// CreateAggregationConfig operation middleware
func (siw *ServerInterfaceWrapper) CreateAggregationConfig(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.CreateAggregationConfig(c)
}

// DeleteAggregationConfig operation middleware
func (siw *ServerInterfaceWrapper) DeleteAggregationConfig(c *gin.Context) {

	var err error

	// ------------- Path parameter "aggregation_config_id" -------------
	var aggregationConfigId string

	err = runtime.BindStyledParameter("simple", false, "aggregation_config_id", c.Param("aggregation_config_id"), &aggregationConfigId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter aggregation_config_id: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.DeleteAggregationConfig(c, aggregationConfigId)
}

// GetAggregationConfig operation middleware
func (siw *ServerInterfaceWrapper) GetAggregationConfig(c *gin.Context) {

	var err error

	// ------------- Path parameter "aggregation_config_id" -------------
	var aggregationConfigId string

	err = runtime.BindStyledParameter("simple", false, "aggregation_config_id", c.Param("aggregation_config_id"), &aggregationConfigId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter aggregation_config_id: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.GetAggregationConfig(c, aggregationConfigId)
}

// UpdateAggregationConfig operation middleware
func (siw *ServerInterfaceWrapper) UpdateAggregationConfig(c *gin.Context) {

	var err error

	// ------------- Path parameter "aggregation_config_id" -------------
	var aggregationConfigId string

	err = runtime.BindStyledParameter("simple", false, "aggregation_config_id", c.Param("aggregation_config_id"), &aggregationConfigId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter aggregation_config_id: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.UpdateAggregationConfig(c, aggregationConfigId)
}

// RegisterHandlerGetAllAggregationsConfig creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetAllAggregationsConfig(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetAllAggregationsConfig(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetAllAggregationsConfig creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetAllAggregationsConfig(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetAllAggregationsConfig)

	router.GET("/admin/aggregations/config/", m...)

	return router
}

// RegisterHandlerCreateAggregationConfig creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerCreateAggregationConfig(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresCreateAggregationConfig(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresCreateAggregationConfig creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresCreateAggregationConfig(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.CreateAggregationConfig)

	router.POST("/admin/aggregations/config/", m...)

	return router
}

// RegisterHandlerDeleteAggregationConfig creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerDeleteAggregationConfig(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresDeleteAggregationConfig(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresDeleteAggregationConfig creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresDeleteAggregationConfig(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.DeleteAggregationConfig)

	router.DELETE("/admin/aggregations/config/:aggregation_config_id", m...)

	return router
}

// RegisterHandlerGetAggregationConfig creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetAggregationConfig(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetAggregationConfig(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetAggregationConfig creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetAggregationConfig(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetAggregationConfig)

	router.GET("/admin/aggregations/config/:aggregation_config_id", m...)

	return router
}

// RegisterHandlerUpdateAggregationConfig creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerUpdateAggregationConfig(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresUpdateAggregationConfig(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresUpdateAggregationConfig creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresUpdateAggregationConfig(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.UpdateAggregationConfig)

	router.PUT("/admin/aggregations/config/:aggregation_config_id", m...)

	return router
}
