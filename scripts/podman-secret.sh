#!/usr/bin/env bash
set -Eeuo pipefail

readonly SECRET_NAME="tourtect_fpt_ai_api_key"

usage() {
  cat <<'EOF'
Usage:
  ./scripts/podman-secret.sh create [path]
  ./scripts/podman-secret.sh remove
  ./scripts/podman-secret.sh status

Without path, `create` reads the secret from stdin. The value is never printed.
Use compose.secrets.yaml together with compose.yaml after creating the secret.
EOF
}

require_podman() {
  command -v podman >/dev/null 2>&1 || {
    echo "podman is required" >&2
    exit 1
  }
}

secret_exists() {
  podman secret inspect "$SECRET_NAME" >/dev/null 2>&1
}

create_secret() {
  local source="${1:--}"

  if secret_exists; then
    echo "Secret $SECRET_NAME already exists; remove it explicitly before replacing it." >&2
    exit 1
  fi

  if [[ "$source" != "-" && ! -r "$source" ]]; then
    echo "Secret file is not readable: $source" >&2
    exit 1
  fi

  if [[ "$source" == "-" ]]; then
    [[ -t 0 ]] && echo "Paste the secret, then press Ctrl-D:" >&2
    podman secret create "$SECRET_NAME" - >/dev/null
  else
    podman secret create "$SECRET_NAME" "$source" >/dev/null
  fi
  echo "Created Podman secret: $SECRET_NAME"
}

require_podman
case "${1:-}" in
  create)
    create_secret "${2:--}"
    ;;
  remove)
    if secret_exists; then
      podman secret rm "$SECRET_NAME" >/dev/null
      echo "Removed Podman secret: $SECRET_NAME"
    else
      echo "Secret is not present: $SECRET_NAME"
    fi
    ;;
  status)
    if secret_exists; then
      echo "Podman secret is present: $SECRET_NAME"
    else
      echo "Podman secret is not present: $SECRET_NAME"
      exit 1
    fi
    ;;
  -h|--help)
    usage
    ;;
  *)
    usage
    exit 2
    ;;
esac
