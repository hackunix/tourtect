-- +goose Up

CREATE TABLE posts (
    post_id               uuid PRIMARY KEY DEFAULT gen_ulid(),
    author_id             uuid NOT NULL REFERENCES principals(principal_id) ON DELETE CASCADE,
    post_type             text NOT NULL CHECK (post_type IN ('discussion','question','review','price_report','scam_report','tip','official_alert','external_link')),
    original_locale       text NOT NULL DEFAULT 'vi-VN' CHECK (original_locale IN ('vi-VN','ko-KR','zh-Hans','en','ru-RU')),
    title                 text NOT NULL CHECK (char_length(title) BETWEEN 1 AND 500),
    body                  text NOT NULL CHECK (char_length(body) >= 1),
    evidence_level        text NOT NULL DEFAULT 'none' CHECK (evidence_level IN ('none','metadata','verified_receipt','verified_source')),
    commercial_disclosure text NOT NULL DEFAULT 'none' CHECK (commercial_disclosure IN ('none','invited','gifted','affiliate','employee','sponsored')),
    moderation_status     text NOT NULL DEFAULT 'draft' CHECK (moderation_status IN ('draft','pending','published','limited','removed','appealed')),
    created_at            timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_posts_author ON posts (author_id);
CREATE INDEX idx_posts_status ON posts (moderation_status);
CREATE INDEX idx_posts_type ON posts (post_type);
CREATE INDEX idx_posts_created ON posts (created_at DESC);
CREATE INDEX idx_posts_title_trgm ON posts USING GIN (title gin_trgm_ops);

CREATE TABLE post_place_links (
    post_id   uuid NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
    place_id  uuid NOT NULL REFERENCES places(place_id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, place_id)
);

CREATE INDEX idx_post_place_links_place ON post_place_links (place_id);

-- +goose Down
DROP TABLE IF EXISTS post_place_links;
DROP TABLE IF EXISTS posts;
