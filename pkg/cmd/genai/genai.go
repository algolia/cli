package genai

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/genai/datasource"
	"github.com/algolia/cli/pkg/cmd/genai/prompt"
	"github.com/algolia/cli/pkg/cmd/genai/response"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewGenaiCmd returns a new command to manage your Algolia GenAI Toolkit.
func NewGenaiCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "genai",
		Aliases: []string{"genai-toolkit"},
		Short:   "Manage your Algolia GenAI Toolkit",
	}

	// Add data source commands
	cmd.AddCommand(datasource.NewDataSourceCmd(f))

	// Add prompt commands
	cmd.AddCommand(prompt.NewPromptCmd(f))

	// Add response commands
	cmd.AddCommand(response.NewResponseCmd(f))

	return cmd
}
