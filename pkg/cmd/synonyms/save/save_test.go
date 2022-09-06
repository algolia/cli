package save

import (
	"fmt"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
				Indice: "legends",
				Synonym: search.NewRegularSynonym(
					"1",
					"jordan", "mj",
				),
				ForwardToReplicas: false,
			},
		},
		{
			name:     "with tty",
			cli:      "legends --id 1 --synonyms jordan,mj",
			tty:      true,
			wantsErr: false,
			wantsOpts: SaveOptions{
				Indice: "legends",
				Synonym: search.NewRegularSynonym(
					"1",
					"jordan", "mj",
				),
				ForwardToReplicas: false,
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
			assert.Equal(t, tt.wantsOpts.Synonym, opts.Synonym)
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
			wantOut:   "✓ Synonym '1' successfully saved with 2 synonyms (jordan, mj) to legends\n",
		},
		{
			name:      "single id, mutiple synonyms, TTY",
			cli:       "legends --id 1 --synonyms jordan,mj,goat,michael,23",
			indice:    "legends",
			synonymID: "1",
			isTTY:     true,
			wantOut:   "✓ Synonym '1' successfully saved with 5 synonyms (jordan, mj, goat, michael, 23) to legends\n",
		},
		{
			name:      "single id, mutiple synonyms, TTY with shorthands",
			cli:       "legends -i 1 -s jordan,mj,goat,michael,23",
			indice:    "legends",
			synonymID: "1",
			isTTY:     true,
			wantOut:   "✓ Synonym '1' successfully saved with 5 synonyms (jordan, mj, goat, michael, 23) to legends\n",
		},
		{
			name:      "single id, oneWaySynonym type, multiple synonyms, TTY",
			cli:       "legends --id 1 --type oneWaySynonym --synonyms jordan,mj,goat,michael --input 23",
			indice:    "legends",
			synonymID: "1",
			isTTY:     true,
			wantOut:   "✓ One way synonym '1' successfully saved with input '23' and 4 synonyms (jordan, mj, goat, michael) to legends\n",
		},
		{
			name:      "single id, placeholder type, one placeholder, multiple replacements, TTY",
			cli:       "legends -i 1 -t placeholder -l jordan -r mj,goat,michael,23",
			indice:    "legends",
			synonymID: "1",
			isTTY:     true,
			wantOut:   "✓ Placeholder synonym '1' successfully saved with placeholder 'jordan' and 4 replacements (mj, goat, michael, 23) to legends\n",
		},
		{
			name:      "single id, altCorrection1 type, one word, multiple corrections, TTY",
			cli:       "legends -i 1 -t altCorrection1 -w jordan -c mj,goat,michael,23",
			indice:    "legends",
			synonymID: "1",
			isTTY:     true,
			wantOut:   "✓ Alt correction 1 synonym '1' successfully saved with word 'jordan' and 4 corrections (mj, goat, michael, 23) to legends\n",
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
