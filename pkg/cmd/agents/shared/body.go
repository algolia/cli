// Package shared holds helpers reused across `algolia agents` subcommands.
// Extract on second use, not pre-emptively. See docs/agents.md.
package shared

import (
	"strings"
)

// SourceLabel returns "stdin" for "-", otherwise the path itself.
func SourceLabel(file string) string {
	if file == "-" {
		return "stdin"
	}
	return file
}

// TrimUTF8BOM strips a leading UTF-8 byte-order-mark if present.
func TrimUTF8BOM(b []byte) []byte {
	const bom = "\xef\xbb\xbf"
	return []byte(strings.TrimPrefix(string(b), bom))
}
