-- name: ClaimOutboxEvents :many
SELECT id, aggregate_type, aggregate_id, event_type, payload, attempts, created_at
FROM outbox_events
WHERE processed_at IS NULL
    AND available_at <= now()
    AND attempts < max_attempts
ORDER BY created_at
FOR UPDATE SKIP LOCKED
LIMIT sqlc.arg('batch_size');

-- name: MarkOutboxEventProcessed :exec
UPDATE outbox_events
SET processed_at = now(), locked_at = NULL, locked_by = NULL
WHERE id = sqlc.arg('id');

-- name: MarkOutboxEventFailed :exec
UPDATE outbox_events
SET attempts = attempts + 1,
    last_error = sqlc.arg('error'),
    available_at = now() + (interval '1 second' * power(2, attempts)),
    locked_at = NULL,
    locked_by = NULL
WHERE id = sqlc.arg('id');

-- name: InsertOutboxEvent :one
INSERT INTO outbox_events (aggregate_type, aggregate_id, event_type, payload)
VALUES (sqlc.arg('aggregate_type'), sqlc.arg('aggregate_id'), sqlc.arg('event_type'), sqlc.arg('payload'))
RETURNING id;
