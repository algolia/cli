package internal

import (
	"errors"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/composition"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/interactive"
)

// ResolveCompositionID returns the id from args, or prompts a selector over
// existing compositions when the argument is omitted on a terminal.
func ResolveCompositionID(
	args []string,
	canPrompt bool,
	prompter interactive.Prompter,
	clientF func() (*composition.APIClient, error),
) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}
	if !canPrompt {
		return "", cmdutil.FlagErrorf("a <composition-id> argument is required")
	}
	return SelectCompositionID(prompter, clientF)
}

// ResolveCompositionAndRuleID resolves a composition id and a rule id from args,
// prompting selectors for whichever is omitted on a terminal. The rule list is
// scoped to the resolved composition.
func ResolveCompositionAndRuleID(
	args []string,
	canPrompt bool,
	prompter interactive.Prompter,
	clientF func() (*composition.APIClient, error),
) (string, string, error) {
	compositionID, err := ResolveCompositionID(args, canPrompt, prompter, clientF)
	if err != nil {
		return "", "", err
	}
	if len(args) >= 2 {
		return compositionID, args[1], nil
	}
	if !canPrompt {
		return "", "", cmdutil.FlagErrorf("a <rule-id> argument is required")
	}
	ruleID, err := SelectRuleID(prompter, clientF, compositionID)
	if err != nil {
		return "", "", err
	}
	return compositionID, ruleID, nil
}

// SelectRuleID lists the rules of a composition and prompts the user to pick one.
func SelectRuleID(prompter interactive.Prompter, clientF func() (*composition.APIClient, error), compositionID string) (string, error) {
	client, err := clientF()
	if err != nil {
		return "", err
	}
	res, err := client.SearchCompositionRules(client.NewApiSearchCompositionRulesRequest(compositionID))
	if err != nil {
		return "", err
	}
	if len(res.Hits) == 0 {
		return "", errors.New("no rules found for composition " + compositionID)
	}
	ids := make([]string, len(res.Hits))
	for i, hit := range res.Hits {
		ids[i] = hit.ObjectID
	}
	return prompter.Select("Rule", ids)
}

// SelectCompositionID lists existing compositions and prompts the user to pick
// one, returning its object id. Used when a command is called without the
// composition-id argument on a terminal.
func SelectCompositionID(prompter interactive.Prompter, clientF func() (*composition.APIClient, error)) (string, error) {
	client, err := clientF()
	if err != nil {
		return "", err
	}
	res, err := client.ListCompositions(client.NewApiListCompositionsRequest())
	if err != nil {
		return "", err
	}
	if len(res.Items) == 0 {
		return "", errors.New("no compositions found")
	}
	ids := make([]string, len(res.Items))
	for i, comp := range res.Items {
		ids[i] = comp.ObjectID
	}
	return prompter.Select("Composition", ids)
}
