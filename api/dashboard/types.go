package dashboard

import "errors"

// OAuthTokenResponse is the response from POST /2/oauth/token.
type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"`
	User         *User  `json:"user,omitempty"`
}

// User represents the authenticated user from the OAuth token response.
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// ApplicationResource is a JSON:API resource wrapper for an application.
type ApplicationResource struct {
	ID         string                `json:"id"`
	Type       string                `json:"type"`
	Attributes ApplicationAttributes `json:"attributes"`
}

// ApplicationAttributes contains the actual application fields.
type ApplicationAttributes struct {
	Name          string `json:"name"`
	ApplicationID string `json:"application_id"`
	APIKey        string `json:"api_key"`
}

// Application is a flattened view of an Algolia application for CLI consumption.
type Application struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	APIKey string `json:"api_key,omitempty"`
}

// PaginationMeta contains page-based pagination metadata.
type PaginationMeta struct {
	TotalCount  int `json:"total_count"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
}

// PaginationLinks contains pagination URLs.
type PaginationLinks struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
}

// ApplicationsResponse is the JSON:API response from GET /1/applications.
type ApplicationsResponse struct {
	Data  []ApplicationResource `json:"data"`
	Meta  PaginationMeta        `json:"meta"`
	Links PaginationLinks       `json:"links"`
}

// SingleApplicationResponse is the JSON:API response from GET /1/application/:id.
type SingleApplicationResponse struct {
	Data ApplicationResource `json:"data"`
}

// CreateApplicationRequest is the payload for POST /1/applications.
type CreateApplicationRequest struct {
	RegionCode string `json:"region_code"`
	Name       string `json:"name"`
}

// Region represents a hosting region from GET /1/hosting/regions.
type Region struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// RegionsResponse is the response from GET /1/hosting/regions.
type RegionsResponse struct {
	RegionCodes []Region `json:"region_codes"`
}

// ErrSessionExpired is returned when an API call gets a 401 Unauthorized.
var ErrSessionExpired = errors.New("session expired")

// ErrClusterUnavailable is returned when a region has no available cluster.
type ErrClusterUnavailable struct {
	Region  string
	Message string
}

func (e *ErrClusterUnavailable) Error() string {
	return e.Message
}

// OAuthErrorResponse is the error format from the OAuth endpoints.
type OAuthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// CreateAPIKeyRequest is the payload for POST /1/applications/{application_id}/api-keys.
type CreateAPIKeyRequest struct {
	ACL         []string `json:"acl"`
	Description string   `json:"description"`
}

// APIKeyResource is a JSON:API resource wrapper for an API key.
type APIKeyResource struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes APIKeyAttributes `json:"attributes"`
}

// APIKeyAttributes contains the actual API key fields.
type APIKeyAttributes struct {
	Value string `json:"value"`
}

// CreateAPIKeyResponse is the JSON:API response from POST /1/applications/{application_id}/api-keys.
type CreateAPIKeyResponse struct {
	Data APIKeyResource `json:"data"`
}

// toApplication flattens a JSON:API resource into a simple Application.
func (r *ApplicationResource) toApplication() Application {
	return Application{
		ID:     r.Attributes.ApplicationID,
		Name:   r.Attributes.Name,
		APIKey: r.Attributes.APIKey,
	}
}
