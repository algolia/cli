package agentstudio

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// sdkPathPrefix is the path prefix the official SDK hard-codes on every
// request. The local client appends "/1/..." to a configurable base URL
// instead, so SDKRequester maps one onto the other.
const sdkPathPrefix = "/agent-studio"

// SDKRequester is a transport.Requester for the official Agent Studio SDK. It
//
//   - rewrites the SDK's fixed "/agent-studio/1/..." paths onto the same base
//     URL the local client uses, so both clients always target the same backend;
//   - never follows redirects, so API keys and request bodies cannot be
//     forwarded to another host (or downgraded from https to http);
//   - turns every non-2xx body into a small {"message": ...} document built
//     from the error detail only, so values the backend echoes back (such as
//     a submitted provider apiKey) never reach the terminal.
type SDKRequester struct {
	client   *http.Client
	basePath string
}

// NewSDKRequester returns a requester for the given resolved base URL (as
// returned by ResolveHost). Only the path of baseURL is used here: the scheme
// and host belong to the SDK's host list (see SplitBaseURL).
func NewSDKRequester(baseURL string) *SDKRequester {
	basePath := ""
	if u, err := url.Parse(baseURL); err == nil {
		basePath = strings.TrimRight(u.Path, "/")
	}
	return &SDKRequester{client: newNoRedirectHTTPClient(), basePath: basePath}
}

// SplitBaseURL returns the scheme and host (including any port) of a resolved
// base URL.
func SplitBaseURL(baseURL string) (scheme, host string, err error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", "", err
	}
	return u.Scheme, u.Host, nil
}

// Request implements transport.Requester. The SDK applies its own per-call
// timeouts through the request context.
func (r *SDKRequester) Request(req *http.Request, _, _ time.Duration) (*http.Response, error) {
	req.URL.Path = r.basePath + strings.TrimPrefix(req.URL.Path, sdkPathPrefix)
	req.URL.RawPath = ""

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err //nolint:wrapcheck // surfaced as-is by the SDK transport
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp, nil
	}

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	_ = resp.Body.Close()
	clean := SanitizeErrorBody(resp.StatusCode, raw)
	resp.Body = io.NopCloser(bytes.NewReader(clean))
	resp.ContentLength = int64(len(clean))
	resp.Header.Set("Content-Type", "application/json")
	return resp, nil
}

// SanitizeErrorBody builds a {"message": ...} error body from a raw
// non-2xx response. The message is the structured detail (FastAPI
// detail[].msg, string detail or Algolia message), followed by an actionable
// hint for authentication, permission and feature-flag failures. Echoed
// request values (detail[].input) are never included.
func SanitizeErrorBody(status int, body []byte) []byte {
	msg := extractDetail(body)
	if sentinel := sentinelFor(status, body); sentinel != nil {
		if msg == "" {
			msg = sentinel.Error()
		} else {
			msg += " (" + sentinel.Error() + ")"
		}
	}
	if msg == "" {
		msg = http.StatusText(status)
	}
	out, _ := json.Marshal(map[string]string{"message": msg})
	return out
}

func newNoRedirectHTTPClient() *http.Client {
	return &http.Client{
		Transport: newHTTPTransport(),
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func newHTTPTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
