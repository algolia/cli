package config

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// NewConfigCmd is the parent for `algolia agents config <verb>`.
//
// Wraps the app-wide /1/configuration endpoints. Today the only
// settable field is maxRetentionDays (data retention for conversations
// and analytics). Kept as a sub-group rather than a flat command
// because the backend's spec models it as a resource (PATCH semantics)
// and future fields will land there. Same pattern as `agents cache`.
func NewConfigCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage app-wide Agent Studio configuration",
		Long: heredoc.Doc(`
			Get and set the app-wide Agent Studio configuration.
			Currently surfaces the data-retention setting; future fields
			will land here as the backend grows them.
		`),
	}
	cmd.AddCommand(newGetCmd(f, nil))
	cmd.AddCommand(newSetCmd(f, nil))
	return cmd
}

// ---------------------------------------------------------------------
// get
// ---------------------------------------------------------------------

type GetOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags
}

func newGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Show the current Agent Studio configuration",
		Long: heredoc.Doc(`
			Calls GET /1/configuration. Note: this endpoint requires the
			` + "`logs`" + ` ACL on the API key, not ` + "`settings`" + ` — the
			single field today (maxRetentionDays) governs log /
			conversation retention, hence the unusual ACL.
		`),
		Example: heredoc.Doc(`
			$ algolia agents config get
			$ algolia agents config get --output json | jq .maxRetentionDays
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			if runF != nil {
				return runF(opts)
			}
			return runGetCmd(opts)
		},
	}
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching configuration")
	cfg, err := client.GetConfiguration(ctx)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	return opts.PrintFlags.Print(opts.IO, cfg)
}

// ---------------------------------------------------------------------
// set
// ---------------------------------------------------------------------

type SetOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	// RetentionDays is only meaningful when the --retention-days flag
	// was actually set; the dispatcher checks cmd.Flags().Changed
	// rather than the value to detect intent. 0 is a valid backend
	// value (allowed set is [0, 30, 60, 90]).
	RetentionDays int
	File          string
	DryRun        bool
	OutputChanged bool
}

func newSetCmd(f *cmdutil.Factory, runF func(*SetOptions) error) *cobra.Command {
	opts := &SetOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "set [--retention-days N] [-F file]",
		Short: "Patch the Agent Studio configuration",
		Long: heredoc.Doc(`
			Patch the app-wide Agent Studio configuration via PATCH
			/1/configuration. Two convenience surfaces:

			  --retention-days N   Set maxRetentionDays. Backend accepts
			                       0, 30, 60, or 90; anything else 422s
			                       with the structured detail.

			  -F file              Send any JSON patch verbatim. Useful
			                       for future fields the CLI doesn't
			                       know about yet.

			Exactly one of the two must be provided. Use --dry-run to
			preview without sending.
		`),
		Example: heredoc.Doc(`
			$ algolia agents config set --retention-days 30
			$ algolia agents config set --retention-days 90 --dry-run
			$ algolia agents config set -F retention.json
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			opts.OutputChanged = cmd.Flags().Changed("output")

			fileSet := cmd.Flags().Changed("file")
			retSet := cmd.Flags().Changed("retention-days")
			if !fileSet && !retSet {
				return cmdutil.FlagErrorf("one of --retention-days or --file is required")
			}
			if fileSet && retSet {
				return cmdutil.FlagErrorf("--retention-days and --file are mutually exclusive")
			}

			if runF != nil {
				return runF(opts)
			}
			return runSetCmd(opts)
		},
	}

	cmd.Flags().
		IntVar(&opts.RetentionDays, "retention-days", 0, "Set maxRetentionDays (0, 30, 60, or 90 — backend-validated)")
	cmd.Flags().
		StringVarP(&opts.File, "file", "F", "", "JSON patch body (use \"-\" for stdin)")
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Validate and print the resolved request body without calling the API")
	opts.PrintFlags.AddFlags(cmd)

	// pflag suppresses "(default ...)" for ints when DefValue == "0".
	// Defaulting to 0 (the natural zero value) keeps the help text
	// clean without bespoke override. The dispatcher distinguishes
	// "user passed --retention-days 0" from "flag omitted" via
	// cmd.Flags().Changed, not the value itself.

	cmd.MarkFlagsMutuallyExclusive("retention-days", "file")
	return cmd
}

func runSetCmd(opts *SetOptions) error {
	body, source, err := buildBody(opts)
	if err != nil {
		return err
	}

	if opts.DryRun {
		return shared.PrintDryRun(opts.IO, opts.PrintFlags, opts.OutputChanged,
			"update_configuration", "PATCH /1/configuration", source, body, nil)
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	opts.IO.StartProgressIndicatorWithLabel("Updating configuration")
	cfg, err := client.UpdateConfiguration(ctx, json.RawMessage(body))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	return opts.PrintFlags.Print(opts.IO, cfg)
}

func buildBody(opts *SetOptions) ([]byte, string, error) {
	if opts.File != "" {
		body, err := cmdutil.ReadFile(opts.File, opts.IO.In)
		if err != nil {
			return nil, "", fmt.Errorf("failed to read configuration body from %s: %w",
				shared.SourceLabel(opts.File), err)
		}
		body = shared.TrimUTF8BOM(body)
		if !json.Valid(body) {
			return nil, "", cmdutil.FlagErrorf("configuration body in %s is not valid JSON",
				shared.SourceLabel(opts.File))
		}
		return body, opts.File, nil
	}
	// Construct the body from --retention-days.
	body, err := json.Marshal(map[string]any{"maxRetentionDays": opts.RetentionDays})
	if err != nil {
		// json.Marshal of a flat map cannot fail in practice. Defend
		// against future schema changes that introduce non-marshalable
		// fields.
		return nil, "", fmt.Errorf("build body from --retention-days: %w", err)
	}
	return body, "--retention-days", nil
}
