#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="$ROOT_DIR/.env"

log() {
  printf '[Tourtect Deploy] %s\n' "$*"
}

error() {
  printf '[Tourtect Deploy] ERROR: %s\n' "$*" >&2
  exit 1
}

# 1. Load env file
if [[ ! -f "$ENV_FILE" ]]; then
  error ".env file not found at $ENV_FILE. Please create it first."
fi

# Simple env loader
load_env_file() {
  local file="$1" line key value
  while IFS= read -r line || [[ -n "$line" ]]; do
    line="${line%$'\r'}"
    [[ -z "$line" || "$line" == \#* || "$line" != *=* ]] && continue
    key="${line%%=*}"
    value="${line#*=}"
    if [[ "$key" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]]; then
      if [[ "$value" == \"*\" && "$value" == *\" ]]; then value="${value:1:${#value}-2}"; fi
      if [[ "$value" == \'*\' && "$value" == *\' ]]; then value="${value:1:${#value}-2}"; fi
      export "$key=$value"
    fi
  done < "$file"
}

load_env_file "$ENV_FILE"

# 2. Check deployment variables
DEPLOY_HOST="${DEPLOY_HOST:-}"
DEPLOY_USER="${DEPLOY_USER:-}"
DEPLOY_PORT="${DEPLOY_PORT:-22}"
DEPLOY_KEY_PATH="${DEPLOY_KEY_PATH:-~/.ssh/id_rsa}"
DEPLOY_PATH="${DEPLOY_PATH:-/var/www/tourtect}"

if [[ -z "$DEPLOY_HOST" || -z "$DEPLOY_USER" ]]; then
  error "DEPLOY_HOST and DEPLOY_USER must be set in .env"
fi

# Expand tilde in key path if present
if [[ "$DEPLOY_KEY_PATH" == ~* ]]; then
  DEPLOY_KEY_PATH="${DEPLOY_KEY_PATH/#\~/$HOME}"
fi

log "Target host: $DEPLOY_USER@$DEPLOY_HOST:$DEPLOY_PORT"
log "Target path: $DEPLOY_PATH"
log "SSH Key: $DEPLOY_KEY_PATH"

# 3. Check key file presence
if [[ ! -f "$DEPLOY_KEY_PATH" ]]; then
  error "SSH private key not found at $DEPLOY_KEY_PATH"
fi

# 4. Test SSH connection
log "Testing SSH connection..."
if ! ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=accept-new -p "$DEPLOY_PORT" -i "$DEPLOY_KEY_PATH" "$DEPLOY_USER@$DEPLOY_HOST" echo "SSH Connection OK" &>/dev/null; then
  error "Could not establish SSH connection to $DEPLOY_USER@$DEPLOY_HOST. Please verify connection details and authorize the public key."
fi
log "SSH Connection test passed!"

# 5. Create remote path
log "Creating remote target path if not exists..."
ssh -p "$DEPLOY_PORT" -i "$DEPLOY_KEY_PATH" "$DEPLOY_USER@$DEPLOY_HOST" "mkdir -p \"$DEPLOY_PATH\""

# 6. Synchronize files
log "Synchronizing files..."
# Check if rsync is available on remote
if ssh -p "$DEPLOY_PORT" -i "$DEPLOY_KEY_PATH" "$DEPLOY_USER@$DEPLOY_HOST" "command -v rsync" >/dev/null 2>&1; then
  log "Using rsync for file transfer"
  rsync -avz -e "ssh -p $DEPLOY_PORT -i $DEPLOY_KEY_PATH" \
    --exclude="node_modules" \
    --exclude=".git" \
    --exclude=".local" \
    --exclude="backups" \
    --exclude="*.log" \
    --exclude="android/.gradle" \
    --exclude="android/build" \
    --exclude="android/*/build" \
    --exclude="web/.next" \
    --exclude="web/node_modules" \
    "$ROOT_DIR/" "$DEPLOY_USER@$DEPLOY_HOST:$DEPLOY_PATH/"
else
  log "rsync not available on remote, falling back to scp (tar over ssh)"
  # Create a tar archive of the project (excluding unwanted dirs) with stable timestamps and extract remotely
  tar --exclude="node_modules" \
      --exclude=".git" \
      --exclude=".local" \
      --exclude="backups" \
      --exclude="*.log" \
      --exclude="android/.gradle" \
      --exclude="android/build" \
      --exclude="android/*/build" \
      --exclude="web/.next" \
      --exclude="web/node_modules" \
      --mtime='2022-01-01' -czf - -C "$ROOT_DIR" . | \
    ssh -p "$DEPLOY_PORT" -i "$DEPLOY_KEY_PATH" "$DEPLOY_USER@$DEPLOY_HOST" "export PATH=\$HOME/.local/bin:\$PATH; tar -xzf - -C \"$DEPLOY_PATH\""
fi

log "Sync completed successfully!"

# 7. Start the application on the host
log "Starting the application on the host..."
ssh -p "$DEPLOY_PORT" -i "$DEPLOY_KEY_PATH" "$DEPLOY_USER@$DEPLOY_HOST" "export PATH=\$HOME/.local/bin:/usr/bin:/usr/local/bin:\$PATH; cd \"$DEPLOY_PATH\" && command -v podman >/dev/null 2>&1 || { echo 'Podman not installed on remote host'; exit 1; } && (systemctl --user start podman.socket 2>/dev/null || nohup podman system service -t 0 >/dev/null 2>&1 &) && podman compose --profile app down && podman compose --profile app up --build -d"

log "Deployment finished successfully!"
