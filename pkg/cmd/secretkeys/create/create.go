package create

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// CreateOptions holds the dependencies and flags for the create command.
type CreateOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)

	Name     string
	AgentIDs []string

	PrintFlags *cmdutil.PrintFlags
}

// NewCreateCmd returns the `secret-keys create` command.
func NewCreateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a secret key",
		Args:  validators.ExactArgsWithMsg(1, "secret-keys create requires a <name> argument."),
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			# Create a secret key
			$ algolia secret-keys create my-key

			# Create a secret key scoped to specific agents
			$ algolia secret-keys create my-key --agent-ids agent_1,agent_2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			return runCreateCmd(opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.AgentIDs, "agent-ids", nil, "Agent IDs this secret key is associated with")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	secretKeyCreate := agentStudio.NewSecretKeyCreate(opts.Name)
	if len(opts.AgentIDs) > 0 {
		secretKeyCreate.AgentIds = opts.AgentIDs
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Creating secret key")

	res, err := client.CreateSecretKey(client.NewApiCreateSecretKeyRequest(secretKeyCreate))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Created secret key %s\n", cs.SuccessIcon(), cs.Bold(res.Id))
	}

	return p.Print(opts.IO, res)
}
