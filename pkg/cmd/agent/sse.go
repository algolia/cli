package agent

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

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
