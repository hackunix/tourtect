package validation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestValidator_ValidExamples(t *testing.T) {
	contractsDir := "../contracts"
	examplesDir := "../examples"
	v := NewPipelineValidator(contractsDir)

	files := []string{
		"official-alert.json",
		"place.json",
		"price-observation.json",
		"scam-pattern.json",
		"news-article.json",
		"video.json",
		"commerce-offer.json",
		"conflicting-sources.json",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			path := filepath.Join(examplesDir, file)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read example file %s: %v", file, err)
			}

			// All our valid/active examples are records conforming to llm-extraction-response schema,
			// or individual records.
			// Let's determine which schema to validate against:
			// If it contains "records", it's a full extraction response; otherwise, it's a common envelope record.
			// Let's validate the individual records or responses against common envelope structures or the generic schemas.
			// In our examples, `official-alert.json`, `place.json` etc are individual envelopes of ExtractedRecord
			// which matches our envelope schema constraints.
			// Wait, let's look at `llm-extraction-response.schema.json` which expects a wrapper response.
			// Since we want to validate these records individually against the record contract structure,
			// let's wrap them or check them using llm-extraction-response schema or validate directly.
			// In `validator.go`, we parse and validate against a specific schema file. Let's validate against `llm-extraction-response.schema.json` by wrapping it, or validate the envelope.
			// Wait! Let's check how the examples are structured: they are individual records with common envelopes.
			// Let's validate them against `llm-extraction-response.schema.json` by placing them inside a mock response!
			// That is very elegant.
			var record map[string]any
			if err := json.Unmarshal(data, &record); err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			var responseData []byte
			if _, isResponse := record["records"]; isResponse {
				responseData = data
			} else {
				// Wrap individual record in LLMExtractionResponse
				response := map[string]any{
					"schema_version": "1.0.0",
					"request_id":     "req-test-123",
					"status":         "completed",
					"records":        []any{record},
				}
				wrapped, err := json.Marshal(response)
				if err != nil {
					t.Fatalf("Failed to wrap record: %v", err)
				}
				responseData = wrapped
			}

			res, err := v.ValidateRecord("llm-extraction-response.schema.json", responseData)
			if err != nil {
				t.Fatalf("Validator execution failed: %v", err)
			}

			// We expect no JSON Schema failures
			for _, failure := range res.Failures {
				if failure.Category == "json_schema" {
					t.Errorf("Unexpected JSON Schema failure: %s - %s", failure.FieldPath, failure.Message)
				}
			}
		})
	}
}

func TestValidator_SafetyValidation(t *testing.T) {
	contractsDir := "../contracts"
	examplesDir := "../examples"
	v := NewPipelineValidator(contractsDir)

	t.Run("hallucinated hotline", func(t *testing.T) {
		path := filepath.Join(examplesDir, "invalid-hallucinated-hotline.json")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		res, err := v.ValidateRecord("llm-extraction-response.schema.json", wrapInResponse(data))
		if err != nil {
			t.Fatalf("Validator failed: %v", err)
		}

		if res.Status != "quarantined" {
			t.Errorf("Expected status to be quarantined, got %s", res.Status)
		}

		foundHotlineFailure := false
		for _, f := range res.Failures {
			if f.Category == "hallucinated_hotline_risk" {
				foundHotlineFailure = true
			}
		}
		if !foundHotlineFailure {
			t.Error("Expected to find hallucinated hotline risk failure")
		}
	})

	t.Run("price outlier", func(t *testing.T) {
		// Create a taxi record with price > 10,000,000 VND
		priceRecord := map[string]any{
			"schema_version": "1.0.0",
			"record_id":      "rec-outlier-price",
			"record_type":    "price_observation",
			"pipeline_stage": "ai_extracted",
			"source": map[string]any{
				"source_id":     "src-1",
				"source_type":   "rss",
				"canonical_url": "https://example.com",
			},
			"content": map[string]any{
				"item_name": "taxi ride",
				"amount": map[string]any{
					"minor_units": "15000000", // 15 million VND
					"currency":    "VND",
					"scale":       0,
				},
				"unit": "trip",
			},
			"provenance": map[string]any{
				"claims": []any{},
			},
		}

		data, _ := json.Marshal(priceRecord)
		res, err := v.ValidateRecord("llm-extraction-response.schema.json", wrapInResponse(data))
		if err != nil {
			t.Fatalf("Validator failed: %v", err)
		}

		foundOutlier := false
		for _, f := range res.Failures {
			if f.Category == "price_outlier" {
				foundOutlier = true
			}
		}
		if !foundOutlier {
			t.Error("Expected to find price outlier failure")
		}
	})

	t.Run("impossible coordinates", func(t *testing.T) {
		// Create a place record with lat/lon in US
		placeRecord := map[string]any{
			"schema_version": "1.0.0",
			"record_id":      "rec-impossible-coords",
			"record_type":    "place",
			"pipeline_stage": "ai_extracted",
			"source": map[string]any{
				"source_id":     "src-1",
				"source_type":   "rss",
				"canonical_url": "https://example.com",
			},
			"content": map[string]any{
				"canonical_name": "Out of Bounds Place",
				"place_type":     "restaurant",
				"coordinates": map[string]any{
					"latitude":  37.7749, // San Francisco lat
					"longitude": -122.4194,
				},
			},
			"provenance": map[string]any{
				"claims": []any{},
			},
		}

		data, _ := json.Marshal(placeRecord)
		res, err := v.ValidateRecord("llm-extraction-response.schema.json", wrapInResponse(data))
		if err != nil {
			t.Fatalf("Validator failed: %v", err)
		}

		foundGeoFailure := false
		for _, f := range res.Failures {
			if f.Category == "impossible_coordinates" {
				foundGeoFailure = true
			}
		}
		if !foundGeoFailure {
			t.Error("Expected to find impossible coordinates failure")
		}
	})
}

func wrapInResponse(recordData []byte) []byte {
	var record map[string]any
	json.Unmarshal(recordData, &record)
	response := map[string]any{
		"schema_version": "1.0.0",
		"request_id":     "req-test-123",
		"status":         "completed",
		"records":        []any{record},
	}
	wrapped, _ := json.Marshal(response)
	return wrapped
}
