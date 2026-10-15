package create

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

// CreateOptions holds the dependencies and flags for the create command.
type CreateOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	PrintFlags        *cmdutil.PrintFlags

	JSON         string
	Name         string
	Instructions string
}

// NewCreateCmd returns the `agents create` command.
func NewCreateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an Agent Studio agent",
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Create from flags
			$ algolia agentstudio agents create --name "support-bot" --instructions "You are a helpful assistant."

			# Create from a JSON object
			$ algolia agentstudio agents create --json '{"name":"support-bot","instructions":"You are a helpful assistant."}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.JSON != "" {
				if opts.Name != "" || opts.Instructions != "" {
					return cmdutil.FlagErrorf("--json cannot be combined with --name or --instructions")
				}
			} else if opts.Name == "" || opts.Instructions == "" {
				return cmdutil.FlagErrorf("`--json` or both `--name` and `--instructions` are required")
			}

			return runCreateCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.JSON, "json", "", "Agent JSON object")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Agent name")
	cmd.Flags().StringVar(&opts.Instructions, "instructions", "", "Agent prompt")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func buildAgent(opts *CreateOptions) (*agentStudio.AgentConfigCreate, error) {
	if opts.JSON == "" {
		return agentStudio.NewEmptyAgentConfigCreate().
			SetName(opts.Name).
			SetInstructions(opts.Instructions), nil
	}

	var body agentStudio.AgentConfigCreate
	if err := json.Unmarshal([]byte(opts.JSON), &body); err != nil {
		return nil, fmt.Errorf("parsing agent JSON: %w", err)
	}
	if body.Name == "" || body.Instructions == "" {
		return nil, cmdutil.FlagErrorf("agent JSON requires `name` and `instructions`")
	}

	return &body, nil
}

func runCreateCmd(opts *CreateOptions) error {
	body, err := buildAgent(opts)
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

	opts.IO.StartProgressIndicatorWithLabel("Creating agent")

	res, err := client.CreateAgent(client.NewApiCreateAgentRequest(body))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Created agent %s\n", cs.SuccessIcon(), res.Id)
	}

	return p.Print(opts.IO, res)
}
