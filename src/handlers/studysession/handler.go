package studysession

import (
	models "go-api/src/models/studysession"
	"go-api/src/services/studysession"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// StudySessionHandler defines the interface for study session API handlers
type StudySessionHandler interface {
	UpsertActiveStudySession(e echo.Context) error
}

// StudySessionHandlerParams defines the dependencies for the study session handler
type StudySessionHandlerParams struct {
	fx.In

	Service studysession.StudySessionService
	Logger  *zap.Logger
}

type studySessionHandler struct {
	service studysession.StudySessionService
	logger  *zap.Logger
}

// NewStudySessionHandler creates a new study session handler with injected dependencies
func NewStudySessionHandler(p StudySessionHandlerParams) StudySessionHandler {
	return &studySessionHandler{
		service: p.Service,
		logger:  p.Logger,
	}
}

// UpsertActiveStudySession handles the creation of a new study session
//
//	@Summary		Create a study session
//	@Description	Create a new study session for the authenticated user
//	@Tags			study-session
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.UpsertActiveStudySessionRequest	true	"Study session data"
//	@Success		200		{string}	models.UpsertActiveStudySessionResponse	"Session details"
//	@Failure		400		{string}	string	"Invalid request"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/study-session/upsert-active [post]
func (h *studySessionHandler) UpsertActiveStudySession(e echo.Context) error {
	var req models.UpsertActiveStudySessionRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid request format")
	}
	ctx := e.Request().Context()

	resp, err := h.service.UpsertActiveStudySession(ctx, req)
	if err != nil {
		h.logger.Error("Failed to create study session", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, "Failed to create study session")
	}

	return e.JSON(http.StatusOK, resp)
}
