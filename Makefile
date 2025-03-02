# Variables
BINARY_NAME=task-tracking-service
COVERAGE_FILE=coverage.out
MAIN_PACKAGE=./cmd/server

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## test: Run all tests
test:
	go test -v ./...

## test-short: Run tests without integration tests
test-short:
	go test -v -short ./...

## test-coverage: Run tests with coverage
test-coverage:
	go test -v -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE)

## test-coverage-text: Run tests with coverage and display in terminal
test-coverage-text:
	go test -v -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE)

## test-watch: Run tests in watch mode (requires reflex)
test-watch:
	reflex -r '\.go$$' go test ./...

## test-clean: Remove test cache and coverage files
test-clean:
	go clean -testcache
	rm -f $(COVERAGE_FILE)

## run: Build and run the application
run:
	go run $(MAIN_PACKAGE)

## build: Build the application
build:
	go build -o $(GOBIN)/$(BINARY_NAME) $(MAIN_PACKAGE)

## clean: Clean up binary files
clean:
	go clean
	rm -f $(GOBIN)/$(BINARY_NAME)

## deps: Download dependencies
deps:
	go mod download
	go mod tidy

## install-tools: Install development tools
install-tools:
	go install github.com/cespare/reflex@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest

## help: Display this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: test test-short test-coverage test-coverage-text test-watch test-clean run build clean deps install-tools help migrate-up migrate-down migrate-create

# Database migration commands
migrate-up:
	go run cmd/migrate/main.go -command up

migrate-down:
	go run cmd/migrate/main.go -command down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/migrations -seq $$name