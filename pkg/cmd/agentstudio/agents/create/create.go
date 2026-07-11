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
	File              string
	PrintFlags        *cmdutil.PrintFlags
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
		Use:   "create -f <file>",
		Short: "Create a draft agent",
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Create an agent from a JSON file
			$ algolia agents create --file agent.json

			# Create an agent from stdin
			$ cat agent.json | algolia agents create --file -
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "JSON file path (use - for stdin)")
	_ = cmd.MarkFlagRequired("file")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	raw, err := cmdutil.ReadFile(opts.File, opts.IO.In)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var agentConfig agentStudio.AgentConfigCreate
	if err := json.Unmarshal(raw, &agentConfig); err != nil {
		return fmt.Errorf("parsing agent JSON: %w", err)
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

	res, err := client.CreateAgent(client.NewApiCreateAgentRequest(&agentConfig))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Created agent %s\n", cs.SuccessIcon(), cs.Bold(res.Id))
	}

	return p.Print(opts.IO, res)
}
