package response

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/genai/response/delete"
	"github.com/algolia/cli/pkg/cmd/genai/response/generate"
	"github.com/algolia/cli/pkg/cmd/genai/response/get"
	"github.com/algolia/cli/pkg/cmd/genai/response/list"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewResponseCmd returns a new command to manage your Algolia GenAI responses.
func NewResponseCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "response",
		Aliases: []string{"responses"},
		Short:   "Manage your GenAI responses",
		Long: heredoc.Doc(`
			Manage your Algolia GenAI responses.

			Responses are generated using your prompts and data sources to answer your queries.
		`),
	}

	cmd.AddCommand(generate.NewGenerateCmd(f, nil))
	cmd.AddCommand(get.NewGetCmd(f, nil))
	cmd.AddCommand(delete.NewDeleteCmd(f, nil))
	cmd.AddCommand(list.NewListCmd(f))

	return cmd
}
