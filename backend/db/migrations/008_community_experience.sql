-- +goose Up

ALTER TABLE posts
    ADD COLUMN region_id text,
    ADD COLUMN structured_data jsonb NOT NULL DEFAULT '{}'::jsonb;

CREATE INDEX idx_posts_region_created ON posts (region_id, created_at DESC);
CREATE INDEX idx_posts_body_fts ON posts USING GIN (to_tsvector('simple', title || ' ' || body));

CREATE TABLE post_comments (
    comment_id        uuid PRIMARY KEY DEFAULT gen_ulid(),
    post_id           uuid NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
    author_id         uuid NOT NULL REFERENCES principals(principal_id) ON DELETE CASCADE,
    parent_comment_id uuid REFERENCES post_comments(comment_id) ON DELETE CASCADE,
    body              text NOT NULL CHECK (char_length(body) BETWEEN 1 AND 10000),
    moderation_status text NOT NULL DEFAULT 'published' CHECK (moderation_status IN ('pending','published','limited','removed')),
    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_comments_post_created ON post_comments (post_id, created_at);
CREATE INDEX idx_comments_parent ON post_comments (parent_comment_id);

CREATE TABLE post_votes (
    post_id      uuid NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
    principal_id uuid NOT NULL REFERENCES principals(principal_id) ON DELETE CASCADE,
    kind         text NOT NULL DEFAULT 'useful' CHECK (kind = 'useful'),
    created_at   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (post_id, principal_id, kind)
);

CREATE TABLE saved_posts (
    principal_id uuid NOT NULL REFERENCES principals(principal_id) ON DELETE CASCADE,
    post_id       uuid NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
    created_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (principal_id, post_id)
);
CREATE INDEX idx_saved_posts_principal_created ON saved_posts (principal_id, created_at DESC);

CREATE TABLE follows (
    follow_id          uuid PRIMARY KEY DEFAULT gen_ulid(),
    principal_id       uuid NOT NULL REFERENCES principals(principal_id) ON DELETE CASCADE,
    target_type        text NOT NULL CHECK (target_type IN ('principal','place')),
    target_principal_id uuid REFERENCES principals(principal_id) ON DELETE CASCADE,
    target_place_id    uuid REFERENCES places(place_id) ON DELETE CASCADE,
    created_at         timestamptz NOT NULL DEFAULT now(),
    CHECK ((target_type = 'principal' AND target_principal_id IS NOT NULL AND target_place_id IS NULL)
        OR (target_type = 'place' AND target_place_id IS NOT NULL AND target_principal_id IS NULL))
);
CREATE UNIQUE INDEX idx_follows_principal_target_principal ON follows (principal_id, target_principal_id) WHERE target_principal_id IS NOT NULL;
CREATE UNIQUE INDEX idx_follows_principal_target_place ON follows (principal_id, target_place_id) WHERE target_place_id IS NOT NULL;

CREATE TABLE notifications (
    notification_id uuid PRIMARY KEY DEFAULT gen_ulid(),
    principal_id     uuid NOT NULL REFERENCES principals(principal_id) ON DELETE CASCADE,
    kind             text NOT NULL,
    actor_id         uuid REFERENCES principals(principal_id) ON DELETE SET NULL,
    post_id          uuid REFERENCES posts(post_id) ON DELETE CASCADE,
    comment_id       uuid REFERENCES post_comments(comment_id) ON DELETE CASCADE,
    message          text NOT NULL,
    read_at          timestamptz,
    created_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX idx_notifications_principal_created ON notifications (principal_id, created_at DESC);

CREATE TABLE post_reports (
    report_id     uuid PRIMARY KEY DEFAULT gen_ulid(),
    post_id       uuid NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
    principal_id  uuid NOT NULL REFERENCES principals(principal_id) ON DELETE CASCADE,
    reason        text NOT NULL CHECK (reason IN ('safety','pii','harassment','spam','misinformation','other')),
    details       text,
    status        text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','reviewed','dismissed','actioned')),
    created_at    timestamptz NOT NULL DEFAULT now(),
    UNIQUE (post_id, principal_id, reason)
);

CREATE TABLE principal_blocks (
    blocker_id uuid NOT NULL REFERENCES principals(principal_id) ON DELETE CASCADE,
    blocked_id uuid NOT NULL REFERENCES principals(principal_id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (blocker_id, blocked_id),
    CHECK (blocker_id <> blocked_id)
);

CREATE TABLE reviews (
    post_id                    uuid PRIMARY KEY REFERENCES posts(post_id) ON DELETE CASCADE,
    place_id                   uuid NOT NULL REFERENCES places(place_id) ON DELETE CASCADE,
    visited_at                 date,
    overall_rating             smallint NOT NULL CHECK (overall_rating BETWEEN 1 AND 5),
    price_transparency_rating  smallint CHECK (price_transparency_rating BETWEEN 1 AND 5),
    service_rating             smallint CHECK (service_rating BETWEEN 1 AND 5),
    safety_rating              smallint CHECK (safety_rating BETWEEN 1 AND 5),
    value_rating               smallint CHECK (value_rating BETWEEN 1 AND 5)
);

CREATE TABLE community_price_reports (
    post_id          uuid PRIMARY KEY REFERENCES posts(post_id) ON DELETE CASCADE,
    item             text NOT NULL,
    amount_minor     bigint NOT NULL CHECK (amount_minor >= 0),
    currency         char(3) NOT NULL,
    unit             text NOT NULL,
    observed_at      timestamptz NOT NULL,
    ingestion_status text NOT NULL DEFAULT 'quarantined' CHECK (ingestion_status IN ('quarantined','accepted','rejected'))
);

CREATE TABLE community_scam_reports (
    post_id             uuid PRIMARY KEY REFERENCES posts(post_id) ON DELETE CASCADE,
    observed_at         timestamptz NOT NULL,
    current_safety_state text NOT NULL,
    triage_status       text NOT NULL DEFAULT 'pending' CHECK (triage_status IN ('pending','reviewed','escalated','closed'))
);

-- +goose Down
DROP TABLE IF EXISTS community_scam_reports;
DROP TABLE IF EXISTS community_price_reports;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS principal_blocks;
DROP TABLE IF EXISTS post_reports;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS follows;
DROP TABLE IF EXISTS saved_posts;
DROP TABLE IF EXISTS post_votes;
DROP TABLE IF EXISTS post_comments;
ALTER TABLE posts DROP COLUMN IF EXISTS structured_data, DROP COLUMN IF EXISTS region_id;
