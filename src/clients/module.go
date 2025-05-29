package clients

import (
	"go-api/src/clients/keycloak"
	"go-api/src/clients/postgres"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		keycloak.NewKeycloakClient,
		postgres.NewPostgresClient,
	),
)
