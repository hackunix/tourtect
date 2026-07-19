package fake

import (
	"context"
	"fmt"

	"github.com/tourtect/backend/internal/ingestion/providers"
)

type FakeProvider struct {
	MockStatus  string
	MockRecords []map[string]any
	MockRefusal *string
}

func NewFakeProvider(status string, records []map[string]any) *FakeProvider {
	return &FakeProvider{
		MockStatus:  status,
		MockRecords: records,
	}
}

func (f *FakeProvider) Extract(
	ctx context.Context,
	request providers.ExtractionRequest,
	outputSchema []byte,
) (*providers.ExtractionResponse, error) {
	if f.MockStatus == "error" {
		return nil, fmt.Errorf("simulated provider error")
	}

	return &providers.ExtractionResponse{
		SchemaVersion:      "1.0.0",
		RequestID:          request.RequestID,
		ProviderResponseID: "fake-response-999",
		Status:             f.MockStatus,
		Records:            f.MockRecords,
		ModelUsage: map[string]any{
			"prompt_tokens":     100,
			"completion_tokens": 50,
		},
		Refusal: f.MockRefusal,
	}, nil
}
