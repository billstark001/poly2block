.PHONY: all build test clean install wasm help

# Default target
all: build

# Build CLI
build:
	@echo "Building CLI..."
	cd cmd/poly2block && go build -ldflags="-s -w" -o poly2block

# Build WASM
wasm:
	@echo "Building WASM..."
	cd wasm && GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o poly2block.wasm
	@echo "WASM binary size:"
	@ls -lh wasm/poly2block.wasm

# Run tests
test:
	@echo "Running tests..."
	cd core && go test -v -race -coverprofile=coverage.txt ./...

# Run tests with coverage report
coverage: test
	@echo "Generating coverage report..."
	cd core && go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report: core/coverage.html"

# Install CLI to GOPATH
install:
	@echo "Installing CLI..."
	cd cmd/poly2block && go install

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f cmd/poly2block/poly2block
	rm -f wasm/poly2block.wasm
	rm -f core/coverage.txt
	rm -f core/coverage.html
	rm -rf dist/

# Build all platforms for release
release:
	@echo "Building release binaries..."
	mkdir -p dist
	# Linux AMD64
	cd cmd/poly2block && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../../dist/poly2block-linux-amd64
	# Linux ARM64
	cd cmd/poly2block && GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ../../dist/poly2block-linux-arm64
	# macOS AMD64
	cd cmd/poly2block && GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ../../dist/poly2block-darwin-amd64
	# macOS ARM64
	cd cmd/poly2block && GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ../../dist/poly2block-darwin-arm64
	# Windows AMD64
	cd cmd/poly2block && GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ../../dist/poly2block-windows-amd64.exe
	# WASM
	cd wasm && GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o ../dist/poly2block.wasm
	@echo "Release binaries built in dist/"

# Generate vanilla palette
palette:
	@echo "Generating vanilla Minecraft palette..."
	cd cmd/poly2block && go run . generate-palette --output ../../dist/vanilla-palette.msgpack
	@echo "Palette saved to dist/vanilla-palette.msgpack"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	cd core && go mod tidy
	cd wasm && go mod tidy
	cd cmd/poly2block && go mod tidy
	go work sync

# Help
help:
	@echo "Available targets:"
	@echo "  all       - Build CLI (default)"
	@echo "  build     - Build CLI binary"
	@echo "  wasm      - Build WASM module"
	@echo "  test      - Run tests"
	@echo "  coverage  - Generate coverage report"
	@echo "  install   - Install CLI to GOPATH"
	@echo "  clean     - Remove build artifacts"
	@echo "  release   - Build all platforms for release"
	@echo "  palette   - Generate vanilla Minecraft palette"
	@echo "  fmt       - Format code"
	@echo "  lint      - Lint code"
	@echo "  tidy      - Tidy module dependencies"
	@echo "  help      - Show this help message"
