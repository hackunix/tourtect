package response

import (
	"fmt"

	"github.com/tourtect/backend/internal/intelligence/model"
	"github.com/tourtect/backend/internal/pricing"
	"github.com/tourtect/backend/internal/safety"
)

type Composer struct{}

func NewComposer() *Composer { return &Composer{} }

func (c *Composer) Price(result *pricing.PriceCheckResult) (string, []string) {
	factors := []string{fmt.Sprintf("entered_amount=%s %s", result.ObservedAmountMinor, result.ObservedCurrency)}
	if result.AlertLevel == "insufficient_data" {
		return "Tourtect does not have enough matching price evidence to make a comparison. You can correct the item, unit, place, or currency and try again.", append(factors, result.Reasons...)
	}
	factors = append(factors,
		fmt.Sprintf("reference_range=%s-%s %s", result.ReferenceP10, result.ReferenceP90, result.ObservedCurrency),
		fmt.Sprintf("sample_size=%d", result.SampleSize), fmt.Sprintf("comparison_scope=%s", result.ComparisonScope))
	switch result.AlertLevel {
	case "high_risk":
		return "This price is significantly above the available reference range. That is a price signal, not an accusation of fraud. Check the unit and route before deciding what to do next.", factors
	case "elevated":
		return "This price is above the available reference range. There may be benign reasons, so verify the unit, route, and any stated fees.", factors
	default:
		return "This price is within the available reference range for the matched cohort.", factors
	}
}

func (c *Composer) Safety(result *safety.AssessmentResult) (string, []string) {
	factors := append([]string{}, result.ExplanationCodes...)
	switch result.Urgency {
	case "critical":
		return "Tourtect detected facts that the rule-first Safety Engine treats as critical. Move toward a safer public place if you can do so without escalating the situation. Use only the verified emergency options shown below.", factors
	case "urgent":
		return "Tourtect detected an urgent safety concern. Follow the approved actions shown below and confirm any consequential action before it occurs.", factors
	case "non_emergency":
		return "No immediate threat was detected from the facts provided. Keep distance, avoid escalation, and update Tourtect if the situation changes.", factors
	default:
		return "Tourtect found no immediate safety signal in the facts provided.", factors
	}
}

func (c *Composer) Degraded(intent string, missing []string, traceID string) model.Response {
	return model.Response{Intent: intent, Message: "Tourtect cannot interpret this automatically right now. You can still use Manual Price Check, the rule-first Safety Assessment, the offline emergency directory, or save a private draft.", Confidence: 0, Evidence: []model.Evidence{}, ToolResults: []model.ToolResult{}, SuggestedActions: []model.SuggestedAction{{ID: "manual-price", Label: "Enter the price manually", Type: "manual_price_check", Target: "/price-check"}, {ID: "manual-safety", Label: "Select safety facts", Type: "manual_safety_assessment", Target: "/safety"}, {ID: "offline-directory", Label: "Open the offline emergency directory", Type: "offline_directory"}, {ID: "private-draft", Label: "Save a private draft", Type: "save_private_draft"}}, SafetyState: "unknown", FactorsConsidered: []string{}, MissingInformation: missing, FallbackUsed: true, TraceID: traceID}
}
