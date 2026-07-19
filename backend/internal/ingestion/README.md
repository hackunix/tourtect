# Tourtect Ingestion Pipeline Layer

This folder houses the provider-independent JSON contracts, validation engine, and provider abstractions for Tourtect's travel-critical content ingestion pipeline.

## Directory Structure

```
backend/internal/ingestion/
├── contracts/
│   ├── common.schema.json
│   ├── discovery-candidate.schema.json
│   ├── source-policy.schema.json
│   ├── fetch-request.schema.json
│   ├── fetch-result.schema.json
│   ├── raw-source-document.schema.json
│   ├── llm-extraction-request.schema.json
│   ├── llm-extraction-response.schema.json
│   ├── normalized-record.schema.json
│   ├── validation-result.schema.json
│   ├── deduplication-result.schema.json
│   ├── entity-link-result.schema.json
│   ├── review-decision.schema.json
│   ├── publishable-record.schema.json
│   ├── refresh-state.schema.json
│   ├── deletion-event.schema.json
│   ├── ingestion-job.schema.json
│   └── ingestion-trace.schema.json
│
├── prompts/
│   ├── system-prompt.md
│   ├── extraction-prompt-template.md
│   ├── normalization-prompt-template.md
│   └── provider-notes.md
│
├── examples/
│   ├── official-alert.json
│   ├── place.json
│   ├── price-observation.json
│   ├── scam-pattern.json
│   ├── news-article.json
│   ├── video.json
│   ├── invalid-hallucinated-hotline.json
│   ├── insufficient-evidence.json
│   ├── commerce-offer.json
│   └── conflicting-sources.json
│
├── validation/
│   ├── validator.go
│   └── validator_test.go
│
└── providers/
    ├── provider.go
    └── fake/
        ├── fake_provider.go
        └── fake_provider_test.go
```

## Running Pipeline Tests

To run the schema validation and fake provider pipeline tests locally, run:

```bash
cd backend
go test ./internal/ingestion/...
```

## Schema Validation in Go

The validation layer utilizes `github.com/xeipuuv/gojsonschema` for strict Draft 2020-12 checks. It resolves relative `$ref` tags locally without making external HTTP requests.

Example Go usage:

```go
package main

import (
	"fmt"
	"github.com/tourtect/backend/internal/ingestion/validation"
)

func main() {
	v := validation.NewPipelineValidator("./internal/ingestion/contracts")
	
	// Validate a record file
	res, err := v.ValidateRecord("llm-extraction-response.schema.json", rawJSONBytes)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Validation Status: %s\n", res.Status)
}
```

## Adding a New Provider Adapter

To integrate a new LLM provider (e.g., Gemini Vertex, OpenAI, FPT AI Factory):
1. Implement the `StructuredExtractionProvider` interface defined in `providers/provider.go`.
2. In your adapter's `Extract` method, format the system prompt and inputs according to the template.
3. Pass the target JSON Schema constraints parameter to the model's native structured outputs API.
4. Clean and parse the output JSON, mapping the model's usage metadata and response records back into Tourtect's provider-independent `ExtractionResponse` struct.
5. Do not leak provider-specific tokens or internal prompt headers into core domain models.
