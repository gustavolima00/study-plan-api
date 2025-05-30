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
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

type StudySessionRepository interface {
	UpsertActiveStudySession(ctx context.Context, userID uuid.UUID, session models.StudySession) (*models.StudySession, error)
	GetUserActiveStudySession(ctx context.Context, userID uuid.UUID) (*models.StudySession, error)
	GetStudySessionByID(ctx context.Context, sessionID uuid.UUID) (*models.StudySession, error)
	AddSessionEvents(ctx context.Context, sessionID uuid.UUID, events []models.SessionEvent) (*models.StudySession, error)
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
		"SELECT * FROM study_sessions WHERE user_id = $1 AND session_state = $2 LIMIT 1",
		userID,
		string(models.SessionStateActive),
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

	var dbEvents []DBSessionEvent
	err = tx.GetContext(
		ctx,
		&dbEvents,
		"SELECT * FROM session_events WHERE session_id = $1 ORDER BY event_time",
		resultSession.ID,
	)

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return resultSession.ToStudySession(dbEvents)
}

func (r *studySessionRepository) GetUserActiveStudySession(ctx context.Context, userID uuid.UUID) (*models.StudySession, error) {
	tx, err := r.pgclient.BeginTransaction(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var dbSession DBStudySession
	err = tx.GetContext(
		ctx,
		&dbSession,
		"SELECT * FROM study_sessions WHERE user_id = $1 AND session_state = $2 LIMIT 1",
		userID,
		string(models.SessionStateActive),
	)

	var dbEvents []DBSessionEvent
	err = tx.GetContext(
		ctx,
		&dbEvents,
		"SELECT * FROM session_events WHERE session_id = $1 ORDER BY event_time",
		dbSession.ID,
	)

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return dbSession.ToStudySession(dbEvents)
}

func (r *studySessionRepository) GetStudySessionByID(ctx context.Context, sessionID uuid.UUID) (*models.StudySession, error) {
	tx, err := r.pgclient.BeginTransaction(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	studySession, err := r.getStudySessionByID(ctx, tx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get study session: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return studySession, nil
}

func (r *studySessionRepository) AddSessionEvents(ctx context.Context, sessionID uuid.UUID, events []models.SessionEvent) (*models.StudySession, error) {
	tx, err := r.pgclient.BeginTransaction(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
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
	if len(events) != 0 {
		_, err = tx.ExecContext(ctx, query, params...)
		if err != nil {
			return nil, fmt.Errorf("failed to insert session events: %w", err)
		}
	}

	studySession, err := r.getStudySessionByID(ctx, tx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get study session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return studySession, nil
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
		EventTime: time.Now().UTC(),
	}
	query = r.sqlQueries["create_session_event"]
	_, err = dbTx.NamedExec(query, startEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to insert session event: %w", err)
	}
	return &resultSession, nil
}

func (r *studySessionRepository) getStudySessionByID(ctx context.Context, dbTx *sqlx.Tx, sessionID uuid.UUID) (*models.StudySession, error) {
	var dbSession DBStudySession
	err := dbTx.GetContext(
		ctx,
		&dbSession,
		"SELECT * FROM study_sessions WHERE id = $1 LIMIT 1",
		sessionID,
	)
	if err != nil {
		return nil, err
	}

	var dbEvents []DBSessionEvent
	err = dbTx.GetContext(
		ctx,
		&dbEvents,
		"SELECT * FROM session_events WHERE session_id = $1 ORDER BY event_time",
		dbSession.ID,
	)
	if err != nil {
		return nil, err
	}
	return dbSession.ToStudySession(dbEvents)
}
