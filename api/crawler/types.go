package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
)

// APIError is an error returned by the Crawler API, carrying the HTTP status.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}

// ErrResponse is a Crawler API error response.
type ErrResponse struct {
	Err Err `json:"error"`
}

// Err is a Crawler API error.
type Err struct {
	Message string         `json:"message"`
	Code    string         `json:"code"`
	Errors  []LabeledError `json:"errors,omitempty"`
}

// LabeledError is a Crawler API labeled error.
type LabeledError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Label   string `json:"label"`
}

// Crawler is a Crawler.
type Crawler struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name"`
	Running    bool   `json:"running,omitempty"`
	Reindexing bool   `json:"reindexing,omitempty"`
	Blocked    bool   `json:"blocked,omitempty"`

	BlockingTaskID string `json:"blockingTaskId,omitempty"`
	BlockingError  string `json:"blockingError,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`

	LastReindexStartedAt time.Time `json:"lastReindexStartedAt,omitempty"`
	LastReindexEndedAt   time.Time `json:"lastReindexEndedAt,omitempty"`

	Config *Config `json:"config,omitempty"`
}

// Config is a Crawler configuration.
type Config struct {
	AppID       string   `json:"appId,omitempty"`
	APIKey      string   `json:"apiKey,omitempty"`
	IndexPrefix string   `json:"indexPrefix,omitempty"`
	Schedule    string   `json:"schedule,omitempty"`
	StartUrls   []string `json:"startUrls,omitempty"`
	Sitemaps    []string `json:"sitemaps,omitempty"`

	ExclusionPatterns []string          `json:"exclusionPatterns,omitempty"`
	IgnoreQueryParams []string          `json:"ignoreQueryParams,omitempty"`
	RenderJavaScript  *RenderJavaScript `json:"renderJavaScript,omitempty"`
	RateLimit         int               `json:"rateLimit,omitempty"`
	ExtraUrls         []string          `json:"extraUrls,omitempty"`
	MaxDepth          int               `json:"maxDepth,omitempty"`
	MaxURLs           int               `json:"maxUrls,omitempty"`

	IgnoreRobotsTxtRules bool `json:"ignoreRobotsTxtRules,omitempty"`
	IgnoreNoIndex        bool `json:"ignoreNoIndex,omitempty"`
	IgnoreNoFollowTo     bool `json:"ignoreNoFollowTo,omitempty"`
	IgnoreCanonicalTo    bool `json:"ignoreCanonicalTo,omitempty"`

	SaveBackup           bool                             `json:"saveBackup,omitempty"`
	InitialIndexSettings map[string]*search.IndexSettings `json:"initialIndexSettings,omitempty"`

	Actions []*Action `json:"actions,omitempty"`
}

type RenderJavaScript struct {
	Enabled bool

	raw json.RawMessage
}

func (r *RenderJavaScript) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil
	}

	switch trimmed[0] {
	case 't', 'f':
		var enabled bool
		if err := json.Unmarshal(trimmed, &enabled); err != nil {
			return err
		}
		r.Enabled = enabled
	case '[':
		r.Enabled = true
	case '{':
		var object struct {
			Enabled *bool `json:"enabled"`
		}
		if err := json.Unmarshal(trimmed, &object); err != nil {
			return err
		}
		r.Enabled = object.Enabled == nil || *object.Enabled
	default:
		return fmt.Errorf(
			"renderJavaScript must be a boolean, an array of URL patterns or an object, got %s",
			trimmed,
		)
	}

	r.raw = bytes.Clone(trimmed)

	return nil
}

func (r *RenderJavaScript) MarshalJSON() ([]byte, error) {
	if len(r.raw) > 0 {
		return bytes.Clone(r.raw), nil
	}

	return json.Marshal(r.Enabled)
}

// Action is a Crawler configuration action.
type Action struct {
	IndexName        string          `json:"indexName"`
	PathsToMatch     []string        `json:"pathsToMatch"`
	SelectorsToMatch []string        `json:"selectorsToMatch,omitempty"`
	FileTypesToMatch []string        `json:"fileTypesToMatch,omitempty"`
	RecordExtractor  RecordExtractor `json:"recordExtractor"`
}

// RecordExtractor is a Crawler configuration record extractor.
type RecordExtractor struct {
	Type   string `json:"__type"`
	Source string `json:"source"`
}

// TestResponse is the response from the crawler crawlers/{id}/test endpoint.
type TestResponse struct {
	StartDate    time.Time   `json:"startDate"`
	EndDate      time.Time   `json:"endDate"`
	Logs         interface{} `json:"logs,omitempty"`
	Records      interface{} `json:"records,omitempty"`
	Links        []string    `json:"links,omitempty"`
	ExternalData interface{} `json:"externalData,omitempty"`
	Error        *Err        `json:"error,omitempty"`
}

// TaskIDResponse is the response when a task is created.
type TaskIDResponse struct {
	TaskID string `json:"taskId"`
}

// StatsResponse is the response from the crawler crawlers/{id}/stats/urls endpoint.
type StatsResponse struct {
	Count int `json:"count"`
	Data  []struct {
		Reason   string `json:"reason"`
		Status   string `json:"status"`
		Category string `json:"category"`
		Readable string `json:"readable"`
		Count    int    `json:"count"`
	} `json:"data"`
}

// CrawlerListItem is a crawler list item.
type CrawlerListItem struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// CrawlersResponse is the response from the crawler crawlers endpoint.
type CrawlersResponse struct {
	Items []*CrawlerListItem `json:"items"`

	// Pagination
	Page         int `json:"page"`
	ItemsPerPage int `json:"itemsPerPage"`
	Total        int `json:"total"`
}
