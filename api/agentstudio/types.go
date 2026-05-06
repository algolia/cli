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

// Conversation mirrors ConversationBaseResponse from the backend (the
// "lightweight, no messages" shape used in PaginatedConversationsResponse).
//
// Feedback and ConversationMetadata are kept as json.RawMessage:
//
//   - Feedback is `[]FeedbackResponse | null` and only populated when
//     the caller passes ?includeFeedback=true. The FeedbackResponse
//     schema has 7 fields with several anyOf nullables; the CLI does
//     not introspect them, it only forwards. Typing them eagerly would
//     be churn-prone for a feature `agents conversations list` only
//     exposes via passthrough.
//   - ConversationMetadata is currently a single nullable timestamp
//     ({cachedAt: …}) but the schema name signals room to grow.
//
// Both are emitted with `omitempty` so default `--output json` for
// list/get isn't bloated by trailing nulls.
type Conversation struct {
	ID                   string          `json:"id"`
	AgentID              string          `json:"agentId"`
	Title                *string         `json:"title,omitempty"`
	CreatedAt            time.Time       `json:"createdAt"`
	UpdatedAt            time.Time       `json:"updatedAt"`
	LastActivityAt       *time.Time      `json:"lastActivityAt,omitempty"`
	UserToken            *string         `json:"userToken,omitempty"`
	IsFromDashboard      bool            `json:"isFromDashboard"`
	MessageCount         int             `json:"messageCount"`
	TotalInputTokens     int             `json:"totalInputTokens"`
	TotalOutputTokens    int             `json:"totalOutputTokens"`
	TotalTokens          int             `json:"totalTokens"`
	Feedback             json.RawMessage `json:"feedback,omitempty"`
	ConversationMetadata json.RawMessage `json:"conversationMetadata,omitempty"`
}

// PaginatedConversationsResponse mirrors PaginatedConversationsResponse.
type PaginatedConversationsResponse struct {
	Data       []Conversation     `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// ListConversationsParams configures GET /1/agents/{id}/conversations.
//
// StartDate/EndDate are YYYY-MM-DD strings sent verbatim — the backend
// validates Pydantic-side (same passthrough convention as
// InvalidateAgentCache.before).
//
// IncludeFeedback toggles ?includeFeedback (backend default: false).
// FeedbackVote is a *int because nil = "no filter" while 0 ("downvote")
// is a meaningful filter value. Backend constraint: 0 <= vote <= 1, and
// the param is silently dropped unless includeFeedback=true is also set.
type ListConversationsParams struct {
	Page            int
	Limit           int
	StartDate       string
	EndDate         string
	IncludeFeedback bool
	FeedbackVote    *int
}

// PurgeConversationsParams configures DELETE /1/agents/{id}/conversations.
// Empty StartDate/EndDate = wipe everything (the CLI layer adds an
// explicit `--all` guardrail; this struct is the wire shape only).
type PurgeConversationsParams struct {
	StartDate string
	EndDate   string
}

// ExportConversationsParams configures GET /1/agents/{id}/conversations/export.
type ExportConversationsParams struct {
	StartDate string
	EndDate   string
}

// AllowedDomain mirrors AllowedDomainResponse.
type AllowedDomain struct {
	ID        string    `json:"id"`
	AppID     string    `json:"appId"`
	AgentID   string    `json:"agentId"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// AllowedDomainListResponse is the un-paginated list shape for domains.
type AllowedDomainListResponse struct {
	Domains []AllowedDomain `json:"domains"`
}

// SecretKey mirrors SecretKeyResponse. Value is the vended secret —
// always treat as sensitive.
type SecretKey struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Value      string     `json:"value"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	IsDefault  bool       `json:"isDefault"`
	AgentIDs   []string   `json:"agentIds"`
}

// PaginatedSecretKeysResponse is the standard paginated envelope.
type PaginatedSecretKeysResponse struct {
	Data       []SecretKey        `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// ListSecretKeysParams configures GET /1/secret-keys.
type ListSecretKeysParams struct {
	Page  int
	Limit int
}

// SecretKeyCreate is the POST body. AgentIDs is optional; when empty
// the field is omitted.
type SecretKeyCreate struct {
	Name     string   `json:"name"`
	AgentIDs []string `json:"agentIds,omitempty"`
}

// SecretKeyPatch is the PATCH body. Both fields are pointers so a
// nil value means "leave unchanged" while a zero value (empty string
// / empty slice) is sent through to the backend.
type SecretKeyPatch struct {
	Name     *string   `json:"name,omitempty"`
	AgentIDs *[]string `json:"agentIds,omitempty"`
}

// FeedbackCreate is the POST body for /1/feedback. Vote is 0
// (downvote) or 1 (upvote); enforced at the CLI layer.
type FeedbackCreate struct {
	MessageID string   `json:"messageId"`
	AgentID   string   `json:"agentId"`
	Vote      int      `json:"vote"`
	Tags      []string `json:"tags,omitempty"`
	Notes     string   `json:"notes,omitempty"`
}

// Feedback mirrors FeedbackResponse.
type Feedback struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agentId"`
	MessageID string    `json:"messageId"`
	Vote      int       `json:"vote"`
	Tags      []string  `json:"tags"`
	Notes     *string   `json:"notes"`
	Model     *string   `json:"model"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserDataResponse mirrors GET /1/user-data/{user_token}. Inner items
// are passed through as raw JSON because their schemas are evolving
// (conversation messages are a discriminated role union, memories
// have an unspecified shape).
type UserDataResponse struct {
	Conversations []json.RawMessage `json:"conversations"`
	Memories      []json.RawMessage `json:"memories"`
}
