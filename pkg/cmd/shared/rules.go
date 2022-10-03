package handler

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmd/rules/shared"
)

type RuleFlagsProvided struct {
	idProvided, enabledProvided, descriptionProvided, patternProvided, anchoringProvided, alternativeProvided,
	contextProvided, promoteObjectIdProvided, promoteObjectIdsProvided, promotePositionProvided, filterPromotesProvided,
	hideProvided, userDataProvided, queryProvided, facetProvided, scoreProvided, disjunctiveProvided, negativeProvided bool
}

func GetRuleFlagsProvided(cmd *cobra.Command) RuleFlagsProvided {
	return RuleFlagsProvided{
		// Options
		idProvided:          cmd.Flags().Changed("id"),
		enabledProvided:     cmd.Flags().Changed("enabled"),
		descriptionProvided: cmd.Flags().Changed("description"),
		// Condition
		patternProvided:     cmd.Flags().Changed("pattern"),
		anchoringProvided:   cmd.Flags().Changed("anchoring"),
		alternativeProvided: cmd.Flags().Changed("alternative"),
		contextProvided:     cmd.Flags().Changed("context"),
		// Consequence
		filterPromotesProvided: cmd.Flags().Changed("filter"),
		hideProvided:           cmd.Flags().Changed("hide"),
		userDataProvided:       cmd.Flags().Changed("userData"),
		// ConsequencePromote
		promoteObjectIdProvided:  cmd.Flags().Changed("promoteObjectId"),
		promoteObjectIdsProvided: cmd.Flags().Changed("promoteObjectIds"),
		promotePositionProvided:  cmd.Flags().Changed("promotePosition"),
		// Consequence Params
		queryProvided: cmd.Flags().Changed("query"),
		// Consequence Params Automatic Facet Filter
		facetProvided:       cmd.Flags().Changed("facet"),
		scoreProvided:       cmd.Flags().Changed("score"),
		disjunctiveProvided: cmd.Flags().Changed("disjunctive"),
		negativeProvided:    cmd.Flags().Changed("negative"),
	}
}

func ValidateRuleFlags(flags shared.RuleFlags, flagsProvided RuleFlagsProvided) error {
	if !flagsProvided.idProvided {
		return fmt.Errorf("a unique rule id is required")
	}

	if flagsProvided.filterPromotesProvided &&
		(!flagsProvided.promoteObjectIdProvided || !flagsProvided.promoteObjectIdsProvided || !flagsProvided.promotePositionProvided) {
		return fmt.Errorf("filterPromotes flag is only used in combination with the promote consequence (required: promoteObjectId, promoteObjectIds & promotePosition)")
	}

	if flagsProvided.anchoringProvided && flags.ConditionAnchoring != string(search.Is) &&
		flags.ConditionAnchoring != string(search.StartsWith) && flags.ConditionAnchoring != string(search.EndsWith) &&
		flags.ConditionAnchoring != string(search.Contains) {
		return fmt.Errorf("anchoring value should be one of is, startsWith, endsWith, contains")
	}

	if !flagsProvided.promoteObjectIdProvided && !flagsProvided.promoteObjectIdsProvided &&
		!flagsProvided.promotePositionProvided && !flagsProvided.hideProvided &&
		!flagsProvided.userDataProvided && !flagsProvided.queryProvided && !flagsProvided.facetProvided &&
		!flagsProvided.scoreProvided && !flagsProvided.disjunctiveProvided && !flagsProvided.negativeProvided {
		return fmt.Errorf("a consequence parameters is required")
	}

	return nil
}
