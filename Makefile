.PHONY: fmt tidy lint test build coverage clean help

# Default target
all: fmt tidy lint test

## fmt: Format the code
fmt:
	gofmt -s -w .
	goimports -local github.com/nuln/dbase -w . || true

## tidy: Tidy go modules
tidy:
	go mod tidy

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## build: Build the project
build:
	go build ./...

## test: Run tests with race detection
test:
	go test -v -race -count=1 ./...

## coverage: Generate test coverage report
coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

## clean: Clean test cache and artifacts
clean:
	rm -f coverage.txt coverage.html
	go clean -testcache

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed -e 's/^## //' | column -t -s ':'
