.PHONY: test
test:
	go test ./...

## Documentation related tasks

docs:
	git clone https://github.com/algolia/cli-docs.git "$@"

.PHONY: docs-bump
docs-bump: docs
	git -C docs pull
	git -C docs rm 'algolia_*.md' 2>/dev/null || true
	go run ./cmd/docs --doc-path docs
	rm -f docs/*.bak
	git -C docs add 'algolia*.md'
	git -C docs commit -m 'update docs' || true
	git -C docs push
