package plans

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/test"
)

func seedToken(t *testing.T) {
	t.Helper()
	keyring.MockInit()
	require.NoError(t, auth.SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken: "test-token",
		ExpiresIn:   3600,
		CreatedAt:   time.Now().Unix(),
	}))
}

func buildTemplate() dashboard.PlanTemplateResource {
	return dashboard.PlanTemplateResource{
		ID:   "build",
		Type: "plan_template",
		Attributes: dashboard.PlanTemplateAttributes{
			Name:          "Build",
			Description:   "Free forever Search & Discovery API.",
			Type:          "free",
			Configuration: dashboard.PlanTemplateConfiguration{Plan: "build"},
		},
	}
}

func growTemplate() dashboard.PlanTemplateResource {
	return dashboard.PlanTemplateResource{
		ID:   "grow",
		Type: "plan_template",
		Attributes: dashboard.PlanTemplateAttributes{
			Name:          "Grow",
			Description:   "Best-in-class Search & Discovery API.",
			Type:          "freeform",
			Freeform:      "$0.50 / 1,000 Requests",
			Configuration: dashboard.PlanTemplateConfiguration{Plan: "grow"},
		},
	}
}

func newServer(t *testing.T, freeOnly bool, userJSON string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/plan-templates/self-serve",
		func(w http.ResponseWriter, _ *http.Request) {
			data := []dashboard.PlanTemplateResource{buildTemplate()}
			if !freeOnly {
				data = append(data, growTemplate())
			}
			require.NoError(t, json.NewEncoder(w).Encode(dashboard.PlanTemplatesResponse{
				Data: data,
			}))
		},
	)
	mux.HandleFunc("/1/user", func(w http.ResponseWriter, _ *http.Request) {
		if userJSON == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		require.NoError(t, json.NewEncoder(w).Encode(json.RawMessage(userJSON)))
	})
	return httptest.NewServer(mux)
}

func newOpts(
	t *testing.T,
	srv *httptest.Server,
	isTTY bool,
	output string,
) (*PlansOptions, *test.CmdInOut) {
	t.Helper()
	seedToken(t)

	f, out := test.NewFactory(isTTY, nil, nil, "")
	pf := cmdutil.NewPrintFlags()
	*pf.OutputFormat = output
	pf.OutputFlagSpecified = func() bool { return output != "" }

	opts := &PlansOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: pf,
		NewDashboardClient: func(string) *dashboard.Client {
			c := dashboard.NewClientWithHTTPClient("test", srv.Client())
			c.APIURL = srv.URL
			return c
		},
	}
	return opts, out
}

func Test_runPlansCmd(t *testing.T) {
	srv := newServer(t, false, `{"has_payment_method": true}`)
	defer srv.Close()

	opts, out := newOpts(t, srv, true, "")
	require.NoError(t, runPlansCmd(opts))

	got := out.String()
	assert.Contains(t, got, "Build")
	assert.Contains(t, got, "Free")
	assert.Contains(t, got, "Grow")
	assert.Contains(t, got, "$0.50 / 1,000 Requests")
	assert.NotContains(t, got, "unavailable")
}

func Test_runPlansCmd_outputJSON(t *testing.T) {
	srv := newServer(t, false, `{"has_payment_method": true}`)
	defer srv.Close()

	opts, out := newOpts(t, srv, false, "json")
	require.NoError(t, runPlansCmd(opts))

	got := out.String()
	assert.Contains(t, got, `"name":"Build"`)
	assert.Contains(t, got, `"price":"Free"`)
	assert.Contains(t, got, `"name":"Grow"`)
	assert.Contains(t, got, `"available":true`)
	assert.NotContains(t, got, "unavailable_reason")
}

func Test_runPlansCmd_noPaymentMethod(t *testing.T) {
	srv := newServer(t, true, `{"has_payment_method": false}`)
	defer srv.Close()

	opts, out := newOpts(t, srv, true, "")
	require.NoError(t, runPlansCmd(opts))

	got := out.String()
	assert.Contains(t, got, "Build")
	assert.Contains(t, got, "Grow")
	assert.Contains(t, got, "Grow Plus")
	assert.Contains(t, got, "unavailable: no payment method on file — add billing to unlock")
}

func Test_runPlansCmd_noPaymentMethod_outputJSON(t *testing.T) {
	srv := newServer(t, true, `{"has_payment_method": false}`)
	defer srv.Close()

	opts, out := newOpts(t, srv, false, "json")
	require.NoError(t, runPlansCmd(opts))

	got := out.String()
	assert.Contains(t, got, `"name":"Build"`)
	assert.Contains(t, got, `"available":true`)
	assert.Contains(t, got, `"id":"grow"`)
	assert.Contains(t, got, `"id":"grow-plus"`)
	assert.Contains(t, got, `"available":false`)
	assert.Contains(t, got, `"unavailable_reason":"no payment method on file"`)
}

func Test_runPlansCmd_userFetchFailed(t *testing.T) {
	srv := newServer(t, true, "")
	defer srv.Close()

	opts, out := newOpts(t, srv, true, "")
	require.NoError(t, runPlansCmd(opts))

	got := out.String()
	assert.Contains(t, got, "Build")
	assert.NotContains(t, got, "Grow")
	assert.NotContains(t, got, "unavailable")
}
