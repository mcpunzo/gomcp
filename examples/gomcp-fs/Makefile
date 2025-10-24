BINARY_NAME = gomcp-filesystem

SRC_DIR = ./cmd

GO_FILES := $(shell find $(SRC_DIR) -name '*.go')

.PHONY: all help build fmt vet clean run test

Default: help

help:
	@echo "Makefile per $(BINARY_NAME)"
	@echo ""
	@echo "Comandi disponibili:"
	@echo "  make all       - Format, vet, test build"
	@echo "  make fmt       - source code Formatting"
	@echo "  make vet       - code static analysis"
	@echo "  make build     - code compilation"
	@echo "  make run       - execute the binary (after build)"
	@echo "  make test      - test"
	@echo "  make clean     - clean build artifacts"

all: fmt vet test build 

fmt:
	@echo "→ Formatting code..."
	@go fmt ./...

vet:
	@echo "→ Running go vet..."
	@go vet ./...

build:
	@echo "→ Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) $(SRC_DIR)

run: build
	@echo "→ Running $(BINARY_NAME)..."
	@./bin/$(BINARY_NAME)

test:
	@echo "→ Running tests..."
	@go test ./... -v --cover

clean:
	@echo "→ Cleaning build artifacts..."
	@rm -rf bin
