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

// envCondition exposes process-env-var checks to testscripts via the
// `[env:VAR]` and `[!env:VAR]` syntax. Truthy means the var is set to
// something other than "" or "0" (matches how callers like
// ALGOLIA_AGENT_STUDIO_E2E=1 set it). Lets a single .txtar file gate
// itself with `[!env:VAR] skip 'reason'` instead of forking the test
// runner per-script.
func envCondition(cond string) (bool, error) {
	const prefix = "env:"
	if !strings.HasPrefix(cond, prefix) {
		return false, fmt.Errorf("unknown condition %q (only env:VAR is registered)", cond)
	}
	v := os.Getenv(strings.TrimPrefix(cond, prefix))
	return v != "" && v != "0", nil
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
		Dir:       dirName,
		Setup:     setupEnv(testEnv),
		Cmds:      setupCmds(testEnv),
		Condition: envCondition,
	})
}

// TestVersion tests the version option
func TestVersion(t *testing.T) {
	runTestsInDir(t, "testscripts/version")
}

// TestIndices test `algolia indices` commands
func TestIndices(t *testing.T) {
	runTestsInDir(t, "testscripts/indices")
}

// TestSettings tests `algolia settings` commands
func TestSettings(t *testing.T) {
	runTestsInDir(t, "testscripts/settings")
}

// TestObjects tests `algolia objects` commands
func TestObjects(t *testing.T) {
	runTestsInDir(t, "testscripts/objects")
}

// TestSynonyms tests `algolia synonyms` commands
func TestSynonyms(t *testing.T) {
	runTestsInDir(t, "testscripts/synonyms")
}

// TestRules tests `algolia rules` commands
func TestRules(t *testing.T) {
	runTestsInDir(t, "testscripts/rules")
}

// TestAgents tests `algolia agents` commands.
//
// Two-tier coverage:
//   - Ungated scripts (`*.txtar` except list): CLI validation helpers
//     that exercise `algolia agents` cobra validators without network calls.
//   - list.txtar: gated on ALGOLIA_AGENT_STUDIO_E2E=1. Read-only
//     smoke against the live backend. Confirms the wire format we
//     parse against still matches what the deployed service emits.
//
// Write-CRUD live coverage (create → update → delete same id) is
// not in here yet — testscript's framework needs an id-extraction
// helper we haven't added.
func TestAgents(t *testing.T) {
	runTestsInDir(t, "testscripts/agents")
}

// TestSearch tests `algolia search`
func TestSearch(t *testing.T) {
	runTestsInDir(t, "testscripts/search")
}

// TestAgentReady tests describe and dry-run contracts.
func TestAgentReady(t *testing.T) {
	runTestsInDir(t, "testscripts/agent-ready")
}
