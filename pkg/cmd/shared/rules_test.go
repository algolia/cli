package handler

import (
	"testing"

	"github.com/algolia/cli/pkg/cmd/rules/shared"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateRuleFlags(t *testing.T) {
	tests := []struct {
		name              string
		ruleFlags         shared.RuleFlags
		ruleFlagsProvided RuleFlagsProvided
		wantsErr          bool
		wantsErrMsg       string
	}{
		{
			name: "wrong condition anchoring",
			ruleFlags: shared.RuleFlags{
				RuleID:             "1",
				ConditionAnchoring: "abc",
			},
			ruleFlagsProvided: RuleFlagsProvided{
				// Options
				idProvided:          true,
				enabledProvided:     false,
				descriptionProvided: false,
				// Condition
				patternProvided:     false,
				anchoringProvided:   true,
				alternativeProvided: false,
				contextProvided:     false,
				// Consequence
				filterPromotesProvided: false,
				hideProvided:           false,
				userDataProvided:       false,
				// ConsequencePromote
				promoteObjectIdProvided:  false,
				promoteObjectIdsProvided: false,
				promotePositionProvided:  false,
				// Consequence Params
				queryProvided: false,
				// Consequence Params Automatic Facet Filter
				facetProvided:       false,
				scoreProvided:       false,
				disjunctiveProvided: false,
				negativeProvided:    false,
			},
			wantsErr:    true,
			wantsErrMsg: "anchoring value should be one of is, startsWith, endsWith, contains",
		},
		{
			name: "filterPromotes without promote consequence",
			ruleFlagsProvided: RuleFlagsProvided{
				// Options
				idProvided:          true,
				enabledProvided:     false,
				descriptionProvided: false,
				// Condition
				patternProvided:     false,
				anchoringProvided:   false,
				alternativeProvided: false,
				contextProvided:     false,
				// Consequence
				filterPromotesProvided: true,
				hideProvided:           false,
				userDataProvided:       false,
				// ConsequencePromote
				promoteObjectIdProvided:  true,
				promoteObjectIdsProvided: true,
				promotePositionProvided:  true,
				// Consequence Params
				queryProvided: false,
				// Consequence Params Automatic Facet Filter
				facetProvided:       false,
				scoreProvided:       false,
				disjunctiveProvided: false,
				negativeProvided:    false,
			},
			wantsErr: false,
		},
		{
			name: "No rule id",
			ruleFlagsProvided: RuleFlagsProvided{
				// Options
				idProvided:          false,
				enabledProvided:     true,
				descriptionProvided: true,
				// Condition
				patternProvided:     true,
				anchoringProvided:   true,
				alternativeProvided: true,
				contextProvided:     true,
				// Consequence
				filterPromotesProvided: true,
				hideProvided:           true,
				userDataProvided:       true,
				// ConsequencePromote
				promoteObjectIdProvided:  true,
				promoteObjectIdsProvided: true,
				promotePositionProvided:  true,
				// Consequence Params
				queryProvided: true,
				// Consequence Params Automatic Facet Filter
				facetProvided:       true,
				scoreProvided:       true,
				disjunctiveProvided: true,
				negativeProvided:    true,
			},
			wantsErr:    true,
			wantsErrMsg: "a unique rule id is required",
		},
		{
			name: "No consequence",
			ruleFlagsProvided: RuleFlagsProvided{
				// Options
				idProvided:          true,
				enabledProvided:     false,
				descriptionProvided: false,
				// Condition
				patternProvided:     false,
				anchoringProvided:   false,
				alternativeProvided: false,
				contextProvided:     false,
				// Consequence
				filterPromotesProvided: false,
				hideProvided:           false,
				userDataProvided:       false,
				// ConsequencePromote
				promoteObjectIdProvided:  false,
				promoteObjectIdsProvided: false,
				promotePositionProvided:  false,
				// Consequence Params
				queryProvided: false,
				// Consequence Params Automatic Facet Filter
				facetProvided:       false,
				scoreProvided:       false,
				disjunctiveProvided: false,
				negativeProvided:    false,
			},
			wantsErr:    true,
			wantsErrMsg: "a consequence parameters is required",
		},
		{
			name: "filterPromotes without full promote consequence",
			ruleFlagsProvided: RuleFlagsProvided{
				// Options
				idProvided:          true,
				enabledProvided:     false,
				descriptionProvided: false,
				// Condition
				patternProvided:     false,
				anchoringProvided:   false,
				alternativeProvided: false,
				contextProvided:     false,
				// Consequence
				filterPromotesProvided: true,
				hideProvided:           false,
				userDataProvided:       false,
				// ConsequencePromote
				promoteObjectIdProvided:  false,
				promoteObjectIdsProvided: true,
				promotePositionProvided:  false,
				// Consequence Params
				queryProvided: false,
				// Consequence Params Automatic Facet Filter
				facetProvided:       false,
				scoreProvided:       false,
				disjunctiveProvided: false,
				negativeProvided:    false,
			},
			wantsErr:    true,
			wantsErrMsg: "filterPromotes flag is only used in combination with the promote consequence (required: promoteObjectId, promoteObjectIds & promotePosition)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRuleFlags(tt.ruleFlags, tt.ruleFlagsProvided)

			if tt.wantsErr {
				assert.EqualError(t, err, tt.wantsErrMsg)
				return
			}
		})
	}
}
