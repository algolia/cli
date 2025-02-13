package delete

import (
	"fmt"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/httpmock/v4"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test/v4"
)

func TestNewDeleteCmd(t *testing.T) {
	tests := []struct {
		name      string
		tty       bool
		cli       string
		wantsErr  bool
		wantsOpts DeleteOptions
	}{
		{
			name:     "no --confirm without tty",
			cli:      "plurals --object-ids 1",
			tty:      false,
			wantsErr: true,
		},
		{
			name:     "--confirm without tty",
			cli:      "plurals --object-ids 1 --confirm",
			tty:      true,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm:      false,
				DictionaryType: search.DICTIONARY_TYPE_PLURALS,
				ObjectIDs: []string{
					"1",
				},
			},
		},
		{
			name:     "no --confirm with tty",
			cli:      "plurals --object-ids 1",
			tty:      true,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm:      true,
				DictionaryType: search.DICTIONARY_TYPE_PLURALS,
				ObjectIDs: []string{
					"1",
				},
			},
		},
		{
			name:     "multiple --object-ids",
			cli:      "plurals --object-ids 1,2 --confirm",
			tty:      false,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm:      false,
				DictionaryType: search.DICTIONARY_TYPE_PLURALS,
				ObjectIDs: []string{
					"1",
					"2",
				},
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

			var opts *DeleteOptions
			cmd := NewDeleteCmd(f, func(o *DeleteOptions) error {
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

			assert.Equal(t, tt.wantsOpts.DictionaryType, opts.DictionaryType)
			assert.Equal(t, tt.wantsOpts.ObjectIDs, opts.ObjectIDs)
			assert.Equal(t, tt.wantsOpts.DoConfirm, opts.DoConfirm)
		})
	}
}

func Test_runDeleteCmd(t *testing.T) {
	tests := []struct {
		name       string
		cli        string
		dictionary search.DictionaryType
		objectIDs  []string
		isTTY      bool
		wantOut    string
	}{
		{
			name:       "single object-id, no TTY",
			cli:        "plurals --object-ids 1 --confirm",
			dictionary: search.DICTIONARY_TYPE_PLURALS,
			objectIDs: []string{
				"1",
			},
			isTTY:   false,
			wantOut: "",
		},
		{
			name:       "single object-id, TTY",
			cli:        "plurals --object-ids 1 --confirm",
			dictionary: search.DICTIONARY_TYPE_PLURALS,
			objectIDs: []string{
				"1",
			},
			isTTY:   true,
			wantOut: "✓ Successfully deleted 1 entry from plurals\n",
		},
		{
			name:       "multiple object-ids, TTY",
			cli:        "plurals --object-ids 1,2 --confirm",
			dictionary: search.DICTIONARY_TYPE_PLURALS,
			objectIDs: []string{
				"1",
				"2",
			},
			isTTY:   true,
			wantOut: "✓ Successfully deleted 2 entries from plurals\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}

			// test is flaky since there's no guarantee of obtaining the right object using a search by objectID
			for _, id := range tt.objectIDs {
				r.Register(
					httpmock.REST(
						"GET",
						fmt.Sprintf("1/dictionaries/%s/search?query=%s", tt.dictionary, id),
					),
					httpmock.JSONResponse(search.SearchDictionaryEntriesResponse{}),
				)
			}
			r.Register(
				httpmock.REST("POST", fmt.Sprintf("1/dictionaries/%s/batch", tt.dictionary)),
				httpmock.JSONResponse(search.UpdatedAtResponse{}),
			)

			f, out := test.NewFactory(tt.isTTY, &r, nil, "")
			cmd := NewDeleteCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOut, out.String())
		})
	}
}
