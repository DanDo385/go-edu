.PHONY: all build run test test-exercise clean help

all: build

build:
	@echo "Building microservice..."
	@go build -o bin/service ./cmd/service

run:
	@echo "Starting microservice..."
	@go run ./cmd/service

test:
	@echo "Running tests..."
	@go test -v ./...

test-exercise:
	@echo "Running exercise tests..."
	@go test -v ./exercise

test-race:
	@echo "Running tests with race detector..."
	@go test -race ./...

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	@go fmt ./...

help:
	@echo "Available targets:"
	@echo "  build         - Build the microservice binary"
	@echo "  run           - Run the microservice"
	@echo "  test          - Run all tests"
	@echo "  test-exercise - Run only exercise tests"
	@echo "  test-race     - Run tests with race detector"
	@echo "  clean         - Remove build artifacts"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  help          - Show this help message"
