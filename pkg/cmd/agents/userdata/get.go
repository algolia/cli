package userdata

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type GetOptions struct {
	IO                *iostreams.IOStreams
	Ctx               context.Context
	AgentStudioClient func() (*agentstudio.Client, error)
	UserToken         string
	OutputFile        string
}

func newGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
	}
	cmd := &cobra.Command{
		Use:   "get <user-token>",
		Short: "Dump all conversations + memories for a user token",
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
			if runF != nil {
				return runF(opts)
			}
			return runGetCmd(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.OutputFile, "output-file", "o", "",
		"Write the raw response to this file (default stdout)")
	return cmd
}

func runGetCmd(opts *GetOptions) error {
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Fetching user data")
	res, err := client.GetUserData(ctxOrBackground(opts.Ctx), opts.UserToken)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	raw, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal user data: %w", err)
	}
	if opts.OutputFile != "" {
		if err := os.WriteFile(opts.OutputFile, raw, 0o600); err != nil {
			return fmt.Errorf("write %s: %w", opts.OutputFile, err)
		}
		if opts.IO.IsStdoutTTY() {
			fmt.Fprintf(opts.IO.Out,
				"Wrote %d byte(s) (%d conversation(s), %d memory record(s)) to %s.\n",
				len(raw), len(res.Conversations), len(res.Memories), opts.OutputFile)
		}
		return nil
	}
	_, err = opts.IO.Out.Write(append(raw, '\n'))
	return err
}
