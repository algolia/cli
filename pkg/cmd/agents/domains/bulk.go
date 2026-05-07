package domains

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmd/agents/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type BulkInsertOptions struct {
	IO                *iostreams.IOStreams
	Ctx               context.Context
	AgentStudioClient func() (*agentstudio.Client, error)
	PrintFlags        *cmdutil.PrintFlags
	AgentID           string
	Domains           []string
	File              string
	DryRun            bool
}

func newBulkInsertCmd(f *cmdutil.Factory, runF func(*BulkInsertOptions) error) *cobra.Command {
	opts := &BulkInsertOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}
	cmd := &cobra.Command{
		Use:   "bulk-insert <agent-id> [--domain x ...] [-F file]",
		Short: "Add multiple allowed domains in one call",
		Long: heredoc.Doc(`
			Add multiple allowed domains in a single request. Provide
			values via repeated --domain flags or a file containing a
			JSON array of strings (use "-" for stdin).
		`),
		Example: heredoc.Doc(`
			$ algolia agents domains bulk-insert <agent-id> --domain a.test --domain b.test
			$ echo '["a.test","b.test"]' | algolia agents domains bulk-insert <agent-id> -F -
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			if opts.File != "" && len(opts.Domains) > 0 {
				return cmdutil.FlagErrorf("--domain and --file are mutually exclusive")
			}
			if opts.File == "" && len(opts.Domains) == 0 {
				return cmdutil.FlagErrorf("provide at least one --domain or --file")
			}
			if opts.File != "" {
				vals, err := readDomainsFromFile(opts.File, opts.IO)
				if err != nil {
					return err
				}
				opts.Domains = vals
			}
			if runF != nil {
				return runF(opts)
			}
			return runBulkInsertCmd(opts)
		},
	}
	cmd.Flags().StringSliceVar(&opts.Domains, "domain", nil, "Domain or pattern (repeatable)")
	cmd.Flags().
		StringVarP(&opts.File, "file", "F", "", "JSON file containing an array of strings (use \"-\" for stdin)")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Print what would be sent without calling the API")
	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runBulkInsertCmd(opts *BulkInsertOptions) error {
	if opts.DryRun {
		fmt.Fprintf(opts.IO.Out,
			"Dry run: would POST /1/agents/%s/allowed-domains/bulk\n  domains (%d): %s\n",
			opts.AgentID, len(opts.Domains), strings.Join(opts.Domains, ", "))
		return nil
	}
	client, err := opts.AgentStudioClient()
	if err != nil {
		return err
	}
	opts.IO.StartProgressIndicatorWithLabel("Bulk-inserting allowed domains")
	res, err := client.BulkInsertAllowedDomains(shared.OrBackground(opts.Ctx), opts.AgentID, opts.Domains)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	return opts.PrintFlags.Print(opts.IO, res)
}

type BulkDeleteOptions struct {
	IO                *iostreams.IOStreams
	Ctx               context.Context
	AgentStudioClient func() (*agentstudio.Client, error)
	AgentID           string
	DomainIDs         []string
	File              string
	DryRun            bool
	DoConfirm         bool
}

func newBulkDeleteCmd(f *cmdutil.Factory, runF func(*BulkDeleteOptions) error) *cobra.Command {
	opts := &BulkDeleteOptions{
		IO:                f.IOStreams,
		AgentStudioClient: f.AgentStudioClient,
	}
	var confirm bool
	cmd := &cobra.Command{
		Use:   "bulk-delete <agent-id> [--domain-id x ...] [-F file] [--confirm]",
		Short: "Remove multiple allowed domains by ID in one call",
		Long: heredoc.Doc(`
			Remove multiple allowed domains by ID. Provide IDs via
			repeated --domain-id flags or a file containing a JSON
			array of strings (use "-" for stdin).
		`),
		Args: validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.AgentID = args[0]
			opts.Ctx = cmd.Context()
			if opts.AgentID == "" {
				return cmdutil.FlagErrorf("agent-id must not be empty")
			}
			if opts.File != "" && len(opts.DomainIDs) > 0 {
				return cmdutil.FlagErrorf("--domain-id and --file are mutually exclusive")
			}
			if opts.File == "" && len(opts.DomainIDs) == 0 {
				return cmdutil.FlagErrorf("provide at least one --domain-id or --file")
			}
			if opts.File != "" {
				vals, err := readDomainsFromFile(opts.File, opts.IO)
				if err != nil {
					return err
				}
				opts.DomainIDs = vals
			}
			doConfirm, err := shared.ResolveConfirm(opts.IO, confirm, opts.DryRun)
			if err != nil {
				return err
			}
			opts.DoConfirm = doConfirm
			if runF != nil {
				return runF(opts)
			}
			return runBulkDeleteCmd(opts)
		},
	}
	cmd.Flags().StringSliceVar(&opts.DomainIDs, "domain-id", nil, "Domain ID (repeatable)")
	cmd.Flags().StringVarP(&opts.File, "file", "F", "", "JSON file containing an array of IDs (use \"-\" for stdin)")
	shared.AddConfirmFlag(cmd, &confirm)
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Print what would be sent without calling the API")
	return cmd
}

func runBulkDeleteCmd(opts *BulkDeleteOptions) error {
	if opts.DryRun {
		fmt.Fprintf(opts.IO.Out,
			"Dry run: would DELETE /1/agents/%s/allowed-domains/bulk\n  ids (%d): %s\n",
			opts.AgentID, len(opts.DomainIDs), strings.Join(opts.DomainIDs, ", "))
		return nil
	}
	if opts.DoConfirm {
		ok, err := shared.Confirm(fmt.Sprintf("Delete %d allowed domain(s)?", len(opts.DomainIDs)))
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

	// Backend returns 204 with no body on bulk delete, so it can't tell us
	// which IDs actually existed. Pre-fetch the list to classify
	// requested IDs as present-and-removed vs already-absent. One extra
	// GET per bulk op — acceptable for an admin-facing command, and the
	// signal saves users from "did my delete actually do anything?".
	var present, absent int
	if opts.IO.IsStdoutTTY() {
		opts.IO.StartProgressIndicatorWithLabel("Inspecting allowed domains")
		current, lerr := client.ListAllowedDomains(ctx, opts.AgentID)
		opts.IO.StopProgressIndicator()
		if lerr == nil {
			existing := make(map[string]struct{}, len(current.Domains))
			for _, d := range current.Domains {
				existing[d.ID] = struct{}{}
			}
			for _, id := range opts.DomainIDs {
				if _, ok := existing[id]; ok {
					present++
				} else {
					absent++
				}
			}
		}
	}

	opts.IO.StartProgressIndicatorWithLabel("Bulk-deleting allowed domains")
	err = client.BulkDeleteAllowedDomains(ctx, opts.AgentID, opts.DomainIDs)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.Out,
			"%s Bulk delete OK — requested %d (removed %d, already absent %d)\n",
			cs.SuccessIcon(), len(opts.DomainIDs), present, absent)
	}
	return nil
}

// readDomainsFromFile reads a JSON array of strings from file/stdin.
func readDomainsFromFile(file string, ios *iostreams.IOStreams) ([]string, error) {
	body, err := shared.ReadJSONFile(ios.In, file)
	if err != nil {
		return nil, err
	}
	var out []string
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, cmdutil.FlagErrorf("%s: expected a JSON array of strings", shared.SourceLabel(file))
	}
	if len(out) == 0 {
		return nil, cmdutil.FlagErrorf("%s: array is empty", shared.SourceLabel(file))
	}
	return out, nil
}
