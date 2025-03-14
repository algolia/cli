package delete

import (
	"fmt"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/pkg/iostreams"
	"github.com/algolia/cli/test"
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
			cli:      "foo --rule-ids 1",
			tty:      false,
			wantsErr: true,
		},
		{
			name:     "--confirm without tty",
			cli:      "foo --rule-ids 1 --confirm",
			tty:      true,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm: false,
				Index:     "foo",
				RuleIDs: []string{
					"1",
				},
				ForwardToReplicas: false,
			},
		},
		{
			name:     "no --confirm with tty",
			cli:      "foo --rule-ids 1",
			tty:      true,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm: true,
				Index:     "foo",
				RuleIDs: []string{
					"1",
				},
				ForwardToReplicas: false,
			},
		},
		{
			name:     "multiple --rule-ids",
			cli:      "foo --rule-ids 1,2 --confirm",
			tty:      false,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm: false,
				Index:     "foo",
				RuleIDs: []string{
					"1",
					"2",
				},
				ForwardToReplicas: false,
			},
		},
		{
			name:     "multiple --rule-ids, forward to replicas",
			cli:      "foo --rule-ids 1,2 --confirm --forward-to-replicas",
			tty:      false,
			wantsErr: false,
			wantsOpts: DeleteOptions{
				DoConfirm: false,
				Index:     "foo",
				RuleIDs: []string{
					"1",
					"2",
				},
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

			assert.Equal(t, tt.wantsOpts.Index, opts.Index)
			assert.Equal(t, tt.wantsOpts.RuleIDs, opts.RuleIDs)
			assert.Equal(t, tt.wantsOpts.ForwardToReplicas, opts.ForwardToReplicas)
			assert.Equal(t, tt.wantsOpts.DoConfirm, opts.DoConfirm)
		})
	}
}

func Test_runDeleteCmd(t *testing.T) {
	tests := []struct {
		name    string
		cli     string
		index   string
		ruleIDs []string
		isTTY   bool
		wantOut string
	}{
		{
			name:  "single rule-id, no TTY",
			cli:   "foo --rule-ids 1 --confirm",
			index: "foo",
			ruleIDs: []string{
				"1",
			},
			isTTY:   false,
			wantOut: "",
		},
		{
			name:  "single rule-id, TTY",
			cli:   "foo --rule-ids 1 --confirm",
			index: "foo",
			ruleIDs: []string{
				"1",
			},
			isTTY:   true,
			wantOut: "✓ Successfully deleted 1 rule from foo\n",
		},
		{
			name:  "multiple rule-ids, TTY",
			cli:   "foo --rule-ids 1,2 --confirm",
			index: "foo",
			ruleIDs: []string{
				"1",
				"2",
			},
			isTTY:   true,
			wantOut: "✓ Successfully deleted 2 rules from foo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			for _, id := range tt.ruleIDs {
				r.Register(
					httpmock.REST("GET", fmt.Sprintf("1/indexes/%s/rules/%s", tt.index, id)),
					httpmock.JSONResponse(search.SearchRulesResponse{}),
				)
				r.Register(
					httpmock.REST("DELETE", fmt.Sprintf("1/indexes/%s/rules/%s", tt.index, id)),
					httpmock.JSONResponse(search.DeletedAtResponse{}),
				)
			}

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
