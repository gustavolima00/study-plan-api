package healthcheck

import (
	"time"

	"go.uber.org/fx"
)

// Healthcheck interface define functions
// that returns the database connection status
// last time the sync was done and the system status
type Healthcheck interface {
	// SetOnlineSince sets the time the system was online
	SetOnlineSince(time.Time)

	// OnlineSince returns the time since the system was online
	OnlineSince() time.Duration
}

// Params defines the dependencies that the healthcheck module needs
type Params struct {
	fx.In
}

type hc struct {
	onlineSince time.Time
}

// New returns an implementation of Healthcheck interface
func New(p Params) Healthcheck {
	return &hc{}
}

func (h *hc) SetOnlineSince(t time.Time) {
	h.onlineSince = t
}

func (h *hc) OnlineSince() time.Duration {
	return time.Since(h.onlineSince)
}
