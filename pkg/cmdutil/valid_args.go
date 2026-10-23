package cmdutil

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/composition"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/algolia/cli/api/crawler"
	"github.com/spf13/cobra"
)

// IndexNames returns a function to list the index names from the given search client.
func IndexNames(
	clientF func() (*search.APIClient, error),
) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := clientF()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		res, err := client.ListIndices(client.NewApiListIndicesRequest())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		names := make([]string, 0, len(res.Items))
		for _, index := range res.Items {
			names = append(names, index.Name)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}

// CompositionIDs returns a function to complete composition object IDs from the given composition client.
// It completes only the first positional argument (the composition id).
func CompositionIDs(
	clientF func() (*composition.APIClient, error),
) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		client, err := clientF()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		res, err := client.ListCompositions(client.NewApiListCompositionsRequest())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		ids := make([]string, 0, len(res.Items))
		for _, comp := range res.Items {
			ids = append(ids, comp.ObjectID)
		}
		return ids, cobra.ShellCompDirectiveNoFileComp
	}
}

// CompositionAndRuleIDs completes the composition id on the first argument and,
// once that is set, the rule ids of that composition on the second.
func CompositionAndRuleIDs(
	clientF func() (*composition.APIClient, error),
) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := clientF()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		switch len(args) {
		case 0:
			res, err := client.ListCompositions(client.NewApiListCompositionsRequest())
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			ids := make([]string, 0, len(res.Items))
			for _, comp := range res.Items {
				ids = append(ids, comp.ObjectID)
			}
			return ids, cobra.ShellCompDirectiveNoFileComp
		case 1:
			res, err := client.SearchCompositionRules(client.NewApiSearchCompositionRulesRequest(args[0]))
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			ids := make([]string, 0, len(res.Hits))
			for _, hit := range res.Hits {
				ids = append(ids, hit.ObjectID)
			}
			return ids, cobra.ShellCompDirectiveNoFileComp
		default:
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
	}
}

// CrawlerIDs returns a function to list the crawler IDs from the given crawler client.
func CrawlerIDs(
	clientF func() (*crawler.Client, error),
) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := clientF()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		items, err := client.ListAll("", "")
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		names := make([]string, 0, len(items))
		for _, crawler := range items {
			names = append(names, fmt.Sprintf("%s\t%s", crawler.ID, crawler.Name))
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}
