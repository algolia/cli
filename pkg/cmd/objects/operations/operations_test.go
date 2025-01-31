package operations

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/algolia/cli/pkg/httpmock"
	"github.com/algolia/cli/test"
)

func Test_runOperationsCmd(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "operations.json")
	err := os.WriteFile(
		tmpFile,
		[]byte(
			`{"action":"addObject","indexName":"index1","body":{"firstname":"Jimmie","lastname":"Barninger"}}`,
		),
		0o600,
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		cli     string
		stdin   string
		wantOut string
		wantErr string
	}{
		{
			name:    "from stdin",
			cli:     "-F -",
			stdin:   `{"action":"addObject","indexName":"index1","body":{"firstname":"Jimmie","lastname":"Barninger"}}`,
			wantOut: "✓ Successfully processed 1 operations in",
		},
		{
			name:    "from file",
			cli:     fmt.Sprintf("-F '%s'", tmpFile),
			wantOut: "✓ Successfully processed 1 operations in",
		},
		{
			name:    "from stdin with invalid JSON",
			cli:     "-F -",
			stdin:   `{"action":"addObject","indexName":"index1","body":{"firstname":"Jimmie","lastname":"Barninger"}},`,
			wantErr: "X Found 1 error (out of 1 operations) while parsing the file:\n  line 1: invalid character ',' after top-level value\n",
		},

		{
			name: "from stdin with invalid JSON (multiple operations)",
			cli:  "-F -",
			stdin: `{"action": "addObject","indexName":"index1"},
			{"test": "bar"}`,
			wantErr: "X Found 2 errors (out of 2 operations) while parsing the file:\n  line 1: invalid character ',' after top-level value\n  missing action\n",
		},
		{
			name:    "from stdin with invalid JSON (1 operation) with --continue-on-error",
			cli:     "-F - --continue-on-error",
			stdin:   `{"action": "addObject"},`,
			wantErr: "X Found 1 error (out of 1 operations) while parsing the file:\n  line 1: invalid character ',' after top-level value\n",
		},
		{
			name: "from stdin with invalid JSON (2 objects) with --continue-on-error",
			cli:  "-F - --continue-on-error",
			stdin: `{"action": "addObject","indexName":"index1"}
			{"action": "deleteObject","indexName":"index2",body:{"objectID": "abc"}}`,
			wantOut: "✓ Successfully processed 1 operations in",
		},
		{
			name:    "missing file flag",
			cli:     "",
			wantErr: "required flag(s) \"file\" not set",
		},
		{
			name:    "non-existant file",
			cli:     "-F /tmp/foo",
			wantErr: "open /tmp/foo: no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpmock.Registry{}
			if tt.wantErr == "" {
				r.Register(
					httpmock.REST("POST", "1/indexes/*/batch"),
					httpmock.JSONResponse(search.MultipleBatchRes{}),
				)
			}
			defer r.Verify(t)

			f, out := test.NewFactory(true, &r, nil, tt.stdin)
			cmd := NewOperationsCmd(f, nil)
			out, err := test.Execute(cmd, tt.cli, out)
			if err != nil {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			assert.Contains(t, out.String(), tt.wantOut)
		})
	}
}

func Test_ValidateBatchOperation(t *testing.T) {
	tests := []struct {
		name       string
		action     string
		body       map[string]interface{}
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "no action",
			action:     "",
			body:       nil,
			wantErr:    true,
			wantErrMsg: "missing action",
		},
		{
			name:       "invalid action",
			action:     "invalid",
			body:       nil,
			wantErr:    true,
			wantErrMsg: "invalid action \"invalid\" (valid actions are addObject, updateObject, partialUpdateObject, partialUpdateObjectNoCreate and deleteObject)",
		},
		{
			name:       "missing objectID for deleteObject action",
			action:     string(search.DeleteObject),
			body:       nil,
			wantErr:    true,
			wantErrMsg: "missing objectID for action deleteObject",
		},
	}

	for _, act := range []string{
		string(search.AddObject), string(search.UpdateObject),
		string(search.PartialUpdateObject), string(search.PartialUpdateObjectNoCreate),
	} {
		tests = append(tests, struct {
			name       string
			action     string
			body       map[string]interface{}
			wantErr    bool
			wantErrMsg string
		}{
			name:    act,
			action:  act,
			body:    nil,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batchOperation := search.BatchOperation{
				Action: search.BatchAction(tt.action),
			}
			if tt.body != nil {
				batchOperation.Body = tt.body
			}

			err := ValidateBatchOperation(search.BatchOperationIndexed{
				IndexName:      "index1",
				BatchOperation: batchOperation,
			})
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
				return
			}
			assert.NoError(t, err)
		})
	}
}
