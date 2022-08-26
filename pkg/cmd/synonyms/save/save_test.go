package save

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

func TestNewSaveCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts SaveOptions
	}{
		{
			name:     "without tty",
			cli:      "legends --id 1 --synonyms jordan,mj",
			tty:      false,
			wantsErr: false,
			wantsOpts: SaveOptions{
				Indice:            "legends",
				SynonymID:         "1",
				Synonyms:          []string{"jordan", "mj"},
				ForwardToReplicas: false,
			},
		},
		{
			name:     "with tty",
			cli:      "legends --id 1 --synonyms jordan,mj",
			tty:      true,
			wantsErr: false,
			wantsOpts: SaveOptions{
				Indice:            "legends",
				SynonymID:         "1",
				Synonyms:          []string{"jordan", "mj"},
				ForwardToReplicas: false,
			},
		},
		{
			name:     "single, --one-way without --input",
			cli:      "legends --id 1 --synonyms jordan,mj --one-way",
			tty:      false,
			wantsErr: true,
		},
		{
			name:     "single, --one-way",
			cli:      "legends -i 1 -s jordan,mj --one-way --input goat",
			tty:      true,
			wantsErr: false,
			wantsOpts: SaveOptions{
				Indice:            "legends",
				SynonymID:         "1",
				Synonyms:          []string{"jordan", "mj"},
				OneWaySynonym:     true,
				SynonymInput:      "goat",
				ForwardToReplicas: false,
			},
		},
		{
			name:     "single, forward to replicas",
			cli:      "legends --id 1 --synonyms jordan,mj --forward-to-replicas",
			tty:      false,
			wantsErr: false,
			wantsOpts: SaveOptions{
				Indice:            "legends",
				Synonyms:          []string{"jordan", "mj"},
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

			var opts *SaveOptions
			cmd := NewSaveCmd(f, func(o *SaveOptions) error {
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
			assert.Equal(t, tt.wantsOpts.SynonymInput, opts.SynonymInput)
			assert.Equal(t, tt.wantsOpts.OneWaySynonym, opts.OneWaySynonym)
			assert.Equal(t, tt.wantsOpts.Synonyms, opts.Synonyms)
			assert.Equal(t, tt.wantsOpts.ForwardToReplicas, opts.ForwardToReplicas)
		})
	}
}

func Test_runSaveCmd(t *testing.T) {
	tests := []struct {
		name      string
		cli       string
		indice    string
		synonymID string
		isTTY     bool
		wantOut   string
	}{
		{
			name:      "single id, two synonyms, no TTY",
			cli:       "legends --id 1 --synonyms jorda,mj",
			indice:    "legends",
			synonymID: "1",
			isTTY:     false,
			wantOut:   "",
		},
		{
			name:      "single id, two synonyms, TTY",
			cli:       "legends --id 1 --synonyms jordan,mj",
			indice:    "legends",
			synonymID: "1",
			isTTY:     true,
			wantOut:   "✓ Synonym '1' successfully created with 2 synonyms (jordan, mj) to legends\n",
		},
		{
			name:      "single id, mutiple synonyms, TTY",
			cli:       "legends --id 1 --synonyms jordan,mj,goat,michael,23",
			indice:    "legends",
			synonymID: "1",
			isTTY:     true,
			wantOut:   "✓ Synonym '1' successfully created with 5 synonyms (jordan, mj, goat, michael, 23) to legends\n",
		},
		{
			name:      "single id, mutiple synonyms, TTY with shorthands",
			cli:       "legends -i 1 -s jordan,mj,goat,michael,23",
			indice:    "legends",
			synonymID: "1",
			isTTY:     true,
			wantOut:   "✓ Synonym '1' successfully created with 5 synonyms (jordan, mj, goat, michael, 23) to legends\n",
		},
		{
			name:      "single id, mutiple synonyms, one way with input",
			cli:       "legends -i 1 -s jordan,mj,goat,23 -o -n michael",
			indice:    "legends",
			synonymID: "1",
			isTTY:     true,
			wantOut:   "✓ One way synonym '1' successfully created with 4 synonyms (jordan, mj, goat, 23) to legends\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			r.Register(httpmock.REST("PUT", fmt.Sprintf("1/indexes/%s/synonyms/%s", tt.indice, tt.synonymID)), httpmock.JSONResponse(search.RegularSynonym{}))
			defer r.Verify(t)

			f, out := test.NewFactory(tt.isTTY, &r, nil, "")
			cmd := NewSaveCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
