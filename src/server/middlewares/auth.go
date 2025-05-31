package middlewares

import (
	"context"
	authmodel "go-api/src/models/auth"
	"go-api/src/models/constants"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (m middlewares) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := getBaerToken(c)
			if token == "" {
				m.logger.Debug("Missing Bearer token")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Bearer token is required",
				})
			}

			// Verificar o token no Keycloak
			userInfo, err := m.authService.GetUserInfo(c.Request().Context(), authmodel.VerifySessionRequest{
				AccessToken: token,
			})
			if err != nil {
				m.logger.Error("Failed to authenticate user", zap.Error(err))
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Invalid or expired token",
				})
			}

			ctx := context.WithValue(c.Request().Context(), constants.ContextKeyUserInfoKey, userInfo)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func getBaerToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
