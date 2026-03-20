package agent

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
)

// Build-time variables injected via ldflags in .goreleaser.yml.
var (
	DefaultAgentID     string
	DefaultAgentAppID  string
	DefaultAgentAPIKey string
)

// AgentOptions holds the configuration for the agent command.
type AgentOptions struct {
	IO *iostreams.IOStreams

	AgentID string
	AppID   string
	APIKey  string
}

// message represents a single message in the conversation.
type message struct {
	ID    string `json:"id,omitempty"`
	Role  string `json:"role"`
	Parts []part `json:"parts"`
}

// part represents a content part within a message.
type part struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text"`
}

// completionRequest is the request body sent to Agent Studio.
type completionRequest struct {
	Messages []message `json:"messages"`
}

// sseEvent represents a parsed SSE data payload.
type sseEvent struct {
	Type      string          `json:"type"`
	ID        string          `json:"id,omitempty"`
	MessageID string          `json:"messageId,omitempty"`
	Delta     string          `json:"delta,omitempty"`
	ToolName  string          `json:"toolName,omitempty"`
	Input     json.RawMessage `json:"input,omitempty"`
}

// suggestCommandInput represents the input for the suggestCommand tool.
type suggestCommandInput struct {
	Command string `json:"command"`
}

func NewAgentCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &AgentOptions{
		IO:      f.IOStreams,
		AgentID: envOrDefault("ALGOLIA_AGENT_ID", DefaultAgentID),
		AppID:   envOrDefault("ALGOLIA_AGENT_APP_ID", DefaultAgentAppID),
		APIKey:  envOrDefault("ALGOLIA_AGENT_API_KEY", DefaultAgentAPIKey),
	}

	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Chat with an AI agent that suggests Algolia CLI commands",
		Long:  "Interactive chat with an AI agent that advises CLI commands for your use case. The agent only prints suggestions — it does not execute commands.",
		Example: heredoc.Doc(`
			$ algolia agent
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAgent(opts)
		},
	}

	auth.DisableAuthCheck(cmd)

	return cmd
}

func runAgent(opts *AgentOptions) error {
	if opts.AgentID == "" || opts.AppID == "" || opts.APIKey == "" {
		return fmt.Errorf("agent credentials are not configured")
	}

	out := opts.IO.Out
	scanner := bufio.NewScanner(os.Stdin)

	cs := opts.IO.ColorScheme()
	separator := cs.Gray(strings.Repeat("─", opts.IO.TerminalWidth()))

	fmt.Fprintln(out, "Algolia CLI Agent (type \"exit\" to quit)")
	fmt.Fprintln(out, separator)

	var history []message
	msgCounter := 0

	for {
		fmt.Fprint(out, "> ")

		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "exit" {
			break
		}

		msgCounter++
		userMsg := message{
			ID:   fmt.Sprintf("alg_msg_%d", msgCounter),
			Role: "user",
			Parts: []part{
				{Text: input},
			},
		}
		history = append(history, userMsg)

		opts.IO.StartProgressIndicator()
		result, err := sendCompletion(opts, history)
		opts.IO.StopProgressIndicator()
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "Error: %s\n", err)
			// Remove the failed user message from history.
			history = history[:len(history)-1]
			continue
		}

		fmt.Fprintln(out, separator)
		fmt.Fprintln(out)
		if result.Text != "" {
			fmt.Fprintln(out, renderMarkdown(opts.IO.ColorScheme(), result.Text))
			fmt.Fprintln(out)
		}
		if result.Command != "" {
			fmt.Fprintf(out, "%s %s\n", cs.Bold("Suggested command:"), cs.Cyan(result.Command))
			fmt.Fprintln(out)
		}
		fmt.Fprintln(out, separator)

		history = append(history, message{
			ID:   result.MessageID,
			Role: "assistant",
			Parts: []part{
				{Type: "text", Text: result.Text},
			},
		})
	}

	return nil
}

// sendCompletion sends the conversation to Agent Studio and streams the response.
func sendCompletion(opts *AgentOptions, messages []message) (completionResult, error) {
	url := fmt.Sprintf(
		"https://%s.algolia.net/agent-studio/1/agents/%s/completions?stream=true&compatibilityMode=ai-sdk-5",
		opts.AppID, opts.AgentID,
	)

	reqBody := completionRequest{
		Messages: messages,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return completionResult{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return completionResult{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-algolia-application-id", opts.AppID)
	req.Header.Set("X-Algolia-API-Key", opts.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return completionResult{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return completionResult{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return parseSSEStream(resp.Body)
}

// completionResult holds the parsed response from the SSE stream.
type completionResult struct {
	Text      string
	MessageID string
	Command   string // optional, from suggestCommand tool
}

// parseSSEStream reads an SSE stream and collects text deltas.
func parseSSEStream(r io.Reader) (completionResult, error) {
	var res completionResult
	var textBuf strings.Builder
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var event sseEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}
		switch event.Type {
		case "start":
			res.MessageID = event.MessageID
		case "text-delta":
			textBuf.WriteString(event.Delta)
		case "tool-input-available":
			if event.ToolName == "suggestCommand" {
				var input suggestCommandInput
				if err := json.Unmarshal(event.Input, &input); err == nil {
					res.Command = input.Command
				}
			}
		}
	}

	res.Text = textBuf.String()
	return res, nil
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
