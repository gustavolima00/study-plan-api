package studysession

import (
	"context"
	models "go-api/src/models/studysession"
	repository "go-api/src/repositories/studysession"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type StudySessionService interface {
	Sample()
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

func (s studySessionService) Sample() {
	ctx := context.TODO()
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	studySessions, err := s.repository.GetUserStudySessions(ctx, userID)
	s.logger.Debug("result", zap.Any("studySessions", studySessions), zap.Error(err))
	studySession, err := s.repository.GetUserActiveStudySession(ctx, userID)
	s.logger.Debug("result", zap.Any("studySession", studySession), zap.Error(err))

	newSession, err := s.repository.CreateOrUpdateUserStudySession(ctx, userID, models.StudySession{
		Date:         time.Now(),
		Notes:        "Some note",
		Title:        "Some title",
		SessionState: models.SessionState("active"),
	})
	s.logger.Debug("result", zap.Any("newSession", newSession), zap.Error(err))
}
