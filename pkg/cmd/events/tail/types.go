package tail

// FetchEventsResponse represents the API response for GET /1/events
type FetchEventsResponse struct {
	Events []EventWrapper `json:"events"`
}

// EventWrapper is an event plus associated data, such as errors or headers
type EventWrapper struct {
	Errors []string `json:"errors"`
	Event  Event    `json:"event"`
	Status int      `json:"status"`
}

// Event represents an Insights event with reduced properties just for printing
type Event struct {
	EventName string `json:"eventName"`
	EventType string `json:"eventType"`
	Index     string `json:"index"`
	UserToken string `json:"userToken"`
	Timestamp int64  `json:"timestamp"`
}
