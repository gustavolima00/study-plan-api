package handlers

import (
	"go-api/src/handlers/auth"
	"go-api/src/handlers/healthcheck"
	"go-api/src/handlers/studysession"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		healthcheck.New,
		auth.NewAuthHandler,
		studysession.NewStudySessionHandler,
	),
)
