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
	DryRun        bool
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
		Short: "Create an LLM provider authentication from a JSON file",
		Long: heredoc.Doc(`
			Create a provider authentication from a JSON file describing
			the ProviderAuthenticationCreate body (name, providerName,
			input). The "input" subobject's shape varies per providerName:

			  - openai / anthropic: {apiKey, baseUrl?}
			  - azure_openai:       {apiKey, azureEndpoint, azureDeployment, apiVersion?}
			  - openai_compatible:  {apiKey, baseUrl, defaultModel}
			  - google_genai / deepseek: {apiKey}

			The file is sent verbatim; field-level validation is the
			backend's job (a 4xx surfaces with the structured detail).

			Use --dry-run to preview the request without sending.

			By default the created provider in the success response is
			masked. Pass --show-secret to render the apiKey verbatim.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers create -F openai-prod.json
			$ cat spec.json | algolia agents providers create -F -
			$ algolia agents providers create -F spec.json --dry-run
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
		BoolVar(&opts.DryRun, "dry-run", false, "Validate and print the resolved request body without calling the API")
	cmd.Flags().
		BoolVar(&opts.Show, "show-secret", false, "Render secret fields verbatim in the success response")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	body, err := readBody(opts.File, opts.IO)
	if err != nil {
		return err
	}

	if opts.DryRun {
		return shared.PrintDryRun(opts.IO, opts.PrintFlags, opts.OutputChanged,
			"create_provider", "POST /1/providers", opts.File, body, nil)
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := ctxOrBackground(opts.Ctx)

	opts.IO.StartProgressIndicatorWithLabel("Creating provider")
	p, err := client.CreateProvider(ctx, json.RawMessage(body))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if !opts.Show {
		p.Input = MaskInput(p.Input)
	}
	return opts.PrintFlags.Print(opts.IO, p)
}
