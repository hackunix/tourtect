package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tourtect/backend/adapters/fptai"
	"github.com/tourtect/backend/internal/intelligence/evidence"
	"github.com/tourtect/backend/internal/intelligence/intent"
	"github.com/tourtect/backend/internal/intelligence/model"
	"github.com/tourtect/backend/internal/intelligence/response"
	"github.com/tourtect/backend/internal/intelligence/retrieval"
	"github.com/tourtect/backend/internal/intelligence/session"
	"github.com/tourtect/backend/internal/intelligence/tools"
	"github.com/tourtect/backend/internal/pricing"
	"github.com/tourtect/backend/internal/safety"
)

var ErrDuplicateMessage = errors.New("duplicate assistant message id")

type Orchestrator struct {
	sessions  *session.Service
	router    *intent.Router
	registry  *tools.Registry
	retrieval *retrieval.Service
	evidence  *evidence.Assembler
	composer  *response.Composer
	now       func() time.Time
}

func New(sessions *session.Service, router *intent.Router, registry *tools.Registry, retrievalService *retrieval.Service, evidenceAssembler *evidence.Assembler, composer *response.Composer) *Orchestrator {
	return &Orchestrator{sessions: sessions, router: router, registry: registry, retrieval: retrievalService, evidence: evidenceAssembler, composer: composer, now: time.Now}
}

func (o *Orchestrator) Handle(ctx context.Context, value *model.Session, input model.Message, traceID string) (model.Response, model.Trace, error) {
	if o.sessions.HasMessage(value, input.ID) {
		return model.Response{}, model.Trace{}, ErrDuplicateMessage
	}
	if _, err := uuid.Parse(traceID); err != nil {
		traceID = uuid.NewString()
	}
	route := o.router.Route(input.InputType, input.Text)
	trace := model.Trace{TraceID: traceID, SessionID: value.ID, Intent: route.Intent, PolicyVersion: "assistant-policy-2026-07-v1", ToolNames: []string{}, ToolDurationsMS: []int64{}, EvidenceIDs: []string{}}
	base := model.Response{ID: uuid.NewString(), Intent: route.Intent, Confidence: route.Confidence, Evidence: []model.Evidence{}, ToolResults: []model.ToolResult{}, SuggestedActions: []model.SuggestedAction{}, SafetyState: "unknown", FactorsConsidered: []string{}, MissingInformation: route.MissingFields, FallbackUsed: true, TraceID: traceID}
	if input.PlaceID != "" {
		value.Context.PlaceID = input.PlaceID
	}
	if input.Locale != "" {
		value.Context.Locale = input.Locale
	}

	switch route.Intent {
	case "price_check", "price_explanation":
		result := o.price(ctx, value, input, base)
		traceFrom(&trace, &result)
		if err := o.sessions.AppendResponse(ctx, value, input.ID, result); err != nil {
			return model.Response{}, trace, err
		}
		return result, trace, nil
	case "emergency_help", "safety_assessment", "scam_pattern_assessment":
		result := o.safety(ctx, value, input, base)
		traceFrom(&trace, &result)
		if err := o.sessions.AppendResponse(ctx, value, input.ID, result); err != nil {
			return model.Response{}, trace, err
		}
		return result, trace, nil
	case "translation", "live_translation":
		result := o.translation(ctx, value, input, base)
		traceFrom(&trace, &result)
		if err := o.sessions.AppendResponse(ctx, value, input.ID, result); err != nil {
			return model.Response{}, trace, err
		}
		return result, trace, nil
	case "place_information", "place_discovery", "community_search", "general_travel_question":
		result := o.place(ctx, value, input, base)
		traceFrom(&trace, &result)
		if err := o.sessions.AppendResponse(ctx, value, input.ID, result); err != nil {
			return model.Response{}, trace, err
		}
		return result, trace, nil
	default:
		result := o.composer.Degraded(route.Intent, route.MissingFields, traceID)
		result.ID = base.ID
		traceFrom(&trace, &result)
		if err := o.sessions.AppendResponse(ctx, value, input.ID, result); err != nil {
			return model.Response{}, trace, err
		}
		return result, trace, nil
	}
}

func (o *Orchestrator) price(ctx context.Context, value *model.Session, input model.Message, result model.Response) model.Response {
	candidate, ok := intent.ExtractPriceCandidate(input.Text)
	if input.InputType == "structured_price_candidate" && len(input.Structured) > 0 {
		if json.Unmarshal(input.Structured, &candidate) == nil && candidate.AmountMinor != "" {
			ok = true
		}
	}
	if !ok {
		result.Message = "What amount and currency were you quoted?"
		result.MissingInformation = []string{"amount", "currency"}
		result.SuggestedActions = []model.SuggestedAction{{ID: "clarify-price", Label: "Add price details", Type: "clarify"}}
		return result
	}
	retrieved, retrievalResult := o.retrieve(ctx, input.Text, value.Context.Locale, first(input.PlaceID, value.Context.PlaceID), result.TraceID)
	result.ToolResults = append(result.ToolResults, retrievalResult)
	if retrievalResult.Status != "failed" {
		result.Evidence = append(result.Evidence, retrieved.Evidence...)
		if retrieved.PlaceID != "" {
			value.Context.PlaceID = retrieved.PlaceID
			value.Context.ApproximateRegion = retrieved.RegionID
		}
	}
	region := value.Context.ApproximateRegion
	if region == "" {
		result.Message = "Which place or region should Tourtect use for the price comparison?"
		result.MissingInformation = []string{"place_or_region"}
		result.SuggestedActions = []model.SuggestedAction{{ID: "clarify-place", Label: "Choose a place", Type: "clarify"}}
		return result
	}
	priceInput := pricing.PriceCheckInput{Vertical: candidate.Vertical, RawItem: candidate.RawItem, AmountMinor: candidate.AmountMinor, Currency: candidate.Currency, Exponent: candidate.Exponent, Unit: candidate.Unit, RegionID: region, ServiceSegment: "standard", VenueType: "transport_vendor", TransactionContext: "verbal_quote", ObservedAt: o.now().UTC(), ExtractionConfidence: .95, UserConfirmed: input.UserConfirmed}
	if candidate.Vertical != "taxi" {
		priceInput.VenueType = "fixed_shop"
	}
	raw, _ := json.Marshal(priceInput)
	toolResult := o.registry.Execute(ctx, "evaluate_price", raw, result.TraceID)
	result.ToolResults = append(result.ToolResults, toolResult)
	if toolResult.Status == "failed" {
		fallback := o.composer.Degraded(result.Intent, nil, result.TraceID)
		fallback.ID = result.ID
		fallback.ToolResults = result.ToolResults
		fallback.Evidence = result.Evidence
		return fallback
	}
	var priceOutput struct {
		EngineResult pricing.PriceCheckResult `json:"engine_result"`
	}
	if json.Unmarshal(toolResult.Output, &priceOutput) != nil || priceOutput.EngineResult.AlertLevel == "" {
		fallback := o.composer.Degraded(result.Intent, nil, result.TraceID)
		fallback.ID = result.ID
		fallback.ToolResults = result.ToolResults
		return fallback
	}
	engineResult := priceOutput.EngineResult
	result.Message, result.FactorsConsidered = o.composer.Price(&engineResult)
	result.Evidence = append(result.Evidence, o.evidence.Price(&engineResult)...)
	result.Confidence = engineResult.Confidence
	result.Freshness = engineResult.Freshness.Format(time.RFC3339)
	result.DatasetVersion = engineResult.DatasetVersion
	result.FallbackUsed = false
	result.SafetyState = "non_emergency"
	result.SuggestedActions = []model.SuggestedAction{{ID: "open-price-check", Label: "Review the full price check", Type: "deep_link", Target: "/price-check"}}
	return result
}

func (o *Orchestrator) safety(ctx context.Context, value *model.Session, input model.Message, result model.Response) model.Response {
	facts := intent.ExtractSafetyFacts(input.Text)
	region := directoryRegion(value.Context.ApproximateRegion)
	assessment := safety.AssessmentInput{ObservedFacts: facts.ObservedFacts, ThreatIndicators: facts.ThreatIndicators, ConfinementIndicators: facts.ConfinementIndicators, CoercionIndicators: facts.CoercionIndicators, AbilityToLeave: facts.AbilityToLeave, RegionID: region}
	raw, _ := json.Marshal(assessment)
	toolResult := o.registry.Execute(ctx, "evaluate_safety", raw, result.TraceID)
	result.ToolResults = append(result.ToolResults, toolResult)
	if toolResult.Status == "failed" {
		fallback := o.composer.Degraded(result.Intent, nil, result.TraceID)
		fallback.ID = result.ID
		fallback.ToolResults = result.ToolResults
		return fallback
	}
	var safetyOutput struct {
		Assessment safety.AssessmentResult `json:"assessment"`
	}
	if json.Unmarshal(toolResult.Output, &safetyOutput) != nil || safetyOutput.Assessment.Urgency == "" {
		fallback := o.composer.Degraded(result.Intent, nil, result.TraceID)
		fallback.ID = result.ID
		fallback.ToolResults = result.ToolResults
		return fallback
	}
	engineResult := safetyOutput.Assessment
	result.Message, result.FactorsConsidered = o.composer.Safety(&engineResult)
	result.Confidence = engineResult.Confidence
	result.Evidence = o.evidence.Safety(&engineResult)
	result.DatasetVersion = engineResult.SafetyDirectoryVersion
	result.FallbackUsed = false
	result.SafetyState = engineResult.Urgency
	result.SuggestedActions = []model.SuggestedAction{{ID: "manual-safety", Label: "Review structured safety facts", Type: "manual_safety_assessment", Target: "/safety"}, {ID: "offline-directory", Label: "Open emergency directory", Type: "offline_directory"}}
	if engineResult.SurfaceEmergencyOptions && len(engineResult.EmergencyContacts) > 0 {
		confirmation := &model.Confirmation{ID: uuid.NewString(), Action: "open_dialer", Title: "Open the phone dialer?", Description: "Tourtect will open the dialer with a verified directory number. It will not place the call.", ExpiresAt: o.now().UTC().Add(5 * time.Minute), Target: "tel:" + engineResult.EmergencyContacts[0].PhoneNumber}
		result.RequestedConfirmation = confirmation
		result.SuggestedActions = append(result.SuggestedActions, model.SuggestedAction{ID: "confirm-dialer", Label: "Open dialer", Type: "confirmation", RequiresConfirmation: true})
	}
	return result
}

func (o *Orchestrator) translation(ctx context.Context, value *model.Session, input model.Message, result model.Response) model.Response {
	if !value.Context.ConsentState.Processing {
		result.Message = "Allow processing for this session before sending text to the configured translation provider."
		result.MissingInformation = []string{"processing_consent"}
		result.SuggestedActions = []model.SuggestedAction{{ID: "phrasebook", Label: "Open offline phrasebook", Type: "deep_link", Target: "/saved?type=phrasebook"}}
		return result
	}
	target := value.Context.TargetLocale
	if target == "" {
		target = "en"
	}
	raw, _ := json.Marshal(fptai.TranslationInput{Text: input.Text, Target: target})
	toolResult := o.registry.Execute(ctx, "translate_text", raw, result.TraceID)
	result.ToolResults = append(result.ToolResults, toolResult)
	if toolResult.Status != "succeeded" {
		fallback := o.composer.Degraded(result.Intent, nil, result.TraceID)
		fallback.ID = result.ID
		fallback.ToolResults = result.ToolResults
		fallback.SuggestedActions = append(fallback.SuggestedActions, model.SuggestedAction{ID: "phrasebook", Label: "Open offline phrasebook", Type: "deep_link", Target: "/saved?type=phrasebook"})
		return fallback
	}
	var translated fptai.Translation
	if json.Unmarshal(toolResult.Output, &translated) != nil {
		return o.composer.Degraded(result.Intent, nil, result.TraceID)
	}
	result.Message = translated.Text
	result.FactorsConsidered = []string{"target_locale=" + target, "critical_tokens_preserved_by_provider_policy"}
	result.FallbackUsed = false
	result.SafetyState = "information"
	result.SuggestedActions = []model.SuggestedAction{}
	return result
}

func (o *Orchestrator) place(ctx context.Context, value *model.Session, input model.Message, result model.Response) model.Response {
	retrieved, retrievalResult := o.retrieve(ctx, input.Text, value.Context.Locale, first(input.PlaceID, value.Context.PlaceID), result.TraceID)
	result.ToolResults = append(result.ToolResults, retrievalResult)
	if retrievalResult.Status == "failed" || len(retrieved.Evidence) == 0 {
		fallback := o.composer.Degraded(result.Intent, []string{"resolvable_place_or_grounded_evidence"}, result.TraceID)
		fallback.ID = result.ID
		return fallback
	}
	if retrieved.PlaceID != "" {
		value.Context.PlaceID = retrieved.PlaceID
		value.Context.ApproximateRegion = retrieved.RegionID
	}
	result.Evidence = retrieved.Evidence
	result.Message = "Tourtect found verified place context and relevant public community knowledge. Review the evidence cards below; community reports are supporting context, not official determinations."
	result.FactorsConsidered = []string{"place_match", "freshness", "public_evidence_only", "source_diversity"}
	result.FallbackUsed = true
	result.SafetyState = "information"
	result.SuggestedActions = []model.SuggestedAction{{ID: "open-place", Label: "Open place details", Type: "deep_link", Target: "/places/" + value.Context.PlaceID}}
	return result
}

func traceFrom(trace *model.Trace, result *model.Response) {
	trace.FallbackUsed = result.FallbackUsed
	trace.RetrievalCount = len(result.Evidence)
	trace.Outcome = "succeeded"
	if result.FallbackUsed {
		trace.Outcome = "degraded"
	}
	for _, t := range result.ToolResults {
		trace.ToolNames = append(trace.ToolNames, t.ToolName)
		trace.ToolDurationsMS = append(trace.ToolDurationsMS, t.DurationMS)
		if t.Status == "failed" {
			trace.ErrorCategory = t.ErrorCategory
		}
	}
	for _, e := range result.Evidence {
		trace.EvidenceIDs = append(trace.EvidenceIDs, e.ID)
	}
}
func first(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
func directoryRegion(region string) string {
	lower := strings.ToLower(region)
	if strings.HasPrefix(lower, "hanoi") || strings.HasPrefix(lower, "ha-noi") {
		return "hanoi"
	}
	return region
}
func (o *Orchestrator) retrieve(ctx context.Context, text, locale, placeID, traceID string) (tools.RetrievalOutput, model.ToolResult) {
	raw, _ := json.Marshal(tools.RetrievalInput{Text: text, Locale: locale, PlaceID: placeID})
	toolResult := o.registry.Execute(ctx, "retrieve_place_context", raw, traceID)
	var output tools.RetrievalOutput
	if toolResult.Status != "failed" {
		_ = json.Unmarshal(toolResult.Output, &output)
	}
	return output, toolResult
}
func (o *Orchestrator) RegistrySpecs() []tools.Spec { return o.registry.Specs() }
func ValidateMessage(input model.Message) error {
	if input.ID == "" {
		return fmt.Errorf("message_id is required")
	}
	if input.InputType != "image_capture" && strings.TrimSpace(input.Text) == "" && len(input.Structured) == 0 {
		return fmt.Errorf("text or structured_data is required")
	}
	if input.InputType == "image_capture" && input.CaptureID == "" {
		return fmt.Errorf("capture_id is required")
	}
	return nil
}
