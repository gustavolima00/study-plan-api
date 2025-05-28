package auth

import (
	authmodel "go-api/src/models/auth"
	"go-api/src/models/constants"
	authservice "go-api/src/services/auth"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// AuthHandler defines the interface for authentication API handlers
type AuthHandler interface {
	// CreateSession authenticates a user and returns tokens
	CreateSession(e echo.Context) error

	// UpdateSession generates new tokens using a refresh token
	UpdateSession(e echo.Context) error

	// FinishSession revokes user tokens and ends the session
	FinishSession(e echo.Context) error

	// GetUser ...
	GetUser(e echo.Context) error
}

// AuthHandlerParams defines the dependencies for the auth module
type AuthHandlerParams struct {
	fx.In

	AuthService authservice.AuthService
	Logger      *zap.Logger
}

type authHandler struct {
	authService authservice.AuthService
	logger      *zap.Logger
}

// NewAuthHandler creates a new auth handler with injected dependencies
func NewAuthHandler(p AuthHandlerParams) AuthHandler {
	return &authHandler{
		authService: p.AuthService,
		logger:      p.Logger,
	}
}

// Login authenticates a user
//
//	@Summary		User login
//	@Description	Authenticate user and return tokens
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		authmodel.CreateSessionRequest	true	"Login credentials"
//	@Success		200		{object}	authmodel.SessionInfo
//	@Failure		400		{string}	string	"Invalid request"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/auth/login [post]
func (h *authHandler) CreateSession(e echo.Context) error {
	var req authmodel.CreateSessionRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid request format")
	}
	ctx := e.Request().Context()

	res, err := h.authService.CreateSession(ctx, req)
	if err != nil {
		h.logger.Error("Failed to create session:", zap.Error(err))
		return e.JSON(http.StatusUnauthorized, "Invalid credentials")
	}

	accessTokens := authmodel.SessionInfo{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
	}

	return e.JSON(http.StatusOK, accessTokens)
}

// Refresh generates new tokens
//
//	@Summary		Refresh tokens
//	@Description	Generate new access token using refresh token
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		authmodel.UpdateSessionRequest	true	"Refresh token"
//	@Success		200		{object}	authmodel.SessionInfo
//	@Failure		400		{string}	string	"Invalid request"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/auth/refresh [post]
func (h *authHandler) UpdateSession(e echo.Context) error {
	var req authmodel.UpdateSessionRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid request format")
	}
	ctx := e.Request().Context()

	res, err := h.authService.UpdateSession(ctx, req)
	if err != nil {
		h.logger.Error("Failed to refresh token", zap.Error(err))
		return e.JSON(http.StatusUnauthorized, "Invalid refresh token")
	}

	accessTokens := authmodel.SessionInfo{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
	}

	return e.JSON(http.StatusOK, accessTokens)
}

// Finish session revoke user tokens
//
//	@Summary		Logout and revoke user tokens
//	@Description	Revoke user tokens and end session
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		authmodel.FinishSessionRequest	true	"Refresh token"
//	@Success		200 	{object}	authmodel.FinishSessionResponse		"Success response"
//	@Failure		400		{string}	string	"Invalid request"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/auth/logout [post]
func (h *authHandler) FinishSession(e echo.Context) error {
	var req authmodel.FinishSessionRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid request format")
	}
	ctx := e.Request().Context()

	res, err := h.authService.FinishSession(ctx, req)
	if err != nil {
		h.logger.Error("Failed to revoke token", zap.Error(err))
		return e.JSON(http.StatusBadRequest, "Invalid refresh token")
	}
	return e.JSON(http.StatusOK, res)
}

// GetUser returns authenticated user information
//
//	@Summary		Get user info
//	@Description	Returns information about the currently authenticated user
//	@Tags			authentication
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200		{object}	authmodel.UserInfo	"User information"
//	@Failure		401		{string}	string			"Unauthorized - Missing or invalid token"
//	@Failure		500		{string}	string			"Internal server error"
//	@Router			/auth/user [get]
func (h *authHandler) GetUser(e echo.Context) error {
	ctx := e.Request().Context()
	userInfo := ctx.Value(constants.ContextKeyUserInfoKey).(*authmodel.UserInfo)
	return e.JSON(http.StatusOK, userInfo)
}
