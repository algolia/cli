#!/usr/bin/env bash
set -euo pipefail

# Override to publish under a different name (e.g. cli_beta) for testing.
PACKAGE_NAME="${PACKAGE_NAME:-cli}"

# Strip leading 'v' from the tag (e.g. v1.8.2 -> 1.8.2)
VERSION="${VERSION#v}"

# Offset the major version by +4 so the npm package starts at 5.x
# (the old @algolia/cli package was at 4.x — a completely different tool)
IFS='.' read -r major minor patch <<< "$VERSION"
VERSION="$((major + 4)).$minor.$patch"

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NPM_DIR="$REPO_ROOT/npm"
DIST_DIR="$REPO_ROOT/dist"

DRY_RUN=""
if [[ "${1:-}" == "--dry-run" ]]; then
  DRY_RUN="--dry-run"
fi

# Provenance attestations need an OIDC issuer — only available in CI runners.
PROVENANCE=""
if [[ -n "${GITHUB_ACTIONS:-}" ]]; then
  PROVENANCE="--provenance"
fi

# Format: "platform-suffix:dist-relative-path"
# The directory under npm/ is npm/cli-<plat>; the published name is
# @algolia/${PACKAGE_NAME}-<plat>.
PLATFORMS=(
  "darwin-x64:algolia_darwin_amd64_v1/algolia"
  "darwin-arm64:algolia_darwin_arm64/algolia"
  "linux-x64:algolia_linux_amd64_v1/algolia"
  "linux-arm64:algolia_linux_arm64/algolia"
  "win32-x64:algolia_windows_amd64_v1/algolia.exe"
  "win32-arm64:algolia_windows_arm64/algolia.exe"
)

# When publishing under a non-default name, rewrite package.json names and the
# run.js shim in place, then revert via trap.
if [[ "$PACKAGE_NAME" != "cli" ]]; then
  BACKUP_DIR=$(mktemp -d)
  FILES_TO_MUTATE=(
    "$NPM_DIR/algolia/package.json"
    "$NPM_DIR/algolia/bin/run.js"
  )
  for entry in "${PLATFORMS[@]}"; do
    FILES_TO_MUTATE+=("$NPM_DIR/cli-${entry%%:*}/package.json")
  done

  cleanup() {
    for f in "${FILES_TO_MUTATE[@]}"; do
      bk="$BACKUP_DIR/$(echo "$f" | tr '/' '_')"
      [[ -f "$bk" ]] && cp "$bk" "$f"
    done
    rm -rf "$BACKUP_DIR"
  }
  trap cleanup EXIT

  for f in "${FILES_TO_MUTATE[@]}"; do
    cp "$f" "$BACKUP_DIR/$(echo "$f" | tr '/' '_')"
    sed -i.bak "s|@algolia/cli|@algolia/$PACKAGE_NAME|g" "$f"
    rm -f "$f.bak"
  done
fi

# Publish platform packages
for entry in "${PLATFORMS[@]}"; do
  plat="${entry%%:*}"
  dist_path="${entry#*:}"
  src="$DIST_DIR/$dist_path"
  dest_dir="$NPM_DIR/cli-$plat/bin"
  binary_name="$(basename "$dist_path")"

  echo "Publishing @algolia/$PACKAGE_NAME-$plat@$VERSION"

  if [[ -z "$DRY_RUN" ]]; then
    cp "$src" "$dest_dir/$binary_name"
    if [[ "$binary_name" != *.exe ]]; then
      chmod +x "$dest_dir/$binary_name"
    fi
  fi

  npm --prefix "$NPM_DIR/cli-$plat" version --no-git-tag-version "$VERSION"
  npm publish "$NPM_DIR/cli-$plat" --access public $PROVENANCE $DRY_RUN
done

# Update coordinator package versions to match and publish
for entry in "${PLATFORMS[@]}"; do
  plat="${entry%%:*}"
  # Pin the optionalDependency version in the coordinator package.json to match the published platform packages
  sed -i.bak "s|\"@algolia/$PACKAGE_NAME-$plat\": \"[^\"]*\"|\"@algolia/$PACKAGE_NAME-$plat\": \"$VERSION\"|g" "$NPM_DIR/algolia/package.json"
  rm -f "$NPM_DIR/algolia/package.json.bak"
done

npm --prefix "$NPM_DIR/algolia" version --no-git-tag-version "$VERSION"

echo "Publishing @algolia/$PACKAGE_NAME@$VERSION"
npm publish "$NPM_DIR/algolia" --access public $PROVENANCE $DRY_RUN
