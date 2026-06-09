package cmdutil

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/api/dashboard"
)

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantClass  string
		wantSource string
		wantStatus int
	}{
		{
			name: "nil",
		},
		{
			name:       "flag error",
			err:        FlagErrorf("bad flag"),
			wantClass:  "validation_error",
			wantSource: ErrorSourceValidation,
		},
		{
			name: "wrapped api error carries status",
			err: fmt.Errorf(
				"create failed: %w",
				&dashboard.APIError{StatusCode: 500, Message: "boom"},
			),
			wantClass:  "http_500",
			wantSource: ErrorSourceAPI,
			wantStatus: 500,
		},
		{
			name:       "session expired",
			err:        fmt.Errorf("call failed: %w", dashboard.ErrSessionExpired),
			wantClass:  "session_expired",
			wantSource: ErrorSourceAPI,
			wantStatus: 401,
		},
		{
			name:       "cluster unavailable",
			err:        &dashboard.ErrClusterUnavailable{Region: "us", Message: "no cluster"},
			wantClass:  "cluster_unavailable",
			wantSource: ErrorSourceAPI,
		},
		{
			name:       "context canceled",
			err:        context.Canceled,
			wantClass:  "canceled",
			wantSource: ErrorSourceLocal,
		},
		{
			name:       "deadline exceeded",
			err:        context.DeadlineExceeded,
			wantClass:  "timeout",
			wantSource: ErrorSourceNetwork,
		},
		{
			name:       "dns error",
			err:        &net.DNSError{Err: "no such host", Name: "example.invalid"},
			wantClass:  "dns_error",
			wantSource: ErrorSourceNetwork,
		},
		{
			name:       "generic error falls back to unknown",
			err:        errors.New("something went wrong"),
			wantClass:  "unknown",
			wantSource: ErrorSourceLocal,
		},
		{
			name:       "wrapped generic error stays unknown",
			err:        fmt.Errorf("outer: %w", errors.New("inner")),
			wantClass:  "unknown",
			wantSource: ErrorSourceLocal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			class, source, status := ClassifyError(tt.err)
			assert.Equal(t, tt.wantClass, class)
			assert.Equal(t, tt.wantSource, source)
			assert.Equal(t, tt.wantStatus, status)

			assert.Equal(t, tt.wantClass, ErrorClass(tt.err))
		})
	}
}
