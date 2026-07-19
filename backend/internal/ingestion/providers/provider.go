package providers

import (
	"context"
)

type ExtractionRequest struct {
	SchemaVersion      string         `json:"schema_version"`
	RequestID          string         `json:"request_id"`
	Task               string         `json:"task"`
	AllowedRecordTypes []string       `json:"allowed_record_types"`
	SourceDocument     any            `json:"source_document"`
	Context            map[string]any `json:"context"`
	Constraints        map[string]any `json:"constraints"`
}

type ExtractionResponse struct {
	SchemaVersion        string         `json:"schema_version"`
	RequestID            string         `json:"request_id"`
	ProviderResponseID   string         `json:"provider_response_id"`
	Status               string         `json:"status"` // completed, partial, insufficient_evidence, provider_refusal, failed
	Records              []map[string]any `json:"records"`
	DocumentClassification map[string]any `json:"document_classification,omitempty"`
	Warnings             []string       `json:"warnings,omitempty"`
	ModelUsage           map[string]any `json:"model_usage,omitempty"`
	Refusal              *string        `json:"refusal,omitempty"`
}

type StructuredExtractionProvider interface {
	Extract(
		ctx context.Context,
		request ExtractionRequest,
		outputSchema []byte,
	) (*ExtractionResponse, error)
}
