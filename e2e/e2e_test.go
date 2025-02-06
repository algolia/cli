//go:build e2e

package e2e_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/algolia/cli/pkg/cmd/root"
	"github.com/cli/go-internal/testscript"
)

// algolia runs the root command of the Algolia CLI
func algolia() int {
	return int(root.Execute())
}

// TestMain sets the executable program so that we don't depend on the compiled binary
func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"algolia": algolia,
	}))
}

// testEnvironment stores the environment variables we need to setup for the tests
type testEnvironment struct {
	AppID  string
	ApiKey string
}

// getEnv reads the environment variables and prints errors for missing ones
func (e *testEnvironment) getEnv() error {
	env := map[string]string{}

	required := []string{
		// The CLI testing Algolia app
		"ALGOLIA_APPLICATION_ID",
		// API key with sufficient permissions to run all tests
		"ALGOLIA_API_KEY",
	}

	var missing []string

	for _, envVar := range required {
		val, ok := os.LookupEnv(envVar)
		if val == "" || !ok {
			missing = append(missing, envVar)
			continue
		}

		env[envVar] = val
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}

	e.AppID = env["ALGOLIA_APPLICATION_ID"]
	e.ApiKey = env["ALGOLIA_API_KEY"]

	return nil
}

// For the `defer` function
var keyT struct{}

// setupEnv sets up the environment variables for the test
func setupEnv(testEnv testEnvironment) func(ts *testscript.Env) error {
	return func(ts *testscript.Env) error {
		ts.Setenv("ALGOLIA_APPLICATION_ID", testEnv.AppID)
		ts.Setenv("ALGOLIA_API_KEY", testEnv.ApiKey)

		ts.Values[keyT] = ts.T()
		return nil
	}
}

// setupCmds sets up custom commands we want to make available in the test scripts
func setupCmds(
	testEnv testEnvironment,
) map[string]func(ts *testscript.TestScript, neg bool, args []string) {
	return map[string]func(ts *testscript.TestScript, neg bool, args []string){
		"defer": func(ts *testscript.TestScript, neg bool, args []string) {
			if neg {
				ts.Fatalf("unsupported ! defer")
			}
			tt, ok := ts.Value(keyT).(testscript.T)
			if !ok {
				ts.Fatalf("%v is not a testscript.T", ts.Value(keyT))
			}
			ts.Defer(func() {
				if err := ts.Exec(args[0], args[1:]...); err != nil {
					tt.FailNow()
				}
			})
		},
	}
}

// runTestsInDir runs all test scripts from a directory
func runTestsInDir(t *testing.T, dirName string) {
	var testEnv testEnvironment
	if err := testEnv.getEnv(); err != nil {
		t.Fatal(err)
	}
	t.Parallel()
	t.Log("Running e2e tests in", dirName)
	testscript.Run(t, testscript.Params{
		Dir:   dirName,
		Setup: setupEnv(testEnv),
		Cmds:  setupCmds(testEnv),
	})
}

// TestVersion tests the version option
func TestVersion(t *testing.T) {
	runTestsInDir(t, "testscripts/version")
}
