package unblock

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
)

// UnblockOptions holds the options for the stats command.
type UnblockOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	CrawlerClient func() (*crawler.Client, error)

	ID        string
	DoConfirm bool

	PrintFlags *cmdutil.PrintFlags
}

// NewUnblockCmd creates and returns a command to unblock a crawler.
func NewUnblockCmd(f *cmdutil.Factory, runF func(*UnblockOptions) error) *cobra.Command {
	opts := &UnblockOptions{
		IO:            f.IOStreams,
		Config:        f.Config,
		CrawlerClient: f.CrawlerClient,
	}

	var confirm bool

	cmd := &cobra.Command{
		Use:               "unblock <crawler_id>",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.CrawlerIDs(opts.CrawlerClient),
		Short:             "Unblock a crawler",
		Long: heredoc.Doc(`
			Unblock a crawler by cancelling the specific task that is currently blocking it.
		`),
		Example: heredoc.Doc(`
			# Unblock the crawler with the ID "my-crawler"
			$ algolia crawler unblock my-crawler
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			if runF != nil {
				return runF(opts)
			}

			if !confirm {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
				}
				opts.DoConfirm = true
			}

			return runUnblockCmd(opts)
		},
	}

	return cmd
}

// runUnblockCmd executes the unblock command.
func runUnblockCmd(opts *UnblockOptions) error {
	client, err := opts.CrawlerClient()
	if err != nil {
		return err
	}

	// Get the crawler and check if it is actually blocked.
	crawler, err := client.Get(opts.ID, false)
	if err != nil {
		return err
	}
	if crawler.BlockingTaskID == "" {
		return fmt.Errorf("crawler %q is not blocked", opts.ID)
	}

	if opts.DoConfirm {
		var confirmed bool
		err := prompt.Confirm(fmt.Sprintf("Are you sure you want to unblock the crawler %q? \nBlocking error is: %s", opts.ID, crawler.BlockingError), &confirmed)
		if err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !confirmed {
			return nil
		}
	}

	// Cancel the task blocking the crawler.
	if err := client.CancelTask(crawler.ID, crawler.BlockingTaskID); err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Unblocked crawler %s\n", cs.SuccessIcon(), cs.Bold(opts.ID))
	}

	return nil
}
