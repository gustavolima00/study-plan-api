package auth

import (
	"context"
	"go-api/src/clients/keycloak"
	model "go-api/src/models/auth"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AuthService interface {
	CreateSession(ctx context.Context, request model.CreateSessionRequest) (*model.SessionInfo, error)
	UpdateSession(ctx context.Context, request model.UpdateSessionRequest) (*model.SessionInfo, error)
	FinishSession(ctx context.Context, request model.FinishSessionRequest) (*model.FinishSessionResponse, error)
	GetUserInfo(ctx context.Context, request model.VerifySessionRequest) (*model.UserInfo, error)
}

type authService struct {
	keycloakClient keycloak.KeycloakClient
	logger         *zap.Logger
}

type AuthServiceParams struct {
	fx.In

	KeycloakClient keycloak.KeycloakClient
	Logger         *zap.Logger
}

func NewAuthService(params AuthServiceParams) AuthService {
	return &authService{
		keycloakClient: params.KeycloakClient,
		logger:         params.Logger,
	}
}

func (s *authService) CreateSession(ctx context.Context, request model.CreateSessionRequest) (*model.SessionInfo, error) {
	res, err := s.keycloakClient.GetOIDCToken(ctx, keycloak.GetOIDCTokenRequest{
		GrantType: "password",
		Username:  request.Username,
		Password:  request.Password,
	})
	if err != nil {
		return nil, err
	}
	return &model.SessionInfo{
		AccessToken:      res.AccessToken,
		ExpiresIn:        res.ExpiresIn,
		RefreshExpiresIn: res.RefreshExpiresIn,
		RefreshToken:     res.RefreshToken,
		TokenType:        res.TokenType,
		SessionState:     res.SessionState,
		Scope:            res.Scope,
	}, nil
}

func (s *authService) UpdateSession(ctx context.Context, request model.UpdateSessionRequest) (*model.SessionInfo, error) {
	token, err := s.keycloakClient.GetOIDCToken(ctx, keycloak.GetOIDCTokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		return nil, err
	}
	return &model.SessionInfo{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}

func (s *authService) FinishSession(ctx context.Context, request model.FinishSessionRequest) (*model.FinishSessionResponse, error) {
	err := s.keycloakClient.RevokeOIDCToken(ctx, keycloak.RevokeOIDCTokenRequest{
		Token: request.AccessToken,
	})
	if err != nil {
		return nil, err
	}
	return &model.FinishSessionResponse{}, nil
}

func (s *authService) GetUserInfo(ctx context.Context, request model.VerifySessionRequest) (*model.UserInfo, error) {
	res, err := s.keycloakClient.IntrospectOIDCToken(ctx, keycloak.IntrospectOIDCTokenRequest{
		AccessToken: request.AccessToken,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Warn("User is nil")
		return nil, nil
	}
	userID, err := uuid.Parse(res.Sub)
	if err != nil {
		return nil, err
	}

	return &model.UserInfo{
		ID:                userID,
		Username:          res.Username,
		Email:             res.Email,
		EmailVerified:     res.EmailVerified,
		Name:              res.Name,
		PreferredUsername: res.PreferredUsername,
		GivenName:         res.GivenName,
		FamilyName:        res.FamilyName,
		ResourceAccess: model.ResourceAccess{
			Account: model.Account{
				Roles: res.ResourceAccess.Account.Roles,
			},
		},
	}, nil
}
