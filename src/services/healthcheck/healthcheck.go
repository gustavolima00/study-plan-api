package healthcheck

import (
	"errors"
	"go-api/src/gateways/postgres"
	"time"

	"go.uber.org/fx"
)

// Service interface define functions
// that returns the database connection status
// last time the sync was done and the system status
type Service interface {
	// SetOnlineSince sets the time the system was online
	SetOnlineSince(time.Time)

	// OnlineSince returns the time since the system was online
	OnlineSince() (time.Duration, error)

	// IsPostgresRunning returns the status of the PostgreSQL connection
	IsPostgresRunning() bool
}

// Params defines the dependencies that the healthcheck module needs
type Params struct {
	fx.In

	PostgresClient postgres.PostgresClient
}

type service struct {
	onlineSince    *time.Time
	postgresClient postgres.PostgresClient
}

// New returns an implementation of Healthcheck interface
func New(p Params) Service {
	now := time.Now()
	return &service{
		onlineSince:    &now,
		postgresClient: p.PostgresClient,
	}
}

func (s *service) SetOnlineSince(t time.Time) {
	s.onlineSince = &t
}

func (s *service) OnlineSince() (time.Duration, error) {
	if s.onlineSince == nil {
		return time.Duration(0), errors.New("online since is not set")
	}
	return time.Since(*s.onlineSince), nil
}

func (s *service) IsPostgresRunning() bool {
	if s.postgresClient == nil {
		return false
	}

	db, err := s.postgresClient.NewConnection()
	if err != nil {
		return false
	}
	defer db.Close()

	return db != nil
}
