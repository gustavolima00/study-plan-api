package middlewares

import (
	"go-api/src/services/auth"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Middlewares interface {
	AuthMiddleware() echo.MiddlewareFunc
}

type middlewares struct {
	logger      *zap.Logger
	authService auth.AuthService
}

type MiddlewaresParams struct {
	fx.In

	Logger      *zap.Logger
	AuthService auth.AuthService
}

func NewMiddlewares(params MiddlewaresParams) Middlewares {
	return &middlewares{
		logger:      params.Logger,
		authService: params.AuthService,
	}
}
