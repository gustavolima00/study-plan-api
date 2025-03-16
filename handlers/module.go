package handlers

import (
	"go.uber.org/fx"

	"go-api/handlers/healthcheck"
)

var Module = fx.Options(
	fx.Provide(
		healthcheck.New,
	),
)
