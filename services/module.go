package services

import (
	"go.uber.org/fx"

	"go-api/services/healthcheck"
)

var Module = fx.Options(
	fx.Provide(
		healthcheck.New,
	),
)
