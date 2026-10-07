# VERSION is the git tag without its leading v, or the commit when there is none; the binary shows it in the badge.
VERSION ?= $(or $(shell git describe --tags --always 2>/dev/null | sed 's/^v//'),dev)

.PHONY: build test

build:
	go build -ldflags "-X goremi/internal/app.Version=$(VERSION)" -o bin/goremi ./cmd/goremi

test:
	go test ./...
	go vet ./...
	@test -z "$$(gofmt -l .)" || (gofmt -l .; exit 1)
