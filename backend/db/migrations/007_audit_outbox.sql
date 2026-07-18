-- +goose Up

CREATE TABLE audit_events (
    event_id      uuid PRIMARY KEY DEFAULT gen_ulid(),
    principal_id  uuid REFERENCES principals(principal_id) ON DELETE SET NULL,
    action        text NOT NULL,
    resource_type text NOT NULL,
    resource_id   text NOT NULL,
    details       jsonb,
    ip_hash       text,
    created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_principal ON audit_events (principal_id);
CREATE INDEX idx_audit_resource ON audit_events (resource_type, resource_id);
CREATE INDEX idx_audit_created ON audit_events (created_at DESC);

CREATE TABLE outbox_events (
    id             uuid PRIMARY KEY DEFAULT gen_ulid(),
    aggregate_type text NOT NULL,
    aggregate_id   text NOT NULL,
    event_type     text NOT NULL,
    payload        jsonb NOT NULL DEFAULT '{}',
    attempts       integer NOT NULL DEFAULT 0,
    max_attempts   integer NOT NULL DEFAULT 5,
    available_at   timestamptz NOT NULL DEFAULT now(),
    locked_at      timestamptz,
    locked_by      text,
    processed_at   timestamptz,
    last_error     text,
    created_at     timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_outbox_pending ON outbox_events (available_at)
    WHERE processed_at IS NULL;
CREATE INDEX idx_outbox_aggregate ON outbox_events (aggregate_type, aggregate_id);

-- +goose Down
DROP TABLE IF EXISTS outbox_events;
DROP TABLE IF EXISTS audit_events;
