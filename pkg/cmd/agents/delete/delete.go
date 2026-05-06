package delete

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

type DeleteOptions struct {
	IO  *iostreams.IOStreams
	Ctx context.Context

	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags

	AgentID   string
	DryRun    bool
	DoConfirm bool
}

func NewDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags(),
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <agent-id> [--confirm]",
		Short: "Delete an Agent Studio agent",
		Long: heredoc.Doc(`
			Soft-delete an Agent Studio agent. The backend keeps the row with
			a deleted flag, so recovery is a platform-side ops concern;
			treat this command as terminal from the CLI's perspective.

			Use --dry-run to fetch the target agent and preview what would
			be deleted, without actually deleting.
		`),
		Example: heredoc.Doc(`
			# Interactive delete (asks for confirmation)
			$ algolia agents delete 11111111-1111-1111-1111-111111111111

			# Skip the prompt (required in non-interactive shells / CI)
			$ algolia agents delete 11111111-1111-1111-1111-111111111111 -y

			# Preview without deleting
			$ algolia agents delete 11111111-1111-1111-1111-111111111111 --dry-run
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}

			// Mirror objects/delete: --confirm/-y is the explicit
			// "I have confirmed" flag. Without it, we either prompt (TTY)
			// or refuse (non-TTY) — but only when a real DELETE will
			// happen; --dry-run is non-destructive so it bypasses the
			// confirmation requirement entirely.
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
			return runDeleteCmd(opts)
		},
	}

	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "Skip confirmation prompt")
	cmd.Flags().
		BoolVar(&opts.DryRun, "dry-run", false, "Fetch and preview the agent without deleting it")

	opts.PrintFlags.AddFlags(cmd)

	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	cs := opts.IO.ColorScheme()
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	// Fetch the agent up front for two reasons:
	//   1. Surface 404 cleanly (matches GET /1/agents/{id}'s shape).
	//   2. Show name + status in the prompt / dry-run output, so the user
	//      can sanity-check they're about to delete the right thing.
	agent, err := client.GetAgent(ctx, opts.AgentID)
	if err != nil {
		return err
	}

	if opts.DryRun {
		fmt.Fprintf(opts.IO.Out, "Dry run: would DELETE /1/agents/%s\n", opts.AgentID)
		fmt.Fprintf(opts.IO.Out, "  name:   %s\n", agent.Name)
		fmt.Fprintf(opts.IO.Out, "  status: %s\n", agent.Status)
		return nil
	}

	if opts.DoConfirm {
		var confirmed bool
		err := prompt.Confirm(
			fmt.Sprintf("Delete agent %q (%s)?", agent.Name, opts.AgentID),
			&confirmed,
		)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	opts.IO.StartProgressIndicatorWithLabel("Deleting agent")
	err = client.DeleteAgent(ctx, opts.AgentID)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Deleted agent %s\n", cs.SuccessIcon(), opts.AgentID)
	}
	return nil
}
