package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(handler http.Handler) (*httptest.Server, *Client) {
	ts := httptest.NewServer(handler)
	client := NewClientWithHTTPClient("test-client-id", ts.Client())
	client.DashboardURL = ts.URL
	client.APIURL = ts.URL
	client.OAuthScope = "scope:test"
	return ts, client
}

func TestAuthorizeURL(t *testing.T) {
	client := &Client{
		DashboardURL: "https://dashboard.example.com",
		ClientID:     "my-client-id",
		OAuthScope:   "scope:test",
	}

	url := client.AuthorizeURL("test-challenge", "http://localhost:12345")
	assert.Contains(t, url, "https://dashboard.example.com/2/oauth/authorize?")
	assert.Contains(t, url, "client_id=my-client-id")
	assert.Contains(t, url, "response_type=code")
	assert.Contains(t, url, "code_challenge=test-challenge")
	assert.Contains(t, url, "code_challenge_method=S256")
	assert.Contains(t, url, "redirect_uri=http")
}

func TestAuthorizationCodeGrant_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/2/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		require.NoError(t, r.ParseForm())
		assert.Equal(t, "authorization_code", r.FormValue("grant_type"))
		assert.Equal(t, "test-client-id", r.FormValue("client_id"))
		assert.Equal(t, "auth-code-123", r.FormValue("code"))
		assert.Equal(t, "verifier-xyz", r.FormValue("code_verifier"))

		require.NoError(t, json.NewEncoder(w).Encode(OAuthTokenResponse{
			AccessToken:  "access-token-123",
			RefreshToken: "refresh-token-456",
			TokenType:    "Bearer",
			ExpiresIn:    7200,
			Scope:        "scope:test",
			User: &User{
				ID:    1,
				Email: "user@test.com",
				Name:  "Test User",
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	resp, err := client.AuthorizationCodeGrant(
		"auth-code-123",
		"verifier-xyz",
		"http://localhost:12345",
	)
	require.NoError(t, err)
	assert.Equal(t, "access-token-123", resp.AccessToken)
	assert.Equal(t, "refresh-token-456", resp.RefreshToken)
	assert.Equal(t, "user@test.com", resp.User.Email)
}

func TestAuthorizationCodeGrant_InvalidGrant(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/2/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		require.NoError(t, json.NewEncoder(w).Encode(OAuthErrorResponse{
			Error:            "invalid_grant",
			ErrorDescription: "Authorization code has expired.",
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.AuthorizationCodeGrant("expired-code", "verifier", "http://localhost:12345")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Authorization code has expired")
}

func TestRefreshToken_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/2/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseForm())
		assert.Equal(t, "refresh_token", r.FormValue("grant_type"))
		assert.Equal(t, "old-refresh-token", r.FormValue("refresh_token"))

		require.NoError(t, json.NewEncoder(w).Encode(OAuthTokenResponse{
			AccessToken:  "new-access-token",
			RefreshToken: "new-refresh-token",
			ExpiresIn:    7200,
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	resp, err := client.RefreshToken("old-refresh-token")
	require.NoError(t, err)
	assert.Equal(t, "new-access-token", resp.AccessToken)
	assert.Equal(t, "new-refresh-token", resp.RefreshToken)
}

func TestListApplications_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		require.NoError(t, json.NewEncoder(w).Encode(ApplicationsResponse{
			Data: []ApplicationResource{
				{
					ID:   "APP1",
					Type: "application",
					Attributes: ApplicationAttributes{
						ApplicationID: "APP1",
						Name:          "My App",
						APIKey:        "key1",
					},
				},
				{
					ID:   "APP2",
					Type: "application",
					Attributes: ApplicationAttributes{
						ApplicationID: "APP2",
						Name:          "Other App",
						APIKey:        "key2",
					},
				},
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	apps, err := client.ListApplications("test-token")
	require.NoError(t, err)
	assert.Len(t, apps, 2)
	assert.Equal(t, "APP1", apps[0].ID)
}

func TestListApplications_Paginated(t *testing.T) {
	mux := http.NewServeMux()
	callCount := 0

	mux.HandleFunc("/1/applications", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		page := r.URL.Query().Get("page")

		switch page {
		case "2":
			require.NoError(t, json.NewEncoder(w).Encode(ApplicationsResponse{
				Data: []ApplicationResource{
					{
						ID:         "APP3",
						Type:       "application",
						Attributes: ApplicationAttributes{ApplicationID: "APP3", Name: "Third App"},
					},
				},
				Meta: PaginationMeta{TotalCount: 3, PerPage: 2, CurrentPage: 2, TotalPages: 2},
			}))
		default:
			require.NoError(t, json.NewEncoder(w).Encode(ApplicationsResponse{
				Data: []ApplicationResource{
					{
						ID:         "APP1",
						Type:       "application",
						Attributes: ApplicationAttributes{ApplicationID: "APP1", Name: "First App"},
					},
					{
						ID:   "APP2",
						Type: "application",
						Attributes: ApplicationAttributes{
							ApplicationID: "APP2",
							Name:          "Second App",
						},
					},
				},
				Meta: PaginationMeta{TotalCount: 3, PerPage: 2, CurrentPage: 1, TotalPages: 2},
			}))
		}
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	apps, err := client.ListApplications("test-token")
	require.NoError(t, err)
	assert.Len(t, apps, 3)
	assert.Equal(t, "APP1", apps[0].ID)
	assert.Equal(t, "APP2", apps[1].ID)
	assert.Equal(t, "APP3", apps[2].ID)
	assert.Equal(t, 2, callCount)
}

func TestListApplications_Unauthorized(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.ListApplications("expired-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session expired")
}

func TestGetApplication_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/application/APP1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		require.NoError(t, json.NewEncoder(w).Encode(SingleApplicationResponse{
			Data: ApplicationResource{
				ID: "APP1", Type: "application",
				Attributes: ApplicationAttributes{
					ApplicationID: "APP1",
					Name:          "My App",
					APIKey:        "api-key-123",
				},
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	app, err := client.GetApplication("test-token", "APP1")
	require.NoError(t, err)
	assert.Equal(t, "APP1", app.ID)
	assert.Equal(t, "api-key-123", app.APIKey)
}

func TestGetApplication_ParsesPlanLabel(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/application/APP1", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(
			`{"data":{"id":"APP1","type":"application","attributes":{"name":"My App","application_id":"APP1","plan":{"name":"v8.5-plg-grow-plus","version":9,"label":"Grow Plus","pay_as_you_go":true}}}}`,
		))
		require.NoError(t, err)
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	app, err := client.GetApplication("test-token", "APP1")
	require.NoError(t, err)
	assert.Equal(t, "Grow Plus", app.PlanLabel)
}

func TestCreateApplication_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload CreateApplicationRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Equal(t, "us", payload.RegionCode)
		assert.Equal(t, "My App", payload.Name)

		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(SingleApplicationResponse{
			Data: ApplicationResource{
				ID: "NEW_APP", Type: "application",
				Attributes: ApplicationAttributes{ApplicationID: "NEW_APP", Name: "My App"},
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	app, err := client.CreateApplication("test-token", "us", "My App")
	require.NoError(t, err)
	assert.Equal(t, "NEW_APP", app.ID)
	assert.Equal(t, "My App", app.Name)
}

func TestCreateAPIKey_ReturnsValueAndUUID(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications/APP1/api-keys", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(CreateAPIKeyResponse{
			Data: APIKeyResource{
				ID:         "key-uuid-123",
				Type:       "api_key",
				Attributes: APIKeyAttributes{Value: "secret-key"},
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	created, err := client.CreateAPIKey("test-token", "APP1", WriteACL, "Algolia CLI")
	require.NoError(t, err)
	assert.Equal(t, "secret-key", created.Value)
	assert.Equal(t, "key-uuid-123", created.UUID)
}

func TestCreateAPIKey_EmptyValueReturnsError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications/APP1/api-keys", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(CreateAPIKeyResponse{
			Data: APIKeyResource{
				ID:         "key-uuid-123",
				Attributes: APIKeyAttributes{Value: ""},
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.CreateAPIKey("test-token", "APP1", WriteACL, "Algolia CLI")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no key was returned")
}

func TestUpdateApplication_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications/APP1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload UpdateApplicationRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Equal(t, "Renamed App", payload.Name)

		require.NoError(t, json.NewEncoder(w).Encode(SingleApplicationResponse{
			Data: ApplicationResource{
				ID: "APP1", Type: "application",
				Attributes: ApplicationAttributes{ApplicationID: "APP1", Name: "Renamed App"},
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	app, err := client.UpdateApplication("test-token", "APP1", "Renamed App")
	require.NoError(t, err)
	assert.Equal(t, "APP1", app.ID)
	assert.Equal(t, "Renamed App", app.Name)
}

func TestUpdateApplication_Unauthorized(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications/APP1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.UpdateApplication("expired-token", "APP1", "Renamed App")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "session expired")
}

func TestUpdateApplication_ErrorStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/applications/APP1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err := w.Write([]byte(`{"errors":[{"title":"name has already been taken"}]}`))
		require.NoError(t, err)
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.UpdateApplication("test-token", "APP1", "Taken Name")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "update application failed with status 422")
	assert.Contains(t, err.Error(), "name has already been taken")
}

func TestGetSelfServePlans_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/plan-templates/self-serve", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		require.NoError(t, json.NewEncoder(w).Encode(PlanTemplatesResponse{
			Data: []PlanTemplateResource{
				{
					ID:   "build",
					Type: "plan_template",
					Attributes: PlanTemplateAttributes{
						Name:        "Build",
						Description: "Free forever Search & Discovery API.",
						Type:        "free",
						Configuration: PlanTemplateConfiguration{
							Plan:        "build",
							AcceptTerms: "Build terms",
						},
					},
				},
				{
					ID:   "grow",
					Type: "plan_template",
					Attributes: PlanTemplateAttributes{
						Name:        "Grow",
						Description: "Best-in-class Search & Discovery API with free tier.",
						Type:        "freeform",
						Freeform:    "$0.50 / 1,000 Requests",
						Configuration: PlanTemplateConfiguration{
							Plan:        "grow",
							AcceptTerms: "Grow terms",
						},
					},
				},
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	plans, err := client.GetSelfServePlans("test-token")
	require.NoError(t, err)
	require.Len(t, plans, 2)

	// The free-type plan keeps its configuration.plan id and prices as "Free".
	assert.Equal(t, "build", plans[0].ID)
	assert.Equal(t, "Build", plans[0].Name)
	assert.Equal(t, "free", plans[0].Type)
	assert.Equal(t, "Free", plans[0].Price)
	assert.Equal(t, "Build terms", plans[0].AcceptTerms)

	// Freeform plans surface the freeform pricing string.
	assert.Equal(t, "grow", plans[1].ID)
	assert.Equal(t, "$0.50 / 1,000 Requests", plans[1].Price)
}

func TestGetSelfServePlans_Unauthorized(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/plan-templates/self-serve", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.GetSelfServePlans("expired-token")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "session expired")
}

func TestGetUser_TopLevel(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/user", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, err := w.Write([]byte(`{"has_payment_method": true, "plan": "grow"}`))
		require.NoError(t, err)
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	user, err := client.GetUser("test-token")
	require.NoError(t, err)
	assert.True(t, user.HasPaymentMethod)
}

func TestGetUser_DataAttributes(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/user", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(
			[]byte(
				`{"data":{"attributes":{"has_payment_method": false, "current_plan": "build"}}}`,
			),
		)
		require.NoError(t, err)
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	user, err := client.GetUser("test-token")
	require.NoError(t, err)
	assert.False(t, user.HasPaymentMethod)
}

func TestChangeApplicationPlan_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/applications/APP1/plan/self-serve",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPatch, r.Method)
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload ChangePlanRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			assert.Equal(t, "grow", payload.Plan)

			require.NoError(t, json.NewEncoder(w).Encode(SingleApplicationResponse{
				Data: ApplicationResource{
					ID: "APP1", Type: "application",
					Attributes: ApplicationAttributes{ApplicationID: "APP1", Name: "My App"},
				},
			}))
		},
	)

	ts, client := newTestClient(mux)
	defer ts.Close()

	app, err := client.ChangeApplicationPlan("test-token", "APP1", "grow")
	require.NoError(t, err)
	assert.Equal(t, "APP1", app.ID)
	assert.Equal(t, "My App", app.Name)
}

func TestChangeApplicationPlan_EmptyBodyFallback(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/applications/APP1/plan/self-serve",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	)

	ts, client := newTestClient(mux)
	defer ts.Close()

	// An empty/204 body still yields a usable result with the known app ID.
	app, err := client.ChangeApplicationPlan("test-token", "APP1", "free")
	require.NoError(t, err)
	assert.Equal(t, "APP1", app.ID)
}

func TestChangeApplicationPlan_Unauthorized(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/applications/APP1/plan/self-serve",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		},
	)

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.ChangeApplicationPlan("expired-token", "APP1", "grow")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "session expired")
}

func TestChangeApplicationPlan_ErrorStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/1/applications/APP1/plan/self-serve",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, err := w.Write([]byte(`{"errors":[{"title":"plan change not allowed"}]}`))
			require.NoError(t, err)
		},
	)

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.ChangeApplicationPlan("test-token", "APP1", "grow")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Couldn't change your application's plan: 422")
	// The response body must never be surfaced in the error, only the status.
	assert.NotContains(t, err.Error(), "plan change not allowed")
}

func TestGetCrawlerUser_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/crawler/user", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		require.NoError(t, json.NewEncoder(w).Encode(DashboardCrawlerUserResponse{
			Data: DashboardCrawlerUserData{
				ID:     "crawler-user-id",
				Email:  "crawler@example.com",
				Name:   "Crawler User",
				APIKey: "crawler-api-key",
			},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	user, err := client.GetCrawlerUser("test-token")
	require.NoError(t, err)
	assert.Equal(t, "crawler-user-id", user.ID)
	assert.Equal(t, "crawler@example.com", user.Email)
	assert.Equal(t, "Crawler User", user.Name)
	assert.Equal(t, "crawler-api-key", user.APIKey)
}

func TestGetCrawlerUser_HTTPError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/crawler/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		detail := "forbidden"
		require.NoError(t, json.NewEncoder(w).Encode(DashboardCrawlerErrorResponse{
			Errors: []DashboardCrawlerError{{
				Status: http.StatusText(http.StatusForbidden),
				Title:  "Forbidden",
				Detail: &detail,
			}},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.GetCrawlerUser("test-token")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get crawler user data: forbidden")
	assert.NotContains(t, err.Error(), "403")
}

func TestGetCrawlerUser_HTTPErrorWithoutDetail(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/crawler/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		require.NoError(t, json.NewEncoder(w).Encode(DashboardCrawlerErrorResponse{
			Errors: []DashboardCrawlerError{{
				Status: http.StatusText(http.StatusForbidden),
				Title:  "Forbidden",
				Detail: nil,
			}},
		}))
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.GetCrawlerUser("test-token")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get crawler user data: Forbidden")
	assert.NotContains(t, err.Error(), "403")
}

func TestGetCrawlerUser_InvalidJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/crawler/user", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`{"data":`))
		require.NoError(t, err)
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.GetCrawlerUser("test-token")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse crawler response")
}

func TestGetCrawlerUser_HTTPErrorInvalidJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/1/crawler/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte(`{"message":`))
		require.NoError(t, err)
	})

	ts, client := newTestClient(mux)
	defer ts.Close()

	_, err := client.GetCrawlerUser("test-token")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse crawler response")
}
