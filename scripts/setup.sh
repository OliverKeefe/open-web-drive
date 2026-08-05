#!/usr/bin/env bash

# Configures the Open Web Drive local development environment.
# Usage: ./setup.sh [--KEY=VALUE]
#
# Flags override preset defaults.
# Writes backend/.env and frontend/.env

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Default flag values.
DATABASE_URL="postgres://admin:passwd@localhost:5433/metadatadb?sslmode=disable"
KEYCLOAK_ISSUER_URL="http://127.0.0.1:8080/realms/owd"
KEYCLOAK_JWKS_URL="http://127.0.0.1:8080/realms/owd/protocol/openid-connect/certs"
PORT="8081"
CORS_ALLOWED_ORIGINS="http://localhost:5173"
S3_ENDPOINT="http://localhost:4566"
AWS_REGION="us-east-1"
MAX_UPLOAD_BYTES="5242880"
BLOB_BUCKET_URL="s3://temp-buck?endpoint=http://localhost:4566&disable_https=true&s3ForcePathStyle=true"
VITE_KEYCLOAK_URL="http://127.0.0.1:8080"
VITE_KEYCLOAK_REALM="gestalt"
VITE_KEYCLOAK_CLIENT_ID="frontend-service"
VITE_BACKEND_BASE_URL="http://127.0.0.1:8081"
VITE_IPFS_GATEWAY_URL="https://bafybeia7wkemsgryogneimjafwwkb33ifwh2oo3djba3lqfeg3lkrqn464.ipfs.w3s.link/"
VITE_MAX_UPLOAD_SIZE_BYTES="16106127360"

BACKEND_VARS=(DATABASE_URL KEYCLOAK_ISSUER_URL KEYCLOAK_JWKS_URL PORT CORS_ALLOWED_ORIGINS \
    S3_ENDPOINT AWS_REGION MAX_UPLOAD_BYTES BLOB_BUCKET_URL)
FRONTEND_VARS=(VITE_KEYCLOAK_URL VITE_KEYCLOAK_REALM VITE_KEYCLOAK_CLIENT_ID \
    VITE_BACKEND_BASE_URL VITE_IPFS_GATEWAY_URL VITE_MAX_UPLOAD_SIZE_BYTES)

usage() {
  cat <<EOF
  Usage: $0 [--KEY=VALUE]

  Overrides default flags. Supported keys:
    backend : ${BACKEND_VARS[*]}
    frontend: ${FRONTEND_VARS[*]}
EOF
}

for arg in "$@"; do
  case "$arg" in
    -h|--help) usage; exit 0;;
    --*=*)
      key="${arg#--}"
      key="${key%%#*=}"
      value="${arg#*=}"
      case " ${BACKEND_VARS[*]} ${FRONTEND_VARS[*]} " in
        *" ${key} "*) ;;
        *) echo "[!] error: unknown variable '$key'" >&2; usage >&2; exit 1;;
      esac
      printf -v "$key" '%s' "$value"
      ;;
    *)
        echo "[!] error: expected --KEY=VALUE, got '$arg'" >&2; usage >&2; exit 1 ;;
  esac
done

# Writes the .env file per service (Frontend, Backend).
write_env() {
  local file="$1"; shift
  printf '# Generated automatically by setup script: scripts/setup.sh' > "$file"
  for var in "$@"; do
    printf '%s=%s\n' "$var" "${!var}" >> "$file"
  done
  echo "[+] created $file"
}

write_env "$REPO_ROOT/backend/.env" "${BACKEND_VARS[@]}"
write_env "$REPO_ROOT/frontend/.env" "${FRONTEND_VARS[@]}"

echo "[+] Local development environment configured successfully."