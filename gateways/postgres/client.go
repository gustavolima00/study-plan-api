package postgres

import (
	"fmt"
	"go-api/config"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PostgresClient interface {
	// NewConnection creates a new PostgreSQL database connection
	NewConnection() (*sqlx.DB, error)
}

type postgresClient struct {
	config *config.Config
	logger *zap.Logger
}

type PostgresClientParams struct {
	fx.In

	Config *config.Config
	Logger *zap.Logger
}

func NewPostgresClient(params PostgresClientParams) PostgresClient {
	return &postgresClient{
		config: params.Config,
		logger: params.Logger,
	}
}

func (c *postgresClient) NewConnection() (*sqlx.DB, error) {
	logger := c.logger.With(
		zap.String("connection_string", c.config.PostgresConnectionString),
		zap.Int("max_open_conns", c.config.PostgresMaxOpenConns),
		zap.Int("max_idle_conns", c.config.PostgresMaxIdleConns),
		zap.Int("conn_max_idle_time", c.config.PostgresMaxIdleTime),
		zap.Int("conn_max_lifetime", c.config.PostgresConnMaxLifetime),
		zap.Int("max_lifetime", c.config.PostgresMaxLifetime),
	)
	db, err := sqlx.Connect("postgres", c.config.PostgresConnectionString)
	if err != nil {
		logger.Error("Failed to connect to PostgreSQL", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(c.config.PostgresMaxOpenConns)
	db.SetMaxIdleConns(c.config.PostgresMaxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(c.config.PostgresMaxIdleTime) * time.Second)
	db.SetConnMaxLifetime(time.Duration(c.config.PostgresMaxLifetime) * time.Second)
	db.SetConnMaxLifetime(time.Duration(c.config.PostgresConnMaxLifetime) * time.Second)

	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping PostgreSQL database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Debug("PostgreSQL connection settings")
	return db, nil
}
