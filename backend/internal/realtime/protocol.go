package realtime

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventEnvelope struct {
	Version     int                    `json:"version"`
	EventID     string                 `json:"event_id"`
	Type        string                 `json:"event_type"`
	SessionID   string                 `json:"session_id"`
	UtteranceID string                 `json:"utterance_id,omitempty"`
	Sequence    int                    `json:"sequence"`
	Timestamp   time.Time              `json:"occurred_at"`
	TraceID     string                 `json:"trace_id"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
}

// UnmarshalJSON accepts the snake_case protocol and the legacy camelCase
// control envelope so deployed clients can migrate without losing PTT input.
func (e *EventEnvelope) UnmarshalJSON(data []byte) error {
	var wire struct {
		Version           int                    `json:"version"`
		EventID           string                 `json:"event_id"`
		LegacyEventID     string                 `json:"eventId"`
		Type              string                 `json:"event_type"`
		LegacyType        string                 `json:"type"`
		SessionID         string                 `json:"session_id"`
		LegacySessionID   string                 `json:"sessionId"`
		UtteranceID       string                 `json:"utterance_id"`
		LegacyUtteranceID string                 `json:"utteranceId"`
		Sequence          int                    `json:"sequence"`
		Timestamp         time.Time              `json:"occurred_at"`
		LegacyTimestamp   time.Time              `json:"timestamp"`
		TraceID           string                 `json:"trace_id"`
		LegacyTraceID     string                 `json:"traceId"`
		Payload           map[string]interface{} `json:"payload"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}

	e.Version = wire.Version
	e.EventID = firstNonEmpty(wire.EventID, wire.LegacyEventID)
	e.Type = firstNonEmpty(wire.Type, wire.LegacyType)
	e.SessionID = firstNonEmpty(wire.SessionID, wire.LegacySessionID)
	e.UtteranceID = firstNonEmpty(wire.UtteranceID, wire.LegacyUtteranceID)
	e.Sequence = wire.Sequence
	e.Timestamp = wire.Timestamp
	if e.Timestamp.IsZero() {
		e.Timestamp = wire.LegacyTimestamp
	}
	e.TraceID = firstNonEmpty(wire.TraceID, wire.LegacyTraceID)
	e.Payload = wire.Payload
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func newEventEnvelope(eventType, sessionID, utteranceID, traceID string, sequence int, payload map[string]interface{}) EventEnvelope {
	if traceID == "" {
		traceID = uuid.NewString()
	}
	return EventEnvelope{
		Version:     1,
		EventID:     uuid.NewString(),
		Type:        eventType,
		SessionID:   sessionID,
		UtteranceID: utteranceID,
		Sequence:    sequence,
		Timestamp:   time.Now().UTC(),
		TraceID:     traceID,
		Payload:     payload,
	}
}

// Client-to-server types
const (
	TypeSessionResume = "session.resume"
	TypePttStarted    = "ptt.started"
	TypeAudioChunk    = "audio.chunk"
	TypePttEnded      = "ptt.ended"
	TypeSessionEnded  = "session.ended"
	TypeClientAck     = "client.ack"
)

// Server-to-client types
const (
	TypeSessionReady        = "session.ready"
	TypeTranscriptPartial   = "transcript.partial"
	TypeTranscriptFinal     = "transcript.final"
	TypeTranslationReady    = "translation.ready"
	TypeIntentDetected      = "intent.detected"
	TypePriceInsightEvent   = "price.insight"
	TypeSafetyAlertEvent    = "safety.alert"
	TypeAssistantSuggestion = "assistant.suggestion"
	TypeProviderDegraded    = "provider.degraded"
	TypeSessionError        = "session.error"
)
