package aievals

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/tourtect/backend/internal/intelligence/intent"
)

type scenario struct {
	ID                     string `json:"id"`
	Input                  string `json:"input"`
	InputType              string `json:"input_type"`
	ExpectedIntent         string `json:"expected_intent"`
	ExpectedAmountMinor    string `json:"expected_amount_minor"`
	ExpectedSafetyOverride bool   `json:"expected_safety_override"`
}

func TestDeterministicIntentAndExtractionFixtures(t *testing.T) {
	b, err := os.ReadFile("scenarios.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures []scenario
	if err := json.Unmarshal(b, &fixtures); err != nil {
		t.Fatal(err)
	}
	router := intent.NewRouter()
	for _, fixture := range fixtures {
		t.Run(fixture.ID, func(t *testing.T) {
			inputType := fixture.InputType
			if inputType == "" {
				inputType = "text"
			}
			route := router.Route(inputType, fixture.Input)
			if route.Intent != fixture.ExpectedIntent {
				t.Fatalf("intent=%s want %s", route.Intent, fixture.ExpectedIntent)
			}
			if route.SafetyOverride != fixture.ExpectedSafetyOverride {
				t.Fatalf("safety_override=%v want %v", route.SafetyOverride, fixture.ExpectedSafetyOverride)
			}
			if fixture.ExpectedAmountMinor != "" {
				candidate, ok := intent.ExtractPriceCandidate(fixture.Input)
				if !ok || candidate.AmountMinor != fixture.ExpectedAmountMinor {
					t.Fatalf("candidate=%+v ok=%v", candidate, ok)
				}
			}
		})
	}
}
