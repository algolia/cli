package models

import (
	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// ModelsOptions holds the dependencies and flags for the models command.
type ModelsOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)

	ProviderID string

	PrintFlags *cmdutil.PrintFlags
}

// NewModelsCmd returns the `providers models` command.
func NewModelsCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ModelsOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "models",
		Short: "List the models available for providers",
		Args:  validators.NoArgs(),
		Annotations: map[string]string{
			"acls": "settings",
		},
		Example: heredoc.Doc(`
			# List models for every supported provider name
			$ algolia providers models

			# List models for a specific configured provider
			$ algolia providers models --provider-id my-provider
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModelsCmd(opts)
		},
	}

	cmd.Flags().
		StringVar(&opts.ProviderID, "provider-id", "", "List the models available for this configured provider")
	_ = cmd.RegisterFlagCompletionFunc("provider-id", cmdutil.ProviderIDs(opts.AgentStudioClient))

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runModelsCmd(opts *ModelsOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching models")

	if opts.ProviderID != "" {
		res, err := client.ListProviderModels(client.NewApiListProviderModelsRequest(opts.ProviderID))
		opts.IO.StopProgressIndicator()
		if err != nil {
			return err
		}
		return p.Print(opts.IO, res)
	}

	res, err := client.ListModels()
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	return p.Print(opts.IO, res)
}
