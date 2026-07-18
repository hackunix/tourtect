-- name: GetLatestSafetyDirectoryVersion :one
SELECT version_id, version, description, published_at, created_at
FROM safety_directory_versions
ORDER BY published_at DESC
LIMIT 1;

-- name: GetSafetyEntriesByRegion :many
SELECT
    entry_id, version_id, region_id, service_name,
    service_type, phone_number, description, operating_hours,
    locale, is_approved, created_at
FROM safety_directory_entries
WHERE version_id = sqlc.arg('version_id')
    AND region_id = sqlc.arg('region_id')
    AND is_approved = true
ORDER BY service_type, service_name;

-- name: GetEmergencyServices :many
SELECT
    entry_id, version_id, region_id, service_name,
    service_type, phone_number, description
FROM safety_directory_entries
WHERE version_id = sqlc.arg('version_id')
    AND region_id = sqlc.arg('region_id')
    AND is_approved = true
    AND service_type IN ('police', 'ambulance', 'tourist_police')
ORDER BY service_type;
