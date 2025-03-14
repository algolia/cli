# Algolia CLI

The Algolia CLI lets you work with your Algolia resources,
such as indices, records, API keys, and synonyms,
and from the command line.

![cli](https://user-images.githubusercontent.com/5702266/153008646-1fd8fbf2-4a4d-4421-b2f2-0886487f3e27.png)

## Documentation

See [Algolia CLI](https://algolia.com/doc/tools/cli/) in the Algolia documentation for setup and usage instructions.

## Installation

### macOS

The Algolia CLI is available on [Homebrew](https://brew.sh/) and as a downloadable binary from the [releases page](https://github.com/algolia/cli/releases).

```sh
brew install algolia/algolia-cli/algolia
```

### Linux

The Algolia CLI is available as a `.deb` package:

```sh
# Select the package appropriate for your platform:
sudo dpkg -i algolia_*.deb
```

as a `.rpm` package:

```sh
# Select the package appropriate for your platform
sudo rpm -i algolia_*.rpm
```

or as a tarball from the [releases page](https://github.com/algolia/cli/releases):

```sh
# Select the archive appropriate for your platform
tar xvf algolia_*_linux_*.tar.gz
```

### Windows

The Algolia CLI is available via [Chocolatey](https://community.chocolatey.org/packages/algolia/) and as a downloadable binary from the [releases page](https://github.com/algolia/cli/releases)

### Community packages

Other packages are maintained by the community, not by Algolia.
If you distribute a package for the Algolia CLI, create a pull request so that we can list it here!

### Build from source

To build the Algolia CLI from source, you'll need:

- Go version 1.23 or later
- [Go task](https://taskfile.dev/)

1. Clone the repo: `git clone https://github.com/kai687/cli.git algolia-cli && cd algolia-cli`
1. Run: `task build`

## Support

If you found an issue with the Algolia CLI,
[open a new GitHub issue](https://github.com/algolia/cli/issues/new),
or join the Algolia community on [Discord](https://alg.li/discord).
