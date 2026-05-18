package delete_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/compositions/delete"
	compinternal "github.com/algolia/cli/pkg/cmd/compositions/internal"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func TestDeleteComposition(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("DELETE", "1/compositions/my-comp"), httpmock.StringResponse(`{"taskID":99}`))
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/99"), httpmock.StringResponse(`{"status":"published"}`))

	compinternal.PollInterval = 1 * time.Millisecond
	compinternal.Timeout = 50 * time.Millisecond
	t.Cleanup(func() {
		compinternal.PollInterval = compinternal.DefaultPollInterval
		compinternal.Timeout = compinternal.DefaultTimeout
	})

	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, "my-comp --confirm", out)
	require.NoError(t, err)

	assert.JSONEq(t, `{"taskID":99}`, strings.TrimSpace(out.String()))
	r.Verify(t)
}

func TestDeleteComposition_WaitsForPublished(t *testing.T) {
	// Verifies polling continues through successive notPublished states.
	r := &httpmock.Registry{}
	r.Register(httpmock.REST("DELETE", "1/compositions/my-comp"), httpmock.StringResponse(`{"taskID":77}`))
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/77"), httpmock.StringResponse(`{"status":"notPublished"}`))
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/77"), httpmock.StringResponse(`{"status":"notPublished"}`))
	r.Register(httpmock.REST("GET", "1/compositions/my-comp/task/77"), httpmock.StringResponse(`{"status":"published"}`))

	compinternal.PollInterval = 1 * time.Millisecond
	compinternal.Timeout = 50 * time.Millisecond
	t.Cleanup(func() {
		compinternal.PollInterval = compinternal.DefaultPollInterval
		compinternal.Timeout = compinternal.DefaultTimeout
	})

	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, "my-comp --confirm", out)
	require.NoError(t, err)

	assert.JSONEq(t, `{"taskID":77}`, strings.TrimSpace(out.String()))
	r.Verify(t)

	taskPolls := 0
	for _, req := range r.Requests {
		if strings.Contains(req.URL.Path, "/task/") {
			taskPolls++
		}
	}
	assert.Equal(t, 3, taskPolls, "expected 3 task status polls (2x notPublished + 1x published)")
}

func TestDeleteComposition_RequiresConfirmation(t *testing.T) {
	// Without --confirm on a non-TTY, the prompt must fail and no HTTP request must be made.
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, "my-comp", out)
	require.Error(t, err)
	assert.Empty(t, out.String())
	assert.Empty(t, r.Requests)
}

func TestDeleteComposition_MissingArg(t *testing.T) {
	r := &httpmock.Registry{}
	f, out := test.NewFactory(false, r, nil, "")
	cmd := delete.NewDeleteCmd(f)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires a <composition-id> argument")
}
