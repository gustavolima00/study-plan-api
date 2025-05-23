package services

import (
	"go-api/services/auth"
	"go-api/services/healthcheck"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		healthcheck.New,
		auth.NewAuthService,
	),
)
