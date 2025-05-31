package studysession

import (
	"net/http"

	models "go-api/src/models/studysession"
	service "go-api/src/services/studysession"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// StudySessionHandler defines the interface for study session API handlers
type StudySessionHandler interface {
	StartStudySession(e echo.Context) error
	GetActiveStudySession(e echo.Context) error
	AddStudySessionEvents(e echo.Context) error
	FinishStudySession(e echo.Context) error
	GetActiveStudySessionEvents(e echo.Context) error
}

// StudySessionHandlerParams defines the dependencies for the study session handler
type StudySessionHandlerParams struct {
	fx.In

	Service service.StudySessionService
	Logger  *zap.Logger
}

type studySessionHandler struct {
	service service.StudySessionService
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
//	@Param			request	body		service.UpsertActiveStudySessionRequest	true	"Study session data"
//	@Success		201		{object}	models.StudySession
//	@Failure		400		{object}	map[string]string
//	@Failure		409		{object}	map[string]string	"Active session already exists"
//	@Failure		500		{object}	map[string]string
//	@Router			/study-session/start [post]
func (h *studySessionHandler) StartStudySession(e echo.Context) error {
	var req service.UpsertActiveStudySessionRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	ctx := e.Request().Context()
	studySession, err := h.service.CreateStudySession(ctx, req)
	if err != nil {
		switch err {
		case models.ErrActiveSessionExists:
			return e.JSON(http.StatusConflict, map[string]string{"error": "Active session already exists"})
		default:
			h.logger.Error("Failed to create study session", zap.Error(err))
			return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create study session"})
		}
	}

	return e.JSON(http.StatusCreated, studySession)
}

// GetActiveStudySession handles retrieving the active study session
//
//	@Summary		Get active study session
//	@Description	Get the user's active study session
//	@Tags			study-session
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	models.StudySession
//	@Failure		404	{object}	map[string]string	"No active session found"
//	@Failure		500	{object}	map[string]string
//	@Router			/study-session [get]
func (h *studySessionHandler) GetActiveStudySession(e echo.Context) error {
	ctx := e.Request().Context()
	studySession, err := h.service.GetActiveStudySession(ctx)
	if err != nil {
		switch err {
		case models.ErrActiveSessionNotFound:
			return e.JSON(http.StatusNotFound, map[string]string{"error": "No active session found"})
		default:
			h.logger.Error("Failed to get active study session", zap.Error(err))
			return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get active study session"})
		}
	}
	return e.JSON(http.StatusOK, studySession)
}

// AddStudySessionEvents handles adding events to the active study session
//
//	@Summary		Add events to active study session
//	@Description	Add events to the user's active study session
//	@Tags			study-session
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		service.AddStudySessionEventsRequest	true	"Session events data"
//	@Success		200		{object}	models.StudySession
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string	"No active session found"
//	@Failure		422		{object}	map[string]string	"Session not active"
//	@Failure		500		{object}	map[string]string
//	@Router			/study-session/events [post]
func (h *studySessionHandler) AddStudySessionEvents(e echo.Context) error {
	var req service.AddStudySessionEventsRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	if len(req.Events) == 0 {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "No events provided"})
	}

	ctx := e.Request().Context()
	studySession, err := h.service.AddStudySessionEvents(ctx, req)
	if err != nil {
		switch err {
		case models.ErrActiveSessionNotFound:
			return e.JSON(http.StatusNotFound, map[string]string{"error": "No active session found"})
		default:
			h.logger.Error("Failed to add study session events", zap.Error(err))
			return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add study session events"})
		}
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
//	@Param			request	body		service.FinishStudySessionRequest	true	"Finish session data"
//	@Success		200		{object}	models.StudySession
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string	"No active session found"
//	@Failure		422		{object}	map[string]string	"Session not active"
//	@Failure		500		{object}	map[string]string
//	@Router			/study-session/finish [post]
func (h *studySessionHandler) FinishStudySession(e echo.Context) error {
	var req service.FinishStudySessionRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	ctx := e.Request().Context()
	studySession, err := h.service.FinishStudySession(ctx, req)
	if err != nil {
		switch err {
		case models.ErrActiveSessionNotFound:
			return e.JSON(http.StatusNotFound, map[string]string{"error": "No active session found"})
		default:
			h.logger.Error("Failed to finish study session", zap.Error(err))
			return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to finish study session"})
		}
	}

	return e.JSON(http.StatusOK, studySession)
}

// GetActiveStudySessionEvents handles retrieving events for the active study session
//
//	@Summary		Get active study session events
//	@Description	Get events for the user's active study session
//	@Tags			study-session
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	[]models.SessionEvent
//	@Failure		404	{object}	map[string]string	"No active session found"
//	@Failure		500	{object}	map[string]string
//	@Router			/study-session/events [get]
func (h *studySessionHandler) GetActiveStudySessionEvents(e echo.Context) error {
	ctx := e.Request().Context()

	events, err := h.service.GetActiveStudySessionEvents(ctx)
	if err != nil {
		switch err {
		case models.ErrActiveSessionNotFound:
			return e.JSON(http.StatusNotFound, map[string]string{"error": "No active session found"})
		default:
			h.logger.Error("Failed to get study session events",
				zap.Error(err),
				zap.String("endpoint", "/study-session/events"),
			)
			return e.JSON(http.StatusInternalServerError,
				map[string]string{"error": "Failed to get study session events"})
		}
	}

	return e.JSON(http.StatusOK, events)
}
