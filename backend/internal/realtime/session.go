package realtime

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

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
	expectedSeq int
	utteranceID string
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
		s.emitError(err.Error())
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
		s.audioBuffer = make([]byte, 0)
		s.startTime = time.Now()

		s.onEvent(EventEnvelope{
			Version:   1,
			Type:      TypeTranscriptPartial,
			SessionID: s.SessionID,
			Sequence:  s.expectedSeq + 1,
			Timestamp: time.Now(),
			Payload: map[string]interface{}{
				"status": "listening",
			},
		})

	case TypeAudioChunk:
		if s.State != StateCapturing {
			return errors.New("cannot receive audio chunk when not capturing")
		}

		// Check duration limits
		if time.Since(s.startTime) > s.maxDuration {
			s.State = StateError
			err := errors.New("max utterance duration exceeded")
			s.emitError(err.Error())
			return err
		}

		// Check buffer limits
		if len(s.audioBuffer)+len(payloadData) > s.maxBufSize {
			s.State = StateError
			err := errors.New("max audio buffer limit exceeded")
			s.emitError(err.Error())
			return err
		}

		s.audioBuffer = append(s.audioBuffer, payloadData...)

	case TypePttEnded:
		if s.State != StateCapturing {
			return errors.New("cannot end PTT when not capturing")
		}
		s.State = StateEnding

		// Trigger ASR & Translation asynchronously
		go s.processUtterance(ctx, s.utteranceID, s.audioBuffer)

		s.State = StateReady

	case TypeSessionEnded:
		s.State = StateIdle
		s.onEvent(EventEnvelope{
			Version:   1,
			Type:      TypeSessionEnded,
			SessionID: s.SessionID,
			Sequence:  s.expectedSeq + 1,
			Timestamp: time.Now(),
		})

	default:
		slog.Warn("Unhandled event type", slog.String("type", env.Type))
	}

	return nil
}

func (s *Session) processUtterance(ctx context.Context, utteranceID string, audioData []byte) {
	// ASR Transcription
	transcript, err := s.asrProvider.Transcribe(ctx, fptai.AudioInput{Data: audioData})
	if err != nil {
		slog.Error("ASR failure", slog.Any("error", err))
		s.onEvent(EventEnvelope{
			Version:     1,
			Type:        TypeProviderDegraded,
			SessionID:   s.SessionID,
			UtteranceID: utteranceID,
			Timestamp:   time.Now(),
			Payload: map[string]interface{}{
				"provider": "asr",
				"error":    err.Error(),
			},
		})
		return
	}

	s.onEvent(EventEnvelope{
		Version:     1,
		Type:        TypeTranscriptFinal,
		SessionID:   s.SessionID,
		UtteranceID: utteranceID,
		Timestamp:   time.Now(),
		Payload: map[string]interface{}{
			"text": transcript.Text,
		},
	})

	// Translation (translation must not wait for downstream price engine)
	translation, err := s.transProvider.Translate(ctx, fptai.TranslationInput{
		Text:   transcript.Text,
		Target: "en", // default locale target
	})
	if err != nil {
		slog.Error("Translation failure", slog.Any("error", err))
		s.onEvent(EventEnvelope{
			Version:     1,
			Type:        TypeProviderDegraded,
			SessionID:   s.SessionID,
			UtteranceID: utteranceID,
			Timestamp:   time.Now(),
			Payload: map[string]interface{}{
				"provider": "translation",
				"error":    err.Error(),
			},
		})
		return
	}

	s.onEvent(EventEnvelope{
		Version:     1,
		Type:        TypeTranslationReady,
		SessionID:   s.SessionID,
		UtteranceID: utteranceID,
		Timestamp:   time.Now(),
		Payload: map[string]interface{}{
			"translated_text": translation.Text,
		},
	})
}

func (s *Session) emitError(msg string) {
	s.onEvent(EventEnvelope{
		Version:   1,
		Type:      TypeSessionError,
		SessionID: s.SessionID,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"error": msg,
		},
	})
}
