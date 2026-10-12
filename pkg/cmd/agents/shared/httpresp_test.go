package shared

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type closeTracker struct {
	io.Reader
	closed bool
}

func (c *closeTracker) Close() error {
	c.closed = true
	return nil
}

func TestRawResponse_ClosesBody(t *testing.T) {
	cases := []struct {
		name   string
		status int
		err    error
	}{
		{"success", http.StatusOK, nil},
		{"http error", http.StatusNotFound, nil},
		{"transport error", http.StatusOK, errors.New("boom")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := &closeTracker{Reader: strings.NewReader("{}")}
			res := &http.Response{StatusCode: tc.status, Body: body}
			_, _ = RawResponse(res, []byte("{}"), tc.err)
			assert.True(t, body.closed)
		})
	}
}

func TestRawResponse_NilResponse(t *testing.T) {
	got, err := RawResponse(nil, []byte(`{"a":1}`), nil)
	require.NoError(t, err)
	assert.JSONEq(t, `{"a":1}`, string(got))
}
