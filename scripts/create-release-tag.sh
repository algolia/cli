#!/usr/bin/env bash
set -euo pipefail

RELEASE_TYPE="${RELEASE_TYPE:-}"
PRERELEASE_BASE="${PRERELEASE_BASE:-minor}"
PRERELEASE_CHANNEL="${PRERELEASE_CHANNEL:-beta}"
CUSTOM_PRERELEASE_CHANNEL="${CUSTOM_PRERELEASE_CHANNEL:-}"
FREEFORM_VERSION="${FREEFORM_VERSION:-}"
MAJOR_CONFIRMATION="${MAJOR_CONFIRMATION:-}"
DRY_RUN="${DRY_RUN:-false}"

SEMVER_RE='^v?([0-9]+)\.([0-9]+)\.([0-9]+)(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$'
STABLE_TAG_RE='^v([0-9]+)\.([0-9]+)\.([0-9]+)$'

fail() {
  echo "Error: $*" >&2
  exit 1
}

is_true() {
  case "$1" in
    true | TRUE | True | 1 | yes | YES | Yes) return 0 ;;
    *) return 1 ;;
  esac
}

validate_release_type() {
  case "$RELEASE_TYPE" in
    fix | minor | major | prerelease | freeform) ;;
    "") fail "RELEASE_TYPE is required." ;;
    *) fail "Unsupported RELEASE_TYPE '$RELEASE_TYPE'." ;;
  esac
}

validate_base() {
  case "$PRERELEASE_BASE" in
    fix | minor | major) ;;
    *) fail "Unsupported PRERELEASE_BASE '$PRERELEASE_BASE'." ;;
  esac
}

normalize_freeform_tag() {
  local version="$1"

  [[ -n "$version" ]] || fail "FREEFORM_VERSION is required when RELEASE_TYPE is freeform."
  [[ "$version" =~ $SEMVER_RE ]] || fail "FREEFORM_VERSION must be semver, for example v1.12.0 or v1.12.0-beta.1; got '$version'."

  version="${version#v}"
  echo "v$version"
}

latest_stable_tag() {
  local latest=""
  local latest_major=-1
  local latest_minor=-1
  local latest_patch=-1
  local tag major minor patch

  while IFS= read -r tag; do
    if [[ "$tag" =~ $STABLE_TAG_RE ]]; then
      major="${BASH_REMATCH[1]}"
      minor="${BASH_REMATCH[2]}"
      patch="${BASH_REMATCH[3]}"

      if (( major > latest_major )) ||
        (( major == latest_major && minor > latest_minor )) ||
        (( major == latest_major && minor == latest_minor && patch > latest_patch )); then
        latest="$tag"
        latest_major="$major"
        latest_minor="$minor"
        latest_patch="$patch"
      fi
    fi
  done < <(git tag --list)

  [[ -n "$latest" ]] || fail "No stable release tags found."
  echo "$latest"
}

bump_tag() {
  local base_tag="$1"
  local bump="$2"

  [[ "$base_tag" =~ $STABLE_TAG_RE ]] || fail "Cannot bump invalid stable tag '$base_tag'."

  local major="${BASH_REMATCH[1]}"
  local minor="${BASH_REMATCH[2]}"
  local patch="${BASH_REMATCH[3]}"

  case "$bump" in
    fix)
      patch=$((patch + 1))
      ;;
    minor)
      minor=$((minor + 1))
      patch=0
      ;;
    major)
      major=$((major + 1))
      minor=0
      patch=0
      ;;
    *)
      fail "Unsupported bump '$bump'."
      ;;
  esac

  echo "v$major.$minor.$patch"
}

resolve_prerelease_channel() {
  local channel="$PRERELEASE_CHANNEL"

  if [[ "$channel" == "custom" ]]; then
    channel="$CUSTOM_PRERELEASE_CHANNEL"
  fi

  [[ -n "$channel" ]] || fail "A custom prerelease channel is required when PRERELEASE_CHANNEL is custom."
  [[ "$channel" =~ ^[A-Za-z][0-9A-Za-z-]*$ ]] || fail "Prerelease channel must start with a letter and contain only letters, numbers, or hyphens; got '$channel'."

  echo "$channel"
}

next_prerelease_tag() {
  local base_tag="$1"
  local channel="$2"
  local base_version="${base_tag#v}"
  local highest=0
  local tag escaped_channel

  escaped_channel="${channel//-/\\-}"

  while IFS= read -r tag; do
    if [[ "$tag" =~ ^v${base_version}-${escaped_channel}\.([0-9]+)$ ]]; then
      if (( BASH_REMATCH[1] > highest )); then
        highest="${BASH_REMATCH[1]}"
      fi
    fi
  done < <(git tag --list)

  echo "v${base_version}-${channel}.$((highest + 1))"
}

tag_exists_local() {
  git rev-parse -q --verify "refs/tags/$1" >/dev/null
}

tag_exists_remote() {
  local tag="$1"
  local status

  if git ls-remote --exit-code --tags origin "refs/tags/$tag" >/dev/null 2>&1; then
    return 0
  else
    status=$?
  fi

  if [[ "$status" == "2" ]]; then
    return 1
  fi

  return "$status"
}

ensure_tag_available() {
  local tag="$1"

  if tag_exists_local "$tag"; then
    fail "Tag '$tag' already exists locally."
  fi

  local remote_status
  if tag_exists_remote "$tag"; then
    fail "Tag '$tag' already exists on origin."
  else
    remote_status=$?
  fi

  if (( remote_status > 1 )); then
    fail "Unable to check whether tag '$tag' exists on origin."
  fi
}

write_outputs() {
  local latest_stable="$1"
  local tag="$2"
  local channel="$3"
  local target_sha="$4"
  local pushed="$5"

  if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
    {
      echo "latest_stable=$latest_stable"
      echo "tag=$tag"
      echo "prerelease_channel=$channel"
      echo "target_sha=$target_sha"
      echo "pushed=$pushed"
    } >>"$GITHUB_OUTPUT"
  fi

  if [[ -n "${GITHUB_STEP_SUMMARY:-}" ]]; then
    {
      echo "## Release tag summary"
      echo
      echo "- Release type: \`$RELEASE_TYPE\`"
      echo "- Latest stable tag: \`$latest_stable\`"
      echo "- Computed tag: \`$tag\`"
      echo "- Target SHA: \`$target_sha\`"
      echo "- Dry run: \`$DRY_RUN\`"
      echo "- Pushed: \`$pushed\`"
      if [[ -n "$channel" ]]; then
        echo "- Prerelease channel: \`$channel\`"
      fi
    } >>"$GITHUB_STEP_SUMMARY"
  fi
}

validate_release_type
latest_stable="$(latest_stable_tag)"
prerelease_channel=""

case "$RELEASE_TYPE" in
  fix | minor | major)
    tag="$(bump_tag "$latest_stable" "$RELEASE_TYPE")"
    ;;
  prerelease)
    validate_base
    prerelease_channel="$(resolve_prerelease_channel)"
    base_tag="$(bump_tag "$latest_stable" "$PRERELEASE_BASE")"
    tag="$(next_prerelease_tag "$base_tag" "$prerelease_channel")"
    ;;
  freeform)
    tag="$(normalize_freeform_tag "$FREEFORM_VERSION")"
    ;;
esac

if [[ "$RELEASE_TYPE" == "major" && "$MAJOR_CONFIRMATION" != "$tag" ]]; then
  fail "Major releases require MAJOR_CONFIRMATION to exactly match '$tag'."
fi

ensure_tag_available "$tag"

target_sha="$(git rev-parse --verify HEAD)"
pushed="false"

if ! is_true "$DRY_RUN"; then
  git tag -a "$tag" "$target_sha" -m "$tag"
  git push origin "$tag"
  pushed="true"
fi

echo "Latest stable tag: $latest_stable"
echo "Computed tag: $tag"
echo "Target SHA: $target_sha"
echo "Dry run: $DRY_RUN"
echo "Pushed: $pushed"

write_outputs "$latest_stable" "$tag" "$prerelease_channel" "$target_sha" "$pushed"
