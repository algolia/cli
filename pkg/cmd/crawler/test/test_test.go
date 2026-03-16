package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
