package telemetry

import (
	"os"
	"strings"

	"github.com/algolia/cli/pkg/utils"
)

const (
	ContextHuman        = "human"
	ContextAgentUnknown = "agent:unknown"
)

var agentEnvVars = []struct {
	envVar string
	name   string
}{
	{"CLAUDECODE", "claude-code"},
	{"CLAUDE_CODE", "claude-code"},
	{"CLAUDE_CODE_SESSION_ID", "claude-code"},
	{"CURSOR_AGENT", "cursor"},
	{"CODEX_THREAD_ID", "codex"},
	{"CODEX_SANDBOX", "codex"},
	{"CODEX_CI", "codex"},
	{"GEMINI_CLI", "gemini"},
	{"COPILOT_CLI", "github-copilot"},
	{"OPENCODE", "opencode"},
	{"OPENCODE_CLIENT", "opencode"},
	{"AMP_CURRENT_THREAD_ID", "amp"},
	{"AUGMENT_AGENT", "augment"},
	{"GOOSE_TERMINAL", "goose"},
	{"CLINE_ACTIVE", "cline"},
	{"ANTIGRAVITY_AGENT", "antigravity"},
	{"PI_CODING_AGENT", "pi"},
	{"KIRO_AGENT_PATH", "kiro"},
}

func DetectCLIContext() string {
	return detectCLIContext(
		os.Getenv,
		utils.IsTerminal(os.Stdin),
		utils.IsTerminal(os.Stdout),
		utils.IsCI(),
	)
}

func detectCLIContext(getenv func(string) string, stdinTTY, stdoutTTY, isCI bool) string {
	for _, agent := range agentEnvVars {
		if getenv(agent.envVar) != "" {
			return "agent:" + agent.name
		}
	}
	for _, key := range []string{"AI_AGENT", "AGENT"} {
		v, _, _ := strings.Cut(strings.ToLower(strings.TrimSpace(getenv(key))), "_")
		if v != "" {
			return "agent:" + v
		}
	}
	if !stdinTTY && !stdoutTTY && !isCI {
		return ContextAgentUnknown
	}
	return ContextHuman
}
