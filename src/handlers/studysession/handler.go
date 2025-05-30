package studysession

import (
	"go-api/src/services/studysession"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// StudySessionHandler defines the interface for study session API handlers
type StudySessionHandler interface {
	StartStudySession(e echo.Context) error
	AddStudySessionEvents(e echo.Context) error
	FinishStudySession(e echo.Context) error
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

// StartStudySession handles the creation of a new study session
//
//	@Summary		Create a study session
//	@Description	Create a new study session for the authenticated user
//	@Tags			study-session
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.UpsertActiveStudySessionRequest	true	"Study session data"
//	@Success		201		{object}	models.StudySession	"Session details"
//	@Failure		400		{string}	string	"Invalid request"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/study-session/start [post]
func (h *studySessionHandler) StartStudySession(e echo.Context) error {
	var req studysession.UpsertActiveStudySessionRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid request format")
	}
	ctx := e.Request().Context()

	studySession, err := h.service.CreateStudySession(ctx, req)
	if err != nil {
		h.logger.Error("Failed to create study session", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, "Failed to create study session")
	}

	return e.JSON(http.StatusCreated, studySession)
}

// AddStudySessionEvents handles adding events to the active study session
//
//	@Summary		Add events to active study session
//	@Description	Add events to the user's active study session
//	@Tags			study-session
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.AddStudySessionEventsRequest	true	"Session events data"
//	@Success		200		{object}	models.StudySession	"Session details"
//	@Failure		400		{string}	string	"Invalid request"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/study-session/add-events [post]
func (h *studySessionHandler) AddStudySessionEvents(e echo.Context) error {
	var req studysession.AddStudySessionEventsRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid request format")
	}
	ctx := e.Request().Context()

	studySession, err := h.service.AddStudySessionEvents(ctx, req)
	if err != nil {
		h.logger.Error("Failed to add study session events", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, "Failed to add study session events")
	}

	return e.JSON(http.StatusOK, studySession)
}

// FinishStudySession handles finishing the active study session
//
//	@Summary		Finish active study session
//	@Description	Finish the user's active study session
//	@Tags			study-session
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.FinishStudySessionRequest	true	"Finish session data"
//	@Success		200		{object}	models.StudySession	"Session details"
//	@Failure		400		{string}	string	"Invalid request"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/study-session/finish [post]
func (h *studySessionHandler) FinishStudySession(e echo.Context) error {
	var req studysession.FinishStudySessionRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid request format")
	}
	ctx := e.Request().Context()

	studySession, err := h.service.FinishStudySession(ctx, req)
	if err != nil {
		h.logger.Error("Failed to finish study session", zap.Error(err))
		return e.JSON(http.StatusInternalServerError, "Failed to finish study session")
	}

	return e.JSON(http.StatusOK, studySession)
}
