package userdata

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/pkg/validators"
)

type DeleteOptions struct {
	IO                *iostreams.IOStreams
	Ctx               context.Context
	AgentStudioClient func() (*agentstudio.Client, error)
	UserToken         string
	DryRun            bool
	DoConfirm         bool
}

func newDeleteCmd(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
	}
	var confirm bool
	cmd := &cobra.Command{
		Use:   "delete <user-token> [--confirm]",
		Short: "Erase ALL conversations + memories tied to a user token (irreversible)",
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.UserToken = args[0]
			opts.Ctx = cmd.Context()
			if opts.UserToken == "" {
				return cmdutil.FlagErrorf("user-token must not be empty")
			}
			if strings.Contains(opts.UserToken, "/") {
				return cmdutil.FlagErrorf("%s", rejectSlashMsg)
			}
			if !confirm && !opts.DryRun {
				if !opts.IO.CanPrompt() {
					return cmdutil.FlagErrorf("--confirm required when non-interactive shell is detected")
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
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Print what would be deleted without calling the API")
	return cmd
}

func runDeleteCmd(opts *DeleteOptions) error {
	if opts.DryRun {
		fmt.Fprintf(opts.IO.Out, "Dry run: would DELETE /1/user-data/%s\n", opts.UserToken)
		return nil
	}
	if opts.DoConfirm {
		var ok bool
		msg := fmt.Sprintf(
			"Erase ALL conversations and memories for user token %q across every agent in this app? This cannot be undone.",
			opts.UserToken,
		)
		if err := prompt.Confirm(msg, &ok); err != nil {
			return fmt.Errorf("failed to prompt: %w", err)
		}
		if !ok {
			return nil
		}
	}
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Erasing user data")
	err = client.DeleteUserData(ctxOrBackground(opts.Ctx), opts.UserToken)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	cs := opts.IO.ColorScheme()
	if opts.IO.IsStdoutTTY() {
		fmt.Fprintf(opts.IO.Out, "%s Erased user data for %s\n", cs.SuccessIcon(), opts.UserToken)
	}
	return nil
}
