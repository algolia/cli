#!/usr/bin/env bash
set -euo pipefail

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

# Format: "npm-package-name:dist-relative-path"
PLATFORMS=(
  "algolia-darwin-x64:algolia_darwin_amd64_v1/algolia"
  "algolia-darwin-arm64:algolia_darwin_arm64/algolia"
  "algolia-linux-x64:algolia_linux_amd64_v1/algolia"
  "algolia-linux-arm64:algolia_linux_arm64/algolia"
  "algolia-win32-x64:algolia_windows_amd64_v1/algolia.exe"
  "algolia-win32-arm64:algolia_windows_arm64/algolia.exe"
)

# Publish platform packages
for entry in "${PLATFORMS[@]}"; do
  pkg="${entry%%:*}"
  dist_path="${entry#*:}"
  src="$DIST_DIR/$dist_path"
  dest_dir="$NPM_DIR/$pkg/bin"
  binary_name="$(basename "$dist_path")"

  echo "Publishing @algolia/$pkg@$VERSION"

  if [[ -z "$DRY_RUN" ]]; then
    cp "$src" "$dest_dir/$binary_name"
    if [[ "$binary_name" != *.exe ]]; then
      chmod +x "$dest_dir/$binary_name"
    fi
  fi

  npm --prefix "$NPM_DIR/$pkg" version --no-git-tag-version "$VERSION"
  npm publish "$NPM_DIR/$pkg" --access public $DRY_RUN
done

# Update coordinator package versions to match and publish
for entry in "${PLATFORMS[@]}"; do
  pkg="${entry%%:*}"
  sed -i.bak "s|\"@algolia/$pkg\": \"[^\"]*\"|\"@algolia/$pkg\": \"$VERSION\"|g" "$NPM_DIR/algolia/package.json"
  rm -f "$NPM_DIR/algolia/package.json.bak"
done

npm --prefix "$NPM_DIR/algolia" version --no-git-tag-version "$VERSION"

echo "Publishing @algolia/cli@$VERSION"
npm publish "$NPM_DIR/algolia" --access public $DRY_RUN
