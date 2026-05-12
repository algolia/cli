package providers

import (
	"context"
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type CreateOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	File          string
	Show          bool
	OutputChanged bool
}

func newCreateCmd(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "create -F <file>",
		Short: "Create an LLM provider authentication",
		Long: heredoc.Doc(`
			Create a provider authentication from a JSON file (-F). The body
			is ProviderAuthenticationCreate JSON (name, providerName, input).
			The "input" shape depends on providerName:

			  - openai / anthropic: {apiKey, baseUrl?}
			  - azure_openai:       {apiKey, azureEndpoint, azureDeployment, apiVersion?}
			  - openai_compatible:  {apiKey, baseUrl, defaultModel}
			  - google_genai / deepseek: {apiKey}

			Put vendor credentials (e.g. input.apiKey) in that file — not on
			the command line.

			By default the created provider in the success response is
			masked. Pass --show-secret to render the apiKey verbatim.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers create -F openai-prod.json
			$ cat spec.json | algolia agents providers create -F -
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			opts.OutputChanged = cmd.Flags().Changed("output")
			if runF != nil {
				return runF(opts)
			}
			return runCreateCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&opts.File, "file", "F", "", "JSON file with the provider body (use \"-\" for stdin)")
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().
		BoolVar(&opts.Show, "show-secret", false, "Render secret fields verbatim in the success response")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	body, err := shared.ReadJSONFile(opts.IO.In, opts.File)
	if err != nil {
		return err
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := shared.OrBackground(opts.Ctx)

	opts.IO.StartProgressIndicatorWithLabel("Creating provider")
	p, err := client.CreateProvider(ctx, json.RawMessage(body))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if !opts.Show {
		p.Input = shared.MaskInput(p.Input)
	}
	return opts.PrintFlags.Print(opts.IO, p)
}
