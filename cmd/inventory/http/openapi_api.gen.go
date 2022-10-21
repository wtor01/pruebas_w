// Package http provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package http

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all distributors
	// (GET /distributors)
	ListDistributors(ctx *gin.Context, params ListDistributorsParams)
	// Get by smarkia id
	// (GET /distributors/smarkia/{smarkiaId})
	GetSmarkiaDistributor(ctx *gin.Context, smarkiaId string)
	// Get by id
	// (GET /distributors/{distributor_id})
	GetDistributor(ctx *gin.Context, distributorId string)
	// Get all measureEquipment
	// (GET /distributors/{distributor_id}/measure-equipments)
	ListMeasureEquipment(ctx *gin.Context, distributorId string, params ListMeasureEquipmentParams)
	// Get by id
	// (GET /distributors/{distributor_id}/measure-equipments/{id})
	GetMeasureEquipment(ctx *gin.Context, distributorId string, id string)
	// Get Meter config
	// (GET /distributors/{distributor_id}/meter_config/)
	GetMeterConfig(ctx *gin.Context, distributorId string, params GetMeterConfigParams)
	// Get by smarkia id
	// (GET /measure-equipments/smarkia/{smarkiaId})
	GetSmarkiaMeasureEquipment(ctx *gin.Context, smarkiaId string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []gin.HandlerFunc
}

// ListDistributors operation middleware
func (siw *ServerInterfaceWrapper) ListDistributors(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListDistributorsParams

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

	// ------------- Optional query parameter "q" -------------
	if paramValue := c.Query("q"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "q", c.Request.URL.Query(), &params.Q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter q: %s", err)})
		return
	}

	// ------------- Optional query parameter "sort" -------------
	if paramValue := c.Query("sort"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "sort", c.Request.URL.Query(), &params.Sort)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter sort: %s", err)})
		return
	}

	siw.Handler.ListDistributors(c, params)
}

// GetSmarkiaDistributor operation middleware
func (siw *ServerInterfaceWrapper) GetSmarkiaDistributor(c *gin.Context) {

	var err error

	// ------------- Path parameter "smarkiaId" -------------
	var smarkiaId string

	err = runtime.BindStyledParameter("simple", false, "smarkiaId", c.Param("smarkiaId"), &smarkiaId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter smarkiaId: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.GetSmarkiaDistributor(c, smarkiaId)
}

// GetDistributor operation middleware
func (siw *ServerInterfaceWrapper) GetDistributor(c *gin.Context) {

	var err error

	// ------------- Path parameter "distributor_id" -------------
	var distributorId string

	err = runtime.BindStyledParameter("simple", false, "distributor_id", c.Param("distributor_id"), &distributorId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter distributor_id: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.GetDistributor(c, distributorId)
}

// ListMeasureEquipment operation middleware
func (siw *ServerInterfaceWrapper) ListMeasureEquipment(c *gin.Context) {

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
	var params ListMeasureEquipmentParams

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

	// ------------- Optional query parameter "q" -------------
	if paramValue := c.Query("q"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "q", c.Request.URL.Query(), &params.Q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter q: %s", err)})
		return
	}

	// ------------- Optional query parameter "sort" -------------
	if paramValue := c.Query("sort"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "sort", c.Request.URL.Query(), &params.Sort)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter sort: %s", err)})
		return
	}

	siw.Handler.ListMeasureEquipment(c, distributorId, params)
}

// GetMeasureEquipment operation middleware
func (siw *ServerInterfaceWrapper) GetMeasureEquipment(c *gin.Context) {

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

	siw.Handler.GetMeasureEquipment(c, distributorId, id)
}

// GetMeterConfig operation middleware
func (siw *ServerInterfaceWrapper) GetMeterConfig(c *gin.Context) {

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
	var params GetMeterConfigParams

	// ------------- Optional query parameter "cups" -------------
	if paramValue := c.Query("cups"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "cups", c.Request.URL.Query(), &params.Cups)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter cups: %s", err)})
		return
	}

	// ------------- Required query parameter "date" -------------
	if paramValue := c.Query("date"); paramValue != "" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Query argument date is required, but not found"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "date", c.Request.URL.Query(), &params.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter date: %s", err)})
		return
	}

	siw.Handler.GetMeterConfig(c, distributorId, params)
}

// GetSmarkiaMeasureEquipment operation middleware
func (siw *ServerInterfaceWrapper) GetSmarkiaMeasureEquipment(c *gin.Context) {

	var err error

	// ------------- Path parameter "smarkiaId" -------------
	var smarkiaId string

	err = runtime.BindStyledParameter("simple", false, "smarkiaId", c.Param("smarkiaId"), &smarkiaId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("Invalid format for parameter smarkiaId: %s", err)})
		return
	}

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.GetSmarkiaMeasureEquipment(c, smarkiaId)
}

// RegisterHandlerListDistributors creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerListDistributors(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresListDistributors(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresListDistributors creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresListDistributors(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.ListDistributors)

	router.GET("/distributors", m...)

	return router
}

// RegisterHandlerGetSmarkiaDistributor creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetSmarkiaDistributor(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetSmarkiaDistributor(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetSmarkiaDistributor creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetSmarkiaDistributor(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetSmarkiaDistributor)

	router.GET("/distributors/smarkia/:smarkiaId", m...)

	return router
}

// RegisterHandlerGetDistributor creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetDistributor(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetDistributor(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetDistributor creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetDistributor(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetDistributor)

	router.GET("/distributors/:distributor_id", m...)

	return router
}

// RegisterHandlerListMeasureEquipment creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerListMeasureEquipment(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresListMeasureEquipment(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresListMeasureEquipment creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresListMeasureEquipment(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.ListMeasureEquipment)

	router.GET("/distributors/:distributor_id/measure-equipments", m...)

	return router
}

// RegisterHandlerGetMeasureEquipment creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetMeasureEquipment(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetMeasureEquipment(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetMeasureEquipment creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetMeasureEquipment(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetMeasureEquipment)

	router.GET("/distributors/:distributor_id/measure-equipments/:id", m...)

	return router
}

// RegisterHandlerGetMeterConfig creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetMeterConfig(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetMeterConfig(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetMeterConfig creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetMeterConfig(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetMeterConfig)

	router.GET("/distributors/:distributor_id/meter_config/", m...)

	return router
}

// RegisterHandlerGetSmarkiaMeasureEquipment creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGetSmarkiaMeasureEquipment(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGetSmarkiaMeasureEquipment(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGetSmarkiaMeasureEquipment creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGetSmarkiaMeasureEquipment(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GetSmarkiaMeasureEquipment)

	router.GET("/measure-equipments/smarkia/:smarkiaId", m...)

	return router
}
