package list

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/genai"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
)

type ListOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	GenAIClient func() (*genai.Client, error)

	PrintFlags *cmdutil.PrintFlags
}

// NewListCmd creates and returns a list command for GenAI prompts.
func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		IO:          f.IOStreams,
		Config:      f.Config,
		GenAIClient: f.GenAIClient,
		PrintFlags:  cmdutil.NewPrintFlags(),
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List GenAI prompts",
		Long: heredoc.Doc(`
			List GenAI prompts.

			Note: This feature is not supported by the Algolia GenAI API yet.
		`),
		Example: heredoc.Doc(`
			# List all prompts
			$ algolia genai prompt list
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.GenAIClient()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	opts.IO.StartProgressIndicatorWithLabel("Fetching prompts")
	response, err := client.ListPrompts()
	opts.IO.StopProgressIndicator()

	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
		fmt.Fprintln(opts.IO.ErrOut, "")
		fmt.Fprintln(opts.IO.ErrOut, "To get a specific prompt, you can use:")
		fmt.Fprintln(opts.IO.ErrOut, "  $ algolia genai prompt get <id>")
		fmt.Fprintln(opts.IO.ErrOut, "")
		fmt.Fprintln(opts.IO.ErrOut, "Alternatively, you can access the Algolia dashboard to view your prompts.")
		return fmt.Errorf("error fetching prompts")
	}

	if opts.PrintFlags.OutputFlagSpecified() {
		printer, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}

		return printer.Print(opts.IO, response)
	}

	if len(response.Prompts) == 0 {
		fmt.Fprintln(opts.IO.Out, "No prompts found")
		return nil
	}

	table := printers.NewTablePrinter(opts.IO)
	table.AddField("ID", nil, nil)
	table.AddField("NAME", nil, nil)
	table.AddField("CREATED", nil, nil)
	table.AddField("UPDATED", nil, nil)
	table.EndRow()

	for _, p := range response.Prompts {
		createdTime := FormatTime(p.CreatedAt)
		updatedTime := FormatTime(p.UpdatedAt)

		table.AddField(p.ID, nil, nil)
		table.AddField(p.Name, nil, nil)
		table.AddField(createdTime, nil, nil)
		table.AddField(updatedTime, nil, nil)
		table.EndRow()
	}

	return table.Render()
}

// FormatTime formats time to a readable format
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
