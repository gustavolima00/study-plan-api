package studysession

import (
	"context"
	"go-api/src/repositories/studysession"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type StudySessionService interface {
	Sample()
}

type studySessionService struct {
	repository studysession.StudySessionRepository
	logger     *zap.Logger
}

type StudySessionServiceParams struct {
	fx.In

	Repository studysession.StudySessionRepository
	Logger     *zap.Logger
}

func NewStudySessionService(p StudySessionServiceParams) StudySessionService {
	return &studySessionService{
		repository: p.Repository,
		logger:     p.Logger,
	}
}

func (s studySessionService) Sample() {
	ctx := context.TODO()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	sessionID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	studySessions, err := s.repository.GetUserStudySessions(ctx, userID)
	s.logger.Debug("result", zap.Any("studySessions", studySessions), zap.Error(err))
	studySession, err := s.repository.GetUserStudySession(ctx, sessionID)
	s.logger.Debug("result", zap.Any("studySession", studySession), zap.Error(err))
}
