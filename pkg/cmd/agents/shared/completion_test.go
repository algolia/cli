package shared

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/iostreams"
)

// nopReadCloser wraps a string in an io.ReadCloser for stdin injection.
func nopReadCloser(s string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(s))
}

func TestBuildMessages_FromMessageConvenience(t *testing.T) {
	got, err := BuildMessages(nopReadCloser(""), MessageInput{Message: "Hello"})
	require.NoError(t, err)
	assert.JSONEq(t, `[{"role":"user","content":"Hello"}]`, string(got))
}

func TestBuildMessages_FromInputFileStdin(t *testing.T) {
	in := `[{"role":"user","content":"Hi"}]`
	got, err := BuildMessages(nopReadCloser(in), MessageInput{InputFile: "-"})
	require.NoError(t, err)
	assert.JSONEq(t, in, string(got))
}

func TestBuildMessages_RejectsBothEmpty(t *testing.T) {
	_, err := BuildMessages(nopReadCloser(""), MessageInput{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "one of --input or --message is required")
}

func TestBuildMessages_RejectsBothSet(t *testing.T) {
	_, err := BuildMessages(nopReadCloser(""), MessageInput{
		InputFile: "-", Message: "x",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "specify either --input or --message")
}

func TestBuildMessages_RejectsInvalidJSON(t *testing.T) {
	_, err := BuildMessages(nopReadCloser(`{not json`), MessageInput{InputFile: "-"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not valid JSON")
}

func TestBuildMessages_RejectsNonArray(t *testing.T) {
	_, err := BuildMessages(nopReadCloser(`{"role":"user"}`), MessageInput{InputFile: "-"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be a JSON array")
}

func TestReadJSONFile_StripsBOMAndValidates(t *testing.T) {
	withBOM := "\xef\xbb\xbf" + `{"name":"x"}`
	got, err := ReadJSONFile(nopReadCloser(withBOM), "-")
	require.NoError(t, err)
	assert.JSONEq(t, `{"name":"x"}`, string(got))
}

func TestReadJSONFile_RejectsInvalidJSON(t *testing.T) {
	_, err := ReadJSONFile(nopReadCloser(`{`), "-")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not valid JSON")
}

func TestMarshalCompletionBody_WithoutConfiguration(t *testing.T) {
	got, err := MarshalCompletionBody(json.RawMessage(`[{"role":"user","content":"x"}]`), nil)
	require.NoError(t, err)
	assert.JSONEq(t, `{"messages":[{"role":"user","content":"x"}]}`, string(got))
}

func TestMarshalCompletionBody_WithConfiguration(t *testing.T) {
	got, err := MarshalCompletionBody(
		json.RawMessage(`[{"role":"user","content":"x"}]`),
		json.RawMessage(`{"model":"gpt-4o-mini"}`),
	)
	require.NoError(t, err)
	assert.JSONEq(
		t,
		`{"messages":[{"role":"user","content":"x"}],"configuration":{"model":"gpt-4o-mini"}}`,
		string(got),
	)
}

func TestMarshalCompletionBody_RejectsEmptyMessages(t *testing.T) {
	_, err := MarshalCompletionBody(nil, nil)
	require.Error(t, err)
}

func TestNormalizeCompatibility(t *testing.T) {
	got, err := NormalizeCompatibility("")
	require.NoError(t, err)
	assert.Equal(t, "ai-sdk-5", string(got))

	for _, in := range []string{"v5", "V5", "ai-sdk-5", "AI-SDK-5", "  AI-SDK-5  "} {
		got, err := NormalizeCompatibility(in)
		require.NoError(t, err, in)
		assert.Equal(t, "ai-sdk-5", string(got), in)
	}
	for _, in := range []string{"v4", "AI-SDK-4"} {
		got, err := NormalizeCompatibility(in)
		require.NoError(t, err, in)
		assert.Equal(t, "ai-sdk-4", string(got), in)
	}
	_, errInvalid := NormalizeCompatibility("nope")
	require.Error(t, errInvalid)
}

func TestRenderCompletion_StreamingEmitsNDJSON(t *testing.T) {
	ios, _, stdout, _ := iostreams.Test()
	body := io.NopCloser(strings.NewReader(strings.Join([]string{
		`data: {"type":"text-delta","delta":"hello"}`,
		``,
		`data: {"type":"text-delta","delta":" world"}`,
		``,
		`data: [DONE]`,
		``,
	}, "\n")))

	require.NoError(t, RenderCompletion(ios, body, "text/event-stream", false))

	scanner := bufio.NewScanner(bytes.NewReader(stdout.Bytes()))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	require.Len(t, lines, 2, "two text-delta events => two NDJSON lines")
	for _, line := range lines {
		var probe map[string]any
		require.NoError(t, json.Unmarshal([]byte(line), &probe), "each line must be valid JSON")
		assert.Equal(t, "text-delta", probe["type"])
	}
}

func TestRenderCompletion_BufferedCopiesVerbatim(t *testing.T) {
	ios, _, stdout, _ := iostreams.Test()
	body := io.NopCloser(strings.NewReader(`{"role":"assistant","content":"hi"}`))

	require.NoError(t, RenderCompletion(ios, body, "application/json", false))
	assert.Equal(t, `{"role":"assistant","content":"hi"}`, stdout.String())
}

// TTY mode renders text-deltas as a single inline assistant reply with
// a trailing newline; tool calls are surfaced as dim annotations; and
// no NDJSON-shaped {"type":...,"data":...} envelope appears.
func TestRenderCompletion_TTYRendersInlineText(t *testing.T) {
	ios, _, stdout, _ := iostreams.Test()
	ios.SetStdoutTTY(true)
	body := io.NopCloser(strings.NewReader(strings.Join([]string{
		`data: {"type":"text-delta","delta":"hello"}`,
		``,
		`data: {"type":"text-delta","delta":" world"}`,
		``,
		`data: [DONE]`,
		``,
	}, "\n")))

	require.NoError(t, RenderCompletion(ios, body, "text/event-stream", false))
	got := stdout.String()
	assert.Contains(t, got, "hello world")
	assert.NotContains(t, got, `"type":"text-delta"`)
}

func TestRenderCompletion_TTYAnnotatesToolCalls(t *testing.T) {
	ios, _, stdout, _ := iostreams.Test()
	ios.SetStdoutTTY(true)
	body := io.NopCloser(strings.NewReader(strings.Join([]string{
		`data: {"type":"text-delta","delta":"thinking..."}`,
		``,
		`data: {"type":"tool-call","toolCallId":"t1","toolName":"search"}`,
		``,
		`data: {"type":"tool-result","toolCallId":"t1"}`,
		``,
		`data: [DONE]`,
		``,
	}, "\n")))

	require.NoError(t, RenderCompletion(ios, body, "text/event-stream", false))
	got := stdout.String()
	assert.Contains(t, got, "thinking...")
	assert.Contains(t, got, "→ tool: search")
	assert.Contains(t, got, "← tool: search")
}

// forceNDJSON=true on a TTY suppresses the rich render — useful when
// users want to see machine output on screen for debugging.
func TestRenderCompletion_TTYForceNDJSON(t *testing.T) {
	ios, _, stdout, _ := iostreams.Test()
	ios.SetStdoutTTY(true)
	body := io.NopCloser(strings.NewReader(strings.Join([]string{
		`data: {"type":"text-delta","delta":"hi"}`,
		``,
		`data: [DONE]`,
		``,
	}, "\n")))

	require.NoError(t, RenderCompletion(ios, body, "text/event-stream", true))
	got := strings.TrimSpace(stdout.String())
	var probe map[string]any
	require.NoError(t, json.Unmarshal([]byte(got), &probe))
	assert.Equal(t, "text-delta", probe["type"])
}

// v4 wire format: text payload is a JSON-encoded string, not an object.
// extractTextDelta must unmarshal the string before printing.
func TestRenderCompletion_TTYHandlesV4Text(t *testing.T) {
	ios, _, stdout, _ := iostreams.Test()
	ios.SetStdoutTTY(true)
	body := io.NopCloser(strings.NewReader(`0:"hello v4"` + "\n"))

	require.NoError(t, RenderCompletion(ios, body, "text/event-stream", false))
	assert.Contains(t, stdout.String(), "hello v4")
}
