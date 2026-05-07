package internal

import (
	"context"
	"fmt"
	"time"

	algoliaComposition "github.com/algolia/algoliasearch-client-go/v4/algolia/composition"

	"github.com/algolia/cli/pkg/iostreams"
)

const (
	// DefaultPollInterval is how often WaitForTask checks task status in production.
	DefaultPollInterval = 2 * time.Second
	// DefaultTimeout is the maximum time WaitForTask will wait before giving up.
	DefaultTimeout = 2 * time.Minute
)

// PollInterval and Timeout are the active timing values used by WaitForTask.
// Override these in tests to avoid slow polling; restore via t.Cleanup.
var (
	PollInterval = DefaultPollInterval
	Timeout      = DefaultTimeout
)

// WaitForTask polls the Compositions API until the given task reaches PUBLISHED
// status, then stops the progress indicator. pollInterval and timeout control
// timing; use DefaultPollInterval and DefaultTimeout in production callers and
// small values (e.g. 1ms / 50ms) in tests.
func WaitForTask(
	io *iostreams.IOStreams,
	client *algoliaComposition.APIClient,
	compositionID string,
	taskID int64,
	pollInterval time.Duration,
	timeout time.Duration,
) error {
	io.StartProgressIndicatorWithLabel("Waiting for task")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			io.StopProgressIndicator()
			return fmt.Errorf("timed out waiting for task %d on composition %s: %w", taskID, compositionID, ctx.Err())
		case <-ticker.C:
			resp, err := client.GetTask(client.NewApiGetTaskRequest(compositionID, taskID))
			if err != nil {
				io.StopProgressIndicator()
				return err
			}
			if resp.Status == algoliaComposition.TASK_STATUS_PUBLISHED {
				io.StopProgressIndicator()
				return nil
			}
			// TASK_STATUS_NOT_PUBLISHED (or any other non-terminal status): keep polling.
		}
	}
}
