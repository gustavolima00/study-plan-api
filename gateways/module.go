package gateways

import (
	"go-api/gateways/keycloak"
	"go-api/gateways/postgres"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		keycloak.NewKeycloakClient,
		keycloak.NewKeycloakGateway,
	),
	fx.Provide(
		postgres.NewPostgresClient,
	),
)
