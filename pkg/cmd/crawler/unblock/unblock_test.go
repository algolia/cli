package unblock

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/test"
)

func TestNewUnblockCmd_confirmFlag(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")

	called := false
	cmd := NewUnblockCmd(f, func(opts *UnblockOptions) error {
		called = true
		assert.Equal(t, "my-crawler", opts.ID)
		return nil
	})

	_, err := test.Execute(cmd, "my-crawler --confirm", out)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
}

func TestNewUnblockCmd_rejectsControlChars(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewUnblockCmd(f, nil)

	_, err := test.Execute(cmd, "\"bad\nid\" --confirm", out)
	if err == nil {
		t.Fatal("expected error")
	}

	assert.EqualError(t, err, "crawler_id must not contain control characters")
}
