package provisionning

import (
	"fmt"
	"net/http"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

// CreateAPIKey creates a new API key.
func (c *Client) CreateAPIKey(appID string, apiKey search.Key) (*APIKey, error) {
	var res APIKey
	path := fmt.Sprintf("applications/%s/api-keys", appID)

	err := c.request(&res, http.MethodPost, path, apiKey, nil)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
