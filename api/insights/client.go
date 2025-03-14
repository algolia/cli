package insights

import (
	"encoding/json"
	"fmt"
	"time"

	algoliaInsights "github.com/algolia/algoliasearch-client-go/v4/algolia/insights"
	"github.com/algolia/algoliasearch-client-go/v4/algolia/transport"
	"github.com/algolia/cli/pkg/version"
)

// Client wraps the default Insights API client so that we can declare methods on it
type Client struct {
	*algoliaInsights.APIClient
}

// NewClient instantiates a new Insights API client
func NewClient(appID, apiKey string, region algoliaInsights.Region) (*Client, error) {
	// Get the default user agent
	userAgent, err := getUserAgentInfo(appID, apiKey, region, version.Version)
	if err != nil {
		return nil, err
	}
	if userAgent == "" {
		return nil, fmt.Errorf("user agent info must not be empty")
	}
	clientConfig := algoliaInsights.InsightsConfiguration{
		Configuration: transport.Configuration{
			AppID:     appID,
			ApiKey:    apiKey,
			UserAgent: userAgent,
		},
	}
	client, err := algoliaInsights.NewClientWithConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}

// GetEvents retrieves a number of events from the Algolia Insights API.
func (c *Client) GetEvents(startDate, endDate time.Time, limit int) (*EventsRes, error) {
	layout := "2006-01-02T15:04:05.000Z"
	params := map[string]any{
		"startDate": startDate.Format(layout),
		"endDate":   endDate.Format(layout),
		"limit":     limit,
	}
	res, err := c.CustomGet(c.NewApiCustomGetRequest("1/events").WithParameters(params))
	if err != nil {
		return nil, err
	}
	tmp, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	var eventsRes EventsRes
	err = json.Unmarshal(tmp, &eventsRes)
	if err != nil {
		return nil, err
	}

	return &eventsRes, err
}

// getUserAgentInfo returns the user agent string for the Insights client in the CLI
func getUserAgentInfo(
	appID string,
	apiKey string,
	region algoliaInsights.Region,
	appVersion string,
) (string, error) {
	client, err := algoliaInsights.NewClient(appID, apiKey, region)
	if err != nil {
		return "", err
	}

	return client.GetConfiguration().UserAgent + fmt.Sprintf("; Algolia CLI (%s)", appVersion), nil
}
