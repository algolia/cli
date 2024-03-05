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
type AddOptions struct {
	config config.IConfig
	IO     *iostreams.IOStreams
}

// NewListCmd returns a new instance of ListCmd
func NewListCmd(f *cmdutil.Factory, runF func(*AddOptions) error) *cobra.Command {
	opts := &AddOptions{
		IO:     f.IOStreams,
		config: f.Config,
	}
	cmd := &cobra.Command{
		Use:   "list",
		Args:  validators.NoArgs(),
		Short: "List the configured profile(s)",
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
func runListCmd(opts *AddOptions) error {
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
		client := search.NewClient(profile.ApplicationID, profile.AdminAPIKey)
		res, err := client.ListIndices()
		if err != nil {
			return fmt.Errorf("could not retrieve indices for profile %q (AppID %q): %w", profile.Name, profile.ApplicationID, err)
		}

		table.AddField(profile.Name, nil, nil)
		table.AddField(profile.ApplicationID, nil, nil)
		table.AddField(fmt.Sprintf("%d", len(res.Items)), nil, nil)
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
