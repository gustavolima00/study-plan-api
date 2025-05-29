package auth

import "github.com/google/uuid"

type CreateSessionRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateSessionRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type FinishSessionRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type FinishSessionResponse struct{}

type VerifySessionRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type SessionInfo struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type Account struct {
	Roles []string `json:"roles"`
}

type ResourceAccess struct {
	Account Account `json:"account"`
}

type UserInfo struct {
	ID                uuid.UUID      `json:"uuid"`
	Username          string         `json:"username"`
	Email             string         `json:"email"`
	EmailVerified     bool           `json:"email_verified"`
	Name              string         `json:"name"`
	PreferredUsername string         `json:"preferred_username"`
	GivenName         string         `json:"given_name"`
	FamilyName        string         `json:"family_name"`
	ResourceAccess    ResourceAccess `json:"resource_access"`
}
