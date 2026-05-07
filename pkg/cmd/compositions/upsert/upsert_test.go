package upsert_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	compinternal "github.com/algolia/cli/pkg/cmd/compositions/internal"
	"github.com/algolia/cli/pkg/cmd/compositions/upsert"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

// validCompBody is a well-formed composition body using the injection behavior schema.
const validCompBody = `{"objectID":"my-comp","behavior":{"injection":{"indices":[{"indexName":"products","maxRecommendations":5}]}}}`

func TestUpsertComposition(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("PUT", "1/compositions/my-comp"), httpmock.StringResponse(`{"taskID":42}`))
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/42"), httpmock.StringResponse(`{"status":"published"}`))

	compinternal.PollInterval = 1 * time.Millisecond
	compinternal.Timeout = 50 * time.Millisecond
	t.Cleanup(func() {
		compinternal.PollInterval = compinternal.DefaultPollInterval
		compinternal.Timeout = compinternal.DefaultTimeout
	})

	f, out := test.NewFactory(false, r, nil, "")
	cmd := upsert.NewUpsertCmd(f)
	out.InBuf.WriteString(validCompBody)
	_, err := test.Execute(cmd, "my-comp --file -", out)
	require.NoError(t, err)

	assert.JSONEq(t, `{"taskID":42}`, strings.TrimSpace(out.String()))
	r.Verify(t)
}

func TestUpsertComposition_WaitsForPublished(t *testing.T) {
	// Verifies that the command polls until published before printing.
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("PUT", "1/compositions/my-comp"), httpmock.StringResponse(`{"taskID":99}`))
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/99"), httpmock.StringResponse(`{"status":"notPublished"}`))
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/99"), httpmock.StringResponse(`{"status":"notPublished"}`))
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/99"), httpmock.StringResponse(`{"status":"published"}`))

	compinternal.PollInterval = 1 * time.Millisecond
	compinternal.Timeout = 50 * time.Millisecond
	t.Cleanup(func() {
		compinternal.PollInterval = compinternal.DefaultPollInterval
		compinternal.Timeout = compinternal.DefaultTimeout
	})

	f, out := test.NewFactory(false, r, nil, "")
	cmd := upsert.NewUpsertCmd(f)
	out.InBuf.WriteString(validCompBody)
	_, err := test.Execute(cmd, "my-comp --file -", out)
	require.NoError(t, err)

	assert.JSONEq(t, `{"taskID":99}`, strings.TrimSpace(out.String()))
	r.Verify(t)

	taskPolls := 0
	for _, req := range r.Requests {
		if strings.Contains(req.URL.Path, "/task/") {
			taskPolls++
		}
	}
	assert.Equal(t, 3, taskPolls, "expected 3 task status polls (2x notPublished + 1x published)")
}

func TestUpsertComposition_MissingFile(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := upsert.NewUpsertCmd(f)
	_, err := test.Execute(cmd, "my-comp", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file")
}

func TestUpsertComposition_InvalidJSON(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := upsert.NewUpsertCmd(f)
	out.InBuf.WriteString(`not-json`)
	_, err := test.Execute(cmd, "my-comp --file -", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing composition JSON")
}

func TestUpsertComposition_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := upsert.NewUpsertCmd(f)
	_, err := test.Execute(cmd, "--file -", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a <composition-id> argument")
}
