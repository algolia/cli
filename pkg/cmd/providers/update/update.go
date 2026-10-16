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
	PrintFlags        *cmdutil.PrintFlags

	ProviderID string
	JSON       string
	Name       string
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
		Use:   "update <provider-id>",
		Short: "Update an Agent Studio LLM provider",
		Long:  "Patches name and/or input. providerName cannot be changed.",
		Args: validators.ExactArgsWithMsg(
			1,
			"agentstudio providers update requires a <provider-id> argument.",
		),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Rename a provider
			$ algolia agentstudio providers update 22222222-2222-4222-8222-222222222222 --name "anthropic-prod"

			# Rotate the API key
			$ algolia agentstudio providers update 22222222-2222-4222-8222-222222222222 --json '{"input":{"apiKey":"<api-key>"}}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProviderID = args[0]

			if opts.JSON != "" {
				if opts.Name != "" {
					return cmdutil.FlagErrorf("--json cannot be combined with --name")
				}
			} else if opts.Name == "" {
				return cmdutil.FlagErrorf("`--json` or `--name` is required")
			}

			return runUpdateCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.JSON, "json", "", "Provider JSON patch object (name and/or input)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Provider display name")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func buildPatch(opts *UpdateOptions) (*agentStudio.ProviderAuthenticationPatch, error) {
	if opts.JSON == "" {
		return agentStudio.NewEmptyProviderAuthenticationPatch().SetName(opts.Name), nil
	}

	var body agentStudio.ProviderAuthenticationPatch
	if err := json.Unmarshal([]byte(opts.JSON), &body); err != nil {
		return nil, fmt.Errorf("parsing provider JSON: %w", err)
	}
	if !body.HasName() && !body.HasInput() {
		return nil, cmdutil.FlagErrorf("provider JSON requires `name` and/or `input`")
	}

	return &body, nil
}

func runUpdateCmd(opts *UpdateOptions) error {
	body, err := buildPatch(opts)
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

	opts.IO.StartProgressIndicatorWithLabel("Updating provider")

	res, err := client.UpdateProvider(client.NewApiUpdateProviderRequest(opts.ProviderID, body))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Updated provider %s\n", cs.SuccessIcon(), opts.ProviderID)
	}

	return p.Print(opts.IO, res)
}
