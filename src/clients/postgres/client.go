package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"go-api/src/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PostgresClient interface {
	QuerySelect(ctx context.Context, result any, sqlQuery string, args ...any) error
	BeginTransaction(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

type postgresClient struct {
	config *config.Config
	logger *zap.Logger
	db     *sqlx.DB
}

type PostgresClientParams struct {
	fx.In

	Config *config.Config
	Logger *zap.Logger
}

func newConnection(config *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", config.PostgresConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(config.PostgresMaxOpenConns)
	db.SetMaxIdleConns(config.PostgresMaxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(config.PostgresMaxIdleTime) * time.Second)
	db.SetConnMaxLifetime(time.Duration(config.PostgresMaxLifetime) * time.Second)
	db.SetConnMaxLifetime(time.Duration(config.PostgresConnMaxLifetime) * time.Second)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func NewPostgresClient(params PostgresClientParams) (PostgresClient, error) {
	db, err := newConnection(params.Config)
	if err != nil {
		return nil, err
	}
	return &postgresClient{
		config: params.Config,
		logger: params.Logger,
		db:     db,
	}, nil
}

func (c *postgresClient) QuerySelect(ctx context.Context, result any, sqlQuery string, args ...any) error {
	return c.db.SelectContext(ctx, result, sqlQuery, args...)
}

func (c *postgresClient) BeginTransaction(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return c.db.BeginTxx(ctx, opts)
}
