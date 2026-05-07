package domains

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
)

// NewDomainsCmd is the parent for `algolia agents domains <verb>`.
func NewDomainsCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domains",
		Short: "Manage Agent Studio per-agent CORS allowlist",
		Long: heredoc.Doc(`
			Manage allowed domains (CORS allowlist) per agent. Domain
			values are forwarded verbatim — patterns like *.example.com
			and exact origins like https://app.example.com are both
			accepted; the backend validates.
		`),
	}
	cmd.AddCommand(newListCmd(f, nil))
	cmd.AddCommand(newGetCmd(f, nil))
	cmd.AddCommand(newCreateCmd(f, nil))
	cmd.AddCommand(newDeleteCmd(f, nil))
	cmd.AddCommand(newBulkInsertCmd(f, nil))
	cmd.AddCommand(newBulkDeleteCmd(f, nil))
	return cmd
}
