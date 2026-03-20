package agent

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
	"github.com/creack/pty"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/google/uuid"

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
	ID       string    `json:"id"`
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

	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "> ",
		HistoryFile: "",
	})
	if err != nil {
		return fmt.Errorf("failed to initialize readline: %w", err)
	}
	defer rl.Close()

	cs := opts.IO.ColorScheme()
	separator := cs.Gray(strings.Repeat("─", opts.IO.TerminalWidth()))

	fmt.Fprintln(out, "Algolia CLI Agent (type \"exit\" to quit)")
	fmt.Fprintln(out, separator)

	conversationID, err := newConversationID()
	if err != nil {
		return err
	}

	var history []message
	msgCounter := 0

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF or interrupt
			break
		}
		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}
		if input == "exit" {
			break
		}
		if input == "/clear" {
			history = nil
			msgCounter = 0
			newID, idErr := newConversationID()
			if idErr != nil {
				fmt.Fprintf(opts.IO.ErrOut, "Error: %s\n", idErr)
				continue
			}
			conversationID = newID
			fmt.Fprintln(out, separator)
			fmt.Fprintln(out, cs.Gray("Conversation cleared."))
			fmt.Fprintln(out, separator)
			continue
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

		opts.IO.StartProgressIndicatorWithLabel("\nThinking...")
		result, err := sendCompletion(opts, conversationID, history)
		opts.IO.StopProgressIndicator()
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "Error: %s\n", err)
			// Remove the failed user message from history.
			history = history[:len(history)-1]
			continue
		}

		fmt.Fprintln(out, separator)

		for {
			fmt.Fprintln(out)
			if result.Text != "" {
				fmt.Fprintln(out, renderMarkdown(cs, result.Text))
			}
			history = append(history, message{
				ID:   result.MessageID,
				Role: "assistant",
				Parts: []part{
					{Type: "text", Text: result.Text},
				},
			})

			if result.Command == "" {
				fmt.Fprintln(out)
				break
			}

			if isSafeCommand(result.Command) {
				fmt.Fprintf(out, "%s\n", cs.Gray(fmt.Sprintf("\033[3mRunning %s\033[0m", result.Command)))
			} else {
				fmt.Fprintf(out, "%s %s\n", cs.Bold("Suggested command:"), cs.Cyan(result.Command))
				fmt.Fprintln(out)
				rl.SetPrompt("Run this command? [Y/n] ")
				confirmLine, confirmErr := rl.Readline()
				rl.SetPrompt("> ")
				if confirmErr != nil {
					break
				}
				answer := strings.TrimSpace(strings.ToLower(confirmLine))
				// Rewrite the prompt line with the actual answer.
				fmt.Fprintf(out, "\033[1A\033[2K")
				if answer == "" || answer == "y" || answer == "yes" {
					fmt.Fprintf(out, "Run this command? %s\n", cs.Green("Y"))
				} else {
					fmt.Fprintf(out, "Run this command? %s\n", cs.Red("n\n"))
					break
				}
			}

			fmt.Fprintln(out)
			cmdOutput, cmdErr := executeCommand(result.Command)
			msgCounter++
			outputText := fmt.Sprintf("Command `%s` was executed.\nOutput:\n\n%s\n", result.Command, cmdOutput)
			if cmdErr != nil {
				outputText += fmt.Sprintf("\nError: %s", cmdErr)
			}
			history = append(history, message{
				ID:   fmt.Sprintf("alg_msg_%d", msgCounter),
				Role: "user",
				Parts: []part{
					{Text: outputText},
				},
			})

			opts.IO.StartProgressIndicatorWithLabel("\nThinking...")
			followUp, err := sendCompletion(opts, conversationID, history)
			opts.IO.StopProgressIndicator()
			if err != nil {
				fmt.Fprintf(opts.IO.ErrOut, "Error: %s\n", err)
				break
			}
			result = followUp
		}

		fmt.Fprintln(out, separator)
	}

	return nil
}

// sendCompletion sends the conversation to Agent Studio and streams the response.
func sendCompletion(opts *AgentOptions, conversationID string, messages []message) (completionResult, error) {
	url := fmt.Sprintf(
		"https://%s.algolia.net/agent-studio/1/agents/%s/completions?stream=true&compatibilityMode=ai-sdk-5",
		opts.AppID, opts.AgentID,
	)

	reqBody := completionRequest{
		ID:       conversationID,
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

// validateCommand checks that a command string does not contain dangerous shell metacharacters.
func validateCommand(command string) error {
	for _, pattern := range []string{"&&", "||", ";", "$(", "`"} {
		if strings.Contains(command, pattern) {
			return fmt.Errorf("command contains disallowed shell operator: %s", pattern)
		}
	}
	return nil
}

// safeCommands lists read-only command prefixes that can be auto-run without confirmation.
var safeCommands = []string{
	"profile list",
	"application list",
	"indices list",
	"apikeys list",
	"search ",
	"objects browse",
	"settings get",
	"rules browse",
	"synonyms browse",
	"dictionary settings get",
	"dictionary entries browse",
	"describe",
	"open",
	"events tail",
	"crawler list",
	"crawler get",
	"crawler stats",
	"indices config export",
	"indices analyze",
}

// isSafeCommand checks if a command is read-only and can be auto-run.
func isSafeCommand(command string) bool {
	// Strip the leading "algolia " to get the subcommand.
	sub := strings.TrimPrefix(command, "algolia ")
	if sub == command {
		return false
	}
	for _, safe := range safeCommands {
		if strings.HasPrefix(sub, safe) {
			return true
		}
	}
	return false
}

// executeCommand runs a command string inside a PTY so that the child process
// sees a real terminal (IsTerminal returns true). Output is tee'd to the user's
// terminal and captured for the agent context.
func executeCommand(command string) (string, error) {
	if command == "" {
		return "", fmt.Errorf("empty command")
	}
	if err := validateCommand(command); err != nil {
		return "", err
	}
	command = replaceAlgoliaBinary(command)
	cmd := exec.Command("sh", "-c", command)
	cmd.Env = append(os.Environ(), "PAGER=cat")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to start command: %w", err)
	}
	defer ptmx.Close()

	// Tee PTY output to both the real terminal and a buffer.
	var buf bytes.Buffer
	_, _ = io.Copy(io.MultiWriter(os.Stdout, &buf), ptmx)

	_ = cmd.Wait()
	return strings.TrimSpace(stripANSI(buf.String())), nil
}

// stripANSI removes ANSI escape sequences from a string and simulates
// carriage return behavior (overwrites the current line).
func stripANSI(s string) string {
	var lines []string
	var cur strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\033' {
			// Skip CSI sequences: ESC [ ... final byte
			if i+1 < len(s) && s[i+1] == '[' {
				j := i + 2
				for j < len(s) && s[j] >= 0x20 && s[j] <= 0x3F {
					j++
				}
				if j < len(s) {
					j++ // skip final byte
				}
				i = j
				continue
			}
			// Skip other ESC sequences (ESC + one byte)
			i += 2
			continue
		}
		if s[i] == '\r' {
			// \r\n is a normal newline, not a spinner overwrite.
			if i+1 < len(s) && s[i+1] == '\n' {
				lines = append(lines, cur.String())
				cur.Reset()
				i += 2
				continue
			}
			// Standalone \r: discard current line content (spinner overwrite)
			cur.Reset()
			i++
			continue
		}
		if s[i] == '\n' {
			lines = append(lines, cur.String())
			cur.Reset()
			i++
			continue
		}
		cur.WriteByte(s[i])
		i++
	}
	if cur.Len() > 0 {
		lines = append(lines, cur.String())
	}
	// Filter out empty lines from spinner artifacts.
	var result []string
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			result = append(result, l)
		}
	}
	return strings.Join(result, "\n")
}

// replaceAlgoliaBinary replaces "algolia" at command positions with the actual binary path
// (e.g. "./algolia" in dev). Only replaces at the start and after pipes.
func replaceAlgoliaBinary(command string) string {
	bin := os.Args[0]
	if bin == "algolia" {
		return command
	}
	if strings.HasPrefix(command, "algolia ") {
		command = bin + command[len("algolia"):]
	}
	command = strings.ReplaceAll(command, "| algolia ", "| "+bin+" ")
	return command
}

func newConversationID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate conversation ID: %w", err)
	}
	return "alg_cnv_" + id.String(), nil
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
