package get

import (
	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// GetOptions holds the dependencies and flags for the get command.
type GetOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	SecretKeyID       string
	PrintFlags        *cmdutil.PrintFlags
}

// NewGetCmd returns the `secret-keys get` command.
func NewGetCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "get <secret-key-id>",
		Short:             "Get a secret key by ID",
		Args:              validators.ExactArgsWithMsg(1, "secret-keys get requires a <secret-key-id> argument."),
		ValidArgsFunction: cmdutil.SecretKeyIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "settings",
		},
		Example: heredoc.Doc(`
			# Get the secret key with ID "my-key"
			$ algolia secret-keys get my-key
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.SecretKeyID = args[0]
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

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching secret key")

	res, err := client.GetSecretKey(client.NewApiGetSecretKeyRequest(opts.SecretKeyID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	return p.Print(opts.IO, res)
}
