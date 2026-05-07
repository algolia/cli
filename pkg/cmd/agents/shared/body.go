// Package shared holds helpers reused across `algolia agents` subcommands.
// Extract on second use, not pre-emptively. See docs/agents.md.
package shared

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
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

// PrintDryRun renders a preview of an HTTP request that would have been
// sent. wantsStructured must be derived from cmd.Flags().Changed("output")
// (not PrintFlags.HasStructuredOutput) — see docs/agents.md "On --dry-run".
// extra adds command-specific keys to the structured summary; pass nil if none.
func PrintDryRun(
	io *iostreams.IOStreams,
	pf *cmdutil.PrintFlags,
	wantsStructured bool,
	action, request, file string,
	body []byte,
	extra map[string]any,
) error {
	if wantsStructured && pf != nil {
		summary := map[string]any{
			"action":  action,
			"request": request,
			"source":  SourceLabel(file),
			"bytes":   len(body),
			"body":    json.RawMessage(body),
			"dryRun":  true,
		}
		for k, v := range extra {
			summary[k] = v
		}
		return pf.Print(io, summary)
	}

	fmt.Fprintf(io.Out, "Dry run: would %s (%d bytes from %s)\n",
		request, len(body), SourceLabel(file))

	pretty, err := json.MarshalIndent(json.RawMessage(body), "", "  ")
	if err != nil {
		pretty = body
	}
	if _, err := io.Out.Write(pretty); err != nil {
		return err
	}
	if len(pretty) == 0 || pretty[len(pretty)-1] != '\n' {
		fmt.Fprintln(io.Out)
	}
	return nil
}
