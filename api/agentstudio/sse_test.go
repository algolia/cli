package agentstudio

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// collect drains a stream into a slice for assertions.
func collect(t *testing.T, body string) []StreamEvent {
	t.Helper()
	var got []StreamEvent
	require.NoError(t, ParseStream(strings.NewReader(body), func(e StreamEvent) error {
		got = append(got, e)
		return nil
	}))
	return got
}

func TestParseStream_V5StandardSSE(t *testing.T) {
	body := strings.Join([]string{
		`data: {"type":"text-delta","delta":"Hello "}`,
		``,
		`data: {"type":"text-delta","delta":"world"}`,
		``,
		`: keepalive`,
		``,
		`data: {"type":"finish","finishReason":"stop"}`,
		``,
		`data: [DONE]`,
		``,
		// anything after [DONE] must be ignored:
		`data: {"type":"text-delta","delta":"never"}`,
		``,
	}, "\n")

	got := collect(t, body)
	require.Len(t, got, 3)

	assert.Equal(t, "text-delta", got[0].Type)
	assert.JSONEq(t, `{"type":"text-delta","delta":"Hello "}`, string(got[0].Data))
	assert.Equal(t, "text-delta", got[1].Type)
	assert.Equal(t, "finish", got[2].Type)
}

func TestParseStream_V4LineDelimited(t *testing.T) {
	body := strings.Join([]string{
		`f:{"messageId":"m1"}`,
		`0:"Hello "`,
		`0:"world"`,
		`9:{"toolCallId":"call_1","toolName":"search","args":{"q":"x"}}`,
		`a:{"toolCallId":"call_1","result":{"hits":[]}}`,
		`d:{"finishReason":"stop"}`,
	}, "\n")

	got := collect(t, body)
	require.Len(t, got, 6)

	assert.Equal(t, "start-step", got[0].Type)
	assert.Equal(t, "text", got[1].Type)
	// v4 text frames are JSON strings; payload preserved verbatim.
	assert.JSONEq(t, `"Hello "`, string(got[1].Data))
	assert.Equal(t, "tool-call", got[3].Type)
	assert.Equal(t, "tool-result", got[4].Type)
	assert.Equal(t, "finish-message", got[5].Type)
}

func TestParseStream_SkipsMalformedLines(t *testing.T) {
	// Mix of: v5 SSE comment, v5 frame, malformed JSON in v5, malformed v4
	// (unknown prefix), valid v4. The parser should skip all bad lines
	// rather than abort the stream — backends evolve faster than parsers.
	body := strings.Join([]string{
		`: comment`,
		`data: {"type":"text-delta","delta":"a"}`,
		`data: not-json`,
		`Q:{"unknown":"prefix"}`,
		`0:"b"`,
		`0:not-json`,
		`data: [DONE]`,
	}, "\n")

	got := collect(t, body)
	require.Len(t, got, 2)
	assert.Equal(t, "text-delta", got[0].Type)
	assert.Equal(t, "text", got[1].Type)
}

func TestParseStream_V5MissingTypeFallsBackToData(t *testing.T) {
	// v5 frames without a `type` field still get emitted with Type="data"
	// so callers can still forward them via NDJSON.
	body := "data: {\"chunk\":\"x\"}\n\n"
	got := collect(t, body)
	require.Len(t, got, 1)
	assert.Equal(t, "data", got[0].Type)
	assert.JSONEq(t, `{"chunk":"x"}`, string(got[0].Data))
}

func TestParseStream_NilGuards(t *testing.T) {
	require.Error(t, ParseStream(nil, func(StreamEvent) error { return nil }))
	require.Error(t, ParseStream(bytes.NewReader([]byte("x")), nil))
}

func TestParseStream_PropagatesOnEventError(t *testing.T) {
	// onEvent returning a sentinel error halts the stream and bubbles up.
	body := strings.Join([]string{
		`0:"a"`,
		`0:"b"`,
		`0:"c"`,
	}, "\n")

	stop := errors.New("caller cancelled")
	count := 0
	err := ParseStream(strings.NewReader(body), func(StreamEvent) error {
		count++
		if count == 2 {
			return stop
		}
		return nil
	})
	require.ErrorIs(t, err, stop)
	assert.Equal(t, 2, count, "should stop on first onEvent error")
}

func TestParseStream_LargeFrame(t *testing.T) {
	// One big text-delta close to the 5 MiB scanner cap. Confirms the
	// custom buffer kicks in (default bufio cap of 64 KiB would fail).
	big := strings.Repeat("a", 200*1024) // 200 KiB
	body := `0:` + `"` + big + `"` + "\n"
	got := collect(t, body)
	require.Len(t, got, 1)
	assert.Equal(t, "text", got[0].Type)
}

func TestParseStream_HandlesV4StreamThatEndsWithoutSentinel(t *testing.T) {
	// v4 has no `[DONE]`; the body just closes. Parser must return nil
	// at io.EOF, having emitted everything it saw.
	body := "0:\"hello\"\nd:{\"finishReason\":\"stop\"}\n"
	got := collect(t, body)
	require.Len(t, got, 2)
	assert.Equal(t, "text", got[0].Type)
	assert.Equal(t, "finish-message", got[1].Type)
}
