-- name: CreateAssistantFeedback :one
INSERT INTO assistant_feedback (
    principal_id, session_id, assistant_message_id, feedback_type, field_name,
    original_ai_output, user_correction, final_confirmed_value,
    provider, model_version, policy_version, tool_results, consent_to_contribute
) VALUES (
    sqlc.arg('principal_id'), sqlc.arg('session_id'), sqlc.arg('assistant_message_id'),
    sqlc.arg('feedback_type'), sqlc.narg('field_name'), sqlc.arg('original_ai_output'),
    sqlc.arg('user_correction'), sqlc.arg('final_confirmed_value'), sqlc.narg('provider'),
    sqlc.narg('model_version'), sqlc.arg('policy_version'), sqlc.arg('tool_results'),
    sqlc.arg('consent_to_contribute')
)
RETURNING feedback_id, moderation_status, created_at;

-- name: CreateAssistantConfirmationAudit :exec
INSERT INTO assistant_confirmation_audit (
    confirmation_id, principal_id, session_id, action, decision, consent_scope, result_id
) VALUES (
    sqlc.arg('confirmation_id'), sqlc.arg('principal_id'), sqlc.arg('session_id'),
    sqlc.arg('action'), sqlc.arg('decision'), sqlc.narg('consent_scope'), sqlc.narg('result_id')
);

-- name: CreateAssistantModelTrace :exec
INSERT INTO assistant_model_traces (
    trace_id, principal_id, session_id, intent, tool_names, tool_durations_ms,
    provider, model_version, policy_version, retrieval_count, evidence_ids,
    outcome, error_category, fallback_used
) VALUES (
    sqlc.arg('trace_id'), sqlc.arg('principal_id'), sqlc.arg('session_id'), sqlc.arg('intent'),
    sqlc.arg('tool_names'), sqlc.arg('tool_durations_ms'), sqlc.narg('provider'),
    sqlc.narg('model_version'), sqlc.arg('policy_version'), sqlc.arg('retrieval_count'),
    sqlc.arg('evidence_ids'), sqlc.arg('outcome'), sqlc.narg('error_category'),
    sqlc.arg('fallback_used')
);
