package handlers

import (
	"go-api/src/handlers/auth"
	"go-api/src/handlers/healthcheck"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		healthcheck.New,
		auth.NewAuthHandler,
	),
)
