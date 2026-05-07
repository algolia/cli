package upsert

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	algoliaComposition "github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/spf13/cobra"

	compinternal "github.com/algolia/cli/pkg/cmd/compositions/internal"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

// UpsertOptions holds dependencies and flags for the rules upsert command.
type UpsertOptions struct {
	Config            config.IConfig
	IO                *iostreams.IOStreams
	CompositionClient func() (*algoliaComposition.APIClient, error)
	CompositionID     string
	ObjectID          string
	File              string
	PrintFlags        *cmdutil.PrintFlags
}

// NewUpsertCmd returns the `compositions rules upsert` command.
func NewUpsertCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &UpsertOptions{
		IO:                f.IOStreams,
		Config:            f.Config,
		CompositionClient: f.CompositionClient,
		PrintFlags:        cmdutil.NewPrintFlags().WithDefaultOutput("json"),
	}

	cmd := &cobra.Command{
		Use:   "upsert <composition-id> <rule-id>",
		Short: "Create or update a composition rule",
		Args:  validators.ExactArgsWithMsg(2, "compositions rules upsert requires a <composition-id> and a <rule-id> argument."),
		Annotations: map[string]string{
			"acls": "editSettings",
		},
		Example: heredoc.Doc(`
			# Upsert a rule from a JSON file
			$ algolia compositions rules upsert my-comp rule-1 --file rule.json

			# Upsert from stdin
			$ cat rule.json | algolia compositions rules upsert my-comp rule-1 --file -
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CompositionID = args[0]
			opts.ObjectID = args[1]
			return runUpsertCmd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "JSON file path (use - for stdin)")
	_ = cmd.MarkFlagRequired("file")

	opts.PrintFlags.AddFlags(cmd)
	return cmd
}

func runUpsertCmd(opts *UpsertOptions) error {
	raw, err := cmdutil.ReadFile(opts.File, opts.IO.In)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var rule algoliaComposition.CompositionRule
	if err := json.Unmarshal(raw, &rule); err != nil {
		return fmt.Errorf("parsing rule JSON: %w", err)
	}

	client, err := opts.CompositionClient()
	if err != nil {
		return err
	}

	p, err := opts.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}

	opts.IO.StartProgressIndicatorWithLabel("Upserting rule")

	res, err := client.PutCompositionRule(client.NewApiPutCompositionRuleRequest(opts.CompositionID, opts.ObjectID, &rule))
	if err != nil {
		opts.IO.StopProgressIndicator()
		return err
	}

	opts.IO.StopProgressIndicator()

	if err := compinternal.WaitForTask(opts.IO, client, opts.CompositionID, res.TaskID, compinternal.PollInterval, compinternal.Timeout); err != nil {
		return err
	}

	return p.Print(opts.IO, res)
}
