package provisionning

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/google/jsonapi"
)

const (
	// DefaultBaseURL is the default base URL for the Algolia Provisionning API.
	DefaultBaseURL = "https://api.dashboard.algolia.com/1/"

	// Supported plans for the provisionning API.
	PlanV8Build = "v8.5-plg-build"
	PlanV8Grow  = "v8.5-plg-grow"
)

// Client provides methods to interact with the Algolia Provisionning API.
type Client struct {
	token  string
	client *http.Client
}

// NewClient returns a new Provisionning API client.
func NewClient(token string) *Client {
	return &Client{
		client: http.DefaultClient,
		token:  token,
	}
}

// NewClientWithHTTPClient returns a new Provisionning API client with a custom HTTP client.
func NewClientWithHTTPClient(token string, client *http.Client) *Client {
	return &Client{
		token:  token,
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

	var errResponse jsonapi.ErrorsPayload
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
			return fmt.Errorf("HTTP %d", resp.StatusCode)
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, errResponse.Errors[0].Detail)
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

		r = io.NopCloser(bytes.NewReader(b))
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

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.api+json")

	// Add URL params
	values := req.URL.Query()
	for k, v := range urlParams {
		values.Set(k, v)
	}
	req.URL.RawQuery = values.Encode()

	return req, err
}

// unmarshalTo unmarshals an HTTP response body to a given interface.
func unmarshalTo(r *http.Response, t interface{}) error {
	// Copy the response body to a buffer to be able to read it multiple times
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	r.Body.Close()
	reader := bytes.NewReader(body)

	// Test if the response is a JSON API response
	data := make(map[string]interface{})
	err = json.NewDecoder(reader).Decode(&data)
	if err != nil {
		return err
	}

	if _, ok := data["data"]; !ok {
		// The response might be a non-JSON API response (ex: hosting regions endpoint)
		reader.Seek(0, 0)
		return json.NewDecoder(reader).Decode(&t)
	}

	// If struct have `PaginatedResponse` type, unmarshal as a list
	if reflect.TypeOf(t).Elem().Kind() == reflect.Slice {
		reader.Seek(0, 0)
		list, err := jsonapi.UnmarshalManyPayload(reader, reflect.TypeOf(t).Elem())
		if err != nil {
			return err
		}
		for _, item := range list {
			reflect.ValueOf(t).Elem().Set(reflect.Append(reflect.ValueOf(t).Elem(), reflect.ValueOf(item).Elem()))
		}
		return nil
	} else {
		reader.Seek(0, 0)
		err = jsonapi.UnmarshalPayload(reader, t)
		if err != nil {
			return err
		}
	}

	return nil
}
