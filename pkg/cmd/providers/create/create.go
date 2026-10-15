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
	PrintFlags        *cmdutil.PrintFlags
	JSON              string
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
		Use:   "create",
		Short: "Create an Agent Studio LLM provider",
		Long:  "Creates a provider from a JSON object. Anthropic keys must be scoped to a single workspace; Agent Studio does not send anthropic-workspace-id.",
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Create an Anthropic provider
			$ algolia agentstudio providers create --json '{"name":"anthropic-prod","providerName":"anthropic","input":{"apiKey":"<api-key>"}}'

			# Create an OpenAI provider
			$ algolia agentstudio providers create --json '{"name":"openai-prod","providerName":"openai","input":{"apiKey":"<api-key>"}}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.JSON == "" {
				return cmdutil.FlagErrorf("`--json` is required")
			}
			return runCreateCmd(opts)
		},
	}

	cmd.Flags().
		StringVar(&opts.JSON, "json", "", "Provider JSON object (name, providerName, input)")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func buildProvider(opts *CreateOptions) (*agentStudio.ProviderAuthenticationCreate, error) {
	var body agentStudio.ProviderAuthenticationCreate
	if err := json.Unmarshal([]byte(opts.JSON), &body); err != nil {
		return nil, fmt.Errorf("parsing provider JSON: %w", err)
	}
	if body.Name == "" || !body.ProviderName.IsValid() {
		return nil, cmdutil.FlagErrorf("provider JSON requires `name` and a valid `providerName`")
	}
	if body.Input.GetActualInstance() == nil {
		return nil, cmdutil.FlagErrorf("provider JSON requires `input`")
	}

	return &body, nil
}

func runCreateCmd(opts *CreateOptions) error {
	body, err := buildProvider(opts)
	if err != nil {
		return err
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

	res, err := client.CreateProvider(client.NewApiCreateProviderRequest(body))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Created provider %s\n", cs.SuccessIcon(), res.Id)
	}

	return p.Print(opts.IO, res)
}
