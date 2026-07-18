-- +goose Up

CREATE TABLE safety_directory_versions (
    version_id   uuid PRIMARY KEY DEFAULT gen_ulid(),
    version      text NOT NULL UNIQUE,
    description  text,
    published_at timestamptz NOT NULL DEFAULT now(),
    created_at   timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE safety_directory_entries (
    entry_id       uuid PRIMARY KEY DEFAULT gen_ulid(),
    version_id     uuid NOT NULL REFERENCES safety_directory_versions(version_id) ON DELETE CASCADE,
    region_id      text NOT NULL,
    service_name   text NOT NULL,
    service_type   text NOT NULL CHECK (service_type IN ('police','ambulance','fire','embassy','hotline','crisis_center','tourist_police')),
    phone_number   text NOT NULL,
    description    text,
    operating_hours text,
    locale         text NOT NULL DEFAULT 'vi-VN',
    is_approved    boolean NOT NULL DEFAULT true,
    created_at     timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_safety_entries_version ON safety_directory_entries (version_id);
CREATE INDEX idx_safety_entries_region ON safety_directory_entries (region_id);
CREATE INDEX idx_safety_entries_type ON safety_directory_entries (service_type);

CREATE TABLE consent_records (
    consent_id   uuid PRIMARY KEY DEFAULT gen_ulid(),
    principal_id uuid REFERENCES principals(principal_id) ON DELETE SET NULL,
    session_id   text,
    scope        text NOT NULL CHECK (scope IN ('process_microphone','process_camera','precise_location','share_incident','contribute_redacted_data')),
    granted      boolean NOT NULL,
    granted_at   timestamptz NOT NULL DEFAULT now(),
    revoked_at   timestamptz,
    ip_hash      text,
    user_agent   text
);

CREATE INDEX idx_consent_principal ON consent_records (principal_id);
CREATE INDEX idx_consent_scope ON consent_records (scope);

-- +goose Down
DROP TABLE IF EXISTS consent_records;
DROP TABLE IF EXISTS safety_directory_entries;
DROP TABLE IF EXISTS safety_directory_versions;
