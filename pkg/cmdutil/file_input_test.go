package cmdutil

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONPayload_Inline(t *testing.T) {
	got, err := JSONPayload(`{"a":1}`, "", nil)
	require.NoError(t, err)
	assert.Equal(t, `{"a":1}`, got)
}

func TestJSONPayload_File(t *testing.T) {
	path := filepath.Join(t.TempDir(), "payload.json")
	require.NoError(t, os.WriteFile(path, []byte("  {\"a\":1}\n"), 0o600))

	got, err := JSONPayload("", path, nil)
	require.NoError(t, err)
	assert.Equal(t, `{"a":1}`, got)
}

func TestJSONPayload_Stdin(t *testing.T) {
	got, err := JSONPayload("", "-", io.NopCloser(strings.NewReader("{\"a\":1}\n")))
	require.NoError(t, err)
	assert.Equal(t, `{"a":1}`, got)
}

func TestJSONPayload_Empty(t *testing.T) {
	got, err := JSONPayload("", "", nil)
	require.NoError(t, err)
	assert.Equal(t, "", got)
}

func TestJSONPayload_Conflict(t *testing.T) {
	_, err := JSONPayload(`{"a":1}`, "payload.json", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--json cannot be combined with --file")
}

func TestJSONPayload_MissingFile(t *testing.T) {
	_, err := JSONPayload("", filepath.Join(t.TempDir(), "missing.json"), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "reading file")
}
