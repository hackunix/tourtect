package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tourtect/backend/adapters/fptai"
	"github.com/tourtect/backend/internal/platform/config"
	"github.com/tourtect/backend/internal/platform/logging"
	"github.com/tourtect/backend/internal/realtime"
)

func main() {
	// 1. Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize logger
	logging.Init(cfg.LogLevel)
	slog.Info("Starting Tourtect Realtime WebSocket Server", slog.String("port", "8081"))

	// 3. Initialize AI providers
	var asrProvider fptai.ASRProvider
	var transProvider fptai.TranslationProvider

	if cfg.FptApiKey == "" {
		slog.Warn("FPT API Key missing. Falling back to Fake AI providers.")
		asrProvider = &fptai.FakeASR{}
		transProvider = &fptai.FakeTranslation{}
	} else {
		fptClient := fptai.NewClient(cfg.FptBaseURL, cfg.FptApiKey, 30*time.Second)
		asrProvider = fptai.NewRealASR(fptClient, "fpt-asr-v1")
		transProvider = fptai.NewRealTranslation(fptClient, "fpt-translation-v1")
	}

	// 4. Initialize WebSocket Upgrade handler
	wsHandler := realtime.NewHandler(asrProvider, transProvider)

	// 5. Wire routing
	mux := http.NewServeMux()
	mux.Handle("/v1/realtime", wsHandler)

	srv := &http.Server{
		Addr:         ":8081", // Realtime port
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go func() {
		slog.Info("Realtime server listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ListenAndServe failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Graceful shutdown coordination
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	slog.Info("Shutting down realtime server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful shutdown failed, forcing close", slog.Any("error", err))
	} else {
		slog.Info("Realtime server stopped successfully")
	}
}
