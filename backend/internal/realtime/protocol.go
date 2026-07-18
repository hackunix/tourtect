package realtime

import (
	"time"
)

type EventEnvelope struct {
	Version     int                    `json:"version"`
	Type        string                 `json:"type"`
	SessionID   string                 `json:"sessionId"`
	UtteranceID string                 `json:"utteranceId,omitempty"`
	Sequence    int                    `json:"sequence"`
	Timestamp   time.Time              `json:"timestamp"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
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
	TypeSessionReady      = "session.ready"
	TypeTranscriptPartial = "transcript.partial"
	TypeTranscriptFinal   = "transcript.final"
	TypeTranslationReady  = "translation.ready"
	TypePriceInsightEvent = "price.insight"
	TypeSafetyAlertEvent  = "safety.alert"
	TypeProviderDegraded  = "provider.degraded"
	TypeSessionError      = "session.error"
)
