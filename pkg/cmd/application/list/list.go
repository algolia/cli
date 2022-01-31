package list

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/utils"
	"github.com/algolia/cli/pkg/validators"
)

// ListOptions represents the options for the list command
type AddOptions struct {
	config *config.Config
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
		Args:  validators.NoArgs,
		Short: "List the configured application(s)",
		Long:  `List the configured application(s).`,
		Example: heredoc.Doc(`
			$ algolia application list
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

// runListCmd executes the add command
func runListCmd(opts *AddOptions) error {

	table := utils.NewTablePrinter(opts.IO)
	if table.IsTTY() {
		table.AddField("NAME", nil, nil)
		table.AddField("ID", nil, nil)
		table.AddField("NUMBER OF INDICES", nil, nil)
		table.EndRow()
	}

	for name, appID := range opts.config.Applications() {
		opts.config.App.Name = name
		key, err := opts.config.App.GetAdminAPIKey()
		if err != nil {
			return err
		}

		client := search.NewClient(appID, key)
		res, err := client.ListIndices()
		if err != nil {
			return err
		}

		table.AddField(name, nil, nil)
		table.AddField(appID, nil, nil)
		table.AddField(fmt.Sprintf("%d", len(res.Items)), nil, nil)
		table.EndRow()
	}
	return table.Render()
}
