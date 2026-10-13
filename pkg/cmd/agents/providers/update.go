package providers

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// UpdateOptions collects inputs for `agents providers update`.
type UpdateOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	ProviderID    string
	File          string
	Name          string
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
		Use:   "update <provider-id> (-F <file> | [--name <name>] [--base-url <url>])",
		Short: "Patch an LLM provider authentication",
		Long: heredoc.Doc(`
			Patch a provider authentication from a JSON file (-F) or with
			--name / --base-url for non-secret fields only.

			The -F body is PATCH JSON (only fields present are updated).
			Rotate a vendor API key by including "input": {"apiKey": "..."}
			in that file — not via flags.

			Do not combine -F with --name or --base-url.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers update <id> -F rename.json
			$ algolia agents providers update <id> -F rotate-key.json
			$ algolia agents providers update <id> --name new-label
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
	cmd.Flags().StringVar(&opts.BaseURL, "base-url", "", `Set base URL inside "input" (shortcut; not with -F)`)
	cmd.Flags().
		BoolVar(&opts.Show, "show-secret", false, "Render secret fields verbatim in the success response")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	var body json.RawMessage
	var err error

	switch {
	case updateInlineFlagsConflictWithFile(opts.File, opts.Name, opts.BaseURL):
		return cmdutil.FlagErrorf("cannot combine -F/--file with --name or --base-url")
	case opts.File != "":
		body, err = shared.ReadJSONFile(opts.IO.In, opts.File)
		if err != nil {
			return err
		}
	default:
		raw, err := marshalSimpleProviderPatch(opts.Name, opts.BaseURL)
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

func marshalSimpleProviderPatch(name, baseURL string) ([]byte, error) {
	patch := map[string]any{}
	if strings.TrimSpace(name) != "" {
		patch["name"] = strings.TrimSpace(name)
	}
	if strings.TrimSpace(baseURL) != "" {
		patch["input"] = map[string]string{"baseUrl": strings.TrimSpace(baseURL)}
	}
	if len(patch) == 0 {
		return nil, cmdutil.FlagErrorf(
			"specify a JSON patch with -F (include input.apiKey to rotate secrets), or --name and/or --base-url",
		)
	}
	return json.Marshal(patch)
}

func updateUsesInlineFlags(name, baseURL string) bool {
	return strings.TrimSpace(name) != "" || strings.TrimSpace(baseURL) != ""
}

func updateInlineFlagsConflictWithFile(file, name, baseURL string) bool {
	if strings.TrimSpace(file) == "" {
		return false
	}
	return updateUsesInlineFlags(name, baseURL)
}
