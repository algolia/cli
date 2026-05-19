# Algolia CLI

The official Algolia CLI lets you manage your Algolia resources — indices, records, API keys, and synonyms — directly from the command line.

> This package installs a prebuilt Go binary. The npm wrapper still requires Node.js to launch the installed command.

## Installation

```sh
npm install -g @algolia/cli
```

Or run without installing:

```sh
npx @algolia/cli --help
```

## Usage

```sh
algolia --help
algolia search --index my-index --query "foo"
algolia indices list
algolia apikeys list
```

## Documentation

Full documentation: [algolia.com/doc/tools/cli](https://algolia.com/doc/tools/cli/)

## Supported platforms

| OS      | x64 | arm64 |
|---------|-----|-------|
| macOS   | ✓   | ✓     |
| Linux   | ✓   | ✓     |
| Windows | ✓   | ✓     |

## Issues

[github.com/algolia/cli/issues](https://github.com/algolia/cli/issues)
