package update

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

type UpdateOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	AgentID       string
	File          string
	DryRun        bool
	OutputChanged bool
}

func NewUpdateCmd(f *cmdutil.Factory, runF func(*UpdateOptions) error) *cobra.Command {
	opts := &UpdateOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "update <agent-id> -F <file>",
		Short: "Update an Agent Studio agent from a JSON patch file",
		Long: heredoc.Doc(`
			Update an existing agent. The file body is a partial
			AgentConfigUpdate — only the fields you want to change need to
			be present.

			Use --dry-run to print the resolved patch body without sending
			it — useful for verifying generated patches in CI.
		`),
		Example: heredoc.Doc(`
			# Rename an agent
			$ echo '{"name":"New name"}' | algolia agents update <id> -F -

			# Apply a patch from a file
			$ algolia agents update <id> -F patch.json

			# Preview without sending
			$ algolia agents update <id> -F patch.json --dry-run
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			opts.OutputChanged = cmd.Flags().Changed("output")
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			if runF != nil {
				return runF(opts)
			}
			return runUpdateCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&opts.File, "file", "F", "", "JSON file with the patch body (use \"-\" for stdin)")
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Validate and print the resolved request body without calling the API")

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	body, err := cmdutil.ReadFile(opts.File, opts.IO.In)
	if err != nil {
		return fmt.Errorf("failed to read patch body from %s: %w", shared.SourceLabel(opts.File), err)
	}
	body = shared.TrimUTF8BOM(body)

	if !json.Valid(body) {
		return cmdutil.FlagErrorf("patch body in %s is not valid JSON", shared.SourceLabel(opts.File))
	}

	if opts.DryRun {
		return shared.PrintDryRun(opts.IO, opts.PrintFlags, opts.OutputChanged,
			"update_agent", "PATCH /1/agents/"+opts.AgentID, opts.File, body,
			map[string]any{"agentId": opts.AgentID})
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Updating agent")
	agent, err := client.UpdateAgent(ctx, opts.AgentID, json.RawMessage(body))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	return opts.PrintFlags.Print(opts.IO, agent)
}
