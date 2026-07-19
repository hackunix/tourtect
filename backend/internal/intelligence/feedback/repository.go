package feedback

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tourtect/backend/generated/database"
	"github.com/tourtect/backend/internal/intelligence/model"
)

type Repository struct{ queries *database.Queries }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{queries: database.New(pool)} }

type FeedbackInput struct {
	PrincipalID, SessionID, AssistantMessageID uuid.UUID
	FeedbackType                               string
	Field, OriginalValue, CorrectedValue       *string
	ConsentToContribute                        bool
	OriginalResponse                           model.Response
}

func (r *Repository) Create(ctx context.Context, input FeedbackInput) (uuid.UUID, string, time.Time, error) {
	original, _ := json.Marshal(input.OriginalResponse)
	correction, _ := json.Marshal(map[string]any{"field": input.Field, "original_value": input.OriginalValue, "corrected_value": input.CorrectedValue})
	finalValue := json.RawMessage(`{}`)
	if input.CorrectedValue != nil {
		finalValue, _ = json.Marshal(map[string]any{"field": input.Field, "value": *input.CorrectedValue, "user_confirmed": true})
	}
	toolResults, _ := json.Marshal(input.OriginalResponse.ToolResults)
	row, err := r.queries.CreateAssistantFeedback(ctx, database.CreateAssistantFeedbackParams{PrincipalID: input.PrincipalID, SessionID: input.SessionID, AssistantMessageID: input.AssistantMessageID, FeedbackType: input.FeedbackType, FieldName: input.Field, OriginalAiOutput: original, UserCorrection: correction, FinalConfirmedValue: finalValue, PolicyVersion: "assistant-policy-2026-07-v1", ToolResults: toolResults, ConsentToContribute: input.ConsentToContribute})
	if err != nil {
		return uuid.Nil, "", time.Time{}, fmt.Errorf("store assistant feedback: %w", err)
	}
	return row.FeedbackID, row.ModerationStatus, row.CreatedAt, nil
}

func (r *Repository) AuditConfirmation(ctx context.Context, principalID, sessionID, confirmationID uuid.UUID, action, decision string, resultID *uuid.UUID) error {
	var result pgtype.UUID
	if resultID != nil {
		result = pgtype.UUID{Bytes: *resultID, Valid: true}
	}
	if err := r.queries.CreateAssistantConfirmationAudit(ctx, database.CreateAssistantConfirmationAuditParams{ConfirmationID: confirmationID, PrincipalID: principalID, SessionID: sessionID, Action: action, Decision: decision, ResultID: result}); err != nil {
		return fmt.Errorf("audit assistant confirmation: %w", err)
	}
	return nil
}

func (r *Repository) CreateTrace(ctx context.Context, principalID uuid.UUID, trace model.Trace) error {
	traceID, err := uuid.Parse(trace.TraceID)
	if err != nil {
		return err
	}
	sessionID, err := uuid.Parse(trace.SessionID)
	if err != nil {
		return err
	}
	evidenceIDs := make([]uuid.UUID, 0, len(trace.EvidenceIDs))
	for _, raw := range trace.EvidenceIDs {
		if id, err := uuid.Parse(raw); err == nil {
			evidenceIDs = append(evidenceIDs, id)
		}
	}
	var errorCategory *string
	if trace.ErrorCategory != "" {
		errorCategory = &trace.ErrorCategory
	}
	return r.queries.CreateAssistantModelTrace(ctx, database.CreateAssistantModelTraceParams{TraceID: traceID, PrincipalID: pgtype.UUID{Bytes: principalID, Valid: true}, SessionID: sessionID, Intent: trace.Intent, ToolNames: trace.ToolNames, ToolDurationsMs: trace.ToolDurationsMS, PolicyVersion: trace.PolicyVersion, RetrievalCount: int32(trace.RetrievalCount), EvidenceIds: evidenceIDs, Outcome: trace.Outcome, ErrorCategory: errorCategory, FallbackUsed: trace.FallbackUsed})
}
