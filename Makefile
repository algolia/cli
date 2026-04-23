ifdef VERSION
VERSION := $(VERSION)
else
VERSION := main
endif

# Run all the tests
test:
	go test ./... -p 1
.PHONY: test

## Build & publish the CLI documentation
DOCS_FOLDER ?= docs
DOCS_GENERATED_PATH = doc/tools/cli/commands
DOCS_REPO_URL = https://github.com/algolia/docs-new.git
DOCS_BRANCH = main

docs:
	git clone $(DOCS_REPO_URL) "$@"

.PHONY: docs-commands-data
docs-commands-data: docs
	git -C $(DOCS_FOLDER) checkout $(DOCS_BRANCH)
	git -C $(DOCS_FOLDER) pull --ff-only origin $(DOCS_BRANCH)
	rm -rf "$(DOCS_FOLDER)/$(DOCS_GENERATED_PATH)"
	mkdir -p "$(DOCS_FOLDER)/$(DOCS_GENERATED_PATH)"
	go run ./cmd/docs --app_data-path "$(DOCS_FOLDER)/$(DOCS_GENERATED_PATH)"
	git -C $(DOCS_FOLDER) add -A "$(DOCS_GENERATED_PATH)"

.PHONY: docs-pr
docs-pr: docs-commands-data
ifndef GITHUB_REF
	$(error GITHUB_REF is not set)
endif
	git -C $(DOCS_FOLDER) checkout -B feat/cli-'$(GITHUB_REF:refs/tags/v%=%)'
	git -C $(DOCS_FOLDER) commit -m 'feat: update cli commands data for $(GITHUB_REF:refs/tags/v%=%) version' || true
	git -C $(DOCS_FOLDER) push --set-upstream origin feat/cli-'$(GITHUB_REF:refs/tags/v%=%)'
	cd $(DOCS_FOLDER); gh pr create -f -b "Changelog: https://github.com/algolia/cli/releases/tag/$(GITHUB_REF:refs/tags/%=%)"

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
	go build -ldflags "\
		-s -w \
		-X=github.com/algolia/cli/pkg/version.Version=$(VERSION) \
		-X=github.com/algolia/cli/api/dashboard.DefaultDashboardURL=$(ALGOLIA_DASHBOARD_URL) \
		-X=github.com/algolia/cli/api/dashboard.DefaultAPIURL=$(ALGOLIA_API_URL) \
		-X=github.com/algolia/cli/pkg/config.DefaultSearchHosts=$(ALGOLIA_SEARCH_HOSTS) \
		-X=github.com/algolia/cli/pkg/auth.DefaultOAuthClientID=$(ALGOLIA_OAUTH_CLIENT_ID) \
		-X 'github.com/algolia/cli/api/dashboard.DefaultOAuthScope=$(ALGOLIA_OAUTH_SCOPE)'" \
		-o algolia cmd/algolia/main.go
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
