# GitStuff Makefile
# 
# Quick start:
#   make build  - Build the application
#   make test   - Run all tests
#   make help   - Show all available commands

.PHONY: build test test-verbose clean help install

# Default target
help:
	@echo "Available commands:"
	@echo "  make build        - Build the gitstuff binary"
	@echo "  make test         - Run all tests"
	@echo "  make test-verbose - Run tests with verbose output"
	@echo "  make clean        - Remove built binaries"
	@echo "  make install      - Build and install to /usr/local/bin"
	@echo "  make help         - Show this help message"

# Build the application
build:
	@echo "Building gitstuff..."
	go build -o gitstuff .
	@echo "✅ Build complete: ./gitstuff"

# Run all tests
test:
	@echo "Running all tests..."
	go test ./cmd ./internal/config ./internal/git ./internal/gitlab
	@echo "✅ All tests passed!"

# Run tests with verbose output
test-verbose:
	@echo "Running all tests with verbose output..."
	go test -v ./cmd ./internal/config ./internal/git ./internal/gitlab

# Clean built binaries
clean:
	@echo "Cleaning up..."
	rm -f gitstuff
	@echo "✅ Cleanup complete"

# Build and install to system PATH
install: build
	@echo "Installing gitstuff to /usr/local/bin..."
	sudo cp gitstuff /usr/local/bin/
	@echo "✅ Installation complete! You can now run 'gitstuff' from anywhere"