package get

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/test"
)

func Test_NewGetCmd_flagBinding(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")

	called := false
	cmd := NewGetCmd(f, func(opts *GetOptions) error {
		called = true
		assert.Equal(t, "a1", opts.ID)
		if assert.NotNil(t, opts.PrintFlags.OutputFormat) {
			assert.Equal(t, "json", *opts.PrintFlags.OutputFormat)
		}
		return nil
	})

	_, err := test.Execute(cmd, "a1", out)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, called)
}
