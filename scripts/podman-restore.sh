#!/usr/bin/env bash
set -Eeuo pipefail

BACKUP_DIR=""
RESTORE_VOLUMES=false
CONFIRMED=false

usage() {
  cat <<'EOF'
Usage: ./scripts/podman-restore.sh --from DIR --confirm-overwrite [--volumes]

Restores postgres.dump into the running PostgreSQL container using pg_restore
--clean. With --volumes, replaces the three exact Tourtect named volumes from
the backup tarballs; all containers using them must already be stopped.
EOF
}

die() { printf '[Tourtect] ERROR: %s\n' "$*" >&2; exit 1; }

while (($#)); do
  case "$1" in
    --from) [[ $# -ge 2 ]] || die "--from needs a directory"; BACKUP_DIR="$2"; shift ;;
    --volumes) RESTORE_VOLUMES=true ;;
    --confirm-overwrite) CONFIRMED=true ;;
    -h|--help) usage; exit 0 ;;
    *) die "Unknown option: $1" ;;
  esac
  shift
done

command -v podman >/dev/null 2>&1 || die "podman is required"
[[ -n "$BACKUP_DIR" && -d "$BACKUP_DIR" ]] || die "Provide an existing backup directory with --from"
$CONFIRMED || die "Restore overwrites data; pass --confirm-overwrite after verifying the target"
[[ -f "$BACKUP_DIR/SHA256SUMS" ]] || die "SHA256SUMS is missing"
(
  cd "$BACKUP_DIR"
  sha256sum --check SHA256SUMS
)

if $RESTORE_VOLUMES; then
  for volume in tourtect-postgres-data tourtect-redis-data tourtect-minio-data; do
    [[ -f "$BACKUP_DIR/$volume.tar" ]] || die "Missing archive: $volume.tar"
    [[ -z "$(podman ps --quiet --filter "volume=$volume")" ]] || die "Volume $volume is still in use"
  done
  for volume in tourtect-postgres-data tourtect-redis-data tourtect-minio-data; do
    podman volume rm --force "$volume" >/dev/null 2>&1 || true
    podman volume create "$volume" >/dev/null
    podman volume import "$volume" "$BACKUP_DIR/$volume.tar"
    printf '[Tourtect] Restored volume: %s\n' "$volume"
  done
  printf '[Tourtect] Cold volumes restored. Start the stack and verify health before serving traffic.\n'
  exit 0
fi

[[ -f "$BACKUP_DIR/postgres.dump" ]] || die "postgres.dump is missing"
postgres_container="$(podman ps --quiet --filter name=^tourtect-postgres$)"
if [[ -z "$postgres_container" ]]; then
  postgres_container="$(podman ps --quiet --filter label=com.docker.compose.service=postgres | head -n 1)"
fi
[[ -n "$postgres_container" ]] || die "PostgreSQL is not running"

postgres_user="${POSTGRES_USER:-tourtect}"
postgres_db="${POSTGRES_DB:-tourtect}"
podman exec -i "$postgres_container" pg_restore --clean --if-exists --no-owner \
  --username="$postgres_user" --dbname="$postgres_db" < "$BACKUP_DIR/postgres.dump"
printf '[Tourtect] PostgreSQL restore completed. Run API smoke tests before serving traffic.\n'
