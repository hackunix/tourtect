package model

import (
	"encoding/json"
	"time"
)

const SessionVersion = 1

type ConsentState struct {
	Processing   bool `json:"processing"`
	Contribution bool `json:"contribution"`
	Publish      bool `json:"publish"`
}

type SessionContext struct {
	Locale             string       `json:"locale"`
	TargetLocale       string       `json:"target_locale,omitempty"`
	PlaceID            string       `json:"place_id,omitempty"`
	ApproximateRegion  string       `json:"approximate_region,omitempty"`
	InteractionMode    string       `json:"interaction_mode"`
	CurrentSafetyState string       `json:"current_safety_state,omitempty"`
	UserConfirmedFacts []string     `json:"user_confirmed_facts,omitempty"`
	ActiveCaptureIDs   []string     `json:"active_capture_ids,omitempty"`
	ConsentState       ConsentState `json:"consent_state"`
}

type Confirmation struct {
	ID          string    `json:"confirmation_id"`
	Action      string    `json:"action"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ExpiresAt   time.Time `json:"expires_at"`
	Consumed    bool      `json:"consumed"`
	Target      string    `json:"target,omitempty"`
}

type Session struct {
	ID                  string                   `json:"session_id"`
	UserID              string                   `json:"user_id"`
	Version             int                      `json:"version"`
	CreatedAt           time.Time                `json:"created_at"`
	UpdatedAt           time.Time                `json:"updated_at"`
	ExpiresAt           time.Time                `json:"expires_at"`
	Context             SessionContext           `json:"context"`
	RecentResponses     []Response               `json:"recent_responses,omitempty"`
	ProcessedMessageIDs []string                 `json:"processed_message_ids,omitempty"`
	Confirmations       map[string]*Confirmation `json:"confirmations,omitempty"`
}

type Route struct {
	Intent         string   `json:"intent"`
	Confidence     float64  `json:"confidence"`
	RequiredTools  []string `json:"required_tools"`
	MissingFields  []string `json:"missing_fields"`
	SafetyOverride bool     `json:"safety_override"`
}

type Message struct {
	ID            string
	InputType     string
	Text          string
	Locale        string
	PlaceID       string
	CaptureID     string
	UserConfirmed bool
	Structured    json.RawMessage
}

type Evidence struct {
	ID            string     `json:"evidence_id"`
	SourceType    string     `json:"source_type"`
	SourceID      string     `json:"source_id"`
	Title         string     `json:"title"`
	Summary       string     `json:"summary"`
	ObservedAt    *time.Time `json:"observed_at,omitempty"`
	Freshness     string     `json:"freshness"`
	EvidenceLevel string     `json:"evidence_level"`
	SourceURL     string     `json:"source_url,omitempty"`
}

type ToolResult struct {
	ID            string          `json:"tool_result_id"`
	ToolName      string          `json:"tool_name"`
	Status        string          `json:"status"`
	DurationMS    int64           `json:"duration_ms"`
	Output        json.RawMessage `json:"output"`
	ErrorCategory string          `json:"error_category,omitempty"`
}

type SuggestedAction struct {
	ID                   string `json:"action_id"`
	Label                string `json:"label"`
	Type                 string `json:"action_type"`
	Target               string `json:"target,omitempty"`
	RequiresConfirmation bool   `json:"requires_confirmation"`
}

type Response struct {
	ID                    string            `json:"assistant_message_id"`
	Intent                string            `json:"intent"`
	Message               string            `json:"message"`
	Confidence            float64           `json:"confidence"`
	Evidence              []Evidence        `json:"evidence"`
	ToolResults           []ToolResult      `json:"tool_results"`
	RequestedConfirmation *Confirmation     `json:"requested_confirmation,omitempty"`
	SuggestedActions      []SuggestedAction `json:"suggested_actions"`
	SafetyState           string            `json:"safety_state"`
	FactorsConsidered     []string          `json:"factors_considered"`
	MissingInformation    []string          `json:"missing_information"`
	Freshness             string            `json:"freshness,omitempty"`
	DatasetVersion        string            `json:"dataset_version,omitempty"`
	FallbackUsed          bool              `json:"fallback_used"`
	TraceID               string            `json:"trace_id"`
}

type Trace struct {
	TraceID         string   `json:"trace_id"`
	SessionID       string   `json:"session_id"`
	Intent          string   `json:"intent"`
	ToolNames       []string `json:"tool_names"`
	ToolDurationsMS []int64  `json:"tool_durations_ms"`
	Provider        string   `json:"provider,omitempty"`
	ModelVersion    string   `json:"model_version,omitempty"`
	PolicyVersion   string   `json:"policy_version"`
	RetrievalCount  int      `json:"retrieval_count"`
	EvidenceIDs     []string `json:"evidence_ids"`
	Outcome         string   `json:"outcome"`
	ErrorCategory   string   `json:"error_category,omitempty"`
	FallbackUsed    bool     `json:"fallback_used"`
}
