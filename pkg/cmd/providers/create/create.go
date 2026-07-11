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

// NewCreateCmd returns the `providers create` command.
func NewCreateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &CreateOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "create -f <file>",
		Short: "Create an LLM provider authentication",
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Create a provider from a JSON file
			$ algolia providers create --file provider.json

			# Create a provider from stdin
			$ cat provider.json | algolia providers create --file -
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

	var providerCreate agentStudio.ProviderAuthenticationCreate
	if err := json.Unmarshal(raw, &providerCreate); err != nil {
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

	opts.IO.StartProgressIndicatorWithLabel("Creating provider")

	res, err := client.CreateProvider(client.NewApiCreateProviderRequest(&providerCreate))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Created provider %s\n", cs.SuccessIcon(), cs.Bold(res.Id))
	}

	return p.Print(opts.IO, res)
}
