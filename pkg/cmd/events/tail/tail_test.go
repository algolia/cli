package tail

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/api/insights"
)

func TestUnseenEventsSkipsDuplicateRequestIDs(t *testing.T) {
	now := time.Now().UTC()
	seenRequestIDs := map[string]time.Time{}
	events := []insights.EventWrapper{
		{
			RequestID: "req-1",
			Event: insights.Event{
				EventName: "first",
				Timestamp: insights.Timestamp{Time: now},
			},
		},
		{
			RequestID: "req-1",
			Event: insights.Event{
				EventName: "duplicate",
				Timestamp: insights.Timestamp{Time: now.Add(time.Second)},
			},
		},
		{
			RequestID: "req-2",
			Event: insights.Event{
				EventName: "second",
				Timestamp: insights.Timestamp{Time: now.Add(2 * time.Second)},
			},
		},
	}

	freshEvents := unseenEvents(events, seenRequestIDs)

	require.Len(t, freshEvents, 2)
	require.Equal(t, "first", freshEvents[0].Event.EventName)
	require.Equal(t, "second", freshEvents[1].Event.EventName)
	require.Contains(t, seenRequestIDs, "req-1")
	require.Contains(t, seenRequestIDs, "req-2")
}

func TestUnseenEventsKeepsEventsWithoutRequestID(t *testing.T) {
	freshEvents := unseenEvents([]insights.EventWrapper{
		{Event: insights.Event{EventName: "first"}},
		{Event: insights.Event{EventName: "second"}},
	}, map[string]time.Time{})

	require.Len(t, freshEvents, 2)
}

func TestPruneSeenRequestIDsRemovesOldEntries(t *testing.T) {
	now := time.Now().UTC()
	seenRequestIDs := map[string]time.Time{
		"stale":  now.Add(-2 * Interval),
		"recent": now.Add(-Interval / 2),
	}

	pruneSeenRequestIDs(seenRequestIDs, now.Add(-Interval))

	require.NotContains(t, seenRequestIDs, "stale")
	require.Contains(t, seenRequestIDs, "recent")
}
