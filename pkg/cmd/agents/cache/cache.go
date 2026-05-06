package cache

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

// NewCacheCmd is the parent for `algolia agents cache <verb>`.
//
// Today there's exactly one verb (`invalidate`); kept as a sub-group
// rather than a flat `agents cache-invalidate` because:
//   - The parent reads as a noun (`cache`), the child as a verb
//     (`invalidate`) — natural language order, matches the established
//     CLI rhythm (`apikeys list`, `objects update`).
//   - Reserves a clean home for follow-ups (`agents cache stats`,
//     `agents cache size`, etc.) without renaming the existing surface.
//   - Mirrors the nested-group pattern Phase 6+ will reuse for
//     `agents providers <verb>` and `agents conversations <verb>`.
func NewCacheCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache",
		Short: "Inspect and invalidate Agent Studio completion caches",
		Long: heredoc.Doc(`
			Manage cached completion responses for an agent.

			Agent Studio caches completions per (agent, request hash). Use
			--no-cache on "agents try" / "agents run" to bypass for a single
			call; use "agents cache invalidate" to drop entries server-side.
		`),
	}

	cmd.AddCommand(newInvalidateCmd(f, nil))
	return cmd
}

// InvalidateOptions configures `algolia agents cache invalidate`.
type InvalidateOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)

	AgentID   string
	Before    string
	DryRun    bool
	DoConfirm bool
}

func newInvalidateCmd(f *cmdutil.Factory, runF func(*InvalidateOptions) error) *cobra.Command {
	opts := &InvalidateOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:   "invalidate <agent-id> [--before YYYY-MM-DD] [--confirm]",
		Short: "Invalidate cached completions for an agent",
		Long: heredoc.Doc(`
			Calls DELETE /1/agents/<id>/cache.

			By default, drops every cached completion for the agent. Pass
			--before YYYY-MM-DD to drop only entries created strictly
			before that date (exclusive — matches the backend's Pydantic
			parsing). Date validation is the backend's job; bad input
			surfaces a structured 422 verbatim.

			Like "agents delete", interactive use prompts to confirm and
			non-interactive use requires --confirm. Use --dry-run to
			preview without deleting.
		`),
		Example: heredoc.Doc(`
			# Wipe all cached completions for an agent (interactive)
			$ algolia agents cache invalidate 11111111-1111-1111-1111-111111111111

			# Drop only entries older than a specific date
			$ algolia agents cache invalidate 11111111-1111-1111-1111-111111111111 --before 2026-01-15

			# Skip the prompt (required in CI)
			$ algolia agents cache invalidate 11111111-1111-1111-1111-111111111111 -y

			# Preview without sending
			$ algolia agents cache invalidate 11111111-1111-1111-1111-111111111111 --dry-run
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}

			// Mirror agents delete: --confirm/-y bypasses prompt.
			// Without it, prompt in TTY, refuse in non-TTY. --dry-run
			// is non-destructive and bypasses the confirmation entirely.
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
			return runInvalidateCmd(opts)
		},
	}

	cmd.Flags().
		StringVar(&opts.Before, "before", "", "Drop entries strictly before this date (YYYY-MM-DD, exclusive)")
	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Skip confirmation prompt")
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Print what would be invalidated without calling the API")

	return cmd
}

func runInvalidateCmd(opts *InvalidateOptions) error {
	if opts.DryRun {
		fmt.Fprintf(opts.IO.Out, "Dry run: would DELETE /1/agents/%s/cache", opts.AgentID)
		if opts.Before != "" {
			fmt.Fprintf(opts.IO.Out, "?before=%s", opts.Before)
		}
		fmt.Fprintln(opts.IO.Out)
		if opts.Before == "" {
			fmt.Fprintln(opts.IO.Out, "  scope: all cached completions for this agent")
		} else {
			fmt.Fprintf(opts.IO.Out, "  scope: cached completions created before %s\n", opts.Before)
		}
		return nil
	}

	if opts.DoConfirm {
		var confirmed bool
		msg := fmt.Sprintf("Invalidate completion cache for agent %s?", opts.AgentID)
		if opts.Before != "" {
			msg = fmt.Sprintf("Invalidate completion cache for agent %s (entries before %s)?",
				opts.AgentID, opts.Before)
		}
		if err := prompt.Confirm(msg, &confirmed); err != nil {
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
	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	opts.IO.StartProgressIndicatorWithLabel("Invalidating agent cache")
	err = client.InvalidateAgentCache(ctx, opts.AgentID, opts.Before)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Invalidated cache for agent %s\n", cs.SuccessIcon(), opts.AgentID)
	}
	return nil
}
