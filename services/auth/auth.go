package auth

import (
	"context"
	"go-api/gateways/keycloak"
	model "go-api/models/auth"

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
	keycloakGateway keycloak.KeycloakGateway
	logger          *zap.Logger
}

type AuthServiceParams struct {
	fx.In

	KeycloakGateway keycloak.KeycloakGateway
	Logger          *zap.Logger
}

func NewAuthService(params AuthServiceParams) AuthService {
	return &authService{
		keycloakGateway: params.KeycloakGateway,
		logger:          params.Logger,
	}
}

func (s *authService) CreateSession(ctx context.Context, request model.CreateSessionRequest) (*model.SessionInfo, error) {
	res, err := s.keycloakGateway.GetOIDCToken(ctx, request.Username, request.Password)
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
	token, err := s.keycloakGateway.RefreshOIDCToken(ctx, request.RefreshToken)
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
	err := s.keycloakGateway.RevokeOIDCToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &model.FinishSessionResponse{}, nil
}

func (s *authService) GetUserInfo(ctx context.Context, request model.VerifySessionRequest) (*model.UserInfo, error) {
	res, err := s.keycloakGateway.IntrospectOIDCToken(ctx, request.AccessToken)
	if err != nil {
		return nil, err
	}
	if res == nil {
		s.logger.Warn("User is nil")
		return nil, nil
	}

	return &model.UserInfo{
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
