package providers

import (
	"context"
	"encoding/json"
	"fmt"

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
	DryRun        bool
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
		Use:   "update <provider-id> -F <file>",
		Short: "Patch an LLM provider authentication from a JSON file",
		Long: heredoc.Doc(`
			Patch a provider authentication. PATCH semantics: only the
			fields in the file are updated. Pass {"name":"new-name"} to
			rename, or {"input":{"apiKey":"sk-NEW"}} to rotate the key.
		`),
		Example: heredoc.Doc(`
			$ algolia agents providers update <id> -F rename.json
			$ algolia agents providers update <id> -F rotate.json --dry-run
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
	_ = cmd.MarkFlagRequired("file")
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Validate and print the resolved request body without calling the API")
	cmd.Flags().
		BoolVar(&opts.Show, "show-secret", false, "Render secret fields verbatim in the success response")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	body, err := shared.ReadJSONFile(opts.IO.In, opts.File)
	if err != nil {
		return err
	}

	if opts.DryRun {
		return shared.PrintDryRun(opts.IO, opts.PrintFlags, opts.OutputChanged,
			"update_provider",
			fmt.Sprintf("PATCH /1/providers/%s", opts.ProviderID),
			opts.File, body, map[string]any{"providerId": opts.ProviderID})
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
