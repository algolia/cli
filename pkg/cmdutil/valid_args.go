package cmdutil

import (
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	v4 "github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/api/crawler"
)

// IndexNames returns a function to list the index names from the given search client.
func V4_IndexNames(
	clientF func() (*v4.APIClient, error),
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

// IndexNames returns a function to list the index names from the given search client.
func IndexNames(
	clientF func() (*search.Client, error),
) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := clientF()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		res, err := client.ListIndices()
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
