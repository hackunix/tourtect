package realtime

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEventEnvelopeUsesSnakeCaseAndAcceptsLegacyControlFields(t *testing.T) {
	var incoming EventEnvelope
	if err := json.Unmarshal([]byte(`{"type":"ptt.started","sessionId":"session-1","utteranceId":"utterance-1","traceId":"trace-1","sequence":1}`), &incoming); err != nil {
		t.Fatalf("unmarshal legacy envelope: %v", err)
	}
	if incoming.Type != TypePttStarted || incoming.SessionID != "session-1" || incoming.UtteranceID != "utterance-1" || incoming.TraceID != "trace-1" {
		t.Fatalf("legacy envelope was not preserved: %+v", incoming)
	}

	outgoing, err := json.Marshal(newEventEnvelope(TypeIntentDetected, "session-1", "utterance-1", "trace-1", 2, nil))
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}
	jsonText := string(outgoing)
	for _, field := range []string{`"event_id"`, `"event_type"`, `"session_id"`, `"utterance_id"`, `"occurred_at"`, `"trace_id"`} {
		if !strings.Contains(jsonText, field) {
			t.Fatalf("missing snake_case field %s in %s", field, jsonText)
		}
	}
	for _, legacy := range []string{`"sessionId"`, `"utteranceId"`, `"traceId"`} {
		if strings.Contains(jsonText, legacy) {
			t.Fatalf("legacy field %s leaked into server event %s", legacy, jsonText)
		}
	}
}

func TestAssistantEventConstantsAreReservedWithoutSyntheticPayloads(t *testing.T) {
	if TypeIntentDetected != "intent.detected" || TypeAssistantSuggestion != "assistant.suggestion" {
		t.Fatal("assistant realtime event names changed unexpectedly")
	}
}
