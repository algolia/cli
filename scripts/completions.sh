#!/bin/sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run cmd/algolia/main.go completion "$sh" >"completions/algolia.$sh"
done