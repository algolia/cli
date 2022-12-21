package browse

import (
	"fmt"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/stretchr/testify/assert"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

/* func TestNewBrowseCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts BrowseOptions
	}{
		{
			name:     "with only --all",
			cli:      "--all",
			tty:      true,
			wantsErr: false,
			wantsOpts: BrowseOptions{
				Dictionnaries: []search.DictionaryName{
					search.Stopwords,
					search.Plurals,
					search.Compounds,
				},
			},
		},
		{
			name:     "with args and --all",
			cli:      "plurals --all",
			tty:      false,
			wantsErr: true,
			wantsOpts: BrowseOptions{
				Dictionnaries: []search.DictionaryName{
					search.Plurals,
				},
			},
		},
		{
			name:      "wrong dictionary name",
			cli:       "foo",
			tty:       false,
			wantsErr:  true,
			wantsOpts: BrowseOptions{},
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

			var opts *BrowseOptions
			cmd := NewBrowseCmd(f, func(o *BrowseOptions) error {
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

			assert.Equal(t, tt.wantsOpts.Dictionnaries, opts.Dictionnaries)
		})
	}
} */

func Test_runBrowseCmd(t *testing.T) {
	tests := []struct {
		name         string
		cli          string
		dictionaries []search.DictionaryName
		entries      bool
		isTTY        bool
		wantOut      string
	}{
		{
			name: "one dictionary",
			cli:  "plurals",
			dictionaries: []search.DictionaryName{
				search.Plurals,
			},
			entries: true,
			isTTY:   false,
			wantOut: "\"plurals\"\n[{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}]\n",
		},
		{
			name: "multiple dictionaries",
			cli:  "plurals compounds",
			dictionaries: []search.DictionaryName{
				search.Plurals,
				search.Compounds,
			},
			entries: true,
			isTTY:   false,
			wantOut: "\"plurals\"\n[{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}]\n\"compounds\"\n[{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}]\n",
		},
		{
			name: "all dictionaries",
			cli:  "--all",
			dictionaries: []search.DictionaryName{
				search.Stopwords,
				search.Plurals,
				search.Compounds,
			},
			entries: true,
			isTTY:   false,
			wantOut: "\"stopwords\"\n[{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}]\n\"plurals\"\n[{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}]\n\"compounds\"\n[{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}]\n",
		},
		{
			name: "one dictionnary with default stopwords",
			cli:  "--all --showDefaultStopwords",
			dictionaries: []search.DictionaryName{
				search.Stopwords,
				search.Plurals,
				search.Compounds,
			},
			entries: true,
			isTTY:   false,
			wantOut: "\"stopwords\"\n[{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}]\n\"plurals\"\n[{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}]\n\"compounds\"\n[{\"Type\":\"custom\",\"ObjectID\":\"\",\"Language\":\"\"}]\n",
		},
		{
			name: "no entries",
			cli:  "plurals",
			dictionaries: []search.DictionaryName{
				search.Plurals,
			},
			entries: false,
			isTTY:   false,
			wantOut: "! No custom entries in plurals dictionary.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			for _, d := range tt.dictionaries {
				var entries []DictionaryEntry
				if tt.entries {
					entries = append(entries, DictionaryEntry{Type: "custom"})
				}
				r.Register(httpmock.REST("POST", fmt.Sprintf("1/dictionaries/%s/search", d)), httpmock.JSONResponse(search.SearchDictionariesRes{
					Hits: entries,
				}))
				r.Register(httpmock.REST("POST", fmt.Sprintf("1/dictionaries/%s/batch", d)), httpmock.JSONResponse(search.TaskStatusRes{}))
			}

			f, out := test.NewFactory(tt.isTTY, &r, nil, "")
			cmd := NewBrowseCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
