package bulkdelete

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

// BulkDeleteOptions holds the dependencies and flags for the bulk-delete command.
type BulkDeleteOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	AgentID           string
	DomainIDs         []string
	DoConfirm         bool
}

// NewBulkDeleteCmd returns the `agents domains bulk-delete` command.
func NewBulkDeleteCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &BulkDeleteOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
	}

	cmd := &cobra.Command{
		Use:               "bulk-delete <agent-id> <domain-id>...",
		Short:             "Delete multiple allowed domains from an agent",
		Args:              validators.AtLeastNArgs(2),
		ValidArgsFunction: cmdutil.AgentIDs(opts.AgentStudioClient),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Delete multiple allowed domains from the agent "my-agent"
			$ algolia agents domains bulk-delete my-agent domain_1 domain_2 --confirm
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.DomainIDs = args[1:]

			if !opts.DoConfirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				var confirmed bool
				err := prompt.Confirm(
					fmt.Sprintf("Delete %d allowed domain(s)?", len(opts.DomainIDs)),
					&confirmed,
				)
				if err != nil {
					return fmt.Errorf("failed to prompt: %w", err)
				}
				if !confirmed {
					return nil
				}
			}

			return runBulkDeleteCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.DoConfirm, "confirm", "y", false, "Skip confirmation prompt")

	return cmd
}

func runBulkDeleteCmd(opts *BulkDeleteOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Deleting allowed domains")

	bulkDelete := agentStudio.NewAllowedDomainBulkDelete(opts.DomainIDs)
	err = client.BulkDeleteAllowedDomains(
		client.NewApiBulkDeleteAllowedDomainsRequest(opts.AgentID, bulkDelete),
	)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out, "%s Deleted %d allowed domain(s)\n", cs.SuccessIcon(), len(opts.DomainIDs))
	}

	return nil
}
