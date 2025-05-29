package server

import (
	"context"
	"go-api/config"
	"go-api/src/services/studysession"
	"log"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Params defines the dependencies that the server module needs
type ServerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *config.Config

	Service studysession.StudySessionService
}

// NewServer returns a pointer to Server
func NewServer(p ServerParams) *echo.Echo {
	e := echo.New()

	p.Service.Sample()

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
