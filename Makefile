.PHONY: all help build fmt vet clean run test

Default: help

help:
	@echo "Makefile per GoMCP"
	@echo ""
	@echo "Comandi disponibili:"
	@echo "  make all       - Format, vet e test"
	@echo "  make fmt       - Format del codice"
	@echo "  make vet       - Analisi statica del codice"
	@echo "  make test      - Esegue i test"
	
all: fmt vet test

fmt:
	@echo "→ Formatting code..."
	@go fmt ./...

vet:
	@echo "→ Running go vet..."
	@go vet ./...

test:
	@echo "→ Running tests..."
	@go test ./... -v --cover

clean:
	@echo "→ Cleaning build artifacts..."
	@rm -rf bin
