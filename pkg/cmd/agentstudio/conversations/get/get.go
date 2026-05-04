package get

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type GetOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	AgentStudioClient func() (*agentstudio.Client, error)

	AgentID        string
	ConversationID string

	PrintFlags *cmdutil.PrintFlags
}

func NewGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "get <agent_id> <conversation_id>",
		Args:  cobra.ExactArgs(2),
		Short: "Get a conversation",
		Annotations: map[string]string{
			"acls": "admin",
		},
		Example: heredoc.Doc(`
			$ algolia agentstudio conversations get a1b2 conv-1
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.ConversationID = args[1]
			if runF != nil {
				return runF(opts)
			}
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
	opts.IO.StartProgressIndicatorWithLabel("Fetching conversation")
	conv, err := client.GetConversation(opts.AgentID, opts.ConversationID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	return p.Print(opts.IO, conv)
}
