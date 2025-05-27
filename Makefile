.PHONY: all deps run test swag-install swag mockery-install mock \
        postgres postgres-stop \
        keycloak keycloak-stop keycloak-configure \
        migrate-install migrate-create migrate-up migrate-down migrate-force

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

# PostgreSQL commands
POSTGRES_IMAGE ?= postgres:15-alpine
POSTGRES_PORT ?= 5432
POSTGRES_USER ?= myuser
POSTGRES_PASSWORD ?= mypassword
POSTGRES_DB ?= mydb
KEYCLOAK_DB ?= keycloak_db

postgres:
	@echo "Starting PostgreSQL with configuration:"
	@echo "-------------------------------------"
	@echo "Image: $(POSTGRES_IMAGE)"
	@echo "Port: $(POSTGRES_PORT)"
	@echo "Username: $(POSTGRES_USER)"
	@echo "Password: $(POSTGRES_PASSWORD)"
	@echo "Application Database: $(POSTGRES_DB)"
	@echo "Keycloak Database: $(KEYCLOAK_DB)"
	@echo "-------------------------------------"
	
	docker run -d --name postgres \
		-p $(POSTGRES_PORT):5432 \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-e POSTGRES_MULTIPLE_DATABASES="$(POSTGRES_DB),$(KEYCLOAK_DB)" \
		-v postgres_data:/var/lib/postgresql/data \
		-v $(PWD)/migrations:/migrations \
		-v $(PWD)/docker-entrypoint-initdb.d:/docker-entrypoint-initdb.d \
		$(POSTGRES_IMAGE)

postgres-stop:
	docker stop postgres || true
	docker rm postgres || true

# Keycloak commands
KEYCLOAK_IMAGE ?= quay.io/keycloak/keycloak:26.2.4
KEYCLOAK_PORT ?= 8088
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
	@echo "Database: $(KEYCLOAK_DB)"
	@echo "-----------------------------------"

	docker run -d --name keycloak \
		-p $(KEYCLOAK_PORT):8080 \
		-e KEYCLOAK_ADMIN=$(KEYCLOAK_ADMIN_USERNAME) \
		-e KEYCLOAK_ADMIN_PASSWORD=$(KEYCLOAK_ADMIN_PASSWORD) \
		-e KC_DB=postgres \
		-e KC_DB_URL=jdbc:postgresql://postgres:5432/$(KEYCLOAK_DB) \
		-e KC_DB_USERNAME=$(POSTGRES_USER) \
		-e KC_DB_PASSWORD=$(POSTGRES_PASSWORD) \
		-e KC_HOSTNAME=localhost \
		--link postgres:postgres \
		$(KEYCLOAK_IMAGE) \
		start-dev

keycloak-stop:
	docker stop keycloak || true
	docker rm keycloak || true

# Keycloak configuration with retry logic
keycloak-configure:
	@echo "Waiting for Keycloak to be ready..."
	@until docker exec keycloak /opt/keycloak/bin/kcadm.sh config credentials --server http://localhost:8080 --realm master --user $(KEYCLOAK_ADMIN_USERNAME) --password $(KEYCLOAK_ADMIN_PASSWORD) >/dev/null; do \
		echo "Keycloak not ready yet - retrying in 5 seconds..."; \
		sleep 5; \
	done
	
	@echo "Configuring Keycloak..."
	docker exec keycloak /opt/keycloak/bin/kcadm.sh config credentials --server http://localhost:8080 --realm master --user $(KEYCLOAK_ADMIN_USERNAME) --password $(KEYCLOAK_ADMIN_PASSWORD)
	docker exec keycloak /opt/keycloak/bin/kcadm.sh create realms -s realm=$(KEYCLOAK_REALM) -s enabled=true
	docker exec keycloak /opt/keycloak/bin/kcadm.sh create clients -r $(KEYCLOAK_REALM) -s clientId=$(KEYCLOAK_CLIENT_ID) -s secret=$(KEYCLOAK_CLIENT_SECRET) -s 'redirectUris=["*"]' -s 'webOrigins=["*"]' -s publicClient=false -s directAccessGrantsEnabled=true -s serviceAccountsEnabled=true
	@echo "Keycloak configuration complete!"
	@echo "Keycloak is ready and configured!"
	@echo "Access the admin console at: http://localhost:$(KEYCLOAK_PORT)/admin"
	@echo "Username: $(KEYCLOAK_ADMIN_USERNAME)"
	@echo "Password: $(KEYCLOAK_ADMIN_PASSWORD)"

# Database migration commands
POSTGRES_CONNECTION_STRING ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-create:
	migrate create -ext sql -dir migrations -seq ${name}

migrate-up:
	migrate -path migrations -database "${POSTGRES_CONNECTION_STRING}" up

migrate-down:
	migrate -path migrations -database "${POSTGRES_CONNECTION_STRING}" down

migrate-force:
	migrate -path migrations -database "${POSTGRES_CONNECTION_STRING}" force ${version}

# Combined commands
setup: deps postgres keycloak migrate-up keycloak-configure 
	@echo "All services are up and configured!"

teardown: keycloak-stop postgres-stop
	@echo "All services are down!"