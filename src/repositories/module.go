package repositories

import (
	"go-api/src/repositories/studysession"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		studysession.NewStudySessionRepository,
	),
)
