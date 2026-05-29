package planchange

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/prompt"
	"github.com/algolia/cli/test"
)

// seedToken installs an in-memory keyring with a valid token so
// auth.EnsureAuthenticated short-circuits without hitting the network.
func seedToken(t *testing.T) {
	t.Helper()
	keyring.MockInit()
	require.NoError(t, auth.SaveToken(&dashboard.OAuthTokenResponse{
		AccessToken: "test-token",
		ExpiresIn:   3600,
		CreatedAt:   time.Now().Unix(),
	}))
}

func samplePlanTemplates() []dashboard.PlanTemplateResource {
	return []dashboard.PlanTemplateResource{
		{
			ID:   "build",
			Type: "plan_template",
			Attributes: dashboard.PlanTemplateAttributes{
				Name:        "Build",
				Description: "Free forever Search & Discovery API.",
				Type:        "free",
				Configuration: dashboard.PlanTemplateConfiguration{
					Plan:        "build",
					AcceptTerms: "Build terms",
				},
			},
		},
		{
			ID:   "grow",
			Type: "plan_template",
			Attributes: dashboard.PlanTemplateAttributes{
				Name:        "Grow",
				Description: "Best-in-class Search & Discovery API.",
				Type:        "freeform",
				Freeform:    "$0.50 / 1,000 Requests",
				Configuration: dashboard.PlanTemplateConfiguration{
					Plan:        "grow",
					AcceptTerms: "Grow terms",
				},
			},
		},
		{
			ID:   "grow-plus",
			Type: "plan_template",
			Attributes: dashboard.PlanTemplateAttributes{
				Name:        "Grow Plus",
				Description: "AI-powered Search & Discovery API.",
				Type:        "freeform",
				Freeform:    "$1.75 / 1,000 Requests",
				Configuration: dashboard.PlanTemplateConfiguration{
					Plan:        "grow-plus",
					AcceptTerms: "Grow Plus terms",
				},
			},
		},
	}
}

type planChangeServer struct {
	*httptest.Server
	patchCalls int
	lastPlan   string
}

// newServer spins up a dashboard stub. userJSON is the raw GET /1/user body;
// an empty string makes /1/user fail so the "user unavailable" fallback is used.
func newServer(t *testing.T, userJSON string) *planChangeServer {
	t.Helper()
	srv := &planChangeServer{}

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/plan-templates/self-serve",
		func(w http.ResponseWriter, _ *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode(dashboard.PlanTemplatesResponse{
				Data: samplePlanTemplates(),
			}))
		},
	)
	mux.HandleFunc("/1/user", func(w http.ResponseWriter, _ *http.Request) {
		if userJSON == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(userJSON))
	})
	mux.HandleFunc(
		"/1/applications/APP1/plan/self-serve",
		func(w http.ResponseWriter, r *http.Request) {
			srv.patchCalls++
			var payload dashboard.ChangePlanRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			srv.lastPlan = payload.Plan
			require.NoError(t, json.NewEncoder(w).Encode(dashboard.SingleApplicationResponse{
				Data: dashboard.ApplicationResource{
					ID: "APP1", Type: "application",
					Attributes: dashboard.ApplicationAttributes{
						ApplicationID: "APP1",
						Name:          "My App",
					},
				},
			}))
		},
	)

	srv.Server = httptest.NewServer(mux)
	return srv
}

func newPrintFlags(output string) *cmdutil.PrintFlags {
	pf := cmdutil.NewPrintFlags()
	*pf.OutputFormat = output
	pf.OutputFlagSpecified = func() bool { return output != "" }
	return pf
}

func newOpts(
	t *testing.T,
	srv *planChangeServer,
	isTTY bool,
) (*Options, *test.CmdInOut) {
	t.Helper()
	seedToken(t)
	t.Setenv("ALGOLIA_APPLICATION_ID", "APP1")

	f, out := test.NewFactory(isTTY, nil, nil, "")
	opts := &Options{
		IO:         f.IOStreams,
		Config:     f.Config,
		PrintFlags: newPrintFlags(""),
		NewDashboardClient: func(string) *dashboard.Client {
			c := dashboard.NewClientWithHTTPClient("test", srv.Client())
			c.APIURL = srv.URL
			return c
		},
	}
	return opts, out
}

func TestRun_WithPlanFlag(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	opts, out := newOpts(t, srv, false)
	opts.Plan = "grow"
	opts.AcceptTerms = true

	require.NoError(t, Run(opts))
	assert.Equal(t, 1, srv.patchCalls)
	assert.Equal(t, "grow", srv.lastPlan)
	assert.Contains(t, out.String(), "Grow")
}

func TestRun_FreeTargetNotBilled(t *testing.T) {
	// No payment method, but the free target must not be blocked.
	srv := newServer(t, `{"has_payment_method": false}`)
	defer srv.Close()

	opts, out := newOpts(t, srv, false)
	opts.Plan = "free"
	opts.AcceptTerms = true

	require.NoError(t, Run(opts))
	assert.Equal(t, 1, srv.patchCalls)
	// "free" maps to the free-type template, whose id is "build".
	assert.Equal(t, "build", srv.lastPlan)
	assert.Contains(t, out.String(), "Build")
}

func TestRun_BillingBlock(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": false}`)
	defer srv.Close()

	opts, _ := newOpts(t, srv, false)
	opts.Plan = "grow"
	opts.AcceptTerms = true

	err := Run(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "payment method")
	assert.Equal(t, 0, srv.patchCalls)
}

func TestRun_ToSDeclineAborts(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	defer prompt.StubConfirm(false)()

	opts, out := newOpts(t, srv, true)
	opts.Plan = "grow"

	require.NoError(t, Run(opts))
	assert.Equal(t, 0, srv.patchCalls)
	assert.Contains(t, out.String(), "aborted")
}

func TestRun_NonInteractiveRequiresPlan(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	opts, _ := newOpts(t, srv, false)
	// No --plan and no TTY.

	err := Run(opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--plan is required")
	assert.Equal(t, 0, srv.patchCalls)
}

func TestRun_InteractivePicker(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	// Pick the second candidate (grow); the picker lists all plans in order.
	origAsk := prompt.SurveyAskOne
	prompt.SurveyAskOne = func(_ survey.Prompt, response interface{}, _ ...survey.AskOpt) error {
		*(response.(*int)) = 1
		return nil
	}
	t.Cleanup(func() { prompt.SurveyAskOne = origAsk })
	defer prompt.StubConfirm(true)()

	opts, _ := newOpts(t, srv, true)

	require.NoError(t, Run(opts))
	assert.Equal(t, 1, srv.patchCalls)
	assert.Equal(t, "grow", srv.lastPlan)
}

func TestRun_DryRunDoesNotCallAPI(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	opts, out := newOpts(t, srv, false)
	opts.Plan = "grow"
	opts.DryRun = true

	require.NoError(t, Run(opts))
	assert.Equal(t, 0, srv.patchCalls)
	assert.Contains(t, out.String(), "Dry run")
	assert.Contains(t, out.String(), "Grow")
}

func TestRun_OutputJSON(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	opts, out := newOpts(t, srv, false)
	opts.Plan = "grow"
	opts.AcceptTerms = true
	opts.PrintFlags = newPrintFlags("json")

	require.NoError(t, Run(opts))
	assert.Equal(t, 1, srv.patchCalls)
	assert.Contains(t, out.String(), `"plan":"grow"`)
	assert.Contains(t, out.String(), `"application_id":"APP1"`)
}
