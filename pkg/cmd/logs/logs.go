package logs

import (
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/logs/list"
	"github.com/algolia/cli/pkg/cmdutil"
)

// NewLogsCmd returns a new command for retrieving logs
func NewLogsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Retrieve your Algolia Search API logs",
	}

	cmd.AddCommand(list.NewListCmd(f, nil))
	return cmd
}
