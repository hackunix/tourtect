package validation

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/xeipuuv/gojsonschema"
)

type ValidationFailure struct {
	FieldPath string `json:"field_path"`
	Category  string `json:"category"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

type ValidationResult struct {
	ValidationID string              `json:"validation_id"`
	RecordID     string              `json:"record_id"`
	Status       string              `json:"status"` // valid, valid_with_warnings, quarantined, rejected
	ValidatedAt  time.Time           `json:"validated_at"`
	Failures     []ValidationFailure `json:"failures"`
	Warnings     []string            `json:"warnings"`
	Metadata     map[string]any      `json:"metadata,omitempty"`
}

var (
	phoneRegex = regexp.MustCompile(`\+?[0-9]{3}-?[0-9]{3,4}-?[0-9]{4}`)
	emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
)

type PipelineValidator struct {
	contractsDir string
}

func NewPipelineValidator(contractsDir string) *PipelineValidator {
	absDir, err := filepath.Abs(contractsDir)
	if err != nil {
		absDir = contractsDir
	}
	return &PipelineValidator{contractsDir: absDir}
}

func (v *PipelineValidator) ValidateRecord(schemaFileName string, rawJSON []byte) (*ValidationResult, error) {
	result := &ValidationResult{
		ValidationID: fmt.Sprintf("val-%d", time.Now().UnixNano()),
		ValidatedAt:  time.Now(),
		Failures:     []ValidationFailure{},
		Warnings:     []string{},
		Metadata:     make(map[string]any),
	}

	// 1. JSON Parse
	var recordMap map[string]any
	if err := json.Unmarshal(rawJSON, &recordMap); err != nil {
		result.Status = "rejected"
		result.Failures = append(result.Failures, ValidationFailure{
			FieldPath: "",
			Category:  "json_schema",
			Code:      "invalid_json",
			Message:   fmt.Sprintf("JSON parsing failed: %v", err),
		})
		return result, nil
	}

	if recordID, ok := recordMap["record_id"].(string); ok {
		result.RecordID = recordID
	} else {
		result.RecordID = "unknown"
	}

	// 2. JSON Schema Validation using gojsonschema
	sl := gojsonschema.NewSchemaLoader()
	files, err := os.ReadDir(v.contractsDir)
	if err != nil {
		result.Status = "rejected"
		result.Failures = append(result.Failures, ValidationFailure{
			FieldPath: "",
			Category:  "json_schema",
			Code:      "contracts_dir_error",
			Message:   fmt.Sprintf("Failed to read contracts dir: %v", err),
		})
		return result, nil
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".schema.json") || f.Name() == schemaFileName {
			continue
		}
		p := filepath.Join(v.contractsDir, f.Name())
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var schemaMap map[string]any
		if err := json.Unmarshal(data, &schemaMap); err == nil {
			if id, ok := schemaMap["$id"].(string); ok {
				sl.AddSchema(id, gojsonschema.NewBytesLoader(data))
			}
		}
	}

	targetPath := filepath.Join(v.contractsDir, schemaFileName)
	targetData, err := os.ReadFile(targetPath)
	if err != nil {
		result.Status = "rejected"
		result.Failures = append(result.Failures, ValidationFailure{
			FieldPath: "",
			Category:  "json_schema",
			Code:      "schema_loading_error",
			Message:   fmt.Sprintf("JSON Schema loading error: %v", err),
		})
		return result, nil
	}

	schema, err := sl.Compile(gojsonschema.NewBytesLoader(targetData))
	if err != nil {
		result.Status = "rejected"
		result.Failures = append(result.Failures, ValidationFailure{
			FieldPath: "",
			Category:  "json_schema",
			Code:      "schema_compile_error",
			Message:   fmt.Sprintf("JSON Schema compile error: %v", err),
		})
		return result, nil
	}

	documentLoader := gojsonschema.NewBytesLoader(rawJSON)
	schemaResult, err := schema.Validate(documentLoader)
	if err != nil {
		result.Status = "rejected"
		result.Failures = append(result.Failures, ValidationFailure{
			FieldPath: "",
			Category:  "json_schema",
			Code:      "schema_validation_error",
			Message:   fmt.Sprintf("Schema validation failed: %v", err),
		})
		return result, nil
	}

	if !schemaResult.Valid() {
		result.Status = "rejected"
		for _, desc := range schemaResult.Errors() {
			result.Failures = append(result.Failures, ValidationFailure{
				FieldPath: desc.Field(),
				Category:  "json_schema",
				Code:      desc.Type(),
				Message:   desc.Description(),
			})
		}
		return result, nil
	}

	// 3. Business Policy / Safety Validation
	if _, isSingleRecord := recordMap["record_type"]; isSingleRecord {
		v.runSafetyChecks(recordMap, result)
	} else if records, ok := recordMap["records"].([]any); ok {
		for _, recAny := range records {
			if rec, ok := recAny.(map[string]any); ok {
				v.runSafetyChecks(rec, result)
			}
		}
	}

	// Set status based on failures/warnings
	if len(result.Failures) > 0 {
		// Determine if quarantined or rejected based on failure types
		hasHighRisk := false
		for _, f := range result.Failures {
			if f.Category == "hallucinated_hotline_risk" || f.Category == "unsupported_accusation" || f.Category == "rights_violation" || f.Category == "personal_data_detected" {
				hasHighRisk = true
			}
		}
		if hasHighRisk {
			result.Status = "quarantined"
		} else {
			result.Status = "rejected"
		}
	} else if len(result.Warnings) > 0 {
		result.Status = "valid_with_warnings"
	} else {
		result.Status = "valid"
	}

	return result, nil
}

func (v *PipelineValidator) runSafetyChecks(record map[string]any, result *ValidationResult) {
	recordType, _ := record["record_type"].(string)

	// A. Check for PII (Personal Data Detected) in content/evidence
	content, hasContent := record["content"].(map[string]any)
	if hasContent {
		for k, val := range content {
			strVal, ok := val.(string)
			if !ok {
				continue
			}
			// Check for unredacted email
			if emailRegex.MatchString(strVal) {
				result.Failures = append(result.Failures, ValidationFailure{
					FieldPath: fmt.Sprintf("/content/%s", k),
					Category:  "personal_data_detected",
					Code:      "pii_email_leaked",
					Message:   "Unredacted email address detected in text field.",
				})
			}
			// Check for unredacted phone (excluding emergency numbers)
			if recordType != "emergency_directory_entry" && phoneRegex.MatchString(strVal) {
				result.Failures = append(result.Failures, ValidationFailure{
					FieldPath: fmt.Sprintf("/content/%s", k),
					Category:  "personal_data_detected",
					Code:      "pii_phone_leaked",
					Message:   "Potential unredacted phone number detected in non-emergency record.",
				})
			}
		}
	}

	// B. Content-specific checks
	switch recordType {
	case "emergency_directory_entry":
		// Hallucinated hotline risk: check provenance claim confidence
		provenance, ok := record["provenance"].(map[string]any)
		if ok {
			claims, ok := provenance["claims"].([]any)
			if ok {
				for i, claimAny := range claims {
					claim, ok := claimAny.(map[string]any)
					if !ok {
						continue
					}
					claimType, _ := claim["claim_type"].(string)
					confVal, _ := claim["confidence"].(float64)
					if claimType == "inferred" || confVal < 0.70 {
						result.Failures = append(result.Failures, ValidationFailure{
							FieldPath: fmt.Sprintf("/provenance/claims/%d", i),
							Category:  "hallucinated_hotline_risk",
							Code:      "low_confidence_inferred_hotline",
							Message:   "Safety-critical emergency number must have explicit provenance with confidence >= 0.70.",
						})
					}
				}
			}
		}

	case "scam_pattern":
		// Unsupported accusation: check pattern_name or common_scenario for specific business accusations
		scenario, _ := content["common_scenario"].(string)
		patternName, _ := content["pattern_name"].(string)
		for _, text := range []string{scenario, patternName} {
			if text == "" {
				continue
			}
			// Simple check for specific company accusations (e.g., claiming a named merchant is a scam)
			lowerText := strings.ToLower(text)
			if strings.Contains(lowerText, "is a scam") || strings.Contains(lowerText, "stole") || strings.Contains(lowerText, "cheat") {
				result.Failures = append(result.Failures, ValidationFailure{
					FieldPath: "/content/common_scenario",
					Category:  "unsupported_accusation",
					Code:      "accusation_guilt_label",
					Message:   "Scam patterns must describe behaviors, not declare legal guilt or accuse named entities.",
				})
			}
		}

	case "price_observation":
		// Check for impossible prices (price outlier)
		amountMap, ok := content["amount"].(map[string]any)
		if ok {
			minorUnits, _ := amountMap["minor_units"].(string)
			currency, _ := amountMap["currency"].(string)
			if currency == "VND" {
				var amount int64
				fmt.Sscanf(minorUnits, "%d", &amount)
				if amount > 10000000 { // 10 million VND limit for normal observations
					result.Failures = append(result.Failures, ValidationFailure{
						FieldPath: "/content/amount/minor_units",
						Category:  "price_outlier",
						Code:      "excessive_vnd_price",
						Message:   "Extracted taxi price exceeds theMateriality Threshold limit of 10,000,000 VND.",
					})
				}
			}
		}

	case "place":
		// Check for impossible coordinates
		coords, ok := content["coordinates"].(map[string]any)
		if ok {
			lat, _ := coords["latitude"].(float64)
			lon, _ := coords["longitude"].(float64)
			// Vietnam bounds: lat 8.0 to 24.0, lon 102.0 to 110.0
			if lat < 8.0 || lat > 24.0 || lon < 102.0 || lon > 110.0 {
				result.Failures = append(result.Failures, ValidationFailure{
					FieldPath: "/content/coordinates",
					Category:  "impossible_coordinates",
					Code:      "outside_vietnam_bounds",
					Message:   fmt.Sprintf("Coordinates (lat: %.4f, lon: %.4f) are outside Vietnam target region bounds.", lat, lon),
				})
			}
		}
	}

	// C. Date checks (Future-date anomaly)
	source, ok := record["source"].(map[string]any)
	if ok {
		pubTimeStr, _ := source["published_at"].(string)
		if pubTimeStr != "" {
			pubTime, err := time.Parse(time.RFC3339, pubTimeStr)
			if err == nil && pubTime.After(time.Now().Add(5*time.Minute)) {
				result.Failures = append(result.Failures, ValidationFailure{
					FieldPath: "/source/published_at",
					Category:  "invalid_date",
					Code:      "future_date_anomaly",
					Message:   fmt.Sprintf("Publication date %s is in the future.", pubTimeStr),
				})
			}
		}
	}

	// D. Robots/Rights checks
	policy, ok := record["policy"].(map[string]any)
	if ok {
		robotsStatus, _ := policy["robots_status"].(string)
		if robotsStatus == "disallowed" {
			result.Failures = append(result.Failures, ValidationFailure{
				FieldPath: "/policy/robots_status",
				Category:  "robots_violation",
				Code:      "robots_disallowed",
				Message:   "Source content violates robots.txt disallow directives.",
			})
		}
		rightsStatus, _ := policy["rights_status"].(string)
		if rightsStatus == "proprietary" {
			result.Failures = append(result.Failures, ValidationFailure{
				FieldPath: "/policy/rights_status",
				Category:  "rights_violation",
				Code:      "unlicensed_proprietary_source",
				Message:   "Cannot store proprietary raw content without explicit licensing agreement.",
			})
		}
	}
}

func IsValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	// Reject private IP / loopback to prevent SSRF
	host := u.Hostname()
	if host == "localhost" || host == "127.0.0.1" || strings.HasPrefix(host, "192.168.") || strings.HasPrefix(host, "10.") {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
