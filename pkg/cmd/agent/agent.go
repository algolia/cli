package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"

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
		Long:  "Interactive chat with an AI agent that can suggest and execute Algolia CLI commands for your use case.",
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

			if isBlockedCommand(result.Command) {
				fmt.Fprintf(out, "%s %s\n", cs.Bold("Suggested command:"), cs.Cyan(result.Command))
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
			cmdOutput = compactJSON(cmdOutput)
			cmdOutput = truncateOutput(cmdOutput, 10)
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

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return completionResult{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return completionResult{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return parseSSEStream(resp.Body)
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