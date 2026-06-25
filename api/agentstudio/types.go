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
// raw JSON; see docs/agents.md ("Pass-through bodies"). Returned by the
// local DuplicateAgent call.
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
// (discriminated union over ProviderName); see docs/agents.md. Returned by
// the local UpdateProvider call.
type Provider struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	ProviderName string          `json:"providerName"`
	Input        json.RawMessage `json:"input"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
	LastUsedAt   *time.Time      `json:"lastUsedAt,omitempty"`
}

// StatusResponse mirrors GET /status.
type StatusResponse map[string]*string

// ModelDefaults mirrors GET /1/providers/models/defaults.
type ModelDefaults map[string]string
