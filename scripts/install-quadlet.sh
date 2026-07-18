#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="$ROOT_DIR/.env"
BUILD_IMAGES=true
START_SERVICES=true
WITH_FPT_SECRET=false

usage() {
  cat <<'EOF'
Usage: ./scripts/install-quadlet.sh [options]

Install the rootless Tourtect Quadlet units for the current user.

Options:
  --env-file PATH       Read deployment values from PATH (default: .env).
  --with-fpt-secret     Attach the existing tourtect_fpt_ai_api_key secret.
  --no-build            Reuse localhost/tourtect-{backend,web}:local images.
  --no-start            Install units without enabling/starting the Web unit.
  -h, --help            Show this help.
EOF
}

die() { printf '[Tourtect] ERROR: %s\n' "$*" >&2; exit 1; }

while (($#)); do
  case "$1" in
    --env-file) [[ $# -ge 2 ]] || die "--env-file needs a path"; ENV_FILE="$2"; shift ;;
    --with-fpt-secret) WITH_FPT_SECRET=true ;;
    --no-build) BUILD_IMAGES=false ;;
    --no-start) START_SERVICES=false ;;
    -h|--help) usage; exit 0 ;;
    *) die "Unknown option: $1" ;;
  esac
  shift
done

command -v podman >/dev/null 2>&1 || die "podman is required"
command -v systemctl >/dev/null 2>&1 || die "systemd user services are required for Quadlet"
[[ -r "$ENV_FILE" ]] || die "Environment file is not readable: $ENV_FILE"

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

load_env_file "$ENV_FILE"

POSTGRES_DB="${POSTGRES_DB:-tourtect}"
POSTGRES_USER="${POSTGRES_USER:-tourtect}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-change_me_postgres}"
MINIO_ROOT_USER="${MINIO_ROOT_USER:-tourtect}"
MINIO_ROOT_PASSWORD="${MINIO_ROOT_PASSWORD:-change_me_minio}"
FPT_AI_BASE_URL="${FPT_AI_BASE_URL:-https://mkp-api.fptcloud.com}"
TZ="${TZ:-Asia/Ho_Chi_Minh}"

for value in "$POSTGRES_DB" "$POSTGRES_USER" "$POSTGRES_PASSWORD" "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD" "$FPT_AI_BASE_URL" "$TZ"; do
  [[ "$value" != *$'\n'* ]] || die "Multiline environment values are not supported"
done

if $WITH_FPT_SECRET; then
  podman secret inspect tourtect_fpt_ai_api_key >/dev/null 2>&1 || \
    die "Create the secret first: ./scripts/podman-secret.sh create [path]"
fi

if $BUILD_IMAGES; then
  printf '[Tourtect] Building backend image\n'
  podman build --tag localhost/tourtect-backend:local --file "$ROOT_DIR/backend/Containerfile" "$ROOT_DIR/backend"
  printf '[Tourtect] Building Web image\n'
  podman build --tag localhost/tourtect-web:local \
    --build-arg NEXT_PUBLIC_API_URL=http://127.0.0.1:8080 \
    --file "$ROOT_DIR/web/Containerfile" "$ROOT_DIR/web"
fi

readonly QUADLET_DIR="$HOME/.config/containers/systemd"
readonly CONFIG_DIR="$HOME/.config/tourtect"
install -d -m 0700 "$QUADLET_DIR" "$CONFIG_DIR"
install -m 0644 "$ROOT_DIR"/deploy/podman/quadlet/*.{container,network,volume} "$QUADLET_DIR/"

umask 077
{
  printf 'POSTGRES_DB=%s\n' "$POSTGRES_DB"
  printf 'POSTGRES_USER=%s\n' "$POSTGRES_USER"
  printf 'POSTGRES_PASSWORD=%s\n' "$POSTGRES_PASSWORD"
  printf 'MINIO_ROOT_USER=%s\n' "$MINIO_ROOT_USER"
  printf 'MINIO_ROOT_PASSWORD=%s\n' "$MINIO_ROOT_PASSWORD"
  printf 'REDIS_PASSWORD=\n'
  printf 'FPT_AI_BASE_URL=%s\n' "$FPT_AI_BASE_URL"
  printf 'TZ=%s\n' "$TZ"
} > "$CONFIG_DIR/tourtect.env"

readonly SECRET_DROPIN="$QUADLET_DIR/tourtect-api.container.d"
if $WITH_FPT_SECRET; then
  install -d -m 0700 "$SECRET_DROPIN"
  {
    printf '[Container]\n'
    printf 'Secret=tourtect_fpt_ai_api_key,type=mount,target=fpt_ai_api_key\n'
    printf 'Environment=FPT_AI_API_KEY_FILE=/run/secrets/fpt_ai_api_key\n'
  } > "$SECRET_DROPIN/10-fpt-secret.conf"
elif [[ -f "$SECRET_DROPIN/10-fpt-secret.conf" ]]; then
  printf '[Tourtect] Existing FPT secret drop-in was left unchanged.\n'
fi

systemctl --user daemon-reload
if $START_SERVICES; then
  systemctl --user enable --now tourtect-web.service
  printf '[Tourtect] Quadlet stack is enabled. Web: http://127.0.0.1:3000\n'
else
  printf '[Tourtect] Quadlet units installed. Start with: systemctl --user start tourtect-web.service\n'
fi
