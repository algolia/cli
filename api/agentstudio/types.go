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
