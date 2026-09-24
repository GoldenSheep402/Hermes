#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

echo "== go mod verify =="
go mod verify

echo "== go test =="
go test ./... -count=1

if command -v govulncheck >/dev/null 2>&1; then
  echo "== govulncheck =="
  govulncheck ./...
else
  echo "== govulncheck skipped (install: go install golang.org/x/vuln/cmd/govulncheck@latest) =="
fi

echo "OK"
