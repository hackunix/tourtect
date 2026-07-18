-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;

-- Helper function to generate UUIDv7 (time-sortable)
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION gen_ulid() RETURNS uuid AS $$
DECLARE
    ts_millis bigint;
    ts_bytes bytea;
    rand_bytes bytea;
    raw bytea;
BEGIN
    ts_millis := (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::bigint;
    ts_bytes := decode(lpad(to_hex(ts_millis), 12, '0'), 'hex');
    rand_bytes := gen_random_bytes(10);
    raw := ts_bytes || rand_bytes;
    -- Set version 7
    raw := set_byte(raw, 6, (get_byte(raw, 6) & 15) | 112);
    -- Set variant 2
    raw := set_byte(raw, 8, (get_byte(raw, 8) & 63) | 128);
    RETURN encode(raw, 'hex')::uuid;
END;
$$ LANGUAGE plpgsql VOLATILE;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION IF EXISTS gen_ulid();
DROP EXTENSION IF EXISTS unaccent;
DROP EXTENSION IF EXISTS pg_trgm;
DROP EXTENSION IF EXISTS postgis;
DROP EXTENSION IF EXISTS pgcrypto;
DROP EXTENSION IF EXISTS "uuid-ossp";
