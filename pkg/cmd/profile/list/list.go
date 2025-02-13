package list

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/printers"
	"github.com/algolia/cli/pkg/validators"
)

// ListOptions represents the options for the list command
type ListOptions struct {
	config config.IConfig
	IO     *iostreams.IOStreams
}

// NewListCmd returns a new instance of ListCmd
func NewListCmd(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:     f.IOStreams,
		config: f.Config,
	}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Args:    validators.NoArgs(),
		Short:   "List the configured profile(s)",
		Example: heredoc.Doc(`
			# List the configured profiles
			$ algolia profile list
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runListCmd(opts)
		},
	}

	return cmd
}

// runListCmd executes the list command
func runListCmd(opts *ListOptions) error {
	profiles := opts.config.ConfiguredProfiles()
	if len(profiles) == 0 {
		fmt.Fprintln(opts.IO.ErrOut, "No configured profiles")
		fmt.Fprintln(opts.IO.ErrOut, "Use `algolia profile add` to add a profile")
		return nil
	}

	table := printers.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("NAME", nil, nil)
		table.AddField("APPLICATION ID", nil, nil)
		table.AddField("NUMBER OF INDICES", nil, nil)
		table.AddField("DEFAULT", nil, nil)
		table.EndRow()
	}

	cs := opts.IO.ColorScheme()

	opts.IO.StartProgressIndicatorWithLabel("Fetching configured profiles")
	for _, profile := range profiles {
		table.AddField(profile.Name, nil, nil)
		table.AddField(profile.ApplicationID, nil, nil)

		apiKey := profile.APIKey
		if apiKey == "" {
			apiKey = profile.AdminAPIKey // Legacy
		}

		client := search.NewClient(profile.ApplicationID, apiKey)
		res, err := client.ListIndices()
		if err != nil {
			table.AddField(err.Error(), nil, nil)
		} else {
			table.AddField(fmt.Sprintf("%d", len(res.Items)), nil, nil)
		}

		if profile.Default {
			table.AddField(cs.SuccessIcon(), nil, nil)
		} else {
			table.AddField("", nil, nil)
		}
		table.EndRow()
	}
	opts.IO.StopProgressIndicator()
	return table.Render()
}
