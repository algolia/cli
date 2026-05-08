package conversations

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type PurgeOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)

	AgentID   string
	StartDate string
	EndDate   string
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
			non-interactive use requires --confirm.
		`),
		Example: heredoc.Doc(`
			# Purge a specific month
			$ algolia agents conversations purge <agent-id> --start-date 2026-01-01 --end-date 2026-01-31

			# Purge everything from a given date onwards (open-ended)
			$ algolia agents conversations purge <agent-id> --start-date 2026-01-01 -y
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
			doConfirm, err := shared.ResolveConfirm(opts.IO, confirm)
			if err != nil {
				return err
			}
			opts.DoConfirm = doConfirm
			if runF != nil {
				return runF(opts)
			}
			return runPurgeCmd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.StartDate, "start-date", "", "Purge conversations >= date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.EndDate, "end-date", "", "Purge conversations <= date (YYYY-MM-DD)")
	shared.AddConfirmFlag(cmd, &confirm)
	return cmd
}

func runPurgeCmd(opts *PurgeOptions) error {
	scope := purgeScope(opts.StartDate, opts.EndDate)

	if opts.DoConfirm {
		ok, err := shared.Confirm(
			fmt.Sprintf("Purge conversations on agent %s — %s ?", opts.AgentID, scope),
		)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}

	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	ctx := shared.OrBackground(opts.Ctx)

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
