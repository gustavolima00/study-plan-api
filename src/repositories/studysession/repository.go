package studysession

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-api/src/clients/postgres"
	models "go-api/src/models/studysession"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type StudySessionRepository interface {
	CreateStudySession(ctx context.Context, session models.StudySession, startTime time.Time) (*models.StudySession, error)
	GetActiveStudySession(ctx context.Context, userID uuid.UUID) (*models.StudySession, error)
	GetActiveStudySessionEvents(ctx context.Context, userID uuid.UUID) ([]models.SessionEvent, error)
	AddActiveStudySessionEvents(ctx context.Context, userID uuid.UUID, events []models.SessionEvent) ([]models.SessionEvent, error)
	FinishActiveStudySession(ctx context.Context, userID uuid.UUID) (*models.StudySession, error)
}

type studySessionRepository struct {
	logger   *zap.Logger
	pgclient postgres.PostgresClient
}

type StudySessionRepositoryParams struct {
	fx.In

	Logger   *zap.Logger
	PGClient postgres.PostgresClient
}

func NewStudySessionRepository(p StudySessionRepositoryParams) (StudySessionRepository, error) {
	return &studySessionRepository{
		logger:   p.Logger,
		pgclient: p.PGClient,
	}, nil
}

func (r *studySessionRepository) CreateStudySession(ctx context.Context, session models.StudySession, startTime time.Time) (*models.StudySession, error) {
	tx, err := r.beginTransaction(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.safeRollback(err)

	existingActiveSession, err := tx.getUserActiveSession(ctx, session.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to check active sessions: %w", err)
	}
	if existingActiveSession != nil {
		return nil, models.ErrActiveSessionExists
	}

	sessionID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to create session uuid: %w", err)
	}
	dbSession := DBStudySession{
		ID:           sessionID.String(),
		UserID:       session.UserID.String(),
		Title:        session.Title,
		Notes:        session.Notes,
		Date:         startTime.UTC(),
		SessionState: string(models.SessionStateActive),
	}

	query := `INSERT INTO 
				study_sessions (id, user_id, title, notes, date, session_state)
				VALUES (:id, :user_id, :title, :notes, :date, :session_state)`
	_, err = tx.NamedExecContext(ctx, query, dbSession)
	if err != nil {
		return nil, fmt.Errorf("failed to create study session: %w", err)
	}

	err = tx.createSessionEvents(ctx, []DBSessionEvent{{
		SessionID: dbSession.ID,
		EventType: string(models.EventTypeStart),
		EventTime: startTime,
	}})
	if err != nil {
		return nil, fmt.Errorf("failed to create start event: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}
	return dbSession.ToStudySession()
}

func (r *studySessionRepository) AddActiveStudySessionEvents(ctx context.Context, userID uuid.UUID, events []models.SessionEvent) ([]models.SessionEvent, error) {
	tx, err := r.beginTransaction(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.safeRollback(err)

	activeSession, err := tx.getUserActiveSession(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to check active sessions: %w", err)
	}
	if activeSession == nil {
		return nil, models.ErrActiveSessionNotFound
	}
	dbEvents := make([]DBSessionEvent, len(events))
	for i, event := range events {
		dbEvents[i] = DBSessionEvent{
			SessionID: activeSession.ID,
			EventType: string(event.EventType),
			EventTime: event.EventTime.UTC(),
		}
	}

	err = tx.createSessionEvents(ctx, dbEvents)
	if err != nil {
		return nil, fmt.Errorf("failed to insert events: %w", err)
	}

	dbSessionEvents, err := tx.getSessionEvents(ctx, activeSession.ID)
	if err != nil {
		return nil, fmt.Errorf("failed get session events: %w", err)
	}
	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}
	sessionEvents := make([]models.SessionEvent, len(dbSessionEvents))
	for i, dbSessionEvent := range dbSessionEvents {
		sessionEvents[i] = dbSessionEvent.ToSessionEvent()
	}

	return sessionEvents, nil
}

func (r *studySessionRepository) FinishActiveStudySession(ctx context.Context, userID uuid.UUID) (*models.StudySession, error) {
	tx, err := r.beginTransaction(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.safeRollback(err)

	activeSession, err := tx.getUserActiveSession(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to check active sessions: %w", err)
	}
	if activeSession == nil {
		return nil, models.ErrActiveSessionNotFound
	}
	activeSession.SessionState = string(models.SessionStateCompleted)

	_, err = tx.ExecContext(ctx,
		"UPDATE study_sessions SET session_state = $1 WHERE id = $2 AND session_state = $3",
		string(models.SessionStateCompleted),
		activeSession.ID,
		string(models.SessionStateActive),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update session state: %w", err)
	}

	err = tx.createSessionEvents(ctx, []DBSessionEvent{{
		SessionID: activeSession.ID,
		EventType: string(models.EventTypeStop),
		EventTime: time.Now().UTC(),
	}})
	if err != nil {
		return nil, fmt.Errorf("failed to create end event: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}

	return activeSession.ToStudySession()
}

func (r *studySessionRepository) GetActiveStudySession(ctx context.Context, userID uuid.UUID) (*models.StudySession, error) {
	tx, err := r.beginTransaction(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.safeRollback(err)

	activeSession, err := tx.getUserActiveSession(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to check active sessions: %w", err)
	}
	if activeSession == nil {
		return nil, models.ErrActiveSessionNotFound
	}
	return activeSession.ToStudySession()
}

func (r *studySessionRepository) GetActiveStudySessionEvents(ctx context.Context, userID uuid.UUID) ([]models.SessionEvent, error) {
	tx, err := r.beginTransaction(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.safeRollback(err)

	activeSession, err := tx.getUserActiveSession(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to check active sessions: %w", err)
	}
	if activeSession == nil {
		return nil, models.ErrActiveSessionNotFound
	}
	dbEvents, err := tx.getSessionEvents(ctx, activeSession.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session events: %w", err)
	}
	events := make([]models.SessionEvent, len(dbEvents))
	for i, dbEvent := range dbEvents {
		events[i] = dbEvent.ToSessionEvent()
	}
	return events, nil
}

type openTransaction struct {
	sqlx.Tx
}

func (r *studySessionRepository) beginTransaction(ctx context.Context, opts *sql.TxOptions) (*openTransaction, error) {
	tx, err := r.pgclient.BeginTransaction(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &openTransaction{
		Tx: *tx,
	}, nil
}

func (tx openTransaction) getUserActiveSession(ctx context.Context, userID string) (*DBStudySession, error) {
	var activeSession DBStudySession
	err := tx.GetContext(ctx, &activeSession,
		"SELECT * FROM study_sessions WHERE user_id = $1 AND session_state = $2",
		userID, string(models.SessionStateActive),
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &activeSession, nil
}

func (tx openTransaction) getSessionEvents(ctx context.Context, sessionID string) ([]DBSessionEvent, error) {
	var sessionEvents []DBSessionEvent
	err := tx.SelectContext(
		ctx,
		&sessionEvents,
		"SELECT * FROM session_events WHERE session_id = $1",
		sessionID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return []DBSessionEvent{}, nil
	}
	if err != nil {
		return nil, err
	}
	return sessionEvents, nil
}

func (tx openTransaction) createSessionEvents(ctx context.Context, events []DBSessionEvent) error {
	if len(events) == 0 {
		return nil
	}
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
			event.SessionID,
			event.EventType,
			event.EventTime,
		)
		paramCount += 3
	}
	_, err := tx.ExecContext(ctx, query, params...)
	return err
}

func (tx openTransaction) safeRollback(err error) {
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p)
	} else if err != nil {
		tx.Rollback()
	}
}
