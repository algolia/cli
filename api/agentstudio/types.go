package agentstudio

import (
	"encoding/json"
	"time"
)

// AgentStatus is the lifecycle state of an agent.
type AgentStatus string

const (
	StatusDraft     AgentStatus = "draft"
	StatusPublished AgentStatus = "published"
)

// Agent mirrors AgentWithVersionResponse from the Agent Studio backend.
//
// Wire format is camelCase (the backend uses Pydantic CamelModel, see
// common/models/agent_config.py:AgentWithVersionResponse).
//
// Config and Tools are kept as raw JSON: the schemas are large, evolve
// frequently, and are not yet stable enough to mirror in Go. CLI commands
// that mutate them (create/update) accept user-supplied JSON files instead
// of building structs in code.
type Agent struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  *string         `json:"description,omitempty"`
	Status       AgentStatus     `json:"status"`
	ProviderID   *string         `json:"providerId,omitempty"`
	Model        *string         `json:"model,omitempty"`
	Instructions string          `json:"instructions"`
	SystemPrompt *string         `json:"systemPrompt,omitempty"`
	Config       json.RawMessage `json:"config,omitempty"`
	Tools        json.RawMessage `json:"tools,omitempty"`
	TemplateType *string         `json:"templateType,omitempty"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    *time.Time      `json:"updatedAt,omitempty"`
	LastUsedAt   *time.Time      `json:"lastUsedAt,omitempty"`
}

// PaginationMetadata mirrors common/models/pagination_metadata.py.
type PaginationMetadata struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalCount int `json:"totalCount"`
	TotalPages int `json:"totalPages"`
}

// PaginatedAgentsResponse mirrors PaginatedAgentsResponse from the backend.
type PaginatedAgentsResponse struct {
	Data       []Agent            `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// ListAgentsParams configures GET /1/agents.
type ListAgentsParams struct {
	// Page is 1-indexed; 0 means "use server default".
	Page int
	// Limit is items per page; 0 means "use server default" (currently 10).
	Limit int
	// ProviderID filters by provider authentication ID. Empty means no filter.
	ProviderID string
}

// ProviderName values mirror the backend's ProviderName enum
// (rag/models/provider.py). These are exposed as constants rather than
// a typed enum because:
//   - The CLI passes the value through verbatim from user JSON.
//   - The backend is the source of truth for which providers are
//     supported on a given deployment (feature gates can disable
//     subsets per app).
//
// Used only for documentation, validation hints, and the
// `agents providers models` discoverability flow.
const (
	ProviderNameOpenAI           = "openai"
	ProviderNameAzureOpenAI      = "azure_openai"
	ProviderNameGoogleGenAI      = "google_genai"
	ProviderNameDeepSeek         = "deepseek"
	ProviderNameOpenAICompatible = "openai_compatible"
	ProviderNameAnthropic        = "anthropic"
)

// AllProviderNames is exported so the cmd layer can build help text and
// flag-validation lists from a single source of truth. Keep in sync
// with the backend's ProviderName enum.
var AllProviderNames = []string{
	ProviderNameOpenAI,
	ProviderNameAzureOpenAI,
	ProviderNameGoogleGenAI,
	ProviderNameDeepSeek,
	ProviderNameOpenAICompatible,
	ProviderNameAnthropic,
}

// Provider mirrors ProviderAuthenticationResponse from the backend.
//
// Input is kept as json.RawMessage (same rationale as Agent.Config /
// Agent.Tools): the input shape is a discriminated union over
// ProviderName with deeply-validated per-variant fields (apiKey,
// baseUrl, azureEndpoint, azureDeployment, defaultModel, ...) that
// evolves as new providers land. Mirroring the variants in Go would
// lie about parity and force a CLI release on every backend bump.
//
// The CLI surfaces a few well-known fields opportunistically (the
// presence of "apiKey" triggers masking) but otherwise pretty-prints
// Input as JSON.
type Provider struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	ProviderName string          `json:"providerName"`
	Input        json.RawMessage `json:"input"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
	LastUsedAt   *time.Time      `json:"lastUsedAt,omitempty"`
}

// PaginatedProvidersResponse mirrors PaginatedProviderAuthenticationsResponse.
type PaginatedProvidersResponse struct {
	Data       []Provider         `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// ListProvidersParams configures GET /1/providers.
type ListProvidersParams struct {
	Page  int
	Limit int
}

// ApplicationConfig mirrors ApplicationConfigResponse.
//
// Single-field today (maxRetentionDays). Kept as a struct rather than
// `int` so the wire shape can grow without breaking callers.
type ApplicationConfig struct {
	MaxRetentionDays int `json:"maxRetentionDays"`
}
