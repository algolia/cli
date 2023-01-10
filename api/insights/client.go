package insights

import (
	"fmt"
	"net/http"
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/call"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/compression"
	_insights "github.com/algolia/algoliasearch-client-go/v3/algolia/insights"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/transport"
)

// Client provides methods to interact with the Algolia Insights API.
type Client struct {
	appID     string
	transport *transport.Transport
}

// NewClient instantiates a new client able to interact with the Algolia
// Insights API.
func NewClient(appID, apiKey string) *Client {
	return NewClientWithConfig(
		_insights.Configuration{
			AppID:  appID,
			APIKey: apiKey,
		},
	)
}

// NewClientWithConfig instantiates a new client able to interact with the
// Algolia Insights API.
func NewClientWithConfig(config _insights.Configuration) *Client {
	var hosts []*transport.StatefulHost

	if config.Hosts == nil {
		hosts = defaultHosts(config.Region)
	} else {
		for _, h := range config.Hosts {
			hosts = append(hosts, transport.NewStatefulHost(h, call.IsReadWrite))
		}
	}

	return &Client{
		appID: config.AppID,
		transport: transport.New(
			hosts,
			config.Requester,
			config.AppID,
			config.APIKey,
			config.ReadTimeout,
			config.WriteTimeout,
			config.Headers,
			config.ExtraUserAgent,
			compression.None,
		),
	}
}

// FetchEvents retrieves events from the Algolia Insights API.
func (c *Client) FetchEvents(startDate, endDate time.Time, limit int) (EventsRes, error) {
	var res EventsRes
	path := fmt.Sprintf("/1/events?startDate=%s&endDate=%s&limit=%d", startDate.Format("2006-01-02T15:04:05.000Z"), endDate.Format("2006-01-02T15:04:05.000Z"), limit)
	err := c.transport.Request(&res, http.MethodGet, path, nil, call.Read, nil)
	return res, err
}
