package cmdutil

import "github.com/algolia/algoliasearch-client-go/v3/algolia/search"

func IndexNames(client search.ClientInterface) ([]string, error) {
	res, err := client.ListIndices()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(res.Items))
	for _, index := range res.Items {
		names = append(names, index.Name)
	}

	return names, nil
}
