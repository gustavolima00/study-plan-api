package studysession

import (
	"context"
	"fmt"
	authmodel "go-api/src/models/auth"
	"go-api/src/models/constants"
	models "go-api/src/models/studysession"
	repository "go-api/src/repositories/studysession"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type StudySessionService interface {
	CreateStudySession(ctx context.Context, request UpsertActiveStudySessionRequest) (*models.StudySession, error)
	AddStudySessionEvents(ctx context.Context, request AddStudySessionEventsRequest) (*models.StudySession, error)
	FinishStudySession(ctx context.Context, request FinishStudySessionRequest) (*models.StudySession, error)
}

type studySessionService struct {
	repository repository.StudySessionRepository
	logger     *zap.Logger
}

type StudySessionServiceParams struct {
	fx.In

	Repository repository.StudySessionRepository
	Logger     *zap.Logger
}

func NewStudySessionService(p StudySessionServiceParams) StudySessionService {
	return &studySessionService{
		repository: p.Repository,
		logger:     p.Logger,
	}
}

func (s studySessionService) CreateStudySession(ctx context.Context, request UpsertActiveStudySessionRequest) (*models.StudySession, error) {
	user := ctx.Value(constants.ContextKeyUserInfoKey).(*authmodel.UserInfo)
	if user == nil {
		s.logger.Error("Failed to create studySession, no user found in context")
		return nil, fmt.Errorf("no user found in context")
	}

	logger := s.logger.With(
		zap.Stringer("user_id", user.ID),
	)
	studySession, err := s.repository.UpsertActiveStudySession(ctx, user.ID, models.StudySession{
		Date:         time.Now(),
		Notes:        request.Notes,
		Title:        request.Title,
		SessionState: models.SessionState("active"),
	})
	if err != nil {
		logger.Error("Failed to create studySession", zap.Error(err))
		return nil, err
	}
	return studySession, nil
}

func (s studySessionService) AddStudySessionEvents(ctx context.Context, request AddStudySessionEventsRequest) (*models.StudySession, error) {
	user := ctx.Value(constants.ContextKeyUserInfoKey).(*authmodel.UserInfo)
	if user == nil {
		s.logger.Error("Failed to create studySession, no user found in context")
		return nil, fmt.Errorf("no user found in context")
	}

	logger := s.logger.With(
		zap.Stringer("user_id", user.ID),
	)
	activeStudySession, err := s.repository.GetUserActiveStudySession(ctx, user.ID)
	if err != nil {
		logger.Error("Failed to get active studySession", zap.Error(err))
		return nil, err
	}
	studySession, err := s.repository.AddSessionEvents(ctx, activeStudySession.UserID, request.Events)
	if err != nil {
		logger.Error("Failed create session events", zap.Error(err))
		return nil, err
	}
	return studySession, nil
}

func (s studySessionService) FinishStudySession(ctx context.Context, request FinishStudySessionRequest) (*models.StudySession, error) {
	user := ctx.Value(constants.ContextKeyUserInfoKey).(*authmodel.UserInfo)
	if user == nil {
		s.logger.Error("Failed to create studySession, no user found in context")
		return nil, fmt.Errorf("no user found in context")
	}
	logger := s.logger.With(
		zap.Stringer("user_id", user.ID),
	)
	activeStudySession, err := s.repository.GetUserActiveStudySession(ctx, user.ID)
	if err != nil {
		logger.Error("Failed to get active studySession", zap.Error(err))
		return nil, err
	}
	studySession, err := s.repository.AddSessionEvents(ctx, activeStudySession.UserID, []models.SessionEvent{
		{
			EventType: models.EventTypeStop,
			EventTime: request.FinishedAt,
		},
	})
	if err != nil {
		logger.Error("Failed add session finish event", zap.Error(err))
		return nil, err
	}
	studySession.SessionState = models.SessionStateCompleted
	studySession, err = s.repository.UpsertActiveStudySession(ctx, activeStudySession.UserID, *studySession)
	if err != nil {
		logger.Error("Failed to upsert finished session", zap.Error(err))
		return nil, err
	}
	return studySession, nil
}
