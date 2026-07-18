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

	"github.com/tourtect/backend/internal/ingestion"
	"github.com/tourtect/backend/internal/platform/config"
	"github.com/tourtect/backend/internal/platform/database"
	"github.com/tourtect/backend/internal/platform/logging"
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
	slog.Info("Starting Tourtect Background Worker Server")

	// 3. Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Database connection failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	// 4. Initialize Outbox processor
	processor := ingestion.NewOutboxProcessor(db.Pool)

	// Run processor in background
	workerCtx, workerCancel := context.WithCancel(context.Background())
	go processor.Start(workerCtx)

	// Stub health server for worker liveness/readiness probes (port 8082)
	mux := http.NewServeMux()
	mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		pingCtx, pingCancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer pingCancel()
		if err := db.Ping(pingCtx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"degraded"}`))
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		}
	})

	srv := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	go func() {
		slog.Info("Worker health server listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Worker health server failed", slog.Any("error", err))
		}
	}()

	// Graceful shutdown coordination
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	slog.Info("Shutting down worker gracefully...")

	workerCancel()
	processor.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	_ = srv.Shutdown(shutdownCtx)
	slog.Info("Worker stopped successfully")
}
