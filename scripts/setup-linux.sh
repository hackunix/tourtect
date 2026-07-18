#!/usr/bin/env bash

set -Eeuo pipefail

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
GO_VERSION="$(awk '/^go / { print $2; exit }' "$ROOT_DIR/backend/go.mod")"
NODE_VERSION="${TOURTECT_NODE_VERSION:-22.18.0}"
CHECK_ONLY=false
INSTALL_BROWSER=false
START_AFTER=false
SKIP_SYSTEM=false

usage() {
  cat <<'EOF'
Usage: ./scripts/setup-linux.sh [options]

Prepare a generic Linux development environment for Tourtect.

Options:
  --check-only           Report missing dependencies without installing anything.
  --install-browser      Install Playwright Chromium for E2E tests.
  --skip-system-packages Do not invoke apt/dnf/pacman/zypper.
  --start                Run scripts/run-local.sh after setup.
  -h, --help             Show this help.

Overrides:
  TOURTECT_NODE_VERSION  Node.js version installed to user space when needed.
EOF
}

while (($#)); do
  case "$1" in
    --check-only) CHECK_ONLY=true ;;
    --install-browser) INSTALL_BROWSER=true ;;
    --skip-system-packages) SKIP_SYSTEM=true ;;
    --start) START_AFTER=true ;;
    -h|--help) usage; exit 0 ;;
    *) printf 'Unknown option: %s\n' "$1" >&2; usage >&2; exit 2 ;;
  esac
  shift
done

log() { printf '\n[Tourtect] %s\n' "$*"; }
warn() { printf '[Tourtect] WARNING: %s\n' "$*" >&2; }
die() { printf '[Tourtect] ERROR: %s\n' "$*" >&2; exit 1; }
has() { command -v "$1" >/dev/null 2>&1; }

version_ge() {
  local current="$1" required="$2"
  [[ "$(printf '%s\n%s\n' "$required" "$current" | sort -V | head -n1)" == "$required" ]]
}

linux_arch() {
  case "$(uname -m)" in
    x86_64|amd64) printf 'amd64' ;;
    aarch64|arm64) printf 'arm64' ;;
    *) die "Unsupported CPU architecture: $(uname -m). Supported: x86_64, aarch64." ;;
  esac
}

node_arch() {
  case "$(linux_arch)" in
    amd64) printf 'x64' ;;
    arm64) printf 'arm64' ;;
  esac
}

sudo_cmd=()
if ((EUID != 0)); then
  has sudo || { $CHECK_ONLY && warn "sudo is missing" || die "sudo is required to install system packages."; }
  sudo_cmd=(sudo)
fi

install_system_packages() {
  $SKIP_SYSTEM && return
  local manager=""
  for candidate in apt-get dnf pacman zypper; do
    if has "$candidate"; then manager="$candidate"; break; fi
  done
  [[ -n "$manager" ]] || { warn "No supported package manager found; continuing with user-space tools."; return; }
  log "Installing base packages with $manager"
  case "$manager" in
    apt-get)
      "${sudo_cmd[@]}" apt-get update
      "${sudo_cmd[@]}" apt-get install -y ca-certificates curl git make xz-utils build-essential podman podman-compose
      ;;
    dnf)
      "${sudo_cmd[@]}" dnf install -y ca-certificates curl git make xz gcc gcc-c++ podman podman-compose
      ;;
    pacman)
      "${sudo_cmd[@]}" pacman -Sy --needed --noconfirm ca-certificates curl git make xz base-devel podman podman-compose
      ;;
    zypper)
      "${sudo_cmd[@]}" zypper --non-interactive install ca-certificates curl git make xz gcc gcc-c++ podman podman-compose
      ;;
  esac
}

write_shell_environment() {
  local config_dir="$HOME/.config/tourtect"
  mkdir -p "$config_dir" "$HOME/.local/bin" "$HOME/.local/opt"
  cat > "$config_dir/env.sh" <<'EOF'
# Tourtect-managed user-space tool paths.
export PATH="$HOME/.local/go/bin:$HOME/.local/node/bin:$HOME/.local/bin:$HOME/go/bin:$PATH"
EOF
  local marker='[ -f "$HOME/.config/tourtect/env.sh" ] && . "$HOME/.config/tourtect/env.sh"'
  touch "$HOME/.profile"
  if ! grep -Fq '.config/tourtect/env.sh' "$HOME/.profile"; then
    printf '\n# Tourtect development tools\n%s\n' "$marker" >> "$HOME/.profile"
  fi
  # shellcheck disable=SC1090
  source "$config_dir/env.sh"
}

download_verified() {
  local url="$1" output="$2" checksum_url="$3"
  local expected actual
  curl --fail --location --silent --show-error "$url" --output "$output"
  expected="$(curl --fail --location --silent --show-error "$checksum_url" | awk '{print $1}')"
  actual="$(sha256sum "$output" | awk '{print $1}')"
  [[ -n "$expected" && "$actual" == "$expected" ]] || die "Checksum verification failed for $url"
}

install_go() {
  local current=""
  has go && current="$(go version | awk '{sub(/^go/, "", $3); print $3}')"
  if [[ -n "$current" ]] && version_ge "$current" "$GO_VERSION"; then
    log "Go $current satisfies required $GO_VERSION"
    return
  fi
  $CHECK_ONLY && { warn "Go >= $GO_VERSION is required (found: ${current:-missing})"; return; }
  local arch archive url temp_dir install_dir
  arch="$(linux_arch)"
  archive="go${GO_VERSION}.linux-${arch}.tar.gz"
  url="https://go.dev/dl/${archive}"
  temp_dir="$(mktemp -d)"
  trap 'rm -rf -- "$temp_dir"' RETURN
  log "Installing Go $GO_VERSION to user space"
  download_verified "$url" "$temp_dir/$archive" "$url.sha256"
  install_dir="$HOME/.local/opt/go-$GO_VERSION"
  rm -rf -- "$install_dir"
  mkdir -p "$install_dir"
  tar -xzf "$temp_dir/$archive" -C "$install_dir" --strip-components=1
  ln -sfn "$install_dir" "$HOME/.local/go"
  trap - RETURN
  rm -rf -- "$temp_dir"
}

install_node() {
  local minimum="20.9.0" current=""
  has node && current="$(node --version | sed 's/^v//')"
  if [[ -n "$current" ]] && version_ge "$current" "$minimum"; then
    log "Node.js $current satisfies required >= $minimum"
    return
  fi
  $CHECK_ONLY && { warn "Node.js >= $minimum is required (found: ${current:-missing})"; return; }
  local arch archive base_url temp_dir expected actual install_dir
  arch="$(node_arch)"
  archive="node-v${NODE_VERSION}-linux-${arch}.tar.xz"
  base_url="https://nodejs.org/dist/v${NODE_VERSION}"
  temp_dir="$(mktemp -d)"
  trap 'rm -rf -- "$temp_dir"' RETURN
  log "Installing Node.js $NODE_VERSION to user space"
  curl --fail --location --silent --show-error "$base_url/$archive" --output "$temp_dir/$archive"
  curl --fail --location --silent --show-error "$base_url/SHASUMS256.txt" --output "$temp_dir/SHASUMS256.txt"
  expected="$(awk -v file="$archive" '$2 == file { print $1 }' "$temp_dir/SHASUMS256.txt")"
  actual="$(sha256sum "$temp_dir/$archive" | awk '{print $1}')"
  [[ -n "$expected" && "$actual" == "$expected" ]] || die "Checksum verification failed for Node.js $NODE_VERSION"
  install_dir="$HOME/.local/opt/node-$NODE_VERSION"
  rm -rf -- "$install_dir"
  mkdir -p "$install_dir"
  tar -xJf "$temp_dir/$archive" -C "$install_dir" --strip-components=1
  ln -sfn "$install_dir" "$HOME/.local/node"
  trap - RETURN
  rm -rf -- "$temp_dir"
}

check_commands() {
  local missing=0 command
  for command in curl git make sha256sum tar podman; do
    if ! has "$command"; then warn "$command is missing"; missing=1; fi
  done
  if has podman && ! podman compose version >/dev/null 2>&1; then
    warn "podman compose provider is unavailable"
    missing=1
  fi
  return "$missing"
}

if $CHECK_ONLY; then
  log "Checking Tourtect Linux environment"
  install_go
  install_node
  check_commands || exit 1
  log "Environment check completed"
  exit 0
fi

install_system_packages
write_shell_environment
install_go
install_node
hash -r
check_commands || die "One or more required commands are still missing. Review warnings above."

log "Installing pinned Web dependencies"
(cd "$ROOT_DIR/web" && npm ci)

log "Installing Go development tools"
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

if [[ ! -f "$ROOT_DIR/.env" ]]; then
  cp "$ROOT_DIR/.env.example" "$ROOT_DIR/.env"
  log "Created .env from .env.example; review local credentials before exposing any port."
else
  log "Keeping existing .env unchanged."
fi

if $INSTALL_BROWSER; then
  log "Installing Playwright Chromium"
  (cd "$ROOT_DIR/web" && npx playwright install chromium)
fi

log "Setup completed. Open a new shell or source $HOME/.config/tourtect/env.sh"
printf '[Tourtect] Start the system with: ./scripts/run-local.sh\n'

if $START_AFTER; then
  exec "$ROOT_DIR/scripts/run-local.sh"
fi
