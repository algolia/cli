ifdef VERSION
VERSION := $(VERSION)
else
VERSION := main
endif

# Run all the tests
test:
	go test ./... -p 1
.PHONY: test

## Build & publish the documentation
docs:
	git clone https://github.com/algolia/doc.git "$@"

.PHONY: docs-commands-data
docs-commands-data: docs
	git -C docs pull
	git -C docs checkout master
	git -C docs rm 'app_data/cli/commands/*.yml' 2>/dev/null || true
	go run ./cmd/docs --app_data-path docs/app_data/cli/commands
	git -C docs add 'app_data/cli/commands/*.yml'

.PHONY: docs-pr
docs-pr: docs-commands-data
ifndef GITHUB_REF
	$(error GITHUB_REF is not set)
endif
	git -C docs checkout -B feat/cli-'$(GITHUB_REF:refs/tags/v%=%)'
	git -C docs commit -m 'feat: update cli commands data for $(GITHUB_REF:refs/tags/v%=%) version' || true
	git -C docs push --set-upstream origin feat/cli-'$(GITHUB_REF:refs/tags/v%=%)'
	cd docs; gh pr create -f -b "Changelog: https://github.com/algolia/cli/releases/tag/$(GITHUB_REF:refs/tags/v%=%)"

## Create a new PR (or update the existing one) to update the API specs
api-specs-pr:
	wget -O ./api/specs/search.yml https://raw.githubusercontent.com/algolia/api-clients-automation/main/specs/bundled/search.yml
	go generate ./...
	if [ -n "$$(git status --porcelain)" ]; then \
		git checkout -b feat/api-specs; \
		git add .; \
		git commit -m 'chore: update search api specs'; \
		git push -f --set-upstream origin feat/api-specs; \
		if ! [ "$$(gh pr list --base main --head feat/api-specs)" ]; then gh pr create --title "Update search api specs" --body "Update search api specs"; fi; \
	fi

# Build the binary
build:
	go generate ./...
	go build -ldflags "-s -w -X=github.com/algolia/cli/pkg/version.Version=$(VERSION)" -o algolia cmd/algolia/main.go
.PHONY: build

## Install & uninstall tasks are here for use on *nix platform only.
prefix  := /usr/local
bindir  := ${prefix}/bin

# Install Algolia CLI
install:
	make build
	install -m755 algolia ${bindir}
.PHONY: install

# Uninstall Algolia CLI
uninstall:
	rm ${bindir}/algolia
.PHONY: uninstall