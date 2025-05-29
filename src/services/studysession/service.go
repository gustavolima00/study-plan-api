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
	UpsertActiveStudySession(ctx context.Context, request models.UpsertActiveStudySessionRequest) (*models.UpsertActiveStudySessionResponse, error)
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

func (s studySessionService) UpsertActiveStudySession(ctx context.Context, request models.UpsertActiveStudySessionRequest) (*models.UpsertActiveStudySessionResponse, error) {
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
	return &models.UpsertActiveStudySessionResponse{
		ID:           studySession.ID,
		Title:        studySession.Title,
		Notes:        studySession.Notes,
		Date:         studySession.Date,
		SessionState: models.SessionState(studySession.SessionState),
	}, nil
}
