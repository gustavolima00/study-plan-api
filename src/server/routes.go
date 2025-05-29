package server

import (
	_ "go-api/.internal/docs"
	// Generate automatically the swagger docs
	"go-api/src/handlers/auth"
	"go-api/src/handlers/healthcheck"
	"go-api/src/handlers/studysession"
	"go-api/src/server/middlewares"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"
)

// Params defines the dependencies for the routes module.
type RegisterRoutesParams struct {
	fx.In

	Echo                *echo.Echo
	Healthcheck         healthcheck.Handler
	AuthHandler         auth.AuthHandler
	StudySessionHandler studysession.StudySessionHandler

	Middlewares middlewares.Middlewares
}

// RegisterRoutes registers the routes for the API.
func RegisterRoutes(p RegisterRoutesParams) {
	p.Echo.GET("/", p.Healthcheck.GetAPIStatus)

	p.Echo.GET("/swagger/*any", echoSwagger.WrapHandler)

	// Authentication routes
	p.Echo.POST("/auth/login", p.AuthHandler.CreateSession)
	p.Echo.POST("/auth/refresh", p.AuthHandler.UpdateSession)
	p.Echo.POST("/auth/logout", p.AuthHandler.FinishSession)
	p.Echo.GET("/auth/user", p.AuthHandler.GetUser, p.Middlewares.AuthMiddleware())

	// StudySession routes
	p.Echo.POST("/study-session/upsert-active", p.StudySessionHandler.UpsertActiveStudySession, p.Middlewares.AuthMiddleware())

}
