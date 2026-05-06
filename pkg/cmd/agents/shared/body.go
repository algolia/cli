// Package shared holds helpers reused across `algolia agents` subcommands
// that pass user-supplied JSON through to the Agent Studio backend
// (currently create + update; future provider/tool commands will reuse the
// same shape). Kept tight on purpose — only extract here when there are
// at least two real callers.
package shared

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
)

// SourceLabel returns "stdin" for "-", otherwise the path itself. Used in
// help text and error messages so users can tell where a body came from.
func SourceLabel(file string) string {
	if file == "-" {
		return "stdin"
	}
	return file
}

// TrimUTF8BOM strips a leading UTF-8 byte-order-mark if present. Common
// when users hand-edit JSON in Notepad or VS Code on Windows; json.Valid
// rejects BOM-prefixed input with a confusing error otherwise.
func TrimUTF8BOM(b []byte) []byte {
	const bom = "\xef\xbb\xbf"
	return []byte(strings.TrimPrefix(string(b), bom))
}

// PrintDryRun renders a preview of an HTTP request that would have been
// sent. When wantsStructured is true (caller derives this from
// `cmd.Flags().Changed("output")` so the user must opt in explicitly,
// not just inherit a `WithDefaultOutput("json")` default), it emits a
// structured summary. Otherwise it prints the request line + the
// resolved JSON body pretty-printed — the human form intentionally
// surfaces the *full* body so users can lint it visually before
// re-running without --dry-run.
//
// extra lets a caller add command-specific keys to the structured summary
// (e.g., "agentId" for update). Pass nil if there are none.
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
		// Body already passed json.Valid upstream; the only realistic
		// failure here is OOM on a giant doc — fall back to raw bytes.
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
