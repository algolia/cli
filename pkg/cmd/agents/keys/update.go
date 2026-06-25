package keys

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	agentStudio "github.com/algolia/algoliasearch-client-go/v4/algolia/agent-studio"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type UpdateOptions struct {
	IO                   *iostreams.IOStreams
	Ctx                  context.Context
	AgentStudioAPIClient func() (*agentStudio.APIClient, error)
	PrintFlags           *cmdutil.PrintFlags
	ID                   string
	Name                 string
	AgentIDs             []string
	NameSet              bool
	AgentIDsSet          bool
	ShowSecret           bool
}

func newUpdateCmd(f *cmdutil.Factory, runF func(*UpdateOptions) error) *cobra.Command {
	opts := &UpdateOptions{
		IO:                   f.IOStreams,
		AgentStudioAPIClient: f.AgentStudioAPIClient,
		PrintFlags:           cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "update <id> [--name <name>] [--agent-id <id> ...]",
		Short: "Update a secret key's name and/or agent scope (admin key required)",
		Long: heredoc.Doc(`
			Update a secret key. Pass --name to rename, repeated
			--agent-id to set the agent allowlist (replaces — not
			merges — the existing list); pass --agent-id="" once to
			clear the list.
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			opts.Ctx = cmd.Context()
			opts.NameSet = cmd.Flags().Changed("name")
			opts.AgentIDsSet = cmd.Flags().Changed("agent-id")
			if opts.ID == "" {
				return cmdutil.FlagErrorf("id must not be empty")
			}
			if !opts.NameSet && !opts.AgentIDsSet {
				return cmdutil.FlagErrorf("provide --name or --agent-id (nothing to update)")
			}
			if runF != nil {
				return runF(opts)
			}
			return runUpdateCmd(opts)
		},
	}
	cmd.Flags().StringVar(&opts.Name, "name", "", "New name (max 128)")
	cmd.Flags().
		StringSliceVar(&opts.AgentIDs, "agent-id", nil, "Replace the agent allowlist (repeatable; pass --agent-id=\"\" to clear)")
	cmd.Flags().
		BoolVar(&opts.ShowSecret, "show-secret", false, "Reveal raw key value in the response (default redacted as ***)")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runUpdateCmd(opts *UpdateOptions) error {
	patch := &agentStudio.SecretKeyPatch{}
	if opts.NameSet {
		n := opts.Name
		patch.Name.Set(&n)
	}
	if opts.AgentIDsSet {
		// Non-nil empty slice clears the allowlist: SecretKeyPatch only omits
		// AgentIds when it is nil, so an empty slice still serializes as [].
		ids := make([]string, 0, len(opts.AgentIDs))
		for _, v := range opts.AgentIDs {
			if v != "" {
				ids = append(ids, v)
			}
		}
		patch.AgentIds = ids
	}
	client, err := opts.AgentStudioAPIClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Updating secret key")
	k, err := client.UpdateSecretKey(
		client.NewApiUpdateSecretKeyRequest(opts.ID, patch),
		agentStudio.WithContext(shared.OrBackground(opts.Ctx)),
	)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	masked := maskKey(*k, opts.ShowSecret)
	return opts.PrintFlags.Print(opts.IO, &masked)
}
