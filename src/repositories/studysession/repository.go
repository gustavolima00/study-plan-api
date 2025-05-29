package studysession

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"go-api/src/clients/postgres"
	models "go-api/src/models/studysession"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gotidy/ptr"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

type StudySessionRepository interface {
	UpsertActiveStudySession(ctx context.Context, userID uuid.UUID, session models.StudySession) (*models.StudySession, error)
	GetUserStudySessions(ctx context.Context, userID uuid.UUID) ([]models.StudySession, error)
	GetUserActiveStudySession(ctx context.Context, userID uuid.UUID) (*models.StudySession, error)
	AddSessionEvents(ctx context.Context, sessionID uuid.UUID, events []models.SessionEvent) error
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

func (r *studySessionRepository) GetUserStudySessions(ctx context.Context, userID uuid.UUID) ([]models.StudySession, error) {
	var rawSessions []DBStudySession
	query := r.sqlQueries["get_user_study_sessions"]
	err := r.pgclient.QuerySelect(ctx, &rawSessions, query, userID.String())
	if err != nil {
		return nil, err
	}
	sessions := make([]models.StudySession, len(rawSessions))
	for i, rawSession := range rawSessions {
		session, err := rawSession.ToStudySession()
		if err != nil {
			return nil, err
		}
		sessions[i] = *session
	}
	return sessions, nil
}

func (r *studySessionRepository) GetUserActiveStudySession(ctx context.Context, userID uuid.UUID) (*models.StudySession, error) {
	var rawSessions []DBStudySession
	query := r.sqlQueries["get_user_active_study_session"]
	err := r.pgclient.QuerySelect(ctx, &rawSessions, query, userID.String())
	if err != nil {
		return nil, err
	}
	if len(rawSessions) == 0 {
		return nil, nil
	}
	return rawSessions[0].ToStudySession()
}

func (r *studySessionRepository) UpsertActiveStudySession(ctx context.Context, userID uuid.UUID, session models.StudySession) (*models.StudySession, error) {
	tx, err := r.pgclient.BeginTransaction(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var existingSession DBStudySession
	err = tx.GetContext(
		ctx,
		&existingSession,
		"SELECT * FROM study_sessions WHERE user_id = $1 AND session_state = 'active' LIMIT 1",
		userID,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check for existing active session: %w", err)
	}
	sessionExists := err == nil

	var resultSession DBStudySession
	if sessionExists {
		updatedSession, err := r.updateStudySession(tx, DBStudySession{
			ID:           existingSession.ID,
			Title:        session.Title,
			Notes:        session.Notes,
			SessionState: string(session.SessionState),
		})
		if err != nil {
			return nil, err
		}
		resultSession = *updatedSession
	} else {
		newSession, err := r.createStudySession(tx, DBStudySession{
			UserID:       userID.String(),
			Title:        session.Title,
			Notes:        session.Notes,
			Date:         time.Now(),
			SessionState: string(models.SessionStateActive),
		})
		if err != nil {
			return nil, err
		}
		resultSession = *newSession
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return resultSession.ToStudySession()
}

func (r *studySessionRepository) AddSessionEvents(ctx context.Context, sessionID uuid.UUID, events []models.SessionEvent) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := r.pgclient.BeginTransaction(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
        INSERT INTO session_events 
            (session_id, event_type, event_time) 
        VALUES 
    `

	var params []any
	paramCount := 1

	for i, event := range events {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d, $%d)",
			paramCount, paramCount+1, paramCount+2)

		params = append(params,
			sessionID,
			string(event.EventType),
			event.EventTime,
		)
		paramCount += 4
	}

	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("failed to insert session events: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *studySessionRepository) updateStudySession(dbTx *sqlx.Tx, dbSession DBStudySession) (*DBStudySession, error) {
	query := r.sqlQueries["update_study_session"]
	rows, err := dbTx.NamedQuery(query, dbSession)
	if err != nil {
		return nil, fmt.Errorf("failed to update study session: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no rows returned after update")
	}

	if err := rows.StructScan(&dbSession); err != nil {
		return nil, fmt.Errorf("failed to scan updated session: %w", err)
	}
	return &dbSession, nil
}

func (r *studySessionRepository) createStudySession(dbTx *sqlx.Tx, dbSession DBStudySession) (*DBStudySession, error) {
	var resultSession DBStudySession
	query := r.sqlQueries["create_study_session"]
	dbSession.Date = time.Now()
	dbSession.SessionState = string(models.SessionStateActive)
	rows, err := dbTx.NamedQuery(query, dbSession)
	if err != nil {
		return nil, fmt.Errorf("failed to create study session: %w", err)
	}

	if !rows.Next() {
		return nil, fmt.Errorf("no rows returned after insert")
	}

	if err := rows.StructScan(&resultSession); err != nil {
		return nil, fmt.Errorf("failed to scan created session: %w", err)
	}
	rows.Close()

	startEvent := DBSessionEvent{
		SessionID: resultSession.ID,
		EventType: string(models.EventTypeStart),
		EventTime: ptr.Of(time.Now().UTC()),
	}
	query = r.sqlQueries["create_session_event"]
	_, err = dbTx.NamedExec(query, startEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to insert session event: %w", err)
	}
	return &resultSession, nil
}
