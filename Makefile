.PHONY: build test

build:
	go build -o bin/goremi ./cmd/goremi

test:
	go test ./...
	go vet ./...
	@test -z "$$(gofmt -l .)" || (gofmt -l .; exit 1)
