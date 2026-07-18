-- name: ListPublishedPosts :many
SELECT
    p.post_id,
    p.author_id,
    p.post_type,
    p.original_locale,
    p.title,
    p.body,
    p.evidence_level,
    p.commercial_disclosure,
    p.moderation_status,
    p.created_at,
    p.updated_at,
    COALESCE(
        ARRAY_AGG(ppl.place_id) FILTER (WHERE ppl.place_id IS NOT NULL),
        '{}'
    ) AS place_ids
FROM posts p
LEFT JOIN post_place_links ppl ON ppl.post_id = p.post_id
WHERE p.moderation_status = 'published'
    AND (sqlc.narg('place_id')::uuid IS NULL OR EXISTS (
        SELECT 1 FROM post_place_links ppl2
        WHERE ppl2.post_id = p.post_id AND ppl2.place_id = sqlc.narg('place_id')
    ))
    AND (sqlc.narg('post_type')::text IS NULL OR p.post_type = sqlc.narg('post_type'))
    AND (sqlc.narg('cursor_id')::uuid IS NULL OR p.post_id < sqlc.narg('cursor_id'))
GROUP BY p.post_id
ORDER BY p.created_at DESC, p.post_id DESC
LIMIT sqlc.arg('page_limit');

-- name: GetPost :one
SELECT
    p.post_id,
    p.author_id,
    p.post_type,
    p.original_locale,
    p.title,
    p.body,
    p.evidence_level,
    p.commercial_disclosure,
    p.moderation_status,
    p.created_at,
    p.updated_at,
    COALESCE(
        ARRAY_AGG(ppl.place_id) FILTER (WHERE ppl.place_id IS NOT NULL),
        '{}'
    ) AS place_ids
FROM posts p
LEFT JOIN post_place_links ppl ON ppl.post_id = p.post_id
WHERE p.post_id = sqlc.arg('post_id')
GROUP BY p.post_id;

-- name: CreateDraftPost :one
INSERT INTO posts (author_id, post_type, original_locale, title, body, moderation_status)
VALUES (sqlc.arg('author_id'), sqlc.arg('post_type'), sqlc.arg('original_locale'), sqlc.arg('title'), sqlc.arg('body'), 'draft')
RETURNING *;

-- name: PublishPost :one
UPDATE posts
SET moderation_status = 'published', updated_at = now()
WHERE post_id = sqlc.arg('post_id') AND moderation_status = 'draft'
RETURNING *;

-- name: LinkPostToPlace :exec
INSERT INTO post_place_links (post_id, place_id) VALUES (sqlc.arg('post_id'), sqlc.arg('place_id'))
ON CONFLICT DO NOTHING;
