package studysession

import (
	"context"
	"fmt"
	authmodel "go-api/src/models/auth"
	"go-api/src/models/constants"
	models "go-api/src/models/studysession"
	repository "go-api/src/repositories/studysession"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type StudySessionService interface {
	CreateStudySession(ctx context.Context, request UpsertActiveStudySessionRequest) (*models.StudySession, error)
	GetActiveStudySession(ctx context.Context) (*models.StudySession, error)
	GetActiveStudySessionEvents(ctx context.Context) ([]models.SessionEvent, error)
	AddStudySessionEvents(ctx context.Context, request AddStudySessionEventsRequest) ([]models.SessionEvent, error)
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
	return s.repository.CreateStudySession(
		ctx,
		models.StudySession{
			Notes:  request.Notes,
			Title:  request.Title,
			UserID: user.ID,
		},
		request.StartedAt,
	)
}

func (s studySessionService) AddStudySessionEvents(ctx context.Context, request AddStudySessionEventsRequest) ([]models.SessionEvent, error) {
	user := ctx.Value(constants.ContextKeyUserInfoKey).(*authmodel.UserInfo)
	if user == nil {
		s.logger.Error("Failed to create studySession, no user found in context")
		return nil, fmt.Errorf("no user found in context")
	}
	return s.repository.AddActiveStudySessionEvents(ctx, user.ID, request.Events)
}

func (s studySessionService) FinishStudySession(ctx context.Context, request FinishStudySessionRequest) (*models.StudySession, error) {
	user := ctx.Value(constants.ContextKeyUserInfoKey).(*authmodel.UserInfo)
	if user == nil {
		s.logger.Error("Failed to create studySession, no user found in context")
		return nil, fmt.Errorf("no user found in context")
	}
	return s.repository.FinishActiveStudySession(ctx, user.ID)
}

func (s studySessionService) GetActiveStudySession(ctx context.Context) (*models.StudySession, error) {
	user := ctx.Value(constants.ContextKeyUserInfoKey).(*authmodel.UserInfo)
	if user == nil {
		s.logger.Error("Failed to create studySession, no user found in context")
		return nil, fmt.Errorf("no user found in context")
	}
	return s.repository.GetActiveStudySession(ctx, user.ID)
}
func (s studySessionService) GetActiveStudySessionEvents(ctx context.Context) ([]models.SessionEvent, error) {
	user := ctx.Value(constants.ContextKeyUserInfoKey).(*authmodel.UserInfo)
	if user == nil {
		s.logger.Error("Failed to create studySession, no user found in context")
		return nil, fmt.Errorf("no user found in context")
	}
	return s.repository.GetActiveStudySessionEvents(ctx, user.ID)
}
