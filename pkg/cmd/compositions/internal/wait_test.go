package internal_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	compinternal "github.com/algolia/cli/pkg/cmd/compositions/internal"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

const (
	testPollInterval = 1 * time.Millisecond
	testTimeout      = 50 * time.Millisecond
)

// TestWaitForTask_ImmediatelyPublished verifies that WaitForTask returns nil
// when the very first poll returns "published".
func TestWaitForTask_ImmediatelyPublished(t *testing.T) {
	r := &httpmock.Registry{}
	r.Register(
		httpmock.REST("GET", "1/compositions/my-comp/task/42"),
		httpmock.StringResponse(`{"status":"published"}`),
	)

	f, _ := test.NewFactory(false, r, nil, "")
	client, err := f.CompositionClient()
	require.NoError(t, err)

	io, _, _, _ := iostreams.Test()
	err = compinternal.WaitForTask(io, client, "my-comp", 42, testPollInterval, testTimeout)
	require.NoError(t, err)
	r.Verify(t)
}

// TestWaitForTask_SuccessiveNotPublished verifies that WaitForTask keeps
// polling through multiple "notPublished" responses before eventually
// succeeding on "published".
func TestWaitForTask_SuccessiveNotPublished(t *testing.T) {
	r := &httpmock.Registry{}
	taskPath := "1/compositions/my-comp/task/42"

	// First two polls return notPublished; third returns published.
	r.Register(httpmock.REST("GET", taskPath), httpmock.StringResponse(`{"status":"notPublished"}`))
	r.Register(httpmock.REST("GET", taskPath), httpmock.StringResponse(`{"status":"notPublished"}`))
	r.Register(httpmock.REST("GET", taskPath), httpmock.StringResponse(`{"status":"published"}`))

	f, _ := test.NewFactory(false, r, nil, "")
	client, err := f.CompositionClient()
	require.NoError(t, err)

	io, _, _, _ := iostreams.Test()
	err = compinternal.WaitForTask(io, client, "my-comp", 42, testPollInterval, testTimeout)
	require.NoError(t, err)

	// All three stubs must have been consumed.
	r.Verify(t)
	assert.Len(t, r.Requests, 3, "expected exactly 3 GetTask requests")
}

// TestWaitForTask_Timeout verifies that WaitForTask returns an error containing
// "timed out" when the task never reaches published within the timeout window.
func TestWaitForTask_Timeout(t *testing.T) {
	r := &httpmock.Registry{}
	taskPath := "1/compositions/my-comp/task/42"

	// Register more stubs than we expect to consume so the registry doesn't
	// run out before the timeout fires.
	for range 100 {
		r.Register(httpmock.REST("GET", taskPath), httpmock.StringResponse(`{"status":"notPublished"}`))
	}

	f, _ := test.NewFactory(false, r, nil, "")
	client, err := f.CompositionClient()
	require.NoError(t, err)

	io, _, _, _ := iostreams.Test()
	err = compinternal.WaitForTask(io, client, "my-comp", 42, testPollInterval, testTimeout)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timed out")
	assert.Contains(t, err.Error(), "42")
	assert.Contains(t, err.Error(), "my-comp")
}

// TestWaitForTask_APIError verifies that WaitForTask surfaces errors from the
// GetTask API call immediately without retrying.
func TestWaitForTask_APIError(t *testing.T) {
	r := &httpmock.Registry{}
	taskPath := "1/compositions/my-comp/task/42"

	r.Register(
		httpmock.REST("GET", taskPath),
		httpmock.ErrorResponse(),
	)

	f, _ := test.NewFactory(false, r, nil, "")
	client, err := f.CompositionClient()
	require.NoError(t, err)

	io, _, _, _ := iostreams.Test()
	err = compinternal.WaitForTask(io, client, "my-comp", 42, testPollInterval, testTimeout)
	require.Error(t, err)
	// Only one request should have been made before returning.
	assert.Len(t, r.Requests, 1, "expected WaitForTask to stop after first API error")
}
