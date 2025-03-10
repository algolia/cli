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

// NewListCmd creates and returns a list command for GenAI responses.
func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		IO:          f.IOStreams,
		Config:      f.Config,
		GenAIClient: f.GenAIClient,
		PrintFlags:  cmdutil.NewPrintFlags(),
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List GenAI responses",
		Long: heredoc.Doc(`
			List GenAI responses.

			Note: This feature might not be supported by the Algolia GenAI API yet.
		`),
		Example: heredoc.Doc(`
			# List all responses
			$ algolia genai response list
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
	opts.IO.StartProgressIndicatorWithLabel("Fetching responses")
	response, err := client.ListResponses()
	opts.IO.StopProgressIndicator()

	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.FailureIcon(), err.Error())
		fmt.Fprintln(opts.IO.ErrOut, "")
		fmt.Fprintln(opts.IO.ErrOut, "To get a specific response, you can use:")
		fmt.Fprintln(opts.IO.ErrOut, "  $ algolia genai response get <id>")
		fmt.Fprintln(opts.IO.ErrOut, "")
		fmt.Fprintln(opts.IO.ErrOut, "Alternatively, you can access the Algolia dashboard to view your responses.")
		return fmt.Errorf("error fetching responses")
	}

	if opts.PrintFlags.OutputFlagSpecified() {
		printer, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}

		return printer.Print(opts.IO, response)
	}

	if len(response.Responses) == 0 {
		fmt.Fprintln(opts.IO.Out, "No responses found")
		return nil
	}

	table := printers.NewTablePrinter(opts.IO)
	table.AddField("ID", nil, nil)
	table.AddField("QUERY", nil, nil)
	table.AddField("DATA SOURCE", nil, nil)
	table.AddField("PROMPT", nil, nil)
	table.AddField("CREATED", nil, nil)
	table.EndRow()

	for _, r := range response.Responses {
		createdTime := FormatTime(r.CreatedAt)

		table.AddField(r.ObjectID, nil, nil)
		table.AddField(truncateString(r.Query, 30), nil, nil)
		table.AddField(r.DataSourceID, nil, nil)
		table.AddField(r.PromptID, nil, nil)
		table.AddField(createdTime, nil, nil)
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

// truncateString truncates a string if it's longer than maxLen
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
