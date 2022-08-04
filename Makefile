# Run all the tests
test:
	go test ./...
.PHONY: test

## Build & publish the documentation
docs:
	git clone https://github.com/algolia/doc.git "$@"

docs-bump: docs
	git -C docs pull
	git -C docs checkout feat/cli 
	git -C docs rm 'app_data/cli/commands/*.yml' 2>/dev/null || true
	go run ./cmd/docs --app_data-path docs/app_data/cli/commands
	git -C docs add 'app_data/cli/commands/*.yml'
	git -C docs commit -m 'update cli commands app data' || true
	git -C docs push
.PHONY: docs-bump

# Build the binary
build:
	go generate ./...
	go build -o algolia cmd/algolia/main.go
.PHONY: build
