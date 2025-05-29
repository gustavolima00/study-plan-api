package studysession

import (
	"time"

	"github.com/google/uuid"
)

type Subject struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type EventType string

type SessionEvent struct {
	EventType   EventType `json:"event_type"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	DeviceInfo  string    `json:"device_info"`
}

type SessionState string

type StudySession struct {
	ID           uuid.UUID      `json:"id"`
	UserID       uuid.UUID      `json:"user_id"`
	Title        string         `json:"title"`
	Notes        string         `json:"notes"`
	Date         time.Time      `json:"date"`
	SessionState SessionState   `json:"session_state"`
	Events       []SessionEvent `json:"events"`
	Subjects     []Subject      `json:"subjects"`
}
