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
	Name          string
	Provider      string
	APIKey        string
	APIKeyStdin   bool
	APIKeyEnv     string
	BaseURL       string
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
		Use: "create (-F <file> | --name <name> --provider <type> " +
			"(--api-key <key> | --api-key-stdin | --api-key-env <var>))",
		Short: "Create an LLM provider authentication",
		Long: heredoc.Doc(`
			Create a provider authentication from a JSON file (-F) or from
			flags for the common case (OpenAI, Anthropic, Google GenAI, or
			DeepSeek with a single API key).

			The -F body is the ProviderAuthenticationCreate JSON (name,
			providerName, input). The "input" subobject's shape varies per
			providerName:

			  - openai / anthropic: {apiKey, baseUrl?}
			  - azure_openai:       {apiKey, azureEndpoint, azureDeployment, apiVersion?}
			  - openai_compatible:  {apiKey, baseUrl, defaultModel}
			  - google_genai / deepseek: {apiKey}

			With flags, only openai, anthropic, google_genai, and deepseek
			are supported; use -F for Azure or openai_compatible.

			Do not combine -F with --name/--provider/--api-key*.

			Prefer --api-key-stdin or --api-key-env over --api-key (shell
			history may record raw flags).

			Use --dry-run to preview the request without sending.

			By default the created provider in the success response is
			masked. Pass --show-secret to render the apiKey verbatim.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers create -F openai-prod.json
			$ algolia agents providers create --name openai-prod --provider openai --api-key-env OPENAI_API_KEY
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
	cmd.Flags().StringVar(&opts.Name, "name", "", "Provider label (shortcut; not with -F)")
	cmd.Flags().StringVar(&opts.Provider, "provider", "", `Provider backend: openai, anthropic, google_genai, or deepseek (shortcut; not with -F)`)
	cmd.Flags().StringVar(&opts.APIKey, "api-key", "", "API credential (shortcut; not with -F)")
	cmd.Flags().BoolVar(&opts.APIKeyStdin, "api-key-stdin", false, "Read API key from stdin (shortcut; not with -F)")
	cmd.Flags().StringVar(&opts.APIKeyEnv, "api-key-env", "", "Read API key from this environment variable (shortcut; not with -F)")
	cmd.Flags().StringVar(&opts.BaseURL, "base-url", "", `Optional OpenAI / Anthropic-compatible base URL (shortcut; not with -F)`)
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Validate and print the resolved request body without calling the API")
	cmd.Flags().
		BoolVar(&opts.Show, "show-secret", false, "Render secret fields verbatim in the success response")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runCreateCmd(opts *CreateOptions) error {
	source := opts.File

	var body json.RawMessage
	var err error

	switch {
	case inlineFlagsConflictWithFile(opts.File, opts.Name, opts.Provider, opts.APIKey, opts.APIKeyEnv, opts.BaseURL, opts.APIKeyStdin):
		return cmdutil.FlagErrorf("cannot combine -F/--file with --name/--provider/--api-key/--api-key-* flags")
	case opts.File != "":
		body, err = shared.ReadJSONFile(opts.IO.In, opts.File)
		if err != nil {
			return err
		}
	default:
		if !createShortcutAttempted(opts.Name, opts.Provider, opts.APIKey, opts.APIKeyEnv, opts.BaseURL, opts.APIKeyStdin) {
			return cmdutil.FlagErrorf(`specify a JSON body with -F, or pass --name, --provider, and an API key flag`)
		}
		if opts.Name == "" || opts.Provider == "" {
			return cmdutil.FlagErrorf("when creating without -F, both --name and --provider are required")
		}
		key, err := resolveProvidedAPIKey(opts.IO.In, opts.APIKey, opts.APIKeyEnv, opts.APIKeyStdin)
		if err != nil {
			return err
		}
		raw, err := marshalSimpleProviderCreate(opts.Name, opts.Provider, opts.BaseURL, key)
		if err != nil {
			return err
		}
		body = raw
		source = "(flags)"
	}

	if opts.DryRun {
		return shared.PrintDryRun(opts.IO, opts.PrintFlags, opts.OutputChanged,
			"create_provider", "POST /1/providers", source, body, nil)
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
