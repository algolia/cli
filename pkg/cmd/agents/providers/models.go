package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/validators"
)

type ModelsOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	ProviderID string
}

func newModelsCmd(f *cmdutil.Factory, runF func(*ModelsOptions) error) *cobra.Command {
	opts := &ModelsOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags(),
	}

	cmd := &cobra.Command{
		Use:   "models [--provider-id <id>]",
		Short: "List models available per provider type, or for a configured provider",
		Long: heredoc.Doc(`
			Without --provider-id: GET /1/providers/models — the static
			catalog of supported models per provider type. Use before
			creating a provider.

			With --provider-id <id>: GET /1/providers/<id>/models — the
			models a specific configured provider exposes, including
			account-specific entries (OpenAI fine-tunes, Azure
			deployments).
		`),
		Example: heredoc.Doc(`
			# What can I configure?
			$ algolia agents providers models

			# What does my OpenAI provider actually expose (incl. fine-tunes)?
			$ algolia agents providers models --provider-id <id>
		`),
		Args: validators.NoArgs(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			// Refuse explicit empty --provider-id (almost certainly a scripting bug).
			if cmd.Flags().Changed("provider-id") && opts.ProviderID == "" {
				return cmdutil.FlagErrorf("--provider-id must not be empty")
			}
			if runF != nil {
				return runF(opts)
			}
			return runModelsCmd(opts)
		},
	}

	cmd.Flags().
		StringVar(&opts.ProviderID, "provider-id", "", "Use the configured-provider route instead of the static catalog")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runModelsCmd(opts *ModelsOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := ctxOrBackground(opts.Ctx)

	if opts.ProviderID == "" {
		return runListCatalog(opts, client, ctx)
	}
	return runListForProvider(opts, client, ctx)
}

func runListCatalog(opts *ModelsOptions, client *agentstudio.Client, ctx context.Context) error {
	opts.IO.StartProgressIndicatorWithLabel("Fetching provider model catalog")
	models, err := client.ListProviderModels(ctx)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.PrintFlags.HasStructuredOutput() {
		return opts.PrintFlags.Print(opts.IO, models)
	}

	providerNames := make([]string, 0, len(models))
	for k := range models {
		providerNames = append(providerNames, k)
	}
	sort.Strings(providerNames)

	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("PROVIDER", nil, nil)
		table.AddField("MODEL", nil, nil)
		table.EndRow()
	}
	for _, prov := range providerNames {
		ms := models[prov]
		sort.Strings(ms)
		for _, m := range ms {
			table.AddField(prov, nil, nil)
			table.AddField(m, nil, nil)
			table.EndRow()
		}
	}
	return table.Render()
}

func runListForProvider(opts *ModelsOptions, client *agentstudio.Client, ctx context.Context) error {
	opts.IO.StartProgressIndicatorWithLabel("Fetching models for provider")
	raw, err := client.ListModelsForProvider(ctx, opts.ProviderID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.PrintFlags.HasStructuredOutput() {
		// Round-trip through a generic decode so --output json re-formats
		// consistently instead of emitting the backend's whitespace.
		var anyV any
		if err := json.Unmarshal(raw, &anyV); err != nil {
			return fmt.Errorf("decode models response: %w", err)
		}
		return opts.PrintFlags.Print(opts.IO, anyV)
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, raw, "", "  "); err != nil {
		_, _ = opts.IO.Out.Write(raw)
		return nil
	}
	_, _ = opts.IO.Out.Write(pretty.Bytes())
	_, _ = opts.IO.Out.Write([]byte("\n"))
	return nil
}
