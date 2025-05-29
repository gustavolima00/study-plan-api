package studysession

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"go-api/src/gateways/postgres"
	models "go-api/src/models/database/studysession"
	"io/fs"
	"path/filepath"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:embed *.sql
var sqlFiles embed.FS

const (
	// SQL Queries
	getUserStudySessionsQueryKey = "get_user_study_sessions"
	getStudySessionQueryKey      = "get_study_session"
)

type StudySessionRepository interface {
	GetUserStudySessions(ctx context.Context, userID uuid.UUID) ([]models.StudySession, error)
	GetUserStudySession(ctx context.Context, sessionID uuid.UUID) (*models.StudySession, error)
}

type studySessionRepository struct {
	logger     *zap.Logger
	pgclient   postgres.PostgresClient
	sqlQueries map[string]string
}

type StudySessionRepositoryParams struct {
	fx.In

	Logger   *zap.Logger
	PGClient postgres.PostgresClient
}

func NewStudySessionRepository(p StudySessionRepositoryParams) (StudySessionRepository, error) {
	queries := make(map[string]string)
	err := fs.WalkDir(sqlFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		queryName := filepath.Base(path[:len(path)-len(filepath.Ext(path))])
		content, err := fs.ReadFile(sqlFiles, path)
		if err != nil {
			return err
		}

		queries[queryName] = string(content)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &studySessionRepository{
		logger:     p.Logger,
		pgclient:   p.PGClient,
		sqlQueries: queries,
	}, nil
}

type rawSession struct {
	models.StudySession

	EventsRaw   json.RawMessage `db:"events"`
	SubjectsRaw json.RawMessage `db:"subjects"`
}

func mapSession(rawSession rawSession) (*models.StudySession, error) {
	if err := json.Unmarshal(rawSession.SubjectsRaw, &rawSession.Subjects); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subjects: %w", err)
	}

	if err := json.Unmarshal(rawSession.EventsRaw, &rawSession.Events); err != nil {
		return nil, fmt.Errorf("failed to unmarshal events: %w", err)
	}
	return &rawSession.StudySession, nil
}

func mapSessions(rawSessions []rawSession) ([]models.StudySession, error) {
	sessions := make([]models.StudySession, len(rawSessions))
	for i, rawSession := range rawSessions {
		result, err := mapSession(rawSession)
		if err != nil {
			return nil, err
		}
		sessions[i] = *result
	}
	return sessions, nil
}

func (r *studySessionRepository) GetUserStudySessions(ctx context.Context, userID uuid.UUID) ([]models.StudySession, error) {
	var rawSessions []rawSession
	query := r.sqlQueries[getUserStudySessionsQueryKey]
	err := r.pgclient.QuerySelect(ctx, &rawSessions, query, userID.String())
	if err != nil {
		return nil, err
	}
	return mapSessions(rawSessions)
}

func (r *studySessionRepository) GetUserStudySession(ctx context.Context, sessionID uuid.UUID) (*models.StudySession, error) {
	var rawSessions []rawSession
	query := r.sqlQueries[getStudySessionQueryKey]
	err := r.pgclient.QuerySelect(ctx, &rawSessions, query, sessionID.String())
	if err != nil {
		return nil, err
	}
	if len(rawSessions) == 0 {
		return nil, nil
	}
	return mapSession(rawSessions[0])
}
