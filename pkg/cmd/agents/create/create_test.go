package create

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

func Test_runCreateCmd_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"name":"Concierge","instructions":"x"}`, string(body))

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id":"abc-123","name":"Concierge","status":"draft",
			"instructions":"x","createdAt":"2025-01-01T00:00:00Z"
		}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, `{"name":"Concierge","instructions":"x"}`)
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewCreateCmd(f, nil)
	result, err := test.Execute(cmd, "-F -", out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"id":"abc-123"`)
}

func Test_runCreateCmd_RejectsInvalidJSON(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, `{not json`)
	cmd := NewCreateCmd(f, nil)
	_, err := test.Execute(cmd, "-F -", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not valid JSON")
}

func Test_runCreateCmd_RejectsInvalidBodyJSON(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewCreateCmd(f, nil)
	_, err := test.Execute(cmd, `--body '{not json}'`, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--body is not valid JSON")
}

func Test_runCreateCmd_FromBody(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/agents", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"name":"Concierge","instructions":"x"}`, string(body))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"abc-123","name":"Concierge"}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewCreateCmd(f, nil)
	result, err := test.Execute(cmd, `--body '{"name":"Concierge","instructions":"x"}'`, out)
	require.NoError(t, err)
	assert.Contains(t, result.String(), `"id":"abc-123"`)
}

func Test_runCreateCmd_RejectsBodyAndFile(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewCreateCmd(f, nil)
	_, err := test.Execute(cmd, `--body {} -F -`, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not both")
}

func Test_runCreateCmd_RequiresBodyOrFile(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewCreateCmd(f, nil)
	_, err := test.Execute(cmd, "", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "one of --body or --file is required")
}
