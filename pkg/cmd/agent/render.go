package agent

import (
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/algolia/cli/pkg/iostreams"
)

// algoliaBlue is the Algolia Blue brand color.
const algoliaBlue = "3970ff"

// renderMarkdown converts a markdown string into ANSI-styled terminal output.
// It handles: headers, bold, italic, inline code, fenced code blocks, and tables.
func renderMarkdown(cs *iostreams.ColorScheme, text string) string {
	lines := strings.Split(text, "\n")
	var out []string
	var tableBuffer []string
	inCodeBlock := false

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock {
			out = append(out, cs.HexToRGB(algoliaBlue, line))
			continue
		}

		if strings.HasPrefix(strings.TrimSpace(line), "|") {
			tableBuffer = append(tableBuffer, line)
			continue
		}

		if len(tableBuffer) > 0 {
			out = append(out, renderTable(cs, tableBuffer))
			tableBuffer = nil
		}

		if strings.HasPrefix(line, "#") {
			stripped := strings.TrimLeft(line, "# ")
			out = append(out, cs.Bold(stripped))
			continue
		}

		line = renderInline(cs, line)
		out = append(out, line)
	}

	if len(tableBuffer) > 0 {
		out = append(out, renderTable(cs, tableBuffer))
	}

	return strings.Join(out, "\n")
}

// Bold regex: **text**
var boldRe = regexp.MustCompile(`\*\*(.+?)\*\*`)

// Italic regex: *text* (but not **text**)
var italicRe = regexp.MustCompile(`(?:^|[^*])\*([^*]+?)\*(?:[^*]|$)`)

// Inline code regex: `text`
var codeRe = regexp.MustCompile("`([^`]+)`")

// renderInline applies bold, italic, and inline code styling to a single line.
func renderInline(cs *iostreams.ColorScheme, line string) string {
	line = boldRe.ReplaceAllStringFunc(line, func(match string) string {
		inner := boldRe.FindStringSubmatch(match)[1]
		return cs.Bold(inner)
	})
	line = codeRe.ReplaceAllStringFunc(line, func(match string) string {
		inner := codeRe.FindStringSubmatch(match)[1]
		return cs.HexToRGB(algoliaBlue, inner)
	})
	line = italicRe.ReplaceAllStringFunc(line, func(match string) string {
		inner := italicRe.FindStringSubmatch(match)[1]
		prefix := ""
		suffix := ""
		if len(match) > 0 && match[0] != '*' {
			prefix = string(match[0])
		}
		if len(match) > 0 && match[len(match)-1] != '*' {
			suffix = string(match[len(match)-1])
		}
		return prefix + cs.Gray(inner) + suffix
	})
	return line
}

// isSeparatorRow checks if a table row is a markdown separator (e.g. |---|---|).
func isSeparatorRow(line string) bool {
	for _, cell := range parseTableRow(line) {
		stripped := strings.Trim(cell, " :-")
		if stripped != "" {
			return false
		}
	}
	return true
}

// parseTableRow splits a markdown table row into cell values.
func parseTableRow(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.Trim(line, "|")
	parts := strings.Split(line, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// renderTable renders markdown table lines using tablewriter.
func renderTable(cs *iostreams.ColorScheme, lines []string) string {
	var header []string
	var dataRows [][]string
	foundSep := false

	for _, line := range lines {
		if isSeparatorRow(line) {
			foundSep = true
			continue
		}
		cells := parseTableRow(line)
		if !foundSep && header == nil {
			header = cells
		} else {
			dataRows = append(dataRows, cells)
		}
	}

	var buf strings.Builder
	table := tablewriter.NewWriter(&buf)

	// Tablewriter auto-formats headers (uppercase + bold).
	if header != nil {
		plain := make([]any, len(header))
		for i, h := range header {
			plain[i] = h
		}
		table.Header(plain...)
	}

	// Apply inline styling to data cells and append rows.
	for _, row := range dataRows {
		styled := make([]any, len(row))
		for i, cell := range row {
			styled[i] = renderInline(cs, cell)
		}
		table.Append(styled...)
	}

	table.Render()

	return strings.TrimRight(buf.String(), "\n")
}
