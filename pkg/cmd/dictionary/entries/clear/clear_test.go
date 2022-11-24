package clear

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

func TestNewClearCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts ClearOptions
	}{
		{
			name:     "no --confirm without tty",
			cli:      "plurals",
			tty:      false,
			wantsErr: true,
			wantsOpts: ClearOptions{
				DoConfirm: true,
				Dictionnaries: []search.DictionaryName{
					search.Plurals,
				},
			},
		},
		{
			name:     "--confirm without tty",
			cli:      "plurals --confirm",
			tty:      false,
			wantsErr: false,
			wantsOpts: ClearOptions{
				DoConfirm: false,
				Dictionnaries: []search.DictionaryName{
					search.Plurals,
				},
			},
		},
		{
			name:     "with only --all",
			cli:      "--all",
			tty:      true,
			wantsErr: false,
			wantsOpts: ClearOptions{
				DoConfirm: true,
				Dictionnaries: []search.DictionaryName{
					search.Stopwords,
					search.Plurals,
					search.Compounds,
				},
			},
		},
		{
			name:     "with args and --all",
			cli:      "plurals --all --confirm",
			tty:      false,
			wantsErr: true,
			wantsOpts: ClearOptions{
				DoConfirm: false,
				Dictionnaries: []search.DictionaryName{
					search.Plurals,
				},
			},
		},
		{
			name:     "wrong dictionary name",
			cli:      "foo --confirm",
			tty:      false,
			wantsErr: true,
			wantsOpts: ClearOptions{
				DoConfirm: false,
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

			var opts *ClearOptions
			cmd := NewClearCmd(f, func(o *ClearOptions) error {
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
			assert.Equal(t, tt.wantsOpts.DoConfirm, opts.DoConfirm)
		})
	}
}

func Test_runDeleteCmd(t *testing.T) {
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
			cli:  "plurals --confirm",
			dictionaries: []search.DictionaryName{
				search.Plurals,
			},
			entries: true,
			isTTY:   false,
			wantOut: "",
		},
		{
			name: "multiple dictionaries",
			cli:  "plurals compounds --confirm",
			dictionaries: []search.DictionaryName{
				search.Plurals,
				search.Compounds,
			},
			entries: true,
			isTTY:   false,
			wantOut: "",
		},
		{
			name: "all dictionaries",
			cli:  "--all --confirm",
			dictionaries: []search.DictionaryName{
				search.Stopwords,
				search.Plurals,
				search.Compounds,
			},
			entries: true,
			isTTY:   false,
			wantOut: "",
		},
		{
			name: "no entries",
			cli:  "plurals --confirm",
			dictionaries: []search.DictionaryName{
				search.Plurals,
			},
			entries: false,
			isTTY:   false,
			wantOut: "! No entries to clear in plurals dictionary.\n",
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
			cmd := NewClearCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
