package update

import (
	"context"
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type UpdateOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioAPIClient func() (*agentStudio.APIClient, error)
	PrintFlags           *cmdutil.PrintFlags

	AgentID       string
	File          string
	OutputChanged bool
}

func NewUpdateCmd(f *cmdutil.Factory, runF func(*UpdateOptions) error) *cobra.Command {
	opts := &UpdateOptions{
		IO:                   f.IOStreams,
		AgentStudioAPIClient: f.AgentStudioAPIClient,
		PrintFlags:           cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "update <agent-id> -F <file>",
		Short: "Update an Agent Studio agent from a JSON patch file",
		Long: heredoc.Doc(`
			Update an existing agent. The file body is a partial
			AgentConfigUpdate — only the fields you want to change need to
			be present.
		`),
		Example: heredoc.Doc(`
			# Rename an agent
			$ echo '{"name":"New name"}' | algolia agents update <id> -F -

			# Apply a patch from a file
			$ algolia agents update <id> -F patch.json
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

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	body, err := shared.ReadJSONFile(opts.IO.In, opts.File)
	if err != nil {
		return err
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	client, err := opts.AgentStudioAPIClient()
	if err != nil {
		return err
	}

	// The update endpoint is a PATCH and the SDK exposes no pass-through PATCH,
	// so decode into the typed partial model. Every AgentConfigUpdate field is
	// optional, so the documented fields survive the round-trip; unrecognized
	// fields are dropped.
	patch := agentStudio.NewEmptyAgentConfigUpdate()
	if err := json.Unmarshal(body, patch); err != nil {
		return cmdutil.FlagErrorf("patch body must be a JSON object: %v", err)
	}

	opts.IO.StartProgressIndicatorWithLabel("Updating agent")
	agent, err := client.UpdateAgent(
		client.NewApiUpdateAgentRequest(opts.AgentID, patch),
		agentStudio.WithContext(ctx),
	)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	return opts.PrintFlags.Print(opts.IO, agent)
}
