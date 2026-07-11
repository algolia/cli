package update

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/utils"
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

	SecretKeyID string
	Name        string
	AgentIDs    []string

	nameChanged     bool
	agentIDsChanged bool

	PrintFlags *cmdutil.PrintFlags
}

// NewUpdateCmd returns the `secret-keys update` command.
func NewUpdateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &UpdateOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "update <secret-key-id>",
		Short:             "Update a secret key",
		Args:              validators.ExactArgsWithMsg(1, "secret-keys update requires a <secret-key-id> argument."),
		ValidArgsFunction: cmdutil.SecretKeyIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			# Rename a secret key
			$ algolia secret-keys update my-key --name new-name

			# Change the agents a secret key is associated with
			$ algolia secret-keys update my-key --agent-ids agent_1,agent_2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.SecretKeyID = args[0]
			opts.nameChanged = cmd.Flags().Changed("name")
			opts.agentIDsChanged = cmd.Flags().Changed("agent-ids")

			if !opts.nameChanged && !opts.agentIDsChanged {
				return cmdutil.FlagErrorf("at least one of `--name` or `--agent-ids` is required")
			}

			return runUpdateCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "New name of the secret key")
	cmd.Flags().StringSliceVar(&opts.AgentIDs, "agent-ids", nil, "Updated agent IDs this secret key is associated with")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	secretKeyPatch := agentStudio.NewSecretKeyPatch()
	if opts.nameChanged {
		secretKeyPatch.Name = *utils.NewNullable(&opts.Name)
	}
	if opts.agentIDsChanged {
		secretKeyPatch.AgentIds = opts.AgentIDs
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Updating secret key")

	res, err := client.UpdateSecretKey(client.NewApiUpdateSecretKeyRequest(opts.SecretKeyID, secretKeyPatch))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Updated secret key %s\n", cs.SuccessIcon(), cs.Bold(opts.SecretKeyID))
	}

	return p.Print(opts.IO, res)
}
