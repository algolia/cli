package providers

import (
	"context"
	"sort"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
)

type DefaultsOptions struct {
	IO                *iostreams.IOStreams
	Ctx               context.Context
	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags
}

func newDefaultsCmd(f *cmdutil.Factory, runF func(*DefaultsOptions) error) *cobra.Command {
	opts := &DefaultsOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags(),
	}
	cmd := &cobra.Command{
		Use:   "defaults",
		Short: "Recommended default model for each provider type",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Ctx = cmd.Context()
			if runF != nil {
				return runF(opts)
			}
			return runDefaultsCmd(opts)
		},
	}
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runDefaultsCmd(opts *DefaultsOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Fetching default models")
	res, err := client.GetProviderModelDefaults(ctxOrBackground(opts.Ctx))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if opts.PrintFlags.HasStructuredOutput() {
		return opts.PrintFlags.Print(opts.IO, res)
	}
	keys := make([]string, 0, len(res))
	for k := range res {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("PROVIDER TYPE", nil, nil)
		table.AddField("DEFAULT MODEL", nil, nil)
		table.EndRow()
	}
	for _, k := range keys {
		table.AddField(k, nil, nil)
		table.AddField(res[k], nil, nil)
		table.EndRow()
	}
	return table.Render()
}
