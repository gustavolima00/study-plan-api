.PHONY: all deps run test swag-install swag mockery-install mock

all: deps run

deps:
	go mod tidy

run: deps
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