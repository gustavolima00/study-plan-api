package routes

import (
	// Generate automatically the swagger docs

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"

	_ "go-api/docs"
)

type Params struct {
	fx.In

	Echo *echo.Echo
}

func Register(p Params) {
	// p.Echo.GET("/", hc.GetAPIStatus)

	p.Echo.GET("/swagger/*any", echoSwagger.WrapHandler)
}
