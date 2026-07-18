-- +goose Up

CREATE TABLE price_snapshots (
    snapshot_id   uuid PRIMARY KEY DEFAULT gen_ulid(),
    vertical      text NOT NULL CHECK (vertical IN ('taxi','exchange','food','tour','street_retail')),
    region_id     text NOT NULL,
    pricing_zone_id text,
    service_segment text NOT NULL DEFAULT 'standard' CHECK (service_segment IN ('budget','standard','premium','luxury','regulated')),
    venue_type    text NOT NULL DEFAULT 'fixed_shop' CHECK (venue_type IN ('fixed_shop','casual_eatery','street_stall','mobile_vendor','market_stall','attraction_concession','transport_vendor','peer_to_peer')),
    unit          text NOT NULL,
    currency      text NOT NULL CHECK (char_length(currency) = 3),
    p10_minor     bigint NOT NULL,
    p50_minor     bigint NOT NULL,
    p90_minor     bigint NOT NULL,
    exponent      integer NOT NULL DEFAULT 0,
    sample_size   integer NOT NULL CHECK (sample_size >= 0),
    independent_source_count integer NOT NULL DEFAULT 1,
    version       text NOT NULL,
    effective_from timestamptz NOT NULL DEFAULT now(),
    effective_to   timestamptz,
    created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_snapshots_lookup ON price_snapshots (vertical, region_id, unit, currency, service_segment, venue_type);
CREATE INDEX idx_snapshots_version ON price_snapshots (version);
CREATE INDEX idx_snapshots_effective ON price_snapshots (effective_from);

CREATE TABLE price_observations (
    observation_id uuid PRIMARY KEY DEFAULT gen_ulid(),
    snapshot_id    uuid REFERENCES price_snapshots(snapshot_id) ON DELETE SET NULL,
    vertical       text NOT NULL CHECK (vertical IN ('taxi','exchange','food','tour','street_retail')),
    region_id      text NOT NULL,
    raw_item       text NOT NULL,
    canonical_item_id text,
    amount_minor   bigint NOT NULL,
    currency       text NOT NULL CHECK (char_length(currency) = 3),
    exponent       integer NOT NULL DEFAULT 0,
    unit           text NOT NULL,
    service_segment text NOT NULL DEFAULT 'standard',
    venue_type     text NOT NULL DEFAULT 'fixed_shop',
    transaction_context text NOT NULL DEFAULT 'posted_price',
    extraction_confidence numeric(4,3) NOT NULL DEFAULT 1.0,
    user_confirmed boolean NOT NULL DEFAULT false,
    source         text NOT NULL DEFAULT 'user_report',
    observed_at    timestamptz NOT NULL DEFAULT now(),
    created_at     timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_observations_snapshot ON price_observations (snapshot_id);
CREATE INDEX idx_observations_lookup ON price_observations (vertical, region_id, unit, currency);
CREATE INDEX idx_observations_observed ON price_observations (observed_at);

-- +goose Down
DROP TABLE IF EXISTS price_observations;
DROP TABLE IF EXISTS price_snapshots;
