#!/bin/sh
set -eu

: "${POSTGRES_PASSWORD:?POSTGRES_PASSWORD is required}"

database_url="${DATABASE_URL:-postgres://${POSTGRES_USER:-tourtect}:${POSTGRES_PASSWORD}@${POSTGRES_HOST:-postgres}:${POSTGRES_PORT:-5432}/${POSTGRES_DB:-tourtect}?sslmode=disable}"
exec goose -dir /app/migrations postgres "$database_url" up
