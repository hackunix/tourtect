package realtime

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"

	"github.com/tourtect/backend/adapters/fptai"
)

type Handler struct {
	asrProvider   fptai.ASRProvider
	transProvider fptai.TranslationProvider
	mu            sync.RWMutex
	sessions      map[string]*Session
}

func NewHandler(asr fptai.ASRProvider, trans fptai.TranslationProvider) *Handler {
	return &Handler{
		asrProvider:   asr,
		transProvider: trans,
		sessions:      make(map[string]*Session),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection to websocket
	opts := &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Dev support
	}
	c, err := websocket.Accept(w, r, opts)
	if err != nil {
		slog.Error("WebSocket upgrade failed", slog.Any("error", err))
		return
	}
	defer c.Close(websocket.StatusInternalError, "the server closed the connection")

	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	// Create sync function to write messages to client safely
	writeMu := sync.Mutex{}
	onEvent := func(env EventEnvelope) {
		writeMu.Lock()
		defer writeMu.Unlock()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		wBytes, err := json.Marshal(env)
		if err != nil {
			slog.Error("failed to marshal server event", slog.Any("error", err))
			return
		}

		err = c.Write(ctx, websocket.MessageText, wBytes)
		if err != nil {
			slog.Error("failed to write websocket event to client", slog.Any("error", err))
		}
	}

	session := NewSession(sessionID, h.asrProvider, h.transProvider, onEvent)

	h.mu.Lock()
	h.sessions[sessionID] = session
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.sessions, sessionID)
		h.mu.Unlock()
	}()

	// Signal session is ready to client
	onEvent(EventEnvelope{
		Version:   1,
		Type:      TypeSessionReady,
		SessionID: sessionID,
		Sequence:  1,
		Timestamp: time.Now(),
	})

	// Framing read loop
	for {
		mt, data, err := c.Read(r.Context())
		if err != nil {
			slog.Info("WebSocket client disconnected or connection closed", slog.Any("error", err))
			break
		}

		if mt == websocket.MessageText {
			var env EventEnvelope
			if err := json.Unmarshal(data, &env); err != nil {
				slog.Warn("Failed to unmarshal websocket control event", slog.Any("error", err))
				continue
			}

			err = session.HandleEvent(r.Context(), env, nil)
			if err != nil {
				slog.Error("Session handle event error", slog.Any("error", err))
			}
		} else if mt == websocket.MessageBinary {
			// Binary frame: stream directly into PTT capture buffer
			// Mock sequence / envelope for streaming audio
			env := EventEnvelope{
				Version:   1,
				Type:      TypeAudioChunk,
				SessionID: sessionID,
				Sequence:  session.expectedSeq + 1,
				Timestamp: time.Now(),
			}
			err = session.HandleEvent(r.Context(), env, data)
			if err != nil {
				slog.Error("Session handle binary frame error", slog.Any("error", err))
			}
		}
	}
}
