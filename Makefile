.PHONY: all deps run test swag-install swag mockery-install mock keycloak

all: deps run

deps:
	go mod tidy

run:
	go run main.go

test:
	go test ./...

swag-install:
	go install github.com/swaggo/swag/cmd/swag@latest

swag:
	swag init -o .internal/docs

mockery-install:
	go install github.com/vektra/mockery/v2@latest

mock:
	mockery

# Keycloak commands
KEYCLOAK_IMAGE ?= quay.io/keycloak/keycloak:26.2.4
KEYCLOAK_PORT ?= 8081
KEYCLOAK_ADMIN_USERNAME ?= admin
KEYCLOAK_ADMIN_PASSWORD ?= admin
KEYCLOAK_REALM ?= myrealm
KEYCLOAK_CLIENT_ID ?= myclient
KEYCLOAK_CLIENT_SECRET ?= mysecret

keycloak:
	@echo "Starting Keycloak with configuration:"
	@echo "-----------------------------------"
	@echo "Image: $(KEYCLOAK_IMAGE)"
	@echo "Port: $(KEYCLOAK_PORT)"
	@echo "Admin Username: $(KEYCLOAK_ADMIN_USERNAME)"
	@echo "Admin Password: $(KEYCLOAK_ADMIN_PASSWORD)"
	@echo "Realm: $(KEYCLOAK_REALM)"
	@echo "Client ID: $(KEYCLOAK_CLIENT_ID)"
	@echo "Client Secret: $(KEYCLOAK_CLIENT_SECRET)"
	@echo "-----------------------------------"

	docker run -d --name keycloak \
		-p $(KEYCLOAK_PORT):8080 \
		-e KEYCLOAK_ADMIN=$(KEYCLOAK_ADMIN_USERNAME) \
		-e KEYCLOAK_ADMIN_PASSWORD=$(KEYCLOAK_ADMIN_PASSWORD) \
		-e KC_HOSTNAME=localhost \
		$(KEYCLOAK_IMAGE) \
		start-dev

keycloak-stop:
	docker stop keycloak || true
	docker rm keycloak || true

keycloak-configure:
	@echo "Configuring Keycloak..."
	docker exec keycloak /opt/keycloak/bin/kcadm.sh config credentials --server http://localhost:8080 --realm master --user $(KEYCLOAK_ADMIN_USERNAME) --password $(KEYCLOAK_ADMIN_PASSWORD)
	docker exec keycloak /opt/keycloak/bin/kcadm.sh create realms -s realm=$(KEYCLOAK_REALM) -s enabled=true
	docker exec keycloak /opt/keycloak/bin/kcadm.sh create clients -r $(KEYCLOAK_REALM) -s clientId=$(KEYCLOAK_CLIENT_ID) -s secret=$(KEYCLOAK_CLIENT_SECRET) -s 'redirectUris=["*"]' -s 'webOrigins=["*"]' -s publicClient=false -s directAccessGrantsEnabled=true -s serviceAccountsEnabled=true
	@echo "Keycloak configuration complete!"
	@echo "Keycloak is ready and configured!"
	@echo "Access the admin console at: http://localhost:$(KEYCLOAK_PORT)/admin"
	@echo "Username: $(KEYCLOAK_ADMIN_USERNAME)"
	@echo "Password: $(KEYCLOAK_ADMIN_PASSWORD)"