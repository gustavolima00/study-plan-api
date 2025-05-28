package services

import (
	"go-api/src/services/auth"
	"go-api/src/services/healthcheck"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		healthcheck.New,
		auth.NewAuthService,
	),
)
