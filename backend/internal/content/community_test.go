package content

import (
	"strings"
	"testing"
)

func TestFeedPoliciesAreExplainableAndOrganic(t *testing.T) {
	modes := []string{"following", "nearby", "latest", "trending", "safety"}
	for _, mode := range modes {
		filter, order, reason := feedPolicy(mode)
		if filter == "" || order == "" || reason == "" {
			t.Fatalf("mode %s has an incomplete policy", mode)
		}
		policy := strings.ToLower(filter + " " + order)
		for _, forbidden := range []string{"advertiser", "affiliate", "business_tier", "sponsored"} {
			if strings.Contains(policy, forbidden) {
				t.Fatalf("mode %s organic policy contains forbidden signal %s", mode, forbidden)
			}
		}
	}
}

func TestSafetyFeedUsesOnlySafetyContentAndEvidence(t *testing.T) {
	filter, order, reason := feedPolicy("safety")
	if !strings.Contains(filter, "official_alert") || !strings.Contains(filter, "scam_report") {
		t.Fatal("safety feed must be restricted to safety content")
	}
	if !strings.Contains(order, "evidence_level") || reason != "safety_priority" {
		t.Fatal("safety feed must be evidence-aware and explainable")
	}
}
