package models

import (
	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

// ModelsOptions holds the dependencies and flags for the models command.
type ModelsOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	PrintFlags        *cmdutil.PrintFlags
	ProviderID        string
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
		Use:   "models [<provider-id>]",
		Short: "List models available per provider type, or for a configured provider",
		Long:  "Without an ID, lists models by provider type. With a provider ID, lists models for that configured provider.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return cmdutil.FlagErrorf(
					"agentstudio providers models takes an optional <provider-id> argument.",
				)
			}
			return nil
		},
		Annotations: map[string]string{
			"acls": "settings",
		},
		Example: heredoc.Doc(`
			# List models by provider type
			$ algolia agentstudio providers models

			# List models for a configured provider
			$ algolia agentstudio providers models 22222222-2222-4222-8222-222222222222
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				opts.ProviderID = args[0]
			}
			return runModelsCmd(opts)
		},
	}

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

	if opts.ProviderID != "" {
		opts.IO.StartProgressIndicatorWithLabel("Fetching provider models")
		res, err := client.ListProviderModels(
			client.NewApiListProviderModelsRequest(opts.ProviderID),
		)
		if err != nil {
			opts.IO.StopProgressIndicator()
			return err
		}
		opts.IO.StopProgressIndicator()
		return p.Print(opts.IO, res)
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching models")
	res, err := client.ListModels()
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}
	opts.IO.StopProgressIndicator()
	return p.Print(opts.IO, res)
}
