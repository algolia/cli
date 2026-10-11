# Contributing to the Algolia CLI

Welcome to the contributing guide for the Algolia CLI!

If this guide does not contain what you are looking for and thus prevents you from contributing, don't hesitate to [open an issue](https://github.com/algolia/cli/issues).

## Reporting an issue

Opening an issue is very effective way to contribute because many users might also be impacted. We'll make sure to fix it quickly if it's technically feasible and doesn't have important side effects for other users.

Before reporting an issue, first check that there is not an already open issue for the same topic using the [issues page](https://github.com/algolia/cli/issues). Don't hesitate to thumb up an issue that corresponds to the problem you have.

To help us solve your issue faster, please include the version of the CLI you are using (`algolia --version`), your OS, and the full command and output.

## Code contribution

For any code contribution, you need to:

- Fork and clone the project
- Create a new branch for what you want to solve (`fix/issue-number`, `feat/name-of-the-feature`)
- Make your changes
- Open a pull request

Then:

- Automatic checks will be run (tests, linters, build)
- A team member will review the pull request

When every check is green and a team member approves, your contribution is merged! 🚀

In your pull request description, explain what the change does and list the CLI commands a reviewer can run to verify it manually.

For coding conventions (architecture, style, testing patterns), see [AGENTS.md](AGENTS.md).

## Commit conventions

This project follows the [conventional changelog](https://conventionalcommits.org/) approach. This means that all commit messages should be formatted using the following scheme:

```
type(scope): description
```

In most cases, we use the following types:

- `fix`: for any resolution of an issue (identified or not)
- `feat`: for any new feature
- `refactor`: for any code change that neither adds a feature nor fixes an issue
- `docs`: for any documentation change or addition
- `chore`: for anything that is not related to the CLI itself (doc, tooling)

Pull requests are squash-merged, so the pull request title must follow the same convention.

Some examples of valid commit messages (used as first lines):

> - feat(config): new credentials model and profile deprecation
> - fix: setup pre-release channel releases
> - chore: update search api spec
> - docs(contributing): reword release section

## Requirements

To run this project, you will need:

- [Go](https://go.dev/) 1.23.4 and [golangci-lint](https://golangci-lint.run/) 1.63.4, as listed in [`.tool-versions`](.tool-versions) (works with `asdf` or `mise`)
- [Task](https://taskfile.dev/), the task runner used for all common commands
- [gofumpt](https://github.com/mvdan/gofumpt) and [golines](https://github.com/segmentio/golines) for formatting

## Launch the dev environment

```sh
git clone https://github.com/algolia/cli.git
cd cli
task build
```

`task build` runs `go generate ./...` first, then produces the `algolia` binary at the root of the repository.

You can then use it by doing `./algolia your_command ...`

Build-time defaults (dashboard URL, OAuth client ID, and so on) are injected from environment variables. If you need non-default values, copy [`.env.example`](.env.example) to `.env` and fill it in.

## Tests

To run all the unit tests:

```sh
task test
```

To run a single package or a single test:

```sh
go test ./pkg/cmd/apikeys/list
go test ./pkg/cmd/apikeys/list -run Test_runListCmd
```

End-to-end tests make real requests to the Algolia API. They need `ALGOLIA_APPLICATION_ID` and `ALGOLIA_API_KEY` in your environment or in the root `.env` file:

```sh
task e2e
```

Keep e2e runs narrow while iterating, for example `go test ./e2e -tags=e2e -run TestIndices`.

If you change flag surfaces or code derived from the API specs, run `go generate ./...` and commit the result.

## Linting and formatting

```sh
task lint
task format
```

Only format the code you touched, and avoid unrelated formatting churn.

## Release

Releases are handled by maintainers and are fully automated from a git tag:

1. Run the [Create release tag](.github/workflows/create-release-tag.yml) workflow from the GitHub Actions tab. Pick a release type (`fix`, `minor`, `major`, `prerelease`) and it computes the next version and pushes the tag. Use the `dry_run` option to preview the computed tag first.
2. The pushed tag triggers the [goreleaser workflow](.github/workflows/releases.yml), which builds the binaries and publishes the GitHub release, the Chocolatey package, the npm packages, and a docs update for stable releases.

There is nothing to run locally and no version file to bump.
