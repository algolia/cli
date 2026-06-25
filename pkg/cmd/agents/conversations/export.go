package conversations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type ExportOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioAPIClient func() (*agentStudio.APIClient, error)
	PrintFlags           *cmdutil.PrintFlags

	AgentID    string
	StartDate  string
	EndDate    string
	OutputFile string
}

func newExportCmd(f *cmdutil.Factory, runF func(*ExportOptions) error) *cobra.Command {
	opts := &ExportOptions{
		IO:                   f.IOStreams,
		AgentStudioAPIClient: f.AgentStudioAPIClient,
		PrintFlags:           cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "export <agent-id> [--start-date YYYY-MM-DD] [--end-date YYYY-MM-DD] [--output-file path]",
		Short: "Export conversations for an agent",
		Long: heredoc.Doc(`
			Dump every matching conversation for the agent.

			The response body is forwarded verbatim — the OpenAPI spec
			leaves the export shape unspecified, so the CLI does not pin
			a Go type. Use --output-file to write directly to disk
			(stdout otherwise).

			Pair with --start-date / --end-date to scope the dump.
		`),
		Example: heredoc.Doc(`
			$ algolia agents conversations export <agent-id> > backup.json
			$ algolia agents conversations export <agent-id> --output-file backup.json
			$ algolia agents conversations export <agent-id> --start-date 2026-01-01 --end-date 2026-01-31
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			if runF != nil {
				return runF(opts)
			}
			return runExportCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Export conversations >= date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "Export conversations <= date (YYYY-MM-DD)")
	cmd.Flags().
		StringVarP(&opts.OutputFile, "output-file", "O", "", "Write the export to this path instead of stdout")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runExportCmd(opts *ExportOptions) error {
	client, err := opts.AgentStudioAPIClient()
	if err != nil {
		return err
	}
	ctx := shared.OrBackground(opts.Ctx)

	req := client.NewApiExportConversationsRequest(opts.AgentID)
	if opts.StartDate != "" {
		req = req.WithStartDate(opts.StartDate)
	}
	if opts.EndDate != "" {
		req = req.WithEndDate(opts.EndDate)
	}

	// Forward the backend payload verbatim — the export shape is unspecified,
	// so use the raw response bytes rather than the typed model.
	opts.IO.StartProgressIndicatorWithLabel("Exporting conversations")
	raw, err := shared.RawResponse(
		client.ExportConversationsWithHTTPInfo(req, agentStudio.WithContext(ctx)),
	)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	// Pretty-print for humans, compact for files (jq-friendly).
	var payload []byte
	if opts.OutputFile != "" {
		payload = raw
	} else {
		var buf bytes.Buffer
		if err := json.Indent(&buf, raw, "", "  "); err != nil {
			payload = raw
		} else {
			payload = buf.Bytes()
		}
	}

	if opts.OutputFile != "" {
		if err := os.WriteFile(opts.OutputFile, payload, 0o600); err != nil {
			return fmt.Errorf("write export to %s: %w", opts.OutputFile, err)
		}
		if opts.IO.IsStdoutTTY() {
			fmt.Fprintf(opts.IO.Out, "Wrote %d byte(s) to %s\n", len(payload), opts.OutputFile)
		}
		return nil
	}

	_, _ = opts.IO.Out.Write(payload)
	if len(payload) == 0 || payload[len(payload)-1] != '\n' {
		_, _ = opts.IO.Out.Write([]byte("\n"))
	}
	return nil
}
