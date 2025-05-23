package server

import (
	// Generate automatically the swagger docs

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"

	_ "go-api/.internal/docs"
	"go-api/handlers/healthcheck"
)

// Params defines the dependencies for the routes module.
type RegisterRoutesParams struct {
	fx.In

	Echo        *echo.Echo
	Healthcheck healthcheck.Handler
}

// RegisterRoutes registers the routes for the API.
func RegisterRoutes(p RegisterRoutesParams) {
	p.Echo.GET("/", p.Healthcheck.GetAPIStatus)

	p.Echo.GET("/swagger/*any", echoSwagger.WrapHandler)
}
