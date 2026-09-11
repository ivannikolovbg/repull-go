#!/usr/bin/env bash
# Regenerate the Repull Go client from the live OpenAPI spec.
#
# Requires:
#   go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
#
# Usage:
#   ./scripts/regen.sh
#
# The snapshot at openapi/v1.json is overwritten with whatever
# https://api.repull.dev/openapi.json returns. Commit both the snapshot and the
# regenerated *.gen.go files in the same change so the SDK and its source of
# truth stay in lockstep.
set -euo pipefail

cd "$(dirname "$0")/.."

SPEC_URL="${REPULL_OPENAPI_URL:-https://api.repull.dev/openapi.json}"
SNAPSHOT="openapi/v1.json"

echo "==> fetching $SPEC_URL"
curl -fsSL "$SPEC_URL" -o "$SNAPSHOT"
echo "    snapshot updated: $SNAPSHOT ($(wc -c <"$SNAPSHOT") bytes)"

if ! command -v oapi-codegen >/dev/null 2>&1; then
  echo "oapi-codegen not on PATH. install with:"
  echo "  go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest"
  exit 1
fi

echo "==> generating types"
( cd repull && oapi-codegen -config cfg-types.yaml ../"$SNAPSHOT" )

echo "==> generating client"
( cd repull && oapi-codegen -config cfg-client.yaml ../"$SNAPSHOT" )

echo "==> go mod tidy"
go mod tidy

echo "==> go build"
go build ./...

echo "==> go vet"
go vet ./...

echo "==> done"
