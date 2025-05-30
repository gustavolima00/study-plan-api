package studysession

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeStart  EventType = "start"
	EventTypePause  EventType = "pause"
	EventTypeResume EventType = "resume"
	EventTypeStop   EventType = "stop"
)

type SessionEvent struct {
	EventType EventType `json:"event_type"`
	EventTime time.Time `json:"event_time"`
}

type SessionState string

const (
	SessionStateActive    SessionState = "active"
	SessionStateCompleted SessionState = "completed"
)

type StudySession struct {
	ID           uuid.UUID      `json:"id"`
	UserID       uuid.UUID      `json:"user_id"`
	Title        string         `json:"title"`
	Notes        string         `json:"notes"`
	Date         time.Time      `json:"date"`
	SessionState SessionState   `json:"session_state"`
	Events       []SessionEvent `json:"events"`
}
