package realtime

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/tourtect/backend/adapters/fptai"
)

type SessionState string

const (
	StateIdle      SessionState = "idle"
	StateReady     SessionState = "ready"
	StateCapturing SessionState = "capturing"
	StateEnding    SessionState = "ending"
	StateError     SessionState = "error"
)

type Session struct {
	SessionID   string
	State       SessionState
	mu          sync.Mutex
	eventMu     sync.Mutex
	expectedSeq int
	serverSeq   int
	utteranceID string
	traceID     string
	audioBuffer []byte
	maxDuration time.Duration
	maxBufSize  int
	startTime   time.Time

	asrProvider   fptai.ASRProvider
	transProvider fptai.TranslationProvider
	onEvent       func(EventEnvelope)
}

func NewSession(id string, asr fptai.ASRProvider, trans fptai.TranslationProvider, onEvent func(EventEnvelope)) *Session {
	return &Session{
		SessionID:     id,
		State:         StateReady,
		maxDuration:   60 * time.Second,
		maxBufSize:    5 * 1024 * 1024, // 5MB limit
		asrProvider:   asr,
		transProvider: trans,
		onEvent:       onEvent,
	}
}

func (s *Session) HandleEvent(ctx context.Context, env EventEnvelope, payloadData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check sequence monotonic ordering
	if env.Sequence != s.expectedSeq+1 {
		s.State = StateError
		err := fmt.Errorf("out-of-order sequence: expected %d, got %d", s.expectedSeq+1, env.Sequence)
		s.emitError(err.Error(), firstNonEmpty(env.TraceID, s.traceID))
		return err
	}
	s.expectedSeq = env.Sequence

	switch env.Type {
	case TypePttStarted:
		if s.State != StateReady {
			return errors.New("cannot start PTT when not in ready state")
		}
		s.State = StateCapturing
		s.utteranceID = env.UtteranceID
		if s.utteranceID == "" {
			s.utteranceID = uuid.NewString()
		}
		s.traceID = env.TraceID
		if s.traceID == "" {
			s.traceID = uuid.NewString()
		}
		s.audioBuffer = make([]byte, 0)
		s.startTime = time.Now()

		s.emit(TypeTranscriptPartial, s.utteranceID, s.traceID, map[string]interface{}{
			"status": "listening",
		})

	case TypeAudioChunk:
		if s.State != StateCapturing {
			return errors.New("cannot receive audio chunk when not capturing")
		}

		// Check duration limits
		if time.Since(s.startTime) > s.maxDuration {
			s.State = StateError
			err := errors.New("max utterance duration exceeded")
			s.emitError(err.Error(), s.traceID)
			return err
		}

		// Check buffer limits
		if len(s.audioBuffer)+len(payloadData) > s.maxBufSize {
			s.State = StateError
			err := errors.New("max audio buffer limit exceeded")
			s.emitError(err.Error(), s.traceID)
			return err
		}

		s.audioBuffer = append(s.audioBuffer, payloadData...)

	case TypePttEnded:
		if s.State != StateCapturing {
			return errors.New("cannot end PTT when not capturing")
		}
		s.State = StateEnding

		// Copy the bounded buffer before processing so a subsequent utterance
		// cannot mutate data still owned by the provider call.
		audioData := append([]byte(nil), s.audioBuffer...)
		go s.processUtterance(ctx, s.utteranceID, s.traceID, audioData)

		s.State = StateReady

	case TypeSessionEnded:
		s.State = StateIdle
		s.emit(TypeSessionEnded, s.utteranceID, s.traceID, nil)

	default:
		slog.Warn("Unhandled event type", slog.String("type", env.Type))
	}

	return nil
}

func (s *Session) processUtterance(ctx context.Context, utteranceID, traceID string, audioData []byte) {
	// ASR Transcription
	transcript, err := s.asrProvider.Transcribe(ctx, fptai.AudioInput{Data: audioData})
	if err != nil {
		slog.Error("ASR failure", slog.Any("error", err))
		s.emitProviderDegraded("asr", utteranceID, traceID)
		return
	}

	s.emit(TypeTranscriptFinal, utteranceID, traceID, map[string]interface{}{
		"text": transcript.Text,
	})

	// Translation (translation must not wait for downstream price engine)
	translation, err := s.transProvider.Translate(ctx, fptai.TranslationInput{
		Text:   transcript.Text,
		Target: "en", // default locale target
	})
	if err != nil {
		slog.Error("Translation failure", slog.Any("error", err))
		s.emitProviderDegraded("translation", utteranceID, traceID)
		return
	}

	s.emit(TypeTranslationReady, utteranceID, traceID, map[string]interface{}{
		"translated_text": translation.Text,
	})
}

func (s *Session) emitError(msg, traceID string) {
	s.emit(TypeSessionError, s.utteranceID, traceID, map[string]interface{}{
		"error": msg,
	})
}

func (s *Session) emitProviderDegraded(provider, utteranceID, traceID string) {
	s.emit(TypeProviderDegraded, utteranceID, traceID, map[string]interface{}{
		"provider":       provider,
		"error_category": "provider_unavailable",
		"message":        "AI processing is currently unavailable.",
	})
}

func (s *Session) emit(eventType, utteranceID, traceID string, payload map[string]interface{}) {
	s.onEvent(s.nextEvent(eventType, utteranceID, traceID, payload))
}

func (s *Session) nextEvent(eventType, utteranceID, traceID string, payload map[string]interface{}) EventEnvelope {
	s.eventMu.Lock()
	defer s.eventMu.Unlock()
	s.serverSeq++
	return newEventEnvelope(eventType, s.SessionID, utteranceID, traceID, s.serverSeq, payload)
}

func (s *Session) nextClientSequence() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.expectedSeq + 1
}
