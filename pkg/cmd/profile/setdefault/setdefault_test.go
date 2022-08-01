package setdefault

import (
	"testing"

	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/test"
	"github.com/bmizerany/assert"
)

func Test_runSetDefaultCmd(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		profiles map[string]bool
		wantsErr string
		wantOut  string
	}{
		{
			name:     "existing default",
			cli:      "foo",
			profiles: map[string]bool{"default": true, "foo": false},
			wantOut:  "✓ Default profile successfuly changed from 'default' to 'foo'.\n",
		},
		{
			name:     "non-existing default",
			cli:      "foo",
			profiles: map[string]bool{"foo": false},
			wantOut:  "✓ Default profile successfuly set to 'foo'.\n",
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
			f, out := test.NewFactory(true, nil, cfg, "")
			cmd := NewSetDefaultCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				assert.Equal(t, tt.wantsErr, err.Error())
				return
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
