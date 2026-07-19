package tools

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tourtect/backend/internal/intelligence/model"
	"github.com/tourtect/backend/internal/intelligence/retrieval"
)

type RetrievalInput struct {
	Text    string `json:"text"`
	Locale  string `json:"locale"`
	PlaceID string `json:"place_id,omitempty"`
}
type RetrievalOutput struct {
	PlaceID   string           `json:"place_id,omitempty"`
	RegionID  string           `json:"region_id,omitempty"`
	PlaceName string           `json:"place_name,omitempty"`
	Evidence  []model.Evidence `json:"evidence"`
}

type RetrievalTool struct{ service *retrieval.Service }

func NewRetrievalTool(service *retrieval.Service) *RetrievalTool {
	return &RetrievalTool{service: service}
}
func (t *RetrievalTool) Spec() Spec {
	return Spec{Name: "retrieve_place_context", Description: "Resolve a place and retrieve bounded public or verified Tourtect evidence", InputSchema: "tools.RetrievalInput", OutputSchema: "tools.RetrievalOutput", Timeout: 3 * time.Second, ErrorBehavior: "return no evidence and ask for place clarification", AuditBehavior: "evidence identifiers, source count, and duration only"}
}
func (t *RetrievalTool) Execute(ctx context.Context, raw json.RawMessage, _ string) (json.RawMessage, string, error) {
	input, err := DecodeInput[RetrievalInput](raw)
	if err != nil {
		return nil, "failed", err
	}
	result, err := t.service.Retrieve(ctx, input.Text, input.Locale, input.PlaceID)
	if err != nil {
		return nil, "failed", err
	}
	output := RetrievalOutput{Evidence: result.Evidence}
	if result.Place != nil {
		output.PlaceID = result.Place.PlaceID.String()
		output.RegionID = result.Place.RegionID
		output.PlaceName = result.Place.Name
	}
	b, err := json.Marshal(output)
	if err != nil {
		return nil, "failed", err
	}
	status := "succeeded"
	if len(output.Evidence) == 0 {
		status = "insufficient_data"
	}
	return b, status, nil
}
