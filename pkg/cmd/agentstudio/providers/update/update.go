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
	ProviderID        string
	File              string
	PrintFlags        *cmdutil.PrintFlags
}

// NewUpdateCmd returns the `providers update` command.
func NewUpdateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &UpdateOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:               "update <provider-id> -f <file>",
		Short:             "Update an LLM provider authentication",
		Args:              validators.ExactArgsWithMsg(1, "providers update requires a <provider-id> argument."),
		ValidArgsFunction: cmdutil.ProviderIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Update a provider from a JSON file
			$ algolia providers update my-provider --file patch.json

			# Update a provider from stdin
			$ cat patch.json | algolia providers update my-provider --file -
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProviderID = args[0]
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

	var providerPatch agentStudio.ProviderAuthenticationPatch
	if err := json.Unmarshal(raw, &providerPatch); err != nil {
		return fmt.Errorf("parsing provider JSON: %w", err)
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Updating provider")

	res, err := client.UpdateProvider(client.NewApiUpdateProviderRequest(opts.ProviderID, &providerPatch))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Updated provider %s\n", cs.SuccessIcon(), cs.Bold(opts.ProviderID))
	}

	return p.Print(opts.IO, res)
}
