package studysession

import (
	"context"
	"embed"
	"go-api/src/gateways/postgres"
	models "go-api/src/models/database/studysession"
	"go-api/src/repositories/util"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

const (
	// SQL Queries
	getUserStudySessionsQueryKey = "get_user_study_sessions"
	getStudySessionEventsKey     = "get_study_session_events"
	getStudySessionSubjectsKey   = "get_study_session_subjects"
)

type StudySessionRepository interface {
	GetUserStudySessions(ctx context.Context, userID uuid.UUID) ([]models.StudySession, error)
	GetStudySessionEvents(ctx context.Context, studySessionID uuid.UUID) ([]models.SessionEvent, error)
	GetStudySessionSubjects(ctx context.Context, studySessionID uuid.UUID) ([]models.Subject, error)
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
	queries, err := util.LoadSQLQueries(sqlFiles)
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
	return util.DBQuery[models.StudySession](ctx, util.DBQueryParams{
		DBConnection: r.pgclient.GetConnection(),
		SqlQuery:     r.sqlQueries[getUserStudySessionsQueryKey],
		Variables: map[string]any{
			"user_id": userID.String(),
		}},
	)
}

func (r *studySessionRepository) GetStudySessionEvents(ctx context.Context, studySessionID uuid.UUID) ([]models.SessionEvent, error) {
	return util.DBQuery[models.SessionEvent](ctx, util.DBQueryParams{
		DBConnection: r.pgclient.GetConnection(),
		SqlQuery:     r.sqlQueries[getStudySessionEventsKey],
		Variables: map[string]any{
			"session_id": studySessionID.String(),
		}},
	)
}

func (r *studySessionRepository) GetStudySessionSubjects(ctx context.Context, studySessionID uuid.UUID) ([]models.Subject, error) {
	return util.DBQuery[models.Subject](ctx, util.DBQueryParams{
		DBConnection: r.pgclient.GetConnection(),
		SqlQuery:     r.sqlQueries[getStudySessionSubjectsKey],
		Variables: map[string]any{
			"session_id": studySessionID.String(),
		}},
	)
}
