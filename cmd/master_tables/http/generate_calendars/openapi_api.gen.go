// Package generate_calendars provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package generate_calendars

import "github.com/gin-gonic/gin"

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Generate all calendars
	// (POST /admin/master-tables/generate/calendars)
	GenerateCalendarsPeriods(ctx *gin.Context)
	// Generate Festive Days
	// (POST /admin/master-tables/generate/festive_days)
	GenerateFestiveDays(ctx *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []gin.HandlerFunc
}

// GenerateCalendarsPeriods operation middleware
func (siw *ServerInterfaceWrapper) GenerateCalendarsPeriods(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.GenerateCalendarsPeriods(c)
}

// GenerateFestiveDays operation middleware
func (siw *ServerInterfaceWrapper) GenerateFestiveDays(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{""})

	siw.Handler.GenerateFestiveDays(c)
}

// RegisterHandlerGenerateCalendarsPeriods creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGenerateCalendarsPeriods(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGenerateCalendarsPeriods(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGenerateCalendarsPeriods creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGenerateCalendarsPeriods(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GenerateCalendarsPeriods)

	router.POST("/admin/master-tables/generate/calendars", m...)

	return router
}

// RegisterHandlerGenerateFestiveDays creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlerGenerateFestiveDays(router *gin.RouterGroup, si ServerInterface) *gin.RouterGroup {
	return RegisterHandlerWithMiddlewaresGenerateFestiveDays(router, si, []gin.HandlerFunc{})
}

// RegisterHandlerWithMiddlewaresGenerateFestiveDays creates http.Handler with additional options
func RegisterHandlerWithMiddlewaresGenerateFestiveDays(router *gin.RouterGroup, si ServerInterface, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: middlewares,
	}

	m := append(middlewares, wrapper.GenerateFestiveDays)

	router.POST("/admin/master-tables/generate/festive_days", m...)

	return router
}
