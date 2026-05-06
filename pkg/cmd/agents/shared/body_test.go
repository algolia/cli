package shared

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/iostreams"
)

func TestSourceLabel(t *testing.T) {
	assert.Equal(t, "stdin", SourceLabel("-"))
	assert.Equal(t, "spec.json", SourceLabel("spec.json"))
	assert.Equal(t, "/abs/path.json", SourceLabel("/abs/path.json"))
}

func TestTrimUTF8BOM(t *testing.T) {
	assert.Equal(t, []byte(`{"a":1}`), TrimUTF8BOM([]byte(`{"a":1}`)))
	withBOM := append([]byte("\xef\xbb\xbf"), []byte(`{"a":1}`)...)
	assert.Equal(t, []byte(`{"a":1}`), TrimUTF8BOM(withBOM))
	assert.Equal(t, []byte{}, TrimUTF8BOM([]byte{}))
}

func TestPrintDryRun_Human(t *testing.T) {
	io, _, stdout, _ := iostreams.Test()
	// pf has a default of "json" but wantsStructured=false => human path.
	pf := cmdutil.NewPrintFlags().WithDefaultOutput("json")

	body := []byte(`{"name":"Concierge","instructions":"x"}`)
	require.NoError(t, PrintDryRun(io, pf, false, "create_agent", "POST /1/agents", "spec.json", body, nil))

	out := stdout.String()
	assert.Contains(t, out, "Dry run: would POST /1/agents")
	assert.Contains(t, out, "(39 bytes from spec.json)")
	// Pretty-printed body must appear so users can lint visually.
	assert.Contains(t, out, `"name": "Concierge"`)
	assert.Contains(t, out, `"instructions": "x"`)
}

func TestPrintDryRun_Structured(t *testing.T) {
	io, _, stdout, _ := iostreams.Test()
	pf := cmdutil.NewPrintFlags().WithDefaultOutput("json")

	body := []byte(`{"name":"Concierge"}`)
	require.NoError(t, PrintDryRun(io, pf, true, "update_agent", "PATCH /1/agents/abc",
		"-", body, map[string]any{"agentId": "abc"}))

	var got map[string]any
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &got))
	assert.Equal(t, "update_agent", got["action"])
	assert.Equal(t, "PATCH /1/agents/abc", got["request"])
	assert.Equal(t, "stdin", got["source"])
	assert.Equal(t, true, got["dryRun"])
	assert.Equal(t, "abc", got["agentId"])
	// body should be the parsed JSON object (not a string).
	assert.IsType(t, map[string]any{}, got["body"])
	assert.Equal(t, "Concierge", got["body"].(map[string]any)["name"])
}
