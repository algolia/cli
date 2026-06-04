package create

import (
	"context"
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

// seedToken installs an in-memory keyring with a valid token.
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

type createServer struct {
	*httptest.Server
	createCalls int
	patchCalls  int
	lastPlan    string
	failPatch   bool
	// freeOnly returns only the free plan, mirroring the API when no billing is on file.
	freeOnly bool
}

// newServer spins up a dashboard stub. An empty userJSON makes /1/user fail.
func newServer(t *testing.T, userJSON string) *createServer {
	t.Helper()
	srv := &createServer{}

	appResponse := dashboard.SingleApplicationResponse{
		Data: dashboard.ApplicationResource{
			ID:   "APP1",
			Type: "application",
			Attributes: dashboard.ApplicationAttributes{
				ApplicationID: "APP1",
				Name:          "My App",
			},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/plan-templates/self-serve",
		func(w http.ResponseWriter, _ *http.Request) {
			templates := samplePlanTemplates()
			if srv.freeOnly {
				templates = templates[:1] // only the free "build" template
			}
			require.NoError(t, json.NewEncoder(w).Encode(dashboard.PlanTemplatesResponse{
				Data: templates,
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
	mux.HandleFunc("/1/hosting/regions", func(w http.ResponseWriter, _ *http.Request) {
		require.NoError(t, json.NewEncoder(w).Encode(dashboard.RegionsResponse{
			RegionCodes: []dashboard.Region{{Code: "CA", Name: "Canada"}},
		}))
	})
	mux.HandleFunc("/1/applications", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		srv.createCalls++
		require.NoError(t, json.NewEncoder(w).Encode(appResponse))
	})
	mux.HandleFunc(
		"/1/applications/APP1/api-keys",
		func(w http.ResponseWriter, _ *http.Request) {
			require.NoError(t, json.NewEncoder(w).Encode(dashboard.CreateAPIKeyResponse{
				Data: dashboard.APIKeyResource{
					ID:         "key",
					Type:       "key",
					Attributes: dashboard.APIKeyAttributes{Value: "test-api-key"},
				},
			}))
		},
	)
	mux.HandleFunc(
		"/1/applications/APP1/plan/self-serve",
		func(w http.ResponseWriter, r *http.Request) {
			srv.patchCalls++
			if srv.failPatch {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var payload dashboard.ChangePlanRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			srv.lastPlan = payload.Plan
			require.NoError(t, json.NewEncoder(w).Encode(appResponse))
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

// newOpts builds CreateOptions wired to the stub server, defaulting to JSON output.
func newOpts(
	t *testing.T,
	srv *createServer,
	isTTY bool,
) (*CreateOptions, *test.CmdInOut, *string) {
	t.Helper()
	seedToken(t)

	f, out := test.NewFactory(isTTY, nil, nil, "")
	opened := new(string)
	opts := &CreateOptions{
		IO:           f.IOStreams,
		Config:       f.Config,
		Name:         "My First Application",
		Region:       "CA",
		nameProvided: true,
		PrintFlags:   newPrintFlags("json"),
		NewDashboardClient: func(string) *dashboard.Client {
			c := dashboard.NewClientWithHTTPClient("test", srv.Client())
			c.APIURL = srv.URL
			c.DashboardURL = "https://dashboard.algolia.com"
			return c
		},
		Browser: func(url string) error {
			*opened = url
			return nil
		},
	}
	return opts, out, opened
}

func TestRun_FreeNonInteractive(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": false}`)
	defer srv.Close()

	opts, out, _ := newOpts(t, srv, false)
	opts.AcceptTerms = true

	require.NoError(t, runCreateCmd(context.Background(), opts))
	assert.Equal(t, 1, srv.createCalls)
	assert.Equal(t, 0, srv.patchCalls)
	assert.Contains(t, out.String(), "APP1")
}

func TestRun_NonInteractiveRequiresAcceptTerms(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	opts, _, _ := newOpts(t, srv, false)

	err := runCreateCmd(context.Background(), opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be accepted")
	assert.Equal(t, 0, srv.createCalls)
}

func TestRun_PaidWithBillingNonInteractive(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	opts, out, _ := newOpts(t, srv, false)
	opts.Plan = "grow"
	opts.AcceptTerms = true

	require.NoError(t, runCreateCmd(context.Background(), opts))
	assert.Equal(t, 1, srv.createCalls)
	assert.Equal(t, 1, srv.patchCalls)
	assert.Equal(t, "grow", srv.lastPlan)
	assert.Contains(t, out.String(), "APP1")
}

func TestRun_PaidWithBillingRequiresAcceptTerms(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	opts, _, _ := newOpts(t, srv, false)
	opts.Plan = "grow"

	err := runCreateCmd(context.Background(), opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be accepted")
	assert.Equal(t, 0, srv.createCalls)
	assert.Equal(t, 0, srv.patchCalls)
}

func TestRun_PaidNoBillingNonInteractive(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": false}`)
	defer srv.Close()

	opts, out, _ := newOpts(t, srv, false)
	opts.Plan = "grow"
	opts.AcceptTerms = true

	err := runCreateCmd(context.Background(), opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "payment method")
	assert.Equal(t, 0, srv.createCalls)
	assert.Equal(t, 0, srv.patchCalls)
	assert.Contains(t, out.String(), "https://dashboard.algolia.com/account/billing/details")
}

func TestRun_PaidNoBillingInteractiveOpensBilling(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": false}`)
	defer srv.Close()

	defer prompt.StubConfirm(true)()

	opts, _, opened := newOpts(t, srv, true)
	opts.Plan = "grow"

	require.NoError(t, runCreateCmd(context.Background(), opts))
	assert.Equal(t, 0, srv.createCalls)
	assert.Equal(
		t,
		"https://dashboard.algolia.com/account/billing/details",
		*opened,
	)
}

func TestRun_PaidNoBillingInteractiveDeclineOpen(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": false}`)
	defer srv.Close()

	defer prompt.StubConfirm(false)()

	opts, _, opened := newOpts(t, srv, true)
	opts.Plan = "grow"

	require.NoError(t, runCreateCmd(context.Background(), opts))
	assert.Equal(t, 0, srv.createCalls)
	assert.Empty(t, *opened)
}

func TestRun_ToSDeclineAborts(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	defer prompt.StubConfirm(false)()

	opts, out, _ := newOpts(t, srv, true)
	opts.Plan = "free"

	require.NoError(t, runCreateCmd(context.Background(), opts))
	assert.Equal(t, 0, srv.createCalls)
	assert.Contains(t, out.String(), "Aborted")
}

func TestRun_AcceptTermsSkipsPromptInteractive(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	// Confirm stubbed to NO; --accept-terms must bypass the prompt.
	defer prompt.StubConfirm(false)()

	opts, out, _ := newOpts(t, srv, true)
	opts.Plan = "free"
	opts.AcceptTerms = true

	require.NoError(t, runCreateCmd(context.Background(), opts))
	assert.Equal(t, 1, srv.createCalls)
	assert.Contains(t, out.String(), "Terms accepted via --accept-terms")
}

func TestRun_PaidPlanHiddenByServerNonInteractive(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": false}`)
	srv.freeOnly = true
	defer srv.Close()

	opts, out, _ := newOpts(t, srv, false)
	opts.Plan = "grow"
	opts.AcceptTerms = true

	err := runCreateCmd(context.Background(), opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "payment method")
	assert.NotContains(t, err.Error(), "invalid plan")
	assert.Equal(t, 0, srv.createCalls)
	assert.Equal(t, 0, srv.patchCalls)
	assert.Contains(t, out.String(), "https://dashboard.algolia.com/account/billing/details")
	assert.Contains(t, out.String(), "--plan grow")
}

func TestRun_PaidPlanHiddenByServerInteractiveOpensBilling(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": false}`)
	srv.freeOnly = true
	defer srv.Close()

	defer prompt.StubConfirm(true)()

	opts, _, opened := newOpts(t, srv, true)
	opts.Plan = "grow"

	require.NoError(t, runCreateCmd(context.Background(), opts))
	assert.Equal(t, 0, srv.createCalls)
	assert.Equal(
		t,
		"https://dashboard.algolia.com/account/billing/details",
		*opened,
	)
}

func TestRun_InvalidPlanErrors(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	opts, _, _ := newOpts(t, srv, false)
	opts.Plan = "bogus"
	opts.AcceptTerms = true

	err := runCreateCmd(context.Background(), opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid plan")
	assert.Equal(t, 0, srv.createCalls)
}

func TestRun_InteractivePickerHidesPaidWithoutBilling(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": false}`)
	defer srv.Close()

	defer prompt.StubConfirm(true)()

	opts, out, _ := newOpts(t, srv, true)

	require.NoError(t, runCreateCmd(context.Background(), opts))
	assert.Equal(t, 1, srv.createCalls)
	assert.Equal(t, 0, srv.patchCalls)
	assert.Contains(t, out.String(), "only the Free plan is available")
	assert.Contains(t, out.String(), "APP1")
}

func TestRun_InteractivePickerSelectsPaid(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	origAsk := prompt.SurveyAskOne
	prompt.SurveyAskOne = func(_ survey.Prompt, response interface{}, _ ...survey.AskOpt) error {
		*(response.(*int)) = 1
		return nil
	}
	t.Cleanup(func() { prompt.SurveyAskOne = origAsk })
	defer prompt.StubConfirm(true)()

	opts, out, _ := newOpts(t, srv, true)

	require.NoError(t, runCreateCmd(context.Background(), opts))
	assert.Equal(t, 1, srv.createCalls)
	assert.Equal(t, 1, srv.patchCalls)
	assert.Equal(t, "grow", srv.lastPlan)
	assert.Contains(t, out.String(), "APP1")
}

func TestRun_DryRunDoesNotCallAPI(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	defer srv.Close()

	opts, out, _ := newOpts(t, srv, false)
	opts.Plan = "grow"
	opts.DryRun = true
	opts.PrintFlags = newPrintFlags("")

	require.NoError(t, runCreateCmd(context.Background(), opts))
	assert.Equal(t, 0, srv.createCalls)
	assert.Equal(t, 0, srv.patchCalls)
	assert.Contains(t, out.String(), "Dry run")
	assert.Contains(t, out.String(), "grow")
}

func TestRun_PlanChangeFailureKeepsFreeApp(t *testing.T) {
	srv := newServer(t, `{"has_payment_method": true}`)
	srv.failPatch = true
	defer srv.Close()

	opts, _, _ := newOpts(t, srv, false)
	opts.Plan = "grow"
	opts.AcceptTerms = true

	err := runCreateCmd(context.Background(), opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to apply")
	assert.Equal(t, 1, srv.createCalls)
	assert.Equal(t, 1, srv.patchCalls)
}

func TestResolveName(t *testing.T) {
	t.Run("explicit flag wins", func(t *testing.T) {
		f, _ := test.NewFactory(true, nil, nil, "")
		opts := &CreateOptions{IO: f.IOStreams, Name: "Explicit", nameProvided: true}
		name, err := resolveName(opts)
		require.NoError(t, err)
		assert.Equal(t, "Explicit", name)
	})

	t.Run("interactive prompt returns entered value", func(t *testing.T) {
		f, _ := test.NewFactory(true, nil, nil, "")
		origAsk := prompt.SurveyAskOne
		prompt.SurveyAskOne = func(_ survey.Prompt, response interface{}, _ ...survey.AskOpt) error {
			*(response.(*string)) = "Typed Name"
			return nil
		}
		t.Cleanup(func() { prompt.SurveyAskOne = origAsk })

		opts := &CreateOptions{IO: f.IOStreams, Name: "My First Application"}
		name, err := resolveName(opts)
		require.NoError(t, err)
		assert.Equal(t, "Typed Name", name)
	})

	t.Run("empty interactive input falls back to default", func(t *testing.T) {
		f, _ := test.NewFactory(true, nil, nil, "")
		origAsk := prompt.SurveyAskOne
		prompt.SurveyAskOne = func(_ survey.Prompt, response interface{}, _ ...survey.AskOpt) error {
			*(response.(*string)) = ""
			return nil
		}
		t.Cleanup(func() { prompt.SurveyAskOne = origAsk })

		opts := &CreateOptions{IO: f.IOStreams, Name: "My First Application"}
		name, err := resolveName(opts)
		require.NoError(t, err)
		assert.Equal(t, "My First Application", name)
	})

	t.Run("non-interactive falls back to default", func(t *testing.T) {
		f, _ := test.NewFactory(false, nil, nil, "")
		opts := &CreateOptions{IO: f.IOStreams, Name: "My First Application"}
		name, err := resolveName(opts)
		require.NoError(t, err)
		assert.Equal(t, "My First Application", name)
	})
}
