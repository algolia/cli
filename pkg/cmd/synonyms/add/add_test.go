package add

import (
	"fmt"
	"testing"

	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
)

func TestNewAddCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts AddOptions
	}{
		{
			name:     "without tty",
			cli:      "legends --id 1 --values jordan,mj",
			tty:      false,
			wantsErr: false,
			wantsOpts: AddOptions{
				Indice:            "legends",
				SynonymID:         "1",
				SynonymValues:     []string{"jordan", "mj"},
				ForwardToReplicas: false,
			},
		},
		{
			name:     "with tty",
			cli:      "legends --id 1 --values jordan,mj",
			tty:      true,
			wantsErr: false,
			wantsOpts: AddOptions{
				Indice:            "legends",
				SynonymID:         "1",
				SynonymValues:     []string{"jordan", "mj"},
				ForwardToReplicas: false,
			},
		},
		{
			name:     "single, forward to replicas",
			cli:      "legends --id 1 --values jordan,mj --forward-to-replicas",
			tty:      false,
			wantsErr: false,
			wantsOpts: AddOptions{
				Indice:            "legends",
				SynonymValues:     []string{"jordan", "mj"},
				SynonymID:         "1",
				ForwardToReplicas: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, stdout, stderr := iostreams.Test()
			if tt.tty {
				io.SetStdinTTY(tt.tty)
				io.SetStdoutTTY(tt.tty)
			}

			f := &cmdutil.Factory{
				IOStreams: io,
			}

			var opts *AddOptions
			cmd := NewAddCmd(f, func(o *AddOptions) error {
				opts = o
				return nil
			})

			args, err := shlex.Split(tt.cli)
			require.NoError(t, err)
			cmd.SetArgs(args)
			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, "", stdout.String())
			assert.Equal(t, "", stderr.String())

			assert.Equal(t, tt.wantsOpts.Indice, opts.Indice)
			assert.Equal(t, tt.wantsOpts.SynonymID, opts.SynonymID)
			assert.Equal(t, tt.wantsOpts.SynonymValues, opts.SynonymValues)
			assert.Equal(t, tt.wantsOpts.ForwardToReplicas, opts.ForwardToReplicas)
		})
	}
}

func Test_runAddCmd(t *testing.T) {
	tests := []struct {
		name          string
		cli           string
		indice        string
		synonymID     string
		synonymValues []string
		isTTY         bool
		wantOut       string
	}{
		{
			name:      "single id, two values, no TTY",
			cli:       "legends --id 1 --values jorda,mj",
			indice:    "legends",
			synonymID: "1",
			isTTY:     false,
			wantOut:   "",
		},
		{
			name:      "single id, two values, TTY",
			cli:       "legends --id 1 --values jordan,mj",
			indice:    "legends",
			synonymID: "1",
			isTTY:     true,
			wantOut:   "✓ Synonym '1' successfully created with 2 values (jordan, mj) from legends\n",
		},
		{
			name:      "single id, mutiple values, TTY",
			cli:       "legends --id 1 --values jordan,mj,goat,michael,23",
			indice:    "legends",
			synonymID: "1",
			isTTY:     true,
			wantOut:   "✓ Synonym '1' successfully created with 5 values (jordan, mj, goat, michael, 23) from legends\n",
		},
		{
			name:      "single id, mutiple values, TTY with shorthands",
			cli:       "legends -i 1 -v jordan,mj,goat,michael,23",
			indice:    "legends",
			synonymID: "1",
			isTTY:     true,
			wantOut:   "✓ Synonym '1' successfully created with 5 values (jordan, mj, goat, michael, 23) from legends\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(httpmock.REST("PUT", fmt.Sprintf("1/indexes/%s/synonyms/%s", tt.indice, tt.synonymID)), httpmock.JSONResponse(search.RegularSynonym{}))
			defer r.Verify(t)

			f, out := test.NewFactory(tt.isTTY, &r, nil, "")
			cmd := NewAddCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
