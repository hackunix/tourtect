-- name: GetPriceSnapshot :one
SELECT
    snapshot_id, vertical, region_id, pricing_zone_id,
    service_segment, venue_type, unit, currency,
    p10_minor, p50_minor, p90_minor, exponent,
    sample_size, independent_source_count, version,
    effective_from, effective_to, created_at
FROM price_snapshots
WHERE vertical = sqlc.arg('vertical')
    AND region_id = sqlc.arg('region_id')
    AND unit = sqlc.arg('unit')
    AND currency = sqlc.arg('currency')
    AND service_segment = sqlc.arg('service_segment')
    AND venue_type = sqlc.arg('venue_type')
    AND effective_from <= sqlc.arg('observed_at')
    AND (effective_to IS NULL OR effective_to > sqlc.arg('observed_at'))
ORDER BY effective_from DESC
LIMIT 1;

-- name: GetPriceSnapshotFallbackRegion :one
SELECT
    snapshot_id, vertical, region_id, pricing_zone_id,
    service_segment, venue_type, unit, currency,
    p10_minor, p50_minor, p90_minor, exponent,
    sample_size, independent_source_count, version,
    effective_from, effective_to, created_at
FROM price_snapshots
WHERE vertical = sqlc.arg('vertical')
    AND unit = sqlc.arg('unit')
    AND currency = sqlc.arg('currency')
    AND effective_from <= sqlc.arg('observed_at')
    AND (effective_to IS NULL OR effective_to > sqlc.arg('observed_at'))
ORDER BY sample_size DESC
LIMIT 1;

-- name: CountRecentObservations :one
SELECT COUNT(*) AS count
FROM price_observations
WHERE vertical = sqlc.arg('vertical')
    AND region_id = sqlc.arg('region_id')
    AND unit = sqlc.arg('unit')
    AND currency = sqlc.arg('currency')
    AND observed_at >= sqlc.arg('since');
