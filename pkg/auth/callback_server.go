package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

const callbackShutdownTimeout = 5 * time.Second

const successHTML = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Algolia CLI</title>
<style>
  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    display: flex; flex-direction: column;
    justify-content: center; align-items: center;
    height: 100vh; margin: 0;
    background: #f5f5fa; color: #21243d;
    position: relative;
  }
  .brand {
    font-size: 1.5rem; font-weight: 700; color: #003dff;
    margin-bottom: 1.5rem;
  }
  .card {
    margin-top: 0.25rem;
    background: #fff;
    border-radius: 4px;
    box-shadow: 0 0 0 1px rgba(35,38,59,.05),0 1px 3px 0 rgba(35,38,59,.15);
    overflow: hidden;
    min-width: 400px;
  }
  .card-header {
    text-align: center; font-size: 1.25rem;
    padding: 1.5rem 2rem 1rem;
  }
  .card-header h2 { margin: 0; font-size: 1.25rem; font-weight: 600; }
  .card-body {
    padding: 0 2rem 1.5rem;
    text-align: center;
    color: #5a5e9a; font-size: 0.95rem;
  }
  .card-body p { margin: 0; }
</style>
</head>
<body>
<h1>Algolia CLI</h1>
<div class="card">
  <div class="card-header">
    <h2>Authentication successful</h2>
  </div>
  <div class="card-body">
    <p>You can close this tab and return to your terminal.</p>
  </div>
</div>
</body>
</html>`

// CallbackResult holds the authorization code (or an error description)
// received from the OAuth redirect.
type CallbackResult struct {
	Code  string
	Error string
}

// StartCallbackServer starts a local HTTP server on a random available port.
// It returns the redirect URI (http://127.0.0.1:{port}) and a channel that
// will receive exactly one CallbackResult when the OAuth redirect arrives.
// The server shuts itself down after handling the first request.
func StartCallbackServer() (redirectURI string, result <-chan CallbackResult, err error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, fmt.Errorf("failed to start callback server: %w", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI = fmt.Sprintf("http://127.0.0.1:%d", port)

	ch := make(chan CallbackResult, 1)

	mux := http.NewServeMux()
	srv := &http.Server{Handler: mux}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if errParam := q.Get("error"); errParam != "" {
			desc := q.Get("error_description")
			if desc == "" {
				desc = errParam
			}
			http.Error(w, desc, http.StatusBadRequest)
			ch <- CallbackResult{Error: desc}
		} else {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, successHTML)
			ch <- CallbackResult{Code: q.Get("code")}
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), callbackShutdownTimeout)
			defer cancel()
			_ = srv.Shutdown(ctx)
		}()
	})

	go func() { _ = srv.Serve(listener) }()

	return redirectURI, ch, nil
}
