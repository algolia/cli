package shared

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

type RuleFlags struct {
	// Options
	RuleID          string
	RuleEnabled     bool
	RuleDescription string
	// TODO: RuleValidity array of tuple
	// Condition
	ConditionPattern     string
	ConditionAnchoring   string
	ConditionAlternative bool
	ConditionContext     string
	// Consequence
	ConsequenceFilterPromotes bool
	ConsequenceHide           []string
	ConsequenceUserData       string
	// ConsequencePromote
	ConsequencePromoteObjectID  string
	ConsequencePromoteObjectIDs []string
	ConsequencePromotePosition  int8
	// Consequence Params
	ConsequenceParamsQuery string
	// Consequence Params Automatic Facet Filter
	ConsequenceParamsAutomaticFacetFilterFacet       string
	ConsequenceParamsAutomaticFacetFilterScore       int8
	ConsequenceParamsAutomaticFacetFilterDisjunctive bool
	ConsequenceParamsAutomaticFacetFilterNegative    bool
	// TODO: Consequence Params Automatic Optional Facet Filter
	// ConsequenceParamsAutomaticFacetOptionalFilterFacet       string
	// ConsequenceParamsAutomaticFacetOptionalFilterScore       int
	// ConsequenceParamsAutomaticFacetOptionalFilterDisjunctive bool
	// ConsequenceParamsAutomaticFacetOptionalFilterNegative    bool
}

func FlagsToRule(flags RuleFlags) (search.Rule, error) {
	// Base options
	rule := search.Rule{
		ObjectID: flags.RuleID,
		Enabled:  opt.Enabled(flags.RuleEnabled),
	}
	if flags.RuleDescription != "" {
		rule.Description = flags.RuleDescription
	}

	// Conditions
	// TODO: handle multiple conditions
	condition := search.RuleCondition{}
	if flags.ConditionPattern != "" {
		condition.Pattern = flags.ConditionPattern
	}
	if flags.ConditionAnchoring != "" {
		switch flags.ConditionAnchoring {
		case string(search.Is):
			condition.Anchoring = search.Is
		case string(search.StartsWith):
			condition.Anchoring = search.StartsWith
		case string(search.EndsWith):
			condition.Anchoring = search.EndsWith
		case string(search.Contains):
			condition.Anchoring = search.Contains
		default:
			return rule, fmt.Errorf("wrong consequence anchoring value")
		}
	}
	if flags.ConditionAlternative {
		condition.Alternatives = search.AlternativesEnabled()
	}
	if flags.ConditionContext != "" {
		condition.Context = flags.ConditionContext
	}

	// Consequences

	rule.Conditions = []search.RuleCondition{condition}

	// rule := search.Rule{
	// 	ObjectID:   flags.RuleID,
	// 	Conditions: []search.RuleCondition{condition},
	// 	Consequence: search.RuleConsequence{
	// 		Params: &search.RuleParams{
	// 			QueryParams: search.QueryParams{
	// 				Filters: opt.Filters("category = 1"),
	// 			},
	// 		},
	// 	},
	// 	Enabled: opt.Enabled(flags.RuleEnabled), // Optionally, to disable the rule

	// }

	return rule, nil
}
