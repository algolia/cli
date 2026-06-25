package userdata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type GetOptions struct {
	IO                   *iostreams.IOStreams
	Ctx                  context.Context
	AgentStudioAPIClient func() (*agentStudio.APIClient, error)
	UserToken            string
	OutputFile           string
}

func newGetCmd(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO:                   f.IOStreams,
		AgentStudioAPIClient: f.AgentStudioAPIClient,
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
	client, err := opts.AgentStudioAPIClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Fetching user data")
	// Forward the backend payload verbatim — this is a GDPR data dump, so the
	// inner conversation/memory schemas are not pinned to a Go type.
	raw, err := shared.RawResponse(client.GetUserDataWithHTTPInfo(
		client.NewApiGetUserDataRequest(opts.UserToken),
		agentStudio.WithContext(shared.OrBackground(opts.Ctx)),
	))
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	// Count items for the human summary without pinning the inner schema.
	var counts struct {
		Conversations []json.RawMessage `json:"conversations"`
		Memories      []json.RawMessage `json:"memories"`
	}
	_ = json.Unmarshal(raw, &counts)

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, raw, "", "  "); err != nil {
		pretty.Reset()
		pretty.Write(raw)
	}
	out := pretty.Bytes()

	if opts.OutputFile != "" {
		if err := os.WriteFile(opts.OutputFile, out, 0o600); err != nil {
			return fmt.Errorf("write %s: %w", opts.OutputFile, err)
		}
		if opts.IO.IsStdoutTTY() {
			fmt.Fprintf(opts.IO.Out,
				"Wrote %d byte(s) (%d conversation(s), %d memory record(s)) to %s.\n",
				len(out), len(counts.Conversations), len(counts.Memories), opts.OutputFile)
		}
		return nil
	}
	_, err = opts.IO.Out.Write(append(out, '\n'))
	return err
}
