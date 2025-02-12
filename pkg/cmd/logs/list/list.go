package list

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type LogOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	PrintFlags *cmdutil.PrintFlags

	Entries   int32
	Start     int32
	LogType   string
	IndexName *string
}

// NewListCmd returns a new command for retrieving logs
func NewListCmd(f *cmdutil.Factory, runF func(*LogOptions) error) *cobra.Command {
	opts := &LogOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List log entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return runLogsCmd(opts)
		},
		Annotations: map[string]string{
			"acls": "logs",
		},
		Example: heredoc.Doc(`
      # Show the latest 5 Search API log entries
      $ algolia logs

      # Show the log entries 11 to 20
      $ algolia logs --entries 10 --start 11

      # Only show log entries with errors
      $ algolia logs --type error
    `),
	}

	opts.PrintFlags.AddFlags(cmd)

	cmd.Flags().Int32VarP(&opts.Entries, "entries", "e", 5, "How many log entries to show")
	cmd.Flags().
		Int32VarP(&opts.Start, "start", "s", 1, "Number of the first log entry to retrieve (starts with 1)")
	cmdutil.StringEnumFlag(
		cmd,
		&opts.LogType,
		"type",
		"t",
		"all",
		[]string{"all", "build", "query", "error"},
		"Type of log entries",
	)

	cmdutil.NilStringFlag(cmd, &opts.IndexName, "index", "i", "Filter logs by index name")

	return cmd
}

func runLogsCmd(opts *LogOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	realLogType, err := search.NewLogTypeFromValue(opts.LogType)
	if err != nil {
		return fmt.Errorf("invalid log type %s: %v", opts.LogType, err)
	}

	request := client.NewApiGetLogsRequest().
		// Offset is 0 based
		WithOffset(opts.Start - 1).
		WithLength(opts.Entries).
		WithType(*realLogType)

	if opts.IndexName != nil {
		request = request.WithIndexName(*opts.IndexName)
	}

	opts.IO.StartProgressIndicatorWithLabel("Retrieving logs")
	res, err := client.GetLogs(request)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	return p.Print(opts.IO, res)
}
