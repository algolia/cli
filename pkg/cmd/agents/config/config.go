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
			Calls GET /1/configuration. Requires the ` + "`logs`" + ` ACL
			(governs log/conversation retention).
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

type SetOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	// RetentionDays: 0 is a valid backend value, so intent is detected
	// via cmd.Flags().Changed, not the value itself.
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
		body, err := shared.ReadJSONFile(opts.IO.In, opts.File)
		if err != nil {
			return nil, "", err
		}
		return body, opts.File, nil
	}
	body, err := json.Marshal(map[string]any{"maxRetentionDays": opts.RetentionDays})
	if err != nil {
		return nil, "", fmt.Errorf("build body from --retention-days: %w", err)
	}
	return body, "--retention-days", nil
}
