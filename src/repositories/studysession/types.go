package studysession

import (
	models "go-api/src/models/studysession"
	"time"

	"github.com/google/uuid"
)

type DBSessionEvent struct {
	ID        string    `db:"id" json:"id"`
	SessionID string    `db:"session_id" json:"session_id"`
	EventType string    `db:"event_type" json:"event_type"`
	EventTime time.Time `db:"event_time" json:"event_time"`
}

type DBStudySession struct {
	ID           string    `db:"id" json:"id"`
	UserID       string    `db:"user_id" json:"user_id"`
	Title        string    `db:"title" json:"title"`
	Notes        string    `db:"notes" json:"notes"`
	Date         time.Time `db:"date" json:"date"`
	SessionState string    `db:"session_state" json:"session_state"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func (e DBSessionEvent) ToSessionEvent() models.SessionEvent {
	return models.SessionEvent{
		EventType: models.EventType(e.EventType),
		EventTime: e.EventTime,
	}
}

func (s DBStudySession) ToStudySession(dbEvents []DBSessionEvent) (*models.StudySession, error) {
	id, err := uuid.Parse(s.ID)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(s.UserID)
	if err != nil {
		return nil, err
	}
	events := make([]models.SessionEvent, len(dbEvents))

	for i, dbEvent := range dbEvents {
		events[i] = dbEvent.ToSessionEvent()
	}
	session := models.StudySession{
		ID:           id,
		UserID:       userID,
		Title:        s.Title,
		Notes:        s.Notes,
		Date:         s.Date,
		SessionState: models.SessionState(s.SessionState),
		Events:       events,
	}
	return &session, nil
}
