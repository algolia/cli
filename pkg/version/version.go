package version

import (
	"fmt"
)

// Version of the CLI.
// This is set to the actual version by GoReleaser, identify by the
// git tag assigned to the release. Versions built from source will
// always show main.
var Version = "main"

// Distribution labels the packaged binary variant (non-release builds).
// Injected via -X for local "algolia-beta" builds; empty for Goreleaser and task build.
var Distribution string

// Template for the version string.
var Template = fmt.Sprintf("algolia version %s\n", Version)
