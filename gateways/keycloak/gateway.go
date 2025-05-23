package keycloak

import (
	"context"
	models "go-api/models/keycloak"

	"go.uber.org/fx"
)

const (
	passwordGrant     = "password"
	refreshTokenGrant = "refresh_token"
)

type KeycloakGateway interface {
	GetOIDCToken(ctx context.Context, username, password string) (*models.GetOIDCTokenResponse, error)
	RefreshOIDCToken(ctx context.Context, refreshToken string) (*models.GetOIDCTokenResponse, error)
	RevokeOIDCToken(ctx context.Context, refreshToken string) error
	IntrospectOIDCToken(ctx context.Context, accessToken string) (*models.IntrospectOIDCTokenResponse, error)
}

type keycloakGateway struct {
	client KeycloakClient
}

type KeycloakGatewayParams struct {
	fx.In

	Client KeycloakClient
}

func NewKeycloakGateway(params KeycloakGatewayParams) (KeycloakGateway, error) {
	return &keycloakGateway{
		client: params.Client,
	}, nil
}

func (g *keycloakGateway) GetOIDCToken(ctx context.Context, username, password string) (*models.GetOIDCTokenResponse, error) {
	request := models.GetOIDCTokenRequest{
		GrantType: passwordGrant,
		Username:  username,
		Password:  password,
	}
	return g.client.GetOIDCToken(ctx, request)
}

func (g *keycloakGateway) RefreshOIDCToken(ctx context.Context, refreshToken string) (*models.GetOIDCTokenResponse, error) {
	request := models.GetOIDCTokenRequest{
		GrantType:    refreshTokenGrant,
		RefreshToken: refreshToken,
	}

	return g.client.GetOIDCToken(ctx, request)
}

func (g *keycloakGateway) RevokeOIDCToken(ctx context.Context, refreshToken string) error {
	request := models.RevokeOIDCTokenRequest{
		Token: refreshToken,
	}

	return g.client.RevokeOIDCToken(ctx, request)
}

func (g *keycloakGateway) IntrospectOIDCToken(ctx context.Context, accessToken string) (*models.IntrospectOIDCTokenResponse, error) {
	request := models.IntrospectOIDCTokenRequest{
		AccessToken: accessToken,
	}

	return g.client.IntrospectOIDCToken(ctx, request)
}
