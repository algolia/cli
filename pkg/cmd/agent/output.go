package agent

import (
	"fmt"
	"strings"
)

// stripANSI removes ANSI escape sequences from a string and simulates
// carriage return behavior (overwrites the current line).
func stripANSI(s string) string {
	var lines []string
	var cur strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\033' {
			// Skip CSI sequences: ESC [ ... final byte
			if i+1 < len(s) && s[i+1] == '[' {
				j := i + 2
				for j < len(s) && s[j] >= 0x20 && s[j] <= 0x3F {
					j++
				}
				if j < len(s) {
					j++ // skip final byte
				}
				i = j
				continue
			}
			// Skip other ESC sequences (ESC + one byte)
			i += 2
			continue
		}
		if s[i] == '\r' {
			// \r\n is a normal newline, not a spinner overwrite.
			if i+1 < len(s) && s[i+1] == '\n' {
				lines = append(lines, cur.String())
				cur.Reset()
				i += 2
				continue
			}
			// Standalone \r: discard current line content (spinner overwrite)
			cur.Reset()
			i++
			continue
		}
		if s[i] == '\n' {
			lines = append(lines, cur.String())
			cur.Reset()
			i++
			continue
		}
		cur.WriteByte(s[i])
		i++
	}
	if cur.Len() > 0 {
		lines = append(lines, cur.String())
	}
	// Filter out empty lines from spinner artifacts.
	var result []string
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			result = append(result, l)
		}
	}
	return strings.Join(result, "\n")
}

// compactJSON attempts to compact each JSON object in the output to a single line.
// Non-JSON lines are left unchanged.
func compactJSON(s string) string {
	lines := strings.Split(s, "\n")
	var result []string
	var jsonBuf strings.Builder
	depth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" && depth == 0 {
			continue
		}

		for _, ch := range trimmed {
			switch ch {
			case '{', '[':
				depth++
			case '}', ']':
				depth--
			}
		}

		if depth > 0 {
			jsonBuf.WriteString(trimmed)
			continue
		}

		if jsonBuf.Len() > 0 {
			jsonBuf.WriteString(trimmed)
			result = append(result, jsonBuf.String())
			jsonBuf.Reset()
		} else {
			result = append(result, trimmed)
		}
	}

	if jsonBuf.Len() > 0 {
		result = append(result, jsonBuf.String())
	}

	return strings.Join(result, "\n")
}

// truncateOutput limits the output to maxLines non-empty lines.
// If truncated, appends a note indicating how many lines were omitted.
func truncateOutput(s string, maxLines int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= maxLines {
		return s
	}
	truncated := strings.Join(lines[:maxLines], "\n")
	return truncated + fmt.Sprintf("\n[... %d more lines truncated]", len(lines)-maxLines)
}
