package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id"
)

// ProblemDetail represents RFC 7807 problem detail response
type ProblemDetail struct {
	Type     string `json:"type"`
	Status   int    `json:"status"`
	Title    string `json:"title"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	TraceID  string `json:"trace_id,omitempty"`
}

func WriteError(w http.ResponseWriter, status int, title, detail, instance, traceID string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)

	prob := ProblemDetail{
		Type:     fmt.Sprintf("https://tourtect.example/errors/%d", status),
		Status:   status,
		Title:    title,
		Detail:   detail,
		Instance: instance,
		TraceID:  traceID,
	}

	_ = json.NewEncoder(w).Encode(prob)
}

// RequestID middleware assigns a UUIDv4 request ID to each request context
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID helper extracts Request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// GetUserID helper extracts User ID from context
func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}

// responseWriterDelegator captures HTTP status code and size for logging
type responseWriterDelegator struct {
	http.ResponseWriter
	status int
	size   int
}

func (d *responseWriterDelegator) WriteHeader(status int) {
	d.status = status
	d.ResponseWriter.WriteHeader(status)
}

func (d *responseWriterDelegator) Write(b []byte) (int, error) {
	if d.status == 0 {
		d.status = http.StatusOK
	}
	n, err := d.ResponseWriter.Write(b)
	d.size += n
	return n, err
}

// Logging middleware logs each request and response using structured slog
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := GetRequestID(r.Context())

		// Redact sensitive headers for logging
		headers := make(map[string][]string)
		for k, v := range r.Header {
			kLower := strings.ToLower(k)
			if strings.Contains(kLower, "auth") || strings.Contains(kLower, "cookie") || strings.Contains(kLower, "key") {
				headers[k] = []string{"[REDACTED]"}
			} else {
				headers[k] = v
			}
		}

		slog.Log(r.Context(), slog.LevelInfo, "HTTP Request Started",
			slog.String("request_id", reqID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Any("headers", headers),
		)

		d := &responseWriterDelegator{ResponseWriter: w}
		next.ServeHTTP(d, r)

		latency := time.Since(start)
		slog.Log(r.Context(), slog.LevelInfo, "HTTP Request Completed",
			slog.String("request_id", reqID),
			slog.Int("status", d.status),
			slog.Int("size", d.size),
			slog.Duration("latency", latency),
		)
	})
}

// PanicRecovery middleware recovers from panics and returns 500 ProblemDetail
func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				reqID := GetRequestID(r.Context())
				slog.Error("HTTP Server Panic",
					slog.Any("error", err),
					slog.String("stack", string(debug.Stack())),
					slog.String("request_id", reqID),
				)
				WriteError(w, http.StatusInternalServerError, "Internal Server Error", "An unexpected error occurred", r.URL.Path, reqID)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Timeout middleware applies context timeout to requests
func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// BodySizeLimit limits the request body size
func BodySizeLimit(limitBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, limitBytes)
			next.ServeHTTP(w, r)
		})
	}
}

// CORS adds standard CORS headers
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// SecurityHeaders adds standard security headers
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "no-referrer-when-downgrade")
		next.ServeHTTP(w, r)
	})
}

// AuthBoundary parses the request authentication credentials or uses a permissive mock author ID for vertical slices
func AuthBoundary(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permissive fallback: Use the default seed alpha user ID for the slice
		authorID := "019078a0-0001-7000-8000-000000000001"
		
		// If authorization header exists, we extract or mock
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "beta" {
				authorID = "019078a0-0001-7000-8000-000000000002"
			}
		}

		ctx := context.WithValue(r.Context(), UserIDKey, authorID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
