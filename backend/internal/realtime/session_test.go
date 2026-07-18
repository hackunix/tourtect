package realtime

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/tourtect/backend/adapters/fptai"
)

func TestRealtimeSession(t *testing.T) {
	ctx := context.Background()

	// Initialize fake providers
	fakeAsr := &fptai.FakeASR{MockText: "Chào anh, taxi giá bao nhiêu?"}
	fakeTrans := &fptai.FakeTranslation{MockText: "Hello, how much is the taxi?"}

	events := make([]EventEnvelope, 0)
	onEvent := func(env EventEnvelope) {
		events = append(events, env)
	}

	t.Run("PTT state transition flow", func(t *testing.T) {
		events = nil
		s := NewSession("sess-1", fakeAsr, fakeTrans, onEvent)

		// 1. Start PTT capture
		err := s.HandleEvent(ctx, EventEnvelope{
			Version:     1,
			Type:        TypePttStarted,
			SessionID:   "sess-1",
			UtteranceID: "utt-1",
			Sequence:    1,
			Timestamp:   time.Now(),
		}, nil)
		if err != nil {
			t.Errorf("Unexpected error starting PTT: %v", err)
		}
		if s.State != StateCapturing {
			t.Errorf("Expected state to be Capturing, got %s", s.State)
		}

		// 2. Stream audio binary data
		err = s.HandleEvent(ctx, EventEnvelope{
			Version:   1,
			Type:      TypeAudioChunk,
			SessionID: "sess-1",
			Sequence:  2,
			Timestamp: time.Now(),
		}, []byte{0x01, 0x02, 0x03})
		if err != nil {
			t.Errorf("Unexpected error handling audio chunk: %v", err)
		}

		// 3. End PTT capture
		err = s.HandleEvent(ctx, EventEnvelope{
			Version:     1,
			Type:        TypePttEnded,
			SessionID:   "sess-1",
			UtteranceID: "utt-1",
			Sequence:    3,
			Timestamp:   time.Now(),
		}, nil)
		if err != nil {
			t.Errorf("Unexpected error ending PTT: %v", err)
		}

		// Sleep briefly to allow async processing to complete
		time.Sleep(50 * time.Millisecond)

		// Verify event triggers
		var hasTranscript, hasTranslation bool
		for _, e := range events {
			if e.Type == TypeTranscriptFinal {
				hasTranscript = true
				txt := e.Payload["text"].(string)
				if txt != "Chào anh, taxi giá bao nhiêu?" {
					t.Errorf("Expected transcript text 'Chào anh, taxi giá bao nhiêu?', got '%s'", txt)
				}
			}
			if e.Type == TypeTranslationReady {
				hasTranslation = true
				trans := e.Payload["translated_text"].(string)
				if trans != "Hello, how much is the taxi?" {
					t.Errorf("Expected translated text 'Hello, how much is the taxi?', got '%s'", trans)
				}
			}
		}

		if !hasTranscript {
			t.Errorf("Expected Final Transcript event to be emitted")
		}
		if !hasTranslation {
			t.Errorf("Expected Translation Ready event to be emitted")
		}
	})

	t.Run("Rejects invalid sequence numbers", func(t *testing.T) {
		events = nil
		s := NewSession("sess-2", fakeAsr, fakeTrans, onEvent)

		// Start sequence at 1
		_ = s.HandleEvent(ctx, EventEnvelope{
			Version:     1,
			Type:        TypePttStarted,
			SessionID:   "sess-2",
			UtteranceID: "utt-2",
			Sequence:    1,
			Timestamp:   time.Now(),
		}, nil)

		// Skip sequence number 2 and send sequence number 3
		err := s.HandleEvent(ctx, EventEnvelope{
			Version:   1,
			Type:      TypeAudioChunk,
			SessionID: "sess-2",
			Sequence:  3, // expects 2
			Timestamp: time.Now(),
		}, []byte{0x01})

		if err == nil {
			t.Errorf("Expected error when receiving out-of-order sequence, got nil")
		}
		if s.State != StateError {
			t.Errorf("Expected session state to transition to error, got %s", s.State)
		}
		
		var hasErrorEvent bool
		for _, e := range events {
			if e.Type == TypeSessionError {
				hasErrorEvent = true
				msg := e.Payload["error"].(string)
				if !strings.Contains(msg, "out-of-order sequence") {
					t.Errorf("Unexpected error message: %s", msg)
				}
			}
		}
		if !hasErrorEvent {
			t.Errorf("Expected SessionError event to be emitted")
		}
	})
}
