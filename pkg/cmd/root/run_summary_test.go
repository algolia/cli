package root

import (
	"strings"
	"testing"
)

func TestSanitizeRunSummaryCommand(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no secrets", "algolia indices list", "algolia indices list"},
		{"api-key equals", "algolia --api-key=abc123 indices list", "algolia --api-key=***c123 indices list"},
		{"application-id equals", "algolia --application-id=MYAPP indices list", "algolia --application-id=***YAPP indices list"},
		{"api-key space", "algolia --api-key abc123 indices list", "algolia --api-key ***c123 indices list"},
		{"profile short", "algolia -p myprofile indices list", "algolia -p ***file indices list"},
		{"short value", "algolia --api-key=ab indices list", "algolia --api-key=**** indices list"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeRunSummaryCommand(tt.in)
			if got != tt.want {
				t.Errorf("sanitizeRunSummaryCommand(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSanitizeRunSummaryError(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no secrets", "connection refused", "connection refused"},
		{"api_key colon", "error: api_key: abc123def456", "error: api_key: ***f456"},
		{"application_id equals", "invalid application_id=MYAPP", "invalid application_id: ***YAPP"},
		{"truncate long", strings.Repeat("x", 600), strings.Repeat("x", 497) + "..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeRunSummaryError(tt.in)
			if got != tt.want {
				t.Errorf("sanitizeRunSummaryError() = %q, want %q", got, tt.want)
			}
		})
	}
}
