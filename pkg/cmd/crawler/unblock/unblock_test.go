package unblock

import (
	"testing"

	"github.com/algolia/cli/api/crawler"
	"github.com/algolia/cli/pkg/httpmock"
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

func TestNewUnblockCmd_dryRunFlag(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")

	called := false
	cmd := NewUnblockCmd(f, func(opts *UnblockOptions) error {
		called = true
		assert.Equal(t, "my-crawler", opts.ID)
		assert.True(t, opts.DryRun)
		return nil
	})

	_, err := test.Execute(cmd, "my-crawler --dry-run", out)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, called)
}

func Test_runUnblockCmd_dryRunJSON(t *testing.T) {
	r := httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "api/1/crawlers/my-crawler"),
		httpmock.JSONResponse(crawler.Crawler{
			ID:             "my-crawler",
			Name:           "My crawler",
			BlockingTaskID: "task-123",
			BlockingError:  "crawler is blocked",
		}),
	)
	defer r.Verify(t)

	f, out := test.NewFactory(false, &r, nil, "")
	cmd := NewUnblockCmd(f, nil)

	out, err := test.Execute(cmd, "my-crawler --dry-run --output json", out)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, out.String(), `"action":"unblock_crawler"`)
	assert.Contains(t, out.String(), `"id":"my-crawler"`)
	assert.Contains(t, out.String(), `"blockingError":"crawler is blocked"`)
	assert.Contains(t, out.String(), `"dryRun":true`)
}
