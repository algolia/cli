package prompt

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/genai/prompt/create"
	"github.com/algolia/cli/pkg/cmd/genai/prompt/delete"
	"github.com/algolia/cli/pkg/cmd/genai/prompt/get"
	"github.com/algolia/cli/pkg/cmd/genai/prompt/list"
	"github.com/algolia/cli/pkg/cmd/genai/prompt/update"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewPromptCmd returns a new command to manage your Algolia GenAI prompts.
func NewPromptCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "prompt",
		Aliases: []string{"prompts"},
		Short:   "Manage your GenAI prompts",
		Long: heredoc.Doc(`
			Manage your Algolia GenAI prompts.

			Prompts define the instructions for how the GenAI toolkit should generate responses.
		`),
	}

	cmd.AddCommand(create.NewCreateCmd(f, nil))
	cmd.AddCommand(update.NewUpdateCmd(f, nil))
	cmd.AddCommand(delete.NewDeleteCmd(f, nil))
	cmd.AddCommand(list.NewListCmd(f))
	cmd.AddCommand(get.NewGetCmd(f, nil))

	return cmd
}
