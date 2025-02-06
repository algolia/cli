# End-to-end tests

These tests run CLI commands like a user would,
built on top of the [`go-internal/testscript`](https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript) package.

They make real API requests,
so they work best in an empty Algolia application.
To run these tests,
you need to set the `ALGOLIA_APPLICATION_ID` and `ALGOLIA_API_KEY` environment variables.
If you're using `devbox`, create a `testing.env` file with these variables.
If you start a development environment with `devbox shell`,
the environment variables will be available for you.

## New tests

The tests use a very simple format.
For more information, run `go doc testscript`.

To add a new scenario, create a new directory under the `testscripts` directory,
and add your files with the extension `txtar`.
Each test directory can have multiple test files.
Multiple directories are tested in parallel.

### Example

A simple 'hello world' testscript may look like this:

```txt
# Test if output is hello
exec echo 'hello'
! stderr .
stdout '^hello\n$'
```

Read the documentation of the `testscript` package for more information.

To add the new directory to the test suite,
add a new function to the file `./e2e/e2e_test.go`.
The function name must begin with `Test`.

```go
// TestHello is a basic example
func TestHello(t *testing.T) {
	RunTestsInDir(t, "testscripts/hello")
}
```

## Notes

Since this makes real real requests to the same Algolia application,
these tests aren't fully isolated from each other.

To make tests interfere less, follow these guidelines:

- Use a unique index name in each `txtar` file.
  For example, use `test-index` in `indices.txtar` and `test-settings` in `settings.txtar`

- Delete indices at the end of your test with `defer`.
  For an example, see `indices.txtar`.

- Don't test for number of indices, or empty lists.
  Different tests might also create their own indices,
  and thus will fail a test that expect a certain number of indices present.
  You can ensure that the index with a given name exists or doesn't exist,
  by searching for the index name's pattern in the standard output.
  Again, see `indices.txtar`.
