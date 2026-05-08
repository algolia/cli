package config

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmd/agents/sharedtest"
	"github.com/algolia/cli/test"
)

// ---------------------------------------------------------------------
// get
// ---------------------------------------------------------------------

func Test_runGetCmd_PrintsConfig(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/configuration", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(`{"maxRetentionDays":60}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewConfigCmd(f)
	result, err := test.Execute(cmd, "get", out)
	require.NoError(t, err)

	var got map[string]int
	require.NoError(t, json.Unmarshal([]byte(result.String()), &got))
	assert.Equal(t, 60, got["maxRetentionDays"])
}

// ---------------------------------------------------------------------
// set
// ---------------------------------------------------------------------

func Test_runSetCmd_RetentionDays_BuildsBodyAndPATCHes(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/configuration", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"maxRetentionDays":30}`, string(body))
		_, _ = w.Write([]byte(`{"maxRetentionDays":30}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewConfigCmd(f)
	_, err := test.Execute(cmd, "set --retention-days 30", out)
	require.NoError(t, err)
}

func Test_runSetCmd_File_RoundTripsBody(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/configuration", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		assert.JSONEq(t, `{"maxRetentionDays":90,"futureField":"x"}`, string(body))
		_, _ = w.Write([]byte(`{"maxRetentionDays":90}`))
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	patchPath := sharedtest.WriteTempJSON(t, "patch.json", `{"maxRetentionDays":90,"futureField":"x"}`)
	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewConfigCmd(f)
	_, err := test.Execute(cmd, "set -F "+patchPath, out)
	require.NoError(t, err)
}

func Test_runSetCmd_DryRunSkipsAPI(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/configuration", func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("backend was called during --dry-run")
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewConfigCmd(f)
	result, err := test.Execute(cmd, "set --retention-days 30 --dry-run", out)
	require.NoError(t, err)
	got := result.String()
	assert.Contains(t, got, "Dry run: would PATCH /1/configuration")
	assert.Contains(t, got, `"maxRetentionDays": 30`)
}

func Test_runSetCmd_RejectsNeitherFlag(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewConfigCmd(f)
	_, err := test.Execute(cmd, "set", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "one of --retention-days or --file is required")
}

func Test_runSetCmd_RejectsBothFlags(t *testing.T) {
	patchPath := sharedtest.WriteTempJSON(t, "patch.json", `{"maxRetentionDays":30}`)
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := NewConfigCmd(f)
	_, err := test.Execute(cmd, "set --retention-days 30 -F "+patchPath, out)
	require.Error(t, err)
	// cobra's mutually-exclusive guard fires; mirrors the "[input message]"
	// assertion style from other agents commands.
	assert.Contains(t, err.Error(), "[file retention-days]")
}

func Test_runSetCmd_PropagatesValidationError(t *testing.T) {
	// Backend rejects retention values not in [0, 30, 60, 90].
	mux := http.NewServeMux()
	mux.HandleFunc("/1/configuration", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write(
			[]byte(
				`{"detail":[{"msg":"maxRetentionDays must be one of [0, 30, 60, 90]","loc":["body","maxRetentionDays"]}]}`,
			),
		)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	f, out := test.NewFactory(false, nil, nil, "")
	f.AgentStudioClient = sharedtest.NewClient(t, ts)

	cmd := NewConfigCmd(f)
	_, err := test.Execute(cmd, "set --retention-days 45", out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be one of")
}
