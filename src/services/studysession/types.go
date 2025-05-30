package studysession

import (
	models "go-api/src/models/studysession"
	"time"
)

type UpsertActiveStudySessionRequest struct {
	Title string `json:"title"`
	Notes string `json:"notes"`
}

type AddStudySessionEventsRequest struct {
	Events []models.SessionEvent `json:"events"`
}

type FinishStudySessionRequest struct {
	FinishedAt time.Time `json:"finished_at"`
}
