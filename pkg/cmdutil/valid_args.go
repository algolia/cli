package cmdutil

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/cobra"
)

// IndexNames returns a function to list the index names from the given search client.
func IndexNames(clientF func() (*search.Client, error)) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
