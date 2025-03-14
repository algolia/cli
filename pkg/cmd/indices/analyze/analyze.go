package analyze

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/internal/analyze"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/validators"
)

type StatsOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.APIClient, error)

	Index        string
	BrowseParams search.BrowseParamsObject
	NoLimit      bool
	Only         string

	PrintFlags *cmdutil.PrintFlags
}

// NewAnalyzeCmd creates and returns an analyze command for index objects
func NewAnalyzeCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &StatsOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
		PrintFlags:   cmdutil.NewPrintFlags(),
	}

	cmd := &cobra.Command{
		Use:               "analyze <index>",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Annotations: map[string]string{
			"acls": "browse,settings",
		},
		Short: "Display records statistics for the specified index",
		Long: heredoc.Doc(`
			This command displays records statistics - frequency of the attributes and their types - for the specified index.
			This can be useful to help you identify individual records (or attributes) within an index that do not conform to the rest of the dataset (e.g. numeric attributes that have null values).

			Per default, the command will only analyze the first 1000 records. You can use the "--no-limit" flag to analyze all the records (this might take a while, depending on the number of records in your index).

			You can also use the "--only" flag to only analyze a specific attribute. In this case, the command will display the frequency of the values for this attribute.
		`),
		Example: heredoc.Doc(`
			# Display records statistics for the "MOVIES" index for the first 1000 records
			$ algolia index analyze MOVIES

			# Display records statistics for the "MOVIES" index without limit
			$ algolia index analyze MOVIES --no-limit

			# Display records statistics for the "MOVIES" index for the first 1000 records and output the result as JSON
			$ algolia index analyze MOVIES -o json

			# Display records statistics for the "MOVIES" index with the "actors" attribute only
			$ algolia index analyze MOVIES --only actors
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Index = args[0]

			browseParamsMap, err := cmdutil.FlagValuesMap(
				cmd.Flags(),
				cmdutil.BrowseParamsObject...)
			if err != nil {
				return err
			}

			// Convert map to object
			tmp, err := json.Marshal(browseParamsMap)
			if err != nil {
				return err
			}
			err = json.Unmarshal(tmp, &opts.BrowseParams)
			if err != nil {
				return err
			}

			return runAnalyzeCmd(opts)
		},
	}

	cmd.Flags().
		BoolVarP(&opts.NoLimit, "no-limit", "n", false, "If set, the command will not limit the number of objects to analyze. Otherwise, the default limit is 1000 objects.")
	cmd.Flags().
		StringVarP(&opts.Only, "only", "", "", "If set, the command will only analyze the specified attribute. Chosen attribute values statistics will be shown in the output.")

	cmdutil.AddBrowseParamsObjectFlags(cmd)
	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runAnalyzeCmd(opts *StatsOptions) error {
	client, err := opts.SearchClient()
	io := opts.IO
	if err != nil {
		return err
	}

	var records []search.Hit
	counter := 1

	// With `--no-limit`, we iterate over all records
	if opts.NoLimit {
		// Get the number of records (just for printing; maybe we can avoid this?)
		res, err := client.SearchSingleIndex(
			client.NewApiSearchSingleIndexRequest(opts.Index).
				WithSearchParams(search.SearchParamsObjectAsSearchParams(search.NewEmptySearchParamsObject().SetQuery("").SetHitsPerPage(0))),
		)
		if err != nil {
			return err
		}
		limit := int(*res.NbHits)
		io.StartProgressIndicatorWithLabel(fmt.Sprintf("Analyzing 0/%d objects", limit))
		err = client.BrowseObjects(
			opts.Index,
			opts.BrowseParams,
			search.WithAggregator(func(response any, err error) {
				if err != nil {
					return
				}
				for _, hit := range response.(*search.BrowseResponse).Hits {
					io.UpdateProgressIndicatorLabel(
						fmt.Sprintf("Analyzing %d/%d objects", counter, limit),
					)
					records = append(records, hit)
					counter++
				}
			}),
		)
		if err != nil {
			io.StopProgressIndicator()
			return err
		}
	} else {
		// Without `--no-limit`, just iterate over the first 1,000 records.
		// Since `BrowseObjects` can't be terminated early, we use the regular Browse method here.
		res, err := client.Browse(
			client.
				NewApiBrowseRequest(opts.Index).
				WithBrowseParams(search.BrowseParamsObjectAsBrowseParams(&opts.BrowseParams)),
		)
		limit := int(*res.NbHits)
		io.StartProgressIndicatorWithLabel(fmt.Sprintf("Analyzing 0/%d objects", limit))
		if err != nil {
			return err
		}
		for _, hit := range res.Hits {
			io.UpdateProgressIndicatorLabel(
				fmt.Sprintf("Analyzing %d/%d objects", counter, limit),
			)
			records = append(records, hit)
			counter++
		}
	}

	settings, err := client.GetSettings(client.NewApiGetSettingsRequest(opts.Index))
	if err != nil {
		io.StopProgressIndicator()
		return err
	}

	stats, err := analyze.ComputeStats(records, *settings, opts.Only)
	if err != nil {
		io.StopProgressIndicator()
		return err
	}

	io.StopProgressIndicator()

	if opts.PrintFlags.OutputFlagSpecified() && opts.PrintFlags.OutputFormat != nil {
		p, err := opts.PrintFlags.ToPrinter()
		if err != nil {
			return err
		}

		// If "only" is specified, we need to format the output differently
		if opts.Only != "" {
			var statsArray []interface{}
			for key := range stats.Attributes {
				for value := range stats.Attributes[key].Values {
					statsArray = append(statsArray, map[string]interface{}{
						"value": value,
						"count": stats.Attributes[key].Values[value],
					})
				}
			}
			return p.Print(io, statsArray)
		}

		return p.Print(io, stats)
	}

	if opts.Only != "" {
		return printSingleAttributeStats(stats, opts)
	}

	return printStats(stats, opts)
}

// printStats prints the global stats for the index in a table format
func printStats(stats *analyze.Stats, opts *StatsOptions) error {
	cs := opts.IO.ColorScheme()
	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("KEY", nil, nil)
		table.AddField("COUNT", nil, nil)
		table.AddField("%", nil, nil)
		table.AddField("TYPES", nil, nil)
		table.AddField("USED IN SETTINGS", nil, nil)
		table.EndRow()
	}

	formatTypes := func(types map[analyze.AttributeType]float64) string {
		var result string
		for key, value := range types {
			if value < 1 {
				result += cs.Red(fmt.Sprintf("%s: %.2f%%", key, value))
			} else if value < 5 {
				result += cs.Yellow(fmt.Sprintf("%s: %.2f%%", key, value))
			} else {
				result += fmt.Sprintf("%s: %.2f%%", key, value)
			}
			result += ", "
		}
		result = result[:len(result)-2] // Remove the last ", "
		return result
	}

	// We need to sort the keys to have a consistent output
	sorted := make([]string, 0, len(stats.Attributes))
	for key := range stats.Attributes {
		sorted = append(sorted, key)
	}
	sort.Strings(sorted)

	for _, key := range sorted {
		// Print colorized output depending on the percentage
		// If <1%: red, if <5%: yellow
		color := func(s string) string { return s }
		if stats.Attributes[key].Percentage < 1 {
			color = cs.Red
		} else if stats.Attributes[key].Percentage < 5 {
			color = cs.Yellow
		}
		value := stats.Attributes[key]
		table.AddField(color(key), nil, nil)
		table.AddField(color(fmt.Sprintf("%d", value.Count)), nil, nil)
		table.AddField(color(fmt.Sprintf("%.2f%%", value.Percentage)), nil, nil)
		table.AddField(formatTypes(value.Types), nil, nil)
		table.AddField(fmt.Sprintf("%v", value.InSettings), nil, nil)

		table.EndRow()
	}

	return table.Render()
}

// printSingleAttributeStats prints the stats for a single attribute in a table format
func printSingleAttributeStats(stats *analyze.Stats, opts *StatsOptions) error {
	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("VALUE", nil, nil)
		table.AddField("COUNT", nil, nil)
		table.AddField("%", nil, nil)
		table.EndRow()
	}

	for key := range stats.Attributes {
		value := stats.Attributes[key]
		// Order the values by count (descending)
		sorted := make([]interface{}, 0, len(value.Values))
		for key := range value.Values {
			sorted = append(sorted, key)
		}
		sort.Slice(sorted, func(i, j int) bool {
			return value.Values[sorted[i]] > value.Values[sorted[j]]
		})

		for _, v := range sorted {
			table.AddField(fmt.Sprintf("%v", v), nil, nil)
			table.AddField(fmt.Sprintf("%d", value.Values[v]), nil, nil)
			table.AddField(
				fmt.Sprintf("%.2f%%", float64(value.Values[v])*100/float64(stats.TotalRecords)),
				nil,
				nil,
			)
			table.EndRow()
		}
	}

	return table.Render()
}
