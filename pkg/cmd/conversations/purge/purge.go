package purge

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

// PurgeOptions holds the dependencies and flags for the purge command.
type PurgeOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	AgentStudioClient func() (*agentStudio.APIClient, error)
	AgentID           string
	StartDate         string
	EndDate           string
	DoConfirm         bool
}

// NewPurgeCmd returns the `conversations purge` command.
func NewPurgeCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &PurgeOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		AgentStudioClient: f.AgentStudioClient,
	}

	cmd := &cobra.Command{
		Use:   "purge <agent-id>",
		Short: "Bulk-delete Agent Studio conversations in a date range",
		Long:  "Deletes conversations for an agent matching --start-date and/or --end-date. At least one date is required so the delete is bounded.",
		Args: validators.ExactArgsWithMsg(
			1,
			"agentstudio conversations purge requires a <agent-id> argument.",
		),
		Annotations: map[string]string{
			"acls": "logs",
		},
		Example: heredoc.Doc(`
			# Purge conversations before a date
			$ algolia agentstudio conversations purge 11111111-1111-4111-8111-111111111111 --end-date 2026-01-01 --confirm
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]

			if opts.StartDate == "" && opts.EndDate == "" {
				return cmdutil.FlagErrorf("`--start-date` or `--end-date` is required")
			}

			if !opts.DoConfirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf(
						"--confirm required when non-interactive shell is detected",
					)
				}
				var confirmed bool
				err := prompt.Confirm(
					fmt.Sprintf("Purge conversations for agent %q?", opts.AgentID),
					&confirmed,
				)
				if err != nil {
					return fmt.Errorf("failed to prompt: %w", err)
				}
				if !confirmed {
					return nil
				}
			}

			return runPurgeCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Created after this date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "Created before this date (YYYY-MM-DD)")
	cmd.Flags().BoolVarP(&opts.DoConfirm, "confirm", "y", false, "Skip confirmation prompt")
	return cmd
}

func runPurgeCmd(opts *PurgeOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	req := client.NewApiDeleteAgentConversationsRequest(opts.AgentID)
	if opts.StartDate != "" {
		req = req.WithStartDate(opts.StartDate)
	}
	if opts.EndDate != "" {
		req = req.WithEndDate(opts.EndDate)
	}

	opts.IO.StartProgressIndicatorWithLabel("Purging conversations")

	err = client.DeleteAgentConversations(req)
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(
			opts.IO.Out,
			"%s Purged conversations for agent %s\n",
			cs.SuccessIcon(),
			opts.AgentID,
		)
	}

	return nil
}
