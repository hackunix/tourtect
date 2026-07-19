-- +goose Up

CREATE TABLE assistant_feedback (
    feedback_id          uuid PRIMARY KEY DEFAULT gen_ulid(),
    principal_id         uuid NOT NULL REFERENCES principals(principal_id) ON DELETE CASCADE,
    session_id           uuid NOT NULL,
    assistant_message_id uuid NOT NULL,
    feedback_type        text NOT NULL CHECK (feedback_type IN ('helpful','not_helpful','correction','translation_incorrect','false_positive','confirm_extraction','contribute_redacted_observation')),
    field_name           text,
    original_ai_output   jsonb NOT NULL DEFAULT '{}'::jsonb,
    user_correction      jsonb NOT NULL DEFAULT '{}'::jsonb,
    final_confirmed_value jsonb NOT NULL DEFAULT '{}'::jsonb,
    provider             text,
    model_version        text,
    policy_version       text NOT NULL,
    tool_results         jsonb NOT NULL DEFAULT '[]'::jsonb,
    consent_to_contribute boolean NOT NULL DEFAULT false,
    moderation_status    text NOT NULL DEFAULT 'quarantined' CHECK (moderation_status IN ('quarantined','approved','rejected')),
    created_at           timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_assistant_feedback_moderation ON assistant_feedback (moderation_status, created_at);
CREATE INDEX idx_assistant_feedback_session ON assistant_feedback (session_id, created_at);

CREATE TABLE assistant_confirmation_audit (
    confirmation_id uuid PRIMARY KEY,
    principal_id    uuid NOT NULL REFERENCES principals(principal_id) ON DELETE CASCADE,
    session_id      uuid NOT NULL,
    action          text NOT NULL,
    decision        text NOT NULL CHECK (decision IN ('confirmed','rejected')),
    consent_scope   text,
    result_id       uuid,
    executed_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_assistant_confirmation_session ON assistant_confirmation_audit (session_id, executed_at);

CREATE TABLE assistant_model_traces (
    trace_id          uuid PRIMARY KEY,
    principal_id      uuid REFERENCES principals(principal_id) ON DELETE SET NULL,
    session_id        uuid NOT NULL,
    intent            text NOT NULL,
    tool_names        text[] NOT NULL DEFAULT '{}',
    tool_durations_ms bigint[] NOT NULL DEFAULT '{}',
    provider          text,
    model_version     text,
    policy_version    text NOT NULL,
    retrieval_count  integer NOT NULL DEFAULT 0,
    evidence_ids     uuid[] NOT NULL DEFAULT '{}',
    outcome           text NOT NULL,
    error_category    text,
    fallback_used     boolean NOT NULL DEFAULT false,
    created_at        timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_assistant_traces_session ON assistant_model_traces (session_id, created_at);
CREATE INDEX idx_assistant_traces_created ON assistant_model_traces (created_at);

-- +goose Down
DROP TABLE IF EXISTS assistant_model_traces;
DROP TABLE IF EXISTS assistant_confirmation_audit;
DROP TABLE IF EXISTS assistant_feedback;
