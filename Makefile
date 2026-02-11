.PHONY: fmt lint test all clean help

# Default target
all: fmt lint test

## fmt: Format the code
fmt:
	go fmt ./...

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## test: Run tests with race detection
test:
	go test -v -race -count=1 ./...

## clean: Clean test cache
clean:
	go clean -testcache

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed -e 's/^## //' | column -t -s ':'
