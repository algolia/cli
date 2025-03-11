package datasource

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/genai/datasource/create"
	"github.com/algolia/cli/pkg/cmd/genai/datasource/delete"
	"github.com/algolia/cli/pkg/cmd/genai/datasource/get"
	"github.com/algolia/cli/pkg/cmd/genai/datasource/list"
	"github.com/algolia/cli/pkg/cmd/genai/datasource/update"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewDataSourceCmd returns a new command to manage your Algolia GenAI data sources.
func NewDataSourceCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "datasource",
		Aliases: []string{"datasources", "data-sources", "data-source", "ds"},
		Short:   "Manage your GenAI data sources",
		Long: heredoc.Doc(`
			Manage your Algolia GenAI data sources.

			Data sources provide the contexts that the toolkit uses to generate relevant responses to your users' queries.
		`),
	}

	cmd.AddCommand(create.NewCreateCmd(f, nil))
	cmd.AddCommand(update.NewUpdateCmd(f, nil))
	cmd.AddCommand(delete.NewDeleteCmd(f, nil))
	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(get.NewGetCmd(f, nil))

	return cmd
}
