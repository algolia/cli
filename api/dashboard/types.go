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
	Name          string          `json:"name"`
	ApplicationID string          `json:"application_id"`
	APIKey        string          `json:"api_key"`
	Plan          ApplicationPlan `json:"plan"`
}

// ApplicationPlan is the plan applied to an application (attributes.plan).
// Label (e.g. "Grow Plus") matches a self-serve plan template's Name.
type ApplicationPlan struct {
	Name       string `json:"name"`
	Label      string `json:"label"`
	Version    int    `json:"version"`
	PayAsYouGo bool   `json:"pay_as_you_go"`
}

// Application is a flattened view of an Algolia application for CLI consumption.
type Application struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	APIKey    string `json:"api_key,omitempty"`
	PlanLabel string `json:"plan_label,omitempty"` // current plan label, e.g. "Grow Plus"
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

// UpdateApplicationRequest is the payload for PATCH /1/applications/{id}.
type UpdateApplicationRequest struct {
	Name string `json:"name"`
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

// APIError is returned for non-2xx dashboard responses. It carries the HTTP
// status so callers (and telemetry) can branch on it, keeping the message.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string { return e.Message }

func (e *APIError) HTTPStatusCode() int { return e.StatusCode }

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

// DashboardCrawlerUserData contains the user information from the crawler API
type DashboardCrawlerUserData struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	APIKey string `json:"apiKey"`
}

// DashboardCrawlerUserResponse is the JSON:API response from GET /1/crawler/user
type DashboardCrawlerUserResponse struct {
	Data DashboardCrawlerUserData `json:"data"`
}

type DashboardCrawlerErrorResponse struct {
	Errors []DashboardCrawlerError `json:"errors"`
}

type DashboardCrawlerError struct {
	Status string  `json:"status"`
	Title  string  `json:"title"`
	Detail *string `json:"detail"`
}

// toApplication flattens a JSON:API resource into a simple Application.
func (r *ApplicationResource) toApplication() Application {
	return Application{
		ID:        r.Attributes.ApplicationID,
		Name:      r.Attributes.Name,
		APIKey:    r.Attributes.APIKey,
		PlanLabel: r.Attributes.Plan.Label,
	}
}

// PlanTypeFree is the attributes.type value that identifies the free-tier plan
// template. The free plan's configuration.plan id is not fixed (it can be
// "build"), so the CLI keys off this type rather than a hard-coded id when it
// needs to map the user-facing "free" choice to a concrete plan.
const PlanTypeFree = "free"

// PlanTemplatesResponse is the JSON:API response from GET /1/plan-templates/self-serve.
type PlanTemplatesResponse struct {
	Data []PlanTemplateResource `json:"data"`
}

// PlanTemplateResource is a JSON:API resource wrapper for a self-serve plan template.
type PlanTemplateResource struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes PlanTemplateAttributes `json:"attributes"`
}

// PlanTemplateAttributes contains the actual plan template fields.
type PlanTemplateAttributes struct {
	Name          string                    `json:"name"`
	Description   string                    `json:"description"`
	Type          string                    `json:"type"`               // e.g. "free" or "freeform"
	Freeform      string                    `json:"freeform,omitempty"` // pricing string for paygo plans
	Configuration PlanTemplateConfiguration `json:"configuration"`
}

// PlanTemplateConfiguration holds the plan identifier and the terms text.
type PlanTemplateConfiguration struct {
	Plan        string `json:"plan"`
	AcceptTerms string `json:"accept_terms"`
}

// Plan is a flattened, CLI-friendly view of a self-serve plan template.
type Plan struct {
	ID          string `json:"id"`           // configuration.plan (e.g. "build", "grow", "grow-plus")
	Name        string `json:"name"`         // human-readable name (e.g. "Grow Plus")
	Description string `json:"description"`  // short description
	Type        string `json:"type"`         // "free" or "freeform"
	Price       string `json:"price"`        // "Free" for the free plan, otherwise the freeform pricing string
	AcceptTerms string `json:"accept_terms"` // ToS text shown before changing the plan
}

// IsFree reports whether the plan is the free-tier plan (no payment method required).
func (p Plan) IsFree() bool {
	return p.Type == PlanTypeFree
}

// toPlan flattens a JSON:API plan template resource into a Plan.
func (r *PlanTemplateResource) toPlan() Plan {
	// The free plan template has no "freeform" pricing field, so present a
	// friendly "Free" instead of an empty string.
	price := r.Attributes.Freeform
	if r.Attributes.Type == PlanTypeFree || price == "" {
		price = "Free"
	}
	return Plan{
		ID:          r.Attributes.Configuration.Plan,
		Name:        r.Attributes.Name,
		Description: r.Attributes.Description,
		Type:        r.Attributes.Type,
		Price:       price,
		AcceptTerms: r.Attributes.Configuration.AcceptTerms,
	}
}

// ChangePlanRequest is the payload for PATCH /1/applications/{id}/plan/self-serve.
//
// Every request body in this client is a plain JSON object (see
// CreateApplicationRequest and UpdateApplicationRequest) rather than a JSON:API
// data.attributes envelope, so we mirror that existing convention here.
type ChangePlanRequest struct {
	Plan string `json:"plan"`
}

// DashboardUser is a flattened view of the authenticated user's account info
// from GET /1/user, exposing only what plan changes need.
type DashboardUser struct {
	HasPaymentMethod bool `json:"has_payment_method"`
}

// userResponse is a forgiving decoder for GET /1/user. The exact shape of this
// endpoint is not documented for the CLI, so we accept "has_payment_method"
// either at the top level or nested under a JSON:API data.attributes object,
// and use whichever is present.
type userResponse struct {
	HasPaymentMethod *bool `json:"has_payment_method"`
	Data             *struct {
		Attributes struct {
			HasPaymentMethod *bool `json:"has_payment_method"`
		} `json:"attributes"`
	} `json:"data"`
}

func (r *userResponse) toUser() DashboardUser {
	u := DashboardUser{}

	switch {
	case r.HasPaymentMethod != nil:
		u.HasPaymentMethod = *r.HasPaymentMethod
	case r.Data != nil && r.Data.Attributes.HasPaymentMethod != nil:
		u.HasPaymentMethod = *r.Data.Attributes.HasPaymentMethod
	}

	return u
}
