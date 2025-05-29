package studysession

import "time"

type Subject struct {
	ID          string     `db:"id" json:"id"`
	UserID      string     `db:"user_id" json:"user_id"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
}

type SessionEvent struct {
	ID          string     `db:"id" json:"id"`
	SessionID   string     `db:"user_id" json:"user_id"`
	EventType   int        `db:"title" json:"title"`
	Description string     `db:"notes" json:"notes"`
	EventTime   *time.Time `db:"event_time" json:"event_time"`
	DeviceInfo  string     `db:"device_info" json:"device_info"`
}

type StudySession struct {
	ID           string    `db:"id" json:"id"`
	UserID       string    `db:"user_id" json:"user_id"`
	Title        string    `db:"title" json:"title"`
	Notes        string    `db:"notes" json:"notes"`
	Date         time.Time `db:"date" json:"date"`
	SessionState int       `db:"session_state" json:"session_state"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	Events       []SessionEvent
	Subjects     []Subject
}
