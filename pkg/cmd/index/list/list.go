package list

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
	"github.com/algolia/cli/pkg/validators"
)

type ListOptions struct {
	Config *config.Config
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Exporter cmdutil.Exporter
}

// NewListCmd creates and returns a list command for indices
func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}
	cmd := &cobra.Command{
		Use:   "list",
		Args:  validators.NoArgs,
		Short: "List indices",
		Example: heredoc.Doc(`
			# List indices
			algolia index list
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCmd(opts)
		},
	}

	cmdutil.AddJSONFlags(cmd, &opts.Exporter, false)

	return cmd
}

func runListCmd(opts *ListOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching indices")
	res, err := client.ListIndices()
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.Exporter != nil {
		return opts.Exporter.Write(opts.IO, res.Items)
	}

	if err := opts.IO.StartPager(); err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "error starting pager: %v\n", err)
	}
	defer opts.IO.StopPager()

	table := utils.NewTablePrinter(opts.IO)
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
		table.AddField(index.Name, nil, nil)
		table.AddField(humanize.Comma(index.Entries), nil, nil)
		table.AddField(humanize.Bytes(uint64(index.DataSize)), nil, nil)
		table.AddField(humanize.Time(index.UpdatedAt), nil, nil)
		table.AddField(humanize.Time(index.CreatedAt), nil, nil)
		table.AddField(index.LastBuildTime.String(), nil, nil)
		table.AddField(index.Primary, nil, nil)
		table.AddField(fmt.Sprintf("%v", index.Replicas), nil, nil)
		table.EndRow()
	}
	return table.Render()
}
