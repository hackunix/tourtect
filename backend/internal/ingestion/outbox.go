package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tourtect/backend/generated/database"
)

type OutboxProcessor struct {
	pool    *pgxpool.Pool
	queries *database.Queries
	stop    chan struct{}
}

func NewOutboxProcessor(pool *pgxpool.Pool) *OutboxProcessor {
	return &OutboxProcessor{
		pool:    pool,
		queries: database.New(pool),
		stop:    make(chan struct{}),
	}
}

func (op *OutboxProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	slog.Info("Background outbox worker processor started")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Outbox processor stopping due to context cancellation...")
			return
		case <-op.stop:
			slog.Info("Outbox processor stopping due to stop signal...")
			return
		case <-ticker.C:
			if err := op.processBatch(ctx); err != nil {
				slog.Error("Failed to process outbox event batch", slog.Any("error", err))
			}
		}
	}
}

func (op *OutboxProcessor) Stop() {
	close(op.stop)
}

func (op *OutboxProcessor) processBatch(ctx context.Context) error {
	// Claim batch in transaction using FOR UPDATE SKIP LOCKED
	tx, err := op.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start claim transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := op.queries.WithTx(tx)

	// Fetch up to 10 pending events
	events, err := qtx.ClaimOutboxEvents(ctx, 10)
	if err != nil {
		return fmt.Errorf("failed to claim outbox events: %w", err)
	}

	if len(events) == 0 {
		return nil
	}

	slog.Info("Claimed outbox events", slog.Int("count", len(events)))

	for _, event := range events {
		// Process each event inside transaction context or individually
		// V1 placeholders: search indexing, notification dispatch
		err := op.handleEvent(ctx, event)
		if err != nil {
			slog.Error("Failed to process outbox event", slog.Any("id", event.ID), slog.Any("error", err))
			// Increment attempts, set lock fields to null, update available_at
			errMsg := err.Error()
			failErr := qtx.MarkOutboxEventFailed(ctx, database.MarkOutboxEventFailedParams{
				ID:    event.ID,
				Error: &errMsg,
			})
			if failErr != nil {
				slog.Error("Failed to mark outbox event failed in database", slog.Any("id", event.ID), slog.Any("error", failErr))
			}
		} else {
			// Mark as processed successfully
			okErr := qtx.MarkOutboxEventProcessed(ctx, event.ID)
			if okErr != nil {
				slog.Error("Failed to mark outbox event processed in database", slog.Any("id", event.ID), slog.Any("error", okErr))
			}
			slog.Info("Successfully processed outbox event", slog.Any("id", event.ID), slog.String("type", event.EventType))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit outbox processing transaction: %w", err)
	}

	return nil
}

func (op *OutboxProcessor) handleEvent(ctx context.Context, event database.ClaimOutboxEventsRow) error {
	slog.Info("Processing job",
		slog.String("id", event.ID.String()),
		slog.String("aggregate_type", event.AggregateType),
		slog.String("event_type", event.EventType),
	)

	// Unmarshal payload
	var payload map[string]interface{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("invalid json payload: %w", err)
	}

	// Simulating handlers based on EventType
	switch event.EventType {
	case "post.created":
		slog.Info("Placeholder: Trigger search indexing for newly created post", slog.String("aggregate_id", event.AggregateID))
		return nil
	case "price_observation.created":
		slog.Info("Placeholder: Process pricing snapshot consolidation", slog.String("aggregate_id", event.AggregateID))
		return nil
	case "safety_incident.triggered":
		slog.Info("Placeholder: Dispatched safety notification broadcast", slog.String("aggregate_id", event.AggregateID))
		return nil
	default:
		return fmt.Errorf("unhandled outbox event type: %s", event.EventType)
	}
}
