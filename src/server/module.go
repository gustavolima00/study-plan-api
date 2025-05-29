package server

import (
	"go-api/src/server/middlewares"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		middlewares.NewMiddlewares,
		NewServer,
	),
	fx.Invoke(
		RegisterRoutes,
	),
)
