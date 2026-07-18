package safety

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tourtect/backend/generated/database"
)

type Engine struct {
	pool    *pgxpool.Pool
	queries *database.Queries
}

func NewEngine(pool *pgxpool.Pool) *Engine {
	return &Engine{
		pool:    pool,
		queries: database.New(pool),
	}
}

type AssessmentInput struct {
	ObservedFacts        []string
	UserReportedState    string
	ThreatIndicators     []string
	InjuryIndicators     []string
	ConfinementIndicators []string
	CoercionIndicators   []string
	AbilityToLeave       *bool
	UserConfirmedFacts   []string
	RegionID             string
}

type AssessmentResult struct {
	Urgency                 string
	SafeActions             []string
	ApprovedActionCodes     []string
	ExplanationCodes        []string
	SilentModeRecommended   bool
	SurfaceEmergencyOptions bool
	EmergencyServiceIDs     []string
	SafetyDirectoryVersion  string
	Confidence              float64
	TraceID                 string
	EmergencyContacts       []Contact
}

type Contact struct {
	Name        string
	PhoneNumber string
	Type        string
}

func (e *Engine) Assess(ctx context.Context, input AssessmentInput, traceID string) (*AssessmentResult, error) {
	if traceID == "" {
		traceID = uuid.New().String()
	}

	result := &AssessmentResult{
		Confidence:             1.0,
		TraceID:                traceID,
		SafeActions:            []string{},
		ApprovedActionCodes:    []string{},
		ExplanationCodes:       []string{},
		EmergencyServiceIDs:    []string{},
		SafetyDirectoryVersion: "unknown",
	}

	// Step 1: Query the active safety directory version from Postgres
	activeVer, err := e.queries.GetLatestSafetyDirectoryVersion(ctx)
	if err != nil {
		slog.Warn("Safety directory unavailable in database, continuing with safety degradation policy")
		result.ExplanationCodes = append(result.ExplanationCodes, "safety_directory_unavailable")
	} else {
		result.SafetyDirectoryVersion = activeVer.Version

		// Step 2: Query contact list for the region
		region := input.RegionID
		if region == "" {
			region = "hanoi" // Default region matching seed data
		}
		entries, err := e.queries.GetSafetyEntriesByRegion(ctx, database.GetSafetyEntriesByRegionParams{
			VersionID: activeVer.VersionID,
			RegionID:  region,
		})
		if err == nil {
			contacts := make([]Contact, 0, len(entries))
			for _, entry := range entries {
				contacts = append(contacts, Contact{
					Name:        entry.ServiceName,
					PhoneNumber: entry.PhoneNumber,
					Type:        entry.ServiceType,
				})
				result.EmergencyServiceIDs = append(result.EmergencyServiceIDs, entry.EntryID.String())
			}
			result.EmergencyContacts = contacts
		}
	}

	// Step 3: Rule-first evaluation
	var isConfinement bool
	var isThreat bool
	var isInjury bool
	var isCoercion bool

	if len(input.ConfinementIndicators) > 0 || (input.AbilityToLeave != nil && !*input.AbilityToLeave) {
		isConfinement = true
	}
	for _, t := range input.ThreatIndicators {
		if t == "weapon" || t == "physical_violence" || t == "confinement" {
			isThreat = true
		}
	}
	if len(input.InjuryIndicators) > 0 {
		isInjury = true
	}
	if len(input.CoercionIndicators) > 0 {
		isCoercion = true
	}

	// Rule 4: Urgency classification mapping
	if isConfinement || isThreat {
		result.Urgency = "critical"
		result.SilentModeRecommended = true
		result.SurfaceEmergencyOptions = true
		result.ExplanationCodes = append(result.ExplanationCodes, "confinement_or_physical_threat_detected")
		result.SafeActions = append(result.SafeActions, "Seek a safe space immediately. Do not make sudden moves or escalate verbal conflict.")
		result.ApprovedActionCodes = append(result.ApprovedActionCodes, "SEEK_REFUGE", "SILENT_ALERT")
	} else if isInjury {
		result.Urgency = "urgent"
		result.SilentModeRecommended = false
		result.SurfaceEmergencyOptions = true
		result.ExplanationCodes = append(result.ExplanationCodes, "injury_detected")
		result.SafeActions = append(result.SafeActions, "Contact local medical emergency services immediately.")
		result.ApprovedActionCodes = append(result.ApprovedActionCodes, "CALL_AMBULANCE")
	} else if isCoercion {
		result.Urgency = "urgent"
		result.SilentModeRecommended = false
		result.SurfaceEmergencyOptions = true
		result.ExplanationCodes = append(result.ExplanationCodes, "coerced_transaction_detected")
		result.SafeActions = append(result.SafeActions, "Refuse transaction safely in public place, or pay minimum amount to leave if unsafe, then file police report.")
		result.ApprovedActionCodes = append(result.ApprovedActionCodes, "REFUSE_PAYMENT", "SEEK_POLICE_SUPPORT")
	} else {
		// Non-emergency or informational
		var hasDispute bool
		for _, f := range input.ObservedFacts {
			if f == "price_dispute" || f == "verbal_disagreement" {
				hasDispute = true
			}
		}

		if hasDispute {
			result.Urgency = "non_emergency"
			result.ExplanationCodes = append(result.ExplanationCodes, "verbal_dispute_no_immediate_threat")
			result.SafeActions = append(result.SafeActions, "Politely decline further service. Walk away to a public area.")
			result.ApprovedActionCodes = append(result.ApprovedActionCodes, "DEPART_VENUE")
		} else {
			result.Urgency = "information"
			result.ExplanationCodes = append(result.ExplanationCodes, "informational_safety_query")
			result.SafeActions = append(result.SafeActions, "Review local rules, standard rates, and approved service providers.")
			result.ApprovedActionCodes = append(result.ApprovedActionCodes, "REVIEW_GUIDELINES")
		}
	}

	// Append emergency phone numbers explicitly from database (never hardcoded/hallucinated)
	for _, c := range result.EmergencyContacts {
		result.SafeActions = append(result.SafeActions, fmt.Sprintf("Call %s at %s", c.Name, c.PhoneNumber))
	}

	return result, nil
}
