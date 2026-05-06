package shared

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/algolia/cli/api/agentstudio"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
)

// MessageInput is what `agents test` and `agents run` accept for their
// `messages` field, in declining priority:
//
//   - InputFile != "": read JSON from file (or "-" for stdin) and use it
//     verbatim. Must already be a JSON array of message objects.
//   - Message != "": wrap as a single-shot user message
//     [{"role":"user","content":"<Message>"}]. Convenient for one-liners.
//
// Exactly one must be non-empty (BuildMessages enforces).
type MessageInput struct {
	InputFile string
	Message   string
}

// BuildMessages resolves a MessageInput into a JSON array suitable for
// the `messages` field of AgentCompletionRequest.
func BuildMessages(stdin io.ReadCloser, in MessageInput) (json.RawMessage, error) {
	hasFile := in.InputFile != ""
	hasMsg := strings.TrimSpace(in.Message) != ""

	switch {
	case hasFile && hasMsg:
		return nil, cmdutil.FlagErrorf("specify either --input or --message, not both")
	case !hasFile && !hasMsg:
		return nil, cmdutil.FlagErrorf("one of --input or --message is required")
	case hasMsg:
		// Marshal handles escaping correctly for arbitrary content.
		body, _ := json.Marshal([]map[string]string{
			{"role": "user", "content": in.Message},
		})
		return body, nil
	}

	raw, err := cmdutil.ReadFile(in.InputFile, stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read messages from %s: %w", SourceLabel(in.InputFile), err)
	}
	raw = TrimUTF8BOM(raw)
	if !json.Valid(raw) {
		return nil, cmdutil.FlagErrorf("messages in %s is not valid JSON", SourceLabel(in.InputFile))
	}
	// Cheap shape check so the user gets a CLI-side error instead of a
	// 422 from the backend's discriminator validator.
	var probe any
	_ = json.Unmarshal(raw, &probe)
	if _, ok := probe.([]any); !ok {
		return nil, cmdutil.FlagErrorf("messages in %s must be a JSON array", SourceLabel(in.InputFile))
	}
	return raw, nil
}

// ReadJSONFile reads a JSON document from a path (or "-" for stdin),
// strips a UTF-8 BOM if present, and validates well-formedness. Used by
// `agents test --config`.
func ReadJSONFile(stdin io.ReadCloser, file string) (json.RawMessage, error) {
	raw, err := cmdutil.ReadFile(file, stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", SourceLabel(file), err)
	}
	raw = TrimUTF8BOM(raw)
	if !json.Valid(raw) {
		return nil, cmdutil.FlagErrorf("%s is not valid JSON", SourceLabel(file))
	}
	return raw, nil
}

// CompletionRequest is the assembled body the CLI POSTs.
//
// Mirrors AgentCompletionRequest in the backend
// (rag/models/agent_completion_request.py). Kept narrow to the fields
// the CLI actually composes; richer structure (algolia.searchParameters,
// tool_approvals, conversation `id`) can be added by hand-writing the
// JSON file and using --input.
type CompletionRequest struct {
	Messages      json.RawMessage `json:"messages"`
	Configuration json.RawMessage `json:"configuration,omitempty"`
}

// NormalizeCompatibility maps user-facing aliases ("v4", "v5") to the
// backend's canonical wire values. Empty defaults to v5 (CLI default;
// see CompletionOptions.Compatibility doc for rationale). Shared
// between `agents test` and `agents run`.
func NormalizeCompatibility(s string) (agentstudio.CompatibilityMode, error) {
	switch s {
	case "", "v5", "ai-sdk-5":
		return agentstudio.CompatV5, nil
	case "v4", "ai-sdk-4":
		return agentstudio.CompatV4, nil
	default:
		return "", cmdutil.FlagErrorf("invalid --compatibility %q (allowed: v4, v5)", s)
	}
}

// MarshalCompletionBody assembles the JSON body. Configuration is
// optional (nil for `agents run`, required for `agents test`).
func MarshalCompletionBody(messages, configuration json.RawMessage) (json.RawMessage, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("messages must not be empty")
	}
	body, err := json.Marshal(CompletionRequest{
		Messages:      messages,
		Configuration: configuration,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal completion body: %w", err)
	}
	return body, nil
}

// RenderCompletion writes a /completions response to io.Out.
//
// Output rules (intentionally simple — power users compose with `jq`):
//
//   - Streaming responses (Content-Type: text/event-stream): one parsed
//     event per line as compact JSON ({"type":"…","data":{…}}). Matches
//     NDJSON conventions used elsewhere in the CLI (e.g. `events tail`).
//     The raw `data` field preserves whatever the backend sent so the
//     v4/v5 distinction is recoverable downstream.
//   - Buffered responses: body is copied to stdout verbatim. The backend
//     already returns a single JSON document; re-encoding would lose
//     canonicalization (key order, number formatting) for no gain.
//
// Cancellation: the caller is responsible for signal handling — pass a
// ctx that is cancelled on SIGINT and the underlying transport will
// tear the stream down cleanly. The body is closed before returning.
func RenderCompletion(ios *iostreams.IOStreams, body io.ReadCloser, contentType string) error {
	defer body.Close()

	if !strings.Contains(contentType, "text/event-stream") {
		_, err := io.Copy(ios.Out, body)
		return err
	}

	enc := json.NewEncoder(ios.Out)
	enc.SetEscapeHTML(false)
	return agentstudio.ParseStream(body, func(e agentstudio.StreamEvent) error {
		return enc.Encode(struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}{Type: e.Type, Data: e.Data})
	})
}
