package insights

import (
	"encoding/json"
	"time"
)

type EventWrapper struct {
	Event     Event    `json:"event"`
	RequestID string   `json:"requestID"`
	Status    int      `json:"status"`
	Errors    []string `json:"errors"`
	Headers   map[string][]string
}

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	*t = Timestamp{time.Unix(0, timestamp*int64(time.Millisecond))}
	return nil
}

type Event struct {
	EventType string    `json:"eventType"`
	EventName string    `json:"eventName"`
	Index     string    `json:"index"`
	UserToken string    `json:"userToken"`
	Timestamp Timestamp `json:"timestamp"`
	ObjectIDs []string  `json:"objectIDs,omitempty"`
	Positions []int     `json:"positions,omitempty"`
	QueryID   string    `json:"queryID,omitempty"`
	Filters   []string  `json:"filters,omitempty"`
}
