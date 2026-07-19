package fake

import (
	"context"
	"testing"

	"github.com/tourtect/backend/internal/ingestion/providers"
)

func TestFakeProvider_Extract(t *testing.T) {
	ctx := context.Background()

	mockRecord := map[string]any{
		"record_id":   "rec-123",
		"record_type": "place_status",
		"content": map[string]any{
			"status":   "temporarily_closed",
			"evidence": "Closed for maintenance",
		},
		"provenance": map[string]any{
			"claims": []any{},
		},
	}

	p := NewFakeProvider("completed", []map[string]any{mockRecord})

	req := providers.ExtractionRequest{
		SchemaVersion:      "1.0.0",
		RequestID:          "req-abc",
		Task:               "extract_travel_critical_data",
		AllowedRecordTypes: []string{"place_status"},
	}

	res, err := p.Extract(ctx, req, []byte("{}"))
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if res.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", res.Status)
	}

	if len(res.Records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(res.Records))
	}

	if res.Records[0]["record_id"] != "rec-123" {
		t.Errorf("Expected record_id 'rec-123', got '%v'", res.Records[0]["record_id"])
	}
}

func TestFakeProvider_Error(t *testing.T) {
	ctx := context.Background()
	p := NewFakeProvider("error", nil)

	req := providers.ExtractionRequest{
		SchemaVersion: "1.0.0",
		RequestID:     "req-abc",
	}

	_, err := p.Extract(ctx, req, []byte("{}"))
	if err == nil {
		t.Error("Expected error from FakeProvider, got nil")
	}
}
