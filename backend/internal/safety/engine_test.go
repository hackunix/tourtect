package safety

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tourtect/backend/generated/database"
)

func TestSafetyEngine(t *testing.T) {
	ctx := context.Background()
	dsn := "postgres://tourtect:change_me_postgres@localhost:5432/tourtect?sslmode=disable"
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Skip("Database not available, skipping safety engine tests")
		return
	}
	defer pool.Close()

	engine := NewEngine(pool)

	// Wrap in a transaction to insert custom directory version + entries and rollback
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	engine.queries = database.New(tx)

	// Clean up database state in this transaction context
	_, _ = tx.Exec(ctx, "DELETE FROM safety_directory_entries")
	_, _ = tx.Exec(ctx, "DELETE FROM safety_directory_versions")

	// Insert test directory version
	versionID := uuid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO safety_directory_versions (version_id, version, description)
		VALUES ($1, 'test-safety-v1', 'Test safety directory')
	`, versionID)
	if err != nil {
		t.Fatalf("failed to insert safety version: %v", err)
	}

	// Insert approved hotline numbers
	_, err = tx.Exec(ctx, `
		INSERT INTO safety_directory_entries (version_id, region_id, service_name, service_type, phone_number, is_approved)
		VALUES
			($1, 'hanoi', 'Police Dept', 'police', '113', true),
			($1, 'hanoi', 'Ambulance Dept', 'ambulance', '115', true),
			($1, 'hanoi', 'Tourist Support', 'tourist_police', '069-942-0626', true)
	`, versionID)
	if err != nil {
		t.Fatalf("failed to insert safety entries: %v", err)
	}

	// Helper for pointer booleans
	boolPtr := func(b bool) *bool { return &b }

	// Test 1: High price without coercion (non_emergency or information)
	t.Run("High price without coercion", func(t *testing.T) {
		res, err := engine.Assess(ctx, AssessmentInput{
			ObservedFacts: []string{"high_price"},
			RegionID:      "hanoi",
		}, "trace-1")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if res.Urgency == "critical" || res.Urgency == "urgent" {
			t.Errorf("Expected low urgency for high price without coercion, got %s", res.Urgency)
		}
	})

	// Test 2: Forced payment (coercion triggers urgent)
	t.Run("Forced payment coercion", func(t *testing.T) {
		res, err := engine.Assess(ctx, AssessmentInput{
			ObservedFacts:      []string{"overcharged"},
			CoercionIndicators: []string{"forced_payment"},
			RegionID:           "hanoi",
		}, "trace-2")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if res.Urgency != "urgent" {
			t.Errorf("Expected urgency 'urgent' for forced payment, got %s", res.Urgency)
		}
	})

	// Test 3: Driver refuses to let the user leave (confinement triggers critical, silent mode recommended)
	t.Run("Refuse to let user leave (confinement)", func(t *testing.T) {
		res, err := engine.Assess(ctx, AssessmentInput{
			ObservedFacts:  []string{"driver_hostile"},
			AbilityToLeave: boolPtr(false),
			RegionID:       "hanoi",
		}, "trace-3")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if res.Urgency != "critical" {
			t.Errorf("Expected urgency 'critical' for inability to leave, got %s", res.Urgency)
		}
		if !res.SilentModeRecommended {
			t.Errorf("Expected silent mode to be recommended for confinement")
		}
	})

	// Test 4: Injury (injury triggers urgent, CALL_AMBULANCE action code)
	t.Run("Injury detected", func(t *testing.T) {
		res, err := engine.Assess(ctx, AssessmentInput{
			ObservedFacts:    []string{"physical_altercation"},
			InjuryIndicators: []string{"bleeding"},
			RegionID:         "hanoi",
		}, "trace-4")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if res.Urgency != "urgent" && res.Urgency != "critical" {
			t.Errorf("Expected urgent or critical urgency for injury, got %s", res.Urgency)
		}
		hasAmbulanceCode := false
		for _, code := range res.ApprovedActionCodes {
			if code == "CALL_AMBULANCE" {
				hasAmbulanceCode = true
			}
		}
		if !hasAmbulanceCode {
			t.Errorf("Expected approved action codes to contain CALL_AMBULANCE")
		}
	})

	// Test 5: Weapon mention (threat weapon triggers critical, silent mode recommended)
	t.Run("Weapon threat", func(t *testing.T) {
		res, err := engine.Assess(ctx, AssessmentInput{
			ObservedFacts:    []string{"confrontation"},
			ThreatIndicators: []string{"weapon"},
			RegionID:         "hanoi",
		}, "trace-5")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if res.Urgency != "critical" {
			t.Errorf("Expected urgency 'critical' for weapon threat, got %s", res.Urgency)
		}
		if !res.SilentModeRecommended {
			t.Errorf("Expected silent mode to be recommended for weapon threat")
		}
	})

	// Test 6: Informational question (triggers informational urgency)
	t.Run("Informational safety question", func(t *testing.T) {
		res, err := engine.Assess(ctx, AssessmentInput{
			ObservedFacts: []string{"asking_guidelines"},
			RegionID:      "hanoi",
		}, "trace-6")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if res.Urgency != "information" {
			t.Errorf("Expected urgency 'information' for basic question, got %s", res.Urgency)
		}
	})

	// Test 7: Missing facts (graceful degradation)
	t.Run("Degrades with missing facts", func(t *testing.T) {
		res, err := engine.Assess(ctx, AssessmentInput{
			ObservedFacts: []string{}, // completely empty
			RegionID:      "hanoi",
		}, "trace-7")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if res.Urgency != "information" {
			t.Errorf("Expected default informational urgency for empty inputs, got %s", res.Urgency)
		}
	})

	// Test 8: Conflicting facts (prioritizes safety first)
	t.Run("Conflicting safety facts", func(t *testing.T) {
		res, err := engine.Assess(ctx, AssessmentInput{
			ObservedFacts:         []string{"everything_ok"},
			ConfinementIndicators: []string{"door_locked"}, // contradicts observed_facts
			AbilityToLeave:        boolPtr(true),          // contradicts confinement indicator
			RegionID:              "hanoi",
		}, "trace-8")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Prioritizes confinement indicator over contradictory fields
		if res.Urgency != "critical" {
			t.Errorf("Expected safety-first resolution (critical) for conflicting inputs, got %s", res.Urgency)
		}
	})

	// Test 9: LLM provider unavailable (proves rule-first evaluation works without provider)
	t.Run("Evaluates without LLM provider", func(t *testing.T) {
		res, err := engine.Assess(ctx, AssessmentInput{
			ObservedFacts:    []string{"incident"},
			ThreatIndicators: []string{"physical_violence"},
			RegionID:         "hanoi",
		}, "trace-9")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if res.Urgency != "critical" {
			t.Errorf("Expected rule evaluation to return critical, got %s", res.Urgency)
		}
	})

	// Test 10: Safety directory unavailable (graceful degradation)
	t.Run("Degrades gracefully when database safety directory is empty", func(t *testing.T) {
		// Temporarily wipe directories to simulate failure
		_, _ = tx.Exec(ctx, "DELETE FROM safety_directory_entries")
		_, _ = tx.Exec(ctx, "DELETE FROM safety_directory_versions")

		res, err := engine.Assess(ctx, AssessmentInput{
			ObservedFacts: []string{"high_price"},
			RegionID:      "hanoi",
		}, "trace-10")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if res.SafetyDirectoryVersion != "unknown" {
			t.Errorf("Expected safety directory version to be 'unknown', got %s", res.SafetyDirectoryVersion)
		}
		foundUnavailableCode := false
		for _, code := range res.ExplanationCodes {
			if code == "safety_directory_unavailable" {
				foundUnavailableCode = true
			}
		}
		if !foundUnavailableCode {
			t.Errorf("Expected explanation codes to reflect safety_directory_unavailable")
		}
	})

	// Test 11: No hallucinated hotline (all numbers come from db)
	t.Run("No hotline hallucination", func(t *testing.T) {
		// Re-insert approved numbers
		_, _ = tx.Exec(ctx, `
			INSERT INTO safety_directory_versions (version_id, version, description)
			VALUES ($1, 'test-safety-v2', 'Clean safety directory')
		`, versionID)
		_, _ = tx.Exec(ctx, `
			INSERT INTO safety_directory_entries (version_id, region_id, service_name, service_type, phone_number, is_approved)
			VALUES ($1, 'hanoi', 'Police Dept', 'police', '113', true)
		`, versionID)

		res, err := engine.Assess(ctx, AssessmentInput{
			ObservedFacts: []string{"dispute"},
			RegionID:      "hanoi",
		}, "trace-11")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Ensure any phone number mentioned in SafeActions matches the database
		for _, action := range res.SafeActions {
			if strings.Contains(action, "Call ") {
				if !strings.Contains(action, "113") {
					t.Errorf("Hallucinated hotline or incorrect phone number printed in action instructions: %s", action)
				}
			}
		}
	})
}
