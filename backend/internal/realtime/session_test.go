package realtime

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/tourtect/backend/adapters/fptai"
)

type stubASR struct {
	text string
	err  error
}

func (s stubASR) Transcribe(context.Context, fptai.AudioInput) (fptai.Transcript, error) {
	return fptai.Transcript{Text: s.text}, s.err
}

type stubTranslation struct {
	text string
	err  error
}

func (s stubTranslation) Translate(context.Context, fptai.TranslationInput) (fptai.Translation, error) {
	return fptai.Translation{Text: s.text}, s.err
}

func TestRealtimeSession(t *testing.T) {
	t.Run("PTT state transition flow", func(t *testing.T) {
		events := make(chan EventEnvelope, 8)
		s := NewSession(
			"sess-1",
			stubASR{text: "Chào anh, taxi giá bao nhiêu?"},
			stubTranslation{text: "Hello, how much is the taxi?"},
			func(env EventEnvelope) { events <- env },
		)

		err := s.HandleEvent(context.Background(), EventEnvelope{
			Version:     1,
			Type:        TypePttStarted,
			SessionID:   "sess-1",
			UtteranceID: "utt-1",
			Sequence:    1,
			Timestamp:   time.Now(),
			TraceID:     "trace-1",
		}, nil)
		if err != nil {
			t.Fatalf("start PTT: %v", err)
		}
		if s.State != StateCapturing {
			t.Fatalf("expected capturing state, got %s", s.State)
		}

		err = s.HandleEvent(context.Background(), EventEnvelope{
			Version:   1,
			Type:      TypeAudioChunk,
			SessionID: "sess-1",
			Sequence:  2,
			Timestamp: time.Now(),
		}, []byte{0x01, 0x02, 0x03})
		if err != nil {
			t.Fatalf("handle audio: %v", err)
		}

		err = s.HandleEvent(context.Background(), EventEnvelope{
			Version:     1,
			Type:        TypePttEnded,
			SessionID:   "sess-1",
			UtteranceID: "utt-1",
			Sequence:    3,
			Timestamp:   time.Now(),
		}, nil)
		if err != nil {
			t.Fatalf("end PTT: %v", err)
		}

		seen := collectUntil(t, events, TypeTranslationReady)
		if got := seen[TypeTranscriptFinal].Payload["text"]; got != "Chào anh, taxi giá bao nhiêu?" {
			t.Fatalf("unexpected transcript: %v", got)
		}
		if got := seen[TypeTranslationReady].Payload["translated_text"]; got != "Hello, how much is the taxi?" {
			t.Fatalf("unexpected translation: %v", got)
		}
		for _, eventType := range []string{TypeTranscriptPartial, TypeTranscriptFinal, TypeTranslationReady} {
			event := seen[eventType]
			if event.EventID == "" || event.TraceID != "trace-1" || event.SessionID != "sess-1" || event.UtteranceID != "utt-1" {
				t.Fatalf("%s missing trace identifiers: %+v", eventType, event)
			}
		}
	})

	t.Run("rejects invalid sequence numbers", func(t *testing.T) {
		events := make(chan EventEnvelope, 4)
		s := NewSession("sess-2", stubASR{}, stubTranslation{}, func(env EventEnvelope) { events <- env })
		_ = s.HandleEvent(context.Background(), EventEnvelope{
			Version:     1,
			Type:        TypePttStarted,
			SessionID:   "sess-2",
			UtteranceID: "utt-2",
			Sequence:    1,
			TraceID:     "trace-2",
		}, nil)

		err := s.HandleEvent(context.Background(), EventEnvelope{
			Version:   1,
			Type:      TypeAudioChunk,
			SessionID: "sess-2",
			Sequence:  3,
		}, []byte{0x01})
		if err == nil {
			t.Fatal("expected out-of-order sequence error")
		}
		if s.State != StateError {
			t.Fatalf("expected error state, got %s", s.State)
		}

		seen := collectUntil(t, events, TypeSessionError)
		message, _ := seen[TypeSessionError].Payload["error"].(string)
		if !strings.Contains(message, "out-of-order sequence") {
			t.Fatalf("unexpected session error: %q", message)
		}
	})

	t.Run("ignores unknown event types without crashing", func(t *testing.T) {
		events := make(chan EventEnvelope, 1)
		s := NewSession("sess-unknown", stubASR{}, stubTranslation{}, func(env EventEnvelope) { events <- env })
		if err := s.HandleEvent(context.Background(), EventEnvelope{
			Type:      "future.event",
			SessionID: "sess-unknown",
			Sequence:  1,
			TraceID:   "trace-unknown",
		}, nil); err != nil {
			t.Fatalf("unknown event must be tolerated: %v", err)
		}
		if s.State != StateReady {
			t.Fatalf("unknown event changed session state to %s", s.State)
		}
		select {
		case event := <-events:
			t.Fatalf("unknown event emitted synthetic output: %+v", event)
		default:
		}
	})
}

func TestProviderDegradationIsSanitized(t *testing.T) {
	tests := []struct {
		name             string
		asr              fptai.ASRProvider
		translation      fptai.TranslationProvider
		expectedProvider string
	}{
		{
			name:             "ASR unavailable",
			asr:              fptai.NewUnavailableASR(),
			translation:      stubTranslation{text: "unused"},
			expectedProvider: "asr",
		},
		{
			name:             "translation internal error",
			asr:              stubASR{text: "hello"},
			translation:      stubTranslation{err: errors.New("secret upstream detail")},
			expectedProvider: "translation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events := make(chan EventEnvelope, 8)
			s := NewSession("sess-degraded", tt.asr, tt.translation, func(env EventEnvelope) { events <- env })
			_ = s.HandleEvent(context.Background(), EventEnvelope{
				Type:        TypePttStarted,
				SessionID:   "sess-degraded",
				UtteranceID: "utt-degraded",
				Sequence:    1,
				TraceID:     "trace-degraded",
			}, nil)
			_ = s.HandleEvent(context.Background(), EventEnvelope{
				Type:      TypeAudioChunk,
				SessionID: "sess-degraded",
				Sequence:  2,
			}, []byte{0x01})
			_ = s.HandleEvent(context.Background(), EventEnvelope{
				Type:        TypePttEnded,
				SessionID:   "sess-degraded",
				UtteranceID: "utt-degraded",
				Sequence:    3,
			}, nil)

			degraded := collectUntil(t, events, TypeProviderDegraded)[TypeProviderDegraded]
			if degraded.EventID == "" || degraded.TraceID != "trace-degraded" || degraded.UtteranceID != "utt-degraded" {
				t.Fatalf("degraded event missing identifiers: %+v", degraded)
			}
			if degraded.Payload["provider"] != tt.expectedProvider {
				t.Fatalf("expected provider %s, got %v", tt.expectedProvider, degraded.Payload["provider"])
			}
			if degraded.Payload["error_category"] != "provider_unavailable" {
				t.Fatalf("unexpected category: %v", degraded.Payload["error_category"])
			}
			if _, exposed := degraded.Payload["error"]; exposed {
				t.Fatal("provider error details must not be exposed")
			}
			if strings.Contains(degraded.Payload["message"].(string), "secret upstream detail") {
				t.Fatal("provider error details leaked through sanitized message")
			}
		})
	}
}

func collectUntil(t *testing.T, events <-chan EventEnvelope, terminalType string) map[string]EventEnvelope {
	t.Helper()
	seen := make(map[string]EventEnvelope)
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	for {
		select {
		case event := <-events:
			seen[event.Type] = event
			if event.Type == terminalType {
				return seen
			}
		case <-timer.C:
			t.Fatalf("timed out waiting for %s; saw %v", terminalType, seen)
		}
	}
}
