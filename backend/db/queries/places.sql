-- name: ListPlaces :many
SELECT
    p.place_id,
    p.name,
    p.category,
    p.region_id,
    p.address,
    ST_Y(p.coordinates::geometry) AS latitude,
    ST_X(p.coordinates::geometry) AS longitude,
    p.post_count,
    p.average_rating,
    p.freshness,
    p.created_at,
    COALESCE(
        ARRAY_AGG(pa.alias) FILTER (WHERE pa.alias IS NOT NULL),
        '{}'
    ) AS aliases
FROM places p
LEFT JOIN place_aliases pa ON pa.place_id = p.place_id
WHERE
    (sqlc.narg('region_id')::text IS NULL OR p.region_id = sqlc.narg('region_id'))
    AND (sqlc.narg('category')::text IS NULL OR p.category = sqlc.narg('category'))
    AND (sqlc.narg('search_query')::text IS NULL OR (
        p.name ILIKE '%' || sqlc.narg('search_query') || '%'
        OR EXISTS (
            SELECT 1 FROM place_aliases pa2
            WHERE pa2.place_id = p.place_id
            AND pa2.alias ILIKE '%' || sqlc.narg('search_query') || '%'
        )
    ))
    AND (sqlc.narg('cursor_id')::uuid IS NULL OR p.place_id < sqlc.narg('cursor_id'))
GROUP BY p.place_id
ORDER BY p.created_at DESC, p.place_id DESC
LIMIT sqlc.arg('page_limit');

-- name: ListPlacesNearby :many
SELECT
    p.place_id,
    p.name,
    p.category,
    p.region_id,
    p.address,
    ST_Y(p.coordinates::geometry) AS latitude,
    ST_X(p.coordinates::geometry) AS longitude,
    p.post_count,
    p.average_rating,
    p.freshness,
    p.created_at,
    ST_Distance(p.coordinates, ST_SetSRID(ST_MakePoint(sqlc.arg('lon'), sqlc.arg('lat')), 4326)::geography) AS distance_m,
    COALESCE(
        ARRAY_AGG(pa.alias) FILTER (WHERE pa.alias IS NOT NULL),
        '{}'
    ) AS aliases
FROM places p
LEFT JOIN place_aliases pa ON pa.place_id = p.place_id
WHERE ST_DWithin(
    p.coordinates,
    ST_SetSRID(ST_MakePoint(sqlc.arg('lon'), sqlc.arg('lat')), 4326)::geography,
    sqlc.arg('radius_m')
)
GROUP BY p.place_id
ORDER BY distance_m ASC
LIMIT sqlc.arg('page_limit');

-- name: GetPlace :one
SELECT
    p.place_id,
    p.name,
    p.category,
    p.region_id,
    p.address,
    p.description,
    p.phone,
    p.website,
    p.opening_hours,
    ST_Y(p.coordinates::geometry) AS latitude,
    ST_X(p.coordinates::geometry) AS longitude,
    p.post_count,
    p.review_count,
    p.average_rating,
    p.freshness,
    p.created_at,
    p.updated_at,
    COALESCE(
        ARRAY_AGG(pa.alias) FILTER (WHERE pa.alias IS NOT NULL),
        '{}'
    ) AS aliases
FROM places p
LEFT JOIN place_aliases pa ON pa.place_id = p.place_id
WHERE p.place_id = sqlc.arg('place_id')
GROUP BY p.place_id;
