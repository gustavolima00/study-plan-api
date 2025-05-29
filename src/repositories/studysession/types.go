package studysession

import (
	"encoding/json"
	"fmt"
	models "go-api/src/models/studysession"
	"time"

	"github.com/google/uuid"
)

type DBSubject struct {
	ID          string     `db:"id" json:"id"`
	UserID      string     `db:"user_id" json:"user_id"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
}

type DBSessionEvent struct {
	ID          string     `db:"id" json:"id"`
	SessionID   string     `db:"user_id" json:"user_id"`
	EventType   string     `db:"title" json:"title"`
	Description string     `db:"notes" json:"notes"`
	EventTime   *time.Time `db:"event_time" json:"event_time"`
}

type DBStudySession struct {
	ID           string          `db:"id" json:"id"`
	UserID       string          `db:"user_id" json:"user_id"`
	Title        string          `db:"title" json:"title"`
	Notes        string          `db:"notes" json:"notes"`
	Date         time.Time       `db:"date" json:"date"`
	SessionState string          `db:"session_state" json:"session_state"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updated_at"`
	EventsRaw    json.RawMessage `db:"events"`
	SubjectsRaw  json.RawMessage `db:"subjects"`
}

func (s *DBStudySession) ToStudySession() (*models.StudySession, error) {
	id, err := uuid.Parse(s.ID)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(s.UserID)
	if err != nil {
		return nil, err
	}
	session := models.StudySession{
		ID:           id,
		UserID:       userID,
		Title:        s.Title,
		Notes:        s.Notes,
		Date:         s.Date,
		SessionState: models.SessionState(s.SessionState),
	}
	if len(s.SubjectsRaw) > 0 {
		if err := json.Unmarshal(s.SubjectsRaw, &session.Subjects); err != nil {
			return nil, fmt.Errorf("failed to unmarshal subjects: %w", err)
		}
	}

	if len(s.EventsRaw) > 0 {
		if err := json.Unmarshal(s.EventsRaw, &session.Events); err != nil {
			return nil, fmt.Errorf("failed to unmarshal events: %w", err)
		}
	}

	return &session, nil
}
