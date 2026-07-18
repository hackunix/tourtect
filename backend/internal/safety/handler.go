package safety

import (
	"encoding/json"
	"net/http"

	"github.com/tourtect/backend/generated/openapi"
	"github.com/tourtect/backend/internal/platform/httpserver"
)

type Handler struct {
	engine *Engine
}

func NewHandler(engine *Engine) *Handler {
	return &Handler{engine: engine}
}

func (h *Handler) CreateSafetyAssessment(w http.ResponseWriter, r *http.Request, params openapi.CreateSafetyAssessmentParams) {
	ctx := r.Context()
	reqID := httpserver.GetRequestID(ctx)

	var req openapi.SafetyAssessmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpserver.WriteError(w, http.StatusUnprocessableEntity, "Unprocessable Entity", "Invalid request body format", r.URL.Path, reqID)
		return
	}

	var state string
	if req.UserReportedState != nil {
		state = *req.UserReportedState
	}

	var threat []string
	if req.ThreatIndicators != nil {
		threat = *req.ThreatIndicators
	}

	var injury []string
	if req.InjuryIndicators != nil {
		injury = *req.InjuryIndicators
	}

	var confinement []string
	if req.ConfinementIndicators != nil {
		confinement = *req.ConfinementIndicators
	}

	var coercion []string
	if req.CoercionIndicators != nil {
		coercion = *req.CoercionIndicators
	}

	var facts []string
	if req.UserConfirmedFacts != nil {
		facts = *req.UserConfirmedFacts
	}

	input := AssessmentInput{
		ObservedFacts:         req.ObservedFacts,
		UserReportedState:     state,
		ThreatIndicators:      threat,
		InjuryIndicators:      injury,
		ConfinementIndicators: confinement,
		CoercionIndicators:    coercion,
		AbilityToLeave:        req.AbilityToLeave,
		UserConfirmedFacts:    facts,
		RegionID:              "hanoi", // default region matching seed
	}

	res, err := h.engine.Assess(ctx, input, reqID)
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, "Internal Server Error", err.Error(), r.URL.Path, reqID)
		return
	}

	// Format response to openapi.SafetyAssessment
	resp := openapi.SafetyAssessment{
		Urgency:                 openapi.SafetyUrgency(res.Urgency),
		SafeActions:            res.SafeActions,
		ApprovedActionCodes:    &res.ApprovedActionCodes,
		ExplanationCodes:       &res.ExplanationCodes,
		SilentModeRecommended:   &res.SilentModeRecommended,
		SurfaceEmergencyOptions: &res.SurfaceEmergencyOptions,
		EmergencyServiceIds:     &res.EmergencyServiceIDs,
		SafetyDirectoryVersion:  res.SafetyDirectoryVersion,
		Confidence:              res.Confidence,
		TraceId:                 res.TraceID,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
