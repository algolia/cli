package unpublish

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

// UnpublishOptions holds the dependencies and flags for the unpublish command.
type UnpublishOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	AgentID           string
	PrintFlags        *cmdutil.PrintFlags
}

// NewUnpublishCmd returns the `agents unpublish` command.
func NewUnpublishCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &UnpublishOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "unpublish <agent-id>",
		Short:             "Unpublish an agent from production",
		Args:              validators.ExactArgsWithMsg(1, "agents unpublish requires an <agent-id> argument."),
		ValidArgsFunction: cmdutil.AgentIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Unpublish the agent with ID "my-agent"
			$ algolia agents unpublish my-agent
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			return runUnpublishCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runUnpublishCmd(opts *UnpublishOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Unpublishing agent")

	res, err := client.UnpublishAgent(client.NewApiUnpublishAgentRequest(opts.AgentID))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Unpublished agent %s\n", cs.SuccessIcon(), cs.Bold(opts.AgentID))
	}

	return p.Print(opts.IO, res)
}
