package gateways

import (
	"go-api/src/gateways/keycloak"
	"go-api/src/gateways/postgres"

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
