package text

import (
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/truncate"
)

const (
	ellipsis            = "..."
	minWidthForEllipsis = len(ellipsis) + 2
)

// DisplayWidth calculates what the rendered width of a string may be
func DisplayWidth(s string) int {
	return ansi.PrintableRuneWidth(s)
}

// Truncate shortens a string to fit the maximum display width
func Truncate(maxWidth int, s string) string {
	w := DisplayWidth(s)
	if w <= maxWidth {
		return s
	}

	tail := ""
	if maxWidth >= minWidthForEllipsis {
		tail = ellipsis
	}

	// Guard against overflow
	if maxWidth < 0 {
		maxWidth = 0
	}

	// Seems to be a false positive from gosec
	// since max(uint) > max(int) && maxWidth >= 0 at this point
	// nolint:gosec
	r := truncate.StringWithTail(s, uint(maxWidth), tail)
	if DisplayWidth(r) < maxWidth {
		r += " "
	}

	return r
}

// TruncateColumn replaces the first new line character with an ellipsis
// and shortens a string to fit the maximum display width
func TruncateColumn(maxWidth int, s string) string {
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		s = s[:i] + ellipsis
	}
	return Truncate(maxWidth, s)
}
