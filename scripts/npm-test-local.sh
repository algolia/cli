#!/usr/bin/env bash
set -euo pipefail

# Local end-to-end test of the npm distribution flow:
#   1. Builds the Go binary for the current platform
#   2. Stages it in the matching npm/cli-<plat> package
#   3. Packs the platform package into a tarball
#   4. Rewrites the coordinator's optionalDependencies to reference that tarball
#   5. Installs the coordinator into a scratch dir
#   6. Runs `npx algolia --version` and `--help` to confirm the shim works
#
# All mutations are reverted on exit (success or failure).

# Override to test publishing under a different name (e.g. cli_beta).
PACKAGE_NAME="${PACKAGE_NAME:-cli}"

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PLAT=$(node -e 'console.log(process.platform + "-" + process.arch)')
PLAT_DIR="$REPO_ROOT/npm/cli-$PLAT"
COORD_DIR="$REPO_ROOT/npm/algolia"
RUN_JS="$COORD_DIR/bin/run.js"

if [[ ! -d "$PLAT_DIR" ]]; then
  echo "No platform package for $PLAT (looked for $PLAT_DIR)" >&2
  exit 1
fi

# Refuse to run if the files we'll mutate already have uncommitted changes —
# otherwise the cleanup trap will restore them to that dirty state instead of
# the committed state.
if ! git -C "$REPO_ROOT" diff --quiet -- "npm/algolia/package.json" "npm/cli-$PLAT/package.json" "npm/algolia/bin/run.js"; then
  echo "npm/algolia/package.json, npm/cli-$PLAT/package.json, or npm/algolia/bin/run.js has uncommitted changes." >&2
  echo "Commit or stash them first; this script will mutate and restore them." >&2
  exit 1
fi

TEST_VERSION=99.0.0
SCRATCH=$(mktemp -d)
TARBALL_DIR=$(mktemp -d)
COORD_PKG_BACKUP=$(mktemp)
PLAT_PKG_BACKUP=$(mktemp)
RUN_JS_BACKUP=$(mktemp)

cp "$COORD_DIR/package.json" "$COORD_PKG_BACKUP"
cp "$PLAT_DIR/package.json"  "$PLAT_PKG_BACKUP"
cp "$RUN_JS"                 "$RUN_JS_BACKUP"

cleanup() {
  echo "==> Cleaning up..."
  cp "$COORD_PKG_BACKUP" "$COORD_DIR/package.json"
  cp "$PLAT_PKG_BACKUP"  "$PLAT_DIR/package.json"
  cp "$RUN_JS_BACKUP"    "$RUN_JS"
  rm -f "$COORD_PKG_BACKUP" "$PLAT_PKG_BACKUP" "$RUN_JS_BACKUP"
  rm -rf "$SCRATCH" "$TARBALL_DIR" "$PLAT_DIR/bin"
}
trap cleanup EXIT

echo "==> Platform: $PLAT"
echo "==> Package name: @algolia/$PACKAGE_NAME"

if [[ "$PACKAGE_NAME" != "cli" ]]; then
  echo "==> Rewriting @algolia/cli -> @algolia/$PACKAGE_NAME..."
  for f in "$PLAT_DIR/package.json" "$COORD_DIR/package.json" "$RUN_JS"; do
    sed -i.bak "s|@algolia/cli|@algolia/$PACKAGE_NAME|g" "$f"
    rm -f "$f.bak"
  done
fi

echo "==> Building binary..."
mkdir -p "$PLAT_DIR/bin"
( cd "$REPO_ROOT" && go build -o "$PLAT_DIR/bin/algolia" ./cmd/algolia )
chmod +x "$PLAT_DIR/bin/algolia"

echo "==> Bumping platform package version to $TEST_VERSION..."
npm --prefix "$PLAT_DIR" version --no-git-tag-version --allow-same-version "$TEST_VERSION" >/dev/null

echo "==> Packing platform package..."
PLAT_TGZ_NAME=$(cd "$TARBALL_DIR" && npm pack "$PLAT_DIR" 2>/dev/null | tail -1)
PLAT_TGZ="$TARBALL_DIR/$PLAT_TGZ_NAME"
echo "    $PLAT_TGZ"

echo "==> Patching coordinator package.json..."
TEST_VERSION="$TEST_VERSION" PLAT="$PLAT" PLAT_TGZ="$PLAT_TGZ" COORD_PKG="$COORD_DIR/package.json" PACKAGE_NAME="$PACKAGE_NAME" node -e '
  const fs = require("fs");
  const p = process.env.COORD_PKG;
  const j = JSON.parse(fs.readFileSync(p));
  j.version = process.env.TEST_VERSION;
  j.optionalDependencies = { ["@algolia/" + process.env.PACKAGE_NAME + "-" + process.env.PLAT]: process.env.PLAT_TGZ };
  fs.writeFileSync(p, JSON.stringify(j, null, 2) + "\n");
'

echo "==> Packing coordinator package..."
COORD_TGZ_NAME=$(cd "$TARBALL_DIR" && npm pack "$COORD_DIR" 2>/dev/null | tail -1)
COORD_TGZ="$TARBALL_DIR/$COORD_TGZ_NAME"
echo "    $COORD_TGZ"

echo "==> Installing into scratch dir $SCRATCH..."
# Install both tarballs explicitly. Real users get the platform pkg via npm's
# optionalDependencies resolution; here we sidestep that (npm is flaky about
# file:-referenced optional deps) and install both directly. The shim's
# runtime behavior — require('@algolia/cli-<plat>') and execFileSync — is
# exercised the same way either path.
( cd "$SCRATCH" && npm init -y >/dev/null && npm install "$COORD_TGZ" "$PLAT_TGZ" )

echo "==> node_modules contents:"
ls "$SCRATCH/node_modules/@algolia"

echo "==> npx algolia --version"
( cd "$SCRATCH" && npx algolia --version )

echo "==> npx algolia --help (first 20 lines)"
( cd "$SCRATCH" && npx algolia --help | head -20 )

echo "==> OK"
