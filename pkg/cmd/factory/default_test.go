package factory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_resolveAgentStudioBaseURL(t *testing.T) {
	tests := []struct {
		name            string
		profileOverride string
		buildDefault    string
		appID           string
		want            string
		wantErr         bool
	}{
		{
			name:            "profile override wins over build default",
			profileOverride: "https://debug.example.com",
			buildDefault:    "https://agent-studio.staging.eu.algolia.com",
			appID:           "betaXYZ",
			want:            "https://debug.example.com",
		},
		{
			name:            "build default wins over cluster-proxy fallback when profile is empty",
			profileOverride: "",
			buildDefault:    "https://agent-studio.staging.eu.algolia.com",
			appID:           "betaXYZ",
			want:            "https://agent-studio.staging.eu.algolia.com",
		},
		{
			name:            "cluster-proxy fallback when both overrides are empty",
			profileOverride: "",
			buildDefault:    "",
			appID:           "APP123",
			want:            "https://APP123.algolia.net/agent-studio",
		},
		{
			name:            "trailing slash on profile override is trimmed",
			profileOverride: "https://debug.example.com/",
			buildDefault:    "",
			appID:           "APP123",
			want:            "https://debug.example.com",
		},
		{
			name:            "missing appID with no overrides errors out",
			profileOverride: "",
			buildDefault:    "",
			appID:           "",
			wantErr:         true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveAgentStudioBaseURL(tc.profileOverride, tc.buildDefault, tc.appID)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
