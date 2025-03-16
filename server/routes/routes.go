package routes

import (
	// Generate automatically the swagger docs

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"

	_ "go-api/docs"
	"go-api/handlers/healthcheck"
)

// Params defines the dependencies for the routes module.
type Params struct {
	fx.In

	Echo        *echo.Echo
	Healthcheck healthcheck.Handler
}

// Register registers the routes for the API.
func Register(p Params) {
	p.Echo.GET("/", p.Healthcheck.GetAPIStatus)

	p.Echo.GET("/swagger/*any", echoSwagger.WrapHandler)
}
