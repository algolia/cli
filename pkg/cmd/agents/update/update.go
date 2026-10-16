package update

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// UpdateOptions holds the dependencies and flags for the update command.
type UpdateOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	PrintFlags        *cmdutil.PrintFlags

	AgentID      string
	JSON         string
	File         string
	Name         string
	Instructions string
}

// NewUpdateCmd returns the `agents update` command.
func NewUpdateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &UpdateOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "update <agent-id>",
		Short: "Update an Agent Studio agent",
		Args: validators.ExactArgsWithMsg(
			1,
			"agentstudio agents update requires a <agent-id> argument.",
		),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Update the agent name
			$ algolia agentstudio agents update 11111111-1111-4111-8111-111111111111 --name "renamed-bot"

			# Update from a JSON object
			$ algolia agentstudio agents update 11111111-1111-4111-8111-111111111111 --json '{"name":"renamed-bot","instructions":"Be brief."}'

			# Update from a JSON file
			$ algolia agentstudio agents update 11111111-1111-4111-8111-111111111111 --file patch.json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]

			hasPayload := opts.JSON != "" || opts.File != ""
			hasFlags := opts.Name != "" || opts.Instructions != ""
			if hasPayload && hasFlags {
				return cmdutil.FlagErrorf(
					"--json and --file cannot be combined with --name or --instructions",
				)
			}
			if !hasPayload && !hasFlags {
				return cmdutil.FlagErrorf(
					"`--json`, `--file`, or at least one of `--name` or `--instructions` is required",
				)
			}

			return runUpdateCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.JSON, "json", "", "Agent JSON patch object")
	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "JSON file path (use - for stdin)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Agent name")
	cmd.Flags().StringVar(&opts.Instructions, "instructions", "", "Agent prompt")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func buildPatch(opts *UpdateOptions) (*agentStudio.AgentConfigUpdate, error) {
	raw, err := cmdutil.JSONPayload(opts.JSON, opts.File, opts.IO.In)
	if err != nil {
		return nil, err
	}
	if raw == "" {
		body := agentStudio.NewEmptyAgentConfigUpdate()
		if opts.Name != "" {
			body.SetName(opts.Name)
		}
		if opts.Instructions != "" {
			body.SetInstructions(opts.Instructions)
		}
		return body, nil
	}

	var body agentStudio.AgentConfigUpdate
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		return nil, fmt.Errorf("parsing agent JSON: %w", err)
	}

	return &body, nil
}

func runUpdateCmd(opts *UpdateOptions) error {
	body, err := buildPatch(opts)
	if err != nil {
		return err
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Updating agent")

	res, err := client.UpdateAgent(client.NewApiUpdateAgentRequest(opts.AgentID, body))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Updated agent %s\n", cs.SuccessIcon(), opts.AgentID)
	}

	return p.Print(opts.IO, res)
}
