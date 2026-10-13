package keys

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
)

type CreateOptions struct {
	IO                   *iostreams.IOStreams
	Ctx                  context.Context
	AgentStudioAPIClient func() (*agentStudio.APIClient, error)
	PrintFlags           *cmdutil.PrintFlags
	Name                 string
	AgentIDs             []string
	ShowSecret           bool
}

func newCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:                   f.IOStreams,
		AgentStudioAPIClient: f.AgentStudioAPIClient,
		PrintFlags:           cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "create --name <name> [--agent-id <id> ...]",
		Short: "Create a secret key (admin key required)",
		Example: heredoc.Doc(`
			$ algolia agents keys create --name web-widget
			$ algolia agents keys create --name shared --agent-id a1 --agent-id a2
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			if opts.Name == "" {
				return cmdutil.FlagErrorf("--name is required")
			}
			if runF != nil {
				return runF(opts)
			}
			return runCreateCmd(opts)
		},
	}
	cmd.Flags().StringVar(&opts.Name, "name", "", "Key name (required, max 128)")
	cmd.Flags().StringSliceVar(&opts.AgentIDs, "agent-id", nil, "Restrict the key to specific agents (repeatable)")
	cmd.Flags().
		BoolVar(&opts.ShowSecret, "show-secret", false, "Reveal raw key value in the response (default redacted as ***)")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	body := &agentStudio.SecretKeyCreate{Name: opts.Name, AgentIds: opts.AgentIDs}
	client, err := opts.AgentStudioAPIClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Creating secret key")
	k, err := client.CreateSecretKey(
		client.NewApiCreateSecretKeyRequest(body),
		agentStudio.WithContext(shared.OrBackground(opts.Ctx)),
	)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	masked := maskKey(*k, opts.ShowSecret)
	return opts.PrintFlags.Print(opts.IO, &masked)
}
