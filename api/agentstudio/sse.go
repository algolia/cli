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
//
// The backend (rag/models/agent_completion_request.py) makes this a
// REQUIRED query parameter with no default. The CLI defaults to v5 — its
// frames are standard SSE (data: <json>\n\n) with an explicit [DONE]
// sentinel; v4 uses a Vercel-AI-specific line format (<type>:<json>\n)
// with no terminator and is harder to use defensively.
type CompatibilityMode string

const (
	CompatV4 CompatibilityMode = "ai-sdk-4"
	CompatV5 CompatibilityMode = "ai-sdk-5"
)

// StreamEvent is a normalized completion frame.
//
// Both v4 and v5 wire formats are reduced to the same shape so callers
// don't have to branch on protocol:
//
//   - Type: a best-effort identifier. For v4 it's the mapped name
//     ("text", "tool-call", …) derived from the single-character prefix.
//     For v5 it's the value of payload.type when present (the v5 stream
//     is a sequence of `{"type": "...", ...}` JSON objects), else "data".
//   - Data: the JSON payload body. For v4 type=text frames the payload is
//     a JSON string ("hello world"); the parser leaves it as RawMessage
//     so callers can decode into the right shape themselves.
//   - Raw: the original line (without trailing newline). Useful for
//     forwarding bytes verbatim into NDJSON pipelines without losing
//     fidelity from a re-marshal round-trip.
type StreamEvent struct {
	Type string
	Data json.RawMessage
	Raw  string
}

// v4 type-code → human name mapping. Mirrors
// rag/utils/ai_sdk/v4/stream.py:Types and the parser in
// tests/acceptance/helpers/completions.py:_V4_TYPE_MAPPING in
// algolia/conversational-ai. Keep in sync if upstream adds codes.
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

// ParseStream reads the body of a streaming /completions response and
// invokes onEvent for each parsed frame, in order.
//
// The parser sniffs the line prefix to support both wire formats:
//
//   - v5 (default): `data: <json>\n\n`. `[DONE]` ends the stream cleanly;
//     the parser stops without error and returns nil.
//   - v4: `<type-code>:<json>\n`. The stream just ends when the body
//     closes (no sentinel); parser returns nil on io.EOF.
//
// Comment lines (`:keepalive`) and blank lines are skipped silently.
//
// Lines that look like one of the two formats but contain malformed JSON
// are skipped (with no error) — backends evolve faster than parsers, and
// crashing the user's pipe on an unrecognized frame is hostile. If
// onEvent returns an error, ParseStream stops and propagates it.
//
// Cancellation is the caller's job: pass a context-bound body (e.g. one
// from http.Request.WithContext) and close it on cancel.
func ParseStream(body io.Reader, onEvent func(StreamEvent) error) error {
	if body == nil {
		return errors.New("agent studio: ParseStream: body is nil")
	}
	if onEvent == nil {
		return errors.New("agent studio: ParseStream: onEvent is nil")
	}

	scanner := bufio.NewScanner(body)
	// SSE frames can be large (a single text-delta with a long token,
	// or a tool-result with a big JSON blob). Default 64 KiB is too tight;
	// match the existing ScanFile cap (5 MiB).
	const maxLine = 1024 * 5120
	scanner.Buffer(make([]byte, 0, 64*1024), maxLine)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, ":") {
			// blank-line frame separator (v5) or SSE comment / keep-alive
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

// parseLine inspects one line and tries v5 then v4. Returns:
//
//   - evt, ok=true: a usable StreamEvent (caller should emit it)
//   - evt, ok=false: line was malformed for both formats, skip silently
//   - stop=true: the [DONE] sentinel was seen; caller should stop
func parseLine(line string) (evt StreamEvent, ok bool, stop bool) {
	// v5: `data: <json>` (after `set_cache_headers` adds `data: ` prefix
	// and `\n\n` separator).
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

	// v4: `<type>:<json>`. Type is single byte, then a colon, then JSON.
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
