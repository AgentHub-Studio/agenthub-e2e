#!/usr/bin/env bash
# Build script para agenthub-e2e Go suite — ADR-006: builds via Docker
set -euo pipefail

GO_IMAGE="golang:1.24-alpine"
CACHE_VOL="$HOME/go/pkg/mod"
CMD="${1:-help}"
shift || true

case "$CMD" in
  compile)
    echo "==> Compilando e2e suite..."
    docker run --rm \
      -v "$(pwd)":/app \
      -v "${CACHE_VOL}":/go/pkg/mod \
      -w /app \
      "${GO_IMAGE}" go build ./...
    echo "==> OK"
    ;;
  test)
    echo "==> Executando testes unitários da suite e2e..."
    docker run --rm \
      -v "$(pwd)":/app \
      -v "${CACHE_VOL}":/go/pkg/mod \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -w /app \
      "${GO_IMAGE}" go test -v ./internal/... "$@"
    echo "==> OK"
    ;;
  e2e)
    echo "==> Executando testes E2E completos (requer serviços rodando)..."
    docker run --rm \
      -v "$(pwd)":/app \
      -v "${CACHE_VOL}":/go/pkg/mod \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -w /app \
      -e E2E=1 \
      -e API_URL="${API_URL:-http://host.docker.internal:8081}" \
      -e AUTH_TOKEN="${AUTH_TOKEN:-}" \
      "${GO_IMAGE}" go test -v -timeout 15m ./tests/... "$@"
    echo "==> OK"
    ;;
  tidy)
    docker run --rm \
      -v "$(pwd)":/app \
      -v "${CACHE_VOL}":/go/pkg/mod \
      -w /app \
      "${GO_IMAGE}" go mod tidy
    ;;
  *)
    echo "Usage: ./build.sh <compile|test|e2e|tidy>"
    ;;
esac
