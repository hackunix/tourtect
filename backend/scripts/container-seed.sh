#!/bin/sh
set -eu

: "${POSTGRES_PASSWORD:?POSTGRES_PASSWORD is required}"
export PGPASSWORD="$POSTGRES_PASSWORD"
exec psql --set=ON_ERROR_STOP=1 \
  --host="${POSTGRES_HOST:-postgres}" \
  --port="${POSTGRES_PORT:-5432}" \
  --username="${POSTGRES_USER:-tourtect}" \
  --dbname="${POSTGRES_DB:-tourtect}" \
  --file=/app/seed/seed.sql
