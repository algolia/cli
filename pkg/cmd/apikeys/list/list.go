package list

import (
	"fmt"
	"sort"
	"time"

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

// NewListCmd creates and returns a list command for API Keys.
func NewListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &ListOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}
	cmd := &cobra.Command{
		Use:   "list",
		Args:  validators.NoArgs,
		Short: "List API keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListCmd(opts)
		},
	}

	cmdutil.AddJSONFlags(cmd, &opts.Exporter)

	return cmd
}

// runListCmd executes the list command
func runListCmd(opts *ListOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching API Keys")
	res, err := client.ListAPIKeys()
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.Exporter != nil {
		return opts.Exporter.Write(opts.IO, res.Keys)
	}

	table := utils.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("KEY", nil, nil)
		table.AddField("DESCRIPTION", nil, nil)
		table.AddField("ACL", nil, nil)
		table.AddField("INDICES", nil, nil)
		table.AddField("VALIDITY", nil, nil)
		table.AddField("MAX HITS PER QUERY", nil, nil)
		table.AddField("MAX QUERIES PER IP PER HOUR", nil, nil)
		table.AddField("REFERERS", nil, nil)
		table.AddField("CREATED AT", nil, nil)
		table.EndRow()
	}

	// Sort API Keys by createdAt
	sort.Slice(res.Keys, func(i, j int) bool {
		return res.Keys[i].CreatedAt.After(res.Keys[j].CreatedAt)
	})

	for _, key := range res.Keys {
		table.AddField(key.Value, nil, nil)
		table.AddField(key.Description, nil, nil)
		table.AddField(fmt.Sprintf("%v", key.ACL), nil, nil)
		table.AddField(fmt.Sprintf("%v", key.Indexes), nil, nil)
		table.AddField(func() string {
			if key.Validity == 0 {
				return "Never expire"
			} else {
				return humanize.Time(time.Now().Add(key.Validity))
			}
		}(), nil, nil)
		table.AddField(humanize.Comma(int64(key.MaxHitsPerQuery)), nil, nil)
		table.AddField(humanize.Comma(int64(key.MaxQueriesPerIPPerHour)), nil, nil)
		table.AddField(fmt.Sprintf("%v", key.Referers), nil, nil)
		table.AddField(humanize.Time(key.CreatedAt), nil, nil)
		table.EndRow()
	}
	return table.Render()
}
