package config

import (
	"go.uber.org/fx"

	"go-api/config/config"
)

var Module = fx.Options(
	fx.Provide(
		config.New,
	),
)
