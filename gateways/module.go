package gateways

import (
	"go-api/gateways/keycloak"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		keycloak.NewKeycloakClient,
		keycloak.NewKeycloakGateway,
	),
)
