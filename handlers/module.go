package handlers

import (
	"go-api/handlers/auth"
	"go-api/handlers/healthcheck"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		healthcheck.New,
		auth.NewAuthHandler,
	),
)
