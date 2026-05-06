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
		Use:   "purge <agent-id> (--start-date YYYY-MM-DD | --end-date YYYY-MM-DD) [--confirm]",
		Short: "Bulk-delete conversations for an agent within a date range",
		Long: heredoc.Doc(`
			Bulk-delete persisted conversations for an agent.

			At least one of --start-date / --end-date is REQUIRED.
			Background: the OpenAPI spec marks both query params as
			optional and reads as if dateless DELETE wipes everything,
			but the live backend rejects dateless requests with
			"400 At least one filter is required." The CLI surfaces
			this as a flag-level error rather than a server round-trip.

			If you genuinely want to wipe every conversation for an
			agent, pass an open-ended range — e.g.
			"--start-date 1970-01-01" or "--end-date 9999-12-31".
			Both bounds are forwarded verbatim (YYYY-MM-DD; backend
			validates and 422s on bad input).

			Like "agents delete", interactive use prompts and
			non-interactive use requires --confirm. --dry-run previews
			the URL and bypasses both.
		`),
		Example: heredoc.Doc(`
			# Purge a specific month
			$ algolia agents conversations purge <agent-id> --start-date 2026-01-01 --end-date 2026-01-31

			# Purge everything from a given date onwards (open-ended)
			$ algolia agents conversations purge <agent-id> --start-date 2026-01-01 -y

			# Preview without sending
			$ algolia agents conversations purge <agent-id> --start-date 2026-01-01 --dry-run
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			if opts.StartDate == "" && opts.EndDate == "" {
				return cmdutil.FlagErrorf(
					"at least one of --start-date / --end-date is required " +
						"(backend rejects dateless purge with 400 \"At least one filter is required\")",
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
	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Print what would be purged without calling the API")
	return cmd
}

func runPurgeCmd(opts *PurgeOptions) error {
	scope := purgeScope(opts.StartDate, opts.EndDate)

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

func purgeScope(start, end string) string {
	switch {
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
