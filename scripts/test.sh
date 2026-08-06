#!/usr/bin/env bash

# Runs tests suites and k8s provisioning tests.
# Usage: ./test.sh [--frontend] [--backend] [--infra]   (default: --all)
#
# - frontend : vitest            frontend/
# - backend  : go vet + go test  backend/
# - infra    : test-k8s.sh       backend/db/ (minikube k8s manifests) , this part is still in-progress

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

RUN_FRONTEND=false
RUN_BACKEND=false
RUN_INFRA=false

usage() {
  cat <<EOF
      Usage: $0 [OPTIONS]

        --frontend  Run frontend tests (vitest)
        --backend   Run backend tests (go vet + go test)
        --infra     Run infra tests (backend/db/test-k8s.sh)
        --all       Run every suite (default)
        -h, --help  Show this help
EOF
}

for arg in "$@"; do
  case "$arg" in
    --all) RUN_FRONTEND=true; RUN_BACKEND=true; RUN_INFRA=true ;;
    --frontend) RUN_FRONTEND=true ;;
    --backend) RUN_BACKEND=true ;;
    --infra) RUN_INFRA=true ;;
    -h|--help) usage; exit 0 ;;
  *) echo "[!] error: invalid option '$arg'" >&2; usage >&2; exit 1 ;;
  esac
done

# No flags given -> default to --all.
if [ "$RUN_FRONTEND" = false ] && [ "$RUN_BACKEND" = false ] && [ "$RUN_INFRA" = false ]; then
  RUN_FRONTEND=true; RUN_BACKEND=true; RUN_INFRA=true
fi

# runs commands in directory, subshell keeps current working dir unchanged though.
run_in() {
  local dir="$1"; shift
  ( cd "$dir" && "$@" )
}

pass_count=0
fail_count=0

# run a phase, records pass/fail but continue going on failure.
run_suite() {
  local name="$1"; shift
  echo; echo "===== $name ====="
  if "$@"; then
    echo "[+] $name: PASS"; pass_count=$((pass_count + 1))
  else
    echo "[!] $name: FAIL"; fail_count=$((fail_count + 1))
  fi
}

if "$RUN_FRONTEND"; then
  run_suite "frontend(vitest)" run_in $"REPO_ROOT/frontend" npx --no-install vitest run_suite
fi

if "$RUN_BACKEND"; then
  run_suite "backend (go vet)" run_in "$REPO_ROOT/backend" go vet ./...
  run_suite "backend (go test)" run_int "$REPO_ROOT/backend" go test ./...
fi

if [ "$RUN_INFRA" = true ]; then
  run_suite "infra tests" run_in "$REPO_ROOT/backend/db" ./test-k8s.sh
fi


echo
echo "[+] test suites passed: $pass_count"
echo "[!] test suites failed: $fail_count"
[ "$fail_count" -eq 0 ]