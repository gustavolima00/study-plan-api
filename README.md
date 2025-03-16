# Go Sample API

This project is an example of a Go API using the `fx` framework. It demonstrates how to set up a simple API with health check endpoints, Swagger documentation, and mock services for testing.

## Table of Contents

- [Setup](#setup)
- [Running the API](#running-the-api)
- [Running Tests](#running-tests)
- [Updating Swagger Docs](#updating-swagger-docs)
- [Updating Mock Files](#updating-mock-files)

## Setup

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/go-api.git
    cd go-api
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

## Running the API

To run the API, use the following command:
```sh
go run main.go
```

The API will be available at `http://localhost:8080`.

## Running Tests

To run the tests, use the following command:
```sh
go test ./...
```

## Updating Swagger Docs

To update the Swagger documentation, follow these steps:

1. Install swag if you haven't already:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
``` 

2. Generate the Swagger docs:

```bash
swag init -o .internal/docs
```

This will generate a new `docs` directory with the updated Swagger files.

## Updating Mock Files

1. Install mockery if you haven't already:

```bash
go install github.com/vektra/mockery/v2@latest
```

2. Generate the mock files:

```bash
mockery
```



