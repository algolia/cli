package agentstudio

import "encoding/json"

// Agent is the full agent representation returned by the API.
// Complex nested fields (config, tools) are kept as raw JSON so the CLI does
// not have to model every provider/tool union variant.
type Agent struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  *string         `json:"description,omitempty"`
	Status       string          `json:"status"`
	ProviderID   *string         `json:"providerId,omitempty"`
	Model        *string         `json:"model,omitempty"`
	Instructions string          `json:"instructions"`
	SystemPrompt *string         `json:"systemPrompt,omitempty"`
	Config       json.RawMessage `json:"config,omitempty"`
	Tools        json.RawMessage `json:"tools,omitempty"`
	TemplateType *string         `json:"templateType,omitempty"`
	CreatedAt    string          `json:"createdAt"`
	UpdatedAt    string          `json:"updatedAt"`
	LastUsedAt   *string         `json:"lastUsedAt,omitempty"`
}

// PaginationMetadata mirrors the spec.
type PaginationMetadata struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalCount int `json:"totalCount"`
	TotalPages int `json:"totalPages"`
}

// PaginatedAgentsResponse is the response of GET /1/agents.
type PaginatedAgentsResponse struct {
	Data       []Agent            `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// AgentConfigCreate is the request body for POST /1/agents.
// Free-form / union-typed fields are passed as raw JSON.
type AgentConfigCreate struct {
	Name         string          `json:"name"`
	Description  *string         `json:"description,omitempty"`
	ProviderID   *string         `json:"providerId,omitempty"`
	Model        *string         `json:"model,omitempty"`
	Instructions string          `json:"instructions"`
	SystemPrompt *string         `json:"systemPrompt,omitempty"`
	TemplateType *string         `json:"templateType,omitempty"`
	Config       json.RawMessage `json:"config,omitempty"`
	Tools        json.RawMessage `json:"tools,omitempty"`
}

// AgentConfigUpdate is the request body for PATCH /1/agents/{id}. All fields
// are optional.
type AgentConfigUpdate struct {
	Name         *string         `json:"name,omitempty"`
	Description  *string         `json:"description,omitempty"`
	ProviderID   *string         `json:"providerId,omitempty"`
	Model        *string         `json:"model,omitempty"`
	Instructions *string         `json:"instructions,omitempty"`
	SystemPrompt *string         `json:"systemPrompt,omitempty"`
	TemplateType *string         `json:"templateType,omitempty"`
	Config       json.RawMessage `json:"config,omitempty"`
	Tools        json.RawMessage `json:"tools,omitempty"`
}

// AgentCompletionRequest is the body of POST /1/agents/{id}/completions.
type AgentCompletionRequest struct {
	ID            *string         `json:"id,omitempty"`
	Configuration json.RawMessage `json:"configuration,omitempty"`
	Messages      json.RawMessage `json:"messages,omitempty"`
	Algolia       json.RawMessage `json:"algolia,omitempty"`
	ToolApprovals json.RawMessage `json:"toolApprovals,omitempty"`
}

// Conversation is a flattened view used by both list and get endpoints.
// Messages is only populated by ConversationFullResponse-shaped payloads.
type Conversation struct {
	ID                   string          `json:"id"`
	AgentID              string          `json:"agentId"`
	Title                *string         `json:"title,omitempty"`
	CreatedAt            string          `json:"createdAt"`
	UpdatedAt            string          `json:"updatedAt"`
	LastActivityAt       *string         `json:"lastActivityAt,omitempty"`
	UserToken            *string         `json:"userToken,omitempty"`
	IsFromDashboard      *bool           `json:"isFromDashboard,omitempty"`
	MessageCount         *int            `json:"messageCount,omitempty"`
	TotalInputTokens     *int            `json:"totalInputTokens,omitempty"`
	TotalOutputTokens    *int            `json:"totalOutputTokens,omitempty"`
	TotalTokens          *int            `json:"totalTokens,omitempty"`
	ConversationMetadata json.RawMessage `json:"conversationMetadata,omitempty"`
	Feedback             json.RawMessage `json:"feedback,omitempty"`
	Messages             json.RawMessage `json:"messages,omitempty"`
}

// PaginatedConversationsResponse is the response of GET /1/agents/{id}/conversations.
type PaginatedConversationsResponse struct {
	Data       []Conversation     `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// ErrResponse models the FastAPI-style validation envelope. `Detail` may be a
// string (auth/permission errors) or an array of structured error entries
// (validation errors).
type ErrResponse struct {
	Detail json.RawMessage `json:"detail,omitempty"`
}
