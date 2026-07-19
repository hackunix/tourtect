package evidence

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tourtect/backend/internal/intelligence/model"
	"github.com/tourtect/backend/internal/pricing"
	"github.com/tourtect/backend/internal/safety"
)

type Assembler struct{ now func() time.Time }

func NewAssembler() *Assembler { return &Assembler{now: time.Now} }

func (a *Assembler) Price(result *pricing.PriceCheckResult) []model.Evidence {
	if result.SnapshotID == "" {
		return []model.Evidence{}
	}
	observed := result.SnapshotEffectiveFrom
	return []model.Evidence{{ID: uuid.NewString(), SourceType: "price_snapshot", SourceID: result.SnapshotID, Title: "Matched price reference", Summary: fmt.Sprintf("Reference range %s–%s %s based on %d observations from %d independent sources.", result.ReferenceP10, result.ReferenceP90, result.ObservedCurrency, result.SampleSize, result.IndependentSourceCount), ObservedAt: &observed, Freshness: age(observed, a.now()), EvidenceLevel: "verified"}}
}

func (a *Assembler) Safety(result *safety.AssessmentResult) []model.Evidence {
	if result.SafetyDirectoryVersion == "" || result.SafetyDirectoryVersion == "unknown" {
		return []model.Evidence{}
	}
	return []model.Evidence{{ID: uuid.NewString(), SourceType: "safety_directory", SourceID: result.SafetyDirectoryVersion, Title: "Verified safety directory", Summary: fmt.Sprintf("Approved emergency directory version %s; %d matching regional services.", result.SafetyDirectoryVersion, len(result.EmergencyContacts)), Freshness: "unknown", EvidenceLevel: "official"}}
}

func age(t, now time.Time) string {
	if t.IsZero() {
		return "unknown"
	}
	d := now.Sub(t)
	if d <= 30*24*time.Hour {
		return "fresh"
	}
	if d <= 180*24*time.Hour {
		return "aging"
	}
	return "stale"
}
