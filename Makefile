# Run all the tests
test:
	go test ./...
.PHONY: test

## Build & publish the documentation
docs:
	git clone https://github.com/algolia/cli-docs.git "$@"

docs-bump: docs
	git -C docs pull
	git -C docs rm 'algolia_*.md' 2>/dev/null || true
	go run ./cmd/docs --doc-path docs
	rm -f docs/*.bak
	git -C docs add 'algolia*.md'
	git -C docs commit -m 'update docs' || true
	git -C docs push
.PHONY: docs-bump

# Build the binary
build:
	go generate ./...
	go build -o algolia cmd/algolia/main.go
.PHONY: build
