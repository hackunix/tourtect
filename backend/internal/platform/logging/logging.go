package logging

import (
	"log/slog"
	"os"
	"strings"
)

var Logger *slog.Logger

func Init(levelStr string) {
	var level slog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: redactAttr,
	}

	Logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(Logger)
}

func redactAttr(groups []string, a slog.Attr) slog.Attr {
	keyLower := strings.ToLower(a.Key)
	if strings.Contains(keyLower, "auth") ||
		strings.Contains(keyLower, "cookie") ||
		strings.Contains(keyLower, "key") ||
		strings.Contains(keyLower, "token") ||
		strings.Contains(keyLower, "password") ||
		strings.Contains(keyLower, "transcript") ||
		strings.Contains(keyLower, "audio") ||
		strings.Contains(keyLower, "incident") {
		return slog.String(a.Key, "[REDACTED]")
	}
	return a
}
