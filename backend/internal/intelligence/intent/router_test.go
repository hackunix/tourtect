package intent

import "testing"

func TestCriticalSafetyOverridesPrice(t *testing.T) {
	route := NewRouter().Route("text", "The driver wants 900,000 VND and will not let me leave.")
	if route.Intent != "emergency_help" || !route.SafetyOverride {
		t.Fatalf("expected critical safety override, got %+v", route)
	}
}

func TestExtractTaxiPrice(t *testing.T) {
	candidate, ok := ExtractPriceCandidate("The driver wants 900,000 VND from Noi Bai to Hoan Kiem.")
	if !ok || candidate.AmountMinor != "900000" || candidate.Currency != "VND" || candidate.Unit != "trip" || candidate.Vertical != "taxi" {
		t.Fatalf("unexpected candidate: %+v", candidate)
	}
}

func TestLowRiskHighPriceDoesNotBecomeEmergency(t *testing.T) {
	route := NewRouter().Route("text", "The taxi price is 900,000 VND, but I can leave.")
	if route.Intent != "price_check" || route.SafetyOverride {
		t.Fatalf("unexpected route: %+v", route)
	}
}
