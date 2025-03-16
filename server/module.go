package server

import (
	"go.uber.org/fx"

	"go-api/server/routes"
	"go-api/server/server"
)

var Module = fx.Options(
	fx.Provide(
		server.New,
	),
	fx.Invoke(
		routes.Register,
	),
)
