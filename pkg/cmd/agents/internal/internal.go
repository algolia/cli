package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
)

// NewInternalCmd is the parent for `algolia agents internal <verb>`.
//
// Verbs here wrap endpoints marked x-hidden in the OpenAPI spec.
// They are exposed for diagnostics / lab usage and are not stable —
// expect breaking changes without notice.
func NewInternalCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "internal",
		Short:  "Hidden Agent Studio diagnostics + memory internals (unstable)",
		Hidden: true,
		Long: heredoc.Doc(`
			Wraps endpoints marked x-hidden in the Agent Studio
			OpenAPI spec. NOT part of the documented public surface;
			may change without notice. Intended for support /
			diagnostics. memorize/ponder/consolidate hit the doubled
			path /1/agents/agents/{id}/<verb> (verified live — the
			single-path equivalents 404).
		`),
	}
	cmd.AddCommand(newStatusCmd(f, nil))
	cmd.AddCommand(newMemorizeCmd(f, nil))
	cmd.AddCommand(newPonderCmd(f, nil))
	cmd.AddCommand(newConsolidateCmd(f, nil))
	return cmd
}

func ctxOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// readJSONBody reads a JSON document from file or stdin and validates it.
// "-" reads stdin.
func readJSONBody(file string, ios *iostreams.IOStreams) ([]byte, error) {
	body, err := cmdutil.ReadFile(file, ios.In)
	if err != nil {
		return nil, fmt.Errorf("read body from %s: %w", shared.SourceLabel(file), err)
	}
	body = shared.TrimUTF8BOM(body)
	if !json.Valid(body) {
		return nil, cmdutil.FlagErrorf("body in %s is not valid JSON", shared.SourceLabel(file))
	}
	return body, nil
}
