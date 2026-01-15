.PHONY: build clean test run demo help install

# Variables
BINARY_NAME=filecmp
BUILD_DIR=build
MAIN_FILE=main.go

# Default target
all: build

# Build the application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Build complete: ./$(BINARY_NAME)"

# Build with optimizations for release
build-release:
	@echo "🚀 Building $(BINARY_NAME) for release..."
	@go build -ldflags="-s -w" -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Release build complete: ./$(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "🌍 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	@echo "✅ Multi-platform builds complete in $(BUILD_DIR)/"
	@ls -la $(BUILD_DIR)/

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod tidy
	@go mod download
	@echo "✅ Dependencies installed"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test ./...
	@echo "✅ Tests complete"

# Run the test diff program
test-diff:
	@echo "🧪 Running diff engine test..."
	@go run test-diff.go

# Run the application with example files
run-files:
	@echo "🚀 Running with example files..."
	@./$(BINARY_NAME) examples/file1.txt examples/file2.txt

# Run the application with example directories
run-dirs:
	@echo "🚀 Running with example directories..."
	@./$(BINARY_NAME) examples/project-v1 examples/project-v2

# Run the application in interactive mode
run:
	@echo "🚀 Running in interactive mode..."
	@./$(BINARY_NAME)

# Show demo script
demo:
	@echo "🎬 Running demo script..."
	@./test-demo.sh

# Show merge functionality demo
demo-merge:
	@echo "🔀 Running merge functionality demo..."
	@./merge-example.sh

# Demo the new merge feature with simple files
demo-merge-simple:
	@echo "🔀 Running simple merge demo..."
	@./demo-merge.sh

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Install the binary to system PATH
install: build
	@echo "📲 Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installation complete. You can now run '$(BINARY_NAME)' from anywhere"

# Show help
help:
	@echo "📚 File Comparison TUI Tool - Makefile Help"
	@echo "==========================================="
	@echo ""
	@echo "Available targets:"
	@echo "  build         Build the application (default)"
	@echo "  build-release Build optimized release version"
	@echo "  build-all     Build for multiple platforms"
	@echo "  deps          Install Go dependencies"
	@echo "  test          Run tests"
	@echo "  test-diff     Run diff engine test"
	@echo "  run           Run in interactive mode"
	@echo "  run-files     Run with example files"
	@echo "  run-dirs      Run with example directories"
	@echo "  demo          Run demo script"
	@echo "  demo-merge    Run comprehensive merge demo"
	@echo "  demo-merge-simple Run simple merge demo"
	@echo "  clean         Clean build artifacts"
	@echo "  install       Install to system PATH"
	@echo "  help          Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build              # Build the application"
	@echo "  make demo               # Show demo with examples"
	@echo "  make demo-merge         # Show merge functionality demo"
	@echo "  make run-files          # Quick test with sample files"
	@echo "  make install            # Install system-wide"
	@echo ""
	@echo "Quick Start:"
	@echo "  1. make deps            # Install dependencies"
	@echo "  2. make build           # Build the application"
	@echo "  3. make demo            # See usage examples"
	@echo "  4. make run-files       # Try it out!"
