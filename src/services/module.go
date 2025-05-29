package services

import (
	"go-api/src/services/auth"
	"go-api/src/services/healthcheck"
	"go-api/src/services/studysession"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		healthcheck.New,
		auth.NewAuthService,
		studysession.NewStudySessionService,
	),
)
