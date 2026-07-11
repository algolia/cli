package update

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
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
	AgentID           string
	File              string
	PrintFlags        *cmdutil.PrintFlags
}

// NewUpdateCmd returns the `agents update` command.
func NewUpdateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &UpdateOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "update <agent-id> -f <file>",
		Short:             "Update an agent",
		Args:              validators.ExactArgsWithMsg(1, "agents update requires an <agent-id> argument."),
		ValidArgsFunction: cmdutil.AgentIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Update an agent from a JSON file
			$ algolia agents update my-agent --file patch.json

			# Update an agent from stdin
			$ cat patch.json | algolia agents update my-agent --file -
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			return runUpdateCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "JSON file path (use - for stdin)")
	_ = cmd.MarkFlagRequired("file")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	raw, err := cmdutil.ReadFile(opts.File, opts.IO.In)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var agentConfig agentStudio.AgentConfigUpdate
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

	opts.IO.StartProgressIndicatorWithLabel("Updating agent")

	res, err := client.UpdateAgent(client.NewApiUpdateAgentRequest(opts.AgentID, &agentConfig))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Updated agent %s\n", cs.SuccessIcon(), cs.Bold(opts.AgentID))
	}

	return p.Print(opts.IO, res)
}
