package list

import (
	"fmt"
	"sort"
	"time"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/apikeys/shared"
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

// NewListCmd creates and returns a list command for API Keys.
func NewListCmd(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.V4_SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags(),
	}
	cmd := &cobra.Command{
		Use:  "list",
		Args: validators.NoArgs(),
		Annotations: map[string]string{
			"acls": "admin",
		},
		Short: "List API keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runListCmd(opts)
		},
	}

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

// runListCmd executes the list command
func runListCmd(opts *ListOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Fetching API Keys")
	res, err := client.ListApiKeys()
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}
		for _, key := range res.Keys {
			keyResult := shared.JSONKey{
				ACL:                    key.Acl,
				CreatedAt:              key.CreatedAt,
				Description:            *key.Description,
				Indexes:                key.Indexes,
				MaxQueriesPerIPPerHour: key.MaxQueriesPerIPPerHour,
				MaxHitsPerQuery:        key.MaxHitsPerQuery,
				Referers:               key.Referers,
				QueryParameters:        key.QueryParameters,
				Validity:               key.Validity,
				Value:                  *key.Value,
			}

			if err := p.Print(opts.IO, keyResult); err != nil {
				return err
			}
		}
		return nil
	}

	table := printers.NewTablePrinter(opts.IO)
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
		return res.Keys[i].CreatedAt > res.Keys[j].CreatedAt
	})

	for _, key := range res.Keys {
		validity := time.Duration(*key.Validity) * time.Second
		createdAt := time.Unix(key.CreatedAt, 0)

		table.AddField(*key.Value, nil, nil)
		table.AddField(*key.Description, nil, nil)
		table.AddField(fmt.Sprintf("%v", key.Acl), nil, nil)
		table.AddField(fmt.Sprintf("%v", key.Indexes), nil, nil)
		table.AddField(func() string {
			if *key.Validity == 0 {
				return "Never expire"
			} else {
				return humanize.Time(time.Now().Add(validity))
			}
		}(), nil, nil)
		if key.MaxHitsPerQuery != nil {
			table.AddField(humanize.Comma(int64(*key.MaxHitsPerQuery)), nil, nil)
		} else {
			table.AddField("0", nil, nil)
		}
		if key.MaxQueriesPerIPPerHour != nil {
			table.AddField(humanize.Comma(int64(*key.MaxQueriesPerIPPerHour)), nil, nil)
		} else {
			table.AddField("0", nil, nil)
		}
		table.AddField(fmt.Sprintf("%v", key.Referers), nil, nil)
		table.AddField(humanize.Time(createdAt), nil, nil)
		table.EndRow()
	}
	return table.Render()
}
