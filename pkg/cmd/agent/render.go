package agent

import (
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"

	"github.com/algolia/cli/pkg/iostreams"
)


// renderMarkdown converts a markdown string into ANSI-styled terminal output
// by parsing it with goldmark and walking the AST.
func renderMarkdown(cs *iostreams.ColorScheme, input string) string {
	source := []byte(input)

	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
	)
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	var out strings.Builder
	renderNode(&out, cs, doc, source)
	return strings.TrimRight(out.String(), "\n")
}

// renderNode recursively walks the AST and writes ANSI-styled text to out.
func renderNode(out *strings.Builder, cs *iostreams.ColorScheme, n ast.Node, source []byte) {
	switch n.Kind() {

	case ast.KindDocument:
		renderChildren(out, cs, n, source)

	case ast.KindHeading:
		var headingText strings.Builder
		renderChildrenTo(&headingText, cs, n, source)
		out.WriteString(cs.Bold(headingText.String()))
		out.WriteString("\n\n")

	case ast.KindParagraph:
		renderChildren(out, cs, n, source)
		out.WriteString("\n\n")

	case ast.KindTextBlock:
		renderChildren(out, cs, n, source)
		out.WriteString("\n")

	case ast.KindText:
		t := n.(*ast.Text)
		out.Write(t.Segment.Value(source))
		if t.SoftLineBreak() {
			out.WriteString("\n")
		}
		if t.HardLineBreak() {
			out.WriteString("\n")
		}

	case ast.KindEmphasis:
		e := n.(*ast.Emphasis)
		var content strings.Builder
		renderChildrenTo(&content, cs, n, source)
		if e.Level == 2 {
			out.WriteString(cs.Bold(content.String()))
		} else {
			out.WriteString(cs.Gray(content.String()))
		}

	case ast.KindCodeSpan:
		var code strings.Builder
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			if child.Kind() == ast.KindText {
				t := child.(*ast.Text)
				code.Write(t.Segment.Value(source))
			}
		}
		out.WriteString(colorCodeSpan(cs, code.String()))

	case ast.KindFencedCodeBlock, ast.KindCodeBlock:
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			seg := lines.At(i)
			line := strings.TrimRight(string(seg.Value(source)), "\n")
			if idx := strings.Index(line, "#"); idx >= 0 {
				cmd := line[:idx]
				comment := line[idx:]
				out.WriteString(cs.Blue(cmd))
				out.WriteString(cs.Green(comment))
			} else {
				out.WriteString(cs.Blue(line))
			}
			out.WriteString("\n")
		}
		out.WriteString("\n")

	case ast.KindList:
		renderChildren(out, cs, n, source)
		out.WriteString("\n")

	case ast.KindListItem:
		out.WriteString("- ")
		// Render list item children inline (skip the paragraph newlines).
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			if child.Kind() == ast.KindParagraph {
				renderChildren(out, cs, child, source)
			} else {
				renderNode(out, cs, child, source)
			}
		}
		out.WriteString("\n")

	case ast.KindThematicBreak:
		out.WriteString("───\n\n")

	case east.KindTable:
		renderTable(out, cs, n, source)

	case ast.KindString:
		s := n.(*ast.String)
		out.Write(s.Value)

	default:
		// For any unhandled node, just render its children.
		renderChildren(out, cs, n, source)
	}
}

// renderChildren renders all children of a node to the shared output.
func renderChildren(out *strings.Builder, cs *iostreams.ColorScheme, n ast.Node, source []byte) {
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		renderNode(out, cs, child, source)
	}
}

// codeTokenRe matches placeholders (<...>) and flags (--word).
var codeTokenRe = regexp.MustCompile(`<[^>]+>|\[[^\]]+\]|--\w[\w-]*|-\w\b`)

// colorCodeSpan colors inline code with different colors per token type:
// - placeholders (<index>) in light blue
// - flags (--query) in cyan
// - everything else in Algolia Blue
func colorCodeSpan(cs *iostreams.ColorScheme, code string) string {
	var out strings.Builder
	last := 0
	for _, match := range codeTokenRe.FindAllStringIndex(code, -1) {
		if match[0] > last {
			out.WriteString(cs.Blue(code[last:match[0]]))
		}
		token := code[match[0]:match[1]]
		out.WriteString(cs.Cyan(token))
		last = match[1]
	}
	if last < len(code) {
		out.WriteString(cs.Blue(code[last:]))
	}
	return out.String()
}

// renderChildrenTo renders all children into a separate builder (for wrapping in styles).
func renderChildrenTo(buf *strings.Builder, cs *iostreams.ColorScheme, n ast.Node, source []byte) {
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		renderNode(buf, cs, child, source)
	}
}

// renderTable collects rows from a goldmark Table node and renders via tablewriter.
func renderTable(out *strings.Builder, cs *iostreams.ColorScheme, n ast.Node, source []byte) {
	var header []string
	var dataRows [][]string

	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		var cells []string
		for cell := child.FirstChild(); cell != nil; cell = cell.NextSibling() {
			var cellBuf strings.Builder
			renderChildrenTo(&cellBuf, cs, cell, source)
			cells = append(cells, cellBuf.String())
		}

		if child.Kind() == east.KindTableHeader {
			header = cells
		} else {
			dataRows = append(dataRows, cells)
		}
	}

	var buf strings.Builder
	table := tablewriter.NewWriter(&buf)

	if header != nil {
		plain := make([]any, len(header))
		for i, h := range header {
			plain[i] = h
		}
		table.Header(plain...)
	}

	for _, row := range dataRows {
		styled := make([]any, len(row))
		for i, cell := range row {
			styled[i] = cell
		}
		_ = table.Append(styled...)
	}

	_ = table.Render()
	out.WriteString(strings.TrimRight(buf.String(), "\n"))
	out.WriteString("\n\n")
}
