package agentstudio

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// CompatibilityMode selects the streaming protocol the backend emits.
// Required server-side; CLI defaults to v5. See docs/agents.md.
type CompatibilityMode string

const (
	CompatV4 CompatibilityMode = "ai-sdk-4"
	CompatV5 CompatibilityMode = "ai-sdk-5"
)

// StreamEvent is a normalised completion frame. v4 and v5 are reduced
// to the same shape so callers don't have to branch on protocol.
type StreamEvent struct {
	Type string
	Data json.RawMessage
	Raw  string
}

// v4 type-code → human name. Mirrors rag/utils/ai_sdk/v4/stream.py:Types
// in algolia/conversational-ai. Keep in sync if upstream adds codes.
var v4TypeNames = map[byte]string{
	'0': "text",
	'2': "data",
	'3': "error",
	'8': "message-annotation",
	'9': "tool-call",
	'a': "tool-result",
	'b': "tool-call-streaming-start",
	'c': "tool-call-delta",
	'd': "finish-message",
	'e': "finish-step",
	'f': "start-step",
	'g': "reasoning",
	'h': "source",
	'i': "redacted-reasoning",
	'j': "reasoning-signature",
	'k': "file",
}

// ParseStream reads a streaming /completions body and invokes onEvent
// for each parsed frame. Malformed frames are skipped silently
// (backends evolve faster than parsers; crashing the user's pipe is
// hostile). The v5 [DONE] sentinel ends the stream cleanly. Callers
// own cancellation via the body's underlying context.
func ParseStream(body io.Reader, onEvent func(StreamEvent) error) error {
	if body == nil {
		return errors.New("agent studio: ParseStream: body is nil")
	}
	if onEvent == nil {
		return errors.New("agent studio: ParseStream: onEvent is nil")
	}

	scanner := bufio.NewScanner(body)
	// SSE frames can be large (long text-deltas, big tool-result blobs).
	const maxLine = 1024 * 5120
	scanner.Buffer(make([]byte, 0, 64*1024), maxLine)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}

		evt, ok, stop := parseLine(line)
		if stop {
			return nil
		}
		if !ok {
			continue
		}
		if err := onEvent(evt); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("agent studio: ParseStream: %w", err)
	}
	return nil
}

// parseLine returns (evt, ok, stop): ok=true => emit; stop=true => [DONE].
func parseLine(line string) (evt StreamEvent, ok bool, stop bool) {
	// v5: `data: <json>`.
	if rest, found := strings.CutPrefix(line, "data: "); found {
		rest = strings.TrimSpace(rest)
		if rest == "" {
			return StreamEvent{}, false, false
		}
		if rest == "[DONE]" {
			return StreamEvent{}, false, true
		}
		var probe map[string]any
		if err := json.Unmarshal([]byte(rest), &probe); err != nil {
			return StreamEvent{}, false, false
		}
		typ, _ := probe["type"].(string)
		if typ == "" {
			typ = "data"
		}
		return StreamEvent{Type: typ, Data: json.RawMessage(rest), Raw: line}, true, false
	}

	// v4: `<type>:<json>` — type is a single byte, then a colon, then JSON.
	if len(line) < 3 || line[1] != ':' {
		return StreamEvent{}, false, false
	}
	name, known := v4TypeNames[line[0]]
	if !known {
		return StreamEvent{}, false, false
	}
	payload := strings.TrimSpace(line[2:])
	if payload == "" {
		return StreamEvent{}, false, false
	}
	if !json.Valid([]byte(payload)) {
		return StreamEvent{}, false, false
	}
	return StreamEvent{Type: name, Data: json.RawMessage(payload), Raw: line}, true, false
}
