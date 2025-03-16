package server

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"go-api/config/config"
)

// Params defines the dependencies that the server module needs
type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *config.Config
}

// New returns a pointer to Server
func New(p Params) *echo.Echo {
	e := echo.New()

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// .SetOnlineSince(time.Now())
			go e.Start(":" + p.Config.Port)
			return nil
		},
		OnStop: func(c context.Context) error {
			log.Println("Stopping server")
			return e.Shutdown(c)
		},
	})

	return e
}
