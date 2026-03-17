package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/test"
)

func Test_NewTestCmd_outputFlag(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")

	called := false
	cmd := NewTestCmd(f, func(opts *TestOptions) error {
		called = true
		assert.Equal(t, "my-crawler", opts.ID)
		assert.Equal(t, "https://example.com", opts.URL)
		if assert.NotNil(t, opts.PrintFlags.OutputFormat) {
			assert.Equal(t, "json", *opts.PrintFlags.OutputFormat)
		}
		return nil
	})

	_, err := test.Execute(cmd, "my-crawler --url https://example.com --output json", out)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
}

func Test_NewTestCmd_dryRunFlag(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "config.json")
	err := os.WriteFile(tmpFile, []byte(`{"appId":"test-app"}`), 0o600)
	require.NoError(t, err)

	f, out := test.NewFactory(false, nil, nil, "")

	called := false
	cmd := NewTestCmd(f, func(opts *TestOptions) error {
		called = true
		assert.Equal(t, "my-crawler", opts.ID)
		assert.Equal(t, "https://example.com", opts.URL)
		assert.True(t, opts.DryRun)
		if assert.NotNil(t, opts.config) {
			assert.Equal(t, "test-app", opts.config.AppID)
		}
		return nil
	})

	_, err = test.Execute(cmd, fmt.Sprintf("my-crawler --url https://example.com -F '%s' --dry-run", tmpFile), out)
	require.NoError(t, err)

	assert.True(t, called)
}

func Test_runTestCmd_dryRunJSON(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "config.json")
	err := os.WriteFile(tmpFile, []byte(`{"appId":"test-app"}`), 0o600)
	require.NoError(t, err)

	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewTestCmd(f, nil)

	out, err = test.Execute(cmd, fmt.Sprintf("my-crawler --url https://example.com -F '%s' --dry-run --output json", tmpFile), out)
	require.NoError(t, err)

	assert.Contains(t, out.String(), `"action":"test_crawler"`)
	assert.Contains(t, out.String(), `"id":"my-crawler"`)
	assert.Contains(t, out.String(), `"url":"https://example.com"`)
	assert.Contains(t, out.String(), `"config":{"appId":"test-app"}`)
	assert.Contains(t, out.String(), `"dryRun":true`)
}
