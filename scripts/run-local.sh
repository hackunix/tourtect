#!/usr/bin/env bash

set -Eeuo pipefail

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
SEED=true
MIGRATE=true
PRODUCTION=false
WITH_SEARCH=false
BACKEND_ONLY=false
API_PID=""

usage() {
  cat <<'EOF'
Usage: ./scripts/run-local.sh [options]

Start Tourtect dependencies, migrate/seed PostgreSQL, then run Go API and Web.

Options:
  --no-seed        Skip synthetic Hanoi seed data.
  --skip-migrate   Skip goose migrations.
  --production     Build and run the production Next.js server.
  --with-search    Also start the optional OpenSearch profile.
  --backend-only   Start API without Next.js.
  -h, --help       Show this help.

Ctrl+C stops API/Web processes. Containers and volumes are intentionally kept.
EOF
}

while (($#)); do
  case "$1" in
    --no-seed) SEED=false ;;
    --skip-migrate) MIGRATE=false ;;
    --production) PRODUCTION=true ;;
    --with-search) WITH_SEARCH=true ;;
    --backend-only) BACKEND_ONLY=true ;;
    -h|--help) usage; exit 0 ;;
    *) printf 'Unknown option: %s\n' "$1" >&2; usage >&2; exit 2 ;;
  esac
  shift
done

log() { printf '\n[Tourtect] %s\n' "$*"; }
die() { printf '[Tourtect] ERROR: %s\n' "$*" >&2; exit 1; }
has() { command -v "$1" >/dev/null 2>&1; }

load_env_file() {
  local file="$1" line key value
  while IFS= read -r line || [[ -n "$line" ]]; do
    line="${line%$'\r'}"
    [[ -z "$line" || "$line" == \#* || "$line" != *=* ]] && continue
    key="${line%%=*}"
    value="${line#*=}"
    [[ "$key" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]] || die "Invalid environment key in $file: $key"
    if [[ "$value" == \"*\" && "$value" == *\" ]]; then value="${value:1:${#value}-2}"; fi
    if [[ "$value" == \'*\' && "$value" == *\' ]]; then value="${value:1:${#value}-2}"; fi
    export "$key=$value"
  done < "$file"
}

if [[ -f "$HOME/.config/tourtect/env.sh" ]]; then
  # shellcheck disable=SC1090
  source "$HOME/.config/tourtect/env.sh"
fi
export PATH="$HOME/.local/go/bin:$HOME/.local/node/bin:$HOME/.local/bin:$HOME/go/bin:$PATH"
hash -r

[[ -f "$ROOT_DIR/.env" ]] || {
  cp "$ROOT_DIR/.env.example" "$ROOT_DIR/.env"
  log "Created .env from .env.example. Local demo credentials are active."
}

load_env_file "$ROOT_DIR/.env"

for command in podman go npm goose curl; do
  has "$command" || die "$command is missing. Run ./scripts/setup-linux.sh first."
done
podman compose version >/dev/null 2>&1 || die "podman compose provider is unavailable."

POSTGRES_DB="${POSTGRES_DB:-tourtect}"
POSTGRES_USER="${POSTGRES_USER:-tourtect}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-change_me_postgres}"
POSTGRES_HOST="${POSTGRES_HOST:-localhost}"
POSTGRES_PORT="${POSTGRES_PORT:-5432}"
PORT="${PORT:-8080}"
WEB_PORT="${WEB_PORT:-3000}"
DATABASE_URL="${DATABASE_URL:-postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable}"
API_URL="${API_URL:-http://127.0.0.1:${PORT}}"
NEXT_PUBLIC_API_URL="${NEXT_PUBLIC_API_URL:-http://127.0.0.1:${PORT}}"
export DATABASE_URL PORT API_URL NEXT_PUBLIC_API_URL

cleanup() {
  local status=$?
  trap - EXIT INT TERM
  if [[ -n "$API_PID" ]] && kill -0 "$API_PID" 2>/dev/null; then
    log "Stopping Tourtect API"
    kill "$API_PID" 2>/dev/null || true
    wait "$API_PID" 2>/dev/null || true
  fi
  printf '\n[Tourtect] Stateful containers are still running. Stop without deleting data: podman compose stop\n'
  exit "$status"
}
trap cleanup EXIT INT TERM

log "Starting PostgreSQL, Redis and MinIO"
(cd "$ROOT_DIR" && podman compose up -d postgres redis minio)
if $WITH_SEARCH; then
  (cd "$ROOT_DIR" && podman compose --profile search up -d opensearch)
fi

log "Waiting for PostgreSQL"
ready=false
for _ in $(seq 1 60); do
  if (cd "$ROOT_DIR" && podman compose exec -T postgres pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" >/dev/null 2>&1); then
    ready=true
    break
  fi
  sleep 1
done
$ready || die "PostgreSQL did not become ready within 60 seconds."

if $MIGRATE; then
  log "Applying database migrations"
  goose -dir "$ROOT_DIR/backend/db/migrations" postgres "$DATABASE_URL" up
fi

if $SEED; then
  log "Loading idempotent synthetic Hanoi seed"
  (cd "$ROOT_DIR" && podman compose exec -T postgres psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" < backend/db/seed/seed.sql)
fi

if curl --fail --silent "$API_URL/health/ready" >/dev/null 2>&1; then
  log "Reusing API already available at $API_URL"
else
  log "Building and starting Go API"
  mkdir -p "$ROOT_DIR/.local/bin" "$ROOT_DIR/.local/logs"
  (cd "$ROOT_DIR/backend" && go build -o "$ROOT_DIR/.local/bin/tourtect-api" ./cmd/api)
  "$ROOT_DIR/.local/bin/tourtect-api" > "$ROOT_DIR/.local/logs/api.log" 2>&1 &
  API_PID=$!
  for _ in $(seq 1 30); do
    if curl --fail --silent "$API_URL/health/ready" >/dev/null 2>&1; then break; fi
    if ! kill -0 "$API_PID" 2>/dev/null; then
      tail -n 80 "$ROOT_DIR/.local/logs/api.log" >&2 || true
      die "API exited before becoming ready."
    fi
    sleep 1
  done
  curl --fail --silent "$API_URL/health/ready" >/dev/null || {
    tail -n 80 "$ROOT_DIR/.local/logs/api.log" >&2 || true
    die "API did not become ready within 30 seconds."
  }
fi

printf '\n[Tourtect] API ready: %s\n' "$API_URL"
if [[ -n "$API_PID" ]]; then printf '[Tourtect] API log: %s\n' "$ROOT_DIR/.local/logs/api.log"; fi

if $BACKEND_ONLY; then
  log "Backend-only mode. Press Ctrl+C to stop the API process."
  while :; do sleep 3600; done
fi

[[ -d "$ROOT_DIR/web/node_modules" ]] || { log "Installing Web dependencies"; (cd "$ROOT_DIR/web" && npm ci); }

if $PRODUCTION; then
  log "Building Next.js production bundle"
  (cd "$ROOT_DIR/web" && npm run build)
  log "Web ready at http://127.0.0.1:$WEB_PORT"
  (cd "$ROOT_DIR/web" && npm run start -- -H 127.0.0.1 -p "$WEB_PORT")
else
  log "Starting Next.js development server at http://127.0.0.1:$WEB_PORT"
  (cd "$ROOT_DIR/web" && npm run dev -- -H 127.0.0.1 -p "$WEB_PORT")
fi
