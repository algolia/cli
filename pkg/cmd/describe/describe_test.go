package describe_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	describeCmd "github.com/algolia/cli/pkg/cmd/describe"
	"github.com/algolia/cli/pkg/cmd/root"
	"github.com/algolia/cli/test"
)

func TestDescribeCommand(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := root.NewRootCmd(f)

	out, err := test.Execute(cmd, "describe search", out)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, out.String(), `"schemaVersion":"v1"`)
	assert.Contains(t, out.String(), `"name":"algolia search"`)
	assert.Contains(t, out.String(), `"commandType":"read"`)
	assert.Contains(t, out.String(), `"query"`)
}

func TestSchemaAlias(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := root.NewRootCmd(f)

	out, err := test.Execute(cmd, "schema", out)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, out.String(), `"schemaVersion":"v1"`)
	assert.Contains(t, out.String(), `"name":"algolia"`)
}

func TestStandaloneDescribeCommandDefaultsToJSON(t *testing.T) {
	f, out := test.NewFactory(false, nil, nil, "")
	cmd := describeCmd.NewDescribeCmd(f)

	out, err := test.Execute(cmd, "", out)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, out.String(), `"schemaVersion":"v1"`)
}
