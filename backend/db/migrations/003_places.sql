-- +goose Up

CREATE TABLE places (
    place_id      uuid PRIMARY KEY DEFAULT gen_ulid(),
    name          text NOT NULL CHECK (char_length(name) BETWEEN 1 AND 500),
    category      text NOT NULL CHECK (char_length(category) BETWEEN 1 AND 100),
    region_id     text NOT NULL,
    address       text,
    description   text,
    phone         text,
    website       text,
    opening_hours text,
    coordinates   geography(Point, 4326) NOT NULL,
    post_count    integer NOT NULL DEFAULT 0,
    review_count  integer NOT NULL DEFAULT 0,
    average_rating numeric(3,2) DEFAULT 0,
    freshness     timestamptz,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_places_region ON places (region_id);
CREATE INDEX idx_places_category ON places (category);
CREATE INDEX idx_places_coordinates ON places USING GIST (coordinates);
CREATE INDEX idx_places_name_trgm ON places USING GIN (name gin_trgm_ops);
CREATE INDEX idx_places_created ON places (created_at);

CREATE TABLE place_aliases (
    alias_id  uuid PRIMARY KEY DEFAULT gen_ulid(),
    place_id  uuid NOT NULL REFERENCES places(place_id) ON DELETE CASCADE,
    alias     text NOT NULL CHECK (char_length(alias) BETWEEN 1 AND 500),
    locale    text NOT NULL DEFAULT 'vi-VN',
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_place_aliases_place ON place_aliases (place_id);
CREATE INDEX idx_place_aliases_trgm ON place_aliases USING GIN (alias gin_trgm_ops);

-- +goose Down
DROP TABLE IF EXISTS place_aliases;
DROP TABLE IF EXISTS places;
