package list

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/validators"
)

type ListOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	PrintFlags *cmdutil.PrintFlags
}

// NewListCmd creates and returns a list command for indices
func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags(),
	}
	cmd := &cobra.Command{
		Use:   "list",
		Args:  validators.NoArgs(),
		Short: "List indices",
		Example: heredoc.Doc(`
			# List indices
			$ algolia indices list
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCmd(opts)
		},
		Annotations: map[string]string{
			"runInWebCLI": "true",
			"acls":        "listIndexes",
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching indices")
	res, err := client.ListIndices(client.NewApiListIndicesRequest())
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
		return p.Print(opts.IO, res)
	}

	if err := opts.IO.StartPager(); err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "error starting pager: %v\n", err)
	}
	defer opts.IO.StopPager()

	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("NAME", nil, nil)
		table.AddField("ENTRIES", nil, nil)
		table.AddField("SIZE", nil, nil)
		table.AddField("UPDATED AT", nil, nil)
		table.AddField("CREATED AT", nil, nil)
		table.AddField("LAST BUILD DURATION", nil, nil)
		table.AddField("PRIMARY", nil, nil)
		table.AddField("REPLICAS", nil, nil)
		table.EndRow()
	}

	for _, index := range res.Items {
		var primary string
		if index.Primary == nil {
			primary = ""
		} else {
			primary = *index.Primary
		}
		updatedAt, err := parseTime(index.UpdatedAt)
		if err != nil {
			return fmt.Errorf("can't parse %s into a time struct", index.UpdatedAt)
		}
		createdAt, err := parseTime(index.CreatedAt)
		if err != nil {
			return fmt.Errorf("can't parse %s into a time struct", index.CreatedAt)
		}
		// Prevent integer overflow
		if index.DataSize < 0 {
			index.DataSize = 0
		}
		table.AddField(index.Name, nil, nil)
		table.AddField(humanize.Comma(int64(index.Entries)), nil, nil)
		table.AddField(humanize.Bytes(uint64(index.DataSize)), nil, nil)
		table.AddField(updatedAt, nil, nil)
		table.AddField(createdAt, nil, nil)
		table.AddField(strconv.Itoa(int(index.LastBuildTimeS))+"s", nil, nil)
		table.AddField(primary, nil, nil)
		table.AddField(fmt.Sprintf("%v", index.Replicas), nil, nil)
		table.EndRow()
	}
	return table.Render()
}

// parseTime parses the string from the API response into a relative time string
func parseTime(timeAsString string) (string, error) {
	const layout = "2006-01-02T15:04:05.999Z"

	// This *should* restore the previous behavior when UpdatedAt is empty
	if timeAsString == "" {
		return "a long while ago", nil
	}

	t, err := time.Parse(layout, timeAsString)
	if err != nil {
		return "", err
	}

	return humanize.Time(t), nil
}
