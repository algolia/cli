package genai

import (
	"fmt"
	"time"
)

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Status  int    `json:"status,omitempty"`
	Details string `json:"details,omitempty"`
}

// APIError represents a more structured error with HTTP and API context
type APIError struct {
	StatusCode int
	Method     string
	Path       string
	Response   ErrorResponse
}

// Error implements the error interface for APIError
func (e *APIError) Error() string {
	if e.Response.Message != "" {
		return fmt.Sprintf("[%d] %s: %s %s", e.StatusCode, e.Response.Message, e.Method, e.Path)
	}
	return fmt.Sprintf("[%d] Error accessing: %s %s", e.StatusCode, e.Method, e.Path)
}

// DeleteResponse represents a successful deletion response
type DeleteResponse struct {
	Message string `json:"message"`
}

// DataSourceResponse represents a response when creating/updating a data source
type DataSourceResponse struct {
	ObjectID string `json:"objectID"`
}

// DataSource represents a data source in the GenAI API
type DataSource struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Config      struct {
		IndexName string `json:"indexName"`
		AppID     string `json:"appId"`
	} `json:"config"`
}

// DataSourceDetails represents a detailed data source returned by the get endpoint
type DataSourceDetails struct {
	Status          int       `json:"status"`
	Name            string    `json:"name"`
	Source          string    `json:"source"`
	Filters         string    `json:"filters,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	LinkedResponses int       `json:"linkedResponses"`
	ObjectID        string    `json:"objectID"`
}

// CreateDataSourceInput represents the input for creating a data source
type CreateDataSourceInput struct {
	Name     string `json:"name"`
	Source   string `json:"source"`
	Filters  string `json:"filters,omitempty"`
	ObjectID string `json:"objectID,omitempty"`
}

// UpdateDataSourceInput represents the input for updating a data source
type UpdateDataSourceInput struct {
	ObjectID string `json:"objectID"`
	Name     string `json:"name,omitempty"`
	Source   string `json:"source,omitempty"`
	Filters  string `json:"filters,omitempty"`
}

// DeleteDataSourcesInput represents the input for deleting data sources
type DeleteDataSourcesInput struct {
	ObjectIDs             []string `json:"objectIDs"`
	DeleteLinkedResponses bool     `json:"deleteLinkedResponses,omitempty"`
}

// ListDataSourcesResponse represents the response from listing data sources
type ListDataSourcesResponse struct {
	DataSources []DataSource `json:"dataSources"`
}

// PromptResponse represents a response when creating/updating a prompt
type PromptResponse struct {
	ObjectID string `json:"objectID"`
}

// Prompt represents a prompt in the GenAI API
type Prompt struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// PromptDetails represents a detailed prompt returned by the get endpoint
type PromptDetails struct {
	Status          int       `json:"status"`
	Name            string    `json:"name"`
	Instructions    string    `json:"instructions"`
	Tone            string    `json:"tone,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	LinkedResponses int       `json:"linkedResponses"`
	ObjectID        string    `json:"objectID"`
}

// CreatePromptInput represents the input for creating a prompt
type CreatePromptInput struct {
	Name         string `json:"name"`
	Instructions string `json:"instructions"`
	Tone         string `json:"tone,omitempty"`
	ObjectID     string `json:"objectID,omitempty"`
}

// UpdatePromptInput represents the input for updating a prompt
type UpdatePromptInput struct {
	ObjectID     string `json:"objectID"`
	Name         string `json:"name,omitempty"`
	Instructions string `json:"instructions,omitempty"`
	Tone         string `json:"tone,omitempty"`
}

// DeletePromptsInput represents the input for deleting prompts
type DeletePromptsInput struct {
	ObjectIDs             []string `json:"objectIDs"`
	DeleteLinkedResponses bool     `json:"deleteLinkedResponses,omitempty"`
}

// ListPromptsResponse represents the response from listing prompts
type ListPromptsResponse struct {
	Prompts []Prompt `json:"prompts"`
}

// GenerateResponseInput represents the input for generating a response
type GenerateResponseInput struct {
	Query                string   `json:"query,omitempty"`
	DataSourceID         string   `json:"dataSourceID"`
	PromptID             string   `json:"promptID"`
	LogRegion            string   `json:"logRegion"`
	ObjectID             string   `json:"objectID,omitempty"`
	NbHits               int      `json:"nbHits,omitempty"`
	AdditionalFilters    string   `json:"additionalFilters,omitempty"`
	WithObjectIDs        []string `json:"withObjectIDs,omitempty"`
	AttributesToRetrieve []string `json:"attributesToRetrieve,omitempty"`
	ConversationID       string   `json:"conversationID,omitempty"`
	Save                 bool     `json:"save,omitempty"`
	UseCache             bool     `json:"useCache,omitempty"`
	Origin               string   `json:"origin,omitempty"`
}

// GenerateResponseOutput represents the output from generating a response
type GenerateResponseOutput struct {
	ObjectID     string    `json:"objectID"`
	Response     string    `json:"response"`
	Query        string    `json:"query"`
	DataSourceID string    `json:"dataSourceID"`
	PromptID     string    `json:"promptID"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
}

// ResponseDetails represents a detailed response returned by the get endpoint
type ResponseDetails struct {
	Status            int       `json:"status"`
	Query             string    `json:"query"`
	DataSourceID      string    `json:"dataSourceID"`
	PromptID          string    `json:"promptID"`
	AdditionalFilters string    `json:"additionalFilters,omitempty"`
	Save              bool      `json:"save"`
	UseCache          bool      `json:"use_cache"`
	Origin            string    `json:"origin"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	ObjectID          string    `json:"objectID"`
	Response          string    `json:"response,omitempty"`
}

// ListResponsesResponse represents the response from listing responses
type ListResponsesResponse struct {
	Responses []GenerateResponseOutput `json:"responses"`
}

// DeleteResponsesInput represents the input for deleting responses
type DeleteResponsesInput struct {
	ObjectIDs []string `json:"objectIDs"`
}
