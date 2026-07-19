package tools

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tourtect/backend/adapters/fptai"
	"github.com/tourtect/backend/internal/pricing"
	"github.com/tourtect/backend/internal/safety"
)

type PriceEvaluator interface {
	Evaluate(context.Context, pricing.PriceCheckInput, string) (*pricing.PriceCheckResult, error)
}
type SafetyEvaluator interface {
	Assess(context.Context, safety.AssessmentInput, string) (*safety.AssessmentResult, error)
}

type PriceTool struct{ engine PriceEvaluator }

func NewPriceTool(engine PriceEvaluator) *PriceTool { return &PriceTool{engine: engine} }
func (t *PriceTool) Spec() Spec {
	return Spec{Name: "evaluate_price", Description: "Evaluate a confirmed or high-confidence price candidate with the deterministic Price Engine", InputSchema: "pricing.PriceCheckInput", OutputSchema: "pricing.PriceCheckResult", Timeout: 3 * time.Second, ErrorBehavior: "return insufficient data; never estimate", AuditBehavior: "trace metadata and snapshot provenance"}
}
func (t *PriceTool) Execute(ctx context.Context, raw json.RawMessage, traceID string) (json.RawMessage, string, error) {
	input, err := DecodeInput[pricing.PriceCheckInput](raw)
	if err != nil {
		return nil, "failed", err
	}
	result, err := t.engine.Evaluate(ctx, input, traceID)
	if err != nil {
		return nil, "failed", err
	}
	output := map[string]any{
		"engine_result": result,
		"insight": map[string]any{
			"alert_level":     result.AlertLevel,
			"observed":        map[string]any{"amount_minor": result.ObservedAmountMinor, "currency": result.ObservedCurrency, "exponent": result.ObservedExponent},
			"reference":       map[string]any{"p10_minor": result.ReferenceP10, "p50_minor": result.ReferenceP50, "p90_minor": result.ReferenceP90, "currency": result.ObservedCurrency, "exponent": result.ObservedExponent, "effective_sample_size": result.SampleSize},
			"deviation_ratio": result.DeviationRatio, "confidence": result.Confidence, "comparison_scope": result.ComparisonScope,
			"freshness": result.Freshness, "reasons": result.Reasons, "possible_benign_explanations": result.PossibleBenignExplanations,
			"dataset_version": result.DatasetVersion, "trace_id": result.TraceID,
		},
	}
	b, err := json.Marshal(output)
	if err != nil {
		return nil, "failed", err
	}
	status := "succeeded"
	if result.AlertLevel == "insufficient_data" {
		status = "insufficient_data"
	}
	return b, status, nil
}

type TranslationTool struct{ provider fptai.TranslationProvider }

func NewTranslationTool(provider fptai.TranslationProvider) *TranslationTool {
	return &TranslationTool{provider: provider}
}
func (t *TranslationTool) Spec() Spec {
	return Spec{Name: "translate_text", Description: "Translate bounded text without changing verified price or safety facts", InputSchema: "fptai.TranslationInput", OutputSchema: "fptai.Translation", RequiredConsent: "processing", Timeout: 12 * time.Second, ErrorBehavior: "provider degraded; show phrasebook", AuditBehavior: "provider/model metadata without raw transcript"}
}
func (t *TranslationTool) Execute(ctx context.Context, raw json.RawMessage, _ string) (json.RawMessage, string, error) {
	input, err := DecodeInput[fptai.TranslationInput](raw)
	if err != nil {
		return nil, "failed", err
	}
	result, err := t.provider.Translate(ctx, input)
	if err != nil {
		return nil, "degraded", err
	}
	b, err := json.Marshal(result)
	if err != nil {
		return nil, "failed", err
	}
	return b, "succeeded", nil
}

type SafetyTool struct{ engine SafetyEvaluator }

func NewSafetyTool(engine SafetyEvaluator) *SafetyTool { return &SafetyTool{engine: engine} }
func (t *SafetyTool) Spec() Spec {
	return Spec{Name: "evaluate_safety", Description: "Evaluate extracted facts with the rule-first Safety Engine", InputSchema: "safety.AssessmentInput", OutputSchema: "safety.AssessmentResult", Timeout: 3 * time.Second, ErrorBehavior: "surface offline safety options; never invent a hotline", AuditBehavior: "trace explanation codes and directory version"}
}
func (t *SafetyTool) Execute(ctx context.Context, raw json.RawMessage, traceID string) (json.RawMessage, string, error) {
	input, err := DecodeInput[safety.AssessmentInput](raw)
	if err != nil {
		return nil, "failed", err
	}
	result, err := t.engine.Assess(ctx, input, traceID)
	if err != nil {
		return nil, "failed", err
	}
	b, err := json.Marshal(result)
	if err != nil {
		return nil, "failed", err
	}
	return b, "succeeded", nil
}
