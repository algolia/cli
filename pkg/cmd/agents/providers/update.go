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

type UpdateOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	ProviderID    string
	File          string
	Name          string
	APIKey        string
	APIKeyStdin   bool
	APIKeyEnv     string
	BaseURL       string
	Show          bool
	OutputChanged bool
}

func newUpdateCmd(f *cmdutil.Factory, runF func(*UpdateOptions) error) *cobra.Command {
	opts := &UpdateOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use: "update <provider-id> (-F <file> | [--name <name>] " +
			"[--api-key <key> | --api-key-stdin | --api-key-env <var>] [--base-url <url>])",
		Short: "Patch an LLM provider authentication",
		Long: heredoc.Doc(`
			Patch a provider authentication. Either pass a JSON body with -F
			(PATCH semantics: only fields in the file are updated) or use
			flags for simple renames and key rotation.

			Do not combine -F with --name/--api-key* or --base-url.

			When using flags, at least one of --name, --api-key,
			--api-key-stdin, --api-key-env, or --base-url is required.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers update <id> -F rename.json
			$ algolia agents providers update <id> --name new-label
			$ algolia agents providers update <id> --api-key-env OPENAI_API_KEY
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProviderID = args[0]
			opts.Ctx = cmd.Context()
			opts.OutputChanged = cmd.Flags().Changed("output")
			if opts.ProviderID == "" {
				return cmdutil.FlagErrorf("provider-id must not be empty")
			}
			if runF != nil {
				return runF(opts)
			}
			return runUpdateCmd(opts)
		},
	}

	cmd.Flags().
		StringVarP(&opts.File, "file", "F", "", "JSON file with the provider patch body (use \"-\" for stdin)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Rename the provider label (shortcut; not with -F)")
	cmd.Flags().StringVar(&opts.APIKey, "api-key", "", "Rotate the API credential (shortcut; not with -F)")
	cmd.Flags().BoolVar(&opts.APIKeyStdin, "api-key-stdin", false, "Read new API key from stdin (shortcut; not with -F)")
	cmd.Flags().StringVar(&opts.APIKeyEnv, "api-key-env", "", "Read new API key from this environment variable (shortcut; not with -F)")
	cmd.Flags().StringVar(&opts.BaseURL, "base-url", "", `Set or clear base URL inside "input" (shortcut; not with -F)`)
	cmd.Flags().
		BoolVar(&opts.Show, "show-secret", false, "Render secret fields verbatim in the success response")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	var body json.RawMessage
	var err error

	switch {
	case updateInlineFlagsConflictWithFile(opts.File, opts.Name, opts.APIKey, opts.APIKeyEnv, opts.BaseURL, opts.APIKeyStdin):
		return cmdutil.FlagErrorf("cannot combine -F/--file with --name/--api-key/--api-key-* or --base-url")
	case opts.File != "":
		body, err = shared.ReadJSONFile(opts.IO.In, opts.File)
		if err != nil {
			return err
		}
	default:
		if !updateUsesInlineFlags(opts.Name, opts.APIKey, opts.APIKeyEnv, opts.BaseURL, opts.APIKeyStdin) {
			return cmdutil.FlagErrorf("specify a JSON patch with -F, or at least one shortcut flag (--name, --api-key, etc.)")
		}
		key, setKey, err := resolveOptionalAPIKey(opts.IO.In, opts.APIKey, opts.APIKeyEnv, opts.APIKeyStdin)
		if err != nil {
			return err
		}
		raw, err := marshalSimpleProviderPatch(opts.Name, opts.BaseURL, key, setKey)
		if err != nil {
			return err
		}
		body = raw
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := shared.OrBackground(opts.Ctx)

	opts.IO.StartProgressIndicatorWithLabel("Updating provider")
	p, err := client.UpdateProvider(ctx, opts.ProviderID, json.RawMessage(body))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if !opts.Show {
		p.Input = shared.MaskInput(p.Input)
	}
	return opts.PrintFlags.Print(opts.IO, p)
}
