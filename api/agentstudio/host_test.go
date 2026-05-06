package agentstudio

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveHost(t *testing.T) {
	tests := []struct {
		name    string
		opts    HostOptions
		want    string
		wantErr error
	}{
		{
			name: "override wins over everything",
			opts: HostOptions{
				Override:      "https://custom.example/",
				Region:        RegionEU,
				ApplicationID: "APP123",
			},
			want: "https://custom.example",
		},
		{
			name: "override preserved verbatim (no scheme rewriting)",
			opts: HostOptions{Override: "http://localhost:8000"},
			want: "http://localhost:8000",
		},
		{
			name: "eu prod",
			opts: HostOptions{Region: RegionEU},
			want: "https://agent-studio.eu.algolia.com",
		},
		{
			name: "us prod",
			opts: HostOptions{Region: RegionUS, Env: EnvProd},
			want: "https://agent-studio.us.algolia.com",
		},
		{
			name: "eu staging",
			opts: HostOptions{Region: RegionEU, Env: EnvStaging},
			want: "https://agent-studio.staging.eu.algolia.com",
		},
		{
			name:    "us staging is not supported",
			opts:    HostOptions{Region: RegionUS, Env: EnvStaging},
			wantErr: ErrStagingNotInRegion,
		},
		{
			name: "region case-insensitive",
			opts: HostOptions{Region: "EU"},
			want: "https://agent-studio.eu.algolia.com",
		},
		{
			name: "env case-insensitive",
			opts: HostOptions{Region: RegionEU, Env: "STAGING"},
			want: "https://agent-studio.staging.eu.algolia.com",
		},
		{
			name:    "unknown region rejected",
			opts:    HostOptions{Region: "apac"},
			wantErr: ErrUnknownRegion,
		},
		{
			name: "cluster-proxy fallback when region missing but appID present",
			opts: HostOptions{ApplicationID: "APP123"},
			want: "https://APP123.algolia.net/agent-studio",
		},
		{
			name:    "no inputs at all returns ErrNoHostResolvable",
			opts:    HostOptions{},
			wantErr: ErrNoHostResolvable,
		},
		{
			name:    "unknown env rejected",
			opts:    HostOptions{Region: RegionEU, Env: "preview"},
			wantErr: nil, // sentinel-less; check via Contains below
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ResolveHost(tc.opts)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tc.wantErr), "got %v, want errors.Is(%v)", err, tc.wantErr)
				return
			}
			if tc.name == "unknown env rejected" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "unknown agent studio env")
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
