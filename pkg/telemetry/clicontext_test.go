package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func fakeEnv(vars map[string]string) func(string) string {
	return func(key string) string { return vars[key] }
}

func TestDetectCLIContext_KnownAgents(t *testing.T) {
	cases := []struct {
		envVar string
		want   string
	}{
		{"CLAUDECODE", "agent:claude-code"},
		{"CLAUDE_CODE", "agent:claude-code"},
		{"CLAUDE_CODE_SESSION_ID", "agent:claude-code"},
		{"CURSOR_AGENT", "agent:cursor"},
		{"CODEX_THREAD_ID", "agent:codex"},
		{"CODEX_SANDBOX", "agent:codex"},
		{"CODEX_CI", "agent:codex"},
		{"GEMINI_CLI", "agent:gemini"},
		{"COPILOT_CLI", "agent:github-copilot"},
		{"OPENCODE", "agent:opencode"},
		{"OPENCODE_CLIENT", "agent:opencode"},
		{"AMP_CURRENT_THREAD_ID", "agent:amp"},
		{"AUGMENT_AGENT", "agent:augment"},
		{"GOOSE_TERMINAL", "agent:goose"},
		{"CLINE_ACTIVE", "agent:cline"},
		{"ANTIGRAVITY_AGENT", "agent:antigravity"},
		{"PI_CODING_AGENT", "agent:pi"},
		{"KIRO_AGENT_PATH", "agent:kiro"},
	}
	for _, c := range cases {
		t.Run(c.envVar, func(t *testing.T) {
			got := detectCLIContext(fakeEnv(map[string]string{c.envVar: "1"}), true, true, false)
			assert.Equal(t, c.want, got)
		})
	}
}

func TestDetectCLIContext_GenericAIAgentVar(t *testing.T) {
	got := detectCLIContext(fakeEnv(map[string]string{"AI_AGENT": "SomeAgent"}), true, true, false)
	assert.Equal(t, "agent:someagent", got)
}

func TestDetectCLIContext_GenericAgentVar(t *testing.T) {
	got := detectCLIContext(fakeEnv(map[string]string{"AGENT": "goose"}), true, true, false)
	assert.Equal(t, "agent:goose", got)
}

func TestDetectCLIContext_AIAgentWinsOverAgent(t *testing.T) {
	env := fakeEnv(map[string]string{"AI_AGENT": "v0", "AGENT": "goose"})
	assert.Equal(t, "agent:v0", detectCLIContext(env, true, true, false))
}

func TestDetectCLIContext_WhitespaceOnlyGenericVar(t *testing.T) {
	got := detectCLIContext(fakeEnv(map[string]string{"AI_AGENT": " "}), true, true, false)
	assert.Equal(t, "human", got)
}

func TestDetectCLIContext_VersionStampedGenericVar(t *testing.T) {
	got := detectCLIContext(fakeEnv(map[string]string{"AI_AGENT": "claude-code_2-1-210_agent"}), true, true, false)
	assert.Equal(t, "agent:claude-code", got)
}

func TestDetectCLIContext_LeadingUnderscoreGenericVar(t *testing.T) {
	got := detectCLIContext(fakeEnv(map[string]string{"AI_AGENT": "_foo"}), true, true, false)
	assert.Equal(t, "human", got)
}

func TestDetectCLIContext_NumericGenericVar(t *testing.T) {
	got := detectCLIContext(fakeEnv(map[string]string{"AGENT": "1"}), true, true, false)
	assert.Equal(t, "human", got)
}

func TestDetectCLIContext_BooleanishGenericVar(t *testing.T) {
	got := detectCLIContext(fakeEnv(map[string]string{"AGENT": "true"}), true, true, false)
	assert.Equal(t, "human", got)
}

func TestDetectCLIContext_OverlongGenericVar(t *testing.T) {
	got := detectCLIContext(fakeEnv(map[string]string{"AI_AGENT": "a-very-long-agent-name-exceeding-32-characters"}), true, true, false)
	assert.Equal(t, "human", got)
}

func TestDetectCLIContext_NamedVarWinsOverGeneric(t *testing.T) {
	env := fakeEnv(map[string]string{"CLAUDECODE": "1", "AI_AGENT": "other"})
	assert.Equal(t, "agent:claude-code", detectCLIContext(env, true, true, false))
}

func TestDetectCLIContext_NoTTYNoCI(t *testing.T) {
	assert.Equal(t, "agent:unknown", detectCLIContext(fakeEnv(nil), false, false, false))
}

func TestDetectCLIContext_NoTTYInCI(t *testing.T) {
	assert.Equal(t, "human", detectCLIContext(fakeEnv(nil), false, false, true))
}

func TestDetectCLIContext_PipedOutputOnly(t *testing.T) {
	assert.Equal(t, "human", detectCLIContext(fakeEnv(nil), true, false, false))
}

func TestDetectCLIContext_Interactive(t *testing.T) {
	assert.Equal(t, "human", detectCLIContext(fakeEnv(nil), true, true, false))
}
