# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary name
BINARY_NAME=stay-for-long

# Server timeouts (in seconds)
READ_TIMEOUT=15
WRITE_TIMEOUT=15
IDLE_TIMEOUT=60

.PHONY: all build clean test test-coverage run lint

all: clean build

mod:
	@go mod tidy
	@go mod vendor

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

run:
	@echo "Starting application..."
	docker-compose up --build

lint:
	@echo "Running linter..."
	golangci-lint run --timeout=5m ./...

mock: ; $(info Generating mock files)
	@./generate-mocks.sh

# Help command
help:
	@echo "Available commands:"
	@echo "  make mod           - Update dependencies"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make run           - Start application with Docker"
	@echo "  make lint          - Run linter"
	@echo "  make mock          - Create new mocks"