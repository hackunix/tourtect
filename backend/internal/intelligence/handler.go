package intelligence

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tourtect/backend/generated/openapi"
	"github.com/tourtect/backend/internal/intelligence/feedback"
	"github.com/tourtect/backend/internal/intelligence/model"
	"github.com/tourtect/backend/internal/intelligence/orchestrator"
	"github.com/tourtect/backend/internal/intelligence/session"
	"github.com/tourtect/backend/internal/platform/httpserver"
)

type Handler struct {
	sessions     *session.Service
	orchestrator *orchestrator.Orchestrator
	feedback     *feedback.Repository
	now          func() time.Time
}

func NewHandler(sessions *session.Service, orchestrator *orchestrator.Orchestrator, feedbackRepo *feedback.Repository) *Handler {
	return &Handler{sessions: sessions, orchestrator: orchestrator, feedback: feedbackRepo, now: time.Now}
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req openapi.CreateAssistantSessionRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return
	}
	userID := httpserver.GetUserID(r.Context())
	mode := "text"
	if req.InteractionMode != nil {
		mode = string(*req.InteractionMode)
	}
	consent := req.ProcessingConsent != nil && *req.ProcessingConsent
	value, err := h.sessions.Create(r.Context(), userID, req.Locale, valueOf(req.TargetLocale), uuidValue(req.PlaceId), valueOf(req.ApproximateRegion), mode, consent)
	if err != nil {
		httpserver.WriteError(w, http.StatusServiceUnavailable, "Assistant session unavailable", err.Error(), r.URL.Path, httpserver.GetRequestID(r.Context()))
		return
	}
	writeJSON(w, http.StatusCreated, publicSession(value))
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	value, err := h.sessions.GetOwned(r.Context(), id.String(), httpserver.GetUserID(r.Context()))
	if err != nil {
		writeSessionError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, publicSession(value))
}
func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	if err := h.sessions.DeleteOwned(r.Context(), id.String(), httpserver.GetUserID(r.Context())); err != nil {
		writeSessionError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	unlock := h.sessions.Lock(id.String())
	defer unlock()
	value, err := h.sessions.GetOwned(r.Context(), id.String(), httpserver.GetUserID(r.Context()))
	if err != nil {
		writeSessionError(w, r, err)
		return
	}
	var req openapi.AssistantMessageRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return
	}
	structured := json.RawMessage(nil)
	if req.StructuredData != nil {
		structured, _ = json.Marshal(req.StructuredData)
	}
	input := model.Message{ID: req.MessageId.String(), InputType: string(req.InputType), Text: valueOf(req.Text), Locale: valueOf(req.Locale), PlaceID: uuidValue(req.PlaceId), CaptureID: uuidValue(req.CaptureId), UserConfirmed: req.UserConfirmed != nil && *req.UserConfirmed, Structured: structured}
	if err := orchestrator.ValidateMessage(input); err != nil {
		httpserver.WriteError(w, http.StatusUnprocessableEntity, "Invalid assistant message", err.Error(), r.URL.Path, httpserver.GetRequestID(r.Context()))
		return
	}
	resp, trace, err := h.orchestrator.Handle(r.Context(), value, input, httpserver.GetRequestID(r.Context()))
	if errors.Is(err, orchestrator.ErrDuplicateMessage) {
		httpserver.WriteError(w, http.StatusConflict, "Duplicate message", "message_id has already been processed", r.URL.Path, httpserver.GetRequestID(r.Context()))
		return
	}
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, "Assistant orchestration failed", err.Error(), r.URL.Path, httpserver.GetRequestID(r.Context()))
		return
	}
	if principal, parseErr := uuid.Parse(httpserver.GetUserID(r.Context())); parseErr == nil {
		if traceErr := h.feedback.CreateTrace(r.Context(), principal, trace); traceErr != nil {
			slog.Warn("failed to store redacted assistant trace", slog.String("trace_id", trace.TraceID), slog.Any("error", traceErr))
		}
	}
	writeJSON(w, http.StatusOK, publicResponse(resp))
}

func (h *Handler) Confirm(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	unlock := h.sessions.Lock(id.String())
	defer unlock()
	value, err := h.sessions.GetOwned(r.Context(), id.String(), httpserver.GetUserID(r.Context()))
	if err != nil {
		writeSessionError(w, r, err)
		return
	}
	var req openapi.AssistantConfirmationRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return
	}
	confirmation, ok := value.Confirmations[req.ConfirmationId.String()]
	if !ok || confirmation.Consumed || h.now().After(confirmation.ExpiresAt) {
		httpserver.WriteError(w, http.StatusConflict, "Confirmation unavailable", "confirmation expired, consumed, or does not belong to this session", r.URL.Path, httpserver.GetRequestID(r.Context()))
		return
	}
	decision := string(req.Decision)
	if decision != "confirmed" && decision != "rejected" {
		httpserver.WriteError(w, http.StatusUnprocessableEntity, "Invalid confirmation", "decision must be confirmed or rejected", r.URL.Path, httpserver.GetRequestID(r.Context()))
		return
	}
	confirmation.Consumed = true
	var resultID *uuid.UUID
	if decision == "confirmed" {
		created := uuid.New()
		resultID = &created
	}
	principalID, _ := uuid.Parse(httpserver.GetUserID(r.Context()))
	if err := h.feedback.AuditConfirmation(r.Context(), principalID, id, req.ConfirmationId, confirmation.Action, decision, resultID); err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, "Confirmation audit failed", err.Error(), r.URL.Path, httpserver.GetRequestID(r.Context()))
		return
	}
	if err := h.sessions.Save(r.Context(), value); err != nil {
		httpserver.WriteError(w, http.StatusServiceUnavailable, "Session unavailable", err.Error(), r.URL.Path, httpserver.GetRequestID(r.Context()))
		return
	}
	result := map[string]any{"confirmation_id": confirmation.ID, "action": confirmation.Action, "status": decision, "executed_at": h.now().UTC()}
	if resultID != nil {
		result["result_id"] = resultID.String()
		result["target"] = confirmation.Target
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) Feedback(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	value, err := h.sessions.GetOwned(r.Context(), id.String(), httpserver.GetUserID(r.Context()))
	if err != nil {
		writeSessionError(w, r, err)
		return
	}
	var req openapi.AssistantFeedbackRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return
	}
	var original *model.Response
	for i := range value.RecentResponses {
		if value.RecentResponses[i].ID == req.AssistantMessageId.String() {
			original = &value.RecentResponses[i]
			break
		}
	}
	if original == nil {
		httpserver.WriteError(w, http.StatusNotFound, "Assistant message not found", "feedback must reference a response in this session", r.URL.Path, httpserver.GetRequestID(r.Context()))
		return
	}
	consent := req.ConsentToContribute != nil && *req.ConsentToContribute
	if string(req.FeedbackType) == "contribute_redacted_observation" && !consent {
		httpserver.WriteError(w, http.StatusUnprocessableEntity, "Contribution consent required", "contribution feedback requires explicit separate consent", r.URL.Path, httpserver.GetRequestID(r.Context()))
		return
	}
	principalID, _ := uuid.Parse(httpserver.GetUserID(r.Context()))
	feedbackID, status, createdAt, err := h.feedback.Create(r.Context(), feedback.FeedbackInput{PrincipalID: principalID, SessionID: id, AssistantMessageID: req.AssistantMessageId, FeedbackType: string(req.FeedbackType), Field: req.Field, OriginalValue: req.OriginalValue, CorrectedValue: req.CorrectedValue, ConsentToContribute: consent, OriginalResponse: *original})
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, "Feedback unavailable", err.Error(), r.URL.Path, httpserver.GetRequestID(r.Context()))
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"feedback_id": feedbackID, "status": status, "created_at": createdAt})
}

func publicSession(value *model.Session) any {
	return map[string]any{"session_id": value.ID, "version": value.Version, "created_at": value.CreatedAt, "updated_at": value.UpdatedAt, "expires_at": value.ExpiresAt, "context": value.Context, "recent_responses": value.RecentResponses}
}
func publicResponse(value model.Response) any {
	var result openapi.AssistantResponse
	b, _ := json.Marshal(value)
	if err := json.Unmarshal(b, &result); err != nil {
		return value
	}
	return result
}
func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256*1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		httpserver.WriteError(w, http.StatusUnprocessableEntity, "Invalid request body", err.Error(), r.URL.Path, httpserver.GetRequestID(r.Context()))
		return err
	}
	return nil
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeSessionError(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusServiceUnavailable
	title := "Assistant session unavailable"
	if errors.Is(err, session.ErrNotFound) {
		status = http.StatusNotFound
		title = "Assistant session not found"
	}
	httpserver.WriteError(w, status, title, err.Error(), r.URL.Path, httpserver.GetRequestID(r.Context()))
}
func valueOf(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
func uuidValue(value *uuid.UUID) string {
	if value == nil {
		return ""
	}
	return value.String()
}
