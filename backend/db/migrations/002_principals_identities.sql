-- +goose Up

CREATE TABLE principals (
    principal_id  uuid PRIMARY KEY DEFAULT gen_ulid(),
    display_name  text NOT NULL CHECK (char_length(display_name) BETWEEN 1 AND 200),
    primary_email text UNIQUE,
    email_verified boolean NOT NULL DEFAULT false,
    status        text NOT NULL DEFAULT 'active' CHECK (status IN ('pending_email_verification','active','suspended','scheduled_for_deletion')),
    locale        text NOT NULL DEFAULT 'vi-VN' CHECK (locale IN ('vi-VN','ko-KR','zh-Hans','en','ru-RU')),
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_principals_email ON principals (primary_email) WHERE primary_email IS NOT NULL;
CREATE INDEX idx_principals_status ON principals (status);

CREATE TABLE identities (
    identity_id   uuid PRIMARY KEY DEFAULT gen_ulid(),
    principal_id  uuid NOT NULL REFERENCES principals(principal_id) ON DELETE CASCADE,
    provider      text NOT NULL CHECK (provider IN ('password','google')),
    issuer        text,
    subject       text,
    email_at_link text,
    password_hash text,
    linked_at     timestamptz NOT NULL DEFAULT now(),
    UNIQUE (provider, issuer, subject)
);

CREATE INDEX idx_identities_principal ON identities (principal_id);

CREATE TABLE sessions (
    session_id    uuid PRIMARY KEY DEFAULT gen_ulid(),
    principal_id  uuid NOT NULL REFERENCES principals(principal_id) ON DELETE CASCADE,
    device_label  text,
    refresh_hash  text,
    created_at    timestamptz NOT NULL DEFAULT now(),
    last_seen_at  timestamptz NOT NULL DEFAULT now(),
    expires_at    timestamptz NOT NULL,
    revoked_at    timestamptz
);

CREATE INDEX idx_sessions_principal ON sessions (principal_id);
CREATE INDEX idx_sessions_expires ON sessions (expires_at) WHERE revoked_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS identities;
DROP TABLE IF EXISTS principals;
