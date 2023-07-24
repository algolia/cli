package crawler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	// DefaultBaseURL is the default base URL for the Algolia Crawler API.
	DefaultBaseURL = "https://crawler.algolia.com/api/1/"
)

// Client provides methods to interact with the Algolia Crawler API.
type Client struct {
	UserID string
	APIKey string

	client *http.Client
}

// NewClient returns a new Crawler API client.
func NewClient(userID, apiKey string) *Client {
	return &Client{
		UserID: userID,
		APIKey: apiKey,
		client: http.DefaultClient,
	}
}

// NewClientWithHTTPClient returns a new Crawler API client with a custom HTTP client.
func NewClientWithHTTPClient(userID, apiKey string, client *http.Client) *Client {
	return &Client{
		UserID: userID,
		APIKey: apiKey,
		client: client,
	}
}

// Request sends an HTTP request and returns an HTTP response.
// It unmarshals the response body to the given interface.
func (c *Client) request(res interface{}, method string, path string, body interface{}, urlParams map[string]string) error {
	r, err := c.buildRequest(method, path, body, urlParams)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		var errResp ErrResponse
		if err := unmarshalTo(resp, &errResp); err != nil {
			return err
		}

		if errResp.Err.Errors != nil {
			var errs []string
			for _, e := range errResp.Err.Errors {
				errs = append(errs, e.Message)
			}
			return fmt.Errorf("%s: %s", errResp.Err.Message, errs)
		}

		return errors.New(errResp.Err.Message)
	}

	if res != nil {
		if err := unmarshalTo(resp, res); err != nil {
			return err
		}
	}

	return nil
}

// buildRequestWithoutBody builds an HTTP request without a body.
func (c *Client) buildRequestWithoutBody(method, url string) (*http.Request, error) {
	return http.NewRequest(method, url, nil)
}

// buildRequestWithBody builds an HTTP request with a body.
func (c *Client) buildRequestWithBody(method, url string, body interface{}) (*http.Request, error) {
	var r io.ReadCloser
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		r = ioutil.NopCloser(bytes.NewReader(b))
	}

	return http.NewRequest(method, url, r)
}

// buildRequest builds an HTTP request.
func (c *Client) buildRequest(method, path string, body interface{}, urlParams map[string]string) (req *http.Request, err error) {
	url := DefaultBaseURL + path

	if body == nil {
		req, err = c.buildRequestWithoutBody(method, url)
	} else {
		req, err = c.buildRequestWithBody(method, url, body)
	}

	req.SetBasicAuth(c.UserID, c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Add URL params
	values := req.URL.Query()
	for k, v := range urlParams {
		values.Set(k, v)
	}
	req.URL.RawQuery = values.Encode()

	return req, err
}

// unmarshalTo unmarshals an HTTP response body to a given interface.
func unmarshalTo(r *http.Response, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// Create creates a new Crawler.
// It returns the Crawler ID if successful.
func (c *Client) Create(name string, config Config) (string, error) {
	var res struct {
		ID string `json:"id"`
	}
	path := "crawlers"

	crawler := &Crawler{
		Name:   name,
		Config: &config,
	}

	err := c.request(&res, http.MethodPost, path, crawler, nil)
	if err != nil {
		return "", err
	}

	return res.ID, nil
}

// List lists Crawlers.
func (c *Client) List(itemsPerPage, page int, name, appID string) (*CrawlersResponse, error) {
	var res CrawlersResponse
	path := "crawlers"
	params := map[string]string{
		"itemsPerPage": fmt.Sprintf("%d", itemsPerPage),
		"page":         fmt.Sprintf("%d", page),
	}

	if name != "" {
		params["name"] = name
	}
	if appID != "" {
		params["appId"] = appID
	}

	err := c.request(&res, http.MethodGet, path, nil, params)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// ListAll lists all Crawlers
func (c *Client) ListAll(name, appID string) ([]*CrawlerListItem, error) {
	var crawlers []*CrawlerListItem

	while := true
	page := 0
	for while {
		page++
		res, err := c.List(20, page, name, appID)
		if err != nil {
			return nil, err
		}

		crawlers = append(crawlers, res.Items...)

		if len(crawlers) >= res.Total {
			while = false
		}
	}

	return crawlers, nil
}

// Get gets a Crawler.
func (c *Client) Get(crawlerID string, withConfig bool) (*Crawler, error) {
	var res Crawler
	path := fmt.Sprintf("crawlers/%s", crawlerID)
	params := map[string]string{
		"withConfig": fmt.Sprintf("%t", withConfig),
	}

	err := c.request(&res, http.MethodGet, path, nil, params)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Run runs a Crawler.
// It returns the Task ID if successful.
func (c *Client) Run(crawlerID string) (string, error) {
	var res TaskIDResponse
	path := fmt.Sprintf("crawlers/%s/run", crawlerID)

	err := c.request(&res, http.MethodPost, path, nil, nil)
	if err != nil {
		return "", err
	}

	return res.TaskID, nil
}

// Pause pauses a Crawler.
// It returns the Task ID if successful.
func (c *Client) Pause(crawlerID string) (string, error) {
	var res TaskIDResponse
	path := fmt.Sprintf("crawlers/%s/pause", crawlerID)

	err := c.request(&res, http.MethodPost, path, nil, nil)
	if err != nil {
		return "", err
	}

	return res.TaskID, nil
}

// Reindex reindexes a Crawler.
// It returns the Task ID if successful.
func (c *Client) Reindex(crawlerID string) (string, error) {
	var res TaskIDResponse
	path := fmt.Sprintf("crawlers/%s/reindex", crawlerID)

	err := c.request(&res, http.MethodPost, path, nil, nil)
	if err != nil {
		return "", err
	}

	return res.TaskID, nil
}

// Stats gets the stats of a Crawler.
func (c *Client) Stats(crawlerID string) (*StatsResponse, error) {
	var res StatsResponse
	path := fmt.Sprintf("crawlers/%s/stats/urls", crawlerID)

	err := c.request(&res, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CrawlURLs crawls the specified URLs on the specified Crawler.
// It returns the Task ID if successful.
func (c *Client) CrawlURLs(crawlerID string, URLs []string, save, saveSpecified bool) (string, error) {
	var res TaskIDResponse
	path := fmt.Sprintf("crawlers/%s/urls/crawl", crawlerID)

	var body interface{}
	if saveSpecified {
		body = struct {
			URLs []string `json:"urls"`
			Save bool     `json:"save"`
		}{URLs, save}
	} else {
		body = struct {
			URL []string `json:"urls"`
		}{URLs}
	}

	err := c.request(&res, http.MethodPost, path, body, nil)
	if err != nil {
		return "", err
	}

	return res.TaskID, nil
}

// Test tests an URL on the specified Crawler.
func (c *Client) Test(crawlerID, URL string, config *Config) (*TestResponse, error) {
	var res TestResponse
	path := fmt.Sprintf("crawlers/%s/test", crawlerID)

	var body interface{}
	if config != nil {
		body = struct {
			URL    string  `json:"url"`
			Config *Config `json:"config"`
		}{URL, config}
	} else {
		body = struct {
			URL string `json:"url"`
		}{URL}
	}

	err := c.request(&res, http.MethodPost, path, body, nil)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// CancelTask cancels a blocking task.
func (c *Client) CancelTask(crawlerID, taskID string) error {
	path := fmt.Sprintf("crawlers/%s/tasks/%s/cancel", crawlerID, taskID)

	err := c.request(nil, http.MethodPost, path, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
