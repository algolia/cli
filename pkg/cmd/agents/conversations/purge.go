package conversations

import (
	"context"
	"fmt"
	"net/url"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

type PurgeOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)

	AgentID   string
	StartDate string
	EndDate   string
	All       bool
	DryRun    bool
	DoConfirm bool
}

func newPurgeCmd(f *cmdutil.Factory, runF func(*PurgeOptions) error) *cobra.Command {
	opts := &PurgeOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
	}
	var confirm bool

	cmd := &cobra.Command{
		Use:   "purge <agent-id> (--all | --start-date YYYY-MM-DD | --end-date YYYY-MM-DD) [--confirm]",
		Short: "Bulk-delete conversations for an agent",
		Long: heredoc.Doc(`
			Bulk-delete persisted conversations for an agent.

			GUARDRAIL: the backend's DELETE /conversations endpoint with
			no date filter wipes EVERY conversation for the agent. To
			make that opt-in (so a typo can never trigger it), the CLI
			refuses to send a dateless purge unless --all is passed
			explicitly.

			With --start-date and/or --end-date the range is forwarded
			verbatim (YYYY-MM-DD; backend validates and 422s on bad input).

			Like "agents delete", interactive use prompts and
			non-interactive use requires --confirm. --dry-run previews
			the URL and bypasses both.
		`),
		Example: heredoc.Doc(`
			# Purge a date range (no --all needed because filter is set)
			$ algolia agents conversations purge <agent-id> --start-date 2026-01-01 --end-date 2026-01-31

			# Wipe everything (requires explicit --all + confirmation)
			$ algolia agents conversations purge <agent-id> --all -y

			# Preview without sending
			$ algolia agents conversations purge <agent-id> --all --dry-run
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			hasFilter := opts.StartDate != "" || opts.EndDate != ""
			if !opts.All && !hasFilter {
				return cmdutil.FlagErrorf(
					"refusing to purge ALL conversations for an agent: pass --all explicitly, or restrict with --start-date / --end-date",
				)
			}
			if opts.All && hasFilter {
				return cmdutil.FlagErrorf(
					"--all is mutually exclusive with --start-date / --end-date",
				)
			}
			if !confirm && !opts.DryRun {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf(
						"--confirm required when non-interactive shell is detected",
					)
				}
				opts.DoConfirm = true
			}
			if runF != nil {
				return runF(opts)
			}
			return runPurgeCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Purge conversations >= date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "Purge conversations <= date (YYYY-MM-DD)")
	cmd.Flags().
		BoolVar(&opts.All, "all", false, "Purge every conversation for this agent (mutually exclusive with date filters)")
	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Print what would be purged without calling the API")
	return cmd
}

func runPurgeCmd(opts *PurgeOptions) error {
	scope := purgeScope(opts.StartDate, opts.EndDate, opts.All)

	if opts.DryRun {
		q := url.Values{}
		if opts.StartDate != "" {
			q.Set("startDate", opts.StartDate)
		}
		if opts.EndDate != "" {
			q.Set("endDate", opts.EndDate)
		}
		path := fmt.Sprintf("/1/agents/%s/conversations", opts.AgentID)
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
		fmt.Fprintf(opts.IO.Out, "Dry run: would DELETE %s\n", path)
		fmt.Fprintf(opts.IO.Out, "  scope: %s\n", scope)
		return nil
	}

	if opts.DoConfirm {
		var confirmed bool
		err := prompt.Confirm(
			fmt.Sprintf("Purge conversations on agent %s — %s ?", opts.AgentID, scope),
			&confirmed,
		)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := ctxOrBackground(opts.Ctx)

	opts.IO.StartProgressIndicatorWithLabel("Purging conversations")
	err = client.PurgeConversations(ctx, opts.AgentID, agentstudio.PurgeConversationsParams{
		StartDate: opts.StartDate,
		EndDate:   opts.EndDate,
	})
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Purged conversations on agent %s (%s)\n",
			cs.SuccessIcon(), opts.AgentID, scope)
	}
	return nil
}

func purgeScope(start, end string, all bool) string {
	switch {
	case all:
		return "ALL conversations"
	case start != "" && end != "":
		return fmt.Sprintf("between %s and %s", start, end)
	case start != "":
		return fmt.Sprintf("from %s onwards", start)
	case end != "":
		return fmt.Sprintf("up to %s", end)
	default:
		return "(unscoped — guardrail bug; please report)"
	}
}
