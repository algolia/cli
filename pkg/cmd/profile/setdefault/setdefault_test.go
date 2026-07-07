package setdefault

import (
	"strings"
	"testing"

	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/test"
	"github.com/bmizerany/assert"
)

func Test_runSetDefaultCmd(t *testing.T) {
	tests := []struct {
		name          string
		cli           string
		profiles      map[string]bool
		hasStateFile  bool
		wantsErr      string
		wantOut       string
		wantErrOut    string
		notWantErrOut string
	}{
		{
			name:          "existing default",
			cli:           "foo",
			profiles:      map[string]bool{"default": true, "foo": false},
			wantOut:       "✓ Default profile successfuly changed from 'default' to 'foo'.\n",
			notWantErrOut: "state.toml",
		},
		{
			name:     "non-existing default",
			cli:      "foo",
			profiles: map[string]bool{"foo": false},
			wantOut:  "✓ Default profile successfuly set to 'foo'.\n",
		},
		{
			name:         "state file exists",
			cli:          "foo",
			profiles:     map[string]bool{"default": true, "foo": false},
			hasStateFile: true,
			wantOut:      "✓ Default profile successfuly changed from 'default' to 'foo'.\n",
			wantErrOut:   "changes to config.toml profiles will be ignored in a future version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p []*config.Profile
			for k, v := range tt.profiles {
				p = append(p, &config.Profile{
					Name:    k,
					Default: v,
				})
			}
			cfg := test.NewConfigStubWithProfiles(p)
			cfg.HasStateFile = tt.hasStateFile
			f, out := test.NewFactory(true, nil, cfg, "")
			cmd := NewSetDefaultCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				assert.Equal(t, tt.wantsErr, err.Error())
				return
			}

			assert.Equal(t, tt.wantOut, out.String())
			if tt.wantErrOut != "" {
				assert.Equal(t, true, strings.Contains(out.Stderr(), tt.wantErrOut))
			}
			if tt.notWantErrOut != "" {
				assert.Equal(t, false, strings.Contains(out.Stderr(), tt.notWantErrOut))
			}
		})
	}
}
