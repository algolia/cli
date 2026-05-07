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

// Agent mirrors AgentWithVersionResponse. Config and Tools are kept as
// raw JSON; see docs/agents.md ("Pass-through bodies").
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

// PaginationMetadata is the standard paginated-response envelope.
type PaginationMetadata struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalCount int `json:"totalCount"`
	TotalPages int `json:"totalPages"`
}

// PaginatedAgentsResponse is the GET /1/agents response.
type PaginatedAgentsResponse struct {
	Data       []Agent            `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// ListAgentsParams configures GET /1/agents. Page/Limit at 0 = server default.
type ListAgentsParams struct {
	Page       int
	Limit      int
	ProviderID string
}

// ProviderName values mirror the backend's ProviderName enum. Kept as
// constants (not a typed enum) because the CLI passes the value through
// verbatim from user JSON.
const (
	ProviderNameOpenAI           = "openai"
	ProviderNameAzureOpenAI      = "azure_openai"
	ProviderNameGoogleGenAI      = "google_genai"
	ProviderNameDeepSeek         = "deepseek"
	ProviderNameOpenAICompatible = "openai_compatible"
	ProviderNameAnthropic        = "anthropic"
)

// AllProviderNames feeds help text and flag-validation lists.
var AllProviderNames = []string{
	ProviderNameOpenAI,
	ProviderNameAzureOpenAI,
	ProviderNameGoogleGenAI,
	ProviderNameDeepSeek,
	ProviderNameOpenAICompatible,
	ProviderNameAnthropic,
}

// Provider mirrors ProviderAuthenticationResponse. Input is raw JSON
// (discriminated union over ProviderName); see docs/agents.md.
type Provider struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	ProviderName string          `json:"providerName"`
	Input        json.RawMessage `json:"input"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
	LastUsedAt   *time.Time      `json:"lastUsedAt,omitempty"`
}

// PaginatedProvidersResponse is the GET /1/providers response.
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
type ApplicationConfig struct {
	MaxRetentionDays int `json:"maxRetentionDays"`
}

// Conversation mirrors ConversationBaseResponse (no messages — the
// lightweight shape used in list responses).
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

// PaginatedConversationsResponse is the list-conversations envelope.
type PaginatedConversationsResponse struct {
	Data       []Conversation     `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// ListConversationsParams configures GET /1/agents/{id}/conversations.
// FeedbackVote is *int because nil = no filter while 0 (downvote) is a
// meaningful value; backend silently drops the param unless
// IncludeFeedback=true.
type ListConversationsParams struct {
	Page            int
	Limit           int
	StartDate       string
	EndDate         string
	IncludeFeedback bool
	FeedbackVote    *int
}

// PurgeConversationsParams configures DELETE /1/agents/{id}/conversations.
// Backend rejects dateless purge — see docs/agents.md gotchas.
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

// SecretKey mirrors SecretKeyResponse. Value is sensitive — always mask
// unless the caller explicitly opts in.
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

// SecretKeyCreate is the POST body. AgentIDs is omitted when empty.
type SecretKeyCreate struct {
	Name     string   `json:"name"`
	AgentIDs []string `json:"agentIds,omitempty"`
}

// SecretKeyPatch is the PATCH body. Pointer fields: nil = leave unchanged,
// non-nil zero value = sent through (clears the field).
type SecretKeyPatch struct {
	Name     *string   `json:"name,omitempty"`
	AgentIDs *[]string `json:"agentIds,omitempty"`
}

// FeedbackCreate is the POST body for /1/feedback. Vote is 0 (downvote)
// or 1 (upvote); enforced at the CLI layer.
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
// are raw JSON (evolving schemas).
type UserDataResponse struct {
	Conversations []json.RawMessage `json:"conversations"`
	Memories      []json.RawMessage `json:"memories"`
}

// StatusResponse mirrors GET /status.
type StatusResponse map[string]*string

// ModelDefaults mirrors GET /1/providers/models/defaults.
type ModelDefaults map[string]string
