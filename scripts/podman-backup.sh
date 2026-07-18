#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
OUTPUT_ROOT="$ROOT_DIR/backups"
COLD_VOLUMES=false

usage() {
  cat <<'EOF'
Usage: ./scripts/podman-backup.sh [--output DIR] [--cold-volumes]

By default, creates a logical PostgreSQL dump from the running stack.
With --cold-volumes, instead exports the three named volumes after verifying
that no running container is using them. Stop the stack first.
EOF
}

die() { printf '[Tourtect] ERROR: %s\n' "$*" >&2; exit 1; }

while (($#)); do
  case "$1" in
    --output) [[ $# -ge 2 ]] || die "--output needs a directory"; OUTPUT_ROOT="$2"; shift ;;
    --cold-volumes) COLD_VOLUMES=true ;;
    -h|--help) usage; exit 0 ;;
    *) die "Unknown option: $1" ;;
  esac
  shift
done

command -v podman >/dev/null 2>&1 || die "podman is required"
timestamp="$(date -u +%Y%m%dT%H%M%SZ)"
backup_dir="$OUTPUT_ROOT/tourtect-$timestamp"
install -d -m 0700 "$backup_dir"

if $COLD_VOLUMES; then
  for volume in tourtect-postgres-data tourtect-redis-data tourtect-minio-data; do
    podman volume exists "$volume" || die "Missing volume: $volume"
    if [[ -n "$(podman ps --quiet --filter "volume=$volume")" ]]; then
      die "Volume $volume is in use. Stop Compose and Quadlet services before --cold-volumes."
    fi
    printf '[Tourtect] Exporting %s\n' "$volume"
    podman volume export "$volume" --output "$backup_dir/$volume.tar"
  done
else
  postgres_container="$(podman ps --quiet --filter label=com.docker.compose.service=postgres | head -n 1)"
  if [[ -z "$postgres_container" ]] && podman container exists tourtect-postgres; then
    postgres_container=tourtect-postgres
  fi
  [[ -n "$postgres_container" ]] || die "PostgreSQL is not running; start Compose or Quadlet first"

  postgres_user="${POSTGRES_USER:-tourtect}"
  postgres_db="${POSTGRES_DB:-tourtect}"
  printf '[Tourtect] Creating PostgreSQL logical backup\n'
  podman exec "$postgres_container" pg_dump --format=custom --username="$postgres_user" --dbname="$postgres_db" > "$backup_dir/postgres.dump"
fi

(
  cd "$backup_dir"
  sha256sum ./* > SHA256SUMS
)
printf '[Tourtect] Backup completed: %s\n' "$backup_dir"
