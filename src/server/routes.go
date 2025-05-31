package server

import (
	_ "go-api/.internal/docs" // Generate automatically the swagger docs
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
	Middlewares         middlewares.Middlewares
}

// RegisterRoutes registers the routes for the API.
func RegisterRoutes(p RegisterRoutesParams) {
	// Base routes
	p.Echo.GET("/", p.Healthcheck.GetAPIStatus)
	p.Echo.GET("/swagger/*any", echoSwagger.WrapHandler)

	// Authentication routes
	authGroup := p.Echo.Group("/auth")
	{
		authGroup.POST("/login", p.AuthHandler.CreateSession)
		authGroup.POST("/refresh", p.AuthHandler.UpdateSession)
		authGroup.POST("/logout", p.AuthHandler.FinishSession)
		authGroup.GET("/user", p.AuthHandler.GetUser, p.Middlewares.AuthMiddleware())
	}

	// StudySession routes
	studySessionGroup := p.Echo.Group("/study-session", p.Middlewares.AuthMiddleware())
	{
		studySessionGroup.POST("/start", p.StudySessionHandler.StartStudySession)
		studySessionGroup.GET("", p.StudySessionHandler.GetActiveStudySession)
		studySessionGroup.GET("/events", p.StudySessionHandler.GetActiveStudySessionEvents)
		studySessionGroup.POST("/events", p.StudySessionHandler.AddStudySessionEvents)
		studySessionGroup.POST("/finish", p.StudySessionHandler.FinishStudySession)
	}
}
