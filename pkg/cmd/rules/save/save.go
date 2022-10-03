package save

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/rules/shared"
	handler "github.com/algolia/cli/pkg/cmd/shared"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/pkg/validators"
)

type SaveOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams

	SearchClient func() (*search.Client, error)

	Indice            string
	ForwardToReplicas bool
	Rule              search.Rule
	SuccessMessage    string
}

// NewSaveCmd creates and returns a save command for rules
func NewSaveCmd(f *cmdutil.Factory, runF func(*SaveOptions) error) *cobra.Command {
	opts := &SaveOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		SearchClient: f.SearchClient,
	}

	flags := &shared.RuleFlags{}

	cmd := &cobra.Command{
		Use:               "save <index> --id <id> TODO: fill usage",
		Args:              validators.ExactArgs(1),
		ValidArgsFunction: cmdutil.IndexNames(opts.SearchClient),
		Short:             "Save a rule to the given index",
		Aliases:           []string{"create", "edit"},
		Long: heredoc.Doc(`
			This command save a rule to the specified index.
			If the rule doesn't exist yet, a new one is created.
		`),
		Example: heredoc.Doc(`
			TODO: write doc
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Indice = args[0]

			flagsHandler := &handler.RuleHandler{
				Flags: flags,
				Cmd:   cmd,
			}

			err := handler.HandleFlags(flagsHandler, opts.IO.CanPrompt())
			if err != nil {
				return err
			}

			rule, err := shared.FlagsToRule(*flags)
			if err != nil {
				return err
			}
			// Correct flags are passed
			opts.Rule = rule

			// err, successMessage := GetSuccessMessage(*flags, opts.Indice)
			// if err != nil {
			// 	return err
			// }
			// opts.SuccessMessage = fmt.Sprintf("%s %s", f.IOStreams.ColorScheme().SuccessIcon(), successMessage)

			// if runF != nil {
			// 	return runF(opts)
			// }

			return runSaveCmd(opts)
		},
	}

	// Common
	cmd.Flags().BoolVarP(&opts.ForwardToReplicas, "forward-to-replicas", "f", false, "Forward the save request to the replicas")
	// Options
	cmd.Flags().StringVarP(&flags.RuleID, "id", "i", "", "objectID of the rule to save")
	cmd.Flags().BoolVarP(&flags.RuleEnabled, "enabled", "e", true, "Whether the Rule is enabled")
	cmd.Flags().StringVarP(&flags.RuleDescription, "description", "d", "", "Description of the rule to save")
	// Condition
	cmd.Flags().StringVarP(&flags.ConditionPattern, "pattern", "n", "", "Condition pattern of the rule to save")
	cmd.Flags().StringVarP(&flags.ConditionAnchoring, "anchoring", "a", "", "Condition anchoring of the rule to save")
	cmd.Flags().BoolVarP(&flags.ConditionAlternative, "alternative", "l", false, "Whether the pattern matches on plurals, synonyms, and typos")
	cmd.Flags().StringVarP(&flags.ConditionContext, "context", "c", "", "Condition context of the rule to save")
	// Consequence
	cmd.Flags().BoolVarP(&flags.ConsequenceFilterPromotes, "filterPromotes", "t", false, "Consequence filter promotes of the rule to save")
	cmd.Flags().StringSliceVarP(&flags.ConsequenceHide, "hide", "b", nil, "Objects ids to hide from hits of the rule to save")
	cmd.Flags().StringVarP(&flags.ConsequenceUserData, "userData", "u", "", "Custom JSON object that will be appended to the userData")
	// Consequence Promote
	cmd.Flags().StringVar(&flags.ConsequencePromoteObjectID, "promoteObjectId", "", "Consequence promote object id of the rule to save")
	cmd.Flags().StringSliceVar(&flags.ConsequencePromoteObjectIDs, "promoteObjectIds", nil, "Consequence promote object ids of the rule to save")
	cmd.Flags().Int8Var(&flags.ConsequencePromotePosition, "promotePosition", 0, "Consequence promote position of the rule to save")
	// Consequence Params
	cmd.Flags().StringVarP(&flags.ConsequenceParamsQuery, "query", "q", "", "Consequence params query of the rule to save")
	// Consequence Params Automatic Facet Filter
	cmd.Flags().StringVar(&flags.ConsequenceParamsAutomaticFacetFilterFacet, "facet", "", "Attribute to filter on: this must match a facet placeholder in the Rule's pattern")
	cmd.Flags().Int8Var(&flags.ConsequenceParamsAutomaticFacetFilterScore, "score", 1, "Score for the filter: typically used for optional or disjunctive filters")
	cmd.Flags().BoolVar(&flags.ConsequenceParamsAutomaticFacetFilterDisjunctive, "disjunctive", false, "Whether the filter is disjunctive (true) or conjunctive (false)")
	cmd.Flags().BoolVar(&flags.ConsequenceParamsAutomaticFacetFilterNegative, "negative", false, "Whether the match is inverted (true) or not (false)")

	return cmd
}

func runSaveCmd(opts *SaveOptions) error {
	client, err := opts.SearchClient()
	if err != nil {
		return err
	}

	indice := client.InitIndex(opts.Indice)
	forwardToReplicas := opt.ForwardToReplicas(opts.ForwardToReplicas)

	_, err = indice.SaveRule(opts.Rule, forwardToReplicas)
	if err != nil {
		err = fmt.Errorf("failed to save rule: %w", err)
		return err
	}

	if opts.IO.IsStdoutTTY() {
		fmt.Fprint(opts.IO.Out, opts.SuccessMessage)
	}

	return nil
}
